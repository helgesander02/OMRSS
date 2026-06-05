package can

import (
	"fmt"
	"src/pkg/config"
	"src/pkg/logger"
)

func GenerateCANFlows(can config.CANFlowConfig, hyperperiod int, CANnode []int) ([]*Flow, []*Flow) {
	// Generate CAN Flows
	importantCANFlows := GenerateImportantCANFlow(can.Important, hyperperiod, CANnode)
	unimportantCANFlows := GenerateUnimportantCANFlow(can.Unimportant, hyperperiod, CANnode)
	logger.Println("Important CAN:", len(importantCANFlows), "Unimportant CAN:", len(unimportantCANFlows), "Complete generating CAN flows.")

	return importantCANFlows, unimportantCANFlows
}

func GenerateImportantCANFlow(impcan int, hyperPeriod int, CANnode []int) []*Flow {
	importantCANFlows := []*Flow{}
	for range impcan {
		importantCAN := configImportantCANFrame()
		source, destination := randomCANDevices(CANnode)

		canFlow := GenerateCANFrames(importantCAN.Period, importantCAN.Deadline, importantCAN.DataSize, hyperPeriod)
		canFlow.Source = source
		canFlow.Destination = destination

		importantCANFlows = append(importantCANFlows, canFlow)
	}

	return importantCANFlows
}

func GenerateUnimportantCANFlow(umimpcan int, hyperPeriod int, CANnode []int) []*Flow {
	unimportantCANFlows := []*Flow{}
	for range umimpcan {
		unimportantCAN := configUnimportantCANFrame()

		// Random End Devices 1. source(Talker) 2. destinations(listener)
		source, destination := randomCANDevices(CANnode)

		canFlow := GenerateCANFrames(unimportantCAN.Period, unimportantCAN.Deadline, unimportantCAN.DataSize, hyperPeriod)
		canFlow.Source = source
		canFlow.Destination = destination

		unimportantCANFlows = append(unimportantCANFlows, canFlow)
	}

	return unimportantCANFlows
}

func GenerateCANFrames(period int, deadline int, datasize float64, hyperPeriod int) *Flow {
	var number int = 0

	flow := newCANFlow(period, deadline, datasize, hyperPeriod)
	for arrivalTime := 0; arrivalTime < hyperPeriod; arrivalTime += period {
		finishTime := arrivalTime + deadline
		name := fmt.Sprint("canframe", number)
		frame := newCANFrame(name, arrivalTime, datasize, deadline, finishTime)
		flow.Frames = append(flow.Frames, frame)
		number += 1
	}

	return flow
}
