package topology

import (
	"os"
	"src/pkg/logger"

	"gopkg.in/yaml.v2"
)

func GenerateTopology(topologyName string, cost float64) *Topology {
	// Read YAML file: ./yaml
	// Create and parse YAML into the Data
	data, err := os.ReadFile("yaml/" + topologyName + ".yaml")
	if err != nil {
		logger.Fatalf("error: %v", err)
	}
	d := Data{}
	err = yaml.Unmarshal([]byte(data), &d)
	if err != nil {
		logger.Fatalf("error: %v", err)
	}

	// Devices(EndStations) Switchs(Bridges)
	topology := NewTopology()

	for s := 0; s < d.Scale.Bridges; s++ {
		switchNode := &Node{ID: s}
		topology.Switch = append(topology.Switch, switchNode)
	}
	defineSwitchConnection(topology, d.BridgeEdges, cost)

	for es := 0; es < d.Scale.EndStations; es++ {
		node := &Node{ID: es + 3000}
		topology.Nodes = append(topology.Nodes, node)
	}
	defineNodesConnection(topology, d.EndStationEdges, cost)

	return topology
}

func defineSwitchConnection(topology *Topology, bridgeEdges []Edge, cost float64) {
	for _, edge := range bridgeEdges {
		topology.AddS2S(edge.Ends[0], edge.Ends[1], cost)
	}

}

func defineNodesConnection(topology *Topology, endstationEdges []Edge, cost float64) {
	for _, edge := range endstationEdges {
		topology.AddnullN2S(edge.Ends[0], edge.Ends[1], cost)
	}

}
