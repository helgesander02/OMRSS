package network

import (
	"fmt"
	"src/network/flow"
	"src/network/graph"
	"src/network/topology"
)

type OMACO_Network struct {
	HyperPeriod  int
	BytesRate    float64
	Bandwidth    float64
	TopologyName string
	BGTSN        int
	BGAVB        int
	Input_TSN    int
	Input_AVB    int
	Topology     *topology.Topology
	FlowSet      *flow.FlowSet
	Graph_Set    *graph.Graphs
}

func new_OMACO_Network(topologyName string, bgTSN int, bgAVB int, inputTSN int, inputAVB int, hyperperiod int, bandwidth float64) *OMACO_Network {
	// 1. Define network parameters
	bw := (bandwidth / 8) * 1e-6 // bytes/us ==> 125 bytes
	bytes_rate := 1. / bw        // The number of bytes that can be transmitted in 1us ==> 1/125
	bw *= float64(hyperperiod)   // The bytes that can be transmitted in 6000us (bytes/us * hyperperiod) ==> 750000 bytes

	Network := &OMACO_Network{
		HyperPeriod:  hyperperiod,
		BytesRate:    bytes_rate,
		Bandwidth:    bw,
		TopologyName: topologyName,
		BGTSN:        bgTSN,
		BGAVB:        bgAVB,
		Input_TSN:    inputTSN,
		Input_AVB:    inputAVB,
	}

	fmt.Println("Define network parameters")
	fmt.Println("----------------------------------------")
	fmt.Printf("HyperPeriod: %d us \n", Network.HyperPeriod)
	fmt.Printf("Bandwidth:  %f bytes/6000us \n", Network.Bandwidth)
	fmt.Printf("BytesRate:  %f BPS \n", Network.BytesRate)
	fmt.Printf("Topology file name: %s \n", Network.TopologyName)
	fmt.Printf("TSN flow: %d, AVB flow: %d \n", Network.Input_TSN+Network.BGTSN, Network.Input_AVB+Network.BGAVB)
	fmt.Println()

	return Network
}

type OSRO_Network struct {
	HyperPeriod     int
	BytesRate       float64
	Bandwidth       float64
	TopologyName    string
	BGTSN           int
	BGAVB           int
	Input_TSN       int
	Input_AVB       int
	Important_CAN   int
	Unimportant_CAN int
	Topology        *topology.Topology
	FlowSet         *flow.FlowSet
	Graph_Set       *graph.Graphs
}

func new_OSRO_Network(topologyName string, bgTSN int, bgAVB int, inputTSN int, inputAVB int, importantCAN int, unimportantCAN int, hyperperiod int, bandwidth float64) *OSRO_Network {
	// 1. Define network parameters
	bw := (bandwidth / 8) * 1e-6 // bytes/us ==> 125 bytes
	bytes_rate := 1. / bw        // The number of bytes that can be transmitted in 1us ==> 1/125
	bw *= float64(hyperperiod)   // The bytes that can be transmitted in 6000us (bytes/us * hyperperiod) ==> 750000 bytes

	Network := &OSRO_Network{
		HyperPeriod:     hyperperiod,
		BytesRate:       bytes_rate,
		Bandwidth:       bw,
		TopologyName:    topologyName,
		BGTSN:           bgTSN,
		BGAVB:           bgAVB,
		Input_TSN:       inputTSN,
		Input_AVB:       inputAVB,
		Important_CAN:   importantCAN,
		Unimportant_CAN: unimportantCAN,
	}

	fmt.Println("Define network parameters")
	fmt.Println("----------------------------------------")
	fmt.Printf("HyperPeriod: %d us \n", Network.HyperPeriod)
	fmt.Printf("Bandwidth:  %f bytes/6000us \n", Network.Bandwidth)
	fmt.Printf("BytesRate:  %f BPS \n", Network.BytesRate)
	fmt.Printf("Topology file name: %s \n", Network.TopologyName)
	fmt.Printf("TSN flow: %d, AVB flow: %d \n", Network.Input_TSN+Network.BGTSN, Network.Input_AVB+Network.BGAVB)
	fmt.Printf("CAN important flow: %d, CAN unimportant flow: %d \n", Network.Important_CAN, Network.Unimportant_CAN)
	fmt.Println()

	return Network
}
