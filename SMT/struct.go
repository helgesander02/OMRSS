package SMT

type Edge struct {
    u, v int
    w    float64
}

type Graph struct {
    n     int      
    Edges []*Edge
	terminal []int 
}
