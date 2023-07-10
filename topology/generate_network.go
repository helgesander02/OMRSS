package topology

func Generate_Network(cost float64) *Topology {
    // Create Topology
    topology := &Topology{}

    // Create Switch Node
    // Switchs (Bridges) number
    for s:=0; s<12; s++ {
        snode := &Node{ID: s}
        topology.Switch = append(topology.Switch, snode)
    } 
    // Switch Connection
    Define_Switch_Connection(topology, cost)   

    // Create null Nodes 
    // END Devices (End Stations) number
    for i:=0; i<6; i++ {
        node := &Node{ID: i}
        topology.Nodes = append(topology.Nodes, node)
    }    
    
    return topology
}

// Define a topology where switches are interconnected with each other
func Define_Switch_Connection(topology *Topology, cost float64) {
    topology.AddS2S(0, 1, cost) 
    topology.AddS2S(0, 4, cost)
    topology.AddS2S(1, 2, cost)
    topology.AddS2S(1, 5, cost)
    topology.AddS2S(2, 3, cost)
    topology.AddS2S(2, 6, cost) 
    topology.AddS2S(3, 7, cost)
    topology.AddS2S(4, 5, cost)
    topology.AddS2S(4, 8, cost)
    topology.AddS2S(5, 6, cost)
    topology.AddS2S(5, 9, cost) 
    topology.AddS2S(6, 7, cost)
    topology.AddS2S(6, 10, cost)
    topology.AddS2S(8, 9, cost)
    topology.AddS2S(9, 10, cost) 
    topology.AddS2S(10, 11, cost)
    topology.AddS2S(11, 7, cost)
}

// Define a topology where null nodes are connected to switches
func Define_NullNodes_Connection(inputID int) int {
    var nodeID int
    if inputID==0 {
        nodeID = 0
    } else if inputID==1 {
        nodeID = 4
    } else if inputID==2 {
        nodeID = 8
    } else if inputID==3 {
        nodeID = 3
    } else if inputID==4 {
        nodeID = 7
    } else {
        nodeID = 11
    }

    return nodeID
}