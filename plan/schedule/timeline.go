// Package schedule — link-timeline / minimal GCL infrastructure.
//
// This file gives the rest of the OSRO pipeline a way to *reason about
// time* on each TSN link, rather than only about average bandwidth (which
// is what objectives.go does today via linkmap[key] += bytes).
//
// What this file is, in paper terms:
//
//   - LinkTimeline is the per-egress-port "busy interval list" inside
//     one hyperperiod. It is the minimal data structure required to
//     enforce Yan et al. 2024 (JSA) §5.1 *Eq. (12) no-congestion
//     constraint*: two TT/TSN frames may not occupy the same link at
//     the same time. Yan's paper does not name the structure GCL —
//     they call it "send offset reservation" — but the semantics are
//     identical to the busy-window list a TSN switch needs to derive
//     an 802.1Qbv Gate Control List.
//
//   - NetworkTimeline aggregates one LinkTimeline per directed link
//     so every routing scorer / EMSO pass / WCD analyser can share
//     a single view of the network's time domain.
//
//   - RegisterTTFlow walks a Route and reserves one window per (hop,
//     instance) following the no-wait policy from Yan §3.4: each
//     downstream hop starts the moment its upstream hop finishes,
//     and the deadline horizon bounds the search exactly as Yan §5.1
//     Eq. (10). When a hop has no free slot inside that horizon we
//     roll back this instance — leaving the timeline atomically
//     consistent so EMSO can probe alternative offsets without
//     fighting the leftover state.
//
// What this file is NOT (yet):
//
//   - It is not a full 802.1Qbv GCL. We do not model "queue 7 open,
//     queue 6 closed" gate states; we only model "this link is busy
//     transmitting at time T". Laursen 2016 §6.3 Algorithm 1 needs
//     the busy windows to compute the worst-case AVB gate-close
//     interval — BusyOverlap() exposes them in the shape that
//     algorithm expects, so the next commit (Laursen-style WCD)
//     can read from here without changing this file.
//
//   - It does not yet enforce Yan §5.1 Eq. (11) *send order* across
//     a route. placeInstance walks the route depth-first and chains
//     hop k+1's preferred start to hop k's End, which gives the
//     no-wait order Yan assumes, but if EMSO later wants to allow
//     queued waiting (Q_F_k > 0), this is the place to relax it.
//
//   - 802.1Qbu preemption is not modelled. Per-window granularity is
//     1 µs (perHop is ceil()'d). On the topologies we run (1 Gbps,
//     ≤1500 B per frame) the rounding adds at most a fraction of a
//     µs of pessimism, which is the conservative direction.
package schedule

import (
	"fmt"
	"math"
	"sort"

	"src/network"
	"src/network/flow/can"
	"src/network/flow/tt"
	"src/pkg/config"
	"src/plan/routes"
)

// -----------------------------------------------------------------------
// Data structures
// -----------------------------------------------------------------------

// Window is one occupied transmission slot on a directed link. All
// timestamps are integer microseconds measured from the start of the
// hyperperiod. [Start, End) is half-open so adjacent windows that share
// an endpoint do not count as overlapping.
type Window struct {
	Start  int    // µs, inclusive
	End    int    // µs, exclusive
	FlowID string // owner; used by ReleaseFlow for EMSO rollback
}

// LinkTimeline is one directed link's busy-window list inside a single
// hyperperiod. Windows are kept sorted by Start ascending so overlap
// probes and inserts are O(log n + occupied window count).
type LinkTimeline struct {
	Key         string // "from>to" — matches objectives.linkmap convention
	HyperPeriod int
	Windows     []Window
}

// NetworkTimeline owns one LinkTimeline per directed link and records
// every flow / instance ID that could not be placed inside its deadline
// horizon. UnplaceableFlows is the diagnostic channel the upcoming OSRO
// experiment harness will read from to compute method-level drop counts.
type NetworkTimeline struct {
	HyperPeriod      int
	Links            map[string]*LinkTimeline
	UnplaceableFlows []string
}

