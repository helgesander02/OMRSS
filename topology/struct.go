package topology

type Node struct {
    ID            int           
    Name          string
	Connections   []*Connection               
}

type Connection struct {
    FromNodeID    int         // strat
    ToNodeID      int         // next
    Cost          float64     // (125,000,000 bytes/s) 1Gbps => (750,000 bytes/6ms) 750,000 bytes under 6ms for each link
}
 
type Topology struct {
    Talker        []*Node
    Switch        []*Node
    Listener      []*Node   
    Nodes         []*Node 
}   