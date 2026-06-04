package flow

import (
	"fmt"
)

func (flows *FlowSet) ShowFrame() {
	TSNFlows := flows.TSNFlows
	AVBFlows := flows.AVBFlows
	number := 1
	for _, flow := range TSNFlows {
		name := fmt.Sprint("TSNflow", number)
		fmt.Println(name)
		for _, frame := range flow.Frames {
			fmt.Printf("%s ArrivalTime:%d DataSize:%f Deadline:%d FinishTime:%d\n",
				frame.Name, frame.ArrivalTime, frame.DataSize, frame.Deadline, frame.FinishTime)
		}
		number += 1

		break
	}

	number = 1
	for _, flow := range AVBFlows {
		name := fmt.Sprint("AVBflow", number)
		fmt.Println(name)
		for _, frame := range flow.Frames {
			fmt.Printf("%s ArrivalTime:%d DataSize:%f Deadline:%d FinishTime:%d\n",
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
		fmt.Printf("Source: %d\n", flow.Source)
		fmt.Printf("Destinations: %v\n", flow.Destinations)
		fmt.Printf("%s : period:%d us, deadline:%d us, datasize:%f bytes\n",
			name, flow.Period, flow.Deadline, flow.DataSize)
		number += 1

		break
	}

	number = 1
	for _, flow := range AVBFlows {
		name := fmt.Sprint("AVBflow", number)
		fmt.Printf("Source: %d\n", flow.Source)
		fmt.Printf("Destinations: %v\n", flow.Destinations)
		fmt.Printf("%s : period:%d us, deadline:%d us, datasize:%f bytes\n",
			name, flow.Period, flow.Deadline, flow.DataSize)
		number += 1

		break
	}
}

func (flows *FlowSet) ShowCANFlow() {
	Method := flows.EncapsulateMethod
	for _, method := range Method {
		number := 1
		fmt.Printf("Method Name: %s\n", method.MethodName)
		for _, flow := range method.CAN2TTFlows {
			name := fmt.Sprint("CAN2TTflow", number)
			fmt.Printf("Source: %d\n", flow.Source)
			fmt.Printf("Destination: %v\n", flow.Destination)
			fmt.Printf("%s : period:%d us, deadline:%d us, datasize:%f bytes\n",
				name, flow.Period, flow.Deadline, flow.DataSize)
			number += 1

			break
		}
	}
}

func (flows *FlowSet) ShowFlows() {
	// Display all flows.
	fmt.Printf("Total Flows:%d ( TSN Flows:%d  AVB Flows:%d )\n",
		len(flows.TSNFlows)+len(flows.AVBFlows), len(flows.TSNFlows), len(flows.AVBFlows))
}
