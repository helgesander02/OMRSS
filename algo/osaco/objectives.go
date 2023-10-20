package osaco

import (
	"src/network"
	"src/routes"
	"time"
)

// Objectives
func Obj(network *network.Network, X *routes.KTrees_set, II *routes.Trees_set, II_prime *routes.Trees_set) [4]float64 {
	S := network.Flow_Set.Input_flow_set()
	S_prime := network.Flow_Set.BG_flow_set()
	var (
		obj                [4]float64
		tsn_failed_count   int = 0
		avb_failed_count   int = 0
		all_rerouted_count int = 0
		avb_wcd_sum        time.Duration
	)

	// O1
	for nth := range II.TSNTrees {
		tsn_failed_count += 1 - Schedulability(0, S.TSNFlows[nth].Deadline)
	}
	for nth := range II_prime.TSNTrees {
		tsn_failed_count += 1 - Schedulability(0, S_prime.TSNFlows[nth].Deadline)
	}

	// O2 and O4
	for nth, tree := range II.AVBTrees {
		wcd := WCD(tree, X, S.AVBFlows[nth], network.Flow_Set)
		avb_wcd_sum += wcd
		avb_failed_count += 1 - Schedulability(wcd, S.AVBFlows[nth].Deadline)
	}
	for nth, tree := range II_prime.AVBTrees {
		wcd := WCD(tree, X, S_prime.AVBFlows[nth], network.Flow_Set)
		avb_wcd_sum += wcd
		avb_failed_count += 1 - Schedulability(wcd, S_prime.AVBFlows[nth].Deadline)
	}

	// O3

	obj[0] = float64(tsn_failed_count)
	obj[1] = float64(avb_failed_count)
	obj[2] = float64(all_rerouted_count)
	obj[3] = float64(avb_wcd_sum)

	return obj
}

func Schedulability(wcd time.Duration, deadline int) int {
	if wcd <= time.Duration(deadline)*time.Microsecond {
		return 1
	}
	return 0
}
