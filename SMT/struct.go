package SMT

type Graphs struct {
    Graphs []*Graph 
}

type Graph struct {
    ID int
    Visited bool
    Cost int
    Path int
    Edges []*Edge
}

type Edge struct {
    Strat    int         
    End      int         
    Cost     int     
}