package plan

import (
	"fmt"
	"src/algo/osaco"
	"src/network"
)

func NewPlan(network *network.Network) *Plan {

	return &Plan{Network: network}
}

func (plan *Plan) InitiatePlan(K int, show_plan bool, show_smt bool, show_osaco bool, time_out int) ([4]float64, [4]float64) {
	Input_SMT, BG_SMT, obj_SMT, Input_OSACO, BG_OSACO, obj_osaco := osaco.Run(plan.Network, K, show_osaco, time_out)

	if show_smt {
		fmt.Println("--- The Steiner Tree final selected routing---")
		Input_SMT.Show_Trees_Set()
		BG_SMT.Show_Trees_Set()
	}
	if show_osaco {
		fmt.Println("--- The OSACO final selected routing ---")
		Input_OSACO.Show_Trees_Set()
		BG_OSACO.Show_Trees_Set()
	}
	if show_plan {
		fmt.Println("Steiner Tree")
		fmt.Println("----------------------------------------")
		Input_SMT.Show_Trees_Set()
		BG_SMT.Show_Trees_Set()

		fmt.Println("OSACO")
		fmt.Println("----------------------------------------")
		Input_OSACO.Show_Trees_Set()
		BG_OSACO.Show_Trees_Set()
	}

	return obj_SMT, obj_osaco
}
