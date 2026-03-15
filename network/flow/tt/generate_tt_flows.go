package tt

import (
	"fmt"
)

func GenerateTTFlows(nnodeLength int, bgTSN int, bgAVB int, inputTSN int, inputAVB int, hyperPeriod int) ([]*Flow, []*Flow) {
	tsnFlows := newTSNFlows()
	avbFlows := newAVBFlows()

	// Generate TT BG Flows, round 1
	tsnFlows = GenerateTTTSNFlow(tsnFlows, nnodeLength, bgTSN, hyperPeriod)
	avbFlows = GenerateTTAVBFlow(avbFlows, nnodeLength, bgAVB, hyperPeriod)
	fmt.Printf("Complete generating round%d bgstreams.\n", 1)

	// Generate TT Input Flows,round 2
	tsnFlows = GenerateTTTSNFlow(tsnFlows, nnodeLength, inputTSN, hyperPeriod)
	avbFlows = GenerateTTAVBFlow(avbFlows, nnodeLength, inputAVB, hyperPeriod)
	fmt.Printf("Complete generating round%d tsnstreams.\n", 2)

	fmt.Println("TSN:", len(tsnFlows), "AVB:", len(avbFlows), "Complete generating TT Flows.")

	return tsnFlows, avbFlows
}

func GenerateTTTSNFlow(flows []*Flow, nnodeLength int, TS int, hyperPeriod int) []*Flow {
	for flow := 0; flow < TS; flow++ {
		tsn := configTSNStream()
		source, destinations := randomTTDevicesForTree(nnodeLength)

		ttFlow := GenerateTTStream(tsn.Period, tsn.Deadline, tsn.DataSize, hyperPeriod)
		ttFlow.Source = source
		ttFlow.Destinations = destinations

		flows = append(flows, ttFlow)
	}
	return flows
}

func GenerateTTAVBFlow(flows []*Flow, nnodeLength int, AS int, hyperPeriod int) []*Flow {
	for flow := 0; flow < AS; flow++ {
		avb := configAVBStream()
		source, destinations := randomTTDevicesForTree(nnodeLength)

		ttFlow := GenerateTTStream(avb.Period, avb.Deadline, avb.DataSize, hyperPeriod)
		ttFlow.Source = source
		ttFlow.Destinations = destinations

		flows = append(flows, ttFlow)
	}
	return flows
}

func GenerateTTStream(period int, deadline int, datasize float64, hyperPeriod int) *Flow {
	var number int = 0

	flow := newTTFlow(period, deadline, datasize, hyperPeriod)
	for arrivalTime := 0; arrivalTime < hyperPeriod; arrivalTime += period {
		finishTime := arrivalTime + deadline
		name := fmt.Sprint("stream", number)
		stream := newTTStream(name, arrivalTime, datasize, deadline, finishTime)
		flow.Streams = append(flow.Streams, stream)
		number += 1
	}

	return flow
}
