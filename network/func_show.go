package network

func (network *Network) ShowNetwork() {
	network.Topology.ShowTopology()
	network.FlowSet.ShowFlows()
	network.FlowSet.ShowTTFlow()
	network.FlowSet.ShowFrame()
	network.GraphSet.ShowGraphs()
}
