package plan

import (
	"fmt"
	"src/plan/routes"
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

	// The timeout of each run is set as 100~1000 ms
	fmt.Println()
	fmt.Println("OSACO")
	fmt.Println("----------------------------------------")
	plan.OSACO.KTrees = routes.Get_OSACO_Routing(plan.Network, plan.SMT.Trees, plan.OSACO.K)
	plan.OSACO.InputTrees = plan.SMT.Trees.Input_Tree_set()
	plan.OSACO.BGTrees = plan.SMT.Trees.BG_Tree_set()
	for i := 0; i < 5; i++ {
		plan.OSACO.Objs_osaco[i] = plan.OSACO.OSACO_Run(plan.Network, i)
	}
	obj_smt, _ := schedule.OBJ(plan.Network, plan.OSACO.KTrees, plan.SMT.Trees.Input_Tree_set(), plan.SMT.Trees.BG_Tree_set())
	plan.SMT.Objs_smt = obj_smt
	obj_mdt, _ := schedule.OBJ(plan.Network, plan.OSACO.KTrees, plan.MDTC.Trees.Input_Tree_set(), plan.MDTC.Trees.BG_Tree_set())
	plan.MDTC.Objs_mdtc = obj_mdt
}

//func (plan *plan2) Initiate_Plan() {
//
//}

//func (plan *plan3) Initiate_Plan() {
//
//}
