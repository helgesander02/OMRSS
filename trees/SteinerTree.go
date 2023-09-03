package trees

import (
	"encoding/json"
	"src/topology"
)

// Bang Ye Wu, Kun-Mao Chao, "Steiner Minimal Trees"
func SteninerTree(v2v *V2V, topology *topology.Topology, Source int, Destinations []int, cost float64) *Tree {
	t := TopologyDeepCopy(topology)        // Duplicate of Topology
	t.AddN2S2N(Source, Destinations, cost) // Undirected Graph

	var Terminal []int
	Terminal = append(Terminal, Source+1000)
	for _, d := range Destinations {
		Terminal = append(Terminal, d+2000)
	}

	// Determine if there is a vertex in the tree.
	// If not, find all shortest paths between terminals and select the path with the minimum cost.
	// Then, add all the vertices from this shortest path to the tree.
	// If there is a vertex, find all shortest paths between the vertex and terminals.
	// Choose the path with the minimum cost and add all the vertices from this shortest path to the tree.
	tree := &Tree{}
	var used_tmal []int
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
						path := v2vedge.GetPath(tmal)
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
			tree.AddTree(BP, cost)
			used_tmal = append(used_tmal, BP[0])
			used_tmal = append(used_tmal, BP[len(BP)-1])

		} else {
			var BP []int
			for _, terminal := range Terminal {
				if TerminalBeenUsed(used_tmal, terminal) {
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
						path := v2vedge.GetPath(terminal)
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
			tree.AddTree(BP, cost)
			used_tmal = append(used_tmal, BP[0])
		}
	}
	return tree
}

func TerminalBeenUsed(A []int, b int) bool {
	for _, a := range A {
		if a == b {
			return true
		}
	}
	return false
}

func (tree *Tree) AddTree(P []int, cost float64) {
	for l := len(P) - 1; l > 0; l-- {
		node1, b1 := tree.InTree(P[l])
		node2, b2 := tree.InTree(P[l-1])

		if !(b1) {
			tree.Nodes = append(tree.Nodes, node1)
		}
		if !(b2) {
			tree.Nodes = append(tree.Nodes, node2)
		}

		b3 := true
		for _, conn := range node1.Connections {

			if conn.ToNodeID == P[l-1] {
				b3 = false
			}
		}
		if b3 {
			connection1 := &Connection{
				FromNodeID: P[l],
				ToNodeID:   P[l-1],
				Cost:       cost,
			}
			node1.Connections = append(node1.Connections, connection1)
			connection2 := &Connection{
				FromNodeID: P[l-1],
				ToNodeID:   P[l],
				Cost:       cost,
			}
			node2.Connections = append(node2.Connections, connection2)
		}
	}
}

func (tree *Tree) InTree(id int) (*Node, bool) {
	for _, node := range tree.Nodes {
		if node.ID == id {
			return node, true
		}
	}
	return &Node{ID: id}, false
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

func TopologyDeepCopy(t1 *topology.Topology) *topology.Topology {
	if buf, err := json.Marshal(t1); err != nil {
		return nil
	} else {
		t2 := &topology.Topology{}
		if err = json.Unmarshal(buf, t2); err != nil {
			return nil
		}
		return t2
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

func (v2vedge *V2VEdge) GetPath(tmal int) [][]int {
	var path [][]int
	for _, graph := range v2vedge.Graphs {
		if graph.ToVertex == tmal {
			path = graph.Path
		}
	}
	return path
}
