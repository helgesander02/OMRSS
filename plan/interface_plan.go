package plan

import (
	"src/network"
	"src/pkg/config"
)

type Plans interface {
	InitiatePlan(*config.Config)
	ShowPlan()
}

func NewPlans(nw *network.Network, cfg *config.Config) Plans {
	switch cfg.Algorithm.Name {
	case "omaco":
		return newOMACOPlan(
			nw,
			cfg.Algorithm.OSACO,
		)

	case "osro":
		return newOSROPlan(
			nw,
			cfg.Algorithm.OSACO,
		)

	default:
		return nil
	}
}
