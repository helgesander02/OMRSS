package network

func (network *Network) Show_Network() {
	network.Topology.Show_Topology()

	network.Flow_Set.Show_Flows()
	network.Flow_Set.Show_Flow()
	network.Flow_Set.Show_Stream()

	network.Graph_Set.Show_Graphs()
}

func (network *OSRO_Network) Show_Network() {

}
