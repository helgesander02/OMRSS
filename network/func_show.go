package network

// ShowNetwork prints the network dump. Summary lines (flow counts, CAN
// encap summary + comparison table) always print; when verbose is true,
// the long sections (topology cost matrix, per-frame timing, per-emit CAN
// detail, per-flow routing graphs) are added.
func (network *Network) ShowNetwork(verbose bool) {
	if verbose {
		network.Topology.ShowTopology()
	}
	network.FlowSet.ShowFlows()
	network.FlowSet.ShowTTFlow()
	if verbose {
		network.FlowSet.ShowFrame()
	}
	// Only OSRO populates EncapsulateMethod; for OMACO this is a no-op.
	network.FlowSet.ShowCANFlow(verbose)
	if verbose {
		network.GraphSet.ShowGraphs()
	}
}
