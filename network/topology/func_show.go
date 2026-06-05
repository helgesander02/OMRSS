package topology

import "src/pkg/logger"

func (topology *Topology) ShowTopology() {
	topology.ShowTalker()
	topology.ShowSwitch()
	topology.ShowListener()
}

func (topology *Topology) ShowTalker() {
	logger.Println("Talker")
	logger.Println("-------------------")
	for _, node := range topology.Talker {
		node.ShowConnection()
	}
}

func (topology *Topology) ShowSwitch() {
	logger.Println("Switch ")
	logger.Println("-------------------")
	for _, node := range topology.Switch {
		node.ShowConnection()
	}
}

func (topology *Topology) ShowListener() {
	logger.Println("Listener")
	logger.Println("-------------------")
	for _, node := range topology.Listener {
		node.ShowConnection()
	}
}

func (node *Node) ShowConnection() {
	for _, conn := range node.Links {
		logger.Printf("%d --> %d cost: %f bytes/us\n", conn.FromNodeID, conn.ToNodeID, conn.Cost)
	}
}
