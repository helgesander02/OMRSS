package plan

import (
	"src/network"
)

type Plans interface {
	Initiate_Plan([4]int)
	Show_Plan()
}

func New_Plans(plan_name string, networks network.Networks, osaco_timeout int, osaco_K int, osaco_P float64) Plans {
	switch plan_name {
	case "omaco":
		return new_OMACO_Plan(networks.(*network.OMACO_Network), osaco_timeout, osaco_K, osaco_P)

	case "osro":
		return new_OSRO_Plan(networks.(*network.OSRO_Network), osaco_timeout, osaco_K, osaco_P)

	default:
		return nil
	}
}
