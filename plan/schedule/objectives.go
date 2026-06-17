package schedule

import (
	"fmt"
	"src/network"
	"src/network/flow/tt"
	"src/pkg/config"
	"src/pkg/logger"
	"src/plan/routes"

	"time"
)

// Objectives. costSetting is read from cfg.GetCostArray() so callers no
// longer thread the same array down through every function level.
//
// As of the Laursen / EMSO / Round 3 commit, OBJ pre-builds a
// NetworkTimeline that contains every native TT instance *and* every
// CAN2TT instance produced by the OSRO encap pipeline. AVB WCD is then
// computed via Laursen 2016 §6.3 Algorithm 1 against that combined
// timeline — so the AVB WCD term (O4) and AVB schedulability (O2) now
// respond to encap-method choices instead of being agnostic to them.
//
// `methodScope` selects which encap method's CAN2TT flows participate
// in this evaluation. An empty string evaluates every method together
// (legacy "mix all" behaviour, used by OMACO which has no encap). A
// non-empty value scopes Round 3 (timeline + schedulability count) to
// flows produced by that single method. OSRO's per-method harness
// passes a non-empty value so each method's effect on native AVB is
// measured in isolation.
//
// CAN2TT schedulability is rolled into O1 (TT-class failures) on the
// pragmatic grounds that downstream callers expect a fixed
// [4]float64 obj layout. A per-class breakdown is logged when
// showflow=true.
func OBJ(network *network.Network, cfg *config.Config, X *routes.KRouteSet, II *routes.RouteSet, II_prime *routes.RouteSet, methodScope string, showflow bool) ([4]float64, int) {
	costSetting := cfg.GetCostArray()
	S := network.FlowSet.InputOMACOFlowSet()
	S_prime := network.FlowSet.BGOMACOFlowSet()
	var (
		obj                 [4]float64
		cost                int
		tsn_failed_count    int           = 0 // O1 (native TT)
		avb_failed_count    int           = 0 // O2
		can2tt_failed_count int           = 0 // folded into O1
		all_rerouted_count  int           = 0 // O3 ... pass
		avb_wcd_sum         time.Duration     // O4
	)
	linkmap := map[string]float64{}

	// Build the timeline + AVB per-link load once for this evaluation.
	// Native TT (Rounds 1+2) and CAN2TT (Round 3 encap output) are both
	// placed on the timeline; selected AVB routes (BG + input) are
	// tallied into avbLoad. WCDWithTimeline reads both so encap-method
	// choices propagate into AVB O2/O4 *and* AVB-on-AVB interference is
	// counted exactly once per actually-selected route.
	timeline := BuildTimelineForRoutes(network, cfg, II, II_prime, methodScope)
	avbLoad := ComputeAVBPerLinkLoad(II, II_prime, network.FlowSet)

	// Round1: Schedule BG flow
	// O1
	for nth, route := range II_prime.TSNRoutes {
		schedulability := schedulability(0, S_prime.TSNFlows[nth], route, linkmap, cfg.Network.Bandwidth, cfg.Network.Hyperperiod)
		tsn_failed_count += 1 - schedulability
	}

	// O2 and O4 — AVB WCD via Laursen
	for nth, route := range II_prime.AVBRoutes {
		wcd := WCDWithTimeline(route, X, S_prime.AVBFlows[nth], network.FlowSet, timeline, avbLoad)
		avb_wcd_sum += wcd
		schedulability := schedulability(wcd, S_prime.AVBFlows[nth], route, linkmap, cfg.Network.Bandwidth, cfg.Network.Hyperperiod)
		avb_failed_count += 1 - schedulability
	}
	// O3 ... pass

	// Round2: Schedule Input flow
	// O1
	for nth, route := range II.TSNRoutes {
		schedulability := schedulability(0, S.TSNFlows[nth], route, linkmap, cfg.Network.Bandwidth, cfg.Network.Hyperperiod)
		tsn_failed_count += 1 - schedulability
	}

	// O2 and O4 — AVB WCD via Laursen
	for nth, route := range II.AVBRoutes {
		wcd := WCDWithTimeline(route, X, S.AVBFlows[nth], network.FlowSet, timeline, avbLoad)
		avb_wcd_sum += wcd
		schedulability := schedulability(wcd, S.AVBFlows[nth], route, linkmap, cfg.Network.Bandwidth, cfg.Network.Hyperperiod)
		avb_failed_count += 1 - schedulability
	}
	// O3 ... pass

	// Round3: Schedule CAN2TT flows (OSRO-only). For each gateway-emit
	// CAN2TT flow we run the same bandwidth schedulability check used
	// for native TT, *plus* we cross-check against the timeline:
	// flows that survived BuildTimelineForRoutes (direct placement or
	// EMSO fallback) are considered timeline-schedulable; flows in
	// UnplaceableFlows lost both their direct slot and their EMSO
	// retries, so we count them here.
	can2tt_failed_count = countCAN2TTFailures(network, cfg, II, linkmap, timeline, methodScope)
	tsn_failed_count += can2tt_failed_count

	obj[0] = float64(tsn_failed_count)                          // O1
	obj[1] = float64(avb_failed_count)                          // O2
	obj[2] = float64(all_rerouted_count)                        // O3 ... pass
	obj[3] = float64(avb_wcd_sum) / float64(time.Millisecond)   // O4 — milliseconds (matches ETFA 2020 Table III convention)

	cost += tsn_failed_count * costSetting[0] // O1
	cost += avb_failed_count * costSetting[1] // O2
	// Cost stays in µs internally for OSACO ranking precision —
	// converting to ms here would round to zero on the small cases
	// where ranking matters most.
	cost += int(avb_wcd_sum/time.Microsecond) * costSetting[3] // O4

	if showflow {
		logger.Println(linkmap)
		logger.Printf(
			"timeline: links=%d reserved_us=%d unplaceable=%d can2tt_fail=%d native_tt_fail=%d\n",
			len(timeline.Links),
			timeline.TotalReservedUs(),
			len(timeline.UnplaceableFlows),
			can2tt_failed_count,
			tsn_failed_count-can2tt_failed_count,
		)
	}

	return obj, cost
}

