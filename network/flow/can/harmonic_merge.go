package can

// Stage 2 of the MAO algorithm (Yan et al. 2024, Algorithm 1, lines 24-29):
// take the Stage 1 TSN messages that did NOT fill up to MTU/2 (the paper's
// set F') and merge those sharing src/dst with harmonic periods into a
// single denser TSN message. Merged-flow properties follow paper Eqs (5)-(7):
//   Period   = gcd(periods)
//   Deadline = min(deadlines)
//   DataSize = sum(payload sizes)

// gcd returns the greatest common divisor of two positive integers.
func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// isHarmonicPeriod reports whether one period is an integer multiple of
// the other. Equal periods count as harmonic.
func isHarmonicPeriod(p1, p2 int) bool {
	if p1 == 0 || p2 == 0 {
		return false
	}
	if p1 < p2 {
		p1, p2 = p2, p1
	}
	return p1%p2 == 0
}

// FlowGroup is a working set of flows being considered for harmonic merge.
type FlowGroup struct {
	Flows       []*Flow
	TotalSize   float64
	MinDeadline int
	GCDPeriod   int
}

// harmonicMerge runs the second aggregation pass of MAO. maxPayload is the
// per-TSN-message payload cap (paper invariant: MTU/2). Only flows whose
// Stage 1 per-emit payload is below maxPayload (i.e., the paper's set F')
// are candidates for merging.
func (method *Method) harmonicMerge(maxPayload float64) {
	// Snapshot each flow's real Stage 1 per-emit payload into flow.DataSize
	// so the "inefficient" classification below reflects the actual emit
	// size, not the original CAN frame size (which is what
	// organizeCAN2TTFlows initialised DataSize to).
	for _, flow := range method.CAN2TTFlows {
		if len(flow.Frames) == 0 {
			continue
		}
		var totalPayload float64
		for _, frame := range flow.Frames {
			totalPayload += frame.DataSize - HeaderBytes
		}
		flow.DataSize = totalPayload / float64(len(flow.Frames))
	}

	// 1. Partition flows into efficient (>= maxPayload) and inefficient
	//    (< maxPayload). Only inefficient flows are merge candidates.
	inefficient := make([]*Flow, 0)
	efficient := make([]*Flow, 0)
	for _, flow := range method.CAN2TTFlows {
		if flow.DataSize < maxPayload {
			inefficient = append(inefficient, flow)
		} else {
			efficient = append(efficient, flow)
		}
	}
	if len(inefficient) == 0 {
		return
	}

	// 2. Greedily group mergeable inefficient flows.
	groups := method.findMergeableGroups(inefficient, maxPayload)

	// 3. For each group with at least two flows, build one merged flow.
	merged := make(map[*Flow]bool)
	newFlows := make([]*Flow, 0)
	for _, group := range groups {
		if len(group.Flows) > 1 {
			newFlows = append(newFlows, method.mergeFlowGroup(group))
			for _, flow := range group.Flows {
				merged[flow] = true
			}
		}
	}

	// 4. Rebuild CAN2TTFlows: efficient flows stay, ungrouped inefficient
	//    flows stay, merged results are appended.
	final := make([]*Flow, 0, len(method.CAN2TTFlows))
	final = append(final, efficient...)
	for _, flow := range inefficient {
		if !merged[flow] {
			final = append(final, flow)
		}
	}
	final = append(final, newFlows...)
	method.CAN2TTFlows = final
}

// findMergeableGroups greedily collects flows into groups where every pair
// shares (src, dst), all periods are pairwise harmonic, and the cumulative
// payload stays within maxPayload.
func (method *Method) findMergeableGroups(flows []*Flow, maxPayload float64) []*FlowGroup {
	groups := make([]*FlowGroup, 0)
	processed := make(map[*Flow]bool)

	for i := 0; i < len(flows); i++ {
		if processed[flows[i]] {
			continue
		}

		group := &FlowGroup{
			Flows:       []*Flow{flows[i]},
			TotalSize:   flows[i].DataSize,
			MinDeadline: flows[i].Deadline,
			GCDPeriod:   flows[i].Period,
		}
		processed[flows[i]] = true

		for j := i + 1; j < len(flows); j++ {
			if processed[flows[j]] {
				continue
			}
			if method.canAddToGroup(group, flows[j], maxPayload) {
				group.Flows = append(group.Flows, flows[j])
				group.TotalSize += flows[j].DataSize
				if flows[j].Deadline < group.MinDeadline {
					group.MinDeadline = flows[j].Deadline
				}
				group.GCDPeriod = gcd(group.GCDPeriod, flows[j].Period)
				processed[flows[j]] = true
			}
		}

		groups = append(groups, group)
	}

	return groups
}

// canAddToGroup checks whether flow may join group: identical (src, dst),
// pairwise harmonic period with every existing member, and the combined
// payload stays within maxPayload (paper invariant PAY <= MTU/2).
func (method *Method) canAddToGroup(group *FlowGroup, flow *Flow, maxPayload float64) bool {
	if len(group.Flows) == 0 {
		return false
	}

	first := group.Flows[0]
	if flow.Source != first.Source || !sameDestinations(flow.Destinations, first.Destinations) {
		return false
	}

	for _, member := range group.Flows {
		if !isHarmonicPeriod(flow.Period, member.Period) {
			return false
		}
	}

	if group.TotalSize+flow.DataSize > maxPayload {
		return false
	}

	return true
}

// mergeFlowGroup builds the merged Flow from a group, then regenerates the
// Frames at the merged (GCD) period. The original Stage 1 Frames are
// discarded by design — the merged Flow represents the post-MAO logical
// view consumed by the downstream scheduler; BytesSent / TTFrameCount on
// the Method still account for the actual Stage 1 emission cost.
func (method *Method) mergeFlowGroup(group *FlowGroup) *Flow {
	if len(group.Flows) == 0 {
		return nil
	}

	first := group.Flows[0]
	merged := &Flow{
		Source:       first.Source,
		Destinations: append([]int{}, first.Destinations...),
		Period:       group.GCDPeriod,   // paper Eq (5)
		Deadline:     group.MinDeadline, // paper Eq (6)
		DataSize:    group.TotalSize,   // paper Eq (7), payload only
		HyperPeriod: first.HyperPeriod,
		Frames:      make([]*Frame, 0),
	}
	merged.Frames = method.regenerateFrames(group, merged.Period, merged.HyperPeriod)
	return merged
}

// regenerateFrames produces one Frame per period boundary in the
// hyperperiod, each carrying the full aggregated payload. This is the
// paper's conservative reservation model (Eq (5) + Eq (7)): bandwidth is
// reserved at the GCD rate even when not every component flow has data
// to send at every tick.
func (method *Method) regenerateFrames(group *FlowGroup, newPeriod int, hyperPeriod int) []*Frame {
	frames := make([]*Frame, 0)
	for currentTime := 0; currentTime < hyperPeriod; currentTime += newPeriod {
		frames = append(frames, &Frame{
			ArrivalTime: currentTime,
			Deadline:    group.MinDeadline,
			DataSize:    group.TotalSize,
			FinishTime:  currentTime + group.MinDeadline,
		})
	}
	return frames
}
