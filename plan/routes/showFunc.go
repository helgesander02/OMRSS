package routes

import "src/pkg/logger"

func (v2v *V2V) Show_V2Vs() {
	for _, v2vedge := range v2v.V2VEdges {
		v2vedge.Show_VertexToVertex()
	}
}

func (v2vedge *V2VEdge) Show_VertexToVertex() {
	for _, graph := range v2vedge.Graphs {
		logger.Printf("From Vertex: %d\n", v2vedge.FromVertex)
		graph.Show_Path()
	}
}

func (graph *Graph) Show_Path() {
	logger.Printf("To Vertex: %d\n", graph.ToVertex)
	for _, path := range graph.Path {
		logger.Printf("%v\n", path)
	}
}

func (trees *KTreesSet) Show_kTrees_Set() {
	tsn := 1
	for _, ktrees := range trees.TSNTrees {
		logger.Printf("\nTSN Tree %d \n", tsn)
		ktrees.ShowKTrees()
		tsn++

		break
	}
	avb := 1
	for _, ktrees := range trees.AVBTrees {
		logger.Printf("\nAVB Tree %d \n", avb)
		ktrees.ShowKTrees()
		avb++

		break
	}
}

func (trees *TreesSet) Show_Trees_Set() {
	tsn := 1
	for _, tree := range trees.TSNTrees {
		logger.Printf("\nTSN Tree %d \n", tsn)
		tree.Show_Tree()
		tsn++

		break
	}
	avb := 1
	for _, tree := range trees.AVBTrees {
		logger.Printf("\nAVB Tree %d \n", avb)
		tree.Show_Tree()
		avb++

		break
	}
}

func (Ktrees *KTrees) ShowKTrees() {
	for index, tree := range Ktrees.Trees {
		logger.Printf("tree%d \n", index)
		logger.Printf("tree weight: %d \n", tree.Weight)
		tree.Show_Tree()
	}
}

func (tree *Tree) Show_Tree() {
	for _, node := range tree.Nodes {
		logger.Println(node.ID)
		for _, c := range node.Connections {
			logger.Printf("%d --> %d \n", c.FromNodeID, c.ToNodeID)
		}
	}
}

func (tree *Tree) Show_Cycle() {
	b, cyclelist := tree.FindCycle()
	if b {
		logger.Println("The MST has cycle")
		logger.Println(cyclelist)

	} else {
		logger.Println("The MST has no cycle")
	}
}

// Show functions for Path-based structures (OSRO)
func (kpaths *KPathsSet) Show_KPaths_Set() {
	tsn := 1
	for _, kpath := range kpaths.TSNPaths {
		logger.Printf("\nTSN K-Path %d \n", tsn)
		kpath.Show_KPath()
		tsn++
		break
	}
	avb := 1
	for _, kpath := range kpaths.AVBPaths {
		logger.Printf("\nAVB K-Path %d \n", avb)
		kpath.Show_KPath()
		avb++
		break
	}
	can := 1
	for _, kpath := range kpaths.CAN2TSNPaths {
		logger.Printf("\nCAN2TSN K-Path %d (Method: %s)\n", can, kpath.Method)
		kpath.Show_KPath()
		can++
		break
	}
}

func (paths *PathsSet) Show_Paths_Set() {
	tsn := 1
	for _, path := range paths.TSNPaths {
		logger.Printf("\nTSN Path %d \n", tsn)
		path.Show_PathStruct()
		tsn++
		break
	}
	avb := 1
	for _, path := range paths.AVBPaths {
		logger.Printf("\nAVB Path %d \n", avb)
		path.Show_PathStruct()
		avb++
		break
	}
	can := 1
	for _, path := range paths.CAN2TSNPaths {
		logger.Printf("\nCAN2TSN Path %d (Method: %s)\n", can, path.Method)
		path.Show_PathStruct()
		can++
		break
	}
}

func (kpath *KPath) Show_KPath() {
	for index, path := range kpath.Paths {
		logger.Printf("Path %d (Weight: %.2f)\n", index, path.Weight)
		path.Show_PathStruct()
	}
}

func (path *Path) Show_PathStruct() {
	logger.Printf("Path IDs: %v\n", path.IDs)
	logger.Printf("Weight: %.2f\n", path.Weight)
	for _, node := range path.Nodes {
		logger.Printf("Node %d:\n", node.ID)
		for _, c := range node.Connections {
			logger.Printf("  %d --> %d\n", c.FromNodeID, c.ToNodeID)
		}
	}
}
