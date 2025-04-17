package flow

import (
	"fmt"
)

var (
	bg_tsnflows_end int
	bg_avbflow_end  int
)

func Generate_Flows(Nnode_length int, bg_tsn int, bg_avb int, input_tsn int, input_avb int, HyperPeriod int) *Flows {
	// Constructing Flows structures
	flow_set := new_Flows()
	bg_tsnflows_end = bg_tsn
	bg_avbflow_end = bg_avb

	// round 1
	Generate_TSNFlow(flow_set, Nnode_length, bg_tsn, HyperPeriod)
	Generate_AVBFlow(flow_set, Nnode_length, bg_avb, HyperPeriod)
	fmt.Printf("Complete generating round%d streams.\n", 1)

	// round 2
	Generate_TSNFlow(flow_set, Nnode_length, input_tsn, HyperPeriod)
	Generate_AVBFlow(flow_set, Nnode_length, input_avb, HyperPeriod)
	fmt.Printf("Complete generating round%d streams.\n", 2)

	return flow_set
}

func Generate_TSNFlow(flows *Flows, Nnode_length int, TS int, HyperPeriod int) {
	for flow := 0; flow < TS; flow++ {
		tsn := config_TSN_Stream()

		// Random End Devices 1. source(Talker) 2. destinations(listener)
		source, destinations := random_TT_Devices_For_Tree(Nnode_length)

		Flow := Generate_stream(tsn.Period, tsn.Deadline, tsn.DataSize, HyperPeriod)
		Flow.Source = source
		Flow.Destinations = destinations

		flows.TSNFlows = append(flows.TSNFlows, Flow)
	}
}

func Generate_AVBFlow(flows *Flows, Nnode_length int, AS int, HyperPeriod int) {
	for flow := 0; flow < AS; flow++ {
		avb := config_AVB_Stream()

		// Random End Devices 1. source(Talker) 2. destinations(listener)
		source, destinations := random_TT_Devices_For_Tree(Nnode_length)

		Flow := Generate_stream(avb.Period, avb.Deadline, avb.DataSize, HyperPeriod)
		Flow.Source = source
		Flow.Destinations = destinations

		flows.AVBFlows = append(flows.AVBFlows, Flow)
	}
}

func Generate_stream(period int, deadline int, datasize float64, HyperPeriod int) *Flow {
	var number int = 0

	flow := new_Flow(period, deadline, datasize, HyperPeriod)
	for ArrivalTime := 0; ArrivalTime < HyperPeriod; ArrivalTime += period {
		FinishTime := ArrivalTime + deadline
		name := fmt.Sprint("stream", number)
		stream := new_Stream(name, ArrivalTime, datasize, deadline, FinishTime)
		flow.Streams = append(flow.Streams, stream)
		number += 1
	}

	return flow
}
