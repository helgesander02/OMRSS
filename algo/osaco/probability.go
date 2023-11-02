package osaco

import (
	"crypto/rand"
	"math/big"
	"src/routes"
)

func probability(X *routes.KTrees_set, VB *Visibility, PRM *Pheromone) (*routes.Trees_set, *routes.Trees_set, [2][]int, [2][]int) {
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

		n := 0
		var arr []int
		for kth := range ktree.Trees {
			probability := (VB.TSN_VB[nth][kth] * PRM.TSN_PRM[nth][kth]) / Denominator
			for j := 0; j < int(probability*100); j++ {
				arr = append(arr, kth)
			}
		}
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(arr))))
		n = arr[int(randomIndex.Int64())]
		t := ktree.Trees[n]

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

		n := 0
		var arr []int
		for kth := range ktree.Trees {
			probability := (VB.AVB_VB[nth][kth] * PRM.AVB_PRM[nth][kth]) / Denominator
			for j := 0; j < int(probability*100); j++ {
				arr = append(arr, kth)
			}
		}
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(arr))))
		n = arr[int(randomIndex.Int64())]
		t := ktree.Trees[n]

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
