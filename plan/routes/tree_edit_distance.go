package routes

import (
	"fmt"
	"math"
	"strconv"
)

// Cost of deleting a node
func costDelete() float64 {
	return 3.0
}

// Cost of inserting a node
func costInsert() float64 {
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

// Recursively Calculate Deletion Cost of Entire Subtree
func treeCostDelete(n *Node, t *Tree, parentID int) float64 {
	c := costDelete()
	children := getChildren(n, t, parentID)
	for _, child := range children {
		c += treeCostDelete(child, t, n.ID)
	}
	return c
}

// Recursively Calculate Insertion Cost of Entire Subtree
func treeCostInsert(n *Node, t *Tree, parentID int) float64 {
	c := costInsert()
	children := getChildren(n, t, parentID)
	for _, child := range children {
		c += treeCostInsert(child, t, n.ID)
	}
	return c
}

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

// Assumes root node ID string has prefix "100".
// This function searches for root nodes in t.Nodes that match this condition,
// then calculates the edit distance for entire trees.
func APTED(t1, t2 *Tree) float64 {
	var root1, root2 *Node
	for _, n := range t1.Nodes {
		if n.ID >= 1000 && n.ID <= 1999 {
			root1 = n
			break
		}
	}
	for _, n := range t2.Nodes {
		if n.ID >= 1000 && n.ID <= 1999 {
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
