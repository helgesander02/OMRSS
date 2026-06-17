package schedule

import (
	"fmt"
	"math"
	"src/network/flow"
	"src/network/flow/tt"
	"src/plan/routes"
	"time"
)

// avbObservationIntervalUs is the IEEE 802.1Qav Class A observation
// interval. The Go AVB spec carries no Class field, so every AVB stream
// is treated as Class A — the more conservative of the two intervals
// (125 μs vs 250 μs for Class B). frame_count = ceil(interval / period).
const avbObservationIntervalUs = 125

// AVBPerLinkLoad maps "from>to" to the total AVB bytes per observation
// interval contributed by every SELECTED AVB stream that traverses the
// link. This is the timeline-aware replacement for the legacy
// interfereFromAVB which used to walk every K-alternative — that
// inflated the byte count by ~K and ignored that only one alternative
// is actually selected by OSACO.
type AVBPerLinkLoad map[string]float64

// ComputeAVBPerLinkLoad walks the SELECTED background + input AVB
// routes (II_prime + II) and tallies the per-link byte contribution.
// Each stream contributes datasize × ceil(observation_interval / period)
// — the AVB Latency Math observation interval, not the full hyperperiod.
//
// The walk follows the route DFS-style from the stream's source. The
// parent filter prevents re-traversing the back-pointer connections
// that IntoTree / ConvertIDsToPath insert for bidirectional access.
func ComputeAVBPerLinkLoad(II, IIprime *routes.RouteSet, flowSet *flow.FlowSet) AVBPerLinkLoad {
	result := make(AVBPerLinkLoad)
	if flowSet == nil {
		return result
	}
	bgFlows := flowSet.BGOMACOFlowSet()
	inputFlows := flowSet.InputOMACOFlowSet()

	addRoute := func(route *routes.Route, avb *tt.Flow) {
		if route == nil || avb == nil {
			return
		}
		bytesPerObs := avbBytesPerObservation(avb)
		src := route.GetNodeByID(avb.Source)
		if src == nil {
			return
		}
		addLoadDirected(src, -1, route, bytesPerObs, result)
	}

	if IIprime != nil {
		for nth, route := range IIprime.AVBRoutes {
			if nth >= len(bgFlows.AVBFlows) {
				break
			}
			addRoute(route, bgFlows.AVBFlows[nth])
		}
	}
	if II != nil {
		for nth, route := range II.AVBRoutes {
			if nth >= len(inputFlows.AVBFlows) {
				break
			}
			addRoute(route, inputFlows.AVBFlows[nth])
		}
	}
	return result
}

// addLoadDirected walks the route DFS from `node` using the parent
// filter so each directed link is hit exactly once on both tree-shaped
// (multicast) and path-shaped (unicast) routes.
func addLoadDirected(node *routes.Node, parentID int, route *routes.Route, bytes float64, result AVBPerLinkLoad) {
	if node == nil {
		return
	}
	for _, conn := range node.Connections {
		if conn.ToNodeID == parentID {
			continue
		}
		result[fmt.Sprintf("%d>%d", node.ID, conn.ToNodeID)] += bytes
		next := route.GetNodeByID(conn.ToNodeID)
		if next != nil {
			addLoadDirected(next, node.ID, route, bytes, result)
		}
	}
}

// avbBytesPerObservation returns the bytes a single AVB stream
// contributes to one link inside the 125 μs Class A observation
// interval (paper-aligned with IEEE 802.1Qav / Laursen 2016).
func avbBytesPerObservation(avb *tt.Flow) float64 {
	if avb == nil {
		return 0
	}
	framesInObs := 1.0
	if avb.Period > 0 {
		framesInObs = math.Ceil(float64(avbObservationIntervalUs) / float64(avb.Period))
		if framesInObs < 1 {
			framesInObs = 1
		}
	}
	return avb.DataSize * framesInObs
}

// WCD is the legacy timeline-free path. OSACO::computeVb calls this
// before any route is selected (visibility heuristic), so we approximate
// the unknown selections by using each stream's first (shortest) K
// alternative. This is a heuristic only — final OBJ scoring uses
// WCDWithTimeline which knows the actual selection.
func WCD(z *routes.Route, krouteSet *routes.KRouteSet, flow *tt.Flow, flowSet *flow.FlowSet) time.Duration {
	end2end := time.Duration(0)
	node := z.GetNodeByID(flow.Source)
	wcd := end2endDelay(node, -1, end2end, z, krouteSet, flow, flowSet, nil, nil)
	return wcd
}

