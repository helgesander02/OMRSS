package graph

import "src/pkg/logger"

func (graphs *Graphs) ShowGraphs() {
	for i, entry := range graphs.Entries {
		if i >= 3 {
			break
		}
		logger.Printf("Graph entry %d  src=%d  dst=%v\n",
			i+1, entry.Source, entry.Destinations)
		entry.Graph.ShowTopology()
	}
}
