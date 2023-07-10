package SMT

import (
	"math"
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

	return graph
} 

func (graph *Graph) GetShortestPath(strat int, terminal int) {
	vertex := graph.FindGraph(strat)
	vertex.Visited = true

	for _, edge := range vertex.Edges {
		nextgraph := graph.FindGraph(edge.End)
		if nextgraph.Visited {continue}
		if nextgraph.Cost >= vertex.Cost + edge.Cost {
			nextgraph.Path = vertex.ID
			nextgraph.Cost = vertex.Cost + edge.Cost
			
			// Find all paths to the terminal and filter out the set of shortest paths
			if nextgraph.ID == terminal {
				if graph.Count == 0 {
					graph.Count = nextgraph.Cost
					graph.AddPath(terminal)

				} else if graph.Count > nextgraph.Cost {
					graph.Count = nextgraph.Cost
					graph.Path = graph.Path[:0]
					graph.AddPath(terminal)

				} else {
					graph.AddPath(terminal)
				}	
			}			
		}	
		graph.GetShortestPath(edge.End, terminal)
	}

	vertex.Visited = false
}

func (graph *Graph) AddPath(terminal int) {
	var path []int
	location := terminal
	vertex := graph.FindGraph(location)
	path = append(path, terminal)

	for vertex.Path != -1 {
		path = append(path, vertex.Path)

		location = vertex.Path
		vertex = graph.FindGraph(location)
	}

	if len(graph.Path) == 0 {
		graph.Path = append(graph.Path, path)
	} else {
		InPath:= true
		for _, P := range graph.Path {
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
    for i, v := range a {
        if v != b[i] {
            return false
        }
    }
    return true
}

func (graph *Graph) FindGraph(id int) *Vertex {
    for _, vertex := range graph.Vertexs {
        if vertex.ID == id {
            return vertex
        }
    }
    return &Vertex{}
}