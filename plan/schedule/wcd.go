package schedule

import (
	"src/network/flow"
	"src/network/flow/tt"
	"src/plan/routes"
	"time"
)

// Worse-Case Delay (legacy, timeline-free path).
//
// This path uses the bandwidth-ratio TT interference model (see
// interfereFromTSN below). It is preserved because OSACO's visibility
// initialiser (algo/osaco.go::computeVb) calls WCD before any routing
// solution has been picked, so there is no timeline to consult.
//
// OBJ uses WCDWithTimeline instead, which models TT interference via
// Laursen 2016 §6.3 Algorithm 1 and so reacts correctly to encap-method
// changes in OSRO.
func WCD(z *routes.Route, krouteSet *routes.KRouteSet, flow *tt.Flow, flowSet *flow.FlowSet) time.Duration {
	end2end := time.Duration(0)
	node := z.GetNodeByID(flow.Source)
	wcd := end2endDelay(node, -1, end2end, z, krouteSet, flow, flowSet, nil)
	return wcd
}

// WCDWithTimeline is the timeline-aware variant. When timeline is
// non-nil, the per-link interference term is computed by
// NetworkTimeline.LaursenInterference (Laursen 2016 §6.3 Algorithm 1)
// instead of the bandwidth-ratio approximation. This is what makes
// OSRO's encap methods produce different AVB WCDs — fifo, mao, wst_mar
// drop different busy-window patterns into the timeline, and the AVB
// frame's worst-case stall responds to those patterns.
func WCDWithTimeline(z *routes.Route, krouteSet *routes.KRouteSet, flow *tt.Flow, flowSet *flow.FlowSet, timeline *NetworkTimeline) time.Duration {
	end2end := time.Duration(0)
	node := z.GetNodeByID(flow.Source)
	wcd := end2endDelay(node, -1, end2end, z, krouteSet, flow, flowSet, timeline)
	return wcd
}

// Use DFS to find all dataflow paths in the Route
// Calculate the End to End Delay for each dataflow path and select the
// maximum one. When `timeline` is non-nil, the per-link TT interference
// term is computed by Laursen 2016 §6.3 Algorithm 1; otherwise it falls
// back to the legacy bandwidth-ratio approximation.
func end2endDelay(node *routes.Node, parentID int, end2end time.Duration, z *routes.Route, krouteSet *routes.KRouteSet, flow *tt.Flow, flowSet *flow.FlowSet, timeline *NetworkTimeline) time.Duration {
	maxE2E := end2end
	for _, link := range node.Connections {
		per_hop := time.Duration(0)
		if link.ToNodeID == parentID {
			continue

		} else {
			// Calculation of latency for a single link
			ownTx := transmitAVBItself(flow.DataSize, link.Cost)
			per_hop += ownTx
			//per_hop += interfereFromBE(conn.Cost)
			per_hop += interfereFromAVB(link, krouteSet, flow.DataSize)
			if timeline != nil {
				per_hop += interfereFromTSNTimeline(link, ownTx, timeline)
			} else {
				per_hop += interfereFromTSN(link, krouteSet, flowSet)
			}
			end2end += per_hop

			nextnode := z.GetNodeByID(link.ToNodeID)
			nextE2E := end2endDelay(nextnode, node.ID, end2end, z, krouteSet, flow, flowSet, timeline)

			if maxE2E < nextE2E {
				maxE2E = nextE2E
			}
		}

		end2end -= per_hop
	}
	return maxE2E
}

// Calculate the transmission time of AVB
func transmitAVBItself(datasize float64, bytesRate float64) time.Duration {
	/// Maximum proportion of bandwidth that AVB streams can occupy.
	const MAX_AVB_SETTING float64 = 0.75
	nanoseconds := datasize * bytesRate * MAX_AVB_SETTING
	duration := time.Duration(int64(nanoseconds))

	return duration
}

// The time occupied by a BE packet before transmission
//func interfereFromBE(bytesRate float64) time.Duration {
//	// Maximum number of bytes in a frame.
//	const MTU float64 = 1500.
//	nanoseconds := MTU * bytesRate
//	duration := time.Duration(int64(nanoseconds))
//
//	return duration
//}

// The time occupied by other AVB packets during transmission
func interfereFromAVB(link *routes.Connection, krouteSet *routes.KRouteSet, datasize float64) time.Duration {
	// Occupied bytes by other AVB
	var occupiedbytes float64
	for _, avb_ktree := range krouteSet.AVBRoutes {
		for _, tree := range avb_ktree.Routes {
			node := tree.GetNodeByID(link.FromNodeID)
			if node != nil {
				for _, conn := range node.Connections {
					if conn.ToNodeID == link.ToNodeID {
						occupiedbytes += datasize
					}
				}
			}
		}
	}
	occupiedbytes -= datasize // Deducting its own datasize

	return transmitAVBItself(occupiedbytes, link.Cost)
}

// interfereFromTSNTimeline implements Laursen 2016 §6.3 Algorithm 1.
// It asks the timeline directly: "given an AVB transmission of `ownTx`
// µs on link, what is the worst-case µs of TT preemption it would
// suffer?" The timeline already contains both native TT and CAN2TT
// busy windows (placed during BuildTimelineForRoutes), so this single
// query covers both interference sources in one go.
//
// The Yan 2024-style argument that "encap method X impacts native AVB
// by Y µs" gets its causal teeth here: different encap methods drop
// different busy-window patterns onto the timeline, and AVB's
// worst-case anchor walk reflects those patterns directly. The flat
// bandwidth-ratio path (interfereFromTSN below) cannot tell methods
// apart because they all consume the same total bytes/hyperperiod.
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

// The known time occupied by TSN packets during transmission
func interfereFromTSN(link *routes.Connection, krouteSet *routes.KRouteSet, flowSet *flow.FlowSet) time.Duration {
	// Occupied bytes by TSN
	var occupiedbytes float64
	for nth, tsn_ktree := range krouteSet.TSNRoutes {
		for _, tree := range tsn_ktree.Routes {
			node := tree.GetNodeByID(link.FromNodeID)
			if node != nil {
				for _, conn := range node.Connections {
					if conn.ToNodeID == link.ToNodeID {
						// occupiedbytes += datasize * (hyperPeriod / period)
						occupiedbytes += flowSet.TSNFlows[nth].DataSize *
							(float64(flowSet.TSNFlows[nth].HyperPeriod) / float64(flowSet.TSNFlows[nth].Period))
					}
				}
			}
		}
	}

	return transmitAVBItself(occupiedbytes, link.Cost)
}
