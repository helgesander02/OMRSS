package plan

import (
	"fmt"
	"src/network"
	"src/plan/algo"
	"src/plan/routes"
)

func NewPlan(network *network.Network) *Plan {

	return &Plan{Network: network}
}

func (plan *Plan) InitiatePlan(show_plan bool, show_smt bool, show_osaco bool) ([5][4]float64, [5][4]float64) {
	var (
		objs_SMT   [5][4]float64 // 100ms{o1, o2, o3, o4} 50ms{o1, o2, o3, o4} 10ms{o1, o2, o3, o4}
		objs_osaco [5][4]float64 // 100ms{o1, o2, o3, o4} 50ms{o1, o2, o3, o4} 10ms{o1, o2, o3, o4}
		timeout    int           = 500
	)

	// The timeout of each run is set as 100~500 ms
	fmt.Println("Steiner Tree")
	fmt.Println("----------------------------------------")
	SMT := algo.SMT_Run(plan.Network, show_smt)

	fmt.Println()
	fmt.Println("OSACO")
	fmt.Println("----------------------------------------")
	X := routes.Get_OSACO_Routing(plan.Network, SMT)
	for i := 0; i < 5; i++ {
		obj_smt, obj_osaco := algo.OSACO_Run(plan.Network, SMT, X, show_osaco, timeout)
		objs_SMT[i] = obj_smt
		objs_osaco[i] = obj_osaco

		timeout -= 100
	}

	if show_plan {
		fmt.Println("Steiner Tree")
		fmt.Println("----------------------------------------")

		fmt.Println("OSACO")
		fmt.Println("----------------------------------------")

	}

	return objs_SMT, objs_osaco
}
