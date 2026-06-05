package flow

import (
	"src/network/flow/can"
	"src/network/flow/tt"
	"src/pkg/config"
)

var (
	bgTSNFlowsEnd int
	bgAVBFlowsEnd int
)

// use the breakpoint to separate input and background flows
func fillFlowBreakpoint(bgTSN int, bgAVB int) {
	bgTSNFlowsEnd = bgTSN
	bgAVBFlowsEnd = bgAVB
}

func GenerateOmacoFlows(flows config.FlowsConfig, hyperPeriod int, nnodeLength int) *FlowSet {

	fillFlowBreakpoint(flows.TSN.Background, flows.AVB.Background)

	tsnFlows, avbFlows := tt.GenerateTTFlows(
		flows.TSN,
		flows.AVB,
		hyperPeriod,
		nnodeLength,
		flows.RoutingMode,
	)

	flowSet := newFlowSet()
	flowSet.TSNFlows = tsnFlows
	flowSet.AVBFlows = avbFlows

	return flowSet
}

func GenerateOsroFlows(flows config.FlowsConfig, hyperPeriod int, nnodeLength int, CANnode []int) *FlowSet {

	fillFlowBreakpoint(flows.TSN.Background, flows.AVB.Background)

	tsnFlows, avbFlows := tt.GenerateTTFlows(
		flows.TSN,
		flows.AVB,
		hyperPeriod,
		nnodeLength,
		flows.RoutingMode,
	)
	methodSet := can.GenerateCAN2TTFlows(
		flows.CAN,
		hyperPeriod,
		CANnode,
	)

	flowSet := newFlowSet()
	flowSet.TSNFlows = tsnFlows
	flowSet.AVBFlows = avbFlows
	flowSet.EncapsulateMethod = methodSet

	return flowSet
}
