package tt

import (
	"fmt"

	"src/pkg/config"
	"src/pkg/logger"
)

func GenerateTTFlows(tsn config.TSNFlowConfig, avb config.AVBFlowConfig, hyperPeriod int, nnodeLength int, mode string) ([]*Flow, []*Flow) {
	tsnFlows := newTSNFlows()
	avbFlows := newAVBFlows()

	// Generate TT BG Flows, round 1
	tsnFlows = GenerateTTTSNFlow(tsnFlows, nnodeLength, tsn.Background, hyperPeriod, mode)
	avbFlows = GenerateTTAVBFlow(avbFlows, nnodeLength, avb.Background, hyperPeriod, mode)
	logger.Printf("Complete generating round%d bgframes.\n", 1)

	// Generate TT Input Flows, round 2
	tsnFlows = GenerateTTTSNFlow(tsnFlows, nnodeLength, tsn.Input, hyperPeriod, mode)
	avbFlows = GenerateTTAVBFlow(avbFlows, nnodeLength, avb.Input, hyperPeriod, mode)
	logger.Printf("Complete generating round%d tsnframes.\n", 2)

	logger.Println("TSN:", len(tsnFlows), "AVB:", len(avbFlows), "Complete generating TT Flows.")

	return tsnFlows, avbFlows
}

func GenerateTTTSNFlow(flows []*Flow, nnodeLength int, TS int, hyperPeriod int, mode string) []*Flow {
	for range TS {
		tsn := configTSNFrame()
		source, destinations := randomTTDevices(nnodeLength, mode)

		ttFlow := GenerateTTFrame(tsn.Period, tsn.Deadline, tsn.DataSize, hyperPeriod)
		ttFlow.Source = source
		ttFlow.Destinations = destinations

		flows = append(flows, ttFlow)
	}
	return flows
}

func GenerateTTAVBFlow(flows []*Flow, nnodeLength int, AS int, hyperPeriod int, mode string) []*Flow {
	for range AS {
		avb := configAVBFrame()
		source, destinations := randomTTDevices(nnodeLength, mode)

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
		name := fmt.Sprint("ttframe", number)
		frame := newTTFrame(name, arrivalTime, datasize, deadline, finishTime)
		flow.Frames = append(flow.Frames, frame)
		number += 1
	}

	return flow
}
