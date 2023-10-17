package flow

func (flows *Flows) Input_flow_set() *Flows {
	Input_flow_set := &Flows{}

	var (
		input_tsn_end int = len(flows.TSNFlows) / 2
		input_avb_end int = len(flows.AVBFlows) / 2
	)

	Input_flow_set.TSNFlows = append(Input_flow_set.TSNFlows, flows.TSNFlows[:input_tsn_end]...)
	Input_flow_set.AVBFlows = append(Input_flow_set.AVBFlows, flows.AVBFlows[:input_avb_end]...)

	return Input_flow_set
}

func (flows *Flows) BG_flow_set() *Flows {
	BG_flow_set := &Flows{}
	var (
		bg_tsn_start int = len(flows.TSNFlows) / 2
		bg_avb_start int = len(flows.AVBFlows) / 2
	)

	BG_flow_set.TSNFlows = append(BG_flow_set.TSNFlows, flows.TSNFlows[bg_tsn_start:]...)
	BG_flow_set.AVBFlows = append(BG_flow_set.AVBFlows, flows.AVBFlows[bg_avb_start:]...)

	return BG_flow_set
}