// WCDWithTimeline is the timeline-aware variant. The per-link AVB
// interference is computed from the precomputed selected-route load
// map, and TT interference is computed via Laursen 2016 §6.3 Algorithm
// 1 directly off the timeline. Both terms react correctly to OSRO's
// per-method scoping because timeline + avbLoad were built only from
// flows in scope.
func WCDWithTimeline(z *routes.Route, krouteSet *routes.KRouteSet, flow *tt.Flow, flowSet *flow.FlowSet, timeline *NetworkTimeline, avbLoad AVBPerLinkLoad) time.Duration {
	end2end := time.Duration(0)
	node := z.GetNodeByID(flow.Source)
	wcd := end2endDelay(node, -1, end2end, z, krouteSet, flow, flowSet, timeline, avbLoad)
	return wcd
}

// end2endDelay walks the route DFS-style and returns the max E2E delay
// across all destinations. The timeline / avbLoad params select between
// the legacy heuristic path (both nil, used by OSACO visibility init)
// and the timeline-aware path (both non-nil, used by OBJ).
func end2endDelay(node *routes.Node, parentID int, end2end time.Duration, z *routes.Route, krouteSet *routes.KRouteSet, flow *tt.Flow, flowSet *flow.FlowSet, timeline *NetworkTimeline, avbLoad AVBPerLinkLoad) time.Duration {
	if node == nil {
		return end2end
	}
	maxE2E := end2end
	for _, link := range node.Connections {
		if link.ToNodeID == parentID {
			continue
		}
		ownTx := transmitAVBItself(flow.DataSize, link.Cost)
		per_hop := ownTx

		// AVB interference
		if avbLoad != nil {
			per_hop += interfereFromAVBSelected(link, flow, avbLoad)
		} else {
			per_hop += interfereFromAVB(link, krouteSet, flowSet, flow)
		}

		// TSN interference
		if timeline != nil {
			per_hop += interfereFromTSNTimeline(link, ownTx, timeline)
		} else {
			per_hop += interfereFromTSN(link, krouteSet, flowSet)
		}

		end2end += per_hop
		nextnode := z.GetNodeByID(link.ToNodeID)
		nextE2E := end2endDelay(nextnode, node.ID, end2end, z, krouteSet, flow, flowSet, timeline, avbLoad)
		if maxE2E < nextE2E {
			maxE2E = nextE2E
		}
		end2end -= per_hop
	}
	return maxE2E
}

// transmitAVBItself returns the credit-shaped per-hop transmission time
// of one AVB frame.
//
// bytesRate is microseconds-per-byte (cfg.Network.ByteRate, e.g. 0.008
// μs/B for 1 Gbps). The misleading comment on topology.Link.Cost calls
// the same value "125 bytes/us"; the actual stored number is 1/that
// rate — see pkg/config/config.go where ByteRate is set to
// 1.0 / ((Bandwidth / 8) * 1e-6).
//
// Raw transmission time: datasize × bytesRate (microseconds).
// CBS inflates that by 1/MAX_AVB_SETTING because the AVB queue is
// credit-capped to MAX_AVB_SETTING of the link bandwidth.
//
// IEEE 802.1BA AVB Latency Math — and Laursen 2016 §6.3 — both divide
// by MAX_AVB_SETTING here. The previous implementation multiplied by
// it (shortening the transmission time instead of inflating it) and
// then treated the resulting microsecond number as nanoseconds when
// it stuffed it into time.Duration, silently dropping a factor of
// 1000. Net effect: ~1800× under-count of every AVB per-hop delay.
func transmitAVBItself(datasize float64, bytesRate float64) time.Duration {
	const MAX_AVB_SETTING = 0.75
	if datasize <= 0 || bytesRate <= 0 {
		return 0
	}
	microseconds := datasize * bytesRate / MAX_AVB_SETTING
	return time.Duration(int64(microseconds * float64(time.Microsecond)))
}

