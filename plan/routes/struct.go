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

// OSRO Path-based structures
type KPaths_set struct {
	TSNPaths     []*KPath
	AVBPaths     []*KPath
	CAN2TSNPaths []*KPath
}

func new_KPaths_Set() *KPaths_set {
	return &KPaths_set{}
}

type KPath struct {
	K      int
	Source int
	Target int
	Paths  []*Path
	Method string
}

func new_KPath(k int, source, target int) *KPath {
	return &KPath{
		K:      k,
		Source: source,
		Target: target,
		Paths:  []*Path{},
	}
}

type Paths_set struct {
	TSNPaths     []*Path
	AVBPaths     []*Path
	CAN2TSNPaths []*Path
}

func new_Paths_Set() *Paths_set {
	return &Paths_set{}
}

type Path struct {
	Method string
	IDs    []int
	Nodes  []*PathNode
	Weight float64
}

func new_Path() *Path {
	return &Path{}
}

type PathNode struct {
	ID          int
	Shape       string
	Connections []*PathConnection
}

type PathConnection struct {
	FromNodeID int
	ToNodeID   int
	Cost       float64
}