// countCAN2TTFailures runs Round 3 of OBJ: for each CAN2TT flow, it
// verifies that
//  1. the bandwidth-based schedulability check still passes (same
//     model used for native TT), and
//  2. the timeline placement succeeded (either directly or via EMSO
//     disaggregation; failure here means even EMSO couldn't fit it).
//
// Either condition failing increments the per-flow counter. When
// `methodScope` is non-empty, only flows produced by that method
// participate in both checks — this is required for OSRO's per-method
// comparison, otherwise schedulability for method A would be polluted
// by methods B…F flows that A would never actually emit.
//
// The flow ordering still matches routes.Get_KPath_Routing's nested
// (method, j) order, and II.CAN2TTRoutes is assumed to have been
// filtered to the same method scope upstream
// (FilterCAN2TTByMethod) — we keep routeIdx walking only when method
// matches scope so the route-to-flow pairing stays aligned.
func countCAN2TTFailures(network *network.Network, cfg *config.Config, II *routes.RouteSet, linkmap map[string]float64, timeline *NetworkTimeline, methodScope string) int {
	if network.FlowSet.EncapsulateMethod == nil || len(II.CAN2TTRoutes) == 0 {
		return 0
	}

	// Map flowID prefixes to UnplaceableFlows for O(1) lookup.
	unplaceable := make(map[string]bool)
	for _, id := range timeline.UnplaceableFlows {
		// strip the per-instance "#k" suffix to get the flow-level key.
		if i := indexByte(id, '#'); i >= 0 {
			unplaceable[id[:i]] = true
		} else {
			unplaceable[id] = true
		}
	}

	failed := 0
	routeIdx := 0
	for _, method := range network.FlowSet.EncapsulateMethod {
		if methodScope != "" && method.MethodName != methodScope {
			// Caller already trimmed II.CAN2TTRoutes to methodScope;
			// keep routeIdx still while we walk over non-matching methods.
			continue
		}
		for j, canFlow := range method.CAN2TTFlows {
			if routeIdx >= len(II.CAN2TTRoutes) {
				return failed
			}
			route := II.CAN2TTRoutes[routeIdx]
			routeIdx++
			if route == nil {
				failed++
				continue
			}

			// (1) Bandwidth check — reuse the native-TT helper by
			// passing the CAN2TT flow's fields through an inline
			// tt.Flow adapter. The adapter is purely structural; both
			// flow types share Source/Destinations/DataSize/Period.
			adapter := &tt.Flow{
				Period:       canFlow.Period,
				Deadline:     canFlow.Deadline,
				DataSize:     canFlow.DataSize,
				HyperPeriod:  canFlow.HyperPeriod,
				Source:       canFlow.Source,
				Destinations: canFlow.Destinations,
			}
			bwOk := schedulability(0, adapter, route, linkmap, cfg.Network.Bandwidth, cfg.Network.Hyperperiod) == 1

			// (2) Timeline check — did this flow's direct placement
			// (or EMSO fallback) leave any unplaceable instance behind?
			tlKey := fmt.Sprintf("can2tt-%s-%d", method.MethodName, j)
			tlOk := !unplaceable[tlKey]

			if !bwOk || !tlOk {
				failed++
			}
		}
	}
	return failed
}

