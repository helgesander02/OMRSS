package flow

import (
	"fmt"
)

func Generate_Flows(Nnode int, tsn int, avb int, HyperPeriod int) *Flows {
	flow_set := &Flows{}
	for round := 0; round < 2; round++ {
		Generate_TSNFlow(flow_set, Nnode, tsn, HyperPeriod)
		Generate_AVBFlow(flow_set, Nnode, avb, HyperPeriod)
		fmt.Printf("Complete generating round%d streams.\n", round+1)

	}

	return flow_set
}

func Generate_TSNFlow(flows *Flows, Nnode int, TS int, HyperPeriod int) {
	TS2 := TS / 2
	for flow := 0; flow < TS2; flow++ {
		tsn := TSN_stream()

		// Random End Devices 1. source(Talker) 2. destinations(listener)
		source, destinations := Random_Devices(Nnode)

		Flow := Generate_stream(tsn.Period, tsn.Deadline, tsn.DataSize, HyperPeriod)
		Flow.Source = source
		Flow.Destinations = destinations

		flows.TSNFlows = append(flows.TSNFlows, Flow)
	}
}

func Generate_AVBFlow(flows *Flows, Nnode int, AS int, HyperPeriod int) {
	AS2 := AS / 2
	for flow := 0; flow < AS2; flow++ {
		avb := AVB_stream()

		// Random End Devices 1. source(Talker) 2. destinations(listener)
		source, destinations := Random_Devices(Nnode)

		Flow := Generate_stream(avb.Period, avb.Deadline, avb.DataSize, HyperPeriod)
		Flow.Source = source
		Flow.Destinations = destinations

		flows.AVBFlows = append(flows.AVBFlows, Flow)
	}
}

func Generate_stream(period int, deadline int, datasize float64, HyperPeriod int) *Flow {
	var (
		ArrivalTime int = 0
		FinishTime  int = 0
		Deadline    int = 0
		number      int = 0
	)

	flow := &Flow{Period: period, Deadline: deadline, DataSize: datasize, HyperPeriod: HyperPeriod}
	for FinishTime < HyperPeriod {
		Deadline += deadline
		FinishTime += period
		name := fmt.Sprint("stream", number)

		stream := &Stream{
			Name:        name,
			ArrivalTime: ArrivalTime,
			DataSize:    datasize,
			Deadline:    Deadline,
			FinishTime:  FinishTime,
		}

		flow.Streams = append(flow.Streams, stream)
		ArrivalTime += period
		number += 1
	}

	return flow
}
