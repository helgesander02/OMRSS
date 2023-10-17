package osaco

import (
	"fmt"
	"src/network"
	"src/routes"
)

// Ching-Chih Chuang et al., "Online Stream-Aware Routing for TSN-Based Industrial Control Systems"
func Run(network *network.Network, K int, show_routes bool, show_osaco bool) {
	// 4. Steiner Tree, K Spanning Tree "routes/generate_trees.go"
	fmt.Printf("\nSteiner Tree and %dth Spanning Tree \n", K)
	fmt.Println("----------------------------------------")
	X := routes.GetRoutes(network, K)
	II := X.Input_SteinerTree_set()
	II_prime := X.BG_SteinerTree_set()

	if show_routes {
		fmt.Println("Steiner Tree")
		II.Show_Trees_Set()
		II_prime.Show_Trees_Set()

		fmt.Println("K Spanning Tree")
		X.Show_kTrees_Set()
	}

	// 5. OSACO
	fmt.Printf("\nOSACO \n")
	fmt.Println("----------------------------------------")
	visibility := CompVB(X, network.Flow_Set)
	fmt.Println(visibility)
}
