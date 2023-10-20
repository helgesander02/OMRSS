package osaco

import (
	"fmt"
	"src/network"
	"src/routes"
)

// Ching-Chih Chuang et al., "Online Stream-Aware Routing for TSN-Based Industrial Control Systems"
func Run(network *network.Network, K int, show_osaco bool) (*routes.Trees_set, *routes.Trees_set, *routes.Trees_set, *routes.Trees_set) {
	// 4. SteinerTree and  5. OSACO
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
		fmt.Println(visibility)
		fmt.Println(pheromone)
	}

	IIV, IIV_prime := Epoch(network, X, visibility, pheromone)

	if show_osaco {
		IIV.Show_Trees_Set()
		IIV_prime.Show_Trees_Set()
	}

	// SteinerTree (Input_SMT, BG_SMT), OSACO (II, II_prime)
	return Input_SMT, BG_SMT, II, II_prime
}

func Epoch(network *network.Network, X *routes.KTrees_set, VB *Visibility, PRM *Pheromone) (*routes.Trees_set, *routes.Trees_set) {
	//p := 6. //0 <= p <= 1
	IIV, IIV_prime, input_k_location, bg_k_location := Probability(X, VB, PRM)
	fmt.Printf("Input select routes %v \n", input_k_location)
	fmt.Printf("BG select routes %v \n", bg_k_location)

	cost := Obj(network, X, IIV, IIV_prime)
	p := 6.

	for nth, ktree := range X.TSNTrees {
		for kth, _ := range ktree.Trees {
			if nth >= (len(X.TSNTrees) / 2) {
				PRM.TSN_PRM[nth][kth] *= p
				if kth == bg_k_location[0][nth] {
					PRM.TSN_PRM[nth][kth] += (1 / cost)
				}
			} else {
				PRM.TSN_PRM[nth][kth] *= p
				if kth == input_k_location[0][nth] {
					PRM.TSN_PRM[nth][kth] += (1 / cost)
				}
			}
		}
	}

	for nth, ktree := range X.AVBTrees {
		for kth, _ := range ktree.Trees {
			if nth >= (len(X.AVBTrees) / 2) {
				PRM.AVB_PRM[nth][kth] *= p
				if kth == bg_k_location[1][nth] {
					PRM.AVB_PRM[nth][kth] += (1 / cost)
				}
			} else {
				PRM.AVB_PRM[nth][kth] *= p
				if kth == input_k_location[1][nth] {
					PRM.AVB_PRM[nth][kth] += (1 / cost)
				}
			}
		}
	}

	return IIV, IIV_prime
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
