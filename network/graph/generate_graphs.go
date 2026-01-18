package graph

import (
	"fmt"
	"src/network/flow"
	"src/network/topology"
)

func Generate_OMACO_Graphs(topology *topology.Topology, flows *flow.Flow_Set, bytes_rate float64) *Graphs {
	// Constructing Graph structures
	graphs := new_Graphs()

	// Generating TSN Graphs
	for _, flow := range flows.TSNFlows {
		t := topology.TopologyDeepCopy()                                // Duplicate of Topology
		t.AddN2S2N_For_Tree(flow.Source, flow.Destinations, bytes_rate) // Undirected Graph
		graphs.TSNGraphs = append(graphs.TSNGraphs, t)
	}
	fmt.Println(len(graphs.TSNGraphs), " TSN Graphs generated.")

	// Generating AVB Graphs
	for _, flow := range flows.AVBFlows {
		t := topology.TopologyDeepCopy()                                // Duplicate of Topology
		t.AddN2S2N_For_Tree(flow.Source, flow.Destinations, bytes_rate) // Undirected Graph
		graphs.AVBGraphs = append(graphs.AVBGraphs, t)
	}
	fmt.Println(len(graphs.AVBGraphs), " AVB Graphs generated.")

	return graphs
}

func Generate_OSRO_Graphs(topology *topology.Topology, flows *flow.Flow_Set, bytes_rate float64) *Graphs {
	// Constructing Graph structures
	graphs := new_Graphs()

	// Generating TSN Graphs
	for _, flow := range flows.TSNFlows {
		t := topology.TopologyDeepCopy()                                   // Duplicate of Topology
		t.AddN2S2N_For_Path(flow.Source, flow.Destinations[0], bytes_rate) // Undirected Graph
		graphs.TSNGraphs = append(graphs.TSNGraphs, t)
	}
	fmt.Println(len(graphs.TSNGraphs), " TSN Graphs generated.")

	// Generating AVB Graphs
	for _, flow := range flows.AVBFlows {
		t := topology.TopologyDeepCopy()                                   // Duplicate of Topology
		t.AddN2S2N_For_Path(flow.Source, flow.Destinations[0], bytes_rate) // Undirected Graph
		graphs.AVBGraphs = append(graphs.AVBGraphs, t)
	}
	fmt.Println(len(graphs.AVBGraphs), " AVB Graphs generated.")

	// Generating CAN2TT Graphs
	for _, method := range flows.EncapsulateMethod {
		for _, can2tsnflow := range method.CAN2TTFlows {
			if !graphs.checkListenerAndTalker(can2tsnflow.Source, can2tsnflow.Destination) {
				t := topology.TopologyDeepCopy()                                             // Duplicate of Topology
				t.AddN2S2N_For_Path(can2tsnflow.Source, can2tsnflow.Destination, bytes_rate) // Undirected Graph
				graphs.CAN2TSNGraphs = append(graphs.CAN2TSNGraphs, t)

			}

		}
	}
	fmt.Println(len(graphs.CAN2TSNGraphs), " CAN2TT Graphs generated.")

	return graphs
}

func (graph *Graphs) checkListenerAndTalker(source int, destination int) bool {
	for _, g := range graph.CAN2TSNGraphs {
		if g.GetListenerAndTalker(source, destination) {
			return true
		}
	}
	return false
}

func (graph *Graphs) GetGarphBySD(source int, destination int) *topology.Topology {
	for _, g := range graph.CAN2TSNGraphs {
		if g.GetListenerAndTalker(source, destination) {
			return g
		}
	}
	return nil
}
