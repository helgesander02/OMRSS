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
		ktrees.ShowKTrees()
		tsn++

		break
	}
	avb := 1
	for _, ktrees := range trees.AVBTrees {
		fmt.Printf("\nAVB Tree %d \n", avb)
		ktrees.ShowKTrees()
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

func (Ktrees *KTrees) ShowKTrees() {
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
	b, cyclelist := tree.FindCycle()
	if b {
		fmt.Println("The MST has cycle")
		fmt.Println(cyclelist)

	} else {
		fmt.Println("The MST has no cycle")
	}
}

// Show functions for Path-based structures (OSRO)
func (kpaths *KPaths_set) Show_KPaths_Set() {
	tsn := 1
	for _, kpath := range kpaths.TSNPaths {
		fmt.Printf("\nTSN K-Path %d \n", tsn)
		kpath.Show_KPath()
		tsn++
		break
	}
	avb := 1
	for _, kpath := range kpaths.AVBPaths {
		fmt.Printf("\nAVB K-Path %d \n", avb)
		kpath.Show_KPath()
		avb++
		break
	}
	can := 1
	for _, kpath := range kpaths.CAN2TSNPaths {
		fmt.Printf("\nCAN2TSN K-Path %d (Method: %s)\n", can, kpath.Method)
		kpath.Show_KPath()
		can++
		break
	}
}

func (paths *Paths_set) Show_Paths_Set() {
	tsn := 1
	for _, path := range paths.TSNPaths {
		fmt.Printf("\nTSN Path %d \n", tsn)
		path.Show_PathStruct()
		tsn++
		break
	}
	avb := 1
	for _, path := range paths.AVBPaths {
		fmt.Printf("\nAVB Path %d \n", avb)
		path.Show_PathStruct()
		avb++
		break
	}
	can := 1
	for _, path := range paths.CAN2TSNPaths {
		fmt.Printf("\nCAN2TSN Path %d (Method: %s)\n", can, path.Method)
		path.Show_PathStruct()
		can++
		break
	}
}

func (kpath *KPath) Show_KPath() {
	for index, path := range kpath.Paths {
		fmt.Printf("Path %d (Weight: %.2f)\n", index, path.Weight)
		path.Show_PathStruct()
	}
}

func (path *Path) Show_PathStruct() {
	fmt.Printf("Path IDs: %v\n", path.IDs)
	fmt.Printf("Weight: %.2f\n", path.Weight)
	for _, node := range path.Nodes {
		fmt.Printf("Node %d:\n", node.ID)
		for _, c := range node.Connections {
			fmt.Printf("  %d --> %d\n", c.FromNodeID, c.ToNodeID)
		}
	}
}
