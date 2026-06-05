package network

import (
	"src/network/flow"
	"src/network/graph"
	"src/network/topology"
	"src/pkg/config"
	"src/pkg/logger"
)

func GenerateNetwork(cfg *config.Config) *Network {
	networkInstance := newNetwork()

	switch cfg.Algorithm.Name {
	case "omaco":
		networkInstance.generateOmacoNetwork(cfg)
	case "osro":
		networkInstance.generateOsroNetwork(cfg)
	default:
		logger.Printf("Unknown algorithm: %s\n", cfg.Algorithm.Name)
	}

	return networkInstance
}

func (network *Network) generateOmacoNetwork(cfg *config.Config) {
	// Generate topology
	logger.Println("Generate Topology")
	logger.Println("----------------------------------------")
	network.Topology = topology.GenerateTopology(
		cfg.Network.Topology,
		cfg.Network.ByteRate,
	)
	logger.Println("Complete Generating Topology.")
	logger.Println()

	// Generate flows
	logger.Println("Generate Flows")
	logger.Println("----------------------------------------")
	network.FlowSet = flow.GenerateOmacoFlows(
		cfg.Network.Flows,
		cfg.Network.Hyperperiod,
		len(network.Topology.Nodes),
	)
	logger.Println("Complete Generating Flows.")
	logger.Println()

	// Simulating graphs using flows in topology
	logger.Println("Simulating Graphs")
	logger.Println("----------------------------------------")
	network.GraphSet = graph.GenerateOmacoGraphs(
		network.Topology,
		network.FlowSet,
		cfg.Network.ByteRate,
	)
	logger.Println("Complete Simulating Graphs.")
	logger.Println()
}

func (network *Network) generateOsroNetwork(cfg *config.Config) {
	// Generate topology
	logger.Println("Generate Topology")
	logger.Println("----------------------------------------")
	network.Topology = topology.GenerateTopology(
		cfg.Network.Topology,
		cfg.Network.ByteRate,
	)
	logger.Println("Complete Generating Topology.")
	logger.Println()

	// Select CAN node
	canNodeSet := network.Topology.SelectCANNodes(cfg.Network.Flows.CAN.Nodes)
	logger.Printf("CAN nodes: %v", canNodeSet)
	logger.Println()

	// Generate flows
	logger.Println("Generate Flows")
	logger.Println("----------------------------------------")
	network.FlowSet = flow.GenerateOsroFlows(
		cfg.Network.Flows,
		cfg.Network.Hyperperiod,
		len(network.Topology.Nodes),
		canNodeSet,
	)
	logger.Println("Complete Generating Flows.")
	logger.Println()

	// Simulating graphs using flows in topology
	logger.Println("Simulating Graphs")
	logger.Println("----------------------------------------")
	network.GraphSet = graph.GenerateOsroGraphs(
		network.Topology,
		network.FlowSet,
		cfg.Network.ByteRate,
	)
	logger.Println("Complete Simulating Graphs.")
	logger.Println()
}
