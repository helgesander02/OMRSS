package topology

type Node struct {
	ID          int
	Connections []*Connection
}

type Connection struct {
	FromNodeID int     // strat
	ToNodeID   int     // next
	Cost       float64 // 1Gbps => (750,000 bytes/6ms) 750,000 bytes under 6ms for each link ==> 125 bytes/ns
}

type Topology struct {
	Talker   []*Node
	Switch   []*Node
	Listener []*Node
	Nodes    []*Node
}
