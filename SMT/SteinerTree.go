package SMT

import (   
    "fmt"
    "math"
    "encoding/json"
    
    "src/topology"
    "src/stream"
)

// Bang Ye Wu, Kun-Mao Chao, "Steiner Minimal Trees"
func SteninerTree(topology *topology.Topology, flows *stream.Flows, cost float64) *Trees {
    trees := &Trees{}
    
    for _, flow := range flows.TSNFlows {
        tree := GetTree(topology, flow, cost)
        trees.TSNTrees = append(trees.TSNTrees, tree)
        break
    }
    fmt.Println("Finish TSN all SteninerTree")

    for _, flow := range flows.AVBFlows {
        tree := GetTree(topology, flow, cost)
        trees.AVBTrees = append(trees.AVBTrees, tree)
        break
    }
    fmt.Println("Finish AVB all SteninerTree")

    return trees
}

func GetTree(topology *topology.Topology, flow *stream.Flow, cost float64) *Tree {
    // Finish Topology
    t := DeepCopy(topology) // Duplicate of Topology
    t.AddT2S(flow.Source, cost)
    t.AddS2L(flow.Destinations, cost)
    
    tree := &Tree{}
    var used_t []int
    for len(used_t) != len(flow.Destinations) {
        if len(tree.Nodes) == 0 {
            // Find the set of all shortest paths from Source to Destinations
            Vertexs_P := make(map[int]*Graph)
            for _, terminal := range flow.Destinations {  
                graph := GetGarph(t)
                graph = Dijkstra(graph, flow.Source+1000, terminal+2000) 

                Vertexs_P[terminal+2000] = graph
            }
            
            // Add the shortest path with the minimum cost to the tree in sequence
            var path [][]int  
            var ut int    
            count := math.MaxInt8
            
            for k, v := range Vertexs_P {
                if v.Count < count {
                    count = v.Count
                    path = v.Path
                    ut = k
                }
            }
            
            b := tree.AddTree(path, cost)
            if b { used_t = append(used_t, ut) } 

        } else {      
            for _, terminal := range flow.Destinations {  
                if InUsed(used_t, terminal+2000) { continue }
                // Find the set of all shortest paths from Vertex to Destinations
                Vertexs_P := make(map[int]*Graph)
                nodecount := math.MaxInt8
                for _, node := range tree.Nodes {
                    graph := GetGarph(t)
                    graph = Dijkstra(graph, node.ID, terminal+2000)
                    if graph.Count == 0 { continue }
                    if nodecount > graph.Count {
                        nodecount = graph.Count
                        Vertexs_P[terminal+2000] = graph
                    }
                }
                // Add the shortest path with the minimum cost to the tree in sequence
                for k, v := range Vertexs_P {
                    b := tree.AddTree(v.Path, cost)
                    if b { used_t = append(used_t, k) }
                }
            }        
        }    
    }

    return tree
}

func InUsed(A []int, b int) bool {
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
        for l:=len(P[0])-1; l>-1; l-- {
            node, b := tree.InTree(P[0][l])
            if P[0][l] == P[0][0] {
                tree.Nodes = append(tree.Nodes, node)

            } else {
                connection := &Connection{
                    FromNodeID: P[0][l],
                    ToNodeID: P[0][l-1],
                    Cost: cost,
                }
                node.Connections = append(node.Connections, connection)
                if !(b) { tree.Nodes = append(tree.Nodes, node) }        
            }     
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

func (graph *Graph) findGraph(id int) *Vertex {
    for _, vertex := range graph.Vertexs {
        if vertex.ID == id {
            return vertex
        }
    }
    return &Vertex{}
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
