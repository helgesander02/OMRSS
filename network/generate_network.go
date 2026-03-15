package network

import (
	"fmt"
	"src/internal/config"
	"src/network/flow"
	"src/network/graph"
	"src/network/topology"
)

func GenerateNetwork(cfg *config.Config) *Network {
	networkInstance := newNetwork()

	switch cfg.Algorithm.Name {
	case "omaco":
		networkInstance.generateOmacoNetwork(cfg)
	case "osro":
		networkInstance.generateOsroNetwork(cfg)
	default:
		fmt.Printf("Unknown algorithm: %s\n", cfg.Algorithm.Name)
	}

	return networkInstance
}

func (network *Network) generateOmacoNetwork(cfg *config.Config) {
	// Generate topology
	fmt.Println("Generate Topology")
	fmt.Println("----------------------------------------")
	network.Topology = topology.GenerateTopology(
		cfg.Network.Topology,
		cfg.Network.ByteRate,
	)
	fmt.Println("Complete Generating Topology.")
	fmt.Println()

	// Generate flows
	fmt.Println("Generate Flows")
	fmt.Println("----------------------------------------")
	network.FlowSet = flow.GenerateOmacoFlows(
		len(network.Topology.Nodes),
		cfg.Network.Flows.TSN.Background,
		cfg.Network.Flows.AVB.Background,
		cfg.Network.Flows.TSN.Input,
		cfg.Network.Flows.AVB.Input,
		cfg.Network.Hyperperiod,
	)
	fmt.Println("Complete Generating Flows.")
	fmt.Println()

	// Simulating graphs using flows in topology
	fmt.Println("Simulating Graphs")
	fmt.Println("----------------------------------------")
	network.GraphSet = graph.GenerateOMACOGraphs(
		network.Topology,
		network.FlowSet,
		cfg.Network.ByteRate,
	)
	fmt.Println("Complete Simulating Graphs.")
	fmt.Println()
}

func (network *Network) generateOsroNetwork(cfg *config.Config) {
	// Generate topology
	fmt.Println("Generate Topology")
	fmt.Println("----------------------------------------")
	network.Topology = topology.GenerateTopology(
		cfg.Network.Topology,
		cfg.Network.ByteRate,
	)
	fmt.Println("Complete Generating Topology.")
	fmt.Println()

	// Select CAN node
	canNodeSet := network.Topology.SelectCANNodes()
	fmt.Printf("CAN nodes: %v", canNodeSet)
	fmt.Println()

	// Generate flows
	fmt.Println("Generate Flows")
	fmt.Println("----------------------------------------")
	network.FlowSet = flow.GenerateOsroFlows(
		len(network.Topology.Nodes),
		cfg.Network.Flows.TSN.Background,
		cfg.Network.Flows.AVB.Background,
		cfg.Network.Flows.TSN.Input,
		cfg.Network.Flows.AVB.Input,
		cfg.Network.Hyperperiod,
		canNodeSet,
		cfg.Network.Flows.CAN.Important,
		cfg.Network.Flows.CAN.Unimportant,
	)
	fmt.Println("Complete Generating Flows.")
	fmt.Println()

	// Simulating graphs using flows in topology
	fmt.Println("Simulating Graphs")
	fmt.Println("----------------------------------------")
	network.GraphSet = graph.GenerateOSROGraphs(
		network.Topology,
		network.FlowSet,
		cfg.Network.ByteRate,
	)
	fmt.Println("Complete Simulating Graphs.")
	fmt.Println()
}
