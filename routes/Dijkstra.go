package routes

import (
	"math"
	"sort"
)

// Sean Chou, "Dijkstra’s Algorithm"
func Dijkstra(graph *Graph, strat int, terminal int) *Graph {
	inf := math.MaxInt8

	for _, vertex := range graph.Vertexs {
		vertex.Visited = false
		if vertex.ID == strat {
			vertex.Cost = 0
		} else {
			vertex.Cost = inf
		}
		vertex.Path = -1
	}

	// depth-first-search return all P(shortest path)
	graph.GetShortestPath(strat, terminal)

	sort.Slice(graph.Path, func(p, q int) bool {
		return len(graph.Path[p]) < len(graph.Path[q])
	})

	return graph
}

func (graph *Graph) GetShortestPath(strat int, terminal int) {
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
		graph.GetShortestPath(edge.End, terminal)
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
			if LoopCompare(P, path) {
				InPath = false
				break
			}
		}
		if InPath {
			graph.Path = append(graph.Path, path)
		}
	}
}

func LoopCompare(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func (graph *Graph) FindVertex(id int) *Vertex {
	for _, vertex := range graph.Vertexs {
		if vertex.ID == id {
			return vertex
		}
	}
	return &Vertex{}
}
