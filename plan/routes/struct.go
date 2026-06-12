package routes

// Route is the single routing structure used everywhere — tree-shaped
// (multicast Steiner) and path-shaped (unicast linear) routings share
// the same node/edge graph. The optional IDs and Method fields are
// populated by path callers (IDs = linear hop sequence; Method = encap
// method tag for CAN2TT). Tree callers leave them zero.
type Route struct {
	Nodes  []*Node
	Weight int    // hop count
	IDs    []int  // path callers only
	Method string // CAN2TT path callers only
}

func newRoute() *Route {
	return &Route{}
}

// KRoute holds K candidate Routes for one (Source, Target) stream so
// OSACO can pick between them via pheromone + visibility.
type KRoute struct {
	K      int
	Source int
	Target int
	Method string
	Routes []*Route
}

func newKRoute(k int, source, target int) *KRoute {
	return &KRoute{
		K:      k,
		Source: source,
		Target: target,
		Routes: []*Route{},
	}
}

// RouteSet groups one Route per flow category. CAN2TT is populated only
// by OSRO; OMACO leaves it empty.
type RouteSet struct {
	TSNRoutes    []*Route
	AVBRoutes    []*Route
	CAN2TTRoutes []*Route
}

func newRouteSet() *RouteSet {
	return &RouteSet{}
}

// KRouteSet is the K-candidate version of RouteSet.
type KRouteSet struct {
	TSNRoutes    []*KRoute
	AVBRoutes    []*KRoute
	CAN2TTRoutes []*KRoute
}

func newKRouteSet() *KRouteSet {
	return &KRouteSet{}
}

type Node struct {
	ID          int
	Connections []*Connection
}

type Connection struct {
	FromNodeID int     // start
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
}
