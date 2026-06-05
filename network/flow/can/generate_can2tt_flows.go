package can

import (
	"time"

	"src/pkg/config"
	"src/pkg/logger"
)

func GenerateCAN2TTFlows(can config.CANFlowConfig, hyperperiod int, CANnode []int) []*Method {
	// step 1: generate CAN flows
	importantCANFlows, unimportantCANFlows := GenerateCANFlows(can, hyperperiod, CANnode)

	// step2: prepare method list
	var methodList = []string{"fifo", "priority", "obo", "wst", "mao"}

	// step3: according to different encapsulation methods, generate CAN2TT flows
	methodSet := newMethodSet()
	for _, methodName := range methodList {
		forwardingEngine := newForwardingEngine()

		if methodName == "mao" {
			for _, impf := range importantCANFlows {
				flowCopy := impf.deepCopyFlow()
				forwardingEngine.aggregateCANFrameByPeriodAndDomain(flowCopy)
			}
			for _, unimpf := range unimportantCANFlows {
				flowCopy := unimpf.deepCopyFlow()
				forwardingEngine.aggregateCANFrameByPeriodAndDomain(flowCopy)
			}

		} else {
			for _, impf := range importantCANFlows {
				flowCopy := impf.deepCopyFlow()
				forwardingEngine.aggregateCANFrameByDomain(flowCopy)
			}
			for _, unimpf := range unimportantCANFlows {
				flowCopy := unimpf.deepCopyFlow()
				forwardingEngine.aggregateCANFrameByDomain(flowCopy)
			}
		}

		start := time.Now()
		method := newMethod(methodName)
		method.EncapsulateCAN2TT(forwardingEngine)
		method.CAN2TSNDelay = time.Since(start)
		methodSet = append(methodSet, method)
	}

	return methodSet
}

type ForwardingEngine struct {
	BUSs []*CANBUS
}

func newForwardingEngine() *ForwardingEngine {
	return &ForwardingEngine{}
}

func (forwardingEngine *ForwardingEngine) aggregateCANFrameByDomain(f *Flow) {
	for _, canBUS := range forwardingEngine.BUSs {
		if canBUS.Source == f.Source && canBUS.Destination == f.Destination {
			canBUS.Frames = append(canBUS.Frames, f.Frames...)
			return
		}
	}
	forwardingEngine.addNewBus(f)
}

func (forwardingEngine *ForwardingEngine) aggregateCANFrameByPeriodAndDomain(f *Flow) {
	for _, canBUS := range forwardingEngine.BUSs {
		if canBUS.Period == f.Period && canBUS.Source == f.Source && canBUS.Destination == f.Destination {
			canBUS.Frames = append(canBUS.Frames, f.Frames...)
			return
		}
	}
	forwardingEngine.addNewBus(f)
}

func (forwardingEngine *ForwardingEngine) addNewBus(f *Flow) {
	canBUS := newCANBUS()
	canBUS.Source = f.Source
	canBUS.Destination = f.Destination
	canBUS.Period = f.Period
	canBUS.Deadline = f.Deadline
	canBUS.DataSize = f.DataSize
	canBUS.HyperPeriod = f.HyperPeriod
	canBUS.Frames = append(canBUS.Frames, f.Frames...)

	forwardingEngine.BUSs = append(forwardingEngine.BUSs, canBUS)
}

func (forwardingEngine *ForwardingEngine) ShowForwardingEngine() {
	logger.Println("CAN2TT Traffic Router:")
	for _, canBUS := range forwardingEngine.BUSs {
		canBUS.ShowCANBUS()
	}
}

type CANBUS struct {
	Flow
}

func newCANBUS() *CANBUS {
	return &CANBUS{}
}

func (canBUS *CANBUS) getFramesByCurrentTime(currentTime int) []*Frame {
	frames := []*Frame{}
	for _, frame := range canBUS.Frames {
		if frame.ArrivalTime == currentTime {
			frames = append(frames, frame)
		}
	}
	return frames
}

func (canBUS *CANBUS) ShowCANBUS() {
	logger.Printf("Queue (%d→%d) frames=%d\n", canBUS.Source, canBUS.Destination, len(canBUS.Frames))
	logger.Printf("Period: %v  ,Deadline: %v ,Datasize: %v\n", canBUS.Period, canBUS.Deadline, canBUS.DataSize)
}
