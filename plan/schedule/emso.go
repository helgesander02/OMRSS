// Package schedule — EMSO disaggregation pass.
//
// This file implements the Exploratory Message Scheduling Optimization
// strategy from Yan et al. 2024 (Journal of Systems Architecture)
// §5.3 Algorithm 2, adapted to the timeline-level data we have here.
//
// What the paper does
// -------------------
// EMSO sits *after* MAO has aggregated CAN frames into TSN frames.
// When the resulting TSN frame F_k cannot meet the three Yan §5.1
// constraints (Eq. 10 flow deadline / Eq. 11 send order / Eq. 12
// no-congestion), Yan's algorithm:
//
//   1. Sorts the CAN frames inside F_k by deadline ascending.
//   2. Pops the smallest-deadline CAN frames out of F_k until the
//      residual F_k can be scheduled. The residual gets a LOOSER
//      deadline because min(deadlines) inside it is now larger.
//   3. Re-packs the popped urgent frames into a fresh F'_k that
//      still carries the original tight deadline but is SMALLER, so
//      it usually finds a slot more easily.
//
// What we do here (and the gap)
// -----------------------------
// At this layer of the pipeline the CAN frames inside a TSN frame have
// already been collapsed into a single (period, deadline, datasize)
// triple — see network/flow/can/encapsulate.go::emitCAN2TTFrame and
// can.Flow.DataSize, which only stores the *aggregate* per-emit
// payload. The constituent CAN-frame deadline list is lost.
//
// Without that list we cannot do Yan's deadline-relaxation trick
// faithfully, so we approximate it: each disaggregation attempt
// splits the unschedulable flow into a *short, tight* part (mimicking
// the urgent CAN frames) and a *larger, relaxed* part (mimicking the
// residual TSN frame). We progressively shrink the urgent portion
// across attempts:
//
//   attempt=1:  urgent = datasize/2, residual = datasize/2 @ 2× deadline
//   attempt=2:  urgent = datasize/3, residual = 2·datasize/3 @ 3× deadline
//   attempt=3:  urgent = datasize/4, residual = 3·datasize/4 @ 4× deadline
//   ...
//
// Both parts are placed on the same route. The residual is allowed a
// looser deadline because in the real protocol the residual carries
// only the larger-deadline CAN frames; the urgent part keeps the
// original deadline. As soon as one (urgent, residual) split fits, we
// declare the flow EMSO-schedulable.
//
// This approximation is conservative on the easy side (a placement
// EMSO finds here would also be findable by real EMSO) but pessimistic
// on the hard side (real EMSO with per-CAN deadlines can sometimes
// place flows we still fail). The gap is documented so the OSRO paper
// can either (a) cite this as "EMSO-light" or (b) extend the encap
// pipeline to preserve constituent CAN deadlines and unlock the real
// algorithm. Either way the public API is stable and the harness can
// upgrade transparently.
//
// We do NOT modify per-CAN-frame composition in network/flow/can/* —
// EMSO operates purely on the timeline view of a single flow, leaving
// the encap output untouched. This keeps OSRO's "compare encap methods
// by their gateway-side metrics" claim independent from "compare
// methods by their post-routing schedulability" claim.
package schedule

import (
	"fmt"

	"src/network/flow/can"
	"src/plan/routes"
)

// EMSOMaxAttempts caps how many disaggregation refinements we try
// before declaring the flow truly unschedulable. Each attempt halves
// the urgent share and multiplies the residual's deadline by one more
// factor, so attempts > log2(datasize) stop helping and only waste
// time. Five is generous for the OSRO test scale (CAN payloads up to
// 1500 B).
const EMSOMaxAttempts = 5

// EMSODisaggregate is the Yan 2024 Algorithm 2 fallback at the
// timeline level. It is called from BuildTimelineForRoutes whenever a
// CAN2TT flow's direct placement fails, and tries up to maxAttempts
// (urgent, residual) splits before giving up.
//
// On success the timeline contains the split sub-flows under
// flowID + "-emso<attempt>-{urgent|residual}" suffixes, and EMSO
// returns true. On failure every reservation made during the
// attempt is rolled back via ReleaseFlowPrefix so the timeline is
// left in the same state as before the call. The caller can then
// count this flow as truly unschedulable (OBJ Round 3 does this).
func EMSODisaggregate(timeline *NetworkTimeline, flow *can.Flow, route *routes.Route, byteRate float64, flowID string, maxAttempts int) bool {
	if timeline == nil || flow == nil || route == nil {
		return false
	}
	if maxAttempts <= 0 {
		return false
	}

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		// urgent share shrinks each round; residual grows and gets a
		// proportionally looser deadline (the "removed-the-tight-CAN-
		// frames" approximation).
		urgentShare := 1.0 / float64(attempt+1)
		residualShare := 1.0 - urgentShare
		residualDeadlineMul := attempt + 1

		urgentSize := flow.DataSize * urgentShare
		residualSize := flow.DataSize * residualShare

		// Once the urgent fragment shrinks below 1 byte we are not
		// modelling a meaningful CAN frame any more; stop here.
		if urgentSize < 1.0 {
			break
		}

		urgent := *flow
		urgent.DataSize = urgentSize
		// urgent keeps the original tight deadline.

		residual := *flow
		residual.DataSize = residualSize
		residual.Deadline = flow.Deadline * residualDeadlineMul
		// Clamp relaxed deadline to one hyperperiod so the placement
		// search horizon never exceeds the timeline's coordinate space.
		if residual.Deadline > timeline.HyperPeriod {
			residual.Deadline = timeline.HyperPeriod
		}

		urgentID := fmt.Sprintf("%s-emso%d-urgent", flowID, attempt)
		residualID := fmt.Sprintf("%s-emso%d-residual", flowID, attempt)

		okU := timeline.RegisterCAN2TTFlow(&urgent, route, byteRate, urgentID)
		okR := timeline.RegisterCAN2TTFlow(&residual, route, byteRate, residualID)

		if okU && okR {
			// Disaggregation succeeded. Both parts are now part of the
			// timeline; subsequent flows / AVB WCD will see their
			// busy windows.
			return true
		}

		// Roll back partial placement and try a finer split.
		timeline.ReleaseFlowPrefix(urgentID)
		timeline.ReleaseFlowPrefix(residualID)
	}
	return false
}
