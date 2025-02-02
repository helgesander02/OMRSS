package plan

import (
	"src/network"
)

type Plans interface {
	Initiate_Plan([4]int)
	Show_Plan()
}

func New_Plans(networks network.Networks, osaco_timeout int, osaco_K int, osaco_P float64) map[string]Plans {
	// plan1 ...
	OMACO := new_OMACO_Plan(networks.(*network.Network), osaco_timeout, osaco_K, osaco_P)

	// Look-up table method
	Plans := map[string]Plans{
		"omaco": OMACO,
		//plan2,
		//plan3,
		// ...
	}

	return Plans
}
