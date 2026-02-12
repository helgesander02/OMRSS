package flow

import (
	"fmt"
)

func (flows *FlowSet) Show_Stream() {
	TSNFlows := flows.TSNFlows
	AVBFlows := flows.AVBFlows
	number := 1
	for _, flow := range TSNFlows {
		name := fmt.Sprint("TSNflow", number)
		fmt.Println(name)
		for _, stream := range flow.Streams {
			fmt.Printf("%s ArrivalTime:%d DataSize:%f Deadline:%d FinishTime:%d\n",
				stream.Name, stream.ArrivalTime, stream.DataSize, stream.Deadline, stream.FinishTime)
		}
		number += 1

		break
	}

	number = 1
	for _, flow := range AVBFlows {
		name := fmt.Sprint("AVBflow", number)
		fmt.Println(name)
		for _, stream := range flow.Streams {
			fmt.Printf("%s ArrivalTime:%d DataSize:%f Deadline:%d FinishTime:%d\n",
				stream.Name, stream.ArrivalTime, stream.DataSize, stream.Deadline, stream.FinishTime)
		}
		number += 1

		break
	}
}

func (flows *FlowSet) Show_TTFlow() {
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

func (flows *FlowSet) Show_CANFlow() {
	Method := flows.EncapsulateMethod
	for _, merhod := range Method {
		number := 1
		fmt.Printf("Method Name: %s\n", merhod.Method_Name)
		for _, flow := range merhod.CAN2TTFlows {
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

func (flows *FlowSet) Show_Flows() {
	// Display all flows.
	fmt.Printf("Total Flows:%d ( TSN Flows:%d  AVB Flows:%d )\n",
		len(flows.TSNFlows)+len(flows.AVBFlows), len(flows.TSNFlows), len(flows.AVBFlows))
}
