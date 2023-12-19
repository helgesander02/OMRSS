package plan

import (
	"fmt"
	"src/network"
	"src/plan/algo"
	"src/plan/routes"
	"src/plan/schedule"
)

func NewPlan(network *network.Network) *Plan {

	return &Plan{Network: network}
}

func (plan *Plan) InitiatePlan(show_plan bool, show_smt bool, show_mdt bool, show_osaco bool) ([4]float64, [4]float64, [5][4]float64) {
	var (
		objs_osaco [5][4]float64 // 1000ms{o1, o2, o3, o4} 800ms{o1, o2, o3, o4} 600ms{o1, o2, o3, o4}, 400ms{o1, o2, o3, o4}, 200ms{o1, o2, o3, o4}
		timeout    int           = 1000
	)

	// The timeout of each run is set as 100~500 ms
	fmt.Println("Steiner Tree")
	fmt.Println("----------------------------------------")
	SMT := algo.SMT_Run(plan.Network, show_smt)

	fmt.Println()
	fmt.Println("MDTC")
	fmt.Println("----------------------------------------")
	MDT := algo.MDT_Run(plan.Network, show_mdt)

	fmt.Println()
	fmt.Println("OSACO")
	fmt.Println("----------------------------------------")
	K_MST := routes.Get_OSACO_Routing(plan.Network, SMT)

	for i := 0; i < 5; i++ {
		obj_osaco := algo.OSACO_Run(plan.Network, SMT, K_MST, timeout, show_osaco)
		objs_osaco[i] = obj_osaco

		timeout -= 200
	}
	obj_smt, _ := schedule.OBJ(plan.Network, K_MST, SMT.Input_Tree_set(), SMT.BG_Tree_set())
	obj_mdt, _ := schedule.OBJ(plan.Network, K_MST, MDT.Input_Tree_set(), MDT.BG_Tree_set())

	if show_plan {
		fmt.Println()
		fmt.Println("Steiner Tree")
		fmt.Println("----------------------------------------")
		SMT.Show_Trees_Set()
		fmt.Println()
		fmt.Println("OSACO")
		fmt.Println("----------------------------------------")
		K_MST.Show_kTrees_Set()
		fmt.Println()
		fmt.Println("MDTC")
		fmt.Println("----------------------------------------")
		MDT.Show_Trees_Set()
	}

	return obj_smt, obj_mdt, objs_osaco
}
