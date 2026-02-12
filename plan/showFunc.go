package plan

import "fmt"

func (plan *OMACO) ShowPlan() {
	fmt.Println()
	fmt.Println("--- The Steiner Tree final selected routing---")
	plan.SMT.Trees.Show_Trees_Set()

	fmt.Println()
	fmt.Println("--- The Distance Tree final selected routing---")
	plan.MDTC.Trees.Show_Trees_Set()

	fmt.Println()
	fmt.Println("--- 5th Spanning Tree ---")
	plan.OSACO.KTrees.Show_kTrees_Set()
	plan.OSACO.Timer[0].TimerExportData()
	plan.OSACO.Timer[1].TimerExportData()
	plan.OSACO.Timer[2].TimerExportData()
	plan.OSACO.Timer[3].TimerExportData()
	plan.OSACO.Timer[4].TimerExportData()

	fmt.Println()
	fmt.Println("--- The OSACO final selected routing ---")
	plan.OSACO.InputTrees.Show_Trees_Set()
	plan.OSACO.BGTrees.Show_Trees_Set()
}

func (plan *OSRO) ShowPlan() {
	fmt.Println()
	fmt.Println("--- The Shortest Path final selected routing---")
	plan.SP.Paths.Show_Paths_Set()

	fmt.Println()
	fmt.Println("--- 5th K-Paths ---")
	plan.OSACO_Path.KPaths.Show_KPaths_Set()

	plan.OSACO_Path.Timer[0].TimerExportData()
	plan.OSACO_Path.Timer[1].TimerExportData()
	plan.OSACO_Path.Timer[2].TimerExportData()
	plan.OSACO_Path.Timer[3].TimerExportData()
	plan.OSACO_Path.Timer[4].TimerExportData()

	fmt.Println()
	fmt.Println("--- The OSACO (Path) final selected routing ---")
	plan.OSACO_Path.InputPaths.Show_Paths_Set()
	plan.OSACO_Path.BGPaths.Show_Paths_Set()
}

//func (plan *Plan3) ShowPlan() {
//
//}
