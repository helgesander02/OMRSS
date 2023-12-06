package network

import (
	"fmt"
	"src/network/flow"
	"src/network/graph"
	"src/network/topology"
)

func Generate_Network(topology_name string, tsn int, avb int, hyperperiod int, bandwidth float64, show_topology bool, show_flows bool, show_graphs bool) *Network {
	bw := (bandwidth / 8) * 1e-6
	bytes_rate := 1. / bw      // The number of bytes that can be transmitted in 1us
	bw *= float64(hyperperiod) // The bytes that can be transmitted in 6000us (bytes/us * hyperperiod)
	Network := &Network{TSN: tsn, AVB: avb, Bytes_Rate: bytes_rate, Bandwidth: bw, HyperPeriod: hyperperiod}

	// 1. Generate Topology
	fmt.Println("Generate Topology")
	fmt.Println("----------------------------------------")
	Topology := topology.Generate_Topology(topology_name, bytes_rate)
	Network.Topology = Topology
	fmt.Println("Topology generation completed.")
	fmt.Println()

	if show_topology {
		Topology.Show_Topology()
	}

	// 2. Generate Flows
	fmt.Println("Generate Flows")
	fmt.Println("----------------------------------------")
	Flow_Set := flow.Generate_Flows(len(Topology.Nodes), tsn, avb, hyperperiod)
	Network.Flow_Set = Flow_Set
	Flow_Set.Show_Flows()
	Flow_Set.Show_Flow()
	fmt.Println()

	if show_flows {
		Flow_Set.Show_Stream()
	}

	// 3. Generate Graphs
	fmt.Println("Generate Graphs")
	fmt.Println("----------------------------------------")
	Graphs_Set := graph.Generate_Graphs(Topology, Flow_Set, bytes_rate)
	Network.Graph_Set = Graphs_Set
	fmt.Println("Complete Generating Graphs")
	fmt.Println()

	if show_graphs {
		Graphs_Set.Show_Graphs()
	}

	return Network
}
