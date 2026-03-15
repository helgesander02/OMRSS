package algo

import (
	"fmt"
	"math"
	"src/internal/config"
	"src/internal/random"
	"src/network"
	"src/network/flow"
	"src/plan/algo_timer"
	"src/plan/routes"
	"src/plan/schedule"
	"time"
)

var (
	bgTSN int
	bgAVB int
)

// Global RNG instance (will be set by plan package)
var rng *random.Generator

// SetRNG sets the random number generator for this package
func SetRNG(r *random.Generator) {
	rng = r
}

func (osaco *OSACO) OSACO_Initial_Settings(network *network.Network, cfg *config.Config, SMT *routes.TreesSet) {
	//// OSACO computing time: Estimate the time it takes to compute routing information
	bgTSN = cfg.Network.Flows.TSN.Background
	bgAVB = cfg.Network.Flows.AVB.Background

	timer := algo_timer.NewTimer()
	timer.TimerStart()
	osaco.KTrees = routes.Get_OSACO_Routing(network, cfg, SMT, osaco.K, osaco.Method_Number)
	timer.TimerEnd()

	osaco.InputTrees = SMT.InputTreeSet(bgTSN, bgAVB)
	osaco.BGTrees = SMT.BGTreeSet(bgTSN, bgAVB)
	osaco.PRM = computePrm(osaco.KTrees)
	osaco.VB = computeVb(osaco.KTrees, network.FlowSet)

	osaco.Timer[0] = algo_timer.NewTimer()
	osaco.Timer[0].TimerMerge(timer)
	osaco.Timer[1] = algo_timer.NewTimer()
	osaco.Timer[1].TimerMerge(timer)
	osaco.Timer[2] = algo_timer.NewTimer()
	osaco.Timer[2].TimerMerge(timer)
	osaco.Timer[3] = algo_timer.NewTimer()
	osaco.Timer[3].TimerMerge(timer)
	osaco.Timer[4] = algo_timer.NewTimer()
	osaco.Timer[4].TimerMerge(timer)
}

// Ching-Chih Chuang et al., "Online Stream-Aware Routing for TSN-Based Industrial Control Systems"
func (osaco *OSACO) OSACO_Run(network *network.Network, cfg *config.Config, timeoutIndex int, costSetting [4]int) [4]float64 {
	// Repeat the execution of epochs within the timeout
	initialobj, initialcost := schedule.OBJ(network, cfg, osaco.KTrees, osaco.InputTrees, osaco.BGTrees, costSetting, false)
	fmt.Println()
	fmt.Printf("initial value: %d \n", initialcost)
	fmt.Printf("O1: %f O2: %f O3: pass O4: %f \n", initialobj[0], initialobj[1], initialobj[3])

	timeout := time.Duration(osaco.Timeout) * time.Millisecond
	startTime := time.Now()
	i := 1
	for {
		fmt.Printf("\nepoch%d:\n", i)
		osaco.Timer[timeoutIndex].TimerStart()
		II := epoch(network, cfg, osaco, timeoutIndex, costSetting)
		osaco.Timer[timeoutIndex].TimerStop()

		_, cost1 := schedule.OBJ(network, cfg, osaco.KTrees, II, osaco.BGTrees, costSetting, false)               // new
		_, cost2 := schedule.OBJ(network, cfg, osaco.KTrees, osaco.InputTrees, osaco.BGTrees, costSetting, false) // old

		if cost1 < cost2 {
			osaco.InputTrees = II
			fmt.Println("Change the selected routing !!")
		}
		i += 1

		if time.Since(startTime) >= timeout {
			break
		}
	}

	resultobj, resultcost := schedule.OBJ(network, cfg, osaco.KTrees, osaco.InputTrees, osaco.BGTrees, costSetting, true)
	fmt.Println()
	fmt.Printf("result value: %d \n", resultcost)
	fmt.Printf("O1: %f O2: %f O3: pass O4: %f \n", resultobj[0], resultobj[1], resultobj[3])
	fmt.Println()

	if resultobj[0] != 0 || resultobj[1] != 0 {
		osaco.Timer[timeoutIndex].TimerMax()
	}

	return resultobj
}

