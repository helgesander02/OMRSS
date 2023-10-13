package routes

import (
	"fmt"
	"src/flow"
	"src/network"
)

func GetRoutes(topology *network.Topology, flow_set *flow.Flows, cost float64, K int) *KTrees_set {
	ktrees_set := &KTrees_set{}
	v2v := &V2V{}

	for _, flow := range flow_set.TSNFlows {
		tree := SteninerTree(v2v, topology, flow.Source, flow.Destinations, cost)
		Ktrees := KSpanningTree(v2v, tree, K, flow.Source, flow.Destinations, cost)
		ktrees_set.TSNTrees = append(ktrees_set.TSNTrees, Ktrees)
	}
	fmt.Printf("Finish %d TSN streams routes\n", len(ktrees_set.TSNTrees))

	for _, flow := range flow_set.AVBFlows {
		tree := SteninerTree(v2v, topology, flow.Source, flow.Destinations, cost)
		Ktrees := KSpanningTree(v2v, tree, K, flow.Source, flow.Destinations, cost)
		ktrees_set.AVBTrees = append(ktrees_set.AVBTrees, Ktrees)
	}
	fmt.Printf("Finish %d AVB streams routes\n", len(ktrees_set.AVBTrees))

	return ktrees_set
}
