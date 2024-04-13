package plan

import (
	"src/network"
	"src/plan/algo"
)

type OMACO struct {
	Network *network.Network
	SMT     *algo.SMT
	MDTC    *algo.MDTC
	OSACO   *algo.OSACO
}

//type plan2 struct {
//	Network *network.Network
//}

//type plan3 struct {
//	Network *network.Network
//}