// indexByte is an inline strings.IndexByte to avoid pulling the package
// just for this one call site.
func indexByte(s string, c byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return i
		}
	}
	return -1
}

// schedulability is the bandwidth-only schedulability check shared by
// every flow class (native TT, AVB, CAN2TT). For AVB, `wcd` is the
// Laursen 2016 worst-case end-to-end delay; for TT/CAN2TT, 0 is passed
// because no WCD analysis is done — the deadline is enforced solely by
// the per-link byte budget.
//
// Previously this routine mutated `linkmap` in place as it walked, and
// only rolled back the LAST link's increment on failure. Upstream
// hops' additions stayed in `linkmap`, polluting subsequent flows and
// causing cascading false failures whenever a multi-hop flow ran out
// of room downstream. The fix below collects every per-link byte
// addition in a separate delta map and commits atomically only if the
// whole walk plus the WCD constraint succeed.
func schedulability(wcd time.Duration, flow *tt.Flow, route *routes.Route, linkmap map[string]float64, bandwidth float64, hyperPeriod int) int {
	r := wcd <= time.Duration(flow.Deadline)*time.Microsecond
	delta := make(map[string]float64)
	node := route.GetNodeByID(flow.Source)
	ok := schedulableDelta(node, -1, flow, route, linkmap, delta, bandwidth, hyperPeriod)
	if !(r && ok) {
		// Discard delta — leave linkmap exactly as we found it.
		return 0
	}
	// Atomic commit.
	for k, v := range delta {
		linkmap[k] += v
	}
	return 1
}

// schedulableDelta walks the route DFS-style and accumulates per-link
// byte additions in `delta`. `linkmap` is read-only — we check
// `linkmap[key] + delta[key]` against the per-link bandwidth budget so
// the combined load (already-committed + this flow's pending) is what
// is enforced. End-station-incident links are skipped, matching the
// duplex assumption that talkers and listeners have dedicated cables.
func schedulableDelta(node *routes.Node, parentID int, flow *tt.Flow, route *routes.Route, linkmap, delta map[string]float64, bandwidth float64, hyperPeriod int) bool {
	if node == nil {
		return true
	}
	for _, link := range node.Connections {
		if link.ToNodeID == parentID {
			continue
		}
		if !(link.FromNodeID == flow.Source || loopcompare(link.ToNodeID, flow.Destinations)) {
			key := fmt.Sprintf("%d>%d", link.FromNodeID, link.ToNodeID)
			if flow.Period <= 0 {
				return false
			}
			add := flow.DataSize * float64(hyperPeriod/flow.Period)
			delta[key] += add
			if linkmap[key]+delta[key] > bandwidth {
				return false
			}
		}
		nextnode := route.GetNodeByID(link.ToNodeID)
		if !schedulableDelta(nextnode, node.ID, flow, route, linkmap, delta, bandwidth, hyperPeriod) {
			return false
		}
	}
	return true
}

func loopcompare(a int, b []int) bool {
	for _, v := range b {
		if a == v {
			return true
		}
	}
	return false
}
