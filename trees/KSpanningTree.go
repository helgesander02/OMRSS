package trees

import (
	"crypto/rand"
	"math/big"
	"encoding/json"
	"sort"
)

// Amal P M, Ajish Kumar K S, "An Algorithm for kth Minimum Spanning Tree"
func KSpanningTree(v2v *V2V, steninertree *Tree, K int, Source int, Destinations []int, cost float64) *KTrees {
	K_MSTS := &KTrees{}        // An array(K MSTS) of length k, which contain k minimum spanning trees
	MST := &Tree{}             // Select a tree with minimum weight
	list_of_trees := &KTrees{} // stores a list of trees

	// first step of the algorithm finds the minimum spanning tree using Prims algorithm,
	// But this step using MST-Steiner algorithm finds the steiner tree
	// step1 (K=1)
	MST = steninertree // Prims ==> MST-Steiner
	MST.Weight = len(steninertree.Nodes) - 1
	K_MSTS.Trees = append(K_MSTS.Trees, MST)

	// step2~K (K=2~K)
	var Terminal []int
	Terminal = append(Terminal, Source+1000)
	for _, d := range Destinations {
		Terminal = append(Terminal, d+2000)
	}

	// Generate all possible trees and then select K of them based on their weight
	for _, terminal := range Terminal {
		v2vedge, _ := v2v.GetV2VEdge(terminal)
		for _, tmal := range Terminal {
			if terminal == tmal {
				continue

			} else {
				// edge E which is not in the Source to Destinations of MST
				allpath := v2vedge.GetPath(tmal)
				for _, E := range allpath {
					AddE2MST := MSTDeepCopy(MST) // Add E to MST
					AddE2MST.AddTree(E, cost)
					// Determine whether a cycle exists within the tree (MSTHasCycle bool, cyclelist int[])
					MSTHasCycle, cyclelist := AddE2MST.FindCyCle()
					// Select edges E’ from the cycle
					E_prime := AddE2MST.GetFeedbackEdgeSet(cyclelist, E)

					if MSTHasCycle {
						//fmt.Printf("\nFrom %d to %d\n", terminal, tmal)
						//fmt.Printf("E: %v\n", E)
						//fmt.Printf("CycleList: %v\n", cyclelist)
						//fmt.Printf("E': %v\n\n", E_prime)
						//MST.Show_Tree()
						MST_prime := MSTDeepCopy(AddE2MST)

						// After removing E' from the AddE2MST, add it to the list_of_trees
						MST_prime.SearchMST(list_of_trees, AddE2MST, E_prime, E, Terminal, cost, K)
						//list_of_trees.Show_KTrees()
					}
				}
			}
		}
	}
	K_MSTS.Select_Min(list_of_trees, K)

	return K_MSTS
}

// Select K Trees
func (K_MSTS *KTrees) Select_Min(list_of_trees *KTrees, K int) {
	treesmap := make(map[int][]*Tree)	
	for _, tree := range list_of_trees.Trees {	
		treesmap[tree.Weight] = append(treesmap[tree.Weight], tree)
	}

	w := K_MSTS.Trees[0].Weight
	for len(K_MSTS.Trees) != K {
		selectq := K - len(K_MSTS.Trees)
		if value,  exists := treesmap[w]; exists {
			if len(value) <= selectq {
				for _, tree := range value {
					K_MSTS.Trees = append(K_MSTS.Trees, tree)
				}
				delete(treesmap, w)

			} else {
				for q:=0; q<selectq; q++ {
					randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(value))))
					index := int(randomIndex.Int64())
					K_MSTS.Trees = append(K_MSTS.Trees, value[index])
					value = append(value[:index], value[index+1:]...)
				}
			}

		} else {
			w += 1
		}
	}
}

