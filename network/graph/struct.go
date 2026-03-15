package graph

import "src/network/topology"

type Graphs struct {
	TSNGraphs     []*topology.Topology
	AVBGraphs     []*topology.Topology
	CAN2TSNGraphs []*topology.Topology
}

func newGraphs() *Graphs {
	return &Graphs{}
}
