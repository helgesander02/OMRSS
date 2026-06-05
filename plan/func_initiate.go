package plan

import (
	"src/pkg/config"
	"src/pkg/logger"
	"src/plan/schedule"
)

func (plan *OMACO) InitiatePlan(costSetting [4]int, cfg *config.Config) {
	// algo run
	logger.Println("Steiner Tree")
	logger.Println("----------------------------------------")
	plan.SMT.SMT_Run(plan.Network, cfg)

	logger.Println()
	logger.Println("MDTC")
	logger.Println("----------------------------------------")
	plan.MDTC.MDTC_Run(plan.Network, cfg)

	logger.Println()
	logger.Println("OSACO")
	logger.Println("----------------------------------------")
	plan.OSACO.OSACO_Initial_Settings(plan.Network, cfg, plan.SMT.Trees)
	// The timeout of each run is set as 100~1000 ms (200ms, 400ms, 600ms, 800ms, 1000ms)
	for i := 0; i < 5; i++ {
		plan.OSACO.Objs_osaco[i] = plan.OSACO.OSACO_Run(plan.Network, cfg, i, costSetting)
	}

	logger.Println()
	logger.Println("OSACO_APTED")
	logger.Println("----------------------------------------")
	plan.OSACO_APTED.OSACO_Initial_Settings(plan.Network, cfg, plan.SMT.Trees)
	// The timeout of each run is set as 100~1000 ms (200ms, 400ms, 600ms, 800ms, 1000ms)
	for i := 0; i < 5; i++ {
		plan.OSACO_APTED.Objs_osaco[i] = plan.OSACO_APTED.OSACO_Run(plan.Network, cfg, i, costSetting)
	}

	obj_smt, _ := schedule.OBJ(
		plan.Network,
		cfg,
		plan.OSACO.KTrees,
		plan.SMT.Trees.InputTreeSet(cfg.Network.Flows.TSN.Background, cfg.Network.Flows.AVB.Background),
		plan.SMT.Trees.BGTreeSet(cfg.Network.Flows.TSN.Background, cfg.Network.Flows.AVB.Background),
		costSetting,
		true,
	)

	obj_mdt, _ := schedule.OBJ(
		plan.Network,
		cfg,
		plan.OSACO.KTrees,
		plan.MDTC.Trees.InputTreeSet(cfg.Network.Flows.TSN.Background, cfg.Network.Flows.AVB.Background),
		plan.MDTC.Trees.BGTreeSet(cfg.Network.Flows.TSN.Background, cfg.Network.Flows.AVB.Background),
		costSetting,
		true,
	)

	plan.SMT.Objs_smt = obj_smt
	plan.MDTC.Objs_mdtc = obj_mdt

	if obj_mdt[0] != 0 || obj_mdt[1] != 0 {
		plan.MDTC.Timer.TimerMax()
	}

}

func (plan *OSRO) InitiatePlan(costSetting [4]int, cfg *config.Config) {
	// algo run
	logger.Println("Shortest Path")
	logger.Println("----------------------------------------")
	plan.SP.SP_Run(plan.Network, cfg)

	logger.Println()
	logger.Println("OSACO (Path-based)")
	logger.Println("----------------------------------------")
	plan.OSACO_Path.OSACO_Initial_Settings_Path(plan.Network, cfg, plan.SP.Paths)

	// The timeout of each run is set as 100~1000 ms (200ms, 400ms, 600ms, 800ms, 1000ms)
	for i := 0; i < 5; i++ {
		plan.OSACO_Path.Objs_osaco[i] = plan.OSACO_Path.OSACO_Run_Path(plan.Network, i, costSetting)
	}

	// Note: Path-based objective calculation (schedule.OBJ_Path) is not yet implemented
	// This will be added when path-based scheduling is completed
	plan.SP.Objs_sp = [4]float64{0, 0, 0, 0}

	logger.Println("OSRO InitiatePlan completed")
}
