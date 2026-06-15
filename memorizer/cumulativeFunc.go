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
	osro := p.(*plan.OSRO)
	for _, method := range osro.Methods {
		OS.ensureMethod(method)

		sp := osro.SPByMethod[method]
		osaco := osro.OSACOByMethod[method]

		// SP cumulative — note SP.Routes / SP.Timer are shared across
		// methods, so each method records the same routing time. The
		// per-method differentiation lives in Objs_sp because OBJ was
		// scored with method scope.
		if sp.Timer != nil {
			*OS.average_time_sp_by_method[method] += sp.Timer.TimerOutputData()
		}
		for j := 0; j < 4; j++ {
			OS.average_obj_sp_by_method[method][j] += sp.Objs_sp[j]
		}

		// OSACO cumulative — every (timeout, objective) is per method.
		for i := 0; i < 5; i++ {
			for j := 0; j < 4; j++ {
				OS.average_objs_osaco_by_method[method][i][j] += osaco.Objs_osaco[i][j]
			}
			if osaco.Timer[i] != nil {
				OS.average_time_osaco_by_method[method][i] += osaco.Timer[i].TimerOutputData()
			}
		}
	}
}
