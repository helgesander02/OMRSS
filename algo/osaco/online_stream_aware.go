package osaco

import (
	"fmt"
	"src/network"
	"src/routes"
	"time"
)

// Ching-Chih Chuang et al., "Online Stream-Aware Routing for TSN-Based Industrial Control Systems"
func Run(network *network.Network, K int, show_osaco bool, time_out int) (*routes.Trees_set, *routes.Trees_set, [4]float64, *routes.Trees_set, *routes.Trees_set, [4]float64) {
	// 4. SteinerTree  5. OSACO
	fmt.Println("OSACO and Steiner Tree")
	fmt.Println("----------------------------------------")
	X := routes.Get_OSACO_Routing(network, K)
	II := X.Input_SteinerTree_set()
	II_prime := X.BG_SteinerTree_set()
	Input_SMT := II
	BG_SMT := II_prime

	pheromone := compute_prm(X)
	visibility := compute_vb(X, network.Flow_Set)

	if show_osaco {
		fmt.Printf("\n--- %dth Spanning Tree ---\n", K)
		X.Show_kTrees_Set()
		fmt.Println("--- Visibility and Pheromone ---")
		fmt.Println(visibility)
		fmt.Println(pheromone)
	}

	// Repeat the execution of epochs within the timeout
	initialobj, initialcost := obj(network, X, II, II_prime)
	fmt.Println()
	fmt.Printf("initial value: %d \n", initialcost)
	fmt.Printf("O1: %f O2: %f O3: pass O4: %f \n", initialobj[0], initialobj[1], initialobj[3])

	timeout := time.Duration(time_out) * time.Millisecond
	startTime := time.Now()
	i := 1
	for {
		fmt.Printf("\nepoch%d:\n", i)
		IIV := epoch(network, X, II_prime, visibility, pheromone)
		_, cost1 := obj(network, X, IIV, II_prime) // new
		_, cost2 := obj(network, X, II, II_prime)  // old
		if cost1 < cost2 {
			II = IIV
			fmt.Println("Change the selected routing !!")
		}
		i += 1

		if time.Since(startTime) >= timeout {
			break
		}
	}

	resultobj, resultcost := obj(network, X, II, II_prime)
	fmt.Println()
	fmt.Printf("result value: %d \n", resultcost)
	fmt.Printf("O1: %f O2: %f O3: pass O4: %f \n", resultobj[0], resultobj[1], resultobj[3])

	// SteinerTree (Input_SMT, BG_SMT, initialobj), OSACO (II, II_prime, resultobj)
	return Input_SMT, BG_SMT, initialobj, II, II_prime, resultobj
}

func epoch(network *network.Network, X *routes.KTrees_set, II_prime *routes.Trees_set, VB *Visibility, PRM *Pheromone) *routes.Trees_set {
	IIV, _, input_k_location, _ := probability(X, VB, PRM)
	//IIV, IIV_prime, input_k_location, bg_k_location := Probability(X, VB, PRM) // BG ... pass
	fmt.Printf("Select input routing %v \n", input_k_location)
	//fmt.Printf("Select background routing %v \n", bg_k_location) // BG ... pass

	_, cost := obj(network, X, IIV, II_prime)
	//obj, cost := Obj(network, X, IIV, IIV_prime) // BG ... pass

	p := 0.9 //0 <= p <= 1
	for nth, ktree := range X.TSNTrees {
		for kth := range ktree.Trees {
			if nth >= (len(X.TSNTrees) / 2) { // BG ... pass
				//PRM.TSN_PRM[nth][kth] *= p
				//if kth == bg_k_location[0][nth] {
				//	PRM.TSN_PRM[nth][kth] += (1 / cost[3])
				//}
			} else { // Input
				PRM.TSN_PRM[nth][kth] *= p
				if kth == input_k_location[0][nth] {
					PRM.TSN_PRM[nth][kth] += float64(1 / cost)
				}
			}
		}
	}

	for nth, ktree := range X.AVBTrees {
		for kth := range ktree.Trees {
			if nth >= (len(X.AVBTrees) / 2) { // BG ... pass
				//PRM.AVB_PRM[nth][kth] *= p
				//if kth == bg_k_location[1][nth] {
				//	PRM.AVB_PRM[nth][kth] += (1 / cost[3])
				//}
			} else { // Input
				PRM.AVB_PRM[nth][kth] *= p
				if kth == input_k_location[1][nth] {
					PRM.AVB_PRM[nth][kth] += float64(1 / cost)
				}
			}
		}
	}

	return IIV
}
