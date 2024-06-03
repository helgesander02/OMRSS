package algo

import (
	"src/network"
	"src/plan/algo_timer"
	"src/plan/routes"
)

func (mtdc *MDTC) MDTC_Run(network *network.Network) {
	// 5. DistanceTree
	mtdc.Timer = algo_timer.NewTimer()
	mtdc.Timer.TimerStart()
	mtdc.Trees = routes.Get_DistanceTree_Routing(network)
	mtdc.Timer.TimerStop()
}
