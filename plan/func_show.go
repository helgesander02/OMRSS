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
	// SP routes are shared across methods (same Yen's K-path output);
	// pick the first method to surface them.
	first := plan.Methods[0]

	logger.Println()
	logger.Println("--- The Shortest Path final selected routing (shared across methods) ---")
	plan.SPByMethod[first].Routes.Show_RouteSet()

	for _, m := range plan.Methods {
		osaco := plan.OSACOByMethod[m]
		logger.Println()
		logger.Printf("--- 5th K-Paths — method: %s ---\n", m)
		osaco.KRoutes.Show_KRouteSet()

		for i := 0; i < 5; i++ {
			osaco.Timer[i].TimerExportData()
		}

		logger.Println()
		logger.Printf("--- The OSACO (Path) final selected routing — method: %s ---\n", m)
		osaco.InputRoutes.Show_RouteSet()
		osaco.BGRoutes.Show_RouteSet()
	}
}

//func (plan *Plan3) ShowPlan() {
//
//}
