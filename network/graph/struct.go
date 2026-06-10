package graph

import "src/network/topology"

type Graphs struct {
	Entries []*GraphEntry
}

type GraphEntry struct {
	Source       int
	Destinations []int
	Graph        *topology.Topology
}

func newGraphs() *Graphs {
	return &Graphs{}
}
