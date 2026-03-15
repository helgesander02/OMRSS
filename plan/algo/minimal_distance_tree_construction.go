package algo

import (
	"src/internal/config"
	"src/network"
	"src/plan/algo_timer"
	"src/plan/routes"
)

func (mtdc *MDTC) MDTC_Run(network *network.Network, cfg *config.Config) {
	mtdc.Timer = algo_timer.NewTimer()
	mtdc.Timer.TimerStart()
	mtdc.Trees = routes.Get_DistanceTree_Routing(network, cfg)
	mtdc.Timer.TimerStop()
}
