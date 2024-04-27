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
	// memorizer1 ...
	OMACO := new_OMACO_Memorizer()

	// Look-up table method
	memorizers := map[string]Memorizers{
		"omaco": OMACO,
		//memorizer2,
		//memorizer3,
		// ...
	}

	return memorizers
}
