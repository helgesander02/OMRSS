package plan

import (
	"fmt"
	"src/plan/schedule"
)

func (plan *OMACO) Initiate_Plan() {
	fmt.Println("Steiner Tree")
	fmt.Println("----------------------------------------")
	plan.SMT.SMT_Run(plan.Network)

	fmt.Println()
	fmt.Println("MDTC")
	fmt.Println("----------------------------------------")
	plan.MDTC.MDTC_Run(plan.Network)

	fmt.Println()
	fmt.Println("OSACO")
	fmt.Println("----------------------------------------")
	// The timeout of each run is set as 100~1000 ms (200ms, 400ms, 600ms, 800ms, 1000ms)
	for i := 0; i < 5; i++ {
		plan.OSACO.Objs_osaco[i] = plan.OSACO.OSACO_Run(plan.Network, plan.SMT.Trees, i)
	}

	obj_smt, _, smt_timer_2 := schedule.OBJ(plan.Network, plan.OSACO.KTrees, plan.SMT.Trees.Input_Tree_set(), plan.SMT.Trees.BG_Tree_set())
	plan.SMT.Objs_smt = obj_smt
	plan.SMT.Timer.TimerMerge(smt_timer_2)

	obj_mdt, _, mdtc_timer_2 := schedule.OBJ(plan.Network, plan.OSACO.KTrees, plan.MDTC.Trees.Input_Tree_set(), plan.MDTC.Trees.BG_Tree_set())
	plan.MDTC.Objs_mdtc = obj_mdt
	plan.MDTC.Timer.TimerMerge(mdtc_timer_2)
}

//func (plan *plan2) Initiate_Plan() {
//
//}

//func (plan *plan3) Initiate_Plan() {
//
//}
