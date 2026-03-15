package routes

type KTreesSet struct {
	TSNTrees []*KTrees
	AVBTrees []*KTrees
}

func newKTreesSet() *KTreesSet {
	return &KTreesSet{}
}

type KTrees struct {
	Trees []*Tree
}

func newKTrees() *KTrees {
	return &KTrees{}
}

type TreesSet struct {
	TSNTrees []*Tree
	AVBTrees []*Tree
}

func newTreesSet() *TreesSet {
	return &TreesSet{}
}

type Tree struct {
	Nodes  []*Node
	Weight int
}

func newTree() *Tree {
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

func newConnection(fromNodeID int, toNodeID int, cost float64) *Connection {
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
type KPathsSet struct {
	TSNPaths     []*KPath
	AVBPaths     []*KPath
	CAN2TSNPaths []*KPath
}

func newKPathsSet() *KPathsSet {
	return &KPathsSet{}
}

type KPath struct {
	K      int
	Source int
	Target int
	Paths  []*Path
	Method string
}

func newKPath(k int, source, target int) *KPath {
	return &KPath{
		K:      k,
		Source: source,
		Target: target,
		Paths:  []*Path{},
	}
}

type PathsSet struct {
	TSNPaths     []*Path
	AVBPaths     []*Path
	CAN2TSNPaths []*Path
}

func newPathsSet() *PathsSet {
	return &PathsSet{}
}

type Path struct {
	Method string
	IDs    []int
	Nodes  []*PathNode
	Weight float64
}

func newPath() *Path {
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
