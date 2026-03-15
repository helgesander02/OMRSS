package topology

import (
	"fmt"
)

func (topology *Topology) ShowTopology() {
	topology.ShowTalker()
	topology.ShowSwitch()
	topology.ShowListener()
}

func (topology *Topology) ShowTalker() {
	fmt.Println("Talker")
	fmt.Println("-------------------")
	for _, node := range topology.Talker {
		node.ShowConnection()
	}
}

func (topology *Topology) ShowSwitch() {
	fmt.Println("Switch ")
	fmt.Println("-------------------")
	for _, node := range topology.Switch {
		node.ShowConnection()
	}
}

func (topology *Topology) ShowListener() {
	fmt.Println("Listener")
	fmt.Println("-------------------")
	for _, node := range topology.Listener {
		node.ShowConnection()
	}
}

func (node *Node) ShowConnection() {
	for _, conn := range node.Links {
		fmt.Printf("%d --> %d cost: %f bytes/us\n", conn.FromNodeID, conn.ToNodeID, conn.Cost)
	}
}
