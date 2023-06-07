package topology

func Generate_Network(cost float64) *Topology {
    // Create Switch Node
    s1 := &Node{ID: 0}
    s2 := &Node{ID: 1}
    s3 := &Node{ID: 2}
    s4 := &Node{ID: 3}

    // Create Topology
    topology := &Topology{
        Switch: []*Node{s1, s2, s3, s4},
    }

    // Switch Connection
    topology.AddS2S(0, 1, cost) 
    topology.AddS2S(0, 2, cost)
    topology.AddS2S(0, 3, cost)
    topology.AddS2S(1, 3, cost)
    topology.AddS2S(2, 3, cost)    

    // Create null Nodes
    for i:=0; i<10; i++ {
        node := &Node{ID: i}
        topology.Nodes = append(topology.Nodes, node)
    }

    
    
    return topology
}