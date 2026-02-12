package plan

import (
	"src/network"
	"src/plan/algo"
)

type OMACO struct {
	Network     *network.OMACO_Network
	SMT         *algo.SMT
	MDTC        *algo.MDTC
	OSACO       *algo.OSACO
	OSACO_APTED *algo.OSACO
}

// Developing the OMACO plan
func new_OMACO_Plan(network *network.OMACO_Network, osacoTimeout int, osacoK int, osacoP float64) *OMACO {
	OMACO := &OMACO{Network: network}

	OMACO.SMT = &algo.SMT{}
	OMACO.MDTC = &algo.MDTC{}
	OMACO.OSACO = &algo.OSACO{Timeout: osacoTimeout, K: osacoK, P: osacoP, Method_Number: 0}
	OMACO.OSACO_APTED = &algo.OSACO{Timeout: osacoTimeout, K: osacoK, P: osacoP, Method_Number: -1}

	return OMACO
}

type OSRO struct {
	Network    *network.OSRO_Network
	SP         *algo.SP
	OSACO_Path *algo.OSACO_Path
}

func new_OSRO_Plan(network *network.OSRO_Network, osacoTimeout int, osacoK int, osacoP float64) *OSRO {
	OSRO := &OSRO{Network: network}

	OSRO.SP = &algo.SP{}
	OSRO.OSACO_Path = &algo.OSACO_Path{Timeout: osacoTimeout, K: osacoK, P: osacoP}

	return OSRO
}
