package can

import (
	"fmt"
)

func GenerateCANFlows(CANnode []int, importantCAN int, unimportantCAN int, hyperPeriod int) ([]*Flow, []*Flow) {
	// Generate CAN Flows
	importantCANFlows := GenerateImportantCANFlow(CANnode, importantCAN, hyperPeriod)
	unimportantCANFlows := GenerateUnimportantCANFlow(CANnode, unimportantCAN, hyperPeriod)
	fmt.Println("Important CAN:", len(importantCANFlows), "Unimportant CAN:", len(unimportantCANFlows), "Complete generating CAN flows.")

	return importantCANFlows, unimportantCANFlows
}

func GenerateImportantCANFlow(CANnode []int, impcan int, hyperPeriod int) []*Flow {
	importantCANFlows := []*Flow{}
	for flow := 0; flow < impcan; flow++ {
		importantCAN := configImportantCANStream()
		source, destination := randomCANDevicesForPath(CANnode)

		canFlow := GenerateCANStreams(importantCAN.Period, importantCAN.Deadline, importantCAN.DataSize, hyperPeriod)
		canFlow.Source = source
		canFlow.Destination = destination

		importantCANFlows = append(importantCANFlows, canFlow)
	}

	return importantCANFlows
}

func GenerateUnimportantCANFlow(CANnode []int, umimpcan int, hyperPeriod int) []*Flow {
	unimportantCANFlows := []*Flow{}
	for flow := 0; flow < umimpcan; flow++ {
		unimportantCAN := configUnimportantCANStream()

		// Random End Devices 1. source(Talker) 2. destinations(listener)
		source, destination := randomCANDevicesForPath(CANnode)

		canFlow := GenerateCANStreams(unimportantCAN.Period, unimportantCAN.Deadline, unimportantCAN.DataSize, hyperPeriod)
		canFlow.Source = source
		canFlow.Destination = destination

		unimportantCANFlows = append(unimportantCANFlows, canFlow)
	}

	return unimportantCANFlows
}

func GenerateCANStreams(period int, deadline int, datasize float64, hyperPeriod int) *Flow {
	var number int = 0

	flow := newCANFlow(period, deadline, datasize, hyperPeriod)
	for arrivalTime := 0; arrivalTime < hyperPeriod; arrivalTime += period {
		finishTime := arrivalTime + deadline
		name := fmt.Sprint("canstream", number)
		stream := newCANStream(name, arrivalTime, datasize, deadline, finishTime)
		flow.Streams = append(flow.Streams, stream)
		number += 1
	}

	return flow
}
