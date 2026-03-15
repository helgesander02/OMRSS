package routes

import (
	"fmt"
	"src/internal/config"
	"src/network"
	"src/network/topology"
)

var v2v_path *V2V = &V2V{} // v2v for shortest path calculation

// Get_ShortestPath_Routing computes shortest paths for all flows in OSRO network
func Get_ShortestPath_Routing(network *network.Network, cfg *config.Config) *PathsSet {
	paths_set := newPathsSet()

	// TSN flows - point-to-point (source to first destination)
	// Note: For OSRO, we use point-to-point paths, not multicast trees
	for nth, flow := range network.FlowSet.TSNFlows {
		dest := flow.Destinations[0] // Use first destination for point-to-point
		path := ShortestPath(v2v_path, network.GraphSet.TSNGraphs[nth], flow.Source, dest, cfg.Network.ByteRate)
		paths_set.TSNPaths = append(paths_set.TSNPaths, path)
	}
	fmt.Printf("Finish Shortest Path %d TSN streams routing\n", len(paths_set.TSNPaths))

	// AVB flows - point-to-point
	for nth, flow := range network.FlowSet.AVBFlows {
		dest := flow.Destinations[0] // Use first destination for point-to-point
		path := ShortestPath(v2v_path, network.GraphSet.AVBGraphs[nth], flow.Source, dest, cfg.Network.ByteRate)
		paths_set.AVBPaths = append(paths_set.AVBPaths, path)
	}
	fmt.Printf("Finish Shortest Path %d AVB streams routing\n", len(paths_set.AVBPaths))

	// CAN2TSN flows - encapsulated flows
	type sd struct{ s, d int }
	usedPath := make(map[sd]*Path)

	for _, method := range network.FlowSet.EncapsulateMethod {
		for _, flow := range method.CAN2TTFlows {
			key := sd{flow.Source, flow.Destination}

			if existingPath, ok := usedPath[key]; ok {
				// Path already computed, clone and set method
				newPath := ClonePath(existingPath)
				newPath.Method = method.MethodName
				paths_set.CAN2TSNPaths = append(paths_set.CAN2TSNPaths, newPath)
			} else {
				// Compute new path
				topo := network.GraphSet.GetGarphBySD(flow.Source, flow.Destination)
				path := ShortestPath(v2v_path, topo, flow.Source, flow.Destination, cfg.Network.ByteRate)
				if path != nil {
					path.Method = method.MethodName
				}
				paths_set.CAN2TSNPaths = append(paths_set.CAN2TSNPaths, path)
				usedPath[key] = path
			}
		}
	}
	fmt.Printf("Finish Shortest Path %d CAN2TSN streams routing\n", len(paths_set.CAN2TSNPaths))

	return paths_set
}

// ShortestPath finds the shortest path between source and destination
func ShortestPath(v2v *V2V, topo *topology.Topology, source int, destination int, cost float64) *Path {
	// Get or compute shortest paths
	v2vedge, firstGet := v2v.GetV2VEdge(source)

	if !v2vedge.InV2VEdge(destination) {
		// Need to compute the path
		graph := GetGarph(topo)
		graph.ToVertex = destination
		graph = Dijkstra(graph, source, destination)
		v2vedge.Graphs = append(v2vedge.Graphs, graph)
	}

	if firstGet {
		v2v.V2VEdges = append(v2v.V2VEdges, v2vedge)
	}

	// Get all shortest paths
	paths := v2vedge.GetV2VPath(destination)
	if len(paths) == 0 {
		return nil
	}

	// Convert first path (shortest) to Path structure
	path := ConvertIDsToPath(paths[0], topo, cost)
	return path
}

// ConvertIDsToPath converts node IDs to Path structure
func ConvertIDsToPath(ids []int, topo *topology.Topology, cost float64) *Path {
	path := newPath()
	path.IDs = ids
	path.Weight = 0

	for i, id := range ids {
		realNode := topo.GetNodeByID(id)
		if realNode == nil {
			fmt.Printf("Warning: Node with ID=%d not found in topology\n", id)
			continue
		}

		node := &PathNode{
			ID:          realNode.ID,
			Shape:       "",
			Connections: make([]*PathConnection, 0),
		}

		// Add forward connection
		if i < len(ids)-1 {
			nextID := ids[i+1]
			conn := findConnectionInTopology(realNode, nextID)
			if conn != nil {
				node.Connections = append(node.Connections, &PathConnection{
					FromNodeID: conn.FromNodeID,
					ToNodeID:   conn.ToNodeID,
					Cost:       cost,
				})
				path.Weight += 1 // hop count
			}
		}

		// Add backward connection (for bidirectional representation)
		if i > 0 {
			prevID := ids[i-1]
			conn := findConnectionInTopology(realNode, prevID)
			if conn != nil {
				node.Connections = append(node.Connections, &PathConnection{
					FromNodeID: conn.FromNodeID,
					ToNodeID:   conn.ToNodeID,
					Cost:       cost,
				})
			}
		}

		path.Nodes = append(path.Nodes, node)
	}

	return path
}

// findConnectionInTopology finds connection from node to target
func findConnectionInTopology(node *topology.Node, toID int) *topology.Link {
	for _, conn := range node.Links {
		if conn.ToNodeID == toID {
			return conn
		}
	}
	return nil
}

// ClonePath creates a deep copy of a path
func ClonePath(src *Path) *Path {
	dup := &Path{
		Method: src.Method,
		IDs:    append([]int{}, src.IDs...),
		Nodes:  make([]*PathNode, len(src.Nodes)),
		Weight: src.Weight,
	}

	for i, node := range src.Nodes {
		newNode := &PathNode{
			ID:          node.ID,
			Shape:       node.Shape,
			Connections: make([]*PathConnection, len(node.Connections)),
		}
		for j, conn := range node.Connections {
			newNode.Connections[j] = &PathConnection{
				FromNodeID: conn.FromNodeID,
				ToNodeID:   conn.ToNodeID,
				Cost:       conn.Cost,
			}
		}
		dup.Nodes[i] = newNode
	}

	return dup
}

// InputPathSet splits paths into input paths
func (paths_set *PathsSet) InputPathSet(bgTSN int, bgAVB int) *PathsSet {
	input := newPathsSet()
	input.TSNPaths = append(input.TSNPaths, paths_set.TSNPaths[bgTSN:]...)
	input.AVBPaths = append(input.AVBPaths, paths_set.AVBPaths[bgAVB:]...)
	input.CAN2TSNPaths = append(input.CAN2TSNPaths, paths_set.CAN2TSNPaths...)
	return input
}

// BGPathSet splits paths into background paths
func (paths_set *PathsSet) BGPathSet(bgTSN int, bgAVB int) *PathsSet {
	bg := newPathsSet()
	bg.TSNPaths = append(bg.TSNPaths, paths_set.TSNPaths[:bgTSN]...)
	bg.AVBPaths = append(bg.AVBPaths, paths_set.AVBPaths[:bgAVB]...)
	return bg
}
