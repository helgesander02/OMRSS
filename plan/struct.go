package plan

import (
	"src/network"
	"src/network/flow/can"
	"src/pkg/config"
	"src/plan/algo"
)

type OMACO struct {
	Network     *network.Network
	SMT         *algo.SMT
	MDTC        *algo.MDTC
	OSACO       *algo.OSACO
	OSACO_APTED *algo.OSACO
}

// Developing the OMACO plan
func newOMACOPlan(nw *network.Network, cfgOSACO config.OSACOConfig) *OMACO {
	OMACO := &OMACO{Network: nw}

	OMACO.SMT = &algo.SMT{}
	OMACO.MDTC = &algo.MDTC{}
	OMACO.OSACO = &algo.OSACO{Timeout: cfgOSACO.Timeout, K: cfgOSACO.K, P: cfgOSACO.PheromoneEvaporation, Method_Number: 0}
	OMACO.OSACO_APTED = &algo.OSACO{Timeout: cfgOSACO.Timeout, K: cfgOSACO.K, P: cfgOSACO.PheromoneEvaporation, Method_Number: -1}

	return OMACO
}

// OSRO holds one SP / OSACO pair per encap method. Each method is
// evaluated in its own scoped run so per-method (O1, O2, O4) reflect
// only that method's CAN2TT footprint on the timeline. Methods carries
// the iteration order (which matches can.MethodList()).
type OSRO struct {
	Network       *network.Network
	Methods       []string
	SPByMethod    map[string]*algo.SP
	OSACOByMethod map[string]*algo.OSACO
}

func newOSROPlan(nw *network.Network, cfgOSACO config.OSACOConfig) *OSRO {
	OSRO := &OSRO{
		Network:       nw,
		Methods:       can.MethodList(),
		SPByMethod:    make(map[string]*algo.SP),
		OSACOByMethod: make(map[string]*algo.OSACO),
	}

	for _, m := range OSRO.Methods {
		OSRO.SPByMethod[m] = &algo.SP{}
		OSRO.OSACOByMethod[m] = &algo.OSACO{
			Timeout:       cfgOSACO.Timeout,
			K:             cfgOSACO.K,
			P:             cfgOSACO.PheromoneEvaporation,
			Method_Number: 0,
			MethodScope:   m,
		}
	}

	return OSRO
}
