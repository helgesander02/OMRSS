package plan

import "src/pkg/logger"

func (plan *OMACO) ShowPlan() {
	logger.Println()
	logger.Println("--- The Steiner Tree final selected routing---")
	plan.SMT.Routes.Show_RouteSet()

	logger.Println()
	logger.Println("--- The Distance Tree final selected routing---")
	plan.MDTC.Routes.Show_RouteSet()

	logger.Println()
	logger.Println("--- 5th Spanning Tree ---")
	plan.OSACO.KRoutes.Show_KRouteSet()
	plan.OSACO.Timer[0].TimerExportData()
	plan.OSACO.Timer[1].TimerExportData()
	plan.OSACO.Timer[2].TimerExportData()
	plan.OSACO.Timer[3].TimerExportData()
	plan.OSACO.Timer[4].TimerExportData()

	logger.Println()
	logger.Println("--- The OSACO final selected routing ---")
	plan.OSACO.InputRoutes.Show_RouteSet()
	plan.OSACO.BGRoutes.Show_RouteSet()
}

func (plan *OSRO) ShowPlan() {
	logger.Println()
	logger.Println("--- The Shortest Path final selected routing---")
	plan.SP.Routes.Show_RouteSet()

	logger.Println()
	logger.Println("--- 5th K-Paths ---")
	plan.OSACO.KRoutes.Show_KRouteSet()

	plan.OSACO.Timer[0].TimerExportData()
	plan.OSACO.Timer[1].TimerExportData()
	plan.OSACO.Timer[2].TimerExportData()
	plan.OSACO.Timer[3].TimerExportData()
	plan.OSACO.Timer[4].TimerExportData()

	logger.Println()
	logger.Println("--- The OSACO (Path) final selected routing ---")
	plan.OSACO.InputRoutes.Show_RouteSet()
	plan.OSACO.BGRoutes.Show_RouteSet()
}

//func (plan *Plan3) ShowPlan() {
//
//}
