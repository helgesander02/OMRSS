package flow

import (
	"src/network/flow/can"
	"src/network/flow/tt"
)

func GenerateOmacoFlows(nnodeLength int, bgTSN int, bgAVB int, inputTSN int, inputAVB int, hyperPeriod int) *FlowSet {

	fillFlowBreakpoint(bgTSN, bgAVB)

	tsnFlows, avbFlows := tt.GenerateTTFlows(nnodeLength, bgTSN, bgAVB, inputTSN, inputAVB, hyperPeriod)

	flowSet := newFlowSet()
	flowSet.TSNFlows = tsnFlows
	flowSet.AVBFlows = avbFlows

	return flowSet
}

func GenerateOsroFlows(nnodeLength int, bgTSN int, bgAVB int, inputTSN int, inputAVB int, hyperPeriod int, CANnode []int, importantCAN int, unimportantCAN int) *FlowSet {

	fillFlowBreakpoint(bgTSN, bgAVB)

	tsnFlows, avbFlows := tt.GenerateTTFlows(nnodeLength, bgTSN, bgAVB, inputTSN, inputAVB, hyperPeriod)
	methodSet := can.GenerateCAN2TTFlows(CANnode, importantCAN, unimportantCAN, hyperPeriod)

	flowSet := newFlowSet()
	flowSet.TSNFlows = tsnFlows
	flowSet.AVBFlows = avbFlows
	flowSet.EncapsulateMethod = methodSet

	return flowSet
}
