package steninertree

type Edge struct {
    u, v, w int
}

type Graph struct {
    n     int      
    Edges []*Edge
	terminal []int 
}