package can

type Encapsulate_Config struct {
	Datasize_Least float64
	Datasize_Max   float64
	Step           int
	CANBandwidth   float64
	BytesPerStep   float64
}

func new_Encapsulate_Config() *Encapsulate_Config {
	return &Encapsulate_Config{
		Datasize_Least: 64.,                    // bytes
		Datasize_Max:   1500.,                  // bytes
		Step:           1000,                   // us
		CANBandwidth:   1000000.0,              // CAN bandwidth 1 Mbps
		BytesPerStep:   (1000000.0 / 8) * 1e-3, // bytes per 1000 us
	}
}

func (method *Method) EncapsulateCAN2TT(can2ttClusterPool *ClusterPool) {
	encap_config := new_Encapsulate_Config()
	if len(method.CAN2TTFlows) == 0 {
		method.organizeCAN2TTFlows(can2ttClusterPool)
	}

	switch method.Method_Name {
	case "obo":
		method.encap_obo(can2ttClusterPool, encap_config)
	case "wst":
		method.encap_wst(can2ttClusterPool, encap_config)
	case "mao":
		method.encap_mao(can2ttClusterPool, encap_config)
	default:
		method.encap_fifo_or_priority(can2ttClusterPool, encap_config)
	}
}

func (method *Method) organizeCAN2TTFlows(can2ttClusterPool *ClusterPool) {
	for _, cluster := range can2ttClusterPool.Clusters {
		if !method.flowExists(cluster.Source, cluster.Destination) {
			can2ttFlow := new_CAN2TTFlow()
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

func (method *Method) encap_fifo_or_priority(can2ttClusterPool *ClusterPool, encap_config *Encapsulate_Config) {
	for _, can_stream_cluster := range can2ttClusterPool.Clusters {
		queue := new_Queue()
		can2ttFlow := method.getCAN2TTFlowByDomain(can_stream_cluster.Source, can_stream_cluster.Destination)

		deadline := 0
		datasize_count := 0.
		for current_time := 0; current_time < can2ttFlow.HyperPeriod; current_time += encap_config.Step {
			queue.appendQueue(can_stream_cluster.getStreamsByCurrentTime(current_time))
			queue.sortQueue(method.Method_Name, current_time)

			// fristly, drop overdue streams
			method.CAN2TT_O1_Drop += queue.checkDrop(current_time)

			// secondly, encapsulate streams
			head := 0
			for head < len(queue.Streams) {
				stream := queue.Streams[head]
				datasize_count += stream.DataSize
				head++
				if deadline == 0 || stream.Deadline < deadline {
					deadline = stream.Deadline
				}

				if datasize_count >= encap_config.Datasize_Max {
					method.flushStream(can2ttFlow, current_time, datasize_count, deadline, encap_config.Datasize_Least)
					datasize_count = 0
					queue.popQueueByHead(head)
					head = 0
					deadline = 0
				}
			}

			if datasize_count > 0 && current_time%can2ttFlow.Period == 0 {
				method.flushStream(can2ttFlow, current_time, datasize_count, deadline, encap_config.Datasize_Least)
				datasize_count = 0
				queue.popQueueByHead(head)
				head = 0
				deadline = 0
			}
		}
		if datasize_count > 0 {
			method.flushStream(can2ttFlow, can2ttFlow.HyperPeriod, datasize_count, deadline, encap_config.Datasize_Least)
			datasize_count = 0
			deadline = 0
		}
	}
}

func (method *Method) encap_obo(can2ttClusterPool *ClusterPool, encap_config *Encapsulate_Config) {
	for _, can_stream_cluster := range can2ttClusterPool.Clusters {
		queue := new_Queue()
		can2ttFlow := method.getCAN2TTFlowByDomain(can_stream_cluster.Source, can_stream_cluster.Destination)

		for current_time := 0; current_time < can2ttFlow.HyperPeriod; current_time += encap_config.Step {
			queue.appendQueue(can_stream_cluster.getStreamsByCurrentTime(current_time))

			// fristly, drop overdue streams
			method.CAN2TT_O1_Drop += queue.checkDrop(current_time)

			// secondly, encapsulate streams
			for len(queue.Streams) > 0 {
				method.flushStream(can2ttFlow, current_time, encap_config.Datasize_Least, queue.Streams[0].Deadline, encap_config.Datasize_Least)
				queue.popQueueByIdx(0)
			}

		}
	}
}

func (method *Method) encap_wst(can2ttClusterPool *ClusterPool, encap_config *Encapsulate_Config) {
	const guardBase = 1800 // µs
	for _, can_stream_cluster := range can2ttClusterPool.Clusters {
		queue := new_Queue()
		can2ttFlow := method.getCAN2TTFlowByDomain(can_stream_cluster.Source, can_stream_cluster.Destination)

		datasize_count := 0.
		deadline := 0
		for current_time := 0; current_time < can2ttFlow.HyperPeriod; current_time += encap_config.Step {
			queue.appendQueue(can_stream_cluster.getStreamsByCurrentTime(current_time))
			queue.sortQueue(method.Method_Name, current_time)

			// fristly, drop overdue streams
			method.CAN2TT_O1_Drop += queue.checkDrop(current_time)

			// secondly, encapsulate streams
			guard := guardBase + len(queue.Streams)*600
			if len(queue.Streams) == 0 {
				continue
			}
			for queue.hasImminent(current_time, guard) {
				head := 0
				for head < len(queue.Streams) && datasize_count+queue.Streams[head].DataSize < encap_config.Datasize_Max {
					if deadline == 0 || queue.Streams[head].Deadline < deadline {
						deadline = queue.Streams[head].Deadline
					}

					datasize_count += queue.Streams[head].DataSize
					head++
				}

				method.flushStream(can2ttFlow, current_time, datasize_count, deadline, encap_config.Datasize_Least)
				queue.popQueueByHead(head)
				datasize_count = 0
				deadline = 0
			}
		}

		if len(queue.Streams) > 0 {
			method.flushStream(can2ttFlow, can2ttFlow.HyperPeriod, datasize_count, deadline, encap_config.Datasize_Least)
			datasize_count = 0
			deadline = 0
		}
	}
}

func (method *Method) encap_mao(can2ttClusterPool *ClusterPool, encap_config *Encapsulate_Config) {
	// step1
	MTU := 1500.0
	for _, can_stream_cluster := range can2ttClusterPool.Clusters {
		queue := new_Queue()
		can2ttFlow := method.getCAN2TTFlowByDomain(can_stream_cluster.Source, can_stream_cluster.Destination)

		deadline := 0
		datasize_count := 0.
		for current_time := 0; current_time < can2ttFlow.HyperPeriod; current_time += encap_config.Step {
			queue.appendQueue(can_stream_cluster.getStreamsByCurrentTime(current_time))
			queue.sortQueue(method.Method_Name, current_time)

			// fristly, drop overdue streams
			method.CAN2TT_O1_Drop += queue.checkDrop(current_time)

			// secondly, encapsulate streams
			head := 0
			for head < len(queue.Streams) {
				stream := queue.Streams[head]
				datasize_count += stream.DataSize
				head++
				if deadline == 0 || stream.Deadline < deadline {
					deadline = stream.Deadline
				}

				if datasize_count >= (MTU / 2) {
					method.flushStream(can2ttFlow, current_time, datasize_count, deadline, encap_config.Datasize_Least)
					datasize_count = 0
					queue.popQueueByHead(head)
					head = 0
					deadline = 0
				}
			}

		}

		if datasize_count >= 0 {
			method.flushStream(can2ttFlow, can2ttFlow.HyperPeriod, datasize_count, deadline, encap_config.Datasize_Least)
			datasize_count = 0
			deadline = 0
		}
	}

	// step2 Harmonic Merge
	for _, can2ttflow := range method.CAN2TTFlows {
		for _, tt_message := range can2ttflow.Streams {
			if tt_message.DataSize < (MTU / 2) {

			}
		}
	}
}

func (method *Method) flushStream(flow *Flow, now int, packedSize float64, dl int, Datasize_Least float64) {
	if packedSize < Datasize_Least {
		packedSize = Datasize_Least
	}
	stream := createCAN2TTStream(now, dl, packedSize+42)
	flow.Streams = append(flow.Streams, stream)

	method.BytesSent += packedSize + 42
	method.TTFrameCount += 1
}
