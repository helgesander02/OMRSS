package plan

import (
	"src/network"
)

type Plan struct {
	Network   *network.Network // Network (1. Topology 2. Flows 3. Graphs)
	Algorithm *Algorithm       // A. OSACO (1. Ktree 2.WCD ...) B. SteinerTree ...
}

type Algorithm struct {
}
