package algo

import (
	"math"
	"src/network"
	"src/network/flow"
	"src/pkg/config"
	"src/pkg/logger"
	"src/plan/algo_timer"
	"src/plan/routes"
	"time"
)

var (
	bgTSN_path int
	bgAVB_path int
)

// OSACO_Initial_Settings_Path initializes OSACO for path-based routing (OSRO)
func (osaco *OSACO_Path) OSACO_Initial_Settings_Path(network *network.Network, cfg *config.Config, SP *routes.PathsSet) {
	bgTSN_path = cfg.Network.Flows.TSN.Background
	bgAVB_path = cfg.Network.Flows.AVB.Background

	timer := algo_timer.NewTimer()
	timer.TimerStart()
	osaco.KPaths = routes.Get_KPath_Routing(network, cfg, SP, osaco.K)
	timer.TimerEnd()

	osaco.InputPaths = SP.InputPathSet(bgTSN_path, bgAVB_path)
	osaco.BGPaths = SP.BGPathSet(bgTSN_path, bgAVB_path)
	osaco.PRM_Path = computePrmPath(osaco.KPaths)
	osaco.VB_Path = computeVbPath(osaco.KPaths, network.FlowSet)

	for i := 0; i < 5; i++ {
		osaco.Timer[i] = algo_timer.NewTimer()
		osaco.Timer[i].TimerMerge(timer)
	}
}

// OSACO_Run_Path runs OSACO algorithm for path-based routing
func (osaco *OSACO_Path) OSACO_Run_Path(network *network.Network, timeoutIndex int, costSetting [4]int) [4]float64 {
	// Note: Path-based scheduling objective function will be implemented
	// when path-based WCD (Worst-Case Delay) calculation is available
	initialobj := [4]float64{0, 0, 0, 0}
	initialcost := 0

	logger.Println()
	logger.Printf("initial value: %d \n", initialcost)
	logger.Printf("O1: %f O2: %f O3: pass O4: %f \n", initialobj[0], initialobj[1], initialobj[3])

	timeout := time.Duration(osaco.Timeout) * time.Millisecond
	startTime := time.Now()
	i := 1

	for {
		logger.Printf("\nepoch%d:\n", i)
		osaco.Timer[timeoutIndex].TimerStart()
		II := epochPath(network, osaco, timeoutIndex, costSetting)
		osaco.Timer[timeoutIndex].TimerStop()

		// Note: Cost calculation will use path-based OBJ function when available
		cost1 := 0 // Placeholder for new path cost
		cost2 := 0 // Placeholder for current path cost

		if cost1 < cost2 {
			osaco.InputPaths = II
			logger.Println("Change the selected routing !!")
		}
		i += 1

		if time.Since(startTime) >= timeout {
			break
		}
	}

	resultobj := [4]float64{0, 0, 0, 0}
	resultcost := 0

	logger.Println()
	logger.Printf("result value: %d \n", resultcost)
	logger.Printf("O1: %f O2: %f O3: pass O4: %f \n", resultobj[0], resultobj[1], resultobj[3])
	logger.Println()

	if resultobj[0] != 0 || resultobj[1] != 0 {
		osaco.Timer[timeoutIndex].TimerMax()
	}

	return resultobj
}

// computePrmPath initializes pheromone values for paths
func computePrmPath(X *routes.KPathsSet) *PheromonePath {
	pheromone := &PheromonePath{}

	for nth, kpath := range X.TSNPaths {
		var prm []float64
		for i := 0; i < len(kpath.Paths); i++ {
			if nth < bgTSN_path {
				prm = append(prm, 0.5)
			} else {
				prm = append(prm, 1.)
			}
		}
		pheromone.TSN_PRM = append(pheromone.TSN_PRM, prm)
	}

	for nth, kpath := range X.AVBPaths {
		var prm []float64
		for i := 0; i < len(kpath.Paths); i++ {
			if nth < bgAVB_path {
				prm = append(prm, 0.5)
			} else {
				prm = append(prm, 1.)
			}
		}
		pheromone.AVB_PRM = append(pheromone.AVB_PRM, prm)
	}

	for _, kpath := range X.CAN2TSNPaths {
		var prm []float64
		for i := 0; i < len(kpath.Paths); i++ {
			prm = append(prm, 1.)
		}
		pheromone.CAN2TSN_PRM = append(pheromone.CAN2TSN_PRM, prm)
	}

	return pheromone
}

