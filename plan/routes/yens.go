package routes

import (
	"math"
	"sort"
	"src/network/topology"
)

// YenKPaths returns K shortest simple paths, sorted by weight
// Uses Yen's algorithm for K-shortest paths
func YenKPaths(graph *Graph, src, dst, K int) [][]int {
	A := [][]int{}

	// Find the first shortest path
	first := oneShortestPath(graph, src, dst)
	if first == nil {
		return nil
	}
	A = append(A, first)
	B := [][]int{}

	// Find K-1 more paths
	for ki := 1; ki < K; ki++ {
		last := A[ki-1]

		// For each node in the previous shortest path (except the last one)
		for i := 0; i < len(last)-1; i++ {
			spurNode := last[i]
			rootPath := last[:i+1]

			// Create a copy of the graph
			graphCopy := graph.Clone()

			// Remove edges that share the same root path
			for _, path := range A {
				if len(path) > i && equalPath(rootPath, path[:i+1]) {
					graphCopy.RemoveEdge(path[i], path[i+1])
				}
			}

			// Remove nodes in root path (except spur node)
			for _, nodeID := range rootPath[:len(rootPath)-1] {
				graphCopy.RemoveVertex(nodeID)
			}

			// Find spur path
			spurPath := oneShortestPath(graphCopy, spurNode, dst)
			if spurPath == nil {
				continue
			}

			// Combine root path and spur path
			fullPath := append(append([]int{}, rootPath[:len(rootPath)-1]...), spurPath...)
			B = append(B, fullPath)
		}

		if len(B) == 0 {
			break
		}

		// Sort B by path length
		sort.Slice(B, func(i, j int) bool {
			return len(B[i]) < len(B[j])
		})

		A = append(A, B[0])
		B = B[1:]
	}

	return A
}

// oneShortestPath finds single shortest path using Dijkstra
func oneShortestPath(graph *Graph, src, dst int) []int {
	const inf = math.MaxInt32
	dist := make(map[int]int)
	prev := make(map[int]int)

	// Initialize distances
	for _, vertex := range graph.Vertexs {
		dist[vertex.ID] = inf
	}
	dist[src] = 0
	visited := make(map[int]bool)

	// Dijkstra's algorithm
	for len(visited) < len(graph.Vertexs) {
		// Find unvisited vertex with minimum distance
		minID, minDist := -1, inf
		for _, vertex := range graph.Vertexs {
			if !visited[vertex.ID] && dist[vertex.ID] < minDist {
				minDist = dist[vertex.ID]
				minID = vertex.ID
			}
		}

		if minID == -1 || minID == dst {
			break
		}

		visited[minID] = true

		// Update distances to neighbors
		vertex := graph.FindVertex(minID)
		for _, edge := range vertex.Edges {
			if visited[edge.End] {
				continue
			}
			alt := dist[minID] + 1
			if alt < dist[edge.End] {
				dist[edge.End] = alt
				prev[edge.End] = minID
			}
		}
	}

	// Check if destination is reachable
	if dist[dst] == inf {
		return nil
	}

	// Reconstruct path
	path := []int{dst}
	for cur := dst; cur != src; {
		cur = prev[cur]
		path = append(path, cur)
	}

	// Reverse path
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path
}

// Clone creates a deep copy of the graph for Yen's algorithm
func (graph *Graph) Clone() *Graph {
	clone := &Graph{
		Vertexs:  make([]*Vertex, 0, len(graph.Vertexs)),
		ToVertex: graph.ToVertex,
		Path:     make([][]int, len(graph.Path)),
	}

	// Copy paths
	for i, p := range graph.Path {
		clone.Path[i] = append([]int{}, p...)
	}

	// Clone vertices
	for _, v := range graph.Vertexs {
		newVertex := &Vertex{
			ID:      v.ID,
			Visited: v.Visited,
			Cost:    v.Cost,
			Path:    v.Path,
			Edges:   make([]*Edge, 0, len(v.Edges)),
		}

		// Clone edges
		for _, e := range v.Edges {
			newVertex.Edges = append(newVertex.Edges, &Edge{
				Strat: e.Strat,
				End:   e.End,
			})
		}

		clone.Vertexs = append(clone.Vertexs, newVertex)
	}

	return clone
}

// RemoveEdge removes a directed edge from the graph
func (graph *Graph) RemoveEdge(from, to int) {
	vertex := graph.FindVertex(from)
	if vertex == nil {
		return
	}

	newEdges := make([]*Edge, 0, len(vertex.Edges))
	for _, edge := range vertex.Edges {
		if !(edge.Strat == from && edge.End == to) {
			newEdges = append(newEdges, edge)
		}
	}
	vertex.Edges = newEdges
}

// RemoveVertex removes a vertex and all its incident edges from the graph
func (graph *Graph) RemoveVertex(id int) {
	// Remove the vertex
	newVertexs := make([]*Vertex, 0, len(graph.Vertexs))
	for _, v := range graph.Vertexs {
		if v.ID != id {
			newVertexs = append(newVertexs, v)
		}
	}
	graph.Vertexs = newVertexs

	// Remove all edges pointing to this vertex
	for _, vertex := range graph.Vertexs {
		newEdges := make([]*Edge, 0, len(vertex.Edges))
		for _, edge := range vertex.Edges {
			if edge.Strat != id && edge.End != id {
				newEdges = append(newEdges, edge)
			}
		}
		vertex.Edges = newEdges
	}
}

// equalPath checks if two paths are equal
func equalPath(a, b []int) bool {
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

// BuildGraphFromTopology converts a topology to a graph structure
func BuildGraphFromTopology(topo *topology.Topology) *Graph {
	graph := &Graph{
		Vertexs: make([]*Vertex, 0),
	}

	// Add all nodes from topology
	addNodes := func(nodes []*topology.Node) {
		for _, node := range nodes {
			vertex := &Vertex{
				ID:    node.ID,
				Edges: make([]*Edge, 0),
			}

			// Add edges. Hop cost is implicit (1) — see routes.Edge doc.
			for _, conn := range node.Links {
				vertex.Edges = append(vertex.Edges, &Edge{
					Strat: conn.FromNodeID,
					End:   conn.ToNodeID,
				})
			}

			graph.Vertexs = append(graph.Vertexs, vertex)
		}
	}

	addNodes(topo.Talker)
	addNodes(topo.Switch)
	addNodes(topo.Listener)

	return graph
}

// ConvertIDsToTree converts a path (node IDs) to a Tree structure
func ConvertIDsToTree(ids []int, topo *topology.Topology, cost float64) *Route {
	tree := newRoute()

	for i := 0; i < len(ids)-1; i++ {
		node1, found1 := tree.CheckNodeByID(ids[i])
		node2, found2 := tree.CheckNodeByID(ids[i+1])

		if !found1 {
			tree.Nodes = append(tree.Nodes, node1)
		}
		if !found2 {
			tree.Nodes = append(tree.Nodes, node2)
		}

		// Check if connection already exists
		hasConnection := false
		for _, conn := range node1.Connections {
			if conn.ToNodeID == ids[i+1] {
				hasConnection = true
				break
			}
		}

		if !hasConnection {
			conn1 := newConnection(ids[i], ids[i+1], cost)
			node1.Connections = append(node1.Connections, conn1)

			conn2 := newConnection(ids[i+1], ids[i], cost)
			node2.Connections = append(node2.Connections, conn2)

			tree.Weight++
		}
	}

	return tree
}
