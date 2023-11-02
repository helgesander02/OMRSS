package routes

import "encoding/json"

// Add to the tree based on the given path
func (tree *Tree) IntoTree(P []int, cost float64) {
	for l := len(P) - 1; l > 0; l-- {
		node1, b1 := tree.CheckNodeByID(P[l])
		node2, b2 := tree.CheckNodeByID(P[l-1])

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

// Remove edges
func (MST_prime *Tree) RemoveEdge(e_prime [2]int) {
	node1 := MST_prime.GetNodeByID(e_prime[0])
	node2 := MST_prime.GetNodeByID(e_prime[1])

	for index, conn1 := range node1.Connections {
		if conn1.ToNodeID == node2.ID {
			node1.Connections = append(node1.Connections[:index], node1.Connections[index+1:]...)
		}
	}
	if len(node1.Connections) == 0 {
		for index, node := range MST_prime.Nodes {
			if node.ID == e_prime[0] {
				MST_prime.Nodes = append(MST_prime.Nodes[:index], MST_prime.Nodes[index+1:]...)
			}
		}
	}

	for index, conn2 := range node2.Connections {
		if conn2.ToNodeID == node1.ID {
			node2.Connections = append(node2.Connections[:index], node2.Connections[index+1:]...)
		}
	}
	if len(node2.Connections) == 0 {
		for index, node := range MST_prime.Nodes {
			if node.ID == e_prime[1] {
				MST_prime.Nodes = append(MST_prime.Nodes[:index], MST_prime.Nodes[index+1:]...)
			}
		}
	}
}

// Determine if it is a tree
func (MST_prime *Tree) CheckIsTree(Terminal []int) bool {
	root := MST_prime.Nodes[0]
	visited := make(map[*Node]bool)

	return DFSTree(MST_prime, root, nil, visited, Terminal) && len(visited) == len(MST_prime.Nodes)
}

func DFSTree(MST_prime *Tree, node *Node, parent *Node, visited map[*Node]bool, Terminal []int) bool {
	if visited[node] {
		return false
	}

	visited[node] = true

	for _, conn := range node.Connections {
		if len(node.Connections) == 1 {
			b := false
			for _, tmal := range Terminal {
				if tmal == node.ID {
					b = true
				}
			}
			if !b {
				return false
			}

		}
		if MST_prime.GetNodeByID(conn.ToNodeID) != parent && !(DFSTree(MST_prime, MST_prime.GetNodeByID(conn.ToNodeID), node, visited, Terminal)) {
			return false
		}
	}

	return true
}

// Verify the existence of a Node in a Tree using its ID
func (tree *Tree) CheckNodeByID(id int) (*Node, bool) {
	for _, node := range tree.Nodes {
		if node.ID == id {
			return node, true
		}
	}
	return &Node{ID: id}, false
}

// Find the nodes in the tree by id
func (tree *Tree) GetNodeByID(id int) *Node {
	for _, node := range tree.Nodes {
		if node.ID == id {
			return node
		}
	}
	return nil
}

// Determine if the tree is the same
func (tree1 *Tree) Compare_Trees(tree2 *Tree) bool {
	if tree1.Weight != tree2.Weight {
		return false
	}

	if len(tree1.Nodes) != len(tree2.Nodes) {
		return false
	}

	for i := 0; i < len(tree1.Nodes); i++ {
		node1 := tree1.GetNodeByID(tree1.Nodes[i].ID)
		node2 := tree2.GetNodeByID(tree1.Nodes[i].ID)
		if node1 == nil || node2 == nil {
			return false
		}

		if !node1.Compare_Nodes(node2) {
			return false
		}
	}

	return true
}

// Determine if the nodes are the same
func (node1 *Node) Compare_Nodes(node2 *Node) bool {
	if node1.ID != node2.ID {
		return false
	}

	if len(node1.Connections) != len(node2.Connections) {
		return false
	}

	for i := 0; i < len(node1.Connections); i++ {
		if !Compare_Connections(node1.Connections, node2.Connections) {
			return false
		}
	}

	return true
}

func Compare_Connections(conn1, conn2 []*Connection) bool {
	i := 0
	for _, c1 := range conn1 {
		for _, c2 := range conn2 {
			if c2.ToNodeID == c1.ToNodeID {
				i += 1
			}
		}
	}

	if i == len(conn1) {
		return true
	} else {
		return false
	}
}

// DeepCopy Tree
func (tree1 *Tree) TreeDeepCopy() *Tree {
	if buf, err := json.Marshal(tree1); err != nil {
		return nil
	} else {
		tree2 := &Tree{}
		if err = json.Unmarshal(buf, tree2); err != nil {
			return nil
		}
		return tree2
	}
}
