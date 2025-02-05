package routes

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// -------------------------
// Operation Cost Functions
// -------------------------
// Cost of deleting a node
func costDelete(n *Node) float64 {
	return 1.0
}

// Cost of inserting a node
func costInsert(n *Node) float64 {
	return 1.0
}

// Cost of renaming (or matching) two nodes:
// Cost is 0 if node IDs are identical, otherwise 1
func costRename(n1, n2 *Node) float64 {
	s1 := strconv.Itoa(n1.ID)
	s2 := strconv.Itoa(n2.ID)
	if s1 == s2 {
		return 0.0
	}
	return 1.0
}

// -------------------------
// Helper Functions
// -------------------------
// Find a Node in the Tree's Nodes by its ID
func getNodeByID(t *Tree, id int) *Node {
	for _, n := range t.Nodes {
		if n.ID == id {
			return n
		}
	}
	return nil
}

// Get all "child nodes" from a node's Connections based on undirected graph characteristics.
// parentID is used to exclude returning to the parent node during DFS,
// initial call can pass -1.
func getChildren(node *Node, t *Tree, parentID int) []*Node {
	var children []*Node
	for _, conn := range node.Connections {
		// Exclude connection pointing to parent node
		if conn.ToNodeID == parentID {
			continue
		}
		child := getNodeByID(t, conn.ToNodeID)
		if child != nil {
			children = append(children, child)
		}
	}
	return children
}

// -------------------------
// Recursively Calculate Deletion Cost of Entire Subtree
// -------------------------
func treeCostDelete(n *Node, t *Tree, parentID int) float64 {
	c := costDelete(n)
	children := getChildren(n, t, parentID)
	for _, child := range children {
		c += treeCostDelete(child, t, n.ID)
	}
	return c
}

// -------------------------
// Recursively Calculate Insertion Cost of Entire Subtree
// -------------------------
func treeCostInsert(n *Node, t *Tree, parentID int) float64 {
	c := costInsert(n)
	children := getChildren(n, t, parentID)
	for _, child := range children {
		c += treeCostInsert(child, t, n.ID)
	}
	return c
}

// -------------------------
// Tree Edit Distance Calculation (Simplified Version Based on APTED Concept)
// -------------------------
// This function calculates the edit distance between subtrees rooted at n1 and n2,
// and uses dynamic programming to calculate matching costs between two "forests".
// parent1 and parent2 are the parent node IDs of n1 and n2 respectively,
// initially passed as -1.
func treeEditDistance(n1, n2 *Node, t1, t2 *Tree, parent1, parent2 int, memo map[string]float64) float64 {
	// Create memo key
	key := strconv.Itoa(n1.ID) + "-" + strconv.Itoa(n2.ID)
	if v, ok := memo[key]; ok {
		return v
	}

	// Get children of n1 and n2
	children1 := getChildren(n1, t1, parent1)
	children2 := getChildren(n2, t2, parent2)
	m := len(children1)
	n := len(children2)

	// Create DP matrix where dp[i][j] represents minimum matching cost
	// for first i children of children1 and first j children of children2 (viewed as forests)
	dp := make([][]float64, m+1)
	for i := 0; i <= m; i++ {
		dp[i] = make([]float64, n+1)
	}

	dp[0][0] = 0
	for i := 1; i <= m; i++ {
		dp[i][0] = dp[i-1][0] + treeCostDelete(children1[i-1], t1, n1.ID)
	}
	for j := 1; j <= n; j++ {
		dp[0][j] = dp[0][j-1] + treeCostInsert(children2[j-1], t2, n2.ID)
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			costDel := dp[i-1][j] + treeCostDelete(children1[i-1], t1, n1.ID)
			costIns := dp[i][j-1] + treeCostInsert(children2[j-1], t2, n2.ID)
			costRen := dp[i-1][j-1] + treeEditDistance(children1[i-1], children2[j-1], t1, t2, n1.ID, n2.ID, memo)
			dp[i][j] = math.Min(costDel, math.Min(costIns, costRen))
		}
	}

	// Add rename cost between root nodes
	dist := costRename(n1, n2) + dp[m][n]
	memo[key] = dist
	return dist
}

// -------------------------
// APTED Main Function
// -------------------------
// Assumes root node ID string has prefix "100".
// This function searches for root nodes in t.Nodes that match this condition,
// then calculates the edit distance for entire trees.
func APTED(t1, t2 *Tree) float64 {
	var root1, root2 *Node
	for _, n := range t1.Nodes {
		if strings.HasPrefix(strconv.Itoa(n.ID), "100") {
			root1 = n
			break
		}
	}
	for _, n := range t2.Nodes {
		if strings.HasPrefix(strconv.Itoa(n.ID), "100") {
			root2 = n
			break
		}
	}
	if root1 == nil || root2 == nil {
		fmt.Println("At least one tree has no root node (ID must start with 100)")
		return -1
	}

	memo := make(map[string]float64)
	return treeEditDistance(root1, root2, t1, t2, -1, -1, memo)
}

