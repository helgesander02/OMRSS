package algo

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"src/network"
	"src/network/flow"
	"src/plan/routes"
	"src/plan/schedule"
	"time"
)

// Ching-Chih Chuang et al., "Online Stream-Aware Routing for TSN-Based Industrial Control Systems"
func (osaco *OSACO) OSACO_Run(network *network.Network, firsttime int) [4]float64 {
	// 6. OSACO
	if firsttime == 0 {
		osaco.PRM = compute_prm(osaco.KTrees)
		osaco.VB = compute_vb(osaco.KTrees, network.Flow_Set)
	}

	// Repeat the execution of epochs within the timeout
	initialobj, initialcost := schedule.OBJ(network, osaco.KTrees, osaco.InputTrees, osaco.BGTrees)
	fmt.Println()
	fmt.Printf("initial value: %d \n", initialcost)
	fmt.Printf("O1: %f O2: %f O3: pass O4: %f \n", initialobj[0], initialobj[1], initialobj[3])

	timeout := time.Duration(osaco.Timeout) * time.Millisecond
	startTime := time.Now()
	i := 1
	for {
		fmt.Printf("\nepoch%d:\n", i)
		II := epoch(network, osaco)
		_, cost1 := schedule.OBJ(network, osaco.KTrees, II, osaco.BGTrees)               // new
		_, cost2 := schedule.OBJ(network, osaco.KTrees, osaco.InputTrees, osaco.BGTrees) // old
		if cost1 < cost2 {
			osaco.InputTrees = II
			fmt.Println("Change the selected routing !!")
		}
		i += 1

		if time.Since(startTime) >= timeout {
			break
		}
	}

	resultobj, resultcost := schedule.OBJ(network, osaco.KTrees, osaco.InputTrees, osaco.BGTrees)
	fmt.Println()
	fmt.Printf("result value: %d \n", resultcost)
	fmt.Printf("O1: %f O2: %f O3: pass O4: %f \n", resultobj[0], resultobj[1], resultobj[3])
	fmt.Println()

	return resultobj
}

func compute_prm(X *routes.KTrees_set) *Pheromone {
	var (
		bg_tsn_start int = len(X.TSNTrees) / 2
		bg_avb_start int = len(X.AVBTrees) / 2
	)

	pheromone := &Pheromone{}
	for nth, ktree := range X.TSNTrees {
		var prm []float64
		for i := 0; i < len(ktree.Trees); i++ {
			if nth >= bg_tsn_start {
				prm = append(prm, 0.5)
			} else {
				prm = append(prm, 1.)
			}
		}
		pheromone.TSN_PRM = append(pheromone.TSN_PRM, prm)
	}

	for nth, ktree := range X.AVBTrees {
		var prm []float64
		for i := 0; i < len(ktree.Trees); i++ {
			if nth >= bg_avb_start {
				prm = append(prm, 0.5)
			} else {
				prm = append(prm, 1.)
			}
		}
		pheromone.AVB_PRM = append(pheromone.AVB_PRM, prm)
	}

	return pheromone
}

func compute_vb(X *routes.KTrees_set, flow_set *flow.Flows) *Visibility {
	Input_flow_set := flow_set.Input_flow_set()
	BG_flow_set := flow_set.BG_flow_set()

	var (
		preference   float64 = 2.
		bg_tsn_start int     = len(X.TSNTrees) / 2
		bg_avb_start int     = len(X.AVBTrees) / 2
	)

	visibility := &Visibility{}
	// OSACO CompVB
	// TSN flow
	for nth, tsn_ktree := range X.TSNTrees {
		var v []float64
		for kth := range tsn_ktree.Trees {
			mult := 1.
			if nth >= bg_tsn_start && kth == 0 {
				mult = preference
			}

			//value := mult / float64(tsn_ktree.Trees[0].Weight) // mult / Tree weight
			value := mult / math.Exp(float64(tsn_ktree.Trees[0].Weight)) // mult / exponential function( Tree weight )
			v = append(v, value)
		}
		visibility.TSN_VB = append(visibility.TSN_VB, v)
	}

	// OSACO CompVB
	// AVB flow
	for nth, avb_ktree := range X.AVBTrees {
		var v []float64
		for kth, z := range avb_ktree.Trees {
			mult := 1.
			if nth >= bg_avb_start && kth == 0 {
				mult = preference
			}

			if nth < bg_avb_start {
				//fmt.Printf("Input flow%d tree%d \n", nth, kth)
				value := mult / float64(schedule.WCD(z, X, Input_flow_set.AVBFlows[nth], flow_set))
				v = append(v, value)

			} else {
				//fmt.Printf("Backgourd flow%d tree%d \n", nth, kth)
				value := mult / float64(schedule.WCD(z, X, BG_flow_set.AVBFlows[nth-bg_avb_start], flow_set))
				v = append(v, value)
			}
		}
		visibility.AVB_VB = append(visibility.AVB_VB, v)
	}

	return visibility
}

