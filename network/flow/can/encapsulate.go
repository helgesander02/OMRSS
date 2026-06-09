package can

const (
	MinPayloadBytes = 64.                    // smallest TT frame payload (bytes)
	MaxPayloadBytes = 1500.                  // largest TT frame payload = MTU (bytes)
	HeaderBytes     = 42.                    // per-TT-frame overhead (Ethernet + IP + UDP)
	Step            = 1000                   // encapsulation time step (µs)
	CANBandwidth    = 1000000.0              // CAN bus bandwidth (1 Mbps)
	BytesPerStep    = (1000000.0 / 8) * 1e-3 // bytes the CAN bus can deliver per Step
)

func (method *Method) EncapsulateCAN2TT(agg *CAN2TTAggregator) {
	if len(method.CAN2TTFlows) == 0 {
		method.organizeCAN2TTFlows(agg)
	}

	switch method.MethodName {
	case MethodOBO:
		method.encapOBO(agg)
	case MethodWST:
		method.encapWST(agg)
	case MethodMAO:
		method.encapMAO(agg)
	case MethodWSTMAR:
		method.encapWSTMAR(agg)
	default:
		method.encapFIFOOrPriority(agg)
	}
}

// organizeCAN2TTFlows materialises one CAN2TT flow per aggregator group.
// The identity key is (src, dst, period) — without the period component,
// MAO's distinct (period, src, dst) groups would collapse into a single
// flow and Stage 2 (harmonic merge) would have nothing to merge.
// FIFO / Priority / WST / OBO each only have a single group per (src, dst)
// in their aggregator, so the extra period dimension is a no-op for them.
func (method *Method) organizeCAN2TTFlows(agg *CAN2TTAggregator) {
	for _, group := range agg.Groups {
		if !method.flowExists(group.Source, group.Destination, group.Period) {
			can2ttFlow := newCAN2TTFlow()
			can2ttFlow.Source = group.Source
			can2ttFlow.Destination = group.Destination
			can2ttFlow.Period = group.Period
			can2ttFlow.Deadline = group.Deadline
			can2ttFlow.DataSize = group.DataSize
			can2ttFlow.HyperPeriod = group.HyperPeriod

			method.CAN2TTFlows = append(method.CAN2TTFlows, can2ttFlow)
		}
	}
}

// getCAN2TTFlow returns the flow whose key matches (src, dst, period).
func (method *Method) getCAN2TTFlow(source int, destination int, period int) *Flow {
	for _, flow := range method.CAN2TTFlows {
		if flow.Source == source && flow.Destination == destination && flow.Period == period {
			return flow
		}
	}
	return nil
}

func (method *Method) flowExists(source int, destination int, period int) bool {
	return method.getCAN2TTFlow(source, destination, period) != nil
}

func (method *Method) encapFIFOOrPriority(agg *CAN2TTAggregator) {
	for _, group := range agg.Groups {
		queue := newQueue()
		can2ttFlow := method.getCAN2TTFlow(group.Source, group.Destination, group.Period)

		deadline := 0
		payloadBytes := 0.
		for currentTime := 0; currentTime < can2ttFlow.HyperPeriod; currentTime += Step {
			queue.appendQueue(group.getFramesByCurrentTime(currentTime))
			queue.sortQueue(method.MethodName, currentTime)

			// firstly, drop overdue frames
			method.CAN2TTO1Drop += queue.checkDrop(currentTime)

			// secondly, drain the queue into the batch
			head := 0
			for head < len(queue.Frames) {
				frame := queue.Frames[head]

				// 1. flush BEFORE adding if this frame would push us past MTU,
				//    so the emitted payload never exceeds MaxPayloadBytes.
				if payloadBytes > 0 && payloadBytes+frame.DataSize > MaxPayloadBytes {
					method.emitCAN2TTFrame(can2ttFlow, currentTime, payloadBytes, deadline)
					queue.popQueueByHead(head)
					payloadBytes = 0
					deadline = 0
					head = 0
					continue
				}

				payloadBytes += frame.DataSize
				if deadline == 0 || frame.Deadline < deadline {
					deadline = frame.Deadline
				}
				head++
			}
			// Always pop what we drained into the batch so the queue starts
			// empty next step (prevents cross-step double-counting of frames
			// that haven't been flushed yet).
			queue.popQueueByHead(head)

			// 2. reach time step
			if payloadBytes > 0 && currentTime%can2ttFlow.Period == 0 {
				method.emitCAN2TTFrame(can2ttFlow, currentTime, payloadBytes, deadline)
				payloadBytes = 0
				deadline = 0
			}
		}
		// 3. reach hyperperiod, send remaining frames
		if payloadBytes > 0 {
			method.emitCAN2TTFrame(can2ttFlow, can2ttFlow.HyperPeriod, payloadBytes, deadline)
			payloadBytes = 0
			deadline = 0
		}
	}
}

