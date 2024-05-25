package algo

import (
	"src/network"
	"src/plan/routes"
)

func (mtdc *MDTC) MDTC_Run(network *network.Network) {
	// 5. DistanceTree
	mtdc.Trees = routes.Get_DistanceTree_Routing(network)
}
