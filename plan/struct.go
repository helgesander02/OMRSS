package plan

import (
	"src/network"
	"src/plan/algo"
)

type OMACO struct {
	Network   *network.OMACO_Network
	SMT       *algo.SMT
	MDTC      *algo.MDTC
	OSACO     *algo.OSACO
	OSACO_IAS *algo.OSACO
}

// Developing the OMACO plan
func new_OMACO_Plan(network *network.OMACO_Network, osaco_timeout int, osaco_K int, osaco_P float64) *OMACO {
	OMACO := &OMACO{Network: network}

	OMACO.SMT = &algo.SMT{}
	OMACO.MDTC = &algo.MDTC{}
	OMACO.OSACO = &algo.OSACO{Timeout: osaco_timeout, K: osaco_K, P: osaco_P, Method_Number: 0}
	OMACO.OSACO_IAS = &algo.OSACO{Timeout: osaco_timeout, K: osaco_K, P: osaco_P, Method_Number: -1}

	return OMACO
}

type OSRO struct {
	Network *network.OSRO_Network
}

func new_OSRO_Plan(network *network.OSRO_Network, osaco_timeout int, osaco_K int, osaco_P float64) *OSRO {
	OSRO := &OSRO{Network: network}

	return OSRO
}

//type plan3 struct {
//	Network *network.Network
//}

// Plan3
// ...
