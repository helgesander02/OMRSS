package routes

import (
	"src/network"
	"src/pkg/config"
	"src/pkg/logger"
)

var v2v *V2V = &V2V{} // v2v is all paths connecting multiple terminals to terminals.

func Get_SteninerTree_Routing(network *network.Network, cfg *config.Config) *RouteSet {
	TreesSet := newRouteSet()

	for _, flow := range network.FlowSet.TSNFlows {
		graph := network.GraphSet.Get(flow.Source, flow.Destinations)
		tree := SteninerTree(v2v, graph, flow.Source, flow.Destinations, cfg.Network.ByteRate)
		TreesSet.TSNRoutes = append(TreesSet.TSNRoutes, tree)
	}
	logger.Printf("Finish Steniner Tree %d TSN streams routing\n", len(TreesSet.TSNRoutes))

	for _, flow := range network.FlowSet.AVBFlows {
		graph := network.GraphSet.Get(flow.Source, flow.Destinations)
		tree := SteninerTree(v2v, graph, flow.Source, flow.Destinations, cfg.Network.ByteRate)
		TreesSet.AVBRoutes = append(TreesSet.AVBRoutes, tree)
	}
	logger.Printf("Finish Steniner Tree %d AVB streams routing\n", len(TreesSet.AVBRoutes))

	return TreesSet
}

func Get_DistanceTree_Routing(network *network.Network, cfg *config.Config) *RouteSet {
	TreesSet := newRouteSet()

	for _, flow := range network.FlowSet.TSNFlows {
		graph := network.GraphSet.Get(flow.Source, flow.Destinations)
		tree := DistanceTree(graph, flow.Source, flow.Destinations, cfg.Network.ByteRate)
		TreesSet.TSNRoutes = append(TreesSet.TSNRoutes, tree)
	}
	logger.Printf("Finish Distance Tree %d TSN streams routing\n", len(TreesSet.TSNRoutes))

	for _, flow := range network.FlowSet.AVBFlows {
		graph := network.GraphSet.Get(flow.Source, flow.Destinations)
		tree := DistanceTree(graph, flow.Source, flow.Destinations, cfg.Network.ByteRate)
		TreesSet.AVBRoutes = append(TreesSet.AVBRoutes, tree)
	}
	logger.Printf("Finish Distance Tree %d AVB streams routing\n", len(TreesSet.AVBRoutes))

	return TreesSet
}

func Get_ShortestPath_Routing(network *network.Network, cfg *config.Config) *RouteSet {
	paths_set := newRouteSet()

	for _, flow := range network.FlowSet.TSNFlows {
		dest := flow.Destinations[0]
		graph := network.GraphSet.Get(flow.Source, flow.Destinations)
		path := ShortestPath(v2v, graph, flow.Source, dest, cfg.Network.ByteRate)
		paths_set.TSNRoutes = append(paths_set.TSNRoutes, path)
	}
	logger.Printf("Finish Shortest Path %d TSN streams routing\n", len(paths_set.TSNRoutes))

	for _, flow := range network.FlowSet.AVBFlows {
		dest := flow.Destinations[0]
		graph := network.GraphSet.Get(flow.Source, flow.Destinations)
		path := ShortestPath(v2v, graph, flow.Source, dest, cfg.Network.ByteRate)
		paths_set.AVBRoutes = append(paths_set.AVBRoutes, path)
	}
	logger.Printf("Finish Shortest Path %d AVB streams routing\n", len(paths_set.AVBRoutes))

	// CAN2TSN flows - encapsulated flows
	type sd struct{ s, d int }
	usedPath := make(map[sd]*Route)

	for _, method := range network.FlowSet.EncapsulateMethod {
		for _, flow := range method.CAN2TTFlows {
			dest := flow.Destinations[0]
			key := sd{flow.Source, dest}

			if existingPath, ok := usedPath[key]; ok {
				// Path already computed, clone and set method
				newPath := CloneRoute(existingPath)
				newPath.Method = method.MethodName
				paths_set.CAN2TTRoutes = append(paths_set.CAN2TTRoutes, newPath)
			} else {
				// Compute new path
				graph := network.GraphSet.Get(flow.Source, flow.Destinations)
				path := ShortestPath(v2v, graph, flow.Source, dest, cfg.Network.ByteRate)
				if path != nil {
					path.Method = method.MethodName
				}
				paths_set.CAN2TTRoutes = append(paths_set.CAN2TTRoutes, path)
				usedPath[key] = path
			}
		}
	}
	logger.Printf("Finish Shortest Path %d CAN2TSN streams routing\n", len(paths_set.CAN2TTRoutes))

	return paths_set
}

func (rs *RouteSet) InputRouteSet(bgTSN, bgAVB int) *RouteSet {
	out := newRouteSet()
	out.TSNRoutes = append(out.TSNRoutes, rs.TSNRoutes[bgTSN:]...)
	out.AVBRoutes = append(out.AVBRoutes, rs.AVBRoutes[bgAVB:]...)
	out.CAN2TTRoutes = append(out.CAN2TTRoutes, rs.CAN2TTRoutes...)
	return out
}

