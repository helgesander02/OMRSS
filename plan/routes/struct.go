package routes

type KTrees_set struct {
	TSNTrees []*KTrees
	AVBTrees []*KTrees
}

func new_KTrees_Set() *KTrees_set {
	return &KTrees_set{}
}

type KTrees struct {
	Trees []*Tree
}

func new_KTrees() *KTrees {
	return &KTrees{}
}

type Trees_set struct {
	TSNTrees []*Tree
	AVBTrees []*Tree
}

func new_Trees_Set() *Trees_set {
	return &Trees_set{}
}

type Tree struct {
	Nodes  []*Node
	Weight int
}

func new_Tree() *Tree {
	return &Tree{}
}

type Node struct {
	ID          int
	Connections []*Connection
}

type Connection struct {
	FromNodeID int     // strat
	ToNodeID   int     // next
	Cost       float64 // 1Gbps => (750,000 bytes/6ms) 750,000 bytes under 6ms for each link ==> 125 bytes/us
}

func new_Connection(fromNodeID int, toNodeID int, cost float64) *Connection {
	return &Connection{
		FromNodeID: fromNodeID,
		ToNodeID:   toNodeID,
		Cost:       cost,
	}
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