// Search for all trees in MST_prime after removing E'
func (MST_prime *Tree) SearchMST(list_of_trees *KTrees, AddE2MST *Tree, E_prime [][2]int, E []int, Terminal []int, cost float64, K int) {
	for _, e_prime := range E_prime {
		MST_prime.RemoveEdge(e_prime)
		if MSTHasCycle, cyclelist := MST_prime.FindCyCle(); MSTHasCycle {
			notree := len(list_of_trees.Trees) 
			E_prime := MST_prime.GetFeedbackEdgeSet(cyclelist, E)
			MST_prime.SearchMST(list_of_trees, MST_prime, E_prime, E, Terminal, cost, K)
			 // Determine whether the quantity of trees in list_of_trees has increased
			if notree < len(list_of_trees.Trees) {
				// Restore all of them
				MST_prime = MSTDeepCopy(AddE2MST)
			} else { 
				// Restore the removed edge
				for _, P := range E_prime {
					p := make([]int, len(P))
					copy(p, P[:])
					MST_prime.AddTree(p, cost)
				}
			}
		}
		if MST_prime.IsTree(Terminal) {
			MST_prime.Weight = len(MST_prime.Nodes) - 1
			list_of_trees.Add(MST_prime, K)
			MST_prime = MSTDeepCopy(AddE2MST)
		}
	}
}

func (MST_prime *Tree) RemoveEdge(e_prime [2]int) {
	node1 := MST_prime.getNodeByID(e_prime[0])
	node2 := MST_prime.getNodeByID(e_prime[1])

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

// Confirm if it is a tree after removing E'
func (MST_prime *Tree) IsTree(Terminal []int) bool {
	root := MST_prime.Nodes[0]
	visited := make(map[*Node]bool)

	return MST_prime.DFSTree(root, nil, visited, Terminal) && len(visited) == len(MST_prime.Nodes)
}

func (MST_prime *Tree) DFSTree(node *Node, parent *Node, visited map[*Node]bool, Terminal []int) bool {
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
		if MST_prime.getNodeByID(conn.ToNodeID) != parent && !(MST_prime.DFSTree(MST_prime.getNodeByID(conn.ToNodeID), node, visited, Terminal)) {
			return false
		}
	}

	return true
}

// Adds tree t into the list, only if tree t is not in the list_of_trees
// If the list contain more than k trees then it removes the tree which have largest weight among them
func (list_of_trees *KTrees) Add(MST *Tree, K int) {
	if len(list_of_trees.Trees) < K {
		if !(list_of_trees.InListOfTrees(MST)) {
			list_of_trees.Trees = append(list_of_trees.Trees, MST)
		}

	} else {
		if list_of_trees.Trees[K-1].Weight > MST.Weight {
			if !(list_of_trees.InListOfTrees(MST)) {
				list_of_trees.Trees = append(list_of_trees.Trees, MST)
			}
		}
	}

	sort.Slice(list_of_trees.Trees, func(p, q int) bool {
		return list_of_trees.Trees[p].Weight < list_of_trees.Trees[q].Weight
	})
}

func (list_of_trees *KTrees) InListOfTrees(MST *Tree) bool {
	for _, tree := range list_of_trees.Trees {
		if compareTrees(tree, MST){
			return true
		}
	}
	return false
}

func compareTrees(tree1, tree2 *Tree) bool {
    if tree1.Weight != tree2.Weight {
        return false
    }

    if len(tree1.Nodes) != len(tree2.Nodes) {
        return false
    }

    for i := 0; i < len(tree1.Nodes); i++ {
		node1 := tree1.getNodeByID(tree1.Nodes[i].ID)
		node2 := tree2.getNodeByID(tree1.Nodes[i].ID)
		if node1 == nil || node2 == nil {
			return false
		}

        if !compareNodes(node1, node2) {
            return false
        }
    }

    return true
}

func compareNodes(node1, node2 *Node) bool {
    if node1.ID != node2.ID {
        return false
    }

    if len(node1.Connections) != len(node2.Connections) {
        return false
    }

    for i := 0; i < len(node1.Connections); i++ {
        if !compareConnections(node1.Connections, node2.Connections) {
            return false
        }
    }

    return true
}

func compareConnections(conn1, conn2 []*Connection) bool {
	i:=0
    for _, c1 := range conn1 {
		for _, c2 := range conn2 {
			if c2.ToNodeID == c1.ToNodeID {
				i+=1
			}
		}
	}

	if i == len(conn1) {
		return true
	} else {
		return false
	}
    
}

// Copy MST
func MSTDeepCopy(MST1 *Tree) *Tree {
	if buf, err := json.Marshal(MST1); err != nil {
		return nil
	} else {
		MST2 := &Tree{}
		if err = json.Unmarshal(buf, MST2); err != nil {
			return nil
		}
		return MST2
	}
}
