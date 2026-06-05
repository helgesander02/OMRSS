package routes

import (
	"src/network"
	"src/network/topology"
	"src/pkg/config"
	"src/pkg/logger"
)

// Get_KPath_Routing computes K shortest paths for all flows in OSRO network
func Get_KPath_Routing(network *network.Network, cfg *config.Config, shortestPaths *PathsSet, K int) *KPathsSet {
	kpaths_set := newKPathsSet()

	// TSN flows
	for nth, flow := range network.FlowSet.TSNFlows {
		dest := flow.Destinations[0]
		kpath := BuildKPath(K, flow.Source, dest, network.GraphSet.TSNGraphs[nth], cfg.Network.ByteRate)
		kpaths_set.TSNPaths = append(kpaths_set.TSNPaths, kpath)
	}
	logger.Printf("Finish K-Path %d TSN streams routing (K=%d)\n", len(kpaths_set.TSNPaths), K)

	// AVB flows
	for nth, flow := range network.FlowSet.AVBFlows {
		dest := flow.Destinations[0]
		kpath := BuildKPath(K, flow.Source, dest, network.GraphSet.AVBGraphs[nth], cfg.Network.ByteRate)
		kpaths_set.AVBPaths = append(kpaths_set.AVBPaths, kpath)
	}
	logger.Printf("Finish K-Path %d AVB streams routing (K=%d)\n", len(kpaths_set.AVBPaths), K)

	// CAN2TSN flows - avoid redundant computation
	type sd struct{ s, d int }
	usedKPath := make(map[sd]*KPath)

	for _, method := range network.FlowSet.EncapsulateMethod {
		for _, flow := range method.CAN2TTFlows {
			key := sd{flow.Source, flow.Destination}

			if existingKPath, ok := usedKPath[key]; ok {
				// KPath already computed, clone and set method
				newKPath := CloneKPath(existingKPath)
				newKPath.Method = method.MethodName
				kpaths_set.CAN2TSNPaths = append(kpaths_set.CAN2TSNPaths, newKPath)
			} else {
				// Compute new KPath
				topo := network.GraphSet.GetGarphBySD(flow.Source, flow.Destination)
				kpath := BuildKPath(K, flow.Source, flow.Destination, topo, cfg.Network.ByteRate)
				if kpath != nil {
					kpath.Method = method.MethodName
				}
				kpaths_set.CAN2TSNPaths = append(kpaths_set.CAN2TSNPaths, kpath)
				usedKPath[key] = kpath
			}
		}
	}
	logger.Printf("Finish K-Path %d CAN2TSN streams routing (K=%d)\n", len(kpaths_set.CAN2TSNPaths), K)

	return kpaths_set
}

// BuildKPath constructs K shortest paths between source and target
func BuildKPath(k int, src, dst int, topo *topology.Topology, cost float64) *KPath {
	// Build graph from topology
	graph := BuildGraphFromTopology(topo)

	// Use Yen's algorithm to find K shortest paths
	pathIDs := YenKPaths(graph, src, dst, k)

	// Create KPath structure
	kpath := newKPath(k, src, dst)
	for _, ids := range pathIDs {
		path := ConvertIDsToPath(ids, topo, cost)
		if path != nil {
			kpath.Paths = append(kpath.Paths, path)
		}
	}

	return kpath
}

// CloneKPath creates a deep copy of a KPath
func CloneKPath(src *KPath) *KPath {
	dup := &KPath{
		K:      src.K,
		Source: src.Source,
		Target: src.Target,
		Method: src.Method,
		Paths:  make([]*Path, len(src.Paths)),
	}

	for i, path := range src.Paths {
		dup.Paths[i] = ClonePath(path)
	}

	return dup
}

// InputKPathSet splits KPaths into input KPaths
func (kpaths_set *KPathsSet) InputKPathSet(bgTSN int, bgAVB int) *KPathsSet {
	input := newKPathsSet()
	input.TSNPaths = append(input.TSNPaths, kpaths_set.TSNPaths[bgTSN:]...)
	input.AVBPaths = append(input.AVBPaths, kpaths_set.AVBPaths[bgAVB:]...)
	input.CAN2TSNPaths = append(input.CAN2TSNPaths, kpaths_set.CAN2TSNPaths...)
	return input
}

// BGKPathSet splits KPaths into background KPaths
func (kpaths_set *KPathsSet) BGKPathSet(bgTSN int, bgAVB int) *KPathsSet {
	bg := newKPathsSet()
	bg.TSNPaths = append(bg.TSNPaths, kpaths_set.TSNPaths[:bgTSN]...)
	bg.AVBPaths = append(bg.AVBPaths, kpaths_set.AVBPaths[:bgAVB]...)
	return bg
}
