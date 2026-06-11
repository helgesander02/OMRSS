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

// Show_Route prints a Route's hop sequence (IDs), weight, and node/edge
// graph. It is the canonical display for both tree-shaped (multicast)
// and path-shaped (unicast) routings.
func (route *Route) Show_Route() {
	logger.Printf("Route IDs: %v\n", route.IDs)
	logger.Printf("Weight: %d\n", route.Weight)
	for _, node := range route.Nodes {
		logger.Printf("Node %d:\n", node.ID)
		for _, c := range node.Connections {
			logger.Printf("  %d --> %d\n", c.FromNodeID, c.ToNodeID)
		}
	}
}

// Show_Cycle reports whether the route contains a cycle.
func (route *Route) Show_Cycle() {
	hasCycle, cycle := route.FindCycle()
	if hasCycle {
		logger.Println("The MST has cycle")
		logger.Println(cycle)
	} else {
		logger.Println("The MST has no cycle")
	}
}

// Show_KRoute prints every alternative in a K-route bundle.
func (kr *KRoute) Show_KRoute() {
	for index, route := range kr.Routes {
		logger.Printf("Alternative %d (Weight: %d)\n", index, route.Weight)
		route.Show_Route()
	}
}

// Show_RouteSet samples one route per category (TSN / AVB / CAN2TT) so
// the user can eyeball the routing shape without dumping everything.
func (rs *RouteSet) Show_RouteSet() {
	if len(rs.TSNRoutes) > 0 {
		logger.Println("\nTSN Route 1")
		rs.TSNRoutes[0].Show_Route()
	}
	if len(rs.AVBRoutes) > 0 {
		logger.Println("\nAVB Route 1")
		rs.AVBRoutes[0].Show_Route()
	}
	if len(rs.CAN2TTRoutes) > 0 {
		logger.Printf("\nCAN2TT Route 1 (Method: %s)\n", rs.CAN2TTRoutes[0].Method)
		rs.CAN2TTRoutes[0].Show_Route()
	}
}

// Show_KRouteSet samples one K-bundle per category.
func (krs *KRouteSet) Show_KRouteSet() {
	if len(krs.TSNRoutes) > 0 {
		logger.Println("\nTSN K-Route 1")
		krs.TSNRoutes[0].Show_KRoute()
	}
	if len(krs.AVBRoutes) > 0 {
		logger.Println("\nAVB K-Route 1")
		krs.AVBRoutes[0].Show_KRoute()
	}
	if len(krs.CAN2TTRoutes) > 0 {
		logger.Printf("\nCAN2TT K-Route 1 (Method: %s)\n", krs.CAN2TTRoutes[0].Method)
		krs.CAN2TTRoutes[0].Show_KRoute()
	}
}
