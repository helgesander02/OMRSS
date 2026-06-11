package routes

import (
	"math"
	"sort"
)

// Dijkstra’s Algorithm.
//
// NOTE: this is misnamed — it is actually a DFS-with-backtracking that
// enumerates every simple path from start to terminal and stores them
// all in graph.Path. The Steiner tree construction relies on having the
// full set of paths available; ShortestPath only consumes graph.Path[0].
// See routes/struct.go for the Graph cache layout.
func Dijkstra(graph *Graph, strat int, terminal int) *Graph {
	// Cost upper bound. The previous implementation used math.MaxInt8
	// (127), which silently truncated any path longer than 127 hops.
	inf := math.MaxInt

	for _, vertex := range graph.Vertexs {
		vertex.Visited = false
		if vertex.ID == strat {
			vertex.Cost = 0
		} else {
			vertex.Cost = inf
		}
		vertex.Path = -1
	}

	getDijkstraShortestPath(graph, strat, terminal)

	sort.Slice(graph.Path, func(p, q int) bool {
		return len(graph.Path[p]) < len(graph.Path[q])
	})

	return graph
}

// DFS graph return all P(shortest path)
func getDijkstraShortestPath(graph *Graph, strat int, terminal int) {
	vertex := graph.FindVertex(strat)
	vertex.Visited = true

	for _, edge := range vertex.Edges {
		nextvertex := graph.FindVertex(edge.End)
		if nextvertex.Visited {
			continue
		}
		if nextvertex.Cost >= vertex.Cost+edge.Cost {
			nextvertex.Path = vertex.ID
			nextvertex.Cost = vertex.Cost + edge.Cost

			// Store all the paths from the vertex 'start' to the vertex 'end'
			if nextvertex.ID == terminal {
				graph.AddPath(terminal)
			}
		}
		getDijkstraShortestPath(graph, edge.End, terminal)
	}

	vertex.Visited = false
}

func (graph *Graph) AddPath(terminal int) {
	var path []int
	location := terminal
	vertex := graph.FindVertex(location)
	path = append(path, terminal)

	for vertex.Path != -1 {
		path = append(path, vertex.Path)

		location = vertex.Path
		vertex = graph.FindVertex(location)
	}

	if len(graph.Path) == 0 {
		graph.Path = append(graph.Path, path)
	} else {
		InPath := true
		for _, P := range graph.Path {
			// If P is already present in the path, do not include it
			if loopCompareComplex(P, path) {
				InPath = false
				break
			}
		}
		if InPath {
			graph.Path = append(graph.Path, path)
		}
	}
}

func (graph *Graph) FindVertex(id int) *Vertex {
	for _, vertex := range graph.Vertexs {
		if vertex.ID == id {
			return vertex
		}
	}
	return &Vertex{}
}
