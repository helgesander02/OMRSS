package flow

import (
	"fmt"
	"src/pkg/logger"
)

func (flows *FlowSet) ShowFrame() {
	TSNFlows := flows.TSNFlows
	AVBFlows := flows.AVBFlows
	number := 1
	for _, flow := range TSNFlows {
		name := fmt.Sprint("TSNflow", number)
		logger.Println(name)
		for _, frame := range flow.Frames {
			logger.Printf("%s ArrivalTime:%d DataSize:%f Deadline:%d FinishTime:%d\n",
				frame.Name, frame.ArrivalTime, frame.DataSize, frame.Deadline, frame.FinishTime)
		}
		number += 1

		break
	}

	number = 1
	for _, flow := range AVBFlows {
		name := fmt.Sprint("AVBflow", number)
		logger.Println(name)
		for _, frame := range flow.Frames {
			logger.Printf("%s ArrivalTime:%d DataSize:%f Deadline:%d FinishTime:%d\n",
				frame.Name, frame.ArrivalTime, frame.DataSize, frame.Deadline, frame.FinishTime)
		}
		number += 1

		break
	}
}

func (flows *FlowSet) ShowTTFlow() {
	TSNFlows := flows.TSNFlows
	AVBFlows := flows.AVBFlows
	number := 1
	for _, flow := range TSNFlows {
		name := fmt.Sprint("TTflow", number)
		logger.Printf("Source: %d\n", flow.Source)
		logger.Printf("Destinations: %v\n", flow.Destinations)
		logger.Printf("%s : period:%d us, deadline:%d us, datasize:%f bytes\n",
			name, flow.Period, flow.Deadline, flow.DataSize)
		number += 1

		break
	}

	number = 1
	for _, flow := range AVBFlows {
		name := fmt.Sprint("AVBflow", number)
		logger.Printf("Source: %d\n", flow.Source)
		logger.Printf("Destinations: %v\n", flow.Destinations)
		logger.Printf("%s : period:%d us, deadline:%d us, datasize:%f bytes\n",
			name, flow.Period, flow.Deadline, flow.DataSize)
		number += 1

		break
	}
}

func (flows *FlowSet) ShowCANFlow() {
	Method := flows.EncapsulateMethod
	for _, method := range Method {
		number := 1
		logger.Printf("Method Name: %s\n", method.MethodName)
		for _, flow := range method.CAN2TTFlows {
			name := fmt.Sprint("CAN2TTflow", number)
			logger.Printf("Source: %d\n", flow.Source)
			logger.Printf("Destination: %v\n", flow.Destination)
			logger.Printf("%s : period:%d us, deadline:%d us, datasize:%f bytes\n",
				name, flow.Period, flow.Deadline, flow.DataSize)
			number += 1

			break
		}
	}
}

func (flows *FlowSet) ShowFlows() {
	// Display all flows.
	logger.Printf("Total Flows:%d ( TSN Flows:%d  AVB Flows:%d )\n",
		len(flows.TSNFlows)+len(flows.AVBFlows), len(flows.TSNFlows), len(flows.AVBFlows))
}
