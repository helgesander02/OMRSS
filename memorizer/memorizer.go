package memorizer

import "src/plan"

type Memorizers interface {
	M_Cumulative(plan.Plans)
	M_Average(int)
	M_Store_Data()
	M_Output_Results()
	M_Store_Files(string, int, int, int)
}

func New_Memorizers() map[string]Memorizers {
	// computer1 ...
	OMACO := New_OMACO_Memorizer()

	// Look-up table method
	Computers := map[string]Memorizers{
		"omaco": OMACO,
		//computer2,
		//computer3,
		// ...
	}

	return Computers
}

func New_OMACO_Memorizer() *OMACO_Memorizer {
	return &OMACO_Memorizer{}
}

//func New_plan2_Memorizer() *Memorizer2 {
//	return
//}

//func New_plan3_Memorizer() *Memorizer3 {
//	return
//}
