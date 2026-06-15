package algo

import (
	"math"
	"src/network"
	"src/network/flow"
	"src/pkg/config"
	"src/pkg/logger"
	"src/pkg/random"
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

// SetRNG sets the random number generator for this package.
func SetRNG(r *random.Generator) {
	rng = r
}

// OSACO_Initial_Settings populates the OSACO state from a pre-built
// K-route set. Callers compute KRoutes themselves so this entry point is
// agnostic to how the K alternatives were produced (Steiner K-spanning
// trees for OMACO, Yen's K-shortest paths for OSRO).
func (osaco *OSACO) OSACO_Initial_Settings(network *network.Network, cfg *config.Config, routeSet *routes.RouteSet, kroutes *routes.KRouteSet, timeround int, prepTimer *algo_timer.Timer) {
	bgTSN = cfg.Network.Flows.TSN.Background
	bgAVB = cfg.Network.Flows.AVB.Background

	osaco.KRoutes = kroutes
	osaco.InputRoutes = routeSet.InputRouteSet(bgTSN, bgAVB)
	osaco.BGRoutes = routeSet.BGRouteSet(bgTSN, bgAVB)
	osaco.PRM = computePrm(osaco.KRoutes)
	osaco.VB = computeVb(osaco.KRoutes, network.FlowSet)

	for i := 0; i < timeround; i++ {
		osaco.Timer[i] = algo_timer.NewTimer()
		if prepTimer != nil {
			osaco.Timer[i].TimerMerge(prepTimer)
		}
	}
}

// OSACO_Run executes the ant-colony loop until its timeout expires.
// Ching-Chih Chuang et al., "Online Stream-Aware Routing for TSN-Based
// Industrial Control Systems".
func (osaco *OSACO) OSACO_Run(network *network.Network, cfg *config.Config, timeoutIndex int) [4]float64 {
	initialobj, initialcost := schedule.OBJ(network, cfg, osaco.KRoutes, osaco.InputRoutes, osaco.BGRoutes, osaco.MethodScope, false)
	logger.Println()
	logger.Printf("initial value: %d \n", initialcost)
	logger.Printf("O1: %f O2: %f O3: pass O4: %f \n", initialobj[0], initialobj[1], initialobj[3])

	timeout := time.Duration(osaco.Timeout) * time.Millisecond
	startTime := time.Now()
	i := 1
	for {
		logger.Printf("\nepoch%d:\n", i)
		osaco.Timer[timeoutIndex].TimerStart()
		II := epoch(network, cfg, osaco, timeoutIndex)
		osaco.Timer[timeoutIndex].TimerStop()

		_, cost1 := schedule.OBJ(network, cfg, osaco.KRoutes, II, osaco.BGRoutes, osaco.MethodScope, false)
		_, cost2 := schedule.OBJ(network, cfg, osaco.KRoutes, osaco.InputRoutes, osaco.BGRoutes, osaco.MethodScope, false)

		if cost1 < cost2 {
			osaco.InputRoutes = II
			logger.Println("Change the selected routing !!")
		}
		i += 1

		if time.Since(startTime) >= timeout {
			break
		}
	}

	resultobj, resultcost := schedule.OBJ(network, cfg, osaco.KRoutes, osaco.InputRoutes, osaco.BGRoutes, osaco.MethodScope, true)
	logger.Println()
	logger.Printf("result value: %d \n", resultcost)
	logger.Printf("O1: %f O2: %f O3: pass O4: %f \n", resultobj[0], resultobj[1], resultobj[3])
	logger.Println()

	if resultobj[0] != 0 || resultobj[1] != 0 {
		osaco.Timer[timeoutIndex].TimerMax()
	}

	return resultobj
}

// computePrm seeds the pheromone trail. Background streams start with 0.5
// (less attractive), input streams with 1.0. CAN2TT streams always start
// at 1.0 because they have no background partition.
func computePrm(X *routes.KRouteSet) *Pheromone {
	pheromone := &Pheromone{}

	for nth, kr := range X.TSNRoutes {
		var prm []float64
		for i := 0; i < len(kr.Routes); i++ {
			if nth < bgTSN {
				prm = append(prm, 0.5)
			} else {
				prm = append(prm, 1.)
			}
		}
		pheromone.TSN_PRM = append(pheromone.TSN_PRM, prm)
	}

	for nth, kr := range X.AVBRoutes {
		var prm []float64
		for i := 0; i < len(kr.Routes); i++ {
			if nth < bgAVB {
				prm = append(prm, 0.5)
			} else {
				prm = append(prm, 1.)
			}
		}
		pheromone.AVB_PRM = append(pheromone.AVB_PRM, prm)
	}

	for _, kr := range X.CAN2TTRoutes {
		var prm []float64
		for i := 0; i < len(kr.Routes); i++ {
			prm = append(prm, 1.)
		}
		pheromone.CAN2TTPRM = append(pheromone.CAN2TTPRM, prm)
	}

	return pheromone
}

// computeVb seeds the visibility heuristic per stream class:
//
//   - TSN  / CAN2TT : 1 / exp(hop weight) — cheap and consistent
//   - AVB           : 1 / WCD(route)      — uses worst-case delay so
//     trees / paths with the same hop count
//     but different interference profiles are
//     told apart. APTED relies on this to
//     surface tree-shape diversity inside one
//     K-bundle.
//
// Input streams' first alternative gets a preference bonus to encourage
// exploitation of the shortest route. Background streams stay at 1.0.
func computeVb(X *routes.KRouteSet, flowSet *flow.FlowSet) *Visibility {
	const preference = 2.0
	inputFlows := flowSet.InputOMACOFlowSet()
	bgFlows := flowSet.BGOMACOFlowSet()

	visibility := &Visibility{}

	// TSN visibility = μ / exp(W₀) where W₀ is the shortest route's hop
	// count (kr.Routes[0].Weight). All K alternatives share the same
	// scaling; only the preference bonus at kth == 0 differentiates them.
	for nth, kr := range X.TSNRoutes {
		var v []float64
		if len(kr.Routes) == 0 {
			visibility.TSN_VB = append(visibility.TSN_VB, v)
			continue
		}
		baseWeight := kr.Routes[0].Weight
		for kth := range kr.Routes {
			mult := 1.0
			if nth < bgTSN && kth == 0 {
				mult = preference
			}
			v = append(v, mult/math.Exp(float64(baseWeight)))
		}
		visibility.TSN_VB = append(visibility.TSN_VB, v)
	}

	for nth, kr := range X.AVBRoutes {
		var v []float64
		for kth, r := range kr.Routes {
			mult := 1.0
			if nth < bgAVB && kth == 0 {
				mult = preference
			}
			var wcd float64
			if nth < bgAVB {
				wcd = float64(schedule.WCD(r, X, bgFlows.AVBFlows[nth], flowSet))
			} else {
				wcd = float64(schedule.WCD(r, X, inputFlows.AVBFlows[nth-bgAVB], flowSet))
			}
			if wcd <= 0 {
				wcd = 1 // guard against divide-by-zero if WCD ever returns 0
			}
			v = append(v, mult/wcd)
		}
		visibility.AVB_VB = append(visibility.AVB_VB, v)
	}

	for _, kr := range X.CAN2TTRoutes {
		var v []float64
		for _, r := range kr.Routes {
			//v = append(v, 1.0/float64(r.Weight))
			v = append(v, 1.0/math.Exp(float64(r.Weight)))
		}
		visibility.CAN2TTVB = append(visibility.CAN2TTVB, v)
	}

	return visibility
}

// probability samples one route per stream weighted by pheromone *
// visibility. The returned input_k_location and bg_k_location are sliced
// per-class — [0]=TSN, [1]=AVB, [2]=CAN2TT — so callers know which
// alternative the ants picked.
func probability(osaco *OSACO) (*routes.RouteSet, *routes.RouteSet, [3][]int, [3][]int) {
	var (
		inputLoc [3][]int
		bgLoc    [3][]int
	)
	II := &routes.RouteSet{}
	IIPrime := &routes.RouteSet{}

	pick := func(weights []float64, count int) int {
		denom := 0.0
		for _, w := range weights {
			denom += w
		}
		var arr []int
		for kth := 0; kth < count; kth++ {
			p := weights[kth] / denom
			for j := 0; j < int(p*100); j++ {
				arr = append(arr, kth)
			}
		}
		return arr[rng.IntN(len(arr))]
	}

	// TSN
	for nth, kr := range osaco.KRoutes.TSNRoutes {
		weights := make([]float64, len(kr.Routes))
		for kth := range kr.Routes {
			weights[kth] = osaco.VB.TSN_VB[nth][kth] * osaco.PRM.TSN_PRM[nth][kth]
		}
		n := pick(weights, len(kr.Routes))
		r := kr.Routes[n]
		if nth < bgTSN {
			bgLoc[0] = append(bgLoc[0], n)
			IIPrime.TSNRoutes = append(IIPrime.TSNRoutes, r)
		} else {
			inputLoc[0] = append(inputLoc[0], n)
			II.TSNRoutes = append(II.TSNRoutes, r)
		}
	}

	// AVB
	for nth, kr := range osaco.KRoutes.AVBRoutes {
		weights := make([]float64, len(kr.Routes))
		for kth := range kr.Routes {
			weights[kth] = osaco.VB.AVB_VB[nth][kth] * osaco.PRM.AVB_PRM[nth][kth]
		}
		n := pick(weights, len(kr.Routes))
		r := kr.Routes[n]
		if nth < bgAVB {
			bgLoc[1] = append(bgLoc[1], n)
			IIPrime.AVBRoutes = append(IIPrime.AVBRoutes, r)
		} else {
			inputLoc[1] = append(inputLoc[1], n)
			II.AVBRoutes = append(II.AVBRoutes, r)
		}
	}

	// CAN2TT (always input-side)
	for nth, kr := range osaco.KRoutes.CAN2TTRoutes {
		weights := make([]float64, len(kr.Routes))
		for kth := range kr.Routes {
			weights[kth] = osaco.VB.CAN2TTVB[nth][kth] * osaco.PRM.CAN2TTPRM[nth][kth]
		}
		n := pick(weights, len(kr.Routes))
		inputLoc[2] = append(inputLoc[2], n)
		II.CAN2TTRoutes = append(II.CAN2TTRoutes, kr.Routes[n])
	}

	return II, IIPrime, inputLoc, bgLoc
}

// epoch runs one ant pass: pick a routing, score it via OBJ, then
// evaporate / reinforce the pheromone trail.
func epoch(network *network.Network, cfg *config.Config, osaco *OSACO, timeoutIndex int) *routes.RouteSet {
	II, _, inputLoc, _ := probability(osaco)
	logger.Printf("Select input routing %v \n", inputLoc)
	osaco.Timer[timeoutIndex].TimerStop()

	obj, cost := schedule.OBJ(network, cfg, osaco.KRoutes, II, osaco.BGRoutes, osaco.MethodScope, false)

	osaco.Timer[timeoutIndex].TimerStart()
	if obj[0] == 0 && obj[1] == 0 {
		osaco.Timer[timeoutIndex].TimerEnd()
	}

	reinforce := 1.0 / float64(cost+1)

	for nth, kr := range osaco.KRoutes.TSNRoutes {
		for kth := range kr.Routes {
			if nth < bgTSN {
				continue
			}
			osaco.PRM.TSN_PRM[nth][kth] *= osaco.P
			if kth == inputLoc[0][nth-bgTSN] {
				osaco.PRM.TSN_PRM[nth][kth] += reinforce
			}
		}
	}

	for nth, kr := range osaco.KRoutes.AVBRoutes {
		for kth := range kr.Routes {
			if nth < bgAVB {
				continue
			}
			osaco.PRM.AVB_PRM[nth][kth] *= osaco.P
			if kth == inputLoc[1][nth-bgAVB] {
				osaco.PRM.AVB_PRM[nth][kth] += reinforce
			}
		}
	}

	for nth, kr := range osaco.KRoutes.CAN2TTRoutes {
		for kth := range kr.Routes {
			osaco.PRM.CAN2TTPRM[nth][kth] *= osaco.P
			if kth == inputLoc[2][nth] {
				osaco.PRM.CAN2TTPRM[nth][kth] += reinforce
			}
		}
	}

	return II
}