// interfereFromAVB is the legacy heuristic used by WCD (OSACO
// visibility init). Each interfering stream is represented by its first
// K alternative since no route has been selected yet. Compared to the
// previous implementation, this fixes three bugs:
//  1. K-fold over-count — old code walked every K alternative;
//  2. wrong byte size — old code used the querying flow's datasize for
//     every neighbour, not the neighbour's own datasize;
//  3. period blindness — old code ignored frame_count(observation,
//     period).
//
// It still cannot tell which alternative will actually be picked by
// OSACO, so for final scoring use WCDWithTimeline + avbLoad instead.
func interfereFromAVB(link *routes.Connection, krouteSet *routes.KRouteSet, flowSet *flow.FlowSet, ownFlow *tt.Flow) time.Duration {
	if link == nil || krouteSet == nil || flowSet == nil {
		return 0
	}
	bgFlows := flowSet.BGOMACOFlowSet()
	inputFlows := flowSet.InputOMACOFlowSet()
	bgAVBLen := len(bgFlows.AVBFlows)

	var occupiedbytes float64
	for nth, avb_ktree := range krouteSet.AVBRoutes {
		if avb_ktree == nil || len(avb_ktree.Routes) == 0 {
			continue
		}
		var other *tt.Flow
		if nth < bgAVBLen {
			other = bgFlows.AVBFlows[nth]
		} else {
			idx := nth - bgAVBLen
			if idx >= len(inputFlows.AVBFlows) {
				continue
			}
			other = inputFlows.AVBFlows[idx]
		}
		if other == nil || other == ownFlow {
			continue
		}
		// Use only the first (shortest) K alternative since OSACO has
		// not yet decided which alternative each stream will take.
		tree := avb_ktree.Routes[0]
		node := tree.GetNodeByID(link.FromNodeID)
		if node == nil {
			continue
		}
		for _, conn := range node.Connections {
			if conn.ToNodeID == link.ToNodeID {
				occupiedbytes += avbBytesPerObservation(other)
				break
			}
		}
	}
	return transmitAVBItself(occupiedbytes, link.Cost)
}

// interfereFromAVBSelected returns the AVB-on-AVB interference time on
// link (from, to) using the precomputed selected-route load map. The
// own flow's contribution is subtracted so it never interferes with
// itself.
func interfereFromAVBSelected(link *routes.Connection, ownFlow *tt.Flow, load AVBPerLinkLoad) time.Duration {
	if link == nil || load == nil {
		return 0
	}
	totalBytes := load[fmt.Sprintf("%d>%d", link.FromNodeID, link.ToNodeID)]
	ownBytes := avbBytesPerObservation(ownFlow)
	otherBytes := totalBytes - ownBytes
	if otherBytes < 0 {
		otherBytes = 0
	}
	return transmitAVBItself(otherBytes, link.Cost)
}

// interfereFromTSNTimeline implements Laursen 2016 §6.3 Algorithm 1.
// See timeline.LaursenInterference for the actual walk.
func interfereFromTSNTimeline(link *routes.Connection, ownTx time.Duration, timeline *NetworkTimeline) time.Duration {
	twcUs := int(ownTx / time.Microsecond)
	if twcUs <= 0 {
		// AVB's own per-hop transmission rounds to zero µs at very
		// small payloads; fall back to 1 µs probe so we don't get a
		// trivially-zero interference estimate.
		twcUs = 1
	}
	imax := timeline.LaursenInterference(link.FromNodeID, link.ToNodeID, twcUs)
	return time.Duration(imax) * time.Microsecond
}

// interfereFromTSN is the legacy heuristic used by WCD when no timeline
// is available. Like interfereFromAVB above, it now uses each stream's
// first K alternative to avoid K-fold over-counting. The byte-per-
// hyperperiod term is unchanged because TSN frames are gate-scheduled,
// not credit-shaped, so the AVB observation interval does not apply —
// every TSN instance inside the hyperperiod must transmit fully.
func interfereFromTSN(link *routes.Connection, krouteSet *routes.KRouteSet, flowSet *flow.FlowSet) time.Duration {
	if link == nil || krouteSet == nil || flowSet == nil {
		return 0
	}
	var occupiedbytes float64
	for nth, tsn_ktree := range krouteSet.TSNRoutes {
		if tsn_ktree == nil || len(tsn_ktree.Routes) == 0 {
			continue
		}
		if nth >= len(flowSet.TSNFlows) {
			continue
		}
		ts := flowSet.TSNFlows[nth]
		if ts == nil || ts.Period <= 0 {
			continue
		}
		tree := tsn_ktree.Routes[0]
		node := tree.GetNodeByID(link.FromNodeID)
		if node == nil {
			continue
		}
		for _, conn := range node.Connections {
			if conn.ToNodeID == link.ToNodeID {
				occupiedbytes += ts.DataSize * (float64(ts.HyperPeriod) / float64(ts.Period))
				break
			}
		}
	}
	return transmitAVBItself(occupiedbytes, link.Cost)
}
