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

func (topology *Topology) AddN2S2N(fromNodeID int, toNodeID int) {
 
}

func (topology *Topology) AddT2S(source int, cost float64) {
    var toNodeID int

    // 決定從哪開始
    //---------------
    if source==0 {
        toNodeID = 0
    } else if source==1 {
        toNodeID = 4
    } else if source==2 {
        toNodeID = 8
    } else if source==3 {
        toNodeID = 3
    } else if source==4 {
        toNodeID = 7
    } else {
        toNodeID = 11
    }
    //---------------
    
    connection := &Connection{
        FromNodeID: source+1000,
        ToNodeID:   toNodeID,
        Cost:       cost,
    }

    
    topology.Nodes[source].Connections = append(topology.Nodes[source].Connections, connection)
    topology.Nodes[source].ID = source+1000

    topology.Talker = append(topology.Talker, topology.Nodes[source]) 

}

func (topology *Topology) AddS2L(destinations []int, cost float64) {   
    var fromNodeID int

    for i := 0; i < len(destinations); i++ {
        // 決定從哪結束
        //---------------
        if destinations[i]==0 {
            fromNodeID = 0
        } else if destinations[i]==1 {
            fromNodeID = 4
        } else if destinations[i]==2 {
            fromNodeID = 8
        } else if destinations[i]==3 {
            fromNodeID = 3
        } else if destinations[i]==4 {
            fromNodeID = 7
        } else {
            fromNodeID = 11
        }
        //---------------

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