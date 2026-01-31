package plan

import (
	"src/network"
)

type Plans interface {
	InitiatePlan([4]int)
	ShowPlan()
}

func NewPlans(planName string, networks network.Networks, osacoTimeout int, osacoK int, osacoP float64) Plans {
	switch planName {
	case "omaco":
		return new_OMACO_Plan(networks.(*network.OMACO_Network), osacoTimeout, osacoK, osacoP)

	case "osro":
		return new_OSRO_Plan(networks.(*network.OSRO_Network), osacoTimeout, osacoK, osacoP)

	default:
		return nil
	}
}
