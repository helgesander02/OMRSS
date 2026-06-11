package graph

import (
	"src/network/flow"
	"src/network/topology"
	"src/pkg/logger"
)

func (g *Graphs) addGraph(topo *topology.Topology, source int, destinations []int, bytesRate float64) {
	if g.Get(source, destinations) != nil {
		return
	}
	t := topo.TopologyDeepCopy()
	t.AddN2S2N(source, destinations, bytesRate)
	g.Entries = append(g.Entries, &GraphEntry{
		Source:       source,
		Destinations: append([]int{}, destinations...), // defensive copy
		Graph:        t,
	})
}

func GenerateOmacoGraphs(topo *topology.Topology, flows *flow.FlowSet, bytesRate float64) *Graphs {
	graphs := newGraphs()

	for _, f := range flows.TSNFlows {
		graphs.addGraph(topo, f.Source, f.Destinations, bytesRate)
	}
	for _, f := range flows.AVBFlows {
		graphs.addGraph(topo, f.Source, f.Destinations, bytesRate)
	}

	logger.Println(len(graphs.Entries), " unique graphs generated (OMACO).")
	return graphs
}

func GenerateOsroGraphs(topo *topology.Topology, flows *flow.FlowSet, bytesRate float64) *Graphs {
	graphs := newGraphs()

	for _, f := range flows.TSNFlows {
		graphs.addGraph(topo, f.Source, f.Destinations, bytesRate)
	}
	for _, f := range flows.AVBFlows {
		graphs.addGraph(topo, f.Source, f.Destinations, bytesRate)
	}
	for _, method := range flows.EncapsulateMethod {
		for _, f := range method.CAN2TTFlows {
			graphs.addGraph(topo, f.Source, f.Destinations, bytesRate)
		}
	}

	logger.Println(len(graphs.Entries), " unique graphs generated (OSRO).")
	return graphs
}

func (g *Graphs) Get(source int, destinations []int) *topology.Topology {
	for _, entry := range g.Entries {
		if entry.Source == source && sameIntSlice(entry.Destinations, destinations) {
			return entry.Graph
		}
	}
	return nil
}

func sameIntSlice(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
