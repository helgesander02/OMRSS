package routes

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sort"
)

// Amal P M, Ajish Kumar K S, "An Algorithm for kth Minimum Spanning Tree"
func KSpanningTree(v2v *V2V, steninertree *Tree, K int, Source int, Destinations []int, cost float64, Method_Number int) *KTrees {
	K_MSTS := new_KTrees()        // An array(K MSTS) of length k, which contain k minimum spanning trees
	list_of_trees := new_KTrees() // stores a list of trees

	// first step of the algorithm finds the minimum spanning tree using Prims algorithm,
	// But this step using MST-Steiner algorithm finds the steiner tree
	// step1 (K=1)
	MST := steninertree // Prims ==> MST-Steiner Select a tree with minimum weight
	MST.Weight = len(steninertree.Nodes) - 1
	K_MSTS.Trees = append(K_MSTS.Trees, MST)

	// step2~K (K=2~K)
	var Terminal []int
	Terminal = append(Terminal, Source)
	Terminal = append(Terminal, Destinations...)

	// Generate all possible trees and then select K of them based on their weight
	for _, terminal := range Terminal {
		v2vedge, _ := v2v.GetV2VEdge(terminal)
		for _, tmal := range Terminal {
			if terminal == tmal {
				continue

			} else {
				// edge E which is not in the Source to Destinations of MST
				allpath := v2vedge.GetV2VPath(tmal)
				for _, E := range allpath {
					// Add E to MST
					AddE2MST := MST.TreeDeepCopy()
					AddE2MST.IntoTree(E, cost)
					// Determine whether a cycle exists within the tree (MSTHasCycle bool, cyclelist int[])
					if MSTHasCycle, cyclelist := AddE2MST.FindCyCle(); MSTHasCycle {
						// Select edges E’ from the cycle
						MST_prime := AddE2MST.TreeDeepCopy()
						E_prime := MST_prime.GetFeedbackEdgeSet(cyclelist, E)

						// After removing E' from the AddE2MST, add it to the list_of_trees
						Traverse_MST(MST_prime, list_of_trees, E_prime, E, Terminal, cost, K)
					}
				}
			}
		}
	}

	if Method_Number == 0 {
		K_MSTS.Select_Min_Weight(list_of_trees, K)

	} else if Method_Number == 1 {
		K_MSTS.Select_Increasing_Arithmetic_Sequence_Weight(list_of_trees, K)

	} else {
		K_MSTS.Select_Average_Arithmetic_Sequence_Weight(list_of_trees, K)

	}
	//K_MSTS.Select_High_Dissimilarity(list_of_trees, K)
	fmt.Printf("list_of_trees: %d\n", len(list_of_trees.Trees))

	return K_MSTS
}

// 1. Remove edges E’
// 2. Traverse the tree to determine if there are any other cycles
// 3. Repeat steps 1 and 2 in DFS
// 4. Determine whether the quantity of trees in list_of_trees has increased
// 4-1. then, Restore all of them
// 4-2. else, Restore the removed edge E’
// 5. Check if it's a tree and then add it to the list of trees
func Traverse_MST(MST_prime *Tree, list_of_trees *KTrees, E_prime [][2]int, E []int, Terminal []int, cost float64, K int) {
	MST_prime_copy := MST_prime.TreeDeepCopy()
	for _, e_prime := range E_prime {
		MST_prime.RemoveEdge(e_prime)
		if MSTHasCycle, cyclelist := MST_prime.FindCyCle(); MSTHasCycle {
			notree := len(list_of_trees.Trees)
			new_E_prime := MST_prime.GetFeedbackEdgeSet(cyclelist, E)
			Traverse_MST(MST_prime, list_of_trees, new_E_prime, E, Terminal, cost, K)

			// Determine whether the quantity of trees in list_of_trees has increased
			if notree < len(list_of_trees.Trees) {
				// Restore all of them
				*MST_prime = *MST_prime_copy.TreeDeepCopy()

			} else {
				// Restore the removed edge
				for _, P := range new_E_prime {
					p := make([]int, len(P))
					copy(p, P[:])
					MST_prime.IntoTree(p, cost)
				}
			}

		} else {
			// Confirm if it is a tree after removing E'
			if MST_prime.CheckIsTree(Terminal) {
				MST_prime.Weight = len(MST_prime.Nodes) - 1
				Add_ListOfTrees(list_of_trees, MST_prime, K)

				// Restore the removed edge
				*MST_prime = *MST_prime_copy.TreeDeepCopy()
			}
		}
	}
}

