package algo

import (
	"src/network"
	"src/plan/algo_timer"
	"src/plan/routes"
)

func (smt *SMT) SMT_Run(network *network.Network) {
	// 4. SteinerTree
	smt.Timer = algo_timer.NewTimer()
	//// SMT computing time: Estimate the time it takes to compute routing information
	smt.Timer.TimerStart()
	smt.Trees = routes.Get_SteninerTree_Routing(network)
	smt.Timer.TimerStop()
}
