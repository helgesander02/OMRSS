package plan

import (
	"fmt"
	"src/plan/schedule"
)

func (plan *OMACO) Initiate_Plan() {
	// algo run
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
	plan.OSACO.OSACO_Initial_Settings(plan.Network, plan.SMT.Trees)
	// The timeout of each run is set as 100~1000 ms (200ms, 400ms, 600ms, 800ms, 1000ms)
	for i := 0; i < 5; i++ {
		plan.OSACO.Objs_osaco[i] = plan.OSACO.OSACO_Run(plan.Network, i)
	}

	fmt.Println()
	fmt.Println("OSACO_IAS")
	fmt.Println("----------------------------------------")
	plan.OSACO_IAS.OSACO_Initial_Settings(plan.Network, plan.SMT.Trees)
	// The timeout of each run is set as 100~1000 ms (200ms, 400ms, 600ms, 800ms, 1000ms)
	for i := 0; i < 5; i++ {
		plan.OSACO_IAS.Objs_osaco[i] = plan.OSACO_IAS.OSACO_Run(plan.Network, i)
	}

	fmt.Println()
	fmt.Println("OSACO_AAS")
	fmt.Println("----------------------------------------")
	plan.OSACO_AAS.OSACO_Initial_Settings(plan.Network, plan.SMT.Trees)
	// The timeout of each run is set as 100~1000 ms (200ms, 400ms, 600ms, 800ms, 1000ms)
	for i := 0; i < 5; i++ {
		plan.OSACO_AAS.Objs_osaco[i] = plan.OSACO_AAS.OSACO_Run(plan.Network, i)
	}

	obj_smt, _ := schedule.OBJ(
		plan.Network,
		plan.OSACO.KTrees,
		plan.SMT.Trees.Input_Tree_set(plan.Network.BG_TSN, plan.Network.BG_AVB),
		plan.SMT.Trees.BG_Tree_set(plan.Network.BG_TSN, plan.Network.BG_AVB),
	)

	obj_mdt, _ := schedule.OBJ(
		plan.Network,
		plan.OSACO.KTrees,
		plan.MDTC.Trees.Input_Tree_set(plan.Network.BG_TSN, plan.Network.BG_AVB),
		plan.MDTC.Trees.BG_Tree_set(plan.Network.BG_TSN, plan.Network.BG_AVB),
	)

	plan.SMT.Objs_smt = obj_smt
	plan.MDTC.Objs_mdtc = obj_mdt

	if obj_mdt[0] != 0 || obj_mdt[1] != 0 {
		plan.MDTC.Timer.TimerMax()
	}

}

//func (plan *plan2) Initiate_Plan() {
//
//}

//func (plan *plan3) Initiate_Plan() {
//
//}
