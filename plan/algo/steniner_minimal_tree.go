package algo

import (
	"src/network"
	"src/pkg/config"
	"src/plan/routes"
)

// SMT_Run computes Steiner Minimal Trees for OMACO network
func (smt *SMT) SMT_Run(network *network.Network, cfg *config.Config) {
	smt.Trees = routes.Get_SteninerTree_Routing(network, cfg)
}
