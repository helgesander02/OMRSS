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
	fmt.Printf("Complete generating round%d bgframes.\n", 1)

	// Generate TT Input Flows,round 2
	tsnFlows = GenerateTTTSNFlow(tsnFlows, nnodeLength, inputTSN, hyperPeriod)
	avbFlows = GenerateTTAVBFlow(avbFlows, nnodeLength, inputAVB, hyperPeriod)
	fmt.Printf("Complete generating round%d tsnframes.\n", 2)

	fmt.Println("TSN:", len(tsnFlows), "AVB:", len(avbFlows), "Complete generating TT Flows.")

	return tsnFlows, avbFlows
}

func GenerateTTTSNFlow(flows []*Flow, nnodeLength int, TS int, hyperPeriod int) []*Flow {
	for flow := 0; flow < TS; flow++ {
		tsn := configTSNFrame()
		source, destinations := randomTTDevicesForTree(nnodeLength)

		ttFlow := GenerateTTFrame(tsn.Period, tsn.Deadline, tsn.DataSize, hyperPeriod)
		ttFlow.Source = source
		ttFlow.Destinations = destinations

		flows = append(flows, ttFlow)
	}
	return flows
}

func GenerateTTAVBFlow(flows []*Flow, nnodeLength int, AS int, hyperPeriod int) []*Flow {
	for flow := 0; flow < AS; flow++ {
		avb := configAVBFrame()
		source, destinations := randomTTDevicesForTree(nnodeLength)

		ttFlow := GenerateTTFrame(avb.Period, avb.Deadline, avb.DataSize, hyperPeriod)
		ttFlow.Source = source
		ttFlow.Destinations = destinations

		flows = append(flows, ttFlow)
	}
	return flows
}

func GenerateTTFrame(period int, deadline int, datasize float64, hyperPeriod int) *Flow {
	var number int = 0

	flow := newTTFlow(period, deadline, datasize, hyperPeriod)
	for arrivalTime := 0; arrivalTime < hyperPeriod; arrivalTime += period {
		finishTime := arrivalTime + deadline
		name := fmt.Sprint("frame", number)
		frame := newTTFrame(name, arrivalTime, datasize, deadline, finishTime)
		flow.Frames = append(flow.Frames, frame)
		number += 1
	}

	return flow
}
