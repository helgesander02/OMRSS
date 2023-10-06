package network

func (topology *Topology) AddS2S(fromNodeID int, toNodeID int, cost float64) {
	connection1 := &Connection{
		FromNodeID: fromNodeID,
		ToNodeID:   toNodeID,
		Cost:       cost,
	}
	connection2 := &Connection{
		FromNodeID: toNodeID,
		ToNodeID:   fromNodeID,
		Cost:       cost,
	}

	topology.Switch[fromNodeID].Connections = append(topology.Switch[fromNodeID].Connections, connection1)
	topology.Switch[toNodeID].Connections = append(topology.Switch[toNodeID].Connections, connection2)
}

// Undirected Graph function
func (topology *Topology) AddN2S2N(source int, destinations []int, cost float64) {
	id := source % 1000
	// 決定從哪開始
	toNodeID := Define_NullNodes_Connection(source - 1000)

	connection1 := &Connection{
		FromNodeID: source,
		ToNodeID:   toNodeID,
		Cost:       cost,
	}

	connection2 := &Connection{
		FromNodeID: toNodeID,
		ToNodeID:   source,
		Cost:       cost,
	}

	topology.Nodes[id].Connections = append(topology.Nodes[id].Connections, connection1)
	topology.Nodes[id].ID = source
	topology.Talker = append(topology.Talker, topology.Nodes[id])
	topology.Switch[toNodeID].Connections = append(topology.Switch[toNodeID].Connections, connection2)

	for i := 0; i < len(destinations); i++ {
		id := destinations[i] % 1000
		// 決定從哪結束
		fromNodeID := Define_NullNodes_Connection(destinations[i] - 2000)

		connection1 := &Connection{
			FromNodeID: fromNodeID,
			ToNodeID:   destinations[i],
			Cost:       cost,
		}

		connection2 := &Connection{
			FromNodeID: destinations[i],
			ToNodeID:   fromNodeID,
			Cost:       cost,
		}

		for _, node := range topology.Switch {
			if node.ID == fromNodeID {
				node.Connections = append(node.Connections, connection1)
				break
			}
		}
		topology.Nodes[id].ID = destinations[i]
		topology.Nodes[id].Connections = append(topology.Nodes[id].Connections, connection2)
		topology.Listener = append(topology.Listener, topology.Nodes[id])
	}
}

// Directed Graph function
func (topology *Topology) AddT2S(source int, cost float64) {
	id := source % 1000
	// 決定從哪開始
	toNodeID := Define_NullNodes_Connection(source - 1000)

	connection := &Connection{
		FromNodeID: source,
		ToNodeID:   toNodeID,
		Cost:       cost,
	}

	topology.Nodes[id].Connections = append(topology.Nodes[id].Connections, connection)
	topology.Nodes[id].ID = source

	topology.Talker = append(topology.Talker, topology.Nodes[id])
}

// Directed Graph function
func (topology *Topology) AddS2L(destinations []int, cost float64) {
	for i := 0; i < len(destinations); i++ {
		id := destinations[i] % 1000
		// 決定從哪結束
		fromNodeID := Define_NullNodes_Connection(destinations[i] - 2000)

		connection := &Connection{
			FromNodeID: fromNodeID,
			ToNodeID:   destinations[i],
			Cost:       cost,
		}

		for _, node := range topology.Switch {
			if node.ID == fromNodeID {
				node.Connections = append(node.Connections, connection)
				break
			}
		}
		topology.Nodes[id].ID = destinations[i]
		topology.Listener = append(topology.Listener, topology.Nodes[id])
	}
}
