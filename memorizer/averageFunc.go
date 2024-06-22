package memorizer

import "time"

func (OC *OMACO_Memorizer) M_Average(test_case int) {
	OC.average_time_mdt = time.Duration(int(OC.average_time_mdt/time.Nanosecond)/test_case) * time.Nanosecond
	for i := 0; i < 5; i++ {
		for j := 0; j < 4; j++ {
			if i == 0 {
				OC.average_obj_smt[j] = OC.average_obj_smt[j] / float64(test_case)
				OC.average_obj_mdt[j] = OC.average_obj_mdt[j] / float64(test_case)
			}
			OC.average_objs_osaco[i][j] = OC.average_objs_osaco[i][j] / float64(test_case)
			OC.average_objs_osaco_ias[i][j] = OC.average_objs_osaco_ias[i][j] / float64(test_case)
			OC.average_objs_osaco_aas[i][j] = OC.average_objs_osaco_aas[i][j] / float64(test_case)
		}
		OC.average_time_osaco[i] = time.Duration(int(OC.average_time_osaco[i]/time.Nanosecond)/test_case) * time.Nanosecond
		OC.average_time_osaco_ias[i] = time.Duration(int(OC.average_time_osaco_ias[i]/time.Nanosecond)/test_case) * time.Nanosecond
		OC.average_time_osaco_aas[i] = time.Duration(int(OC.average_time_osaco_aas[i]/time.Nanosecond)/test_case) * time.Nanosecond
	}
}

//func (mm2 *Memorizer2) M_Average(p plan.Plans) {
//
//}

//func (mm3 *Memorizer3) M_Average(p plan.Plans) {
//
//}
