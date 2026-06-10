package network

func (network *Network) ShowNetwork(verbose bool) {
	if verbose {
		network.Topology.ShowTopology()
	}

	network.FlowSet.ShowFlows()
	network.FlowSet.ShowTTFlow()
	if verbose {
		network.FlowSet.ShowFrame()
	}

	network.FlowSet.ShowCANFlow(verbose)
	if verbose {
		network.GraphSet.ShowGraphs()
	}
}
