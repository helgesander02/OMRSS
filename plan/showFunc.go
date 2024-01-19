package plan

import "fmt"

func (plan *Plan) Show_Plan() {
	fmt.Println()
	fmt.Println("--- The Steiner Tree final selected routing---")
	plan.SMT.Trees.Show_Trees_Set()

	fmt.Println()
	fmt.Println("--- The Distance Tree final selected routing---")
	plan.MDTC.Trees.Show_Trees_Set()

	fmt.Println()
	fmt.Println("--- 5th Spanning Tree ---")
	plan.OSACO.KTrees.Show_kTrees_Set()

	fmt.Println()
	fmt.Println("--- The OSACO final selected routing ---")
	plan.OSACO.InputTrees.Show_Trees_Set()
	plan.OSACO.BGTrees.Show_Trees_Set()
}
