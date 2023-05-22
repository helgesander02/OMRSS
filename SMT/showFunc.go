package SMT

import (
	"fmt"
)

func (graphs *Graphs) Show_Graph() {
    for _, g := range graphs.Graphs {
        fmt.Printf("%d\t", g.ID)
        fmt.Printf("%v\t", g.Visited)
        fmt.Printf("%d\t", g.Cost)
        fmt.Printf("%d\n", g.Path)
		//for _, e := range g.Edges {
			//fmt.Printf("%d-->%d cost:%d\n", e.Strat, e.End, e.Cost)
		//}
    } 
}