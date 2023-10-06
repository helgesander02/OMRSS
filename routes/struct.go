package routes

type KTrees_set struct {
	TSNTrees []*KTrees
	AVBTrees []*KTrees
}

type KTrees struct {
	Trees []*Tree
}

type Trees_set struct {
	TSNTrees []*Tree
	AVBTrees []*Tree
}

type Tree struct {
	Nodes  []*Node
	Weight int
}

type Node struct {
	ID          int
	Connections []*Connection
}

type Connection struct {
	FromNodeID int     // strat
	ToNodeID   int     // next
	Cost       float64 // (125,000,000 bytes/s) 1Gbps => (750,000 bytes/6ms) 750,000 bytes under 6ms for each link ==> 125 bytes/ns
}

type V2V struct {
	V2VEdges []*V2VEdge
}

type V2VEdge struct {
	FromVertex int
	Graphs     []*Graph
}

type Graph struct {
	Vertexs  []*Vertex
	ToVertex int
	Path     [][]int
}

type Vertex struct {
	ID      int
	Visited bool
	Cost    int
	Path    int
	Edges   []*Edge
}

type Edge struct {
	Strat int
	End   int
	Cost  int
}
