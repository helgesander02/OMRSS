package algo

import (
	"fmt"
	"src/network"
	"src/plan/routes"
)

func SMT_Run(network *network.Network, show_smt bool) *routes.Trees_set {
	// 4. SteinerTree
	SMT := routes.Get_SteninerTree_Routing(network)

	if show_smt {
		fmt.Println()
		fmt.Println("--- The Steiner Tree final selected routing---")
		Input_SMT := SMT.Input_SteinerTree_set()
		BG_SMT := SMT.BG_SteinerTree_set()
		Input_SMT.Show_Trees_Set()
		BG_SMT.Show_Trees_Set()
	}

	return SMT
}
