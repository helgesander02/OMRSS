package algo

import (
	"src/network"
	"src/plan/routes"
)

func (smt *SMT) SMT_Run(network *network.OMACO_Network) {
	// 4. SteinerTree
	smt.Trees = routes.Get_SteninerTree_Routing(network)
}
