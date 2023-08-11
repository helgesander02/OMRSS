package trees

import (
	"fmt"

	"src/flow"
	"src/topology"
)

func GetTrees(topology *topology.Topology, flows *flow.Flows, cost float64, K int) *Trees {
	trees := &Trees{}
	v2v := &V2V{}

	for _, flow := range flows.TSNFlows {
		tree := SteninerTree(v2v, topology, flow.Source, flow.Destinations, cost)
		Ktrees := KSpanningTree(v2v, tree, K, flow.Source, flow.Destinations, cost)
		trees.TSNTrees = append(trees.TSNTrees, Ktrees)

		break
	}
	fmt.Println("Finish TSN all Trees")

	//for _, flow := range flows.AVBFlows {
	//tree := SteninerTree(v2v, topology, flow.Source, flow.Destinations, cost)
	//Ktrees := KSpanningTree(v2v, tree, K, flow.Source, flow.Destinations)
	//trees.AVBTrees = append(trees.AVBTrees, Ktrees)

	//}
	//fmt.Println("Finish AVB all Trees")

	return trees
}
