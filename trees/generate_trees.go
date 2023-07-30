package trees

import (   
    "fmt"
    
    "src/topology"
    "src/flow"
)

func GetTrees(topology *topology.Topology, flows *flow.Flows, cost float64, K int) *Trees {
    trees := &Trees{}
    
    for _, flow := range flows.TSNFlows {
        tree := SteninerTree(topology, flow, cost)
        Ktrees := KSpanningTree(tree, K)
        trees.TSNTrees = append(trees.TSNTrees, Ktrees)
        break

    }
    fmt.Println("Finish TSN all Trees")

    for _, flow := range flows.AVBFlows {
        tree := SteninerTree(topology, flow, cost)
        Ktrees := KSpanningTree(tree, K)
        trees.AVBTrees = append(trees.AVBTrees, Ktrees)
        break
        
    }
    fmt.Println("Finish AVB all Trees")

    return trees
}