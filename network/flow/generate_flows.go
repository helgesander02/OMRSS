package flow

import (
	"src/network/flow/can"
	"src/network/flow/tt"
)

var (
	bg_tsnflows_end int
	bg_avbflows_end int
)

func Generate_OMACO_Flows(Nnode_length int, bg_tsn int, bg_avb int, input_tsn int, input_avb int, HyperPeriod int) *Flow_Set {
	bg_tsnflows_end = bg_tsn
	bg_avbflows_end = bg_avb

	tsn_flows, avb_flows := tt.Generate_TT_Flows(Nnode_length, bg_tsn, bg_avb, input_tsn, input_avb, HyperPeriod)

	flow_set := new_Flow_Set()
	flow_set.TSNFlows = tsn_flows
	flow_set.AVBFlows = avb_flows

	return flow_set
}

func (flows *Flow_Set) Input_OMACO_Flow_Set() *Flow_Set {
	Input_flow_set := new_Flow_Set()
	Input_flow_set.TSNFlows = append(Input_flow_set.TSNFlows, flows.TSNFlows[bg_tsnflows_end:]...)
	Input_flow_set.AVBFlows = append(Input_flow_set.AVBFlows, flows.AVBFlows[bg_avbflows_end:]...)

	return Input_flow_set
}

func (flows *Flow_Set) BG_OMACO_Flow_Set() *Flow_Set {
	BG_flow_set := new_Flow_Set()
	BG_flow_set.TSNFlows = append(BG_flow_set.TSNFlows, flows.TSNFlows[:bg_tsnflows_end]...)
	BG_flow_set.AVBFlows = append(BG_flow_set.AVBFlows, flows.AVBFlows[:bg_avbflows_end]...)

	return BG_flow_set
}

func Generate_OSRO_Flows(CANnode []int, importantCAN int, unimportantCAN int, Nnode_length int, bg_tsn int, bg_avb int, input_tsn int, input_avb int, HyperPeriod int) *Flow_Set {
	bg_tsnflows_end = bg_tsn
	bg_avbflows_end = bg_avb

	tsn_flows, avb_flows := tt.Generate_TT_Flows(Nnode_length, bg_tsn, bg_avb, input_tsn, input_avb, HyperPeriod)
	method_set := can.Generate_CAN2TT_Flows(CANnode, importantCAN, unimportantCAN, HyperPeriod)

	flow_set := new_Flow_Set()
	flow_set.TSNFlows = tsn_flows
	flow_set.AVBFlows = avb_flows
	flow_set.EncapsulateMethod = method_set

	return flow_set
}
