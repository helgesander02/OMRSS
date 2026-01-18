package topology

func (topology *Topology) AddS2S(fromNodeID int, toNodeID int, cost float64) {
	connection1 := new_Connection(fromNodeID, toNodeID, cost)
	connection2 := new_Connection(toNodeID, fromNodeID, cost)
	topology.Switch[fromNodeID].Connections = append(topology.Switch[fromNodeID].Connections, connection1)
	topology.Switch[toNodeID].Connections = append(topology.Switch[toNodeID].Connections, connection2)
}

func (topology *Topology) AddnullN2S(fromNodeID int, toNodeID int, cost float64) {
	connection1 := new_Connection(fromNodeID, toNodeID, cost)
	topology.Nodes[fromNodeID%1000].Connections = append(topology.Nodes[fromNodeID%1000].Connections, connection1)
}

// Undirected Graph function
func (topology *Topology) AddN2S2N_For_Tree(source int, destinations []int, cost float64) {
	id := source % 1000
	fromNode := topology.GetNodeByID(id + 3000)
	fromNode.ID = source
	fromNode.Connections[0].FromNodeID = source
	topology.Talker = append(topology.Talker, fromNode)

	toNodeID := fromNode.Connections[0].ToNodeID
	connection2 := new_Connection(toNodeID, source, cost)

	topology.Switch[toNodeID].Connections = append(topology.Switch[toNodeID].Connections, connection2)

	for i := 0; i < len(destinations); i++ {
		id := destinations[i] % 1000
		fromNode := topology.GetNodeByID(id + 3000)
		fromNode.ID = destinations[i]
		fromNode.Connections[0].FromNodeID = destinations[i]
		topology.Listener = append(topology.Listener, fromNode)

		toNodeID := fromNode.Connections[0].ToNodeID
		connection1 := new_Connection(toNodeID, destinations[i], cost)
		topology.Switch[toNodeID].Connections = append(topology.Switch[toNodeID].Connections, connection1)
	}
}

func (topology *Topology) AddN2S2N_For_Path(source int, destinations int, cost float64) {

	sid := source % 1000
	sfromNode := topology.GetNodeByID(sid + 3000)
	sfromNode.ID = source
	sfromNode.Connections[0].FromNodeID = source
	topology.Talker = append(topology.Talker, sfromNode)

	stoNodeID := sfromNode.Connections[0].ToNodeID
	connection2 := new_Connection(stoNodeID, source, cost)

	topology.Switch[stoNodeID].Connections = append(topology.Switch[stoNodeID].Connections, connection2)

	did := destinations % 1000
	dfromNode := topology.GetNodeByID(did + 3000)
	dfromNode.ID = destinations
	dfromNode.Connections[0].FromNodeID = destinations
	topology.Listener = append(topology.Listener, dfromNode)

	dtoNodeID := dfromNode.Connections[0].ToNodeID
	connection1 := new_Connection(dtoNodeID, destinations, cost)
	topology.Switch[dtoNodeID].Connections = append(topology.Switch[dtoNodeID].Connections, connection1)

}

// Directed Graph function
func (topology *Topology) AddT2S(source int, cost float64) {
	id := source % 1000
	fromNode := topology.GetNodeByID(id + 3000)

	fromNode.ID = source
	fromNode.Connections[0].FromNodeID = source
	topology.Talker = append(topology.Talker, fromNode)
}

func (topology *Topology) AddS2Ls(destinations []int, cost float64) {
	for i := 0; i < len(destinations); i++ {
		id := destinations[i] % 1000
		toNode := topology.GetNodeByID(id + 3000)
		toNode.ID = destinations[i]
		topology.Listener = append(topology.Listener, toNode)

		fromNodeID := toNode.Connections[0].ToNodeID
		connection := new_Connection(fromNodeID, destinations[i], cost)
		topology.Switch[fromNodeID].Connections = append(topology.Switch[fromNodeID].Connections, connection)
	}
}

func (topology *Topology) AddS2L(destination int, cost float64) {
	id := destination % 1000
	toNode := topology.GetNodeByID(id + 3000)
	toNode.ID = destination
	topology.Listener = append(topology.Listener, toNode)

	fromNodeID := toNode.Connections[0].ToNodeID
	connection := new_Connection(fromNodeID, destination, cost)
	topology.Switch[fromNodeID].Connections = append(topology.Switch[fromNodeID].Connections, connection)
}
