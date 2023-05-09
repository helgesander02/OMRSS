package SMT

type Node struct {
    ID string
    Visited bool
    Cost float64
    Path string
}

type NodeTable struct {
    Nodes []*Node
}
