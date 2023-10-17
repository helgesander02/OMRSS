package routes

import (
	"fmt"
	"src/network"
)

func GetRoutes(network *network.Network, K int) *KTrees_set {
	ktrees_set := &KTrees_set{}
	v2v := &V2V{}

	for nth, flow := range network.Flow_Set.TSNFlows {
		tree := SteninerTree(v2v, network.Graph_Set.TSNGraphs[nth], flow.Source, flow.Destinations, network.Bytes_Rate)
		Ktrees := KSpanningTree(v2v, tree, K, flow.Source, flow.Destinations, network.Bytes_Rate)
		ktrees_set.TSNTrees = append(ktrees_set.TSNTrees, Ktrees)
	}
	fmt.Printf("Finish %d TSN streams routes\n", len(ktrees_set.TSNTrees))

	for nth, flow := range network.Flow_Set.AVBFlows {
		tree := SteninerTree(v2v, network.Graph_Set.AVBGraphs[nth], flow.Source, flow.Destinations, network.Bytes_Rate)
		Ktrees := KSpanningTree(v2v, tree, K, flow.Source, flow.Destinations, network.Bytes_Rate)
		ktrees_set.AVBTrees = append(ktrees_set.AVBTrees, Ktrees)
	}
	fmt.Printf("Finish %d AVB streams routes\n", len(ktrees_set.AVBTrees))

	return ktrees_set
}
