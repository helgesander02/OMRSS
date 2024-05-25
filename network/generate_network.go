package network

import (
	"fmt"
	"src/network/flow"
	"src/network/graph"
	"src/network/topology"
)

func Generate_Network(topology_name string, tsn int, avb int, hyperperiod int, bandwidth float64) *Network {
	// 1. Define network parameters
	Network := new_Network(topology_name, tsn, avb, hyperperiod, bandwidth)

	// 2. Generate topology
	fmt.Println("Generate Topology")
	fmt.Println("----------------------------------------")
	Network.Topology = topology.Generate_Topology(Network.TopologyName, Network.BytesRate)
	fmt.Println("Complete Generating Topology.")
	fmt.Println()

	// 3. Generate flows
	fmt.Println("Generate Flows")
	fmt.Println("----------------------------------------")
	Network.Flow_Set = flow.Generate_Flows(len(Network.Topology.Nodes), Network.BG_TSN, Network.BG_AVB, Network.Input_TSN, Network.Input_AVB, Network.HyperPeriod)
	fmt.Println("Complete Generating Flows.")
	fmt.Println()

	// 4. Simulating graphs using flows in topology
	fmt.Println("Simulating Graphs")
	fmt.Println("----------------------------------------")
	Network.Graph_Set = graph.Generate_Graphs(Network.Topology, Network.Flow_Set, Network.BytesRate)
	fmt.Println("Complete Simulating Graphs.")
	fmt.Println()

	return Network
}

func new_Network(topology_name string, tsn int, avb int, hyperperiod int, bandwidth float64) *Network {
	bw := (bandwidth / 8) * 1e-6 // bytes/us ==> 125 bytes
	bytes_rate := 1. / bw        // The number of bytes that can be transmitted in 1us ==> 1/125
	bw *= float64(hyperperiod)   // The bytes that can be transmitted in 6000us (bytes/us * hyperperiod) ==> 750000 bytes

	Network := &Network{
		HyperPeriod:  hyperperiod,
		BytesRate:    bytes_rate,
		Bandwidth:    bw,
		TopologyName: topology_name,
		BG_TSN:       35,
		BG_AVB:       15,
		Input_TSN:    tsn,
		Input_AVB:    avb,
	}

	fmt.Println("Define network parameters")
	fmt.Println("----------------------------------------")
	fmt.Printf("HyperPeriod: %d us \n", Network.HyperPeriod)
	fmt.Printf("Bandwidth:  %f bytes/6000us \n", Network.Bandwidth)
	fmt.Printf("BytesRate:  %f BPS \n", Network.BytesRate)
	fmt.Printf("Topology file name: %s \n", Network.TopologyName)
	fmt.Printf("TSN flow: %d, AVB flow: %d \n", Network.Input_TSN+Network.BG_TSN, Network.Input_AVB+Network.BG_AVB)
	fmt.Println()

	return Network
}
