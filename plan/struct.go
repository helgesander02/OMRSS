package plan

import (
	"src/network"
	"src/plan/algo"
)

type Plan struct {
	Network *network.Network
	SMT     *algo.SMT
	MDTC    *algo.MDTC
	OSACO   *algo.OSACO
}
