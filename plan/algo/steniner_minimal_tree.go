package algo

import (
	"src/network"
	"src/plan/routes"
)

// SMT_Run computes Steiner Minimal Trees for OMACO network
func (smt *SMT) SMT_Run(network *network.OMACO_Network) {
	smt.Trees = routes.Get_SteninerTree_Routing(network)
}
