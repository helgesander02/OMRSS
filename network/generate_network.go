package network

import (
	"fmt"
	"src/network/flow"
	"src/network/graph"
	"src/network/topology"
)

func (network *Network) Generate_Network() {
	// 2. Generate topology
	fmt.Println("Generate Topology")
	fmt.Println("----------------------------------------")
	network.Topology = topology.Generate_Topology(network.TopologyName, network.BytesRate)
	fmt.Println("Complete Generating Topology.")
	fmt.Println()

	// 3. Generate flows
	fmt.Println("Generate Flows")
	fmt.Println("----------------------------------------")
	network.Flow_Set = flow.Generate_Flows(len(network.Topology.Nodes), network.BG_TSN, network.BG_AVB, network.Input_TSN, network.Input_AVB, network.HyperPeriod)
	fmt.Println("Complete Generating Flows.")
	fmt.Println()

	// 4. Simulating graphs using flows in topology
	fmt.Println("Simulating Graphs")
	fmt.Println("----------------------------------------")
	network.Graph_Set = graph.Generate_Graphs(network.Topology, network.Flow_Set, network.BytesRate)
	fmt.Println("Complete Simulating Graphs.")
	fmt.Println()
}

func (network *OSRO_Network) Generate_Network() {
	// 2. Generate topology
	fmt.Println("Generate Topology")
	fmt.Println("----------------------------------------")
	network.Topology = topology.Generate_Topology(network.TopologyName, network.BytesRate)
	fmt.Println("Complete Generating Topology.")
	fmt.Println()

	// select CAN node
	CAN_Node_Set := network.Topology.Select_CAN_Node_Set()
	fmt.Printf("CAN nodes: %v", CAN_Node_Set)
	fmt.Println()

}

//func (network *Network) Generate_Network() *Network {
//
//}
