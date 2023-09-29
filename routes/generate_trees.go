package routes

import (
	"fmt"
	"src/flow"
	"src/network"
)

func GetRoutes(topology *network.Topology, Input_flow_set *flow.Flows, BG_flow_set *flow.Flows, cost float64, K int) (*KTrees_set, *Trees_set, *Trees_set) {
	ktrees_set := &KTrees_set{}
	Input_tree_set := &Trees_set{}
	BG_tree_set := &Trees_set{}
	v2v := &V2V{}

	for round := 0; round < 2; round++ {
		if round == 0 {
			for _, flow := range Input_flow_set.TSNFlows {
				tree := SteninerTree(v2v, topology, flow.Source, flow.Destinations, cost)
				Input_tree_set.TSNTrees = append(Input_tree_set.TSNTrees, tree)
				Ktrees := KSpanningTree(v2v, tree, K, flow.Source, flow.Destinations, cost)
				ktrees_set.TSNTrees = append(ktrees_set.TSNTrees, Ktrees)
			}
			fmt.Printf("Finish %d input TSN streams trees\n", len(Input_tree_set.TSNTrees))

			for _, flow := range Input_flow_set.AVBFlows {
				tree := SteninerTree(v2v, topology, flow.Source, flow.Destinations, cost)
				Input_tree_set.AVBTrees = append(Input_tree_set.AVBTrees, tree)
				Ktrees := KSpanningTree(v2v, tree, K, flow.Source, flow.Destinations, cost)
				ktrees_set.AVBTrees = append(ktrees_set.AVBTrees, Ktrees)
			}
			fmt.Printf("Finish %d input AVB streams trees\n", len(Input_tree_set.AVBTrees))

		} else {
			for _, flow := range BG_flow_set.TSNFlows {
				tree := SteninerTree(v2v, topology, flow.Source, flow.Destinations, cost)
				BG_tree_set.TSNTrees = append(BG_tree_set.TSNTrees, tree)
				Ktrees := KSpanningTree(v2v, tree, K, flow.Source, flow.Destinations, cost)
				ktrees_set.TSNTrees = append(ktrees_set.TSNTrees, Ktrees)
			}
			fmt.Printf("Finish %d background TSN streams trees\n", len(BG_tree_set.TSNTrees))

			for _, flow := range BG_flow_set.AVBFlows {
				tree := SteninerTree(v2v, topology, flow.Source, flow.Destinations, cost)
				BG_tree_set.AVBTrees = append(BG_tree_set.AVBTrees, tree)
				Ktrees := KSpanningTree(v2v, tree, K, flow.Source, flow.Destinations, cost)
				ktrees_set.AVBTrees = append(ktrees_set.AVBTrees, Ktrees)
			}
			fmt.Printf("Finish %d background AVB streams trees\n", len(BG_tree_set.AVBTrees))

		}
	}

	return ktrees_set, Input_tree_set, BG_tree_set
}
