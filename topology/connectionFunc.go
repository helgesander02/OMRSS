package topology

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
    // 決定從哪開始
    var toNodeID int
    toNodeID = Define_NullNodes_Connection(source)
    
    connection1 := &Connection{
        FromNodeID: source+1000,
        ToNodeID:   toNodeID,
        Cost:       cost,
    }

    connection2 := &Connection{
        FromNodeID: toNodeID,
        ToNodeID:   source+1000,
        Cost:       cost,
    }
   
    topology.Nodes[source].Connections = append(topology.Nodes[source].Connections, connection1)
    topology.Nodes[source].ID = source+1000
    topology.Talker = append(topology.Talker, topology.Nodes[source])
    topology.Switch[toNodeID].Connections = append(topology.Switch[toNodeID].Connections, connection2)    

    for i := 0; i < len(destinations); i++ {
        // 決定從哪結束
        var fromNodeID int
        fromNodeID = Define_NullNodes_Connection(destinations[i])

        connection1 := &Connection{
            FromNodeID: fromNodeID,
            ToNodeID:   destinations[i]+2000,
            Cost:       cost,
        }

        connection2 := &Connection{
            FromNodeID: destinations[i]+2000,
            ToNodeID:   fromNodeID,
            Cost:       cost,
        }

        for _, node := range topology.Switch {
            if node.ID == fromNodeID {
                node.Connections = append(node.Connections, connection1)
                break
            }
        }     
        topology.Nodes[destinations[i]].ID = destinations[i]+2000
        topology.Nodes[destinations[i]].Connections = append(topology.Nodes[destinations[i]].Connections, connection2)
        topology.Listener = append(topology.Listener, topology.Nodes[destinations[i]])     
    } 
}

// Directed Graph function
func (topology *Topology) AddT2S(source int,  cost float64) {
    // 決定從哪開始
    var toNodeID int
    toNodeID = Define_NullNodes_Connection(source)
    
    connection := &Connection{
        FromNodeID: source+1000,
        ToNodeID:   toNodeID,
        Cost:       cost,
    }
    
    topology.Nodes[source].Connections = append(topology.Nodes[source].Connections, connection)
    topology.Nodes[source].ID = source+1000

    topology.Talker = append(topology.Talker, topology.Nodes[source])
}

// Directed Graph function
func (topology *Topology) AddS2L(destinations []int, cost float64) {   
    for i := 0; i < len(destinations); i++ {
        // 決定從哪結束
        var fromNodeID int
        fromNodeID = Define_NullNodes_Connection(destinations[i])

        connection := &Connection{
            FromNodeID: fromNodeID,
            ToNodeID:   destinations[i]+2000,
            Cost:       cost,
        }

        for _, node := range topology.Switch {
            if node.ID == fromNodeID {
                node.Connections = append(node.Connections, connection)
                break
            }
        }     
        topology.Nodes[destinations[i]].ID = destinations[i]+2000
        topology.Listener = append(topology.Listener, topology.Nodes[destinations[i]])     
    } 
}