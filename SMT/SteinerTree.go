package SMT

import (    
    "src/topology"
    "src/stream"
    "fmt"
)

func SteninerTree(topology *topology.Topology, flows *stream.Flows) {
    for _, flow := range flows.TSNFlows {
        Feature(topology, flow)

        break
    }

    //for _, flow := range flows.AVBFlows {
        //Feature(topology, flow)
    //}
}

func Feature(topology *topology.Topology, flow *stream.Flow) {
    t := topology
    t.AddT2S(flow.Source)
    t.AddS2L(flow.Destinations)

    for _, terminal := range flow.Destinations {
        fmt.Printf("%d --> %d Shortest Path:\n", flow.Source, terminal)
        graphs := AddGarph(t, terminal+2000)
        graphs = Dijkstra(graphs, flow.Source+1000, terminal+2000)
        graphs.Show_Graph()
    }
}

func AddGarph(topology *topology.Topology, terminal int) *Graphs {
    graphs := &Graphs{}
    //Talker
    for _, t := range topology.Talker {
        gt := &Graph{}
        gt.ID = t.ID
        gt.AddEdge(t.Connections)
        graphs.Graphs = append(graphs.Graphs, gt)
    }

    //Switch
    for _, s := range topology.Switch {
        gs := &Graph{}
        gs.ID = s.ID
        gs.AddEdge(s.Connections)
        graphs.Graphs = append(graphs.Graphs, gs)
    }

    //Listener 
    gl := &Graph{}
    gl.ID = terminal
    graphs.Graphs = append(graphs.Graphs, gl)

    for _, l := range topology.Listener {
        if l.ID == terminal {
            if (terminal-2000)<5{
                graph := graphs.findGraph(0)
                graph.AddEdge(l.Connections)
            } else {
                graph := graphs.findGraph(3)
                graph.AddEdge(l.Connections)
            } 
        }        
    }
    return graphs
}

func (graph *Graph) AddEdge(connections []*topology.Connection){
    for _, c := range connections {
        edge := &Edge{
            Strat: c.FromNodeID,
            End: c.ToNodeID,
            Cost: 1,
        }
        graph.Edges = append(graph.Edges, edge)  
    }
}

func (graphs *Graphs) findGraph(id int) *Graph {
    for _, graph := range graphs.Graphs {
        if graph.ID == id {
            return graph
        }
    }
    return &Graph{}
}

