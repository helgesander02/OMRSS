package plan

import (
	"src/network"
	"src/plan/algo"
)

type OMACO struct {
	Network   *network.Network
	SMT       *algo.SMT
	MDTC      *algo.MDTC
	OSACO     *algo.OSACO
	OSACO_IAS *algo.OSACO
	OSACO_AAS *algo.OSACO
}

// Developing the OMACO plan
func New_OMACO_Plan(network *network.Network, osaco_timeout int, osaco_K int, osaco_P float64) *OMACO {
	OMACO := &OMACO{Network: network}

	OMACO.SMT = &algo.SMT{}
	OMACO.MDTC = &algo.MDTC{}
	OMACO.OSACO = &algo.OSACO{Timeout: osaco_timeout, K: osaco_K, P: osaco_P, Method_Number: 0}
	OMACO.OSACO_IAS = &algo.OSACO{Timeout: osaco_timeout, K: osaco_K, P: osaco_P, Method_Number: 1}
	OMACO.OSACO_AAS = &algo.OSACO{Timeout: osaco_timeout, K: osaco_K, P: osaco_P, Method_Number: 2}

	return OMACO
}

//type plan2 struct {
//	Network *network.Network
//}

// Plan2

//type plan3 struct {
//	Network *network.Network
//}

// Plan3
// ...
