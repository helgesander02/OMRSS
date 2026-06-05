package graph

import (
	"src/network/flow"
	"src/network/topology"
	"src/pkg/logger"
)

func GenerateOmacoGraphs(topology *topology.Topology, flows *flow.FlowSet, bytesRate float64) *Graphs {
	// Constructing Graph structures
	graphs := newGraphs()

	// Generating TSN Graphs
	for _, flow := range flows.TSNFlows {
		t := topology.TopologyDeepCopy()                               // Duplicate of Topology
		t.AddN2S2N_For_Tree(flow.Source, flow.Destinations, bytesRate) // Undirected Graph
		graphs.TSNGraphs = append(graphs.TSNGraphs, t)
	}
	logger.Println(len(graphs.TSNGraphs), " TSN Graphs generated.")

	// Generating AVB Graphs
	for _, flow := range flows.AVBFlows {
		t := topology.TopologyDeepCopy()                               // Duplicate of Topology
		t.AddN2S2N_For_Tree(flow.Source, flow.Destinations, bytesRate) // Undirected Graph
		graphs.AVBGraphs = append(graphs.AVBGraphs, t)
	}
	logger.Println(len(graphs.AVBGraphs), " AVB Graphs generated.")

	return graphs
}

func GenerateOsroGraphs(topology *topology.Topology, flows *flow.FlowSet, bytesRate float64) *Graphs {
	// Constructing Graph structures
	graphs := newGraphs()

	// Generating TSN Graphs
	for _, flow := range flows.TSNFlows {
		t := topology.TopologyDeepCopy()                                  // Duplicate of Topology
		t.AddN2S2N_For_Path(flow.Source, flow.Destinations[0], bytesRate) // Undirected Graph
		graphs.TSNGraphs = append(graphs.TSNGraphs, t)
	}
	logger.Println(len(graphs.TSNGraphs), " TSN Graphs generated.")

	// Generating AVB Graphs
	for _, flow := range flows.AVBFlows {
		t := topology.TopologyDeepCopy()                                  // Duplicate of Topology
		t.AddN2S2N_For_Path(flow.Source, flow.Destinations[0], bytesRate) // Undirected Graph
		graphs.AVBGraphs = append(graphs.AVBGraphs, t)
	}
	logger.Println(len(graphs.AVBGraphs), " AVB Graphs generated.")

	// Generating CAN2TT Graphs
	for _, method := range flows.EncapsulateMethod {
		for _, can2tsnflow := range method.CAN2TTFlows {
			if !graphs.checkListenerAndTalker(can2tsnflow.Source, can2tsnflow.Destination) {
				t := topology.TopologyDeepCopy()                                            // Duplicate of Topology
				t.AddN2S2N_For_Path(can2tsnflow.Source, can2tsnflow.Destination, bytesRate) // Undirected Graph
				graphs.CAN2TSNGraphs = append(graphs.CAN2TSNGraphs, t)

			}

		}
	}
	logger.Println(len(graphs.CAN2TSNGraphs), " CAN2TT Graphs generated.")

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
