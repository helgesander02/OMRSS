package flow

func (flows *FlowSet) InputOMACOFlowSet() *FlowSet {
	inputFlowSet := newFlowSet()
	inputFlowSet.TSNFlows = append(inputFlowSet.TSNFlows, flows.TSNFlows[bgTSNFlowsEnd:]...)
	inputFlowSet.AVBFlows = append(inputFlowSet.AVBFlows, flows.AVBFlows[bgAVBFlowsEnd:]...)

	return inputFlowSet
}

func (flows *FlowSet) BGOMACOFlowSet() *FlowSet {
	bgFlowSet := newFlowSet()
	bgFlowSet.TSNFlows = append(bgFlowSet.TSNFlows, flows.TSNFlows[:bgTSNFlowsEnd]...)
	bgFlowSet.AVBFlows = append(bgFlowSet.AVBFlows, flows.AVBFlows[:bgAVBFlowsEnd]...)

	return bgFlowSet
}