// NewNetworkTimeline allocates an empty timeline bound to one hyperperiod.
func NewNetworkTimeline(hyperPeriod int) *NetworkTimeline {
	return &NetworkTimeline{
		HyperPeriod: hyperPeriod,
		Links:       make(map[string]*LinkTimeline),
	}
}

// LinkKey is the canonical "from>to" identifier shared with
// objectives.linkmap so both representations agree on link identity.
func LinkKey(from, to int) string {
	return fmt.Sprintf("%d>%d", from, to)
}

// EnsureLink lazily creates and returns the LinkTimeline for (from, to).
func (nt *NetworkTimeline) EnsureLink(from, to int) *LinkTimeline {
	key := LinkKey(from, to)
	if lt, ok := nt.Links[key]; ok {
		return lt
	}
	lt := &LinkTimeline{
		Key:         key,
		HyperPeriod: nt.HyperPeriod,
	}
	nt.Links[key] = lt
	return lt
}

// -----------------------------------------------------------------------
// Overlap / reserve primitives
// -----------------------------------------------------------------------

// CheckOverlap reports whether [start, start+duration) collides with any
// already-reserved window on link (from, to). Half-open semantics: a
// window ending exactly at `start` does not collide. Cyclic wrap past
// HyperPeriod is *not* handled here; callers that want wrap should
// pre-modulo and probe twice.
func (nt *NetworkTimeline) CheckOverlap(from, to, start, duration int) bool {
	if duration <= 0 {
		return false
	}
	lt, ok := nt.Links[LinkKey(from, to)]
	if !ok {
		return false
	}
	end := start + duration

	// First window whose Start >= start.
	idx := sort.Search(len(lt.Windows), func(i int) bool {
		return lt.Windows[i].Start >= start
	})

	// The window immediately before idx may straddle `start`.
	if idx > 0 && lt.Windows[idx-1].End > start {
		return true
	}
	// The window at idx may begin before our `end`.
	if idx < len(lt.Windows) && lt.Windows[idx].Start < end {
		return true
	}
	return false
}

// Reserve inserts [start, start+duration) on (from, to) only if it does
// not collide with any existing window. Returns true on success. Each
// call adds at most one window — repeated calls with the same flowID on
// the same link append a new window per call (one per flow instance).
func (nt *NetworkTimeline) Reserve(from, to, start, duration int, flowID string) bool {
	if duration <= 0 {
		return true
	}
	if start < 0 || start+duration > nt.HyperPeriod {
		return false
	}
	if nt.CheckOverlap(from, to, start, duration) {
		return false
	}
	lt := nt.EnsureLink(from, to)
	w := Window{Start: start, End: start + duration, FlowID: flowID}
	idx := sort.Search(len(lt.Windows), func(i int) bool {
		return lt.Windows[i].Start >= start
	})
	lt.Windows = append(lt.Windows, Window{}) // grow by one
	copy(lt.Windows[idx+1:], lt.Windows[idx:])
	lt.Windows[idx] = w
	return true
}

