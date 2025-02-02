package schedule

import (
	"fmt"
	"src/network"
	"src/network/flow"
	"src/plan/routes"

	"time"
)

// Objectives
func OBJ(network *network.Network, X *routes.KTrees_set, II *routes.Trees_set, II_prime *routes.Trees_set, cost_setting [4]int) ([4]float64, int) {
	S := network.Flow_Set.Input_flow_set()
	S_prime := network.Flow_Set.BG_flow_set()
	var (
		obj                [4]float64
		cost               int
		tsn_failed_count   int           = 0 // O1
		avb_failed_count   int           = 0 // O2
		all_rerouted_count int           = 0 // O3 ... pass
		avb_wcd_sum        time.Duration     // O4
	)
	linkmap := map[string]float64{}

	// Round1: Schedule BG flow
	// O1
	for nth, route := range II_prime.TSNTrees {
		schedulability := schedulability(0, S_prime.TSNFlows[nth], route, linkmap, network.Bandwidth, network.HyperPeriod)
		tsn_failed_count += 1 - schedulability
		//fmt.Printf("BackGround TSN route%d: %b \n", nth, schedulability)
	}

	// O2 and O4
	for nth, route := range II_prime.AVBTrees {
		wcd := WCD(route, X, S_prime.AVBFlows[nth], network.Flow_Set)
		avb_wcd_sum += wcd
		schedulability := schedulability(wcd, S_prime.AVBFlows[nth], route, linkmap, network.Bandwidth, network.HyperPeriod)
		avb_failed_count += 1 - schedulability
		//fmt.Printf("BackGround AVB route%d: %b \n", nth, schedulability)
	}
	// O3 ... pass

	// Round2: Schedule Input flow
	// O1
	for nth, route := range II.TSNTrees {
		schedulability := schedulability(0, S.TSNFlows[nth], route, linkmap, network.Bandwidth, network.HyperPeriod)
		tsn_failed_count += 1 - schedulability
		//fmt.Printf("Input TSN route%d: %b \n", nth, schedulability)
	}

	// O2 and O4
	for nth, route := range II.AVBTrees {
		wcd := WCD(route, X, S.AVBFlows[nth], network.Flow_Set)
		avb_wcd_sum += wcd
		schedulability := schedulability(wcd, S.AVBFlows[nth], route, linkmap, network.Bandwidth, network.HyperPeriod)
		avb_failed_count += 1 - schedulability
		//fmt.Printf("Input AVB route%d: %b \n", nth, schedulability)
	}
	// O3 ... pass

	obj[0] = float64(tsn_failed_count)               // O1
	obj[1] = float64(avb_failed_count)               // O2
	obj[2] = float64(all_rerouted_count)             // O3 ... pass
	obj[3] = float64(avb_wcd_sum / time.Microsecond) // O4

	cost += tsn_failed_count * cost_setting[0]                  // O1
	cost += avb_failed_count * cost_setting[1]                  // O2
	cost += int(avb_wcd_sum/time.Microsecond) * cost_setting[3] // O4

	return obj, cost
}

func schedulability(wcd time.Duration, flow *flow.Flow, route *routes.Tree, linkmap map[string]float64, bandwidth float64, hyperPeriod int) int {
	r := wcd <= time.Duration(flow.Deadline)*time.Microsecond
	node := route.GetNodeByID(flow.Source)
	schedulable, _ := schedulable(node, -1, flow, route, linkmap, bandwidth, hyperPeriod)

	if r && schedulable {
		return 1
	}
	return 0
}

func schedulable(node *routes.Node, parentID int, flow *flow.Flow, route *routes.Tree, linkmap map[string]float64, bandwidth float64, hyperPeriod int) (bool, map[string]float64) {
	for _, link := range node.Connections {
		if link.ToNodeID == parentID {
			continue

		} else {
			//// Duplex
			if !(link.FromNodeID == flow.Source || loopcompare(link.ToNodeID, flow.Destinations)) {
				key := fmt.Sprintf("%d>%d", link.FromNodeID, link.ToNodeID)
				linkmap[key] += flow.DataSize * float64((hyperPeriod / flow.Period))
				if linkmap[key] > bandwidth {
					return false, linkmap
				}
			}

			//// Simplex
			//if !(link.FromNodeID == flow.Source || loopcompare(link.ToNodeID, flow.Destinations)) {
			//	key := ""
			//	key1 := fmt.Sprintf("%d>%d", link.FromNodeID, link.ToNodeID)
			//	key2 := fmt.Sprintf("%d>%d", link.ToNodeID, link.FromNodeID)
			//	if _, ok := linkmap[key1]; !ok {
			//		if _, ok := linkmap[key2]; !ok {
			//			key = key1
			//		} else {
			//			key = key2
			//		}
			//
			//	} else {
			//		key = key1
			//	}
			//
			//	linkmap[key] += flow.DataSize * float64((hyperPeriod / flow.Period))
			//	if linkmap[key] > bandwidth {
			//		return false, linkmap
			//	}
			//}

			nextnode := route.GetNodeByID(link.ToNodeID)
			schedulable, updatedLinkmap := schedulable(nextnode, node.ID, flow, route, linkmap, bandwidth, hyperPeriod)
			if !schedulable {
				return false, updatedLinkmap
			}
		}
	}
	return true, linkmap
}

func loopcompare(a int, b []int) bool {
	for _, v := range b {
		if a == v {
			return true
		}
	}
	return false
}
