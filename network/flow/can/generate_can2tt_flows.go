package can

import (
	"fmt"
	"time"
)

func Generate_CAN2TT_Flows(CANnode []int, importantCAN int, unimportantCAN int, hyperperiod int) []*Method {
	// step 1: generate CAN flows
	ImportantCANFlows, UnimportantCANFlows := Generate_CAN_Flows(CANnode, importantCAN, unimportantCAN, hyperperiod)

	// step2: prepare method list
	var method_list = []string{"fifo", "priority", "obo", "wst", "mao"}

	// step3: according to different encapsulation methods, generate CAN2TT flows
	Method_Set := new_Method_Set()
	for _, method_name := range method_list {
		can2ttClusterPool := new_ClusterPool()

		if method_name == "mao" {
			for _, impf := range ImportantCANFlows {
				flowCopy := impf.deepcopyFlow()
				can2ttClusterPool.organizeCANStreamByPeriodAndDomain(flowCopy)
			}
			for _, unimpf := range UnimportantCANFlows {
				flowCopy := unimpf.deepcopyFlow()
				can2ttClusterPool.organizeCANStreamByPeriodAndDomain(flowCopy)
			}

		} else {
			for _, impf := range ImportantCANFlows {
				flowCopy := impf.deepcopyFlow()
				can2ttClusterPool.organizeCANStreamByDomain(flowCopy)
			}
			for _, unimpf := range UnimportantCANFlows {
				flowCopy := unimpf.deepcopyFlow()
				can2ttClusterPool.organizeCANStreamByDomain(flowCopy)
			}
		}

		start := time.Now()
		method := new_Method(method_name)
		method.EncapsulateCAN2TT(can2ttClusterPool)
		method.CAN2TSN_Delay = time.Since(start)
		Method_Set = append(Method_Set, method)
	}

	return Method_Set
}

type ClusterPool struct {
	Clusters []*CANStreamCluster
}

func new_ClusterPool() *ClusterPool {
	return &ClusterPool{}
}

func (can2ttClusterPool *ClusterPool) organizeCANStreamByDomain(f *Flow) {
	for _, cluster := range can2ttClusterPool.Clusters {
		if cluster.Source == f.Source && cluster.Destination == f.Destination {
			cluster.Streams = append(cluster.Streams, f.Streams...)
			return
		}
	}
	can2ttClusterPool.addNewCluster(f)
}

func (can2ttClusterPool *ClusterPool) organizeCANStreamByPeriodAndDomain(f *Flow) {
	for _, cluster := range can2ttClusterPool.Clusters {
		if cluster.Period == f.Period && cluster.Source == f.Source && cluster.Destination == f.Destination {
			cluster.Streams = append(cluster.Streams, f.Streams...)
			return
		}
	}
	can2ttClusterPool.addNewCluster(f)
}

func (can2ttClusterPool *ClusterPool) addNewCluster(f *Flow) {
	cluster := new_StreamCluster()
	cluster.Source = f.Source
	cluster.Destination = f.Destination
	cluster.Period = f.Period
	cluster.Deadline = f.Deadline
	cluster.DataSize = f.DataSize
	cluster.HyperPeriod = f.HyperPeriod
	cluster.Streams = append(cluster.Streams, f.Streams...)

	can2ttClusterPool.Clusters = append(can2ttClusterPool.Clusters, cluster)
}

func (can2ttClusterPool *ClusterPool) Show_ClusterPool() {
	fmt.Println("CAN2TT Traffic Router:")
	for _, cluster := range can2ttClusterPool.Clusters {
		cluster.Show_StreamCluster()
	}
}

type CANStreamCluster struct {
	Flow
}

func new_StreamCluster() *CANStreamCluster {
	return &CANStreamCluster{}
}

func (cluster *CANStreamCluster) getStreamsByCurrentTime(current_time int) []*Stream {
	streams := []*Stream{}
	for _, stream := range cluster.Streams {
		if stream.ArrivalTime == current_time {
			streams = append(streams, stream)
		}
	}
	return streams
}

func (cluster *CANStreamCluster) Show_StreamCluster() {
	fmt.Printf("Queue (%d→%d) streams=%d\n", cluster.Source, cluster.Destination, len(cluster.Streams))
	fmt.Printf("Period: %v  ,Deadline: %v ,Datasize: %v\n", cluster.Period, cluster.Deadline, cluster.DataSize)
}
