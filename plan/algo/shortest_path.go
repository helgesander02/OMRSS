package algo

import (
	"src/network"
	"src/plan/algo_timer"
	"src/plan/routes"
)

// SP_Run computes shortest paths for OSRO network
func (sp *SP) SP_Run(network *network.OSRO_Network) {
	sp.Timer = algo_timer.NewTimer()
	sp.Timer.TimerStart()
	sp.Paths = routes.Get_ShortestPath_Routing(network)
	sp.Timer.TimerStop()
}
