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

func (method *Method) EncapsulateCAN2TT(forwardingEngine *ForwardingEngine) {
	encapConfig := newEncapsulateConfig()
	if len(method.CAN2TTFlows) == 0 {
		method.organizeCAN2TTFlows(forwardingEngine)
	}

	switch method.MethodName {
	case "obo":
		method.encapOBO(forwardingEngine, encapConfig)
	case "wst":
		method.encapWST(forwardingEngine, encapConfig)
	case "mao":
		method.encapMAO(forwardingEngine, encapConfig)
	default:
		method.encapFIFOOrPriority(forwardingEngine, encapConfig)
	}
}

func (method *Method) organizeCAN2TTFlows(forwardingEngine *ForwardingEngine) {
	for _, canBUS := range forwardingEngine.BUSs {
		if !method.flowExists(canBUS.Source, canBUS.Destination) {
			can2ttFlow := newCAN2TTFlow()
			can2ttFlow.Source = canBUS.Source
			can2ttFlow.Destination = canBUS.Destination
			can2ttFlow.Period = canBUS.Period
			can2ttFlow.Deadline = canBUS.Deadline
			can2ttFlow.DataSize = canBUS.DataSize
			can2ttFlow.HyperPeriod = canBUS.HyperPeriod

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

func (method *Method) encapFIFOOrPriority(forwardingEngine *ForwardingEngine, encapConfig *EncapsulateConfig) {
	for _, canBUS := range forwardingEngine.BUSs {
		queue := newQueue()
		can2ttFlow := method.getCAN2TTFlowByDomain(canBUS.Source, canBUS.Destination)

		deadline := 0
		datasizeCount := 0.
		for currentTime := 0; currentTime < can2ttFlow.HyperPeriod; currentTime += encapConfig.Step {
			queue.appendQueue(canBUS.getFramesByCurrentTime(currentTime))
			queue.sortQueue(method.MethodName, currentTime)

			// firstly, drop overdue frames
			method.CAN2TTO1Drop += queue.checkDrop(currentTime)

			// secondly, encapsulate frames
			head := 0
			for head < len(queue.Frames) {
				frame := queue.Frames[head]
				datasizeCount += frame.DataSize
				head++
				if deadline == 0 || frame.Deadline < deadline {
					deadline = frame.Deadline
				}

				if datasizeCount >= encapConfig.DatasizeMax {
					method.flushFrame(can2ttFlow, currentTime, datasizeCount, deadline, encapConfig.DatasizeLeast)
					datasizeCount = 0
					queue.popQueueByHead(head)
					head = 0
					deadline = 0
				}
			}

			if datasizeCount > 0 && currentTime%can2ttFlow.Period == 0 {
				method.flushFrame(can2ttFlow, currentTime, datasizeCount, deadline, encapConfig.DatasizeLeast)
				datasizeCount = 0
				queue.popQueueByHead(head)
				head = 0
				deadline = 0
			}
		}
		if datasizeCount > 0 {
			method.flushFrame(can2ttFlow, can2ttFlow.HyperPeriod, datasizeCount, deadline, encapConfig.DatasizeLeast)
			datasizeCount = 0
			deadline = 0
		}
	}
}

func (method *Method) encapOBO(forwardingEngine *ForwardingEngine, encapConfig *EncapsulateConfig) {
	for _, canBUS := range forwardingEngine.BUSs {
		queue := newQueue()
		can2ttFlow := method.getCAN2TTFlowByDomain(canBUS.Source, canBUS.Destination)

		for currentTime := 0; currentTime < can2ttFlow.HyperPeriod; currentTime += encapConfig.Step {
			queue.appendQueue(canBUS.getFramesByCurrentTime(currentTime))

			// firstly, drop overdue frames
			method.CAN2TTO1Drop += queue.checkDrop(currentTime)

			// secondly, encapsulate frames
			for len(queue.Frames) > 0 {
				method.flushFrame(can2ttFlow, currentTime, encapConfig.DatasizeLeast, queue.Frames[0].Deadline, encapConfig.DatasizeLeast)
				queue.popQueueByIdx(0)
			}

		}
	}
}

func (method *Method) encapWST(forwardingEngine *ForwardingEngine, encapConfig *EncapsulateConfig) {
	const guardBase = 1800 // µs
	for _, canBUS := range forwardingEngine.BUSs {
		queue := newQueue()
		can2ttFlow := method.getCAN2TTFlowByDomain(canBUS.Source, canBUS.Destination)

		datasizeCount := 0.
		deadline := 0
		for currentTime := 0; currentTime < can2ttFlow.HyperPeriod; currentTime += encapConfig.Step {
			queue.appendQueue(canBUS.getFramesByCurrentTime(currentTime))
			queue.sortQueue(method.MethodName, currentTime)

			// firstly, drop overdue frames
			method.CAN2TTO1Drop += queue.checkDrop(currentTime)

			// secondly, encapsulate frames
			guard := guardBase + len(queue.Frames)*600
			if len(queue.Frames) == 0 {
				continue
			}
			for queue.hasImminent(currentTime, guard) {
				head := 0
				for head < len(queue.Frames) && datasizeCount+queue.Frames[head].DataSize < encapConfig.DatasizeMax {
					if deadline == 0 || queue.Frames[head].Deadline < deadline {
						deadline = queue.Frames[head].Deadline
					}

					datasizeCount += queue.Frames[head].DataSize
					head++
				}

				method.flushFrame(can2ttFlow, currentTime, datasizeCount, deadline, encapConfig.DatasizeLeast)
				queue.popQueueByHead(head)
				datasizeCount = 0
				deadline = 0
			}
		}

		if len(queue.Frames) > 0 {
			method.flushFrame(can2ttFlow, can2ttFlow.HyperPeriod, datasizeCount, deadline, encapConfig.DatasizeLeast)
			datasizeCount = 0
			deadline = 0
		}
	}
}

func (method *Method) encapMAO(forwardingEngine *ForwardingEngine, encapConfig *EncapsulateConfig) {
	// step1: MAO Aggregation
	MTU := 1500.0
	for _, canBUS := range forwardingEngine.BUSs {
		queue := newQueue()
		can2ttFlow := method.getCAN2TTFlowByDomain(canBUS.Source, canBUS.Destination)

		deadline := 0
		datasizeCount := 0.
		for currentTime := 0; currentTime < can2ttFlow.HyperPeriod; currentTime += encapConfig.Step {
			queue.appendQueue(canBUS.getFramesByCurrentTime(currentTime))
			queue.sortQueue(method.MethodName, currentTime)

			// firstly, drop overdue frames
			method.CAN2TTO1Drop += queue.checkDrop(currentTime)

			// secondly, encapsulate frames
			head := 0
			for head < len(queue.Frames) {
				frame := queue.Frames[head]
				datasizeCount += frame.DataSize
				head++
				if deadline == 0 || frame.Deadline < deadline {
					deadline = frame.Deadline
				}

				if datasizeCount >= (MTU / 2) {
					method.flushFrame(can2ttFlow, currentTime, datasizeCount, deadline, encapConfig.DatasizeLeast)
					datasizeCount = 0
					queue.popQueueByHead(head)
					head = 0
					deadline = 0
				}
			}

		}

		if datasizeCount >= 0 {
			method.flushFrame(can2ttFlow, can2ttFlow.HyperPeriod, datasizeCount, deadline, encapConfig.DatasizeLeast)
			datasizeCount = 0
			deadline = 0
		}
	}

	// step2: Harmonic Merge
	method.harmonicMerge(MTU)
}

func (method *Method) flushFrame(flow *Flow, now int, packedSize float64, dl int, datasizeLeast float64) {
	if packedSize < datasizeLeast {
		packedSize = datasizeLeast
	}
	frame := createCAN2TTFrame(now, dl, packedSize+42)
	flow.Frames = append(flow.Frames, frame)

	method.BytesSent += packedSize + 42
	method.TTFrameCount += 1
}
