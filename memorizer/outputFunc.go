package memorizer

import "fmt"

func (OC *OMACO_Memorizer) M_Output_Results() {
	fmt.Println()
	fmt.Println("--- The experimental results are as follows ---")
	fmt.Println("The average objective result for the Steiner Tree:")
	fmt.Printf("O1: %f O2: %f O3: pass O4: %f \n", OC.average_obj_smt[0], OC.average_obj_smt[1], OC.average_obj_smt[3])
	fmt.Println("The average objective result for the MDTC:")
	fmt.Printf("O1: %f O2: %f O3: pass O4: %f \n", OC.average_obj_mdt[0], OC.average_obj_mdt[1], OC.average_obj_mdt[3])
	fmt.Println("The average objective result for OSACO:")
	fmt.Printf("1000ms: O1: %f O2: %f O3: pass O4: %f \n", OC.average_objs_osaco[4][0], OC.average_objs_osaco[4][1], OC.average_objs_osaco[4][3])
	fmt.Printf("800ms: O1: %f O2: %f O3: pass O4: %f \n", OC.average_objs_osaco[3][0], OC.average_objs_osaco[3][1], OC.average_objs_osaco[3][3])
	fmt.Printf("600ms: O1: %f O2: %f O3: pass O4: %f \n", OC.average_objs_osaco[2][0], OC.average_objs_osaco[2][1], OC.average_objs_osaco[2][3])
	fmt.Printf("400ms: O1: %f O2: %f O3: pass O4: %f \n", OC.average_objs_osaco[1][0], OC.average_objs_osaco[1][1], OC.average_objs_osaco[1][3])
	fmt.Printf("200ms: O1: %f O2: %f O3: pass O4: %f \n", OC.average_objs_osaco[0][0], OC.average_objs_osaco[0][1], OC.average_objs_osaco[0][3])
	fmt.Println()
}

//func (mm2 *Memorizer2) M_Output_Results()(test_case int) {
//
//}

//func (mmt3 *Memorizer2) M_Output_Results()(test_case int) {
//
//}
