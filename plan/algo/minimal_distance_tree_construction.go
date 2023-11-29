package algo

import (
	"fmt"
	"src/network"
	"src/plan/routes"
)

func MDT_Run(network *network.Network, show_mdt bool) *routes.Trees_set {
	// 5. DistanceTree
	MDT := routes.Get_DistanceTree_Routing(network)

	if show_mdt {
		fmt.Println()
		fmt.Println("--- The Distance Tree final selected routing---")
		Input_MDT := MDT.Input_Tree_set()
		BG_MDT := MDT.BG_Tree_set()
		Input_MDT.Show_Trees_Set()
		BG_MDT.Show_Trees_Set()
	}

	return MDT
}
