package network

import (
	"src/network/flow"
	"src/network/graph"
	"src/network/topology"
)

type Network struct {
	HyperPeriod  int
	BytesRate    float64
	Bandwidth    float64
	TopologyName string
	BG_TSN       int
	BG_AVB       int
	Input_TSN    int
	Input_AVB    int
	Topology     *topology.Topology
	Flow_Set     *flow.Flows
	Graph_Set    *graph.Graphs
}
