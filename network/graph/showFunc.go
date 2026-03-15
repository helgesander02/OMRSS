package graph

func (graphs *Graphs) ShowGraphs() {
	for _, graph := range graphs.TSNGraphs {
		graph.ShowTopology()
		break
	}

	for _, graph := range graphs.AVBGraphs {
		graph.ShowTopology()
		break
	}

	for _, graph := range graphs.CAN2TSNGraphs {
		graph.ShowTopology()
		break
	}
}
