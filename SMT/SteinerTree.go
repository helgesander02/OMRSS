package SMT

import (   
    "fmt"
    "math"
    "encoding/json"
    
    "src/topology"
    "src/flow"
)

// Bang Ye Wu, Kun-Mao Chao, "Steiner Minimal Trees"
func SteninerTree(topology *topology.Topology, flows *flow.Flows, cost float64) *Trees {
    trees := &Trees{}
    
    for _, flow := range flows.TSNFlows {
        tree := GetTree(topology, flow, cost)
        trees.TSNTrees = append(trees.TSNTrees, tree)
    }
    fmt.Println("Finish TSN all SteninerTree")

    for _, flow := range flows.AVBFlows {
        tree := GetTree(topology, flow, cost)
        trees.AVBTrees = append(trees.AVBTrees, tree)
    }
    fmt.Println("Finish AVB all SteninerTree")

    return trees
}

func GetTree(topology *topology.Topology, flow *flow.Flow, cost float64) *Tree {
    t := DeepCopy(topology) // Duplicate of Topology
    t.AddN2S2N(flow.Source, flow.Destinations, cost) // Undirected Graph

    var Terminal []int
    Terminal = append(Terminal, flow.Source+1000)
    for _, d := range flow.Destinations {
        Terminal = append(Terminal, d+2000)
    }
    
    // ---------------------------------------------------------------------------------------------------
    // Determine if there is a vertex in the tree. 
    // If not, find all shortest paths between terminals and select the path with the minimum cost. 
    // Then, add all the vertices from this shortest path to the tree. 
    // If there is a vertex, find all shortest paths between the vertex and terminals. 
    // Choose the path with the minimum cost and add all the vertices from this shortest path to the tree.
    // ---------------------------------------------------------------------------------------------------
    tree := &Tree{}
    var used_tmal []int
    for len(used_tmal) != len(Terminal) {
        if len(tree.Nodes) == 0 {
            // Find the set of all shortest paths from Source to Destinations
            Vertexs_P := make(map[int]*Graph)
            for _, terminal := range Terminal {  
                tmalcount := math.MaxInt8
                for _, tmal :=  range Terminal {
                    graph := GetGarph(t)
                    graph = Dijkstra(graph, terminal, tmal)
                    
                    if graph.Count == 0 { continue }
                    if tmalcount > graph.Count {
                        tmalcount = graph.Count
                        Vertexs_P[terminal] = graph
                    }
                }                
            }
            
            // Add the shortest path with the minimum cost to the tree in sequence
            count := math.MaxInt8
            var path [][]int
            for _, v := range Vertexs_P {
                if count > v.Count {
                    count = v.Count
                    path = v.Path
                }
            }
            b := tree.AddTree(path, cost)
            if b { 
                used_tmal = append(used_tmal, path[0][0]) 
                used_tmal = append(used_tmal, path[0][len(path[0])-1])
            }

        } else {      
            for _, terminal := range Terminal {  
                if TerminalBeenUsed(used_tmal, terminal) { continue }
                // Find the set of all shortest paths from Vertex to Destinations
                Vertexs_P := make(map[int]*Graph)
                nodecount := math.MaxInt8
                for _, node := range tree.Nodes {
                    graph := GetGarph(t)
                    graph = Dijkstra(graph, node.ID, terminal)
                    if graph.Count == 0 { continue }
                    if nodecount > graph.Count {
                        nodecount = graph.Count
                        Vertexs_P[terminal] = graph
                    }
                }
                // Add the shortest path with the minimum cost to the tree in sequence
                for k, v := range Vertexs_P {
                    if len(used_tmal)+1 == len(Terminal) {
                        var path [][]int
                        path = append(path, v.Path[0])
                        b := tree.AddTree(path, cost)
                        if b { used_tmal = append(used_tmal, k) }

                    } else {
                        b := tree.AddTree(v.Path, cost)
                        if b { used_tmal = append(used_tmal, k) }

                    }  
                }
            }        
        } 
    }

    return tree
}

func TerminalBeenUsed(A []int, b int) bool {
    for _, a := range A {
        if a == b {
            return true
        }
    }
    return false
}

func (tree *Tree) AddTree(P [][]int, cost float64) bool {   
    if len(P) == 1 {
        //fmt.Println(P, "true")
        for l:=len(P[0])-1; l>0; l-- {
            node1, b1 := tree.InTree(P[0][l])
            if !(b1) { tree.Nodes = append(tree.Nodes, node1) }  
            connection1 := &Connection{
                FromNodeID: P[0][l],
                ToNodeID: P[0][l-1],
                Cost: cost,
            }
            node1.Connections = append(node1.Connections, connection1)         

            node2, b2 := tree.InTree(P[0][l-1]) 
            if !(b2) { tree.Nodes = append(tree.Nodes, node2) }
            connection2 := &Connection{
                FromNodeID: P[0][l-1],
                ToNodeID: P[0][l],
                Cost: cost,
            }
            node2.Connections = append(node2.Connections, connection2)           
        }
        return true

    } else {
        //fmt.Println(P, "false")
        return false
    }
} 
    
func (tree *Tree) InTree(id int) (*Node, bool) {
    for _, node := range tree.Nodes{
        if node.ID == id  {
            return node, true
        }
    }
    return &Node{ID: id}, false
}

func GetGarph(topology *topology.Topology) *Graph {
    graph := &Graph{}
    //Talker
    for _, t := range topology.Talker {
        gt := &Vertex{}
        gt.ID = t.ID
        gt.AddEdge(t.Connections)
        graph.Vertexs = append(graph.Vertexs, gt)
    }

    //Switch
    for _, s := range topology.Switch {
        gs := &Vertex{}
        gs.ID = s.ID
        gs.AddEdge(s.Connections)
        graph.Vertexs = append(graph.Vertexs, gs)
    }

    //Listener 
    for _, l := range topology.Listener {
        gl := &Vertex{}
        gl.ID = l.ID
        gl.AddEdge(l.Connections)
        graph.Vertexs = append(graph.Vertexs, gl)
    }

    return graph
}

func (vertex *Vertex) AddEdge(connections []*topology.Connection){
    for _, c := range connections {
        edge := &Edge{
            Strat: c.FromNodeID,
            End: c.ToNodeID,
            Cost: 1,
        }
        vertex.Edges = append(vertex.Edges, edge)  
    }
}

func DeepCopy(t1 *topology.Topology) (*topology.Topology) {
    if buf, err := json.Marshal(t1); err != nil {
        return nil
    } else {
        t2 := &topology.Topology{}
        if err = json.Unmarshal(buf, t2); err != nil {
            return nil
        }
        return t2
    }
}
