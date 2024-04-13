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

	// 1. Constructing Topology structures
	topology := newTopology()

	// 2. Create Switch Node
	// 2.1 Switchs (Bridges) number
	// 2.2 Switch Connection
	for s := 0; s < d.Scale.Bridges; s++ {
		snode := &Node{ID: s}
		topology.Switch = append(topology.Switch, snode)
	}
	define_switch_connection(topology, d.BridgeEdges, cost)

	// 3. Create null Nodes
	// 3.1 END Devices (End Stations) number
	// 3.2 END Devices Connection
	for es := 0; es < d.Scale.EndStations; es++ {
		node := &Node{ID: es + 3000}
		topology.Nodes = append(topology.Nodes, node)
	}
	define_nodes_connection(topology, d.EndStationEdges, cost)

	return topology
}

func newTopology() *Topology {
	return &Topology{}
}

// Define a topology where switches are interconnected with each other
func define_switch_connection(topology *Topology, bridge_edges []Edge, cost float64) {
	for _, edge := range bridge_edges {
		topology.AddS2S(edge.Ends[0], edge.Ends[1], cost)
	}

}

// Define a topology where null nodes are connected to switches
func define_nodes_connection(topology *Topology, endstation_edges []Edge, cost float64) {
	for _, edge := range endstation_edges {
		topology.AddnullN2S(edge.Ends[0], edge.Ends[1], cost)
	}

}
