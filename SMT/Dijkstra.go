package SMT

import (
	"math"
)

func Dijkstra(graphs *Graphs, strat int, terminal int) *Graphs {
	inf := math.MaxInt8

	for _, graph := range graphs.Graphs {
		graph.Visited = false
		if graph.ID == strat {
			graph.Cost = 0
		} else {
			graph.Cost = inf
		}	
		graph.Path = -1
	}

	graphs.GetShortestPath(strat, terminal)

	return graphs
} 

func (graphs *Graphs) GetShortestPath(strat int, terminal int) {
	graph := graphs.findGraph(strat)
	graph.Visited = true

	for _, edge := range graph.Edges {
		nextgraph := graphs.findGraph(edge.End)
		if nextgraph.Visited {continue}
		if nextgraph.Cost > graph.Cost + edge.Cost {
			nextgraph.Path = graph.ID
			nextgraph.Cost = graph.Cost + edge.Cost		
		}
		
		graphs.GetShortestPath(edge.End, terminal)
	}

	graph.Visited = false
}
