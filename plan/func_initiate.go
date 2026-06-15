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
	osacoKRoutes := routes.Get_KTree_Routing(plan.Network, cfg, plan.SMT.Routes, plan.OSACO.K, plan.OSACO.Method_Number)
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
	aptedKRoutes := routes.Get_KTree_Routing(plan.Network, cfg, plan.SMT.Routes, plan.OSACO_APTED.K, plan.OSACO_APTED.Method_Number)
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
		"", // OMACO has no CAN2TT
		true,
	)

	obj_mdt, _ := schedule.OBJ(
		plan.Network,
		cfg,
		plan.OSACO.KRoutes,
		plan.MDTC.Routes.InputRouteSet(cfg.Network.Flows.TSN.Background, cfg.Network.Flows.AVB.Background),
		plan.MDTC.Routes.BGRouteSet(cfg.Network.Flows.TSN.Background, cfg.Network.Flows.AVB.Background),
		"", // OMACO has no CAN2TT
		true,
	)

	plan.SMT.Objs_smt = obj_smt
	plan.MDTC.Objs_mdtc = obj_mdt

	if obj_mdt[0] != 0 || obj_mdt[1] != 0 {
		plan.MDTC.Timer.TimerMax()
	}

}

// OSRO.InitiatePlan runs the full SP + OSACO pipeline once per encap
// method. The shared (and expensive) work — shortest-path computation,
// K-path routing — is done once for the entire flow set, then the
// per-method scoping happens at the schedule layer: each method's
// OSACO sees a KRouteSet trimmed to only its own CAN2TT routes via
// FilterCAN2TTByMethod, and every OBJ scoring call carries its method
// name in MethodScope.
//
// This is what makes OSRO's per-method comparison meaningful: when
// methodA's OSACO scores an (II, II') candidate, the AVB WCDs and
// CAN2TT failure count it sees come *only* from methodA's gateway
// emission — not from a soup of all six methods' frames stacked on
// the same timeline.
func (plan *OSRO) InitiatePlan(cfg *config.Config) {
	// SP and K-path computation are method-agnostic so we share them
	// across every per-method run. The SP run also primes plan.SP for
	// the per-method SP objective evaluation below.
	logger.Println("Shortest Path")
	logger.Println("----------------------------------------")
	sharedSP := plan.SPByMethod[plan.Methods[0]]
	sharedSP.SP_Run(plan.Network, cfg)

	pathTimer := algo_timer.NewTimer()
	pathTimer.TimerStart()
	pathKRoutes := routes.Get_KPath_Routing(plan.Network, cfg, sharedSP.Routes, plan.OSACOByMethod[plan.Methods[0]].K)
	pathTimer.TimerEnd()

	// 5 timeouts × 1 OSACO_Run per method, mirroring the original
	// timeround=5 loop semantics.
	timeround := 5

	for _, method := range plan.Methods {
		logger.Println()
		logger.Printf("OSACO (Path-based) — method: %s\n", method)
		logger.Println("----------------------------------------")

		// Per-method scope: trim KRouteSet to this method's CAN2TT
		// candidate trees, then OSACO's probability walk and timeline
		// build will only see this method's CAN2TT.
		methodKRoutes := pathKRoutes.FilterCAN2TTByMethod(method)
		methodRoutes := sharedSP.Routes.FilterCAN2TTByMethod(method)

		osacoForMethod := plan.OSACOByMethod[method]
		osacoForMethod.OSACO_Initial_Settings(plan.Network, cfg, methodRoutes, methodKRoutes, timeround, pathTimer)
		for i := 0; i < timeround; i++ {
			osacoForMethod.Objs_osaco[i] = osacoForMethod.OSACO_Run(plan.Network, cfg, i)
		}

		// Per-method SP baseline objective — scope the same way OSACO
		// does so the SP and OSACO numbers are directly comparable.
		spForMethod := plan.SPByMethod[method]
		if spForMethod != sharedSP {
			// Cheap copy: share the routing result, give the method its
			// own Timer for downstream accounting.
			spForMethod.Routes = sharedSP.Routes
			spForMethod.Timer = sharedSP.Timer
		}
		spInput := methodRoutes.InputRouteSet(cfg.Network.Flows.TSN.Background, cfg.Network.Flows.AVB.Background)
		spBG := methodRoutes.BGRouteSet(cfg.Network.Flows.TSN.Background, cfg.Network.Flows.AVB.Background)
		spForMethod.Objs_sp, _ = schedule.OBJ(plan.Network, cfg, osacoForMethod.KRoutes, spInput, spBG, method, true)
	}

	logger.Println("OSRO InitiatePlan completed (per-method)")
}
