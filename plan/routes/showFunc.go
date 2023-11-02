package routes

import (
	"fmt"
)

func (v2v *V2V) Show_V2Vs() {
	for _, v2vedge := range v2v.V2VEdges {
		v2vedge.Show_VertexToVertex()
	}
}

func (v2vedge *V2VEdge) Show_VertexToVertex() {
	for _, graph := range v2vedge.Graphs {
		fmt.Printf("From Vertex: %d\n", v2vedge.FromVertex)
		graph.Show_Path()
	}
}

func (graph *Graph) Show_Path() {
	fmt.Printf("To Vertex: %d\n", graph.ToVertex)
	for _, path := range graph.Path {
		fmt.Printf("%v\n", path)
	}
}

func (trees *KTrees_set) Show_kTrees_Set() {
	tsn := 1
	for _, ktrees := range trees.TSNTrees {
		fmt.Printf("\nTSN Tree %d \n", tsn)
		ktrees.Show_KTrees()
		tsn++

		break
	}
	avb := 1
	for _, ktrees := range trees.AVBTrees {
		fmt.Printf("\nAVB Tree %d \n", avb)
		ktrees.Show_KTrees()
		avb++

		break
	}
}

func (trees *Trees_set) Show_Trees_Set() {
	tsn := 1
	for _, tree := range trees.TSNTrees {
		fmt.Printf("\nTSN Tree %d \n", tsn)
		tree.Show_Tree()
		tsn++

		break
	}
	avb := 1
	for _, tree := range trees.AVBTrees {
		fmt.Printf("\nAVB Tree %d \n", avb)
		tree.Show_Tree()
		avb++

		break
	}
}

func (Ktrees *KTrees) Show_KTrees() {
	for index, tree := range Ktrees.Trees {
		fmt.Printf("tree%d \n", index)
		fmt.Printf("tree weight: %d \n", tree.Weight)
		tree.Show_Tree()
	}
}

func (tree *Tree) Show_Tree() {
	for _, node := range tree.Nodes {
		fmt.Println(node.ID)
		for _, c := range node.Connections {
			fmt.Printf("%d --> %d \n", c.FromNodeID, c.ToNodeID)
		}
	}
}

func (tree *Tree) Show_Cycle() {
	b, cyclelist := tree.FindCyCle()
	if b {
		fmt.Println("The MST has cycle")
		fmt.Println(cyclelist)

	} else {
		fmt.Println("The MST has no cycle")
	}
}
