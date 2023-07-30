package  trees

import (   
    //"fmt"
)

// Amal P M, Ajish Kumar K S, "An Algorithm for kth Minimum Spanning Tree"
func KSpanningTree(steninertree *Tree, K int) *KTrees {
	K_MSTS := &KTrees{}
	MST := &Tree{}
	list_of_trees := &KTrees{}

	// first step of the algorithm finds the minimum spanning tree using Prims algorithm, 
	// But this step using MST-Steiner algorithm finds the steiner tree
	MST = steninertree // Prims ==> MST-Steiner
	MST.Weight = len(steninertree.Nodes)-1
	list_of_trees.Add(MST)
	MST = list_of_trees.Select_min()
	K_MSTS.Trees = append(K_MSTS.Trees, MST)
	for i:=2; i<=K; i++ {
		list_of_trees.Remove(MST)
		

	}

	return K_MSTS
}

func (list_of_trees *KTrees) Add(MST *Tree) {
	list_of_trees.Trees = append(list_of_trees.Trees, MST)

}

func (list_of_trees *KTrees) Remove(MST *Tree) {
	for index, tree := range list_of_trees.Trees {
		var b bool = true
		if tree.Weight == MST.Weight {
			for i:=0; i<len(tree.Nodes); i++ {
				if tree.Nodes[i] != MST.Nodes[i] {
					b = false
				}
			}
			if b {
				list_of_trees.Trees = append(list_of_trees.Trees[:index], list_of_trees.Trees[index+1:]...)
				break
			}
		}
	}
}

// Select the tree with the smallest weight from the list_of_trees. 
// If there are multiple trees with the same minimum weight, choose based on K.
func (list_of_trees *KTrees) Select_min() *Tree {
	t := list_of_trees.Trees[0]
	w := t.Weight
	for _, tree := range list_of_trees.Trees {
		if tree.Weight < w {
			t = tree
		}
	}
	//t.Show_Tree()
	return t
}