func probability(osaco *OSACO) (*routes.Trees_set, *routes.Trees_set, [2][]int, [2][]int) {
	var (
		input_k_location [2][]int
		bg_k_location    [2][]int
		bg_tsn_start     int = len(osaco.KTrees.TSNTrees) / 2
		bg_avb_start     int = len(osaco.KTrees.AVBTrees) / 2
	)

	II := &routes.Trees_set{}
	II_prime := &routes.Trees_set{}
	for nth, ktree := range osaco.KTrees.TSNTrees {
		Denominator := 0.
		for kth := range ktree.Trees {
			Denominator += osaco.VB.TSN_VB[nth][kth] * osaco.PRM.TSN_PRM[nth][kth]
		}

		n := 0
		var arr []int
		for kth := range ktree.Trees {
			probability := (osaco.VB.TSN_VB[nth][kth] * osaco.PRM.TSN_PRM[nth][kth]) / Denominator
			for j := 0; j < int(probability*100); j++ {
				arr = append(arr, kth)
			}
		}
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(arr))))
		n = arr[int(randomIndex.Int64())]
		t := ktree.Trees[n]

		if nth >= bg_tsn_start {
			bg_k_location[0] = append(bg_k_location[0], n)
			II_prime.TSNTrees = append(II_prime.TSNTrees, t)
		} else {
			input_k_location[0] = append(input_k_location[0], n)
			II.TSNTrees = append(II.TSNTrees, t)
		}
	}

	for nth, ktree := range osaco.KTrees.AVBTrees {
		Denominator := 0.
		for kth := range ktree.Trees {
			Denominator += osaco.VB.AVB_VB[nth][kth] * osaco.PRM.AVB_PRM[nth][kth]
		}

		n := 0
		var arr []int
		for kth := range ktree.Trees {
			probability := (osaco.VB.AVB_VB[nth][kth] * osaco.PRM.AVB_PRM[nth][kth]) / Denominator
			for j := 0; j < int(probability*100); j++ {
				arr = append(arr, kth)
			}
		}
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(arr))))
		n = arr[int(randomIndex.Int64())]
		t := ktree.Trees[n]

		if nth >= bg_avb_start {
			bg_k_location[1] = append(bg_k_location[1], n)
			II_prime.AVBTrees = append(II_prime.AVBTrees, t)
		} else {
			input_k_location[1] = append(input_k_location[1], n)
			II.AVBTrees = append(II.AVBTrees, t)
		}
	}

	return II, II_prime, input_k_location, bg_k_location
}

func epoch(network *network.Network, osaco *OSACO) *routes.Trees_set {
	II, _, input_k_location, _ := probability(osaco)
	//II, II_prime, input_k_location, bg_k_location := Probability(osaco.KTrees, osaco.VB, osaco.PRM) // BG ... pass
	fmt.Printf("Select input routing %v \n", input_k_location)
	//fmt.Printf("Select background routing %v \n", bg_k_location) // BG ... pass

	_, cost := schedule.OBJ(network, osaco.KTrees, II, osaco.BGTrees)
	//obj, cost := Obj(network, X, II, II_prime) // BG ... pass

	for nth, ktree := range osaco.KTrees.TSNTrees {
		for kth := range ktree.Trees {
			if nth >= (len(osaco.KTrees.TSNTrees) / 2) { // BG ... pass
				//osaco.PRM.TSN_PRM[nth][kth] *= osaco.P
				//if kth == bg_k_location[0][nth] {
				//	osaco.PRM.TSN_PRM[nth][kth] += (1 / cost[3])
				//}
			} else { // Input
				osaco.PRM.TSN_PRM[nth][kth] *= osaco.P
				if kth == input_k_location[0][nth] {
					osaco.PRM.TSN_PRM[nth][kth] += float64(1 / cost)
				}
			}
		}
	}

	for nth, ktree := range osaco.KTrees.AVBTrees {
		for kth := range ktree.Trees {
			if nth >= (len(osaco.KTrees.AVBTrees) / 2) { // BG ... pass
				//osaco.PRM.AVB_PRM[nth][kth] *= osaco.P
				//if kth == bg_k_location[1][nth] {
				//	osaco.PRM.AVB_PRM[nth][kth] += (1 / cost[3])
				//}
			} else { // Input
				osaco.PRM.AVB_PRM[nth][kth] *= osaco.P
				if kth == input_k_location[1][nth] {
					osaco.PRM.AVB_PRM[nth][kth] += float64(1 / cost)
				}
			}
		}
	}

	return II
}
