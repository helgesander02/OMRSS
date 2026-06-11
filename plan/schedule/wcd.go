package schedule

import (
	"src/network/flow"
	"src/network/flow/tt"
	"src/plan/routes"
	"time"
)

// Worse-Case Delay
func WCD(z *routes.Route, krouteSet *routes.KRouteSet, flow *tt.Flow, flowSet *flow.FlowSet) time.Duration {
	end2end := time.Duration(0)
	node := z.GetNodeByID(flow.Source)
	wcd := end2endDelay(node, -1, end2end, z, krouteSet, flow, flowSet)
	//logger.Printf("max wcd: %v \n", wcd)

	return wcd
}

// Use DFS to find all dataflow paths in the Route
// Calculate the End to End Delay for each dataflow path and select the maximum one
func end2endDelay(node *routes.Node, parentID int, end2end time.Duration, z *routes.Route, krouteSet *routes.KRouteSet, flow *tt.Flow, flowSet *flow.FlowSet) time.Duration {
	//logger.Printf("%d: %v \n", node.ID, end2end)
	maxE2E := end2end
	for _, link := range node.Connections {
		per_hop := time.Duration(0)
		if link.ToNodeID == parentID {
			continue

		} else {
			// Calculation of latency for a single link
			per_hop += transmitAVBItself(flow.DataSize, link.Cost)
			//per_hop += interfereFromBE(conn.Cost)
			per_hop += interfereFromAVB(link, krouteSet, flow.DataSize)
			per_hop += interfereFromTSN(link, krouteSet, flowSet)
			end2end += per_hop

			nextnode := z.GetNodeByID(link.ToNodeID)
			nextE2E := end2endDelay(nextnode, node.ID, end2end, z, krouteSet, flow, flowSet)

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
