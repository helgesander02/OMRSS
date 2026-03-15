package memorizer

import "src/plan"

func (OC *OmacoMemorizer) MCumulative(p plan.Plans) {
	OC.average_time_mdt += p.(*plan.OMACO).MDTC.Timer.TimerOutputData()
	for i := 0; i < 5; i++ {
		for j := 0; j < 4; j++ {
			if i == 0 {
				OC.average_obj_smt[j] += p.(*plan.OMACO).SMT.Objs_smt[j]
				OC.average_obj_mdt[j] += p.(*plan.OMACO).MDTC.Objs_mdtc[j]
			}
			OC.average_objs_osaco[i][j] += p.(*plan.OMACO).OSACO.Objs_osaco[i][j]
			OC.average_objs_osaco_apted[i][j] += p.(*plan.OMACO).OSACO_APTED.Objs_osaco[i][j]
		}
		OC.average_time_osaco[i] += p.(*plan.OMACO).OSACO.Timer[i].TimerOutputData()
		OC.average_time_osaco_apted[i] += p.(*plan.OMACO).OSACO_APTED.Timer[i].TimerOutputData()
	}
}

func (OS *OsroMemorizer) MCumulative(p plan.Plans) {
	OS.average_time_mdt += p.(*plan.OSRO).SP.Timer.TimerOutputData()
	for i := 0; i < 5; i++ {
		for j := 0; j < 4; j++ {
			if i == 0 {
				OS.average_obj_smt[j] += p.(*plan.OSRO).SP.Objs_sp[j]
			}
			OS.average_objs_osaco[i][j] += p.(*plan.OSRO).OSACO_Path.Objs_osaco[i][j]
		}
		OS.average_time_osaco[i] += p.(*plan.OSRO).OSACO_Path.Timer[i].TimerOutputData()
	}
}