// Adds tree t into the list, only if tree t is not in the list_of_trees
// If the list contain more than k trees then it removes the tree which have largest weight among them
func Add_ListOfTrees(list_of_trees *KTrees, MST *Tree, K int) {
	MST_copy := MST.TreeDeepCopy()
	if !(In_ListOfTrees(list_of_trees, MST_copy)) {
		list_of_trees.Trees = append(list_of_trees.Trees, MST_copy)
	}

	sort.Slice(list_of_trees.Trees, func(p, q int) bool {
		return list_of_trees.Trees[p].Weight < list_of_trees.Trees[q].Weight
	})
}

func In_ListOfTrees(list_of_trees *KTrees, MST *Tree) bool {
	for _, tree := range list_of_trees.Trees {
		if tree.Compare_Trees(MST) {
			return true
		}
	}
	return false
}

// 1. Select the minimum weight KTrees [0, 1, 2, 3]
func (K_MSTS *KTrees) Select_Min_Weight(list_of_trees *KTrees, K int) {
	if len(list_of_trees.Trees) >= K {
		treesmap := make(map[int][]*Tree)
		for _, tree := range list_of_trees.Trees {
			treesmap[tree.Weight] = append(treesmap[tree.Weight], tree)
		}

		w := K_MSTS.Trees[0].Weight
		for len(K_MSTS.Trees) != K {
			selectq := K - len(K_MSTS.Trees)
			if value, exists := treesmap[w]; exists {
				if len(value) <= selectq {
					K_MSTS.Trees = append(K_MSTS.Trees, value...)
					delete(treesmap, w)

				} else {
					for q := 0; q < selectq; q++ {
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
	} else {
		K_MSTS.Trees = append(K_MSTS.Trees, list_of_trees.Trees...)
	}
}

// 2. Select Increasing Arithmetic Sequence Weight KTrees [1, 3 ,5 ,7] or [0, 2, 4, 6], up=2
func (K_MSTS *KTrees) Select_Increasing_Arithmetic_Sequence_Weight(list_of_trees *KTrees, K int) {
	var ArithmeticSequence int = 2
	if len(list_of_trees.Trees) >= K {
		if len(list_of_trees.Trees) >= ArithmeticSequence*(K-1) {
			for idx := 0; idx < ArithmeticSequence*(K-1); idx += ArithmeticSequence {
				K_MSTS.Trees = append(K_MSTS.Trees, list_of_trees.Trees[idx])
			}
		} else {
			K_MSTS.Trees = append(K_MSTS.Trees, list_of_trees.Trees[:K-1]...)
		}
	} else {
		K_MSTS.Trees = append(K_MSTS.Trees, list_of_trees.Trees...)
	}

}

// 3. Select Average Arithmetic Sequence Weight KTrees  [0, up, 2up, 3up], up=len(list_of_trees)/K
func (K_MSTS *KTrees) Select_Average_Arithmetic_Sequence_Weight(list_of_trees *KTrees, K int) {
	var ArithmeticSequence int = int(float64(len(list_of_trees.Trees)) / float64(K-1))
	if len(list_of_trees.Trees) >= K {
		for idx := 0; idx < len(list_of_trees.Trees); idx += ArithmeticSequence {
			K_MSTS.Trees = append(K_MSTS.Trees, list_of_trees.Trees[idx])
		}
	} else {
		K_MSTS.Trees = append(K_MSTS.Trees, list_of_trees.Trees...)
	}

}

// 4. Select high dissimilarity KTrees
func (K_MSTS *KTrees) Select_High_Dissimilarity(list_of_trees *KTrees, K int) {

}
