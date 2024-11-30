package network

import (
	"fmt"
	"src/network/flow"
	"src/network/graph"
	"src/network/topology"
)

type Network struct {
	HyperPeriod  int
	BytesRate    float64
	Bandwidth    float64
	TopologyName string
	BG_TSN       int
	BG_AVB       int
	Input_TSN    int
	Input_AVB    int
	Topology     *topology.Topology
	Flow_Set     *flow.Flows
	Graph_Set    *graph.Graphs
}

func new_OMACO_Network(topology_name string, bg_tsn int, bg_avb int, input_tsn int, input_avb int, hyperperiod int, bandwidth float64) *Network {
	// 1. Define network parameters
	bw := (bandwidth / 8) * 1e-6 // bytes/us ==> 125 bytes
	bytes_rate := 1. / bw        // The number of bytes that can be transmitted in 1us ==> 1/125
	bw *= float64(hyperperiod)   // The bytes that can be transmitted in 6000us (bytes/us * hyperperiod) ==> 750000 bytes

	Network := &Network{
		HyperPeriod:  hyperperiod,
		BytesRate:    bytes_rate,
		Bandwidth:    bw,
		TopologyName: topology_name,
		BG_TSN:       bg_tsn,
		BG_AVB:       bg_avb,
		Input_TSN:    input_tsn,
		Input_AVB:    input_avb,
	}

	fmt.Println("Define network parameters")
	fmt.Println("----------------------------------------")
	fmt.Printf("HyperPeriod: %d us \n", Network.HyperPeriod)
	fmt.Printf("Bandwidth:  %f bytes/6000us \n", Network.Bandwidth)
	fmt.Printf("BytesRate:  %f BPS \n", Network.BytesRate)
	fmt.Printf("Topology file name: %s \n", Network.TopologyName)
	fmt.Printf("TSN flow: %d, AVB flow: %d \n", Network.Input_TSN+Network.BG_TSN, Network.Input_AVB+Network.BG_AVB)
	fmt.Println()

	return Network
}

type OSRO_Network struct {
	HyperPeriod     int
	BytesRate       float64
	Bandwidth       float64
	TopologyName    string
	BG_TSN          int
	BG_AVB          int
	Input_TSN       int
	Input_AVB       int
	Important_CAN   int
	Unimportant_CAN int
	Topology        *topology.Topology
	Flow_Set        *flow.Flows
	Graph_Set       *graph.Graphs
}

func new_OSRO_Network(topology_name string, bg_tsn int, bg_avb int, input_tsn int, input_avb int, important_can int, unimportant_can int, hyperperiod int, bandwidth float64) *OSRO_Network {
	// 1. Define network parameters
	bw := (bandwidth / 8) * 1e-6 // bytes/us ==> 125 bytes
	bytes_rate := 1. / bw        // The number of bytes that can be transmitted in 1us ==> 1/125
	bw *= float64(hyperperiod)   // The bytes that can be transmitted in 6000us (bytes/us * hyperperiod) ==> 750000 bytes

	Network := &OSRO_Network{
		HyperPeriod:     hyperperiod,
		BytesRate:       bytes_rate,
		Bandwidth:       bw,
		TopologyName:    topology_name,
		BG_TSN:          bg_tsn,
		BG_AVB:          bg_avb,
		Input_TSN:       input_tsn,
		Input_AVB:       input_avb,
		Important_CAN:   important_can,
		Unimportant_CAN: unimportant_can,
	}

	fmt.Println("Define network parameters")
	fmt.Println("----------------------------------------")
	fmt.Printf("HyperPeriod: %d us \n", Network.HyperPeriod)
	fmt.Printf("Bandwidth:  %f bytes/6000us \n", Network.Bandwidth)
	fmt.Printf("BytesRate:  %f BPS \n", Network.BytesRate)
	fmt.Printf("Topology file name: %s \n", Network.TopologyName)
	fmt.Printf("TSN flow: %d, AVB flow: %d \n", Network.Input_TSN+Network.BG_TSN, Network.Input_AVB+Network.BG_AVB)
	fmt.Printf("CAN important flow: %d, CAN unimportant flow: %d \n", Network.Important_CAN, Network.Unimportant_CAN)
	fmt.Println()

	return Network
}
