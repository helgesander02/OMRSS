package plan

import (
	"fmt"
	"src/network"
	"src/plan/algo/osaco"
)

func NewPlan(network *network.Network) *Plan {

	return &Plan{Network: network}
}

func (plan *Plan) InitiatePlan(K int, show_plan bool, show_smt bool, show_osaco bool, time_out int) ([3][4]float64, [3][4]float64) {
	var (
		objs_SMT   [3][4]float64 // 100ms{o1, o2, o3, o4} 50ms{o1, o2, o3, o4} 10ms{o1, o2, o3, o4}
		objs_osaco [3][4]float64 // 100ms{o1, o2, o3, o4} 50ms{o1, o2, o3, o4} 10ms{o1, o2, o3, o4}
	)

	// timeout 100ms
	obj_smt_100, obj_osaco_100 := osaco.Run(plan.Network, K, show_smt, show_osaco, time_out)
	objs_SMT[0] = obj_smt_100
	objs_osaco[0] = obj_osaco_100

	// timeout 50ms
	obj_smt_50, obj_osaco_50 := osaco.Run(plan.Network, K, show_smt, show_osaco, 50)
	objs_SMT[1] = obj_smt_50
	objs_osaco[1] = obj_osaco_50

	// timeout 10ms
	obj_smt_10, obj_osaco_10 := osaco.Run(plan.Network, K, show_smt, show_osaco, 10)
	objs_SMT[2] = obj_smt_10
	objs_osaco[2] = obj_osaco_10

	if show_plan {
		fmt.Println("Steiner Tree")
		fmt.Println("----------------------------------------")

		fmt.Println("OSACO")
		fmt.Println("----------------------------------------")

	}

	return objs_SMT, objs_osaco
}
