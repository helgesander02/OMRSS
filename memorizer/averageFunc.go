package memorizer

import "time"

func (OC *OMACO_Memorizer) MAverage(testCase int) {
	OC.average_time_mdt = time.Duration(int(OC.average_time_mdt/time.Nanosecond)/testCase) * time.Nanosecond
	for i := 0; i < 5; i++ {
		for j := 0; j < 4; j++ {
			if i == 0 {
				OC.average_obj_smt[j] = OC.average_obj_smt[j] / float64(testCase)
				OC.average_obj_mdt[j] = OC.average_obj_mdt[j] / float64(testCase)
			}
			OC.average_objs_osaco[i][j] = OC.average_objs_osaco[i][j] / float64(testCase)
			OC.average_objs_osaco_ias[i][j] = OC.average_objs_osaco_ias[i][j] / float64(testCase)
		}
		OC.average_time_osaco[i] = time.Duration(int(OC.average_time_osaco[i]/time.Nanosecond)/testCase) * time.Nanosecond
		OC.average_time_osaco_ias[i] = time.Duration(int(OC.average_time_osaco_ias[i]/time.Nanosecond)/testCase) * time.Nanosecond
	}
}

func (mm2 *OSRO_Memorizer) MAverage(testCase int) {

}

//func (mm3 *Memorizer3) MAverage(p plan.Plans) {
//
//}
