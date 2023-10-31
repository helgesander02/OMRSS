package routes

import (
	"math"
	"sort"
)

// Dijkstra’s Algorithm
// https://medium.com/%E6%8A%80%E8%A1%93%E7%AD%86%E8%A8%98/%E5%9F%BA%E7%A4%8E%E6%BC%94%E7%AE%97%E6%B3%95%E7%B3%BB%E5%88%97-graph-%E8%B3%87%E6%96%99%E7%B5%90%E6%A7%8B%E8%88%87dijkstras-algorithm-6134f62c1fc2
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

	GetShortestPath(graph, strat, terminal)

	sort.Slice(graph.Path, func(p, q int) bool {
		return len(graph.Path[p]) < len(graph.Path[q])
	})

	return graph
}

// DFS graph return all P(shortest path)
func GetShortestPath(graph *Graph, strat int, terminal int) {
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
		GetShortestPath(graph, edge.End, terminal)
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
			if loopcompare(P, path) {
				InPath = false
				break
			}
		}
		if InPath {
			graph.Path = append(graph.Path, path)
		}
	}
}

func loopcompare(a, b []int) bool {
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
