package network

import (
	"fmt"
	"src/network/flow"
	"src/network/graph"
	"src/network/topology"
)

func Generate_Network(tsn int, avb int, hyperperiod int, bandwidth float64, show_topology bool, show_flows bool, show_graphs bool) *Network {
	bytes_rate := 1. / ((bandwidth / 8) * 1e-6)
	Network := &Network{TSN: tsn, AVB: avb, Bytes_Rate: bytes_rate, HyperPeriod: hyperperiod}

	// 1. Generate Topology "generate_topology.go"
	fmt.Println("Generate Topology")
	fmt.Println("----------------------------------------")
	Topology := topology.Generate_Topology(bytes_rate)
	Network.Topology = Topology
	fmt.Println("Topology generation completed.")

	if show_topology {
		Topology.Show_Topology()
	}

	// 2. Generate Flows "generate_flows.go"
	fmt.Println("\nGenerate Flows")
	fmt.Println("----------------------------------------")
	Flow_Set := flow.Generate_Flows(len(Topology.Nodes), tsn, avb, hyperperiod)
	Network.Flow_Set = Flow_Set

	if show_flows {
		Flow_Set.Show_Flows()
		Flow_Set.Show_Flow()
		Flow_Set.Show_Stream()
	}

	// 3. Generate Graphs
	fmt.Println("\nGenerate Graphs")
	fmt.Println("----------------------------------------")
	Graphs_Set := graph.Generate_Graphs(Topology, Flow_Set, bytes_rate)
	Network.Graph_Set = Graphs_Set

	if show_graphs {
		Graphs_Set.Show_Graphs()
	}

	return Network
}
