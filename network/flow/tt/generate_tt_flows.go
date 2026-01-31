package tt

import (
	"fmt"
)

func Generate_TT_Flows(Nnode_length int, bgTSN int, bgAVB int, inputTSN int, inputAVB int, HyperPeriod int) ([]*Flow, []*Flow) {
	tsn_flows := new_TSN_flows()
	avb_flows := new_AVB_flows()

	// Generate TT BG Flows, round 1
	tsn_flows = Generate_TT_TSNFlow(tsn_flows, Nnode_length, bgTSN, HyperPeriod)
	avb_flows = Generate_TT_AVBFlow(avb_flows, Nnode_length, bgAVB, HyperPeriod)
	fmt.Printf("Complete generating round%d bgstreams.\n", 1)

	// Generate TT Input Flows,round 2
	tsn_flows = Generate_TT_TSNFlow(tsn_flows, Nnode_length, inputTSN, HyperPeriod)
	avb_flows = Generate_TT_AVBFlow(avb_flows, Nnode_length, inputAVB, HyperPeriod)
	fmt.Printf("Complete generating round%d tsnstreams.\n", 2)

	fmt.Println("TSN:", len(tsn_flows), "AVB:", len(avb_flows), "Complete generating TT Flows.")

	return tsn_flows, avb_flows
}

func Generate_TT_TSNFlow(flows []*Flow, Nnode_length int, TS int, HyperPeriod int) []*Flow {
	for flow := 0; flow < TS; flow++ {
		tsn := config_TSN_Stream()
		source, destinations := random_TT_Devices_For_Tree(Nnode_length)

		Flow := Generate_TT_Stream(tsn.Period, tsn.Deadline, tsn.DataSize, HyperPeriod)
		Flow.Source = source
		Flow.Destinations = destinations

		flows = append(flows, Flow)
	}
	return flows
}

func Generate_TT_AVBFlow(flows []*Flow, Nnode_length int, AS int, HyperPeriod int) []*Flow {
	for flow := 0; flow < AS; flow++ {
		avb := config_AVB_Stream()
		source, destinations := random_TT_Devices_For_Tree(Nnode_length)

		Flow := Generate_TT_Stream(avb.Period, avb.Deadline, avb.DataSize, HyperPeriod)
		Flow.Source = source
		Flow.Destinations = destinations

		flows = append(flows, Flow)
	}
	return flows
}

func Generate_TT_Stream(period int, deadline int, datasize float64, HyperPeriod int) *Flow {
	var number int = 0

	flow := new_TTFlow(period, deadline, datasize, HyperPeriod)
	for ArrivalTime := 0; ArrivalTime < HyperPeriod; ArrivalTime += period {
		FinishTime := ArrivalTime + deadline
		name := fmt.Sprint("stream", number)
		stream := new_TTStream(name, ArrivalTime, datasize, deadline, FinishTime)
		flow.Streams = append(flow.Streams, stream)
		number += 1
	}

	return flow
}
