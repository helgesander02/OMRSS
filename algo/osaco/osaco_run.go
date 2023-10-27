package osaco

import (
	"fmt"
	"src/network"
	"src/routes"
)

// Ching-Chih Chuang et al., "Online Stream-Aware Routing for TSN-Based Industrial Control Systems"
func Run(network *network.Network, K int, show_osaco bool) (*routes.Trees_set, *routes.Trees_set, *routes.Trees_set, *routes.Trees_set) {
	// 4. SteinerTree  5. OSACO
	fmt.Printf("\nOSACO \n")
	fmt.Println("----------------------------------------")
	X := routes.GetRoutes(network, K)
	II := X.Input_SteinerTree_set()
	II_prime := X.BG_SteinerTree_set()
	Input_SMT := II
	BG_SMT := II_prime

	pheromone := CompPRM(X)
	visibility := CompVB(X, network.Flow_Set)

	if show_osaco {
		fmt.Printf("\nSteiner Tree and %dth Spanning Tree \n", K)
		X.Show_kTrees_Set()
		fmt.Println("Visibility and Pheromone")
		fmt.Println(visibility)
		fmt.Println(pheromone)
	}

	fmt.Printf("\nepoch%d:\n", 1)
	IIV := Epoch(network, X, II_prime, visibility, pheromone)

	if show_osaco {
		IIV.Show_Trees_Set()
	}

	// SteinerTree (Input_SMT, BG_SMT), OSACO (II, II_prime)
	return Input_SMT, BG_SMT, II, II_prime
}

func Epoch(network *network.Network, X *routes.KTrees_set, II_prime *routes.Trees_set, VB *Visibility, PRM *Pheromone) *routes.Trees_set {
	IIV, _, input_k_location, _ := Probability(X, VB, PRM)
	// IIV, IIV_prime, input_k_location, bg_k_location := Probability(X, VB, PRM) // BG ... pass
	fmt.Printf("Select input routes %v \n", input_k_location)
	// fmt.Printf("Select background routes %v \n", bg_k_location) // BG ... pass

	cost := Obj(network, X, IIV, II_prime)
	// cost := Obj(network, X, IIV, IIV_prime) // BG ... pass
	fmt.Printf("O1: %f O2: %f O3: pass O4: %f\n", cost[0], cost[1], cost[3])

	p := 6. //0 <= p <= 1
	for nth, ktree := range X.TSNTrees {
		for kth := range ktree.Trees {
			if nth >= (len(X.TSNTrees) / 2) { // BG ... pass
				//PRM.TSN_PRM[nth][kth] *= p
				//if kth == bg_k_location[0][nth] {
				//PRM.TSN_PRM[nth][kth] += (1 / cost[3])
				//}
			} else { // Input
				PRM.TSN_PRM[nth][kth] *= p
				if kth == input_k_location[0][nth] {
					PRM.TSN_PRM[nth][kth] += (1 / cost[3])
				}
			}
		}
	}

	for nth, ktree := range X.AVBTrees {
		for kth := range ktree.Trees {
			if nth >= (len(X.AVBTrees) / 2) { // BG ... pass
				//PRM.AVB_PRM[nth][kth] *= p
				//if kth == bg_k_location[1][nth] {
				//PRM.AVB_PRM[nth][kth] += (1 / cost[3])
				//}
			} else { // Input
				PRM.AVB_PRM[nth][kth] *= p
				if kth == input_k_location[1][nth] {
					PRM.AVB_PRM[nth][kth] += (1 / cost[3])
				}
			}
		}
	}

	return IIV
}

func Probability(X *routes.KTrees_set, VB *Visibility, PRM *Pheromone) (*routes.Trees_set, *routes.Trees_set, [2][]int, [2][]int) {
	var (
		input_k_location [2][]int
		bg_k_location    [2][]int
		bg_tsn_start     int = len(X.TSNTrees) / 2
		bg_avb_start     int = len(X.AVBTrees) / 2
	)

	IIV := &routes.Trees_set{}
	IIV_prime := &routes.Trees_set{}
	for nth, ktree := range X.TSNTrees {
		Denominator := 0.
		for kth := range ktree.Trees {
			Denominator += VB.TSN_VB[nth][kth] * PRM.TSN_PRM[nth][kth]
		}

		t := ktree.Trees[0]
		n := 0
		p := 0.
		for kth, tree := range ktree.Trees {
			probability := (VB.TSN_VB[nth][kth] * PRM.TSN_PRM[nth][kth]) / Denominator
			if Denominator == 0. {
				probability = 0.
			}

			if probability > p {
				p = probability
				n = kth
				t = tree
			}
		}

		if nth >= bg_tsn_start {
			bg_k_location[0] = append(bg_k_location[0], n)
			IIV_prime.TSNTrees = append(IIV_prime.TSNTrees, t)
		} else {
			input_k_location[0] = append(input_k_location[0], n)
			IIV.TSNTrees = append(IIV.TSNTrees, t)
		}
	}

	for nth, ktree := range X.AVBTrees {
		Denominator := 0.
		for kth := range ktree.Trees {
			Denominator += VB.AVB_VB[nth][kth] * PRM.AVB_PRM[nth][kth]
		}
		t := ktree.Trees[0]
		n := 0
		p := 0.
		for kth, tree := range ktree.Trees {
			probability := (VB.AVB_VB[nth][kth] * PRM.AVB_PRM[nth][kth]) / Denominator
			if Denominator == 0. {
				probability = 0.
			}

			if probability > p {
				p = probability
				n = kth
				t = tree
			}
		}

		if nth >= bg_avb_start {
			bg_k_location[1] = append(bg_k_location[1], n)
			IIV_prime.AVBTrees = append(IIV_prime.AVBTrees, t)
		} else {
			input_k_location[1] = append(input_k_location[1], n)
			IIV.AVBTrees = append(IIV.AVBTrees, t)
		}
	}

	return IIV, IIV_prime, input_k_location, bg_k_location
}
