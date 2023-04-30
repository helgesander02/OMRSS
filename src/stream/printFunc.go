package stream

import (
	"fmt"
)

func (flows *Flows) Print_() {
	fmt.Printf("")
}

func (flows *Flows) Print_Stream() {
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
		number+=1
	}

	number = 1
	for _, flow := range AVBFlows {
		name := fmt.Sprint("AVBflow", number)
		fmt.Println(name)
		for _, stream := range flow.Streams {
			fmt.Printf("%s ArrivalTime:%d DataSize:%f Deadline:%d FinishTime:%d\n",
			 	stream.Name, stream.ArrivalTime, stream.DataSize, stream.Deadline, stream.FinishTime)
		} 
		number+=1
	}
}

func (flows *Flows) Print_Flow() {
	TSNFlows := flows.TSNFlows
	AVBFlows := flows.AVBFlows
	number := 1
	for _, flow := range TSNFlows {
		name := fmt.Sprint("TSNflow", number)
		fmt.Printf("Source: %d\n", flow.Source)
        fmt.Printf("Destinations: %v\n", flow.Destinations)
		fmt.Printf("%s : period:%d us, deadline:%d us, datasize:%f bytes\n",
		name ,flow.Streams[0].FinishTime ,flow.Streams[0].Deadline ,flow.Streams[0].DataSize)	
		number+=1
	}

	number = 1
	for _, flow := range AVBFlows {
		name := fmt.Sprint("AVBflow", number)
		fmt.Printf("Source: %d\n", flow.Source)
        fmt.Printf("Destinations: %v\n", flow.Destinations)
		fmt.Printf("%s : period:%d us, deadline:%d us, datasize:%f bytes\n",
		name ,flow.Streams[0].FinishTime ,flow.Streams[0].Deadline ,flow.Streams[0].DataSize)	
		number+=1
	}

}

func (flows *Flows) Print_Flows() {	
	// Display all flows.
	fmt.Printf("Total Flows:%d ( TSN Flows:%d  AVB Flows:%d )\n", 
	len(flows.TSNFlows)+len(flows.AVBFlows), len(flows.TSNFlows), len(flows.AVBFlows))	
}