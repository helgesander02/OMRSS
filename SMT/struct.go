package SMT

type Trees struct {
    TSNTrees []*Tree
    AVBTrees []*Tree
}

type Tree struct {
    Nodes []*Node
}

type Node struct {
    ID int
    Connections []*Connection
}

type Connection struct {
    FromNodeID    int         // strat
    ToNodeID      int         // next
    Cost          float64     // (125,000,000 bytes/s) 1Gbps => (750,000 bytes/6ms) 750,000 bytes under 6ms for each link
}

type Graph struct {
    Vertexs []*Vertex
    Count int
    Path [][]int 
}

type Vertex struct {
    ID int
    Visited bool
    Cost int
    Path int
    Edges []*Edge
}

type Edge struct {
    Strat    int         
    End      int         
    Cost     int     
}