package memorizer

import "time"

func (OC *OmacoMemorizer) MAverage(testCase int) {
	OC.average_time_mdt = time.Duration(int(OC.average_time_mdt/time.Nanosecond)/testCase) * time.Nanosecond
	for i := 0; i < 5; i++ {
		for j := 0; j < 4; j++ {
			if i == 0 {
				OC.average_obj_smt[j] = OC.average_obj_smt[j] / float64(testCase)
				OC.average_obj_mdt[j] = OC.average_obj_mdt[j] / float64(testCase)
			}
			OC.average_objs_osaco[i][j] = OC.average_objs_osaco[i][j] / float64(testCase)
			OC.average_objs_osaco_apted[i][j] = OC.average_objs_osaco_apted[i][j] / float64(testCase)
		}
		OC.average_time_osaco[i] = time.Duration(int(OC.average_time_osaco[i]/time.Nanosecond)/testCase) * time.Nanosecond
		OC.average_time_osaco_apted[i] = time.Duration(int(OC.average_time_osaco_apted[i]/time.Nanosecond)/testCase) * time.Nanosecond
	}
}

func (OS *OsroMemorizer) MAverage(testCase int) {
	if testCase <= 0 {
		return
	}
	for _, method := range OS.methods {
		// SP averages
		spTime := OS.average_time_sp_by_method[method]
		*spTime = time.Duration(int(*spTime/time.Nanosecond)/testCase) * time.Nanosecond
		spObj := OS.average_obj_sp_by_method[method]
		for j := 0; j < 4; j++ {
			spObj[j] = spObj[j] / float64(testCase)
		}

		// OSACO averages
		osacoObj := OS.average_objs_osaco_by_method[method]
		osacoTime := OS.average_time_osaco_by_method[method]
		for i := 0; i < 5; i++ {
			for j := 0; j < 4; j++ {
				osacoObj[i][j] = osacoObj[i][j] / float64(testCase)
			}
			osacoTime[i] = time.Duration(int(osacoTime[i]/time.Nanosecond)/testCase) * time.Nanosecond
		}
	}
}

//func (mm3 *Memorizer3) MAverage(p plan.Plans) {
//
//}
