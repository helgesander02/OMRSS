package graph

import (
	"src/network/flow"
	"src/network/topology"
)

func Generate_Graphs(topology *topology.Topology, flows *flow.Flows, bytes_rate float64) *Graphs {
	graphs := &Graphs{}

	for _, flow := range flows.TSNFlows {
		t := topology.TopologyDeepCopy()                       // Duplicate of Topology
		t.AddN2S2N(flow.Source, flow.Destinations, bytes_rate) // Undirected Graph
		graphs.TSNGraphs = append(graphs.TSNGraphs, t)
	}

	for _, flow := range flows.AVBFlows {
		t := topology.TopologyDeepCopy()                       // Duplicate of Topology
		t.AddN2S2N(flow.Source, flow.Destinations, bytes_rate) // Undirected Graph
		graphs.AVBGraphs = append(graphs.AVBGraphs, t)
	}

	return graphs
}