// -------------------------
// Example Testing
// -------------------------
func APTED_Testing() {
	// -------------------------
	// Define tree1
	// Root node: ID starts with "100"
	// Intermediate nodes: e.g., 0 and 2 (random non-repeating)
	// Leaf nodes: ID starts with "200"
	// Example structure:
	//        1001
	//       /    \
	//      0      2
	//      |      |
	//    2003   2002
	//
	// Connection settings: conn1001_0 and conn1001_2 stored in 1001,
	//                      conn0_1001 and conn0_2003 stored in 0,
	//                      conn2_1001 and conn2_2002 stored in 2,
	//                      conn2002_2 stored in 2002,
	//                      conn2003_0 stored in 2003.
	n1001 := &Node{ID: 1001}
	n0 := &Node{ID: 0}
	n2 := &Node{ID: 2}
	n2002 := &Node{ID: 2002}
	n2003 := &Node{ID: 2003}

	conn1001_0 := &Connection{FromNodeID: n1001.ID, ToNodeID: n0.ID, Cost: 1.0}
	conn1001_2 := &Connection{FromNodeID: n1001.ID, ToNodeID: n2.ID, Cost: 1.0}
	conn2_1001 := &Connection{FromNodeID: n2.ID, ToNodeID: n1001.ID, Cost: 1.0}
	conn2_2002 := &Connection{FromNodeID: n2.ID, ToNodeID: n2002.ID, Cost: 1.0}
	conn0_1001 := &Connection{FromNodeID: n0.ID, ToNodeID: n1001.ID, Cost: 1.0}
	conn0_2003 := &Connection{FromNodeID: n0.ID, ToNodeID: n2003.ID, Cost: 1.0}
	conn2002_2 := &Connection{FromNodeID: n2002.ID, ToNodeID: n2.ID, Cost: 1.0}
	conn2003_0 := &Connection{FromNodeID: n2003.ID, ToNodeID: n0.ID, Cost: 1.0}

	n1001.Connections = []*Connection{conn1001_0, conn1001_2}
	n2.Connections = []*Connection{conn2_1001, conn2_2002}
	n0.Connections = []*Connection{conn0_1001, conn0_2003}
	n2002.Connections = []*Connection{conn2002_2}
	n2003.Connections = []*Connection{conn2003_0}

	tree1 := &Tree{
		Nodes:  []*Node{n1001, n0, n2, n2002, n2003},
		Weight: 0,
	}

	// -------------------------
	// Define tree2
	// Example structure:
	//         1001
	//         /
	//        0
	//       / \
	//      2   2003
	//       \
	//        3
	//         \
	//         2002
	//
	// Connection settings:
	//   n1001 (1001) stores connection to 0,
	//   n0 (0) stores connections to 1001, 2003, 2,
	//   n2 (2) stores connections to 0, 3,
	//   n3 (3) stores connections to 2, 2002,
	//   n2002 (2002) stores connection to 3,
	//   n2003 (2003) stores connection to 0.
	n10012 := &Node{ID: 1001}
	n02 := &Node{ID: 0}
	n22 := &Node{ID: 2}
	n32 := &Node{ID: 3}
	n20022 := &Node{ID: 2002}
	n20032 := &Node{ID: 2003}

	conn1001_02 := &Connection{FromNodeID: n10012.ID, ToNodeID: n02.ID, Cost: 1.0}
	conn0_10012 := &Connection{FromNodeID: n02.ID, ToNodeID: n10012.ID, Cost: 1.0}
	conn0_20032 := &Connection{FromNodeID: n02.ID, ToNodeID: n20032.ID, Cost: 1.0}
	conn0_22 := &Connection{FromNodeID: n02.ID, ToNodeID: n22.ID, Cost: 1.0}
	conn2_02 := &Connection{FromNodeID: n22.ID, ToNodeID: n02.ID, Cost: 1.0}
	conn2_32 := &Connection{FromNodeID: n22.ID, ToNodeID: n32.ID, Cost: 1.0}
	conn3_22 := &Connection{FromNodeID: n32.ID, ToNodeID: n22.ID, Cost: 1.0}
	conn3_20022 := &Connection{FromNodeID: n32.ID, ToNodeID: n20022.ID, Cost: 1.0}
	conn2002_32 := &Connection{FromNodeID: n20022.ID, ToNodeID: n32.ID, Cost: 1.0}
	conn2003_02 := &Connection{FromNodeID: n20032.ID, ToNodeID: n02.ID, Cost: 1.0}

	n10012.Connections = []*Connection{conn1001_02}
	n22.Connections = []*Connection{conn2_02, conn2_32}
	n02.Connections = []*Connection{conn0_10012, conn0_20032, conn0_22}
	n32.Connections = []*Connection{conn3_22, conn3_20022}
	n20022.Connections = []*Connection{conn2002_32}
	n20032.Connections = []*Connection{conn2003_02}

	tree2 := &Tree{
		Nodes:  []*Node{n10012, n02, n22, n32, n20022, n20032},
		Weight: 0,
	}

	// -------------------------
	// Calculate edit distance between tree1 and tree2
	// -------------------------
	distance := APTED(tree1, tree2)
	fmt.Printf("Edit distance between tree1 and tree2: %v\n", distance)
}
