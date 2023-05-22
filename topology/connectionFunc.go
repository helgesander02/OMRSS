package topology

import (
    "fmt"
)

func (c *Connection) IsValid(nodes []*Node) bool {
    for _, node := range nodes {
        if node.ID == c.ToNodeID {
            return true
        }
    }
    return false
}

func (topology *Topology) AddS2S(fromNodeID int, toNodeID int) {
    const cost float64 = 750000 // 750,000 bytes/6ms

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

    if connection1.IsValid(topology.Switch) {
        topology.Switch[fromNodeID].Connections = append(topology.Switch[fromNodeID].Connections, connection1)
		topology.Switch[toNodeID].Connections = append(topology.Switch[toNodeID].Connections, connection2)
    } 
}

func (topology *Topology) AddT2S(source int) {
    const cost float64 = 750000 // 750,000 bytes/6ms
    var toNodeID int

    if source>4 {
        toNodeID = 3
    } else {
        toNodeID = 0
    }
    
    connection := &Connection{
        FromNodeID: source+1000,
        ToNodeID:   toNodeID,
        Cost:       cost,
    }

    topology.Nodes[source].Connections = append(topology.Nodes[source].Connections, connection)
    topology.Nodes[source].ID = source+1000
    topology.Nodes[source].Name = fmt.Sprintf("Talker%d", source)

    topology.Talker = append(topology.Talker, topology.Nodes[source]) 

}

func (topology *Topology) AddS2L(destinations []int) {   
    const cost float64 = 750000 // 750,000 bytes/6ms 
    var fromNodeID int

    for i := 0; i < len(destinations); i++ {
        if destinations[i]>4 {
            fromNodeID = 3
        } else {
            fromNodeID = 0
        }

        connection := &Connection{
            FromNodeID: fromNodeID,
            ToNodeID:   destinations[i]+2000,
            Cost:       cost,
        }

        topology.Nodes[destinations[i]].Connections = append(topology.Nodes[destinations[i]].Connections, connection)
        topology.Nodes[destinations[i]].ID = destinations[i]+2000
        topology.Nodes[destinations[i]].Name = fmt.Sprintf("Listener%d", destinations[i]) 

        topology.Listener = append(topology.Listener, topology.Nodes[destinations[i]])
    } 
}

