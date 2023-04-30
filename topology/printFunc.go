package topology

import (
    "fmt"
)

func (topology *Topology) Print_() {
    fmt.Printf("")
}

func (topology *Topology) Print_Topology() {
    for _, node := range topology.Talker {
        fmt.Printf("%s \n", node.Name)
        for _,conn := range node.Connections {
            fmt.Printf("%d --> %d cost: %f bytes/s\n", conn.FromNodeID, conn.ToNodeID, conn.Cost)
        }   
    } 
    for _, node := range topology.Switch {
        fmt.Printf("%s \n", node.Name)
        for _,conn := range node.Connections {
            fmt.Printf("%d --> %d cost: %f bytes/s\n", conn.FromNodeID, conn.ToNodeID, conn.Cost)
        }   
    } 
    for _, node := range topology.Listener {
        fmt.Printf("%s \n", node.Name)
        for _,conn := range node.Connections {
            fmt.Printf("%d --> %d cost: %f bytes/s\n", conn.FromNodeID, conn.ToNodeID, conn.Cost)
        }   
    } 
}