func (method *Method) encapOBO(agg *CAN2TTAggregator) {
	for _, group := range agg.Groups {
		queue := newQueue()
		can2ttFlow := method.getCAN2TTFlow(group.Source, group.Destination, group.Period)

		for currentTime := 0; currentTime < can2ttFlow.HyperPeriod; currentTime += Step {
			queue.appendQueue(group.getFramesByCurrentTime(currentTime))

			// firstly, drop overdue frames
			method.CAN2TTO1Drop += queue.checkDrop(currentTime)

			// secondly, encapsulate frames
			for len(queue.Frames) > 0 {
				// 1. one by one
				method.emitCAN2TTFrame(can2ttFlow, currentTime, MinPayloadBytes, queue.Frames[0].Deadline)
				queue.popQueueByIdx(0)
			}

		}
	}
}

func (method *Method) encapWST(agg *CAN2TTAggregator) {
	const (
		guardBase     = 1800                // µs base for urgency guard
		softThreshold = MaxPayloadBytes / 2 // emit-when-reached threshold (MTU/2)
	)

	for _, group := range agg.Groups {
		queue := newQueue()
		can2ttFlow := method.getCAN2TTFlow(group.Source, group.Destination, group.Period)

		payloadBytes := 0.
		deadline := 0      // min relative deadline of batched frames (written into the emitted TT frame)
		minFinishTime := 0 // min absolute deadline of batched frames (used for urgency check)
		batchSize := 0     // number of CAN frames currently held in the batch

		for currentTime := 0; currentTime < can2ttFlow.HyperPeriod; currentTime += Step {
			queue.appendQueue(group.getFramesByCurrentTime(currentTime))
			queue.sortQueue(method.MethodName, currentTime)

			// firstly, drop overdue frames
			method.CAN2TTO1Drop += queue.checkDrop(currentTime)

			// secondly, drain the queue into the batch
			head := 0
			for head < len(queue.Frames) {
				frame := queue.Frames[head]

				// Trigger (1) — hard MTU limit: flush BEFORE adding the frame
				// so the emitted payload never exceeds MaxPayloadBytes.
				if payloadBytes > 0 && payloadBytes+frame.DataSize > MaxPayloadBytes {
					method.emitCAN2TTFrame(can2ttFlow, currentTime, payloadBytes, deadline)
					queue.popQueueByHead(head)
					payloadBytes = 0
					deadline = 0
					minFinishTime = 0
					batchSize = 0
					head = 0
					continue
				}

				payloadBytes += frame.DataSize
				if deadline == 0 || frame.Deadline < deadline {
					deadline = frame.Deadline
				}
				if minFinishTime == 0 || frame.FinishTime < minFinishTime {
					minFinishTime = frame.FinishTime
				}
				batchSize++
				head++
			}
			// Pop everything we drained so the queue is empty when the next
			// step starts. This is what prevents the same frame from being
			// re-counted into payloadBytes on subsequent iterations.
			queue.popQueueByHead(head)

			// Decide whether to emit the current batch this step.
			if payloadBytes > 0 {
				guard := guardBase + batchSize*600
				urgent := minFinishTime-currentTime <= guard
				saturated := payloadBytes >= softThreshold

				if urgent || saturated {
					method.emitCAN2TTFrame(can2ttFlow, currentTime, payloadBytes, deadline)
					payloadBytes = 0
					deadline = 0
					minFinishTime = 0
					batchSize = 0
				}
				// else: hold the batch over to the next step.
			}
		}

		// Hyperperiod tail flush: emit whatever is still held in the batch.
		if payloadBytes > 0 {
			method.emitCAN2TTFrame(can2ttFlow, can2ttFlow.HyperPeriod, payloadBytes, deadline)
			payloadBytes = 0
			deadline = 0
			minFinishTime = 0
			batchSize = 0
		}
	}
}

func (method *Method) encapMAO(agg *CAN2TTAggregator) {
	// MAO emits a TT frame as soon as the aggregated payload reaches MTU/2;
	// harmonicMerge afterwards will combine these small emits when possible.
	const maoThreshold = MaxPayloadBytes / 2

	// step1: MAO Aggregation
	for _, group := range agg.Groups {
		queue := newQueue()
		can2ttFlow := method.getCAN2TTFlow(group.Source, group.Destination, group.Period)

		deadline := 0
		payloadBytes := 0.
		for currentTime := 0; currentTime < can2ttFlow.HyperPeriod; currentTime += Step {
			queue.appendQueue(group.getFramesByCurrentTime(currentTime))
			queue.sortQueue(method.MethodName, currentTime)

			// firstly, drop overdue frames
			method.CAN2TTO1Drop += queue.checkDrop(currentTime)

			// secondly, drain the queue into the batch
			head := 0
			for head < len(queue.Frames) {
				frame := queue.Frames[head]

				// Flush BEFORE adding if this frame would push us past the
				// MAO threshold, so the emitted payload stays <= MTU/2.
				if payloadBytes > 0 && payloadBytes+frame.DataSize > maoThreshold {
					method.emitCAN2TTFrame(can2ttFlow, currentTime, payloadBytes, deadline)
					queue.popQueueByHead(head)
					payloadBytes = 0
					deadline = 0
					head = 0
					continue
				}

				payloadBytes += frame.DataSize
				if deadline == 0 || frame.Deadline < deadline {
					deadline = frame.Deadline
				}
				head++
			}
			// Always pop what we drained into the batch so the queue starts
			// empty next step (prevents cross-step double-counting of frames
			// that haven't reached the MAO threshold yet).
			queue.popQueueByHead(head)
		}

		// tail flush: only emit if there is actually something to send.
		if payloadBytes > 0 {
			method.emitCAN2TTFrame(can2ttFlow, can2ttFlow.HyperPeriod, payloadBytes, deadline)
			payloadBytes = 0
			deadline = 0
		}
	}

	// step2: Harmonic Merge — paper enforces PAY <= MTU/2 throughout MAO,
	// so the merged-flow size cap is MaxPayloadBytes/2, not MaxPayloadBytes.
	method.harmonicMerge(MaxPayloadBytes / 2)
}

// encapWSTMAR aggregates CAN frames inside one (src, dst) bucket using three
// layered policies. The name MAR stands for "Multi-tier Adaptive Release":
//
//  1. EDF ordering — the (src, dst) queue is sorted by FinishTime ascending,
//     so the most urgent frame is always at the head. Frames with identical
//     FinishTime form a natural cohort (same slack) and are pulled together
//     during a single drain pass.
//
//  2. Slack-tiered emit timing — the drain loop always packs up to the hard
//     MTU (so urgent frames never starve), but WHEN to emit the batch depends
//     on the smallest slack currently held:
//
//     slack ≤ slackUrgent → emit immediately, no size check
//     slack ≤ slackSoft   → emit once payload ≥ MTU/2 (soft target)
//     slack >  slackSoft  → keep holding until the batch hits MTU
//
//  3. Peek-ahead merge — if a soft-tier emit would fire but the next CAN
//     arrival in this (src, dst) is close enough that the batch can survive
//     waiting for it (slack at next arrival is still above slackUrgent), the
//     emit is suppressed and the batch is held to merge the upcoming wave.
//     This is the move that turns MAR into a winner over MAO: it exploits the
//     deterministic periodicity of CAN traffic to combine consecutive cohorts
//     into a single TT frame, without sacrificing safety.
//
// MAO buckets by arrival rate (period); MAR orders by remaining time and uses
// that remaining time, combined with the next known arrival, to decide both
// the emit threshold and the emit moment dynamically.
func (method *Method) encapWSTMAR(agg *CAN2TTAggregator) {
	const (
		slackUrgent = 2000                  // µs — emit as-is, do not wait
		slackSoft   = 8000                  // µs — emit when payload ≥ softTarget
		softTarget  = MaxPayloadBytes / 2.0 // 750 B
	)

	for _, group := range agg.Groups {
		queue := newQueue()
		can2ttFlow := method.getCAN2TTFlow(group.Source, group.Destination, group.Period)

		payloadBytes := 0.0
		deadline := 0      // min relative deadline of batched frames
		minFinishTime := 0 // min absolute deadline of batched frames (drives emit timing)

		for currentTime := 0; currentTime < can2ttFlow.HyperPeriod; currentTime += Step {
			queue.appendQueue(group.getFramesByCurrentTime(currentTime))
			queue.sortQueue(MethodWSTMAR, currentTime)

			// firstly, drop overdue frames
			method.CAN2TTO1Drop += queue.checkDrop(currentTime)

			// secondly, drain the queue up to the hard MTU. The slack-aware
			// gate is in the emit decision below, NOT here — limiting the
			// drain would starve urgent frames piled up across steps.
			head := 0
			for head < len(queue.Frames) {
				frame := queue.Frames[head]

				// Hard MTU emit: flush before adding so the emitted payload
				// never exceeds MaxPayloadBytes.
				if payloadBytes > 0 && payloadBytes+frame.DataSize > MaxPayloadBytes {
					method.emitCAN2TTFrame(can2ttFlow, currentTime, payloadBytes, deadline)
					queue.popQueueByHead(head)
					payloadBytes = 0
					deadline = 0
					minFinishTime = 0
					head = 0
					continue
				}

				payloadBytes += frame.DataSize
				if deadline == 0 || frame.Deadline < deadline {
					deadline = frame.Deadline
				}
				if minFinishTime == 0 || frame.FinishTime < minFinishTime {
					minFinishTime = frame.FinishTime
				}
				head++
			}
			queue.popQueueByHead(head)

			// Decide whether to emit the current batch this step.
			if payloadBytes > 0 {
				slack := minFinishTime - currentTime

				shouldEmit := false
				switch {
				case slack <= slackUrgent:
					// urgent: do not wait, send whatever we have.
					shouldEmit = true
				case slack <= slackSoft:
					// medium pressure: emit at soft target to limit risk.
					shouldEmit = payloadBytes >= softTarget
				default:
					// loose: keep packing toward MTU. Hard MTU emit was
					// already handled inside the drain loop above.
					shouldEmit = false
				}

				// Peek-ahead override (soft tier): if a soft-tier emit would fire
				// but the next CAN arrival is close enough that the batch can
				// still hold a safe slack margin after it lands, hold the batch
				// so the upcoming wave can be merged into the same TT frame.
				if shouldEmit && slack > slackUrgent {
					nextArrival := group.nextArrivalAfter(currentTime)
					if nextArrival > 0 && nextArrival < can2ttFlow.HyperPeriod {
						slackAtNext := minFinishTime - nextArrival
						if slackAtNext > slackUrgent {
							shouldEmit = false
						}
					}
				}

				// Peek-ahead override (urgent tier): even when the batch is
				// already urgent, if the next CAN arrival lands at or before
				// minFinishTime (i.e. before any held frame would be dropped),
				// hold one more step so the new wave can be merged in. Drops
				// are still avoided because checkDrop fires only when
				// currentTime > FinishTime, so emitting AT minFinishTime is
				// still safe. This is the move that fuses consecutive
				// monoperiodic CAN waves into a single TT frame.
				if shouldEmit && slack <= slackUrgent {
					nextArrival := group.nextArrivalAfter(currentTime)
					if nextArrival > 0 && nextArrival <= minFinishTime {
						shouldEmit = false
					}
				}

				if shouldEmit {
					method.emitCAN2TTFrame(can2ttFlow, currentTime, payloadBytes, deadline)
					payloadBytes = 0
					deadline = 0
					minFinishTime = 0
				}
				// else: hold the batch over to the next step.
			}
		}

		// Hyperperiod tail flush: emit whatever is still held in the batch.
		if payloadBytes > 0 {
			method.emitCAN2TTFrame(can2ttFlow, can2ttFlow.HyperPeriod, payloadBytes, deadline)
			payloadBytes = 0
			deadline = 0
			minFinishTime = 0
		}
	}
}

func (method *Method) emitCAN2TTFrame(flow *Flow, currentTime int, payloadBytes float64, deadline int) {
	if payloadBytes < MinPayloadBytes {
		payloadBytes = MinPayloadBytes
	}

	frame := createCAN2TTFrame(currentTime, deadline, payloadBytes+HeaderBytes)
	flow.Frames = append(flow.Frames, frame)

	method.BytesSent += payloadBytes + HeaderBytes
	method.TTFrameCount += 1
}
