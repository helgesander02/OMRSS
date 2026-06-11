package routes

import (
	"sort"
	"src/pkg/logger"
	"src/pkg/random"
)

// Global RNG instance (will be set by plan package)
var rng *random.Generator

// SetRNG sets the random number generator for this package
func SetRNG(r *random.Generator) {
	rng = r
}

// Amal P M, Ajish Kumar K S, "An Algorithm for kth Minimum Spanning Tree"
func KSpanningTree(v2v *V2V, steninertree *Route, K int, Source int, Destinations []int, cost float64, Method_Number int) *KRoute {
	K_MSTS := newKRoute(0, 0, 0)        // An array(K MSTS) of length k, which contain k minimum spanning trees
	list_of_trees := newKRoute(0, 0, 0) // stores a list of trees

	// first step of the algorithm finds the minimum spanning tree using Prims algorithm,
	// But this step using MST-Steiner algorithm finds the steiner tree
	// step1 (K=1)
	MST := steninertree // Prims ==> MST-Steiner Select a tree with minimum weight
	MST.Weight = len(steninertree.Nodes) - 1
	K_MSTS.Routes = append(K_MSTS.Routes, MST)

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
					AddE2MST := MST.RouteDeepCopy()
					AddE2MST.IntoTree(E, cost)
					// Determine whether a cycle exists within the tree (MSTHasCycle bool, cyclelist int[])
					if MSTHasCycle, cyclelist := AddE2MST.FindCycle(); MSTHasCycle {
						// Select edges E’ from the cycle
						MST_prime := AddE2MST.RouteDeepCopy()
						E_prime := MST_prime.GetFeedbackEdgeSet(cyclelist, E)
						// After removing E' from the AddE2MST, add it to the list_of_trees
						Traverse_MST(MST_prime, list_of_trees, E_prime, E, Terminal, cost, K)
					}
				}
			}
		}
	}

	if Method_Number == 0 {
		K_MSTS.SelectMinWeight(list_of_trees.Routes, K)

	} else if Method_Number == 1 {
		K_MSTS.SelectIncreasingArithmeticSequenceWeight(list_of_trees.Routes, K)

	} else if Method_Number == 2 {
		K_MSTS.SelectAverageArithmeticSequenceWeight(list_of_trees.Routes, K)

	} else if Method_Number == 3 {
		//K_MSTS.SelectTreeEditDistance(list_of_trees.Routes, K)
		K_MSTS.SelectMinWeightAndTreeEditDistance(list_of_trees.Routes, K)

	} else {
		K_MSTS.SelectTreeEditDistance(list_of_trees.Routes, K)
		//K_MSTS.SelectMinWeightAndTreeEditDistance(list_of_trees.Routes, K)
	}

	logger.Printf("list_of_trees: %d\n", len(list_of_trees.Routes))
	//K_MSTS.Show_KRoute()

	return K_MSTS
}

// 1. Remove edges E’
// 2. Traverse the tree to determine if there are any other cycles
// 3. Repeat steps 1 and 2 in DFS
// 4. Determine whether the quantity of trees in list_of_trees has increased
// 4-1. then, Restore all of them
// 4-2. else, Restore the removed edge E’
// 5. Check if it's a tree and then add it to the list of trees
func Traverse_MST(MST_prime *Route, list_of_trees *KRoute, E_prime [][2]int, E []int, Terminal []int, cost float64, K int) {
	MST_prime_copy := MST_prime.RouteDeepCopy()
	for _, e_prime := range E_prime {
		MST_prime.RemoveEdge(e_prime)
		if MSTHasCycle, cyclelist := MST_prime.FindCycle(); MSTHasCycle {
			notree := len(list_of_trees.Routes)
			new_E_prime := MST_prime.GetFeedbackEdgeSet(cyclelist, E)
			Traverse_MST(MST_prime, list_of_trees, new_E_prime, E, Terminal, cost, K)

			// Determine whether the quantity of trees in list_of_trees has increased
			if notree < len(list_of_trees.Routes) {
				// Restore all of them
				*MST_prime = *MST_prime_copy.RouteDeepCopy()

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
				*MST_prime = *MST_prime_copy.RouteDeepCopy()
			}
		}
	}
}

// Adds tree t into the list, only if tree t is not in the list_of_trees
// If the list contain more than k trees then it removes the tree which have largest weight among them
func Add_ListOfTrees(list_of_trees *KRoute, MST *Route, K int) {
	MST_copy := MST.RouteDeepCopy()
	if !(In_ListOfTrees(list_of_trees, MST_copy)) {
		list_of_trees.Routes = append(list_of_trees.Routes, MST_copy)
	}
	//if len(list_of_trees.Routes) < K {
	//	if !(In_ListOfTrees(list_of_trees, MST_copy)) {
	//		list_of_trees.Routes = append(list_of_trees.Routes, MST_copy)
	//	}
	//} else {
	///	if list_of_trees.Routes[K-1].Weight > MST_copy.Weight {
	//		if !(In_ListOfTrees(list_of_trees, MST_copy)) {
	//			list_of_trees.Routes = append(list_of_trees.Routes, MST_copy)
	//		}
	//	}
	//}

	sort.Slice(list_of_trees.Routes, func(p, q int) bool {
		return list_of_trees.Routes[p].Weight < list_of_trees.Routes[q].Weight
	})
}

func In_ListOfTrees(list_of_trees *KRoute, MST *Route) bool {
	for _, tree := range list_of_trees.Routes {
		if tree.CompareTrees(MST) {
			return true
		}
	}
	return false
}

