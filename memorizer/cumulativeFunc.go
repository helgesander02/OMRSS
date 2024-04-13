package memorizer

import "src/plan"

func (OC *OMACO_Memorizer) M_Cumulative(p plan.Plans) {
	for i := 0; i < 5; i++ {
		for j := 0; j < 4; j++ {
			if i == 0 {
				OC.average_obj_smt[j] += p.(*plan.OMACO).SMT.Objs_smt[j]
				OC.average_obj_mdt[j] += p.(*plan.OMACO).MDTC.Objs_mdtc[j]
			}
			OC.average_objs_osaco[i][j] += p.(*plan.OMACO).OSACO.Objs_osaco[i][j]
		}
	}
}

//func (cpt2 *Computer2) M_Cumulative(p plan.Plans) {
//
//}

//func (cpt3 *Computer3) M_Cumulative(p plan.Plans) {
//
//}
