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

    // Create null Nodes 
    // END Devices (End Stations) number
    for i:=0; i<6; i++ {
        node := &Node{ID: i}
        topology.Nodes = append(topology.Nodes, node)
    }    
    
    return topology
}