package memorizer

import "src/pkg/logger"

func (OC *OmacoMemorizer) MOutputResults() {
	logger.Println()
	logger.Println("--- The experimental results are as follows ---")
	logger.Println("The average objective result for the Steiner Tree:")
	logger.Printf("O1: %f O2: %f O3: pass O4: %f \n", OC.average_obj_smt[0], OC.average_obj_smt[1], OC.average_obj_smt[3])
	logger.Println("The average objective result for the MDTC:")
	logger.Printf("O1: %f O2: %f O3: pass O4: %f \n", OC.average_obj_mdt[0], OC.average_obj_mdt[1], OC.average_obj_mdt[3])
	logger.Printf("Computering time: %v\n", OC.average_time_mdt)
	logger.Println("The average objective result for OSACO:")
	logger.Printf("1000ms: O1: %f O2: %f O3: pass O4: %f \n", OC.average_objs_osaco[4][0], OC.average_objs_osaco[4][1], OC.average_objs_osaco[4][3])
	logger.Printf("800ms: O1: %f O2: %f O3: pass O4: %f \n", OC.average_objs_osaco[3][0], OC.average_objs_osaco[3][1], OC.average_objs_osaco[3][3])
	logger.Printf("600ms: O1: %f O2: %f O3: pass O4: %f \n", OC.average_objs_osaco[2][0], OC.average_objs_osaco[2][1], OC.average_objs_osaco[2][3])
	logger.Printf("400ms: O1: %f O2: %f O3: pass O4: %f \n", OC.average_objs_osaco[1][0], OC.average_objs_osaco[1][1], OC.average_objs_osaco[1][3])
	logger.Printf("200ms: O1: %f O2: %f O3: pass O4: %f \n", OC.average_objs_osaco[0][0], OC.average_objs_osaco[0][1], OC.average_objs_osaco[0][3])
	logger.Printf("Computering time: %v\n", OC.average_time_osaco[4])
	logger.Println("The average objective result for OSACO_APTED:")
	logger.Printf("1000ms: O1: %f O2: %f O3: pass O4: %f \n", OC.average_objs_osaco_apted[4][0], OC.average_objs_osaco_apted[4][1], OC.average_objs_osaco_apted[4][3])
	logger.Printf("800ms: O1: %f O2: %f O3: pass O4: %f \n", OC.average_objs_osaco_apted[3][0], OC.average_objs_osaco_apted[3][1], OC.average_objs_osaco_apted[3][3])
	logger.Printf("600ms: O1: %f O2: %f O3: pass O4: %f \n", OC.average_objs_osaco_apted[2][0], OC.average_objs_osaco_apted[2][1], OC.average_objs_osaco_apted[2][3])
	logger.Printf("400ms: O1: %f O2: %f O3: pass O4: %f \n", OC.average_objs_osaco_apted[1][0], OC.average_objs_osaco_apted[1][1], OC.average_objs_osaco_apted[1][3])
	logger.Printf("200ms: O1: %f O2: %f O3: pass O4: %f \n", OC.average_objs_osaco_apted[0][0], OC.average_objs_osaco_apted[0][1], OC.average_objs_osaco_apted[0][3])
	logger.Printf("Computering time: %v\n", OC.average_time_osaco_apted[4])
	logger.Println()
}

func (OS *OsroMemorizer) MOutputResults() {
	logger.Println()
	logger.Println("--- The experimental results are as follows ---")
	logger.Println("The average objective result for the Shortest Path:")
	logger.Printf("O1: %f O2: %f O3: pass O4: %f \n", OS.average_obj_smt[0], OS.average_obj_smt[1], OS.average_obj_smt[3])
	logger.Printf("Computering time: %v\n", OS.average_time_mdt)
	logger.Println("The average objective result for OSACO (Path-based):")
	logger.Printf("1000ms: O1: %f O2: %f O3: pass O4: %f \n", OS.average_objs_osaco[4][0], OS.average_objs_osaco[4][1], OS.average_objs_osaco[4][3])
	logger.Printf("800ms: O1: %f O2: %f O3: pass O4: %f \n", OS.average_objs_osaco[3][0], OS.average_objs_osaco[3][1], OS.average_objs_osaco[3][3])
	logger.Printf("600ms: O1: %f O2: %f O3: pass O4: %f \n", OS.average_objs_osaco[2][0], OS.average_objs_osaco[2][1], OS.average_objs_osaco[2][3])
	logger.Printf("400ms: O1: %f O2: %f O3: pass O4: %f \n", OS.average_objs_osaco[1][0], OS.average_objs_osaco[1][1], OS.average_objs_osaco[1][3])
	logger.Printf("200ms: O1: %f O2: %f O3: pass O4: %f \n", OS.average_objs_osaco[0][0], OS.average_objs_osaco[0][1], OS.average_objs_osaco[0][3])
	logger.Printf("Computering time: %v\n", OS.average_time_osaco[4])
	logger.Println()
}
