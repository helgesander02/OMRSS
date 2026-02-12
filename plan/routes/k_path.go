package routes

import (
	"fmt"
	"src/network"
	"src/network/topology"
)

// Get_KPath_Routing computes K shortest paths for all flows in OSRO network
func Get_KPath_Routing(network *network.OSRO_Network, shortestPaths *Paths_set, K int) *KPaths_set {
	kpaths_set := new_KPaths_Set()

	// TSN flows
	for nth, flow := range network.FlowSet.TSNFlows {
		dest := flow.Destinations[0]
		kpath := BuildKPath(K, flow.Source, dest, network.Graph_Set.TSNGraphs[nth], network.BytesRate)
		kpaths_set.TSNPaths = append(kpaths_set.TSNPaths, kpath)
	}
	fmt.Printf("Finish K-Path %d TSN streams routing (K=%d)\n", len(kpaths_set.TSNPaths), K)

	// AVB flows
	for nth, flow := range network.FlowSet.AVBFlows {
		dest := flow.Destinations[0]
		kpath := BuildKPath(K, flow.Source, dest, network.Graph_Set.AVBGraphs[nth], network.BytesRate)
		kpaths_set.AVBPaths = append(kpaths_set.AVBPaths, kpath)
	}
	fmt.Printf("Finish K-Path %d AVB streams routing (K=%d)\n", len(kpaths_set.AVBPaths), K)

	// CAN2TSN flows - avoid redundant computation
	type sd struct{ s, d int }
	usedKPath := make(map[sd]*KPath)

	for _, method := range network.FlowSet.EncapsulateMethod {
		for _, flow := range method.CAN2TTFlows {
			key := sd{flow.Source, flow.Destination}

			if existingKPath, ok := usedKPath[key]; ok {
				// KPath already computed, clone and set method
				newKPath := CloneKPath(existingKPath)
				newKPath.Method = method.Method_Name
				kpaths_set.CAN2TSNPaths = append(kpaths_set.CAN2TSNPaths, newKPath)
			} else {
				// Compute new KPath
				topo := network.Graph_Set.GetGarphBySD(flow.Source, flow.Destination)
				kpath := BuildKPath(K, flow.Source, flow.Destination, topo, network.BytesRate)
				if kpath != nil {
					kpath.Method = method.Method_Name
				}
				kpaths_set.CAN2TSNPaths = append(kpaths_set.CAN2TSNPaths, kpath)
				usedKPath[key] = kpath
			}
		}
	}
	fmt.Printf("Finish K-Path %d CAN2TSN streams routing (K=%d)\n", len(kpaths_set.CAN2TSNPaths), K)

	return kpaths_set
}

// BuildKPath constructs K shortest paths between source and target
func BuildKPath(k int, src, dst int, topo *topology.Topology, cost float64) *KPath {
	// Build graph from topology
	graph := BuildGraphFromTopology(topo)

	// Use Yen's algorithm to find K shortest paths
	pathIDs := YenKPaths(graph, src, dst, k)

	// Create KPath structure
	kpath := new_KPath(k, src, dst)
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
func (kpaths_set *KPaths_set) InputKPathSet(bgTSN int, bgAVB int) *KPaths_set {
	input := new_KPaths_Set()
	input.TSNPaths = append(input.TSNPaths, kpaths_set.TSNPaths[bgTSN:]...)
	input.AVBPaths = append(input.AVBPaths, kpaths_set.AVBPaths[bgAVB:]...)
	input.CAN2TSNPaths = append(input.CAN2TSNPaths, kpaths_set.CAN2TSNPaths...)
	return input
}

// BGKPathSet splits KPaths into background KPaths
func (kpaths_set *KPaths_set) BGKPathSet(bgTSN int, bgAVB int) *KPaths_set {
	bg := new_KPaths_Set()
	bg.TSNPaths = append(bg.TSNPaths, kpaths_set.TSNPaths[:bgTSN]...)
	bg.AVBPaths = append(bg.AVBPaths, kpaths_set.AVBPaths[:bgAVB]...)
	return bg
}
