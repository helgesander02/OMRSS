package network

func (network *Network) ShowNetwork() {
	network.Topology.ShowTopology()
	network.FlowSet.ShowFlows()
	network.FlowSet.ShowTTFlow()
	network.FlowSet.ShowStream()
	network.GraphSet.ShowGraphs()
}
