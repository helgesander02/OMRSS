package network

import (
	"src/network/flow"
	"src/network/graph"
	"src/network/topology"
)

type Network struct {
	TSN         int
	AVB         int
	BytesRate   float64
	Bandwidth   float64
	HyperPeriod int
	Topology    *topology.Topology
	Flow_Set    *flow.Flows
	Graph_Set   *graph.Graphs
}
