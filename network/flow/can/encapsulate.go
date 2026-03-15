package can

type EncapsulateConfig struct {
	DatasizeLeast float64
	DatasizeMax   float64
	Step          int
	CANBandwidth  float64
	BytesPerStep  float64
}

func newEncapsulateConfig() *EncapsulateConfig {
	return &EncapsulateConfig{
		DatasizeLeast: 64.,                    // bytes
		DatasizeMax:   1500.,                  // bytes
		Step:          1000,                   // us
		CANBandwidth:  1000000.0,              // CAN bandwidth 1 Mbps
		BytesPerStep:  (1000000.0 / 8) * 1e-3, // bytes per 1000 us
	}
}

func (method *Method) EncapsulateCAN2TT(can2ttClusterPool *ClusterPool) {
	encapConfig := newEncapsulateConfig()
	if len(method.CAN2TTFlows) == 0 {
		method.organizeCAN2TTFlows(can2ttClusterPool)
	}

	switch method.MethodName {
	case "obo":
		method.encapOBO(can2ttClusterPool, encapConfig)
	case "wst":
		method.encapWST(can2ttClusterPool, encapConfig)
	case "mao":
		method.encapMAO(can2ttClusterPool, encapConfig)
	default:
		method.encapFIFOOrPriority(can2ttClusterPool, encapConfig)
	}
}

func (method *Method) organizeCAN2TTFlows(can2ttClusterPool *ClusterPool) {
	for _, cluster := range can2ttClusterPool.Clusters {
		if !method.flowExists(cluster.Source, cluster.Destination) {
			can2ttFlow := newCAN2TTFlow()
			can2ttFlow.Source = cluster.Source
			can2ttFlow.Destination = cluster.Destination
			can2ttFlow.Period = cluster.Period
			can2ttFlow.Deadline = cluster.Deadline
			can2ttFlow.DataSize = cluster.DataSize
			can2ttFlow.HyperPeriod = cluster.HyperPeriod

			method.CAN2TTFlows = append(method.CAN2TTFlows, can2ttFlow)
		}
	}
}

func (method *Method) getCAN2TTFlowByDomain(source int, destination int) *Flow {
	for _, flow := range method.CAN2TTFlows {
		if flow.Source == source && flow.Destination == destination {
			return flow
		}
	}
	return nil
}

func (method *Method) flowExists(source int, destination int) bool {
	for _, flow := range method.CAN2TTFlows {
		if flow.Source == source && flow.Destination == destination {
			return true
		}
	}
	return false
}

func (method *Method) encapFIFOOrPriority(can2ttClusterPool *ClusterPool, encapConfig *EncapsulateConfig) {
	for _, canStreamCluster := range can2ttClusterPool.Clusters {
		queue := newQueue()
		can2ttFlow := method.getCAN2TTFlowByDomain(canStreamCluster.Source, canStreamCluster.Destination)

		deadline := 0
		datasizeCount := 0.
		for currentTime := 0; currentTime < can2ttFlow.HyperPeriod; currentTime += encapConfig.Step {
			queue.appendQueue(canStreamCluster.getStreamsByCurrentTime(currentTime))
			queue.sortQueue(method.MethodName, currentTime)

			// firstly, drop overdue streams
			method.CAN2TTO1Drop += queue.checkDrop(currentTime)

			// secondly, encapsulate streams
			head := 0
			for head < len(queue.Streams) {
				stream := queue.Streams[head]
				datasizeCount += stream.DataSize
				head++
				if deadline == 0 || stream.Deadline < deadline {
					deadline = stream.Deadline
				}

				if datasizeCount >= encapConfig.DatasizeMax {
					method.flushStream(can2ttFlow, currentTime, datasizeCount, deadline, encapConfig.DatasizeLeast)
					datasizeCount = 0
					queue.popQueueByHead(head)
					head = 0
					deadline = 0
				}
			}

			if datasizeCount > 0 && currentTime%can2ttFlow.Period == 0 {
				method.flushStream(can2ttFlow, currentTime, datasizeCount, deadline, encapConfig.DatasizeLeast)
				datasizeCount = 0
				queue.popQueueByHead(head)
				head = 0
				deadline = 0
			}
		}
		if datasizeCount > 0 {
			method.flushStream(can2ttFlow, can2ttFlow.HyperPeriod, datasizeCount, deadline, encapConfig.DatasizeLeast)
			datasizeCount = 0
			deadline = 0
		}
	}
}

