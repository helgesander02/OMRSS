package topology

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

func Generate_Topology(topology_name string, cost float64) *Topology {
	// 1. Read YAML file
	// 2. Create Data
	// 3. Parse YAML into the Data
	data, err := ioutil.ReadFile("yaml/" + topology_name + ".yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	d := Data{}
	err = yaml.Unmarshal([]byte(data), &d)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// 1. Create Topology
	topology := &Topology{}

	// 2. Create Switch Node
	// 2.1 Switchs (Bridges) number
	// 2.2 Switch Connection
	for s := 0; s < d.Scale.Bridges; s++ {
		snode := &Node{ID: s}
		topology.Switch = append(topology.Switch, snode)
	}
	Define_Switch_Connection(topology, d.BridgeEdges, cost)

	// 3. Create null Nodes
	// 3.1 END Devices (End Stations) number
	// 3.2 END Devices Connection
	for i := 0; i < d.Scale.EndStations; i++ {
		node := &Node{ID: i + 3000}
		topology.Nodes = append(topology.Nodes, node)
	}
	Define_Nodes_Connection(topology, d.EndStationEdges, cost)

	return topology
}

// Define a topology where switches are interconnected with each other
func Define_Switch_Connection(topology *Topology, bridge_edges []Edge, cost float64) {
	for _, edge := range bridge_edges {
		topology.AddS2S(edge.Ends[0], edge.Ends[1], cost)
	}

}

// Define a topology where null nodes are connected to switches
func Define_Nodes_Connection(topology *Topology, endstation_edges []Edge, cost float64) {
	for _, edge := range endstation_edges {
		topology.AddnullN2S(edge.Ends[0], edge.Ends[1], cost)
	}

}
