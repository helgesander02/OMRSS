package trees

import (
	"fmt"
)

func (graph *Graph) Show_Graph() {
    for _, g := range graph.Vertexs {
        fmt.Printf("%d\t", g.ID)
        fmt.Printf("%v\t", g.Visited)
        fmt.Printf("%d\t", g.Cost)
        fmt.Printf("%d\n", g.Path)
		for _, e := range g.Edges {
			fmt.Printf("%d-->%d cost:%d\n", e.Strat, e.End, e.Cost)
		}
    } 
}

func (trees *Trees) Show_Trees() {
    tsn := 1
    for _,ktrees := range trees.TSNTrees{  
        fmt.Printf("TSN Tree %d \n", tsn)
        ktrees.Show_KTrees()
        tsn++

        break  
    }
    avb := 1
    for _,ktrees := range trees.AVBTrees{
        fmt.Printf("AVB Tree %d \n", avb)
        ktrees.Show_KTrees()
        avb++

        break
    }
}

func (Ktrees *KTrees) Show_KTrees() {
    for _, tree := range Ktrees.Trees{
        fmt.Printf("tree weight: %d \n", tree.Weight)
        tree.Show_Tree()
    }
}

func (tree *Tree) Show_Tree() {
    for _, node := range tree.Nodes{
        fmt.Println(node.ID)
        for _, c := range node.Connections {
            fmt.Printf("%d --> %d \n", c.FromNodeID, c.ToNodeID)
        }      
    }
}