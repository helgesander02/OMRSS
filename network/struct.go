package network

import (
	"src/network/flow"
	"src/network/graph"
	"src/network/topology"
)

type Network struct {
	Topology *topology.Topology
	FlowSet  *flow.FlowSet
	GraphSet *graph.Graphs
}

func newNetwork() *Network {
	return &Network{}
}
