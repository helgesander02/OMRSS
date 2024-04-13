package plan

import (
	"src/network"
	"src/plan/algo"
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

// Developing the OMACO plan
func New_OMACO_Plan(network *network.Network, osaco_timeout int, osaco_K int, osaco_P float64) *OMACO {
	OMACO := &OMACO{Network: network}

	OMACO.SMT = &algo.SMT{}
	OMACO.MDTC = &algo.MDTC{}
	OMACO.OSACO = &algo.OSACO{Timeout: osaco_timeout, K: osaco_K, P: osaco_P}

	return OMACO
}

// Plan2
// Plan3
// ...
