package graph

func (graphs *Graphs) Show_Graphs() {
	for _, graph := range graphs.TSNGraphs {
		graph.Show_Topology()
		break
	}

	for _, graph := range graphs.AVBGraphs {
		graph.Show_Topology()
		break
	}
}