func computePrm(X *routes.KTreesSet) *Pheromone {
	pheromone := &Pheromone{}

	for nth, ktree := range X.TSNTrees {
		var prm []float64
		for i := 0; i < len(ktree.Trees); i++ {
			if nth < bgTSN {
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
			if nth < bgAVB {
				prm = append(prm, 0.5)
			} else {
				prm = append(prm, 1.)
			}
		}
		pheromone.AVB_PRM = append(pheromone.AVB_PRM, prm)
	}

	return pheromone
}

func computeVb(X *routes.KTreesSet, flowSet *flow.FlowSet) *Visibility {
	var preference float64 = 2.
	Input_flowSet := flowSet.InputOMACOFlowSet()
	BG_flowSet := flowSet.BGOMACOFlowSet()

	visibility := &Visibility{}
	// OSACO CompVB
	// TSN flow
	for nth, tsn_ktree := range X.TSNTrees {
		var v []float64
		for kth := range tsn_ktree.Trees {
			mult := 1.
			if nth < bgTSN && kth == 0 {
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
			if nth < bgAVB && kth == 0 {
				mult = preference
			}

			if nth >= bgAVB {
				//fmt.Printf("Input flow%d tree%d \n", nth, kth)
				value := mult / float64(schedule.WCD(z, X, Input_flowSet.AVBFlows[nth-bgAVB], flowSet))
				v = append(v, value)

			} else {
				//fmt.Printf("Backgourd flow%d tree%d \n", nth, kth)
				value := mult / float64(schedule.WCD(z, X, BG_flowSet.AVBFlows[nth], flowSet))
				v = append(v, value)
			}
		}
		visibility.AVB_VB = append(visibility.AVB_VB, v)
	}

	return visibility
}

func probability(osaco *OSACO) (*routes.TreesSet, *routes.TreesSet, [2][]int, [2][]int) {
	var (
		input_k_location [2][]int // (tsn k index, avb k index)
		bg_k_location    [2][]int // (tsn k index, avb k index)
	)

	II := &routes.TreesSet{}
	II_prime := &routes.TreesSet{}
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
				// if kth == 5 => arr[0,0,0,0,0,0,...,1,1,1,...,2,2,2,2,..,3,3,3,...,4,4,4,4,...] len(arr) ~ 100
				arr = append(arr, kth)
			}
		}
		randomIndex := rng.IntN(len(arr))
		n = arr[randomIndex]
		t := ktree.Trees[n]

		if nth < bgTSN {
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
				// if kth == 5 => arr[0,0,0,0,0,0,...,1,1,1,...,2,2,2,2,..,3,3,3,...,4,4,4,4,...] len(arr) ~ 100
				arr = append(arr, kth)
			}
		}
		randomIndex := rng.IntN(len(arr))
		n = arr[randomIndex]
		t := ktree.Trees[n]

		if nth < bgAVB {
			bg_k_location[1] = append(bg_k_location[1], n)
			II_prime.AVBTrees = append(II_prime.AVBTrees, t)
		} else {
			input_k_location[1] = append(input_k_location[1], n)
			II.AVBTrees = append(II.AVBTrees, t)
		}
	}

	return II, II_prime, input_k_location, bg_k_location
}

func epoch(network *network.Network, cfg *config.Config, osaco *OSACO, timeoutIndex int, costSetting [4]int) *routes.TreesSet {
	II, _, input_k_location, _ := probability(osaco)
	//II, II_prime, input_k_location, bg_k_location := Probability(osaco.KTrees, osaco.VB, osaco.PRM) // BG ... pass
	fmt.Printf("Select input routing %v \n", input_k_location)
	//fmt.Printf("Select background routing %v \n", bg_k_location) // BG ... pass
	osaco.Timer[timeoutIndex].TimerStop()
	obj_list, cost := schedule.OBJ(network, cfg, osaco.KTrees, II, osaco.BGTrees, costSetting, false)
	//obj, cost := Obj(network, X, II, II_prime) // BG ... pass
	osaco.Timer[timeoutIndex].TimerStart()

	if obj_list[0] == 0 && obj_list[1] == 0 {
		osaco.Timer[timeoutIndex].TimerEnd()
	}

	for nth, ktree := range osaco.KTrees.TSNTrees {
		for kth := range ktree.Trees {
			if nth < bgTSN { // BG ... pass
				//osaco.PRM.TSN_PRM[nth][kth] *= osaco.P
				//if kth == bg_k_location[0][nth] {
				//	osaco.PRM.TSN_PRM[nth][kth] += (1 / cost[3])
				//}
			} else { // Input
				osaco.PRM.TSN_PRM[nth][kth] *= osaco.P
				if kth == input_k_location[0][nth-bgTSN] {
					osaco.PRM.TSN_PRM[nth][kth] += float64(1 / cost)
				}
			}
		}
	}

	for nth, ktree := range osaco.KTrees.AVBTrees {
		for kth := range ktree.Trees {
			if nth < bgAVB { // BG ... pass
				//osaco.PRM.AVB_PRM[nth][kth] *= osaco.P
				//if kth == bg_k_location[1][nth] {
				//	osaco.PRM.AVB_PRM[nth][kth] += (1 / cost[3])
				//}
			} else { // Input
				osaco.PRM.AVB_PRM[nth][kth] *= osaco.P
				if kth == input_k_location[1][nth-bgAVB] {
					osaco.PRM.AVB_PRM[nth][kth] += float64(1 / cost)
				}
			}
		}
	}

	return II
}