// FindFreeOffset returns the smallest offset in [preferredStart, maxStart]
// at which `duration` µs fits on (from, to) without overlap. Returns -1
// if no such offset exists or if duration would exceed the hyperperiod.
//
// This is the building block EMSO will call when it relaxes Eq. (10)
// to widen the deadline horizon (Yan 2024 Algorithm 2 line 12,
// "O_F_k ← O_F_k + 1" — we do the same probe but skip past known busy
// windows rather than incrementing by 1 µs at a time).
func (nt *NetworkTimeline) FindFreeOffset(from, to, preferredStart, maxStart, duration int) int {
	if duration <= 0 {
		return preferredStart
	}
	candidate := preferredStart
	if candidate < 0 {
		candidate = 0
	}
	if candidate+duration > nt.HyperPeriod {
		return -1
	}

	lt, ok := nt.Links[LinkKey(from, to)]
	if !ok {
		if candidate <= maxStart {
			return candidate
		}
		return -1
	}

	for {
		if candidate > maxStart || candidate+duration > nt.HyperPeriod {
			return -1
		}
		if !nt.CheckOverlap(from, to, candidate, duration) {
			return candidate
		}
		// Jump past the window that is currently in the way.
		idx := sort.Search(len(lt.Windows), func(i int) bool {
			return lt.Windows[i].Start >= candidate
		})
		if idx > 0 && lt.Windows[idx-1].End > candidate {
			candidate = lt.Windows[idx-1].End
			continue
		}
		if idx < len(lt.Windows) {
			candidate = lt.Windows[idx].End
			continue
		}
		return -1
	}
}

// -----------------------------------------------------------------------
// Release / query helpers
// -----------------------------------------------------------------------

// ReleaseFlow removes every window whose FlowID matches `flowID` from
// every link. EMSO calls this when it gives up on a TSN frame's current
// offset and wants to re-place it from scratch.
func (nt *NetworkTimeline) ReleaseFlow(flowID string) {
	for _, lt := range nt.Links {
		out := lt.Windows[:0]
		for _, w := range lt.Windows {
			if w.FlowID != flowID {
				out = append(out, w)
			}
		}
		lt.Windows = out
	}
}

// ReleaseFlowPrefix removes every window whose FlowID begins with the
// given prefix. Useful when a flow uses per-instance suffixes (e.g.
// "tt-input-3#0", "tt-input-3#1", ...) and EMSO wants to roll the whole
// flow back at once.
func (nt *NetworkTimeline) ReleaseFlowPrefix(prefix string) {
	for _, lt := range nt.Links {
		out := lt.Windows[:0]
		for _, w := range lt.Windows {
			if !hasPrefix(w.FlowID, prefix) {
				out = append(out, w)
			}
		}
		lt.Windows = out
	}
}

// hasPrefix is a tiny inline reimplementation so we don't pull strings.
func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

// LaursenInterference computes the worst-case TT interference µs that
// an AVB frame of payload `twcUs` µs would experience on link (from, to),
// following Laursen 2016 §6.3 Algorithm 1 verbatim:
//
//   - Try anchoring the AVB transmission start at the Start of every
//     existing TT window (those are the worst-case arrival times — AVB
//     arrives just as a gate is closing).
//   - From each anchor, walk subsequent TT windows: gaps between them
//     are AVB free-transmit time (decrement remaining payload), the
//     windows themselves are AVB stall time (accumulate into icurrent).
//   - Stop the walk for an anchor when remaining payload hits zero,
//     i.e. the AVB frame has had enough free µs to finish.
//   - Return max icurrent across all anchors.
//
// This is strictly stronger than a flat BusyOverlap([start, start+twc])
// query: when TT bursts extend past start+twc but still preempt AVB,
// flat overlap understates the interference; the recursive walk
// captures the full inflated transmission window.
func (nt *NetworkTimeline) LaursenInterference(from, to, twcUs int) int {
	if twcUs <= 0 {
		return 0
	}
	lt, ok := nt.Links[LinkKey(from, to)]
	if !ok {
		return 0
	}
	if len(lt.Windows) == 0 {
		return 0
	}

	imax := 0
	for i := range lt.Windows {
		ic := walkLaursenFromAnchor(lt.Windows, i, lt.Windows[i].Start, twcUs)
		if ic > imax {
			imax = ic
		}
	}
	return imax
}

