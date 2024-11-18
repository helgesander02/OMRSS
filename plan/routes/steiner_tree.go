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
		tree      *Tree = new_Tree()
		used_tmal []int
	)

	for len(used_tmal) != len(Terminal) {
		if len(tree.Nodes) == 0 {
			// Find the set of all shortest paths from Vertex to Vertexs
			BP := make(map[int][][]int)
			for _, terminal1 := range Terminal {
				v2vedge, firstget := v2v.GetV2VEdge(terminal1) //Define a function to get the v2vedge form v2v
				for _, terminal2 := range Terminal {
					if terminal1 != terminal2 {
						saveShortestPathsToBP(BP, v2vedge, terminal2, terminal1, t)
					}
				}
				if firstget {
					v2v.V2VEdges = append(v2v.V2VEdges, v2vedge)
				}
			}

			// Add the shortest path with the minimum cost to the tree in sequence
			bp := chooseBestPath(BP, tree, cost)
			used_tmal = append(used_tmal, bp[0])
			used_tmal = append(used_tmal, bp[len(bp)-1])

		} else {
			BP := make(map[int][][]int)
			for _, terminal1 := range Terminal {
				if terminal_been_used(used_tmal, terminal1) {
					continue
				}
				// Find the set of all shortest paths from Vertex to Vertexs
				for _, node := range tree.Nodes {
					v2vedge, firstget := v2v.GetV2VEdge(node.ID)
					if terminal1 != node.ID {
						saveShortestPathsToBP(BP, v2vedge, terminal1, node.ID, t)
					}
					if firstget {
						v2v.V2VEdges = append(v2v.V2VEdges, v2vedge)
					}
				}
			}
			// Add the shortest path with the minimum cost to the tree in sequence
			bp := chooseBestPath(BP, tree, cost)
			used_tmal = append(used_tmal, bp[0])
		}
	}
	return tree
}

func chooseBestPath(BP map[int][][]int, tree *Tree, cost float64) []int {
	var (
		minPaths [][]int
		minPath  []int
	)
	minW := 0

	// Choosing shortest paths
	for key, val := range BP {
		if minW == 0 {
			minW = key
			minPaths = val
		}
		if minW > key {
			minW = key
			minPaths = val
		}
	}

	// Choosing to join the shortest path with the minimum weight in the tree
	if len(tree.Nodes) == 0 {
		tree.IntoTree(minPaths[0], cost)
		return minPaths[0]

	} else {
		for _, path := range minPaths {
			tree_copy := tree.TreeDeepCopy()
			tree_copy.IntoTree(path, cost)
			if len(minPath) == 0 {
				minPath = path
				minW = len(tree_copy.Nodes)
			} else {
				if minW > len(tree_copy.Nodes) {
					minPath = path
					minW = len(tree_copy.Nodes)
				}
			}
		}
		tree.IntoTree(minPath, cost)
		return minPath
	}
}

func saveShortestPathsToBP(BP map[int][][]int, v2vedge *V2VEdge, vertex1 int, vertex2 int, t *topology.Topology) {
	// Check if this path has already been taken
	if !(v2vedge.InV2VEdge(vertex2)) {
		graph := GetGarph(t)
		graph.ToVertex = vertex1
		graph = Dijkstra(graph, vertex2, vertex1)
		v2vedge.Graphs = append(v2vedge.Graphs, graph)
		addShortestPaths(BP, graph.Path)

	} else {
		path := v2vedge.GetV2VPath(vertex1)
		addShortestPaths(BP, path)

	}
}

func addShortestPaths(BP map[int][][]int, paths [][]int) {
	for _, path := range paths {
		BP[len(path)] = append(BP[len(path)], path)
	}
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

func (v2vedge *V2VEdge) InV2VEdge(terminal int) bool {
	for _, graph := range v2vedge.Graphs {
		if graph.ToVertex == terminal {
			return true
		}
	}
	return false
}

func (v2vedge *V2VEdge) GetV2VPath(terminal int) [][]int {
	var path [][]int
	for _, graph := range v2vedge.Graphs {
		if graph.ToVertex == terminal {
			path = graph.Path
		}
	}
	return path
}