// computeVbPath computes visibility values for paths
func computeVbPath(X *routes.KPathsSet, flowSet *flow.FlowSet) *VisibilityPath {
	var preference float64 = 2.

	visibility := &VisibilityPath{}

	// TSN flow visibility
	for nth, kpath := range X.TSNPaths {
		var v []float64
		for kth := range kpath.Paths {
			mult := 1.
			if nth < bgTSN_path && kth == 0 {
				mult = preference
			}

			// Use path weight (hop count) for visibility
			value := mult / math.Exp(float64(kpath.Paths[kth].Weight))
			v = append(v, value)
		}
		visibility.TSN_VB = append(visibility.TSN_VB, v)
	}

	// AVB flow visibility
	for nth, kpath := range X.AVBPaths {
		var v []float64
		for kth := range kpath.Paths {
			mult := 1.
			if nth < bgAVB_path && kth == 0 {
				mult = preference
			}

			// Use path weight for visibility
			value := mult / math.Exp(float64(kpath.Paths[kth].Weight))
			v = append(v, value)
		}
		visibility.AVB_VB = append(visibility.AVB_VB, v)
	}

	// CAN2TSN flow visibility
	for _, kpath := range X.CAN2TSNPaths {
		var v []float64
		for _, path := range kpath.Paths {
			mult := 1.
			value := mult / math.Exp(float64(path.Weight))
			v = append(v, value)
		}
		visibility.CAN2TSN_VB = append(visibility.CAN2TSN_VB, v)
	}

	return visibility
}

// probabilityPath selects paths based on pheromone and visibility
func probabilityPath(osaco *OSACO_Path) (*routes.PathsSet, *routes.PathsSet, [3][]int, [3][]int) {
	var (
		input_k_location [3][]int // (tsn k index, avb k index, can2tsn k index)
		bg_k_location    [3][]int // (tsn k index, avb k index, can2tsn k index)
	)

	II := &routes.PathsSet{}
	II_prime := &routes.PathsSet{}

	// TSN paths
	for nth, kpath := range osaco.KPaths.TSNPaths {
		Denominator := 0.
		for kth := range kpath.Paths {
			Denominator += osaco.VB_Path.TSN_VB[nth][kth] * osaco.PRM_Path.TSN_PRM[nth][kth]
		}

		n := 0
		var arr []int
		for kth := range kpath.Paths {
			probability := (osaco.VB_Path.TSN_VB[nth][kth] * osaco.PRM_Path.TSN_PRM[nth][kth]) / Denominator
			for j := 0; j < int(probability*100); j++ {
				arr = append(arr, kth)
			}
		}
		randomIndex := rng.IntN(len(arr))
		n = arr[randomIndex]
		p := kpath.Paths[n]

		if nth < bgTSN_path {
			bg_k_location[0] = append(bg_k_location[0], n)
			II_prime.TSNPaths = append(II_prime.TSNPaths, p)
		} else {
			input_k_location[0] = append(input_k_location[0], n)
			II.TSNPaths = append(II.TSNPaths, p)
		}
	}

	// AVB paths
	for nth, kpath := range osaco.KPaths.AVBPaths {
		Denominator := 0.
		for kth := range kpath.Paths {
			Denominator += osaco.VB_Path.AVB_VB[nth][kth] * osaco.PRM_Path.AVB_PRM[nth][kth]
		}

		n := 0
		var arr []int
		for kth := range kpath.Paths {
			probability := (osaco.VB_Path.AVB_VB[nth][kth] * osaco.PRM_Path.AVB_PRM[nth][kth]) / Denominator
			for j := 0; j < int(probability*100); j++ {
				arr = append(arr, kth)
			}
		}
		randomIndex := rng.IntN(len(arr))
		n = arr[randomIndex]
		p := kpath.Paths[n]

		if nth < bgAVB_path {
			bg_k_location[1] = append(bg_k_location[1], n)
			II_prime.AVBPaths = append(II_prime.AVBPaths, p)
		} else {
			input_k_location[1] = append(input_k_location[1], n)
			II.AVBPaths = append(II.AVBPaths, p)
		}
	}

	// CAN2TSN paths
	for nth, kpath := range osaco.KPaths.CAN2TSNPaths {
		Denominator := 0.
		for kth := range kpath.Paths {
			Denominator += osaco.VB_Path.CAN2TSN_VB[nth][kth] * osaco.PRM_Path.CAN2TSN_PRM[nth][kth]
		}

		n := 0
		var arr []int
		for kth := range kpath.Paths {
			probability := (osaco.VB_Path.CAN2TSN_VB[nth][kth] * osaco.PRM_Path.CAN2TSN_PRM[nth][kth]) / Denominator
			for j := 0; j < int(probability*100); j++ {
				arr = append(arr, kth)
			}
		}
		randomIndex := rng.IntN(len(arr))
		n = arr[randomIndex]
		p := kpath.Paths[n]

		input_k_location[2] = append(input_k_location[2], n)
		II.CAN2TSNPaths = append(II.CAN2TSNPaths, p)
	}

	return II, II_prime, input_k_location, bg_k_location
}