func (rs *RouteSet) BGRouteSet(bgTSN, bgAVB int) *RouteSet {
	out := newRouteSet()
	out.TSNRoutes = append(out.TSNRoutes, rs.TSNRoutes[:bgTSN]...)
	out.AVBRoutes = append(out.AVBRoutes, rs.AVBRoutes[:bgAVB]...)
	return out
}

func Get_KTree_Routing(network *network.Network, cfg *config.Config, SMT *RouteSet, K int, Method_Number int) *KRouteSet {
	ktrees_set := newKRouteSet()

	for nth, flow := range network.FlowSet.TSNFlows {
		Ktrees := KSpanningTree(v2v, SMT.TSNRoutes[nth], K, flow.Source, flow.Destinations, cfg.Network.ByteRate, Method_Number)
		ktrees_set.TSNRoutes = append(ktrees_set.TSNRoutes, Ktrees)
	}
	logger.Printf("Finish OSACO %d TSN streams routing\n", len(ktrees_set.TSNRoutes))

	for nth, flow := range network.FlowSet.AVBFlows {
		Ktrees := KSpanningTree(v2v, SMT.AVBRoutes[nth], K, flow.Source, flow.Destinations, cfg.Network.ByteRate, Method_Number)
		ktrees_set.AVBRoutes = append(ktrees_set.AVBRoutes, Ktrees)
	}
	logger.Printf("Finish OSACO %d AVB streams routing\n", len(ktrees_set.AVBRoutes))

	return ktrees_set
}

func Get_KPath_Routing(network *network.Network, cfg *config.Config, shortestPaths *RouteSet, K int) *KRouteSet {
	kpaths_set := newKRouteSet()

	// TSN flows
	for _, flow := range network.FlowSet.TSNFlows {
		dest := flow.Destinations[0]
		topo := network.GraphSet.Get(flow.Source, flow.Destinations)
		kpath := BuildKRoute(K, flow.Source, dest, topo, cfg.Network.ByteRate)
		kpaths_set.TSNRoutes = append(kpaths_set.TSNRoutes, kpath)
	}
	logger.Printf("Finish K-Path %d TSN streams routing (K=%d)\n", len(kpaths_set.TSNRoutes), K)

	// AVB flows
	for _, flow := range network.FlowSet.AVBFlows {
		dest := flow.Destinations[0]
		topo := network.GraphSet.Get(flow.Source, flow.Destinations)
		kpath := BuildKRoute(K, flow.Source, dest, topo, cfg.Network.ByteRate)
		kpaths_set.AVBRoutes = append(kpaths_set.AVBRoutes, kpath)
	}
	logger.Printf("Finish K-Path %d AVB streams routing (K=%d)\n", len(kpaths_set.AVBRoutes), K)

	// CAN2TSN flows - avoid redundant computation
	type sd struct{ s, d int }
	usedKPath := make(map[sd]*KRoute)

	for _, method := range network.FlowSet.EncapsulateMethod {
		for _, flow := range method.CAN2TTFlows {
			dest := flow.Destinations[0]
			key := sd{flow.Source, dest}

			if existingKPath, ok := usedKPath[key]; ok {
				// KPath already computed, clone and set method
				newKRoute := CloneKPath(existingKPath)
				newKRoute.Method = method.MethodName
				kpaths_set.CAN2TTRoutes = append(kpaths_set.CAN2TTRoutes, newKRoute)
			} else {
				// Compute new KPath
				topo := network.GraphSet.Get(flow.Source, flow.Destinations)
				kpath := BuildKRoute(K, flow.Source, dest, topo, cfg.Network.ByteRate)
				if kpath != nil {
					kpath.Method = method.MethodName
				}
				kpaths_set.CAN2TTRoutes = append(kpaths_set.CAN2TTRoutes, kpath)
				usedKPath[key] = kpath
			}
		}
	}
	logger.Printf("Finish K-Path %d CAN2TSN streams routing (K=%d)\n", len(kpaths_set.CAN2TTRoutes), K)

	return kpaths_set
}

func (krs *KRouteSet) InputKRouteSet(bgTSN, bgAVB int) *KRouteSet {
	out := newKRouteSet()
	out.TSNRoutes = append(out.TSNRoutes, krs.TSNRoutes[bgTSN:]...)
	out.AVBRoutes = append(out.AVBRoutes, krs.AVBRoutes[bgAVB:]...)
	out.CAN2TTRoutes = append(out.CAN2TTRoutes, krs.CAN2TTRoutes...)
	return out
}

func (krs *KRouteSet) BGKRouteSet(bgTSN, bgAVB int) *KRouteSet {
	out := newKRouteSet()
	out.TSNRoutes = append(out.TSNRoutes, krs.TSNRoutes[:bgTSN]...)
	out.AVBRoutes = append(out.AVBRoutes, krs.AVBRoutes[:bgAVB]...)
	return out
}
