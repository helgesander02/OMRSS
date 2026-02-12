package plan

import (
	"fmt"
	"src/plan/schedule"
)

func (plan *OMACO) InitiatePlan(costSetting [4]int) {
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
		plan.OSACO.Objs_osaco[i] = plan.OSACO.OSACO_Run(plan.Network, i, costSetting)
	}

	fmt.Println()
	fmt.Println("OSACO_APTED")
	fmt.Println("----------------------------------------")
	plan.OSACO_APTED.OSACO_Initial_Settings(plan.Network, plan.SMT.Trees)
	// The timeout of each run is set as 100~1000 ms (200ms, 400ms, 600ms, 800ms, 1000ms)
	for i := 0; i < 5; i++ {
		plan.OSACO_APTED.Objs_osaco[i] = plan.OSACO_APTED.OSACO_Run(plan.Network, i, costSetting)
	}

	obj_smt, _ := schedule.OBJ(
		plan.Network,
		plan.OSACO.KTrees,
		plan.SMT.Trees.InputTreeSet(plan.Network.BGTSN, plan.Network.BGAVB),
		plan.SMT.Trees.BGTreeSet(plan.Network.BGTSN, plan.Network.BGAVB),
		costSetting,
		true,
	)

	obj_mdt, _ := schedule.OBJ(
		plan.Network,
		plan.OSACO.KTrees,
		plan.MDTC.Trees.InputTreeSet(plan.Network.BGTSN, plan.Network.BGAVB),
		plan.MDTC.Trees.BGTreeSet(plan.Network.BGTSN, plan.Network.BGAVB),
		costSetting,
		true,
	)

	plan.SMT.Objs_smt = obj_smt
	plan.MDTC.Objs_mdtc = obj_mdt

	if obj_mdt[0] != 0 || obj_mdt[1] != 0 {
		plan.MDTC.Timer.TimerMax()
	}

}

func (plan *OSRO) InitiatePlan(costSetting [4]int) {
	// algo run
	fmt.Println("Shortest Path")
	fmt.Println("----------------------------------------")
	plan.SP.SP_Run(plan.Network)

	fmt.Println()
	fmt.Println("OSACO (Path-based)")
	fmt.Println("----------------------------------------")
	plan.OSACO_Path.OSACO_Initial_Settings_Path(plan.Network, plan.SP.Paths)

	// The timeout of each run is set as 100~1000 ms (200ms, 400ms, 600ms, 800ms, 1000ms)
	for i := 0; i < 5; i++ {
		plan.OSACO_Path.Objs_osaco[i] = plan.OSACO_Path.OSACO_Run_Path(plan.Network, i, costSetting)
	}

	// Note: Path-based objective calculation (schedule.OBJ_Path) is not yet implemented
	// This will be added when path-based scheduling is completed
	plan.SP.Objs_sp = [4]float64{0, 0, 0, 0}

	fmt.Println("OSRO InitiatePlan completed")
}