func (method *Method) encapOBO(can2ttClusterPool *ClusterPool, encapConfig *EncapsulateConfig) {
	for _, canStreamCluster := range can2ttClusterPool.Clusters {
		queue := newQueue()
		can2ttFlow := method.getCAN2TTFlowByDomain(canStreamCluster.Source, canStreamCluster.Destination)

		for currentTime := 0; currentTime < can2ttFlow.HyperPeriod; currentTime += encapConfig.Step {
			queue.appendQueue(canStreamCluster.getStreamsByCurrentTime(currentTime))

			// firstly, drop overdue streams
			method.CAN2TTO1Drop += queue.checkDrop(currentTime)

			// secondly, encapsulate streams
			for len(queue.Streams) > 0 {
				method.flushStream(can2ttFlow, currentTime, encapConfig.DatasizeLeast, queue.Streams[0].Deadline, encapConfig.DatasizeLeast)
				queue.popQueueByIdx(0)
			}

		}
	}
}

func (method *Method) encapWST(can2ttClusterPool *ClusterPool, encapConfig *EncapsulateConfig) {
	const guardBase = 1800 // µs
	for _, canStreamCluster := range can2ttClusterPool.Clusters {
		queue := newQueue()
		can2ttFlow := method.getCAN2TTFlowByDomain(canStreamCluster.Source, canStreamCluster.Destination)

		datasizeCount := 0.
		deadline := 0
		for currentTime := 0; currentTime < can2ttFlow.HyperPeriod; currentTime += encapConfig.Step {
			queue.appendQueue(canStreamCluster.getStreamsByCurrentTime(currentTime))
			queue.sortQueue(method.MethodName, currentTime)

			// firstly, drop overdue streams
			method.CAN2TTO1Drop += queue.checkDrop(currentTime)

			// secondly, encapsulate streams
			guard := guardBase + len(queue.Streams)*600
			if len(queue.Streams) == 0 {
				continue
			}
			for queue.hasImminent(currentTime, guard) {
				head := 0
				for head < len(queue.Streams) && datasizeCount+queue.Streams[head].DataSize < encapConfig.DatasizeMax {
					if deadline == 0 || queue.Streams[head].Deadline < deadline {
						deadline = queue.Streams[head].Deadline
					}

					datasizeCount += queue.Streams[head].DataSize
					head++
				}

				method.flushStream(can2ttFlow, currentTime, datasizeCount, deadline, encapConfig.DatasizeLeast)
				queue.popQueueByHead(head)
				datasizeCount = 0
				deadline = 0
			}
		}

		if len(queue.Streams) > 0 {
			method.flushStream(can2ttFlow, can2ttFlow.HyperPeriod, datasizeCount, deadline, encapConfig.DatasizeLeast)
			datasizeCount = 0
			deadline = 0
		}
	}
}

func (method *Method) encapMAO(can2ttClusterPool *ClusterPool, encapConfig *EncapsulateConfig) {
	// step1: MAO Aggregation (已實現)
	MTU := 1500.0
	for _, canStreamCluster := range can2ttClusterPool.Clusters {
		queue := newQueue()
		can2ttFlow := method.getCAN2TTFlowByDomain(canStreamCluster.Source, canStreamCluster.Destination)

		deadline := 0
		datasizeCount := 0.
		for currentTime := 0; currentTime < can2ttFlow.HyperPeriod; currentTime += encapConfig.Step {
			queue.appendQueue(canStreamCluster.getStreamsByCurrentTime(currentTime))
			queue.sortQueue(method.MethodName, currentTime)

			// firstly, drop overdue streams
			method.CAN2TTO1Drop += queue.checkDrop(currentTime)

			// secondly, encapsulate streams
			head := 0
			for head < len(queue.Streams) {
				stream := queue.Streams[head]
				datasizeCount += stream.DataSize
				head++
				if deadline == 0 || stream.Deadline < deadline {
					deadline = stream.Deadline
				}

				if datasizeCount >= (MTU / 2) {
					method.flushStream(can2ttFlow, currentTime, datasizeCount, deadline, encapConfig.DatasizeLeast)
					datasizeCount = 0
					queue.popQueueByHead(head)
					head = 0
					deadline = 0
				}
			}

		}

		if datasizeCount >= 0 {
			method.flushStream(can2ttFlow, can2ttFlow.HyperPeriod, datasizeCount, deadline, encapConfig.DatasizeLeast)
			datasizeCount = 0
			deadline = 0
		}
	}

	// step2: Harmonic Merge
	method.harmonicMerge(MTU)
}

func (method *Method) flushStream(flow *Flow, now int, packedSize float64, dl int, datasizeLeast float64) {
	if packedSize < datasizeLeast {
		packedSize = datasizeLeast
	}
	stream := createCAN2TTStream(now, dl, packedSize+42)
	flow.Streams = append(flow.Streams, stream)

	method.BytesSent += packedSize + 42
	method.TTFrameCount += 1
}