// walkLaursenFromAnchor walks the per-link window list from the
// startIdx-th window forward, modelling AVB transmission starting at
// `anchor` with `twcUs` of payload still to send. Returns total µs of
// TT-induced stall (= Laursen's icurrent for this anchor).
func walkLaursenFromAnchor(windows []Window, startIdx, anchor, twcUs int) int {
	icurrent := 0
	remaining := twcUs
	pos := anchor

	for i := startIdx; i < len(windows); i++ {
		w := windows[i]
		if w.End <= pos {
			continue
		}
		// Free gap before this window: AVB transmits, consume remaining.
		if w.Start > pos {
			gap := w.Start - pos
			if gap >= remaining {
				return icurrent
			}
			remaining -= gap
			pos = w.Start
		}
		// Inside the window: AVB stalled, gate closed.
		stallStart := pos
		if stallStart < w.Start {
			stallStart = w.Start
		}
		icurrent += w.End - stallStart
		pos = w.End
	}
	// No more windows but still have remaining payload — that final
	// stretch is free, no additional interference.
	return icurrent
}

// BusyOverlap returns the total µs that [start, end) intersects busy
// windows on (from, to). Cheaper but weaker than LaursenInterference:
// it ignores TT windows that overlap the AVB transmission inflated by
// preemption. Use BusyOverlap for diagnostics / quick checks, and
// LaursenInterference for the real WCD bound.
func (nt *NetworkTimeline) BusyOverlap(from, to, start, end int) int {
	if end <= start {
		return 0
	}
	lt, ok := nt.Links[LinkKey(from, to)]
	if !ok {
		return 0
	}
	total := 0
	// Skip past windows that end before the query starts.
	startIdx := sort.Search(len(lt.Windows), func(i int) bool {
		return lt.Windows[i].End > start
	})
	for i := startIdx; i < len(lt.Windows); i++ {
		w := lt.Windows[i]
		if w.Start >= end {
			break
		}
		s := w.Start
		if s < start {
			s = start
		}
		e := w.End
		if e > end {
			e = end
		}
		if e > s {
			total += e - s
		}
	}
	return total
}

// LinkLoadUs is a diagnostic: total reserved µs on (from, to). Not part
// of the scheduler; only useful for inspecting how full a link is.
func (nt *NetworkTimeline) LinkLoadUs(from, to int) int {
	lt, ok := nt.Links[LinkKey(from, to)]
	if !ok {
		return 0
	}
	total := 0
	for _, w := range lt.Windows {
		total += w.End - w.Start
	}
	return total
}

// TotalReservedUs is a diagnostic: summed busy µs across every link in
// the network. Useful to compare methods at the same load level.
func (nt *NetworkTimeline) TotalReservedUs() int {
	total := 0
	for _, lt := range nt.Links {
		for _, w := range lt.Windows {
			total += w.End - w.Start
		}
	}
	return total
}

// -----------------------------------------------------------------------
// Flow placement
// -----------------------------------------------------------------------

// RegisterTTFlow places a native TSN TT flow on the timeline.
// See registerFlowGeneric for the placement semantics.
func (nt *NetworkTimeline) RegisterTTFlow(flow *tt.Flow, route *routes.Route, byteRate float64, flowID string) bool {
	if flow == nil {
		return false
	}
	return nt.registerFlowGeneric(
		flow.Period, flow.Deadline, flow.DataSize, flow.Source,
		route, byteRate, flowID,
	)
}

// RegisterCAN2TTFlow places a CAN-derived TSN flow on the timeline.
// can.Flow and tt.Flow share the same scheduling-relevant shape
// (Period, Deadline, DataSize, Source) so the placement reduces to the
// same generic helper as native TT. This is what makes OSRO Round 3
// work — once a CAN2TT flow is on the timeline it interferes with
// AVB exactly the same way a native TT flow does.
func (nt *NetworkTimeline) RegisterCAN2TTFlow(flow *can.Flow, route *routes.Route, byteRate float64, flowID string) bool {
	if flow == nil {
		return false
	}
	return nt.registerFlowGeneric(
		flow.Period, flow.Deadline, flow.DataSize, flow.Source,
		route, byteRate, flowID,
	)
}

