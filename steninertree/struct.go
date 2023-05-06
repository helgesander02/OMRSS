package steninertree

type Edge struct {
    u, v, w int
}

type Graph struct {
    edges []Edge
    n, m  int
}