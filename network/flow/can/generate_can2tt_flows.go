package can

import (
	"time"

	"src/pkg/config"
	"src/pkg/logger"
)

const (
	MethodFIFO     = "fifo"
	MethodPriority = "priority"
	MethodOBO      = "obo"
	MethodWST      = "wst"
	MethodMAO      = "mao"
	MethodWSTMAR   = "wst_mar" // Multi-tier Adaptive Release: EDF order + slack-tiered emit + peek-ahead merge
)

var methodList = []string{
	MethodFIFO,
	MethodPriority,
	MethodOBO,
	MethodWST,
	MethodMAO,
	MethodWSTMAR,
}

// MethodList exposes the encap method ordering to other packages. Stays
// in sync with what GenerateCAN2TTFlows produces, so OSRO's per-method
// harness can iterate the same set the gateway actually emitted.
func MethodList() []string {
	out := make([]string, len(methodList))
	copy(out, methodList)
	return out
}

func GenerateCAN2TTFlows(can config.CANFlowConfig, hyperperiod int, CANnode []int, mode string) []*Method {
	// step 1: generate CAN flows
	importantCANFlows, unimportantCANFlows := GenerateCANFlows(can, hyperperiod, CANnode)

	// step 2: aggregate CAN frames per encapsulation method,
	// then run the encapsulation pass for each.
	methodSet := newMethodSet()
	for _, methodName := range methodList {
		agg := newCAN2TTAggregator()

		if methodName == MethodMAO {
			// MAO needs frames split per (Period, Source, Destination).
			for _, impf := range importantCANFlows {
				agg.aggregateByPeriodAndSrcDst(impf.deepCopyFlow())
			}
			for _, unimpf := range unimportantCANFlows {
				agg.aggregateByPeriodAndSrcDst(unimpf.deepCopyFlow())
			}
		} else {
			// All other strategies aggregate per (Source, Destination).
			for _, impf := range importantCANFlows {
				agg.aggregateBySrcDst(impf.deepCopyFlow())
			}
			for _, unimpf := range unimportantCANFlows {
				agg.aggregateBySrcDst(unimpf.deepCopyFlow())
			}
		}

		start := time.Now()
		method := newMethod(methodName)
		method.EncapsulateCAN2TT(agg)
		method.CAN2TTDelay = time.Since(start)
		methodSet = append(methodSet, method)
	}

	return methodSet
}

type CAN2TTAggregator struct {
	Groups []*CANFrameGroup
}

func newCAN2TTAggregator() *CAN2TTAggregator {
	return &CAN2TTAggregator{}
}

func (agg *CAN2TTAggregator) aggregateBySrcDst(f *Flow) {
	for _, group := range agg.Groups {
		if group.Source == f.Source && sameDestinations(group.Destinations, f.Destinations) {
			group.Frames = append(group.Frames, f.Frames...)
			return
		}
	}
	agg.addGroup(f)
}

func (agg *CAN2TTAggregator) aggregateByPeriodAndSrcDst(f *Flow) {
	for _, group := range agg.Groups {
		if group.Period == f.Period && group.Source == f.Source && sameDestinations(group.Destinations, f.Destinations) {
			group.Frames = append(group.Frames, f.Frames...)
			return
		}
	}
	agg.addGroup(f)
}

func (agg *CAN2TTAggregator) addGroup(f *Flow) {
	group := newCANFrameGroup()
	group.Source = f.Source
	group.Destinations = append([]int{}, f.Destinations...) // defensive copy
	group.Period = f.Period
	group.Deadline = f.Deadline
	group.DataSize = f.DataSize
	group.HyperPeriod = f.HyperPeriod
	group.Frames = append(group.Frames, f.Frames...)

	agg.Groups = append(agg.Groups, group)
}

// sameDestinations compares two destination lists for aggregator keying.
// Order matters because tt.Flow / graph entries treat destination order as
// significant; staying consistent here avoids surprising merges.
func sameDestinations(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func (agg *CAN2TTAggregator) Show() {
	logger.Println("CAN2TT Aggregator:")
	for _, group := range agg.Groups {
		group.Show()
	}
}

type CANFrameGroup struct {
	Flow
}

func newCANFrameGroup() *CANFrameGroup {
	return &CANFrameGroup{}
}

func (group *CANFrameGroup) getFramesByCurrentTime(currentTime int) []*Frame {
	frames := []*Frame{}
	for _, frame := range group.Frames {
		if frame.ArrivalTime == currentTime {
			frames = append(frames, frame)
		}
	}
	return frames
}

// nextArrivalAfter returns the earliest ArrivalTime in the group that is
// strictly greater than currentTime. Used by MAR's peek-ahead merge.
// Returns -1 when no future arrival exists.
func (group *CANFrameGroup) nextArrivalAfter(currentTime int) int {
	next := -1
	for _, frame := range group.Frames {
		if frame.ArrivalTime <= currentTime {
			continue
		}
		if next == -1 || frame.ArrivalTime < next {
			next = frame.ArrivalTime
		}
	}
	return next
}

func (group *CANFrameGroup) Show() {
	logger.Printf("Group (%d→%v) frames=%d\n", group.Source, group.Destinations, len(group.Frames))
	logger.Printf("Period: %v  ,Deadline: %v ,Datasize: %v\n", group.Period, group.Deadline, group.DataSize)
}
