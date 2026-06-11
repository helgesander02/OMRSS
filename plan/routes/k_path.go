package routes

import (
	"src/network/topology"
)

// BuildKRoute constructs K shortest paths between source and target
func BuildKRoute(k int, src, dst int, topo *topology.Topology, cost float64) *KRoute {
	// Build graph from topology
	graph := BuildGraphFromTopology(topo)

	// Use Yen's algorithm to find K shortest paths
	pathIDs := YenKPaths(graph, src, dst, k)

	// Create KPath structure
	kpath := newKRoute(k, src, dst)
	for _, ids := range pathIDs {
		path := ConvertIDsToPath(ids, topo, cost)
		if path != nil {
			kpath.Routes = append(kpath.Routes, path)
		}
	}

	return kpath
}

// CloneKPath creates a deep copy of a KPath
func CloneKPath(src *KRoute) *KRoute {
	dup := &KRoute{
		K:      src.K,
		Source: src.Source,
		Target: src.Target,
		Method: src.Method,
		Routes: make([]*Route, len(src.Routes)),
	}

	for i, path := range src.Routes {
		dup.Routes[i] = CloneRoute(path)
	}

	return dup
}
