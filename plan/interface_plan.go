package plan

import (
	"src/network"
	"src/pkg/config"
)

type Plans interface {
	InitiatePlan([4]int, *config.Config)
	ShowPlan()
}

func NewPlans(nw *network.Network, cfg *config.Config) Plans {
	switch cfg.Algorithm.Name {
	case "omaco":
		return newOMACOPlan(
			nw,
			cfg.Algorithm.OSACO.Timeout,
			cfg.Algorithm.OSACO.KTrees,
			cfg.Algorithm.OSACO.PheromoneEvaporation,
		)

	case "osro":
		return newOSROPlan(
			nw,
			cfg.Algorithm.OSACO.Timeout,
			cfg.Algorithm.OSACO.KTrees,
			cfg.Algorithm.OSACO.PheromoneEvaporation,
		)

	default:
		return nil
	}
}