// registerFlowGeneric is the underlying placement routine used by both
// RegisterTTFlow and RegisterCAN2TTFlow. It assumes Yan 2024 §3.4
// no-wait dispatch: hop k+1's reservation starts the moment hop k's
// reservation ends. Each instance is placed atomically — if any hop
// within an instance has no slot inside its deadline horizon (Eq. 10),
// every reservation made for that instance is rolled back so EMSO can
// probe again from a clean state.
//
// Returns true only if every instance of every hop succeeded. Partial
// successes (some instances placed, others not) still return false but
// leave the placed instances reserved so other flows can see them.
//
// `byteRate` is microseconds per byte (cfg.Network.ByteRate). The
// per-hop transmission time is ceil(DataSize × byteRate) µs — the
// ceiling is the conservative direction so the timeline never under-
// reserves what the underlying link will actually spend transmitting.
func (nt *NetworkTimeline) registerFlowGeneric(period, deadline int, dataSize float64, source int, route *routes.Route, byteRate float64, flowID string) bool {
	if route == nil || period <= 0 {
		return false
	}
	perHop := int(math.Ceil(dataSize * byteRate))
	if perHop <= 0 {
		perHop = 1
	}
	if period > nt.HyperPeriod {
		return false
	}
	instances := nt.HyperPeriod / period

	src := route.GetNodeByID(source)
	if src == nil {
		nt.UnplaceableFlows = append(nt.UnplaceableFlows, flowID+"#nosrc")
		return false
	}

	allPlaced := true
	for k := 0; k < instances; k++ {
		baseStart := k * period
		instID := fmt.Sprintf("%s#%d", flowID, k)
		if !nt.placeInstance(src, -1, baseStart, perHop, instID, route, deadline) {
			// Atomic rollback for this instance only — other instances
			// of this flow keep their reservations.
			nt.ReleaseFlow(instID)
			nt.UnplaceableFlows = append(nt.UnplaceableFlows, instID)
			allPlaced = false
		}
	}
	return allPlaced
}

// placeInstance walks the route DFS-style, mirroring the schedulability
// and WCD walks in objectives.go and wcd.go. For each link out of `node`
// it picks the smallest legal offset ≥ currentStart whose hop fits
// inside [baseStart, baseStart+deadline]. On success it recurses into
// the downstream subtree with the chosen End as the next currentStart.
// On any failure the caller (RegisterTTFlow) calls ReleaseFlow to undo
// this instance's partial reservations.
func (nt *NetworkTimeline) placeInstance(node *routes.Node, parentID, currentStart, perHop int, flowID string, route *routes.Route, deadline int) bool {
	for _, link := range node.Connections {
		if link.ToNodeID == parentID {
			continue
		}
		maxStart := currentStart + deadline - perHop
		if maxStart < currentStart {
			return false
		}
		slot := nt.FindFreeOffset(link.FromNodeID, link.ToNodeID, currentStart, maxStart, perHop)
		if slot < 0 {
			return false
		}
		if !nt.Reserve(link.FromNodeID, link.ToNodeID, slot, perHop, flowID) {
			return false
		}
		nextNode := route.GetNodeByID(link.ToNodeID)
		if nextNode == nil {
			continue // leaf / end device — nothing further to walk
		}
		if !nt.placeInstance(nextNode, node.ID, slot+perHop, perHop, flowID, route, deadline) {
			return false
		}
	}
	return true
}

// -----------------------------------------------------------------------
// Population helper — wire OBJ's flow / route ordering into the timeline
// -----------------------------------------------------------------------

