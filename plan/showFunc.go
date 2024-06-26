package plan

import "fmt"

func (plan *OMACO) Show_Plan() {
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

//func (plan *Plan2) Show_Plan() {
//
//}

//func (plan *Plan3) Show_Plan() {
//
//}