// 1. Select the minimum weight KTrees [0, 1, 2, 3]
func (K_MSTS *KRoute) SelectMinWeight(list_of_trees []*Route, K int) {
	//logger.Printf("list_of_trees: %d  K: %d\n", len(list_of_trees.Routes), K)
	if len(list_of_trees) >= K {
		treesmap := make(map[int][]*Route)
		w := list_of_trees[0].Weight
		for _, tree := range list_of_trees {
			treesmap[tree.Weight] = append(treesmap[tree.Weight], tree)
			//logger.Printf("Tree Weight %d \n", tree.Weight)
		}

		for len(K_MSTS.Routes) != K {
			selectq := K - len(K_MSTS.Routes)
			if selectq == 0 {
				break
			}

			if value, exists := treesmap[w]; exists {
				if len(value) <= selectq {
					K_MSTS.Routes = append(K_MSTS.Routes, value...)

				} else {
					// Randomly select from trees with same weight
					for q := 0; q < selectq; q++ {
						index := rng.IntN(len(value))
						K_MSTS.Routes = append(K_MSTS.Routes, value[index])
						// Remove selected tree to avoid duplicates
						value[index] = value[len(value)-1]
						value = value[:len(value)-1]
					}
				}

			}
			w += 1
		}
	} else {
		K_MSTS.Routes = append(K_MSTS.Routes, list_of_trees...)
	}
}

// 2. Select Increasing Arithmetic Sequence Weight KTrees [1, 3 ,5 ,7] or [0, 2, 4, 6], up=2
func (K_MSTS *KRoute) SelectIncreasingArithmeticSequenceWeight(list_of_trees []*Route, K int) {
	//logger.Printf("list_of_trees: %d  K: %d\n", len(list_of_trees.Routes), K)
	var ArithmeticSequence int = 2
	if len(list_of_trees) >= K {
		if len(list_of_trees) >= ArithmeticSequence*(K-1) {
			for idx := 0; idx < ArithmeticSequence*(K-1); idx += ArithmeticSequence {
				K_MSTS.Routes = append(K_MSTS.Routes, list_of_trees[idx])
			}
		} else {
			K_MSTS.Routes = append(K_MSTS.Routes, list_of_trees[:K-1]...)
		}
	} else {
		K_MSTS.Routes = append(K_MSTS.Routes, list_of_trees...)
	}

}

// 3. Select Average Arithmetic Sequence Weight KTrees  [0, up, 2up, 3up], up=len(list_of_trees)/K
func (K_MSTS *KRoute) SelectAverageArithmeticSequenceWeight(list_of_trees []*Route, K int) {
	var ArithmeticSequence int = int(float64(len(list_of_trees)) / float64(K-1))
	if len(list_of_trees) >= K {
		for idx := 0; idx < len(list_of_trees); idx += ArithmeticSequence {
			K_MSTS.Routes = append(K_MSTS.Routes, list_of_trees[idx])
		}
	} else {
		K_MSTS.Routes = append(K_MSTS.Routes, list_of_trees...)
	}

}

// 4. Select Edit_Distance
// Mateusz Pawlik, Nikolaus Augsten, "APTED: Tree edit distance: Robust and memory-efficient"
func (K_MSTS *KRoute) SelectTreeEditDistance(list_of_trees []*Route, K int) {
	// matrix = len(list_of_trees.Routes) * K
	numCandidates := len(list_of_trees)
	distances := make([][]float64, numCandidates)
	for i := 0; i < numCandidates; i++ {
		distances[i] = make([]float64, K)
	}

	cumulative := make([]float64, numCandidates)
	selectedIndices := make(map[int]bool)

	// Start from the 1st round since seed tree is already selected
	// Iterate K-1 times to select remaining trees
	for round := 1; round < K; round++ {
		newTree := K_MSTS.Routes[round-1]

		// For each candidate tree (not yet selected), calculate new distance and accumulate
		for i, candidate := range list_of_trees {
			if selectedIndices[i] {
				continue
			}
			d := APTED(newTree, candidate)
			distances[i][round] = d
			cumulative[i] += d
		}

		// Find the candidate with the highest cumulative distance among unselected trees
		maxScore := -1.0
		maxIdx := -1
		for i, score := range cumulative {
			if !selectedIndices[i] && score > maxScore {
				maxScore = score
				maxIdx = i
			}
		}

		// If a candidate is found, add to K_MSTS
		if maxIdx != -1 {
			K_MSTS.Routes = append(K_MSTS.Routes, list_of_trees[maxIdx])
			selectedIndices[maxIdx] = true
			//logger.Println(cumulative)
			//logger.Print("Selected: ", maxIdx, " ", maxScore, "\n")
		} else {
			break
		}
	}
}

func (K_MSTS *KRoute) SelectMinWeightAndTreeEditDistance(list_of_trees []*Route, K int) {
	limit := (K - 1) * 1
	if len(list_of_trees) > limit {
		treesmap := make(map[int][]*Route)
		w := list_of_trees[0].Weight
		for _, tree := range list_of_trees {
			treesmap[tree.Weight] = append(treesmap[tree.Weight], tree)
		}

		minweight_trees := []*Route{}
		for len(minweight_trees) < limit {
			if value, exists := treesmap[w]; exists {
				minweight_trees = append(minweight_trees, value...)
			}
			w += 1
		}
		K_MSTS.SelectTreeEditDistance(minweight_trees, K)

	} else {
		if len(list_of_trees) >= K {
			K_MSTS.SelectTreeEditDistance(list_of_trees, K)

		} else {
			K_MSTS.Routes = append(K_MSTS.Routes, list_of_trees...)
		}
	}
}
