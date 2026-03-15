package topology

func (topology *Topology) AddS2S(fromNodeID int, toNodeID int, cost float64) {
	connection1 := NewLink(fromNodeID, toNodeID, cost)
	connection2 := NewLink(toNodeID, fromNodeID, cost)
	topology.Switch[fromNodeID].Links = append(topology.Switch[fromNodeID].Links, connection1)
	topology.Switch[toNodeID].Links = append(topology.Switch[toNodeID].Links, connection2)
}

func (topology *Topology) AddnullN2S(fromNodeID int, toNodeID int, cost float64) {
	connection1 := NewLink(fromNodeID, toNodeID, cost)
	topology.Nodes[fromNodeID%1000].Links = append(topology.Nodes[fromNodeID%1000].Links, connection1)
}

// Undirected Graph function
func (topology *Topology) AddN2S2N_For_Tree(source int, destinations []int, cost float64) {
	id := source % 1000
	fromNode := topology.GetNodeByID(id + 3000)
	fromNode.ID = source
	fromNode.Links[0].FromNodeID = source
	topology.Talker = append(topology.Talker, fromNode)

	toNodeID := fromNode.Links[0].ToNodeID
	connection2 := NewLink(toNodeID, source, cost)

	topology.Switch[toNodeID].Links = append(topology.Switch[toNodeID].Links, connection2)

	for i := 0; i < len(destinations); i++ {
		id := destinations[i] % 1000
		fromNode := topology.GetNodeByID(id + 3000)
		fromNode.ID = destinations[i]
		fromNode.Links[0].FromNodeID = destinations[i]
		topology.Listener = append(topology.Listener, fromNode)

		toNodeID := fromNode.Links[0].ToNodeID
		connection1 := NewLink(toNodeID, destinations[i], cost)
		topology.Switch[toNodeID].Links = append(topology.Switch[toNodeID].Links, connection1)
	}
}

func (topology *Topology) AddN2S2N_For_Path(source int, destinations int, cost float64) {
	sid := source % 1000
	sfromNode := topology.GetNodeByID(sid + 3000)
	sfromNode.ID = source
	sfromNode.Links[0].FromNodeID = source
	topology.Talker = append(topology.Talker, sfromNode)

	stoNodeID := sfromNode.Links[0].ToNodeID
	connection2 := NewLink(stoNodeID, source, cost)

	topology.Switch[stoNodeID].Links = append(topology.Switch[stoNodeID].Links, connection2)

	did := destinations % 1000
	dfromNode := topology.GetNodeByID(did + 3000)
	dfromNode.ID = destinations
	dfromNode.Links[0].FromNodeID = destinations
	topology.Listener = append(topology.Listener, dfromNode)

	dtoNodeID := dfromNode.Links[0].ToNodeID
	connection1 := NewLink(dtoNodeID, destinations, cost)
	topology.Switch[dtoNodeID].Links = append(topology.Switch[dtoNodeID].Links, connection1)

}

// Directed Graph function
func (topology *Topology) AddT2S(source int, cost float64) {
	id := source % 1000
	fromNode := topology.GetNodeByID(id + 3000)

	fromNode.ID = source
	fromNode.Links[0].FromNodeID = source
	topology.Talker = append(topology.Talker, fromNode)
}

func (topology *Topology) AddS2Ls(destinations []int, cost float64) {
	for i := 0; i < len(destinations); i++ {
		id := destinations[i] % 1000
		toNode := topology.GetNodeByID(id + 3000)
		toNode.ID = destinations[i]
		topology.Listener = append(topology.Listener, toNode)

		fromNodeID := toNode.Links[0].ToNodeID
		connection := NewLink(fromNodeID, destinations[i], cost)
		topology.Switch[fromNodeID].Links = append(topology.Switch[fromNodeID].Links, connection)
	}
}

func (topology *Topology) AddS2L(destination int, cost float64) {
	id := destination % 1000
	toNode := topology.GetNodeByID(id + 3000)
	toNode.ID = destination
	topology.Listener = append(topology.Listener, toNode)

	fromNodeID := toNode.Links[0].ToNodeID
	connection := NewLink(fromNodeID, destination, cost)
	topology.Switch[fromNodeID].Links = append(topology.Switch[fromNodeID].Links, connection)
}
