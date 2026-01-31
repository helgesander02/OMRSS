package network

func (network *OMACO_Network) ShowNetwork() {
	network.Topology.Show_Topology()

	network.FlowSet.Show_Flows()
	network.FlowSet.Show_Flow()
	network.FlowSet.Show_Stream()

	network.Graph_Set.Show_Graphs()
}

func (network *OSRO_Network) ShowNetwork() {
	network.Topology.Show_Topology()

	network.FlowSet.Show_Flows()
	network.FlowSet.Show_Flow()
	network.FlowSet.Show_Stream()

	network.Graph_Set.Show_Graphs()
}
