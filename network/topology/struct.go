package topology

import "src/pkg/random"

var rng *random.Generator

func FillRNG(r *random.Generator) {
	rng = r
}

type Node struct {
	ID    int
	Links []*Link
}

type Link struct {
	FromNodeID int     // start
	ToNodeID   int     // next
	Cost       float64 // 1Gbps => (750,000 bytes/6ms) 750,000 bytes under 6ms for each link ==> 125 bytes/us
}

func NewLink(fromNodeID int, toNodeID int, cost float64) *Link {
	return &Link{
		FromNodeID: fromNodeID,
		ToNodeID:   toNodeID,
		Cost:       cost,
	}
}

type Topology struct {
	Talker   []*Node
	Switch   []*Node
	Listener []*Node
	Nodes    []*Node
}

func NewTopology() *Topology {
	return &Topology{}
}

type Edge struct {
	Ends []int `yaml:"ends"`
}

type Data struct {
	Scale struct {
		EndStations int `yaml:"end_stations"`
		Bridges     int `yaml:"bridges"`
	} `yaml:"scale"`
	EndStationEdges []Edge `yaml:"end_station_edges"`
	BridgeEdges     []Edge `yaml:"bridge_edges"`
}