// BuildTimelineForRoutes constructs a NetworkTimeline populated with
// every native TT flow from both background and input route sets, plus
// every CAN-derived TT flow produced by the OSRO encapsulation
// pipeline. AVB flows are intentionally not placed — they are
// credit-shaped (CBS), not gate-scheduled, so they do not produce GCL
// windows. AVB WCD logic reads this timeline via LaursenInterference
// to discover the gate-close intervals that TT (native + CAN2TT)
// creates.
//
// `methodScope` controls which encap method's CAN2TT flows participate.
// When empty, every method's flows are placed together (legacy "mix all"
// behaviour). When set to a method name (e.g. "mao"), only that method's
// CAN2TT flows are placed — the others are skipped entirely. This is
// what makes OSRO's per-method evaluation meaningful: each method's
// effect on native AVB WCD is measured against a timeline that only
// contains the busy windows that method would actually emit.
//
// Native TT ordering matches objectives.OBJ Round 1 then Round 2.
// CAN2TT ordering matches routes.Get_KPath_Routing, but we now walk
// `II.CAN2TTRoutes` in lockstep with method.CAN2TTFlows so the
// route-to-flow pairing still works after FilterCAN2TTByMethod has
// trimmed II.CAN2TTRoutes to a single method's slice.
//
// EMSO fallback: when a CAN2TT flow's full payload cannot fit on its
// route inside the deadline horizon, we call EMSODisaggregate to try
// the Yan 2024 §5.3 disaggregate-and-retry strategy. If EMSO also
// fails the flow is recorded in UnplaceableFlows; OBJ Round 3 will
// then count it as a CAN2TT schedulability failure.
func BuildTimelineForRoutes(network *network.Network, cfg *config.Config, II, IIprime *routes.RouteSet, methodScope string) *NetworkTimeline {
	nt := NewNetworkTimeline(cfg.Network.Hyperperiod)
	if network == nil || II == nil || IIprime == nil {
		return nt
	}

	S := network.FlowSet.InputOMACOFlowSet()
	Sprime := network.FlowSet.BGOMACOFlowSet()

	// Round 1: background TT flows (native TSN).
	for nth, r := range IIprime.TSNRoutes {
		if nth >= len(Sprime.TSNFlows) {
			break
		}
		nt.RegisterTTFlow(
			Sprime.TSNFlows[nth],
			r,
			cfg.Network.ByteRate,
			fmt.Sprintf("tt-bg-%d", nth),
		)
	}
	// Round 2: input TT flows (native TSN).
	for nth, r := range II.TSNRoutes {
		if nth >= len(S.TSNFlows) {
			break
		}
		nt.RegisterTTFlow(
			S.TSNFlows[nth],
			r,
			cfg.Network.ByteRate,
			fmt.Sprintf("tt-input-%d", nth),
		)
	}

	// Round 3: CAN2TT flows (gateway encap output). Walk methods in
	// the same nested order routingFunc.Get_KPath_Routing used, and
	// align route consumption with the active method scope.
	if network.FlowSet.EncapsulateMethod != nil {
		routeIdx := 0
		for _, method := range network.FlowSet.EncapsulateMethod {
			if methodScope != "" && method.MethodName != methodScope {
				// Caller filtered II.CAN2TTRoutes to methodScope already,
				// so we must NOT advance routeIdx for non-matching methods.
				continue
			}
			for j, canFlow := range method.CAN2TTFlows {
				if routeIdx >= len(II.CAN2TTRoutes) {
					break
				}
				route := II.CAN2TTRoutes[routeIdx]
				routeIdx++
				if route == nil {
					continue
				}
				flowID := fmt.Sprintf("can2tt-%s-%d", method.MethodName, j)
				if nt.RegisterCAN2TTFlow(canFlow, route, cfg.Network.ByteRate, flowID) {
					continue
				}
				// Direct placement failed — invoke EMSO disaggregation.
				if !EMSODisaggregate(nt, canFlow, route, cfg.Network.ByteRate, flowID, EMSOMaxAttempts) {
					// EMSO also failed; UnplaceableFlows already recorded
					// by registerFlowGeneric, OBJ Round 3 will fail this
					// flow in the schedulability count.
				}
			}
		}
	}

	return nt
}
