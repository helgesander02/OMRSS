package topology

import (
    "fmt"
)

func (topology *Topology) Show_Topology() {
    fmt.Println("Talker")
    fmt.Println("-------------------") 
    for _, node := range topology.Talker {
        for _,conn := range node.Connections {
            fmt.Printf("%d --> %d cost: %f bytes/s\n", conn.FromNodeID, conn.ToNodeID, conn.Cost)
        }   
    }
    fmt.Println("Switch ")
    fmt.Println("-------------------") 
    for _, node := range topology.Switch {
        for _,conn := range node.Connections {
            fmt.Printf("%d --> %d cost: %f bytes/s\n", conn.FromNodeID, conn.ToNodeID, conn.Cost)
        }   
    } 
    fmt.Println("Listener")
    fmt.Println("-------------------") 
    for _, node := range topology.Listener {
        for _,conn := range node.Connections {
            fmt.Printf("%d --> %d cost: %f bytes/s\n", conn.FromNodeID, conn.ToNodeID, conn.Cost)
        }   
    } 
}