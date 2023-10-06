package routes

//tutorialspoint, "Detect Cycle in a an Undirected Graph"
func (tree *Tree) FindCyCle() (bool, []int) {
	var cyclelist []int
	for _, node := range tree.Nodes {
		visited := make(map[int]bool)
		if tree.DFSCyCle(node, visited, -1, node.ID) {
			cyclelist = append(cyclelist, node.ID)

		}
	}
	if len(cyclelist) != 0 {
		return true, cyclelist
	} else {
		return false, cyclelist
	}
}

func (tree *Tree) DFSCyCle(node *Node, visited map[int]bool, parentID int, startID int) bool {
	visited[node.ID] = true
	for _, conn := range node.Connections {
		if conn.ToNodeID == parentID {
			continue
		}
		if conn.ToNodeID == startID {
			return true
		}
		if visited[conn.ToNodeID] {
			continue
		}
		toNode := tree.GetNodeByID(conn.ToNodeID)
		if tree.DFSCyCle(toNode, visited, node.ID, startID) {
			return true
		}
	}
	return false
}

func (MST_prime *Tree) GetFeedbackEdgeSet(cyclelist []int, E []int) [][2]int {
	var E_prime [][2]int
	for _, cycle := range cyclelist {
		node := MST_prime.GetNodeByID(cycle)
		for _, conn := range node.Connections {
			if !(InCycleList(cyclelist, conn.ToNodeID)) {
				continue
			} else {
				var nodeconn [2]int
				nodeconn[0] = conn.FromNodeID
				nodeconn[1] = conn.ToNodeID
				if InE(E, nodeconn) {
					continue
				}
				if !(InEPrime(E_prime, nodeconn)) {
					E_prime = append(E_prime, nodeconn)
				}
			}
		}
	}

	return E_prime
}

func InCycleList(cyclelist []int, ToNodeID int) bool {
	for _, id := range cyclelist {
		if id == ToNodeID {
			return true
		}
	}
	return false
}

func InE(E []int, nodeconn [2]int) bool {
	for index, id := range E {
		if id == nodeconn[0] {
			if E[index+1] == nodeconn[1] {
				return true
			}
			if E[index-1] == nodeconn[1] {
				return true
			}
		}
	}
	return false
}

func InEPrime(E_prime [][2]int, nodeconn [2]int) bool {
	for _, e_prime := range E_prime {
		if e_prime[0] == nodeconn[0] && e_prime[1] == nodeconn[1] {
			return true
		}
		if e_prime[0] == nodeconn[1] && e_prime[1] == nodeconn[0] {
			return true
		}
	}
	return false
}

func (tree *Tree) GetNodeByID(id int) *Node {
	for _, node := range tree.Nodes {
		if node.ID == id {
			return node
		}
	}
	return nil
}
