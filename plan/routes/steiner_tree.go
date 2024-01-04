package routes

import (
	"src/network/topology"
)

// Bang Ye Wu, Kun-Mao Chao, "Steiner Minimal Trees"
func SteninerTree(v2v *V2V, t *topology.Topology, Source int, Destinations []int, cost float64) *Tree {
	var Terminal []int
	Terminal = append(Terminal, Source)
	Terminal = append(Terminal, Destinations...)

	// Determine if there is a vertex in the tree.
	// If not, find all shortest paths between terminals and select the path with the minimum cost.
	// Then, add all the vertices from this shortest path to the tree.
	// If there is a vertex, find all shortest paths between the vertex and terminals.
	// Choose the path with the minimum cost and add all the vertices from this shortest path to the tree.
	var (
		tree      *Tree = &Tree{}
		used_tmal []int
	)

	for len(used_tmal) != len(Terminal) {
		if len(tree.Nodes) == 0 {
			// Find the set of all shortest paths from Vertex to Vertexs
			var BP []int
			for _, terminal := range Terminal {
				v2vedge, firstget := v2v.GetV2VEdge(terminal) //Define a function to get the v2vedge form v2v
				for _, tmal := range Terminal {
					if terminal == tmal {
						continue

					} else if !(v2vedge.InV2VEdge(tmal)) {
						graph := GetGarph(t)
						graph.ToVertex = tmal
						graph = Dijkstra(graph, terminal, tmal)
						v2vedge.Graphs = append(v2vedge.Graphs, graph)
						// Choose the path that is the shortest and least expensive from Vertex to Vertexs in BP
						if len(BP) == 0 {
							BP = graph.Path[0]
						} else {
							if len(BP) > len(graph.Path[0]) {
								BP = graph.Path[0]
							}
						}

					} else {
						path := v2vedge.GetV2VPath(tmal)
						if len(BP) == 0 {
							BP = path[0]
						} else {
							if len(BP) > len(path[0]) {
								BP = path[0]
							}
						}
					}
				}
				if firstget {
					v2v.V2VEdges = append(v2v.V2VEdges, v2vedge)
				}
			}

			// Add the shortest path with the minimum cost to the tree in sequence
			tree.IntoTree(BP, cost)
			used_tmal = append(used_tmal, BP[0])
			used_tmal = append(used_tmal, BP[len(BP)-1])

		} else {
			var BP []int
			for _, terminal := range Terminal {
				if terminal_been_used(used_tmal, terminal) {
					continue
				}
				// Find the set of all shortest paths from Vertex to Vertexs
				for _, node := range tree.Nodes {
					v2vedge, firstget := v2v.GetV2VEdge(node.ID)
					if terminal == node.ID {
						continue

					} else if !(v2vedge.InV2VEdge(node.ID)) {
						graph := GetGarph(t)
						graph.ToVertex = terminal
						graph = Dijkstra(graph, node.ID, terminal)
						v2vedge.Graphs = append(v2vedge.Graphs, graph)
						// Choose the path that is the shortest and least expensive from Vertex to Vertexs in BP
						if len(BP) == 0 {
							BP = graph.Path[0]
						} else {
							if len(BP) > len(graph.Path[0]) {
								BP = graph.Path[0]
							}
						}

					} else {
						path := v2vedge.GetV2VPath(terminal)
						if len(BP) == 0 {
							BP = path[0]
						} else {
							if len(BP) > len(path[0]) {
								BP = path[0]
							}
						}
					}

					if firstget {
						v2v.V2VEdges = append(v2v.V2VEdges, v2vedge)
					}
				}
			}
			// Add the shortest path with the minimum cost to the tree in sequence
			tree.IntoTree(BP, cost)
			used_tmal = append(used_tmal, BP[0])
		}
	}
	return tree
}

func terminal_been_used(A []int, b int) bool {
	for _, a := range A {
		if a == b {
			return true
		}
	}
	return false
}

func GetGarph(topology *topology.Topology) *Graph {
	graph := &Graph{}
	//Talker
	for _, t := range topology.Talker {
		gt := &Vertex{}
		gt.ID = t.ID
		gt.AddEdge(t.Connections)
		graph.Vertexs = append(graph.Vertexs, gt)
	}

	//Switch
	for _, s := range topology.Switch {
		gs := &Vertex{}
		gs.ID = s.ID
		gs.AddEdge(s.Connections)
		graph.Vertexs = append(graph.Vertexs, gs)
	}

	//Listener
	for _, l := range topology.Listener {
		gl := &Vertex{}
		gl.ID = l.ID
		gl.AddEdge(l.Connections)
		graph.Vertexs = append(graph.Vertexs, gl)
	}

	return graph
}

func (vertex *Vertex) AddEdge(connections []*topology.Connection) {
	for _, c := range connections {
		edge := &Edge{
			Strat: c.FromNodeID,
			End:   c.ToNodeID,
			Cost:  1,
		}
		vertex.Edges = append(vertex.Edges, edge)
	}
}

func (v2v *V2V) GetV2VEdge(terminal int) (*V2VEdge, bool) {
	for _, edge := range v2v.V2VEdges {
		if edge.FromVertex == terminal {
			return edge, false
		}
	}
	return &V2VEdge{FromVertex: terminal}, true
}

func (v2vedge *V2VEdge) InV2VEdge(tmal int) bool {
	for _, graph := range v2vedge.Graphs {
		if graph.ToVertex == tmal {
			return true
		}
	}
	return false
}

func (v2vedge *V2VEdge) GetV2VPath(tmal int) [][]int {
	var path [][]int
	for _, graph := range v2vedge.Graphs {
		if graph.ToVertex == tmal {
			path = graph.Path
		}
	}
	return path
}