// epochPath performs one epoch of OSACO for paths
func epochPath(network *network.Network, osaco *OSACO_Path, timeoutIndex int, costSetting [4]int) *routes.PathsSet {
	II, _, input_k_location, _ := probabilityPath(osaco)

	logger.Printf("Select input routing %v \n", input_k_location)

	osaco.Timer[timeoutIndex].TimerStop()

	// Note: Objective calculation will be implemented with path-based scheduling
	obj_list := [4]float64{0, 0, 0, 0}
	cost := 0

	osaco.Timer[timeoutIndex].TimerStart()

	if obj_list[0] == 0 && obj_list[1] == 0 {
		osaco.Timer[timeoutIndex].TimerEnd()
	}

	// Update pheromones for TSN
	for nth, kpath := range osaco.KPaths.TSNPaths {
		for kth := range kpath.Paths {
			if nth < bgTSN_path {
				// BG ... pass
			} else {
				osaco.PRM_Path.TSN_PRM[nth][kth] *= osaco.P
				if kth == input_k_location[0][nth-bgTSN_path] {
					osaco.PRM_Path.TSN_PRM[nth][kth] += float64(1 / (cost + 1))
				}
			}
		}
	}

	// Update pheromones for AVB
	for nth, kpath := range osaco.KPaths.AVBPaths {
		for kth := range kpath.Paths {
			if nth < bgAVB_path {
				// BG ... pass
			} else {
				osaco.PRM_Path.AVB_PRM[nth][kth] *= osaco.P
				if kth == input_k_location[1][nth-bgAVB_path] {
					osaco.PRM_Path.AVB_PRM[nth][kth] += float64(1 / (cost + 1))
				}
			}
		}
	}

	// Update pheromones for CAN2TSN
	for nth, kpath := range osaco.KPaths.CAN2TSNPaths {
		for kth := range kpath.Paths {
			osaco.PRM_Path.CAN2TSN_PRM[nth][kth] *= osaco.P
			if kth == input_k_location[2][nth] {
				osaco.PRM_Path.CAN2TSN_PRM[nth][kth] += float64(1 / (cost + 1))
			}
		}
	}

	return II
}
