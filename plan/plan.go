package plan

import (
	"fmt"
	"src/network"
	"src/plan/algo"
	"src/plan/routes"
	"src/plan/schedule"
)

func NewPlan(network *network.Network, osaco_timeout int, osaco_K int, osaco_P float64) *Plan {
	Plan := &Plan{Network: network}

	Plan.SMT = &algo.SMT{}
	Plan.MDTC = &algo.MDTC{}
	Plan.OSACO = &algo.OSACO{Timeout: osaco_timeout, K: osaco_K, P: osaco_P}

	return Plan
}

func (plan *Plan) InitiatePlan() {
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
