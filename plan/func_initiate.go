package plan

import (
	"src/pkg/config"
	"src/pkg/logger"
	"src/plan/algo_timer"
	"src/plan/routes"
	"src/plan/schedule"
)

func (plan *OMACO) InitiatePlan(cfg *config.Config) {
	// algo run
	logger.Println("Steiner Tree")
	logger.Println("----------------------------------------")
	plan.SMT.SMT_Run(plan.Network, cfg)

	logger.Println()
	logger.Println("MDTC")
	logger.Println("----------------------------------------")
	plan.MDTC.MDTC_Run(plan.Network, cfg)

	// The timeout of each run is set as 100~1000 ms (200ms, 400ms, 600ms, 800ms, 1000ms)
	timeround := 5

	logger.Println()
	logger.Println("OSACO")
	logger.Println("----------------------------------------")

	osacoTimer := algo_timer.NewTimer()
	osacoTimer.TimerStart()
	osacoKRoutes := routes.Get_OSACO_Routing(plan.Network, cfg, plan.SMT.Routes, plan.OSACO.K, plan.OSACO.Method_Number)
	osacoTimer.TimerEnd()

	plan.OSACO.OSACO_Initial_Settings(plan.Network, cfg, plan.SMT.Routes, osacoKRoutes, timeround, osacoTimer)
	for i := 0; i < timeround; i++ {
		plan.OSACO.Objs_osaco[i] = plan.OSACO.OSACO_Run(plan.Network, cfg, i)
	}

	logger.Println()
	logger.Println("OSACO_APTED")
	logger.Println("----------------------------------------")

	aptedTimer := algo_timer.NewTimer()
	aptedTimer.TimerStart()
	aptedKRoutes := routes.Get_OSACO_Routing(plan.Network, cfg, plan.SMT.Routes, plan.OSACO_APTED.K, plan.OSACO_APTED.Method_Number)
	aptedTimer.TimerEnd()

	plan.OSACO_APTED.OSACO_Initial_Settings(plan.Network, cfg, plan.SMT.Routes, aptedKRoutes, timeround, aptedTimer)
	for i := 0; i < timeround; i++ {
		plan.OSACO_APTED.Objs_osaco[i] = plan.OSACO_APTED.OSACO_Run(plan.Network, cfg, i)
	}

	obj_smt, _ := schedule.OBJ(
		plan.Network,
		cfg,
		plan.OSACO.KRoutes,
		plan.SMT.Routes.InputRouteSet(cfg.Network.Flows.TSN.Background, cfg.Network.Flows.AVB.Background),
		plan.SMT.Routes.BGRouteSet(cfg.Network.Flows.TSN.Background, cfg.Network.Flows.AVB.Background),
		true,
	)

	obj_mdt, _ := schedule.OBJ(
		plan.Network,
		cfg,
		plan.OSACO.KRoutes,
		plan.MDTC.Routes.InputRouteSet(cfg.Network.Flows.TSN.Background, cfg.Network.Flows.AVB.Background),
		plan.MDTC.Routes.BGRouteSet(cfg.Network.Flows.TSN.Background, cfg.Network.Flows.AVB.Background),
		true,
	)

	plan.SMT.Objs_smt = obj_smt
	plan.MDTC.Objs_mdtc = obj_mdt

	if obj_mdt[0] != 0 || obj_mdt[1] != 0 {
		plan.MDTC.Timer.TimerMax()
	}

}

func (plan *OSRO) InitiatePlan(cfg *config.Config) {
	// algo run
	logger.Println("Shortest Path")
	logger.Println("----------------------------------------")
	plan.SP.SP_Run(plan.Network, cfg)

	// The timeout of each run is set as 100~1000 ms (200ms, 400ms, 600ms, 800ms, 1000ms)
	timeround := 5

	logger.Println()
	logger.Println("OSACO (Path-based)")
	logger.Println("----------------------------------------")

	pathTimer := algo_timer.NewTimer()
	pathTimer.TimerStart()
	pathKRoutes := routes.Get_KPath_Routing(plan.Network, cfg, plan.SP.Routes, plan.OSACO.K)
	pathTimer.TimerEnd()

	plan.OSACO.OSACO_Initial_Settings(plan.Network, cfg, plan.SP.Routes, pathKRoutes, timeround, pathTimer)
	for i := 0; i < timeround; i++ {
		plan.OSACO.Objs_osaco[i] = plan.OSACO.OSACO_Run(plan.Network, cfg, i)
	}

	// Shortest-path baseline objective: hand the SP results straight to
	// schedule.OBJ now that Route / RouteSet / KRouteSet are unified.
	spInput := plan.SP.Routes.InputRouteSet(cfg.Network.Flows.TSN.Background, cfg.Network.Flows.AVB.Background)
	spBG := plan.SP.Routes.BGRouteSet(cfg.Network.Flows.TSN.Background, cfg.Network.Flows.AVB.Background)
	plan.SP.Objs_sp, _ = schedule.OBJ(plan.Network, cfg, plan.OSACO.KRoutes, spInput, spBG, true)

	logger.Println("OSRO InitiatePlan completed")
}
