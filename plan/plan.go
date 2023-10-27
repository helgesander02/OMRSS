package plan

import (
	"fmt"
	"src/algo/osaco"
	"src/network"
)

func NewPlan(network *network.Network) *Plan {

	return &Plan{Network: network}
}

func (plan *Plan) InitiatePlan(K int, show_plan bool, show_osaco bool) {
	Input_SMT, BG_SMT, Input_OSACO, BG_OSACO := osaco.Run(plan.Network, K, show_osaco)

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
}
