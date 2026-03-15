package can

import (
	"fmt"
	"time"
)

func GenerateCAN2TTFlows(CANnode []int, importantCAN int, unimportantCAN int, hyperperiod int) []*Method {
	// step 1: generate CAN flows
	importantCANFlows, unimportantCANFlows := GenerateCANFlows(CANnode, importantCAN, unimportantCAN, hyperperiod)

	// step2: prepare method list
	var methodList = []string{"fifo", "priority", "obo", "wst", "mao"}

	// step3: according to different encapsulation methods, generate CAN2TT flows
	methodSet := newMethodSet()
	for _, methodName := range methodList {
		can2ttClusterPool := newClusterPool()

		if methodName == "mao" {
			for _, impf := range importantCANFlows {
				flowCopy := impf.deepCopyFlow()
				can2ttClusterPool.organizeCANStreamByPeriodAndDomain(flowCopy)
			}
			for _, unimpf := range unimportantCANFlows {
				flowCopy := unimpf.deepCopyFlow()
				can2ttClusterPool.organizeCANStreamByPeriodAndDomain(flowCopy)
			}

		} else {
			for _, impf := range importantCANFlows {
				flowCopy := impf.deepCopyFlow()
				can2ttClusterPool.organizeCANStreamByDomain(flowCopy)
			}
			for _, unimpf := range unimportantCANFlows {
				flowCopy := unimpf.deepCopyFlow()
				can2ttClusterPool.organizeCANStreamByDomain(flowCopy)
			}
		}

		start := time.Now()
		method := newMethod(methodName)
		method.EncapsulateCAN2TT(can2ttClusterPool)
		method.CAN2TSNDelay = time.Since(start)
		methodSet = append(methodSet, method)
	}

	return methodSet
}

type ClusterPool struct {
	Clusters []*CANStreamCluster
}

func newClusterPool() *ClusterPool {
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
	cluster := newStreamCluster()
	cluster.Source = f.Source
	cluster.Destination = f.Destination
	cluster.Period = f.Period
	cluster.Deadline = f.Deadline
	cluster.DataSize = f.DataSize
	cluster.HyperPeriod = f.HyperPeriod
	cluster.Streams = append(cluster.Streams, f.Streams...)

	can2ttClusterPool.Clusters = append(can2ttClusterPool.Clusters, cluster)
}

func (can2ttClusterPool *ClusterPool) ShowClusterPool() {
	fmt.Println("CAN2TT Traffic Router:")
	for _, cluster := range can2ttClusterPool.Clusters {
		cluster.ShowStreamCluster()
	}
}

type CANStreamCluster struct {
	Flow
}

func newStreamCluster() *CANStreamCluster {
	return &CANStreamCluster{}
}

func (cluster *CANStreamCluster) getStreamsByCurrentTime(currentTime int) []*Stream {
	streams := []*Stream{}
	for _, stream := range cluster.Streams {
		if stream.ArrivalTime == currentTime {
			streams = append(streams, stream)
		}
	}
	return streams
}

func (cluster *CANStreamCluster) ShowStreamCluster() {
	fmt.Printf("Queue (%d→%d) streams=%d\n", cluster.Source, cluster.Destination, len(cluster.Streams))
	fmt.Printf("Period: %v  ,Deadline: %v ,Datasize: %v\n", cluster.Period, cluster.Deadline, cluster.DataSize)
}
