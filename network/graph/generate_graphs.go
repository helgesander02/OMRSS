package graph

import (
	"encoding/json"
	"fmt"
	"src/network/flow"
	"src/network/topology"
)

func Generate_Graphs(topology *topology.Topology, flows *flow.Flows, bytes_rate float64) *Graphs {
	graphs := &Graphs{}

	for _, flow := range flows.TSNFlows {
		t := TopologyDeepCopy(topology)                        // Duplicate of Topology
		t.AddN2S2N(flow.Source, flow.Destinations, bytes_rate) // Undirected Graph
		graphs.TSNGraphs = append(graphs.TSNGraphs, t)
	}

	for _, flow := range flows.AVBFlows {
		t := TopologyDeepCopy(topology)                        // Duplicate of Topology
		t.AddN2S2N(flow.Source, flow.Destinations, bytes_rate) // Undirected Graph
		graphs.AVBGraphs = append(graphs.AVBGraphs, t)
	}
	fmt.Println("Complete Generating Graphs")

	return graphs
}

func TopologyDeepCopy(t1 *topology.Topology) *topology.Topology {
	if buf, err := json.Marshal(t1); err != nil {
		return nil
	} else {
		t2 := &topology.Topology{}
		if err = json.Unmarshal(buf, t2); err != nil {
			return nil
		}
		return t2
	}
}
