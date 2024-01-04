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
	TSN          int
	AVB          int
	Topology     *topology.Topology
	Flow_Set     *flow.Flows
	Graph_Set    *graph.Graphs
}
