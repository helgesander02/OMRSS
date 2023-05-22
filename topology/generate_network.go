package topology

func Generate_Network() *Topology {
    // 100 Mbps
    const (
        cost float64 = 750000 // 750,000 bytes/6ms
    )
    
    // Create Switch Node
    s1 := &Node{ID: 0, Name: "Switch0"}
    s2 := &Node{ID: 1, Name: "Switch1"}
    s3 := &Node{ID: 2, Name: "Switch2"}
    s4 := &Node{ID: 3, Name: "Switch3"}

    // Create Topology
    topology := &Topology{
        Switch: []*Node{s1, s2, s3, s4},
    }

    // Create null Nodes
    for i:=0; i<10; i++ {
        node := &Node{ID: i}
        topology.Nodes = append(topology.Nodes, node)
    }

    // Switch Connection
    topology.AddS2S(0, 1) // ConnectionFunc AddS2S
    topology.AddS2S(0, 2)
    topology.AddS2S(0, 3)
    topology.AddS2S(1, 3)
    topology.AddS2S(2, 3)    
    
    return topology
}