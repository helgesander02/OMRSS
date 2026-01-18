package can

import (
	"fmt"
)

func Generate_CAN_Flows(CANnode []int, importantCAN int, unimportantCAN int, HyperPeriod int) ([]*Flow, []*Flow) {
	// Generate CAN Flows
	ImportantCANFlows := Generate_Important_CANFlow(CANnode, importantCAN, HyperPeriod)
	UnimportantCANFlows := Generate_Unimportant_CANFlow(CANnode, unimportantCAN, HyperPeriod)
	fmt.Println("Important CAN:", len(ImportantCANFlows), "Unimportant CAN:", len(UnimportantCANFlows), "Complete generating CAN flows.")

	return ImportantCANFlows, UnimportantCANFlows
}

func Generate_Important_CANFlow(CANnode []int, impcan int, HyperPeriod int) []*Flow {
	ImportantCANFlows := []*Flow{}
	for flow := 0; flow < impcan; flow++ {
		importantCAN := config_ImportantCAN_Stream()
		source, destination := random_CAN_Devices_For_Path(CANnode)

		Flow := Generate_CAN_Streams(importantCAN.Period, importantCAN.Deadline, importantCAN.DataSize, HyperPeriod)
		Flow.Source = source
		Flow.Destination = destination

		ImportantCANFlows = append(ImportantCANFlows, Flow)
	}

	return ImportantCANFlows
}

func Generate_Unimportant_CANFlow(CANnode []int, umimpcan int, HyperPeriod int) []*Flow {
	UnimportantCANFlows := []*Flow{}
	for flow := 0; flow < umimpcan; flow++ {
		unimportantCAN := config_UnimportantCAN_Stream()

		// Random End Devices 1. source(Talker) 2. destinations(listener)
		source, destination := random_CAN_Devices_For_Path(CANnode)

		Flow := Generate_CAN_Streams(unimportantCAN.Period, unimportantCAN.Deadline, unimportantCAN.DataSize, HyperPeriod)
		Flow.Source = source
		Flow.Destination = destination

		UnimportantCANFlows = append(UnimportantCANFlows, Flow)
	}

	return UnimportantCANFlows
}

func Generate_CAN_Streams(period int, deadline int, datasize float64, HyperPeriod int) *Flow {
	var number int = 0

	flow := new_CANFlow(period, deadline, datasize, HyperPeriod)
	for ArrivalTime := 0; ArrivalTime < HyperPeriod; ArrivalTime += period {
		FinishTime := ArrivalTime + deadline
		name := fmt.Sprint("canstream", number)
		stream := new_CANStream(name, ArrivalTime, datasize, deadline, FinishTime)
		flow.Streams = append(flow.Streams, stream)
		number += 1
	}

	return flow
}
