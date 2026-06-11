package plan

import (
	"src/network"
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

type OSRO struct {
	Network *network.Network
	SP      *algo.SP
	OSACO   *algo.OSACO
}

func newOSROPlan(nw *network.Network, cfgOSACO config.OSACOConfig) *OSRO {
	OSRO := &OSRO{Network: nw}

	OSRO.SP = &algo.SP{}
	OSRO.OSACO = &algo.OSACO{Timeout: cfgOSACO.Timeout, K: cfgOSACO.K, P: cfgOSACO.PheromoneEvaporation, Method_Number: 0}

	return OSRO
}
