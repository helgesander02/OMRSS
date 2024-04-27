package routes

import (
	"src/network/topology"
)

// QINGHAN YU et al., "Online Scheduling for Dynamic VM Migration in Multicast Time-Sensitive Networks"
func DistanceTree(Graph *topology.Topology, Source int, Destinations []int, cost float64) *Tree {
	MDTC := new_Tree()
	for _, destination := range Destinations {
		MDPC := DP_TSP(Graph, Source, destination)
		MDTC.IntoTree(MDPC, cost)
	}

	return MDTC
}

// Traveling Salesperson Problem - Dynamic Programming
// https://www.youtube.com/watch?v=XaXsJJh-Q5Y&ab_channel=AbdulBari
func DP_TSP(Graph *topology.Topology, source int, destination int) []int {
	var (
		path    []int // result
		visited []int // The feedback path of subproblems
	)
	node := Graph.GetNodeByID(source)
	path = get_tsp_shortestpath(Graph, node, visited, destination)

	return path
}

// 1. List all possible paths from source to destination, without repeating any nodes
// 2. Divide all combinations into subproblems
//    ex: G(i, set{}) = min(Cij + G(j, set{}), Cik + G(k, set{}), ...)
// 3. The shortest path is completed using the feedback from subproblems
//    ex: G(S, {1}) = min(CS1 + G(1, {2,3,4})) <= G(1, {2,3,4}) = min(C12 + G(2, {4}), C13 + G(3, {4}), C14 + G(4, {D})) <= G(2, {4}) , G(3, {4})
//
//		  		    		- 2 - 4 - D  C12 + G(2, {4}) cost=3
//    G(1, {2,3,4}) =  min{	- 3 - 4 - D  C13 + G(3, {4}) cost=3  }
//				   			- 4 - D      C14 + G(4, {D}) cost=2
//
//	  path = S - 1 - 4 - D
func get_tsp_shortestpath(Graph *topology.Topology, node *topology.Node, visited []int, end int) []int {
	var (
		paths [][]int
	)

	visited = append(visited, node.ID)
	if node.ID == end {
		return visited
	}

	for _, conn := range node.Connections {
		if loopcompare_simplex(conn.ToNodeID, visited) {
			continue
		}

		nextnode := Graph.GetNodeByID(conn.ToNodeID)
		path := get_tsp_shortestpath(Graph, nextnode, visited, end)
		if len(path) != 0 {
			paths = append(paths, path)
		}
	}

	return select_min_path(paths)
}

func select_min_path(paths [][]int) []int {
	var (
		min_path   []int
		min_length = 100
	)

	if len(paths) == 0 {
		return min_path
	}

	for _, path := range paths {
		if len(path) < min_length {
			min_length = len(path)
			min_path = make([]int, len(path))
			copy(min_path, path)
		}
	}
	return min_path
}
