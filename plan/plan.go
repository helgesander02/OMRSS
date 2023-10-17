package plan

import (
	"src/algo/osaco"
	"src/network"
)

func NewPlan(network *network.Network) *Plan {

	return &Plan{Network: network}
}

func (plan *Plan) InitiatePlan(show_routes bool, K int, show_osaco bool) {
	osaco.Run(plan.Network, K, show_routes, show_osaco)
}
