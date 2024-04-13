package memorizer

func (OC *OMACO_Memorizer) M_Average(test_case int) {
	for i := 0; i < 5; i++ {
		for j := 0; j < 4; j++ {
			if i == 0 {
				OC.average_obj_smt[j] = OC.average_obj_smt[j] / float64(test_case)
				OC.average_obj_mdt[j] = OC.average_obj_mdt[j] / float64(test_case)
			}
			OC.average_objs_osaco[i][j] = OC.average_objs_osaco[i][j] / float64(test_case)
		}
	}
}

//func (cpt2 *Computer2) M_Average(p plan.Plans) {
//
//}

//func (cpt3 *Computer3) M_Average(p plan.Plans) {
//
//}
