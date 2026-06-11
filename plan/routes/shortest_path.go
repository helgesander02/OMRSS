package routes

import (
	"src/network/topology"
	"src/pkg/logger"
)

func ShortestPath(v2v *V2V, topo *topology.Topology, source int, destination int, cost float64) *Route {
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

func ConvertIDsToPath(ids []int, topo *topology.Topology, cost float64) *Route {
	path := newRoute()
	path.IDs = ids
	path.Weight = 0

	for i, id := range ids {
		realNode := topo.GetNodeByID(id)
		if realNode == nil {
			logger.Printf("Warning: Node with ID=%d not found in topology\n", id)
			continue
		}

		node := &Node{
			ID:          realNode.ID,
			Connections: make([]*Connection, 0),
		}

		// Add forward connection
		if i < len(ids)-1 {
			nextID := ids[i+1]
			conn := findConnectionInTopology(realNode, nextID)
			if conn != nil {
				node.Connections = append(node.Connections, &Connection{
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
				node.Connections = append(node.Connections, &Connection{
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

func findConnectionInTopology(node *topology.Node, toID int) *topology.Link {
	for _, conn := range node.Links {
		if conn.ToNodeID == toID {
			return conn
		}
	}
	return nil
}

func CloneRoute(src *Route) *Route {
	dup := &Route{
		Method: src.Method,
		IDs:    append([]int{}, src.IDs...),
		Nodes:  make([]*Node, len(src.Nodes)),
		Weight: src.Weight,
	}

	for i, node := range src.Nodes {
		newNode := &Node{
			ID:          node.ID,
			Connections: make([]*Connection, len(node.Connections)),
		}
		for j, conn := range node.Connections {
			newNode.Connections[j] = &Connection{
				FromNodeID: conn.FromNodeID,
				ToNodeID:   conn.ToNodeID,
				Cost:       conn.Cost,
			}
		}
		dup.Nodes[i] = newNode
	}

	return dup
}
