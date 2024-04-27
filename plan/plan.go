package plan

import (
	"src/network"
)

type Plans interface {
	Initiate_Plan()
	Show_Plan()
}

func New_Plans(network *network.Network, osaco_timeout int, osaco_K int, osaco_P float64) map[string]Plans {
	// plan1 ...
	OMACO := New_OMACO_Plan(network, osaco_timeout, osaco_K, osaco_P)

	// Look-up table method
	Plans := map[string]Plans{
		"omaco": OMACO,
		//plan2,
		//plan3,
		// ...
	}

	return Plans
}
