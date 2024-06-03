package memorizer

import (
	"fmt"
	"log"
	"os"
)

func (OC *OMACO_Memorizer) M_Store_Files(topology_name string, test_case int, tsn int, avb int) {
	// Check if the directory exists
	dirName := "result"
	_, err := os.Stat(dirName)
	if err != nil {
		if os.IsNotExist(err) {
			// Create the directory if it doesn't exist
			log.Printf("Creating directory %s\n", dirName)
			err := os.MkdirAll(dirName, os.ModePerm)
			if err != nil {
				log.Fatalf("Failed to create directory %s: %v\n", dirName, err)
			}
		} else {
			// Handle other errors accessing the directory
			log.Fatalf("Failed to check directory %s: %v\n", dirName, err)
		}
	}

	// Change the working directory
	log.Printf("Switching to directory %s\n", dirName)
	err = os.Chdir(dirName)
	if err != nil {
		log.Fatalf("Failed to switch to directory %s: %v\n", dirName, err)
	}

	text := "--- The experimental results are as follows --- \n"
	text += "The average objective result for the Steiner Tree:\n"
	text += fmt.Sprintf("O1: %f O2: %f O3: pass O4: %f \n", OC.average_obj_smt[0], OC.average_obj_smt[1], OC.average_obj_smt[3])
	text += "The average objective result for the MDTC:\n"
	text += fmt.Sprintf("O1: %f O2: %f O3: pass O4: %f \n", OC.average_obj_mdt[0], OC.average_obj_mdt[1], OC.average_obj_mdt[3])
	text += fmt.Sprintf("Computering time: %v\n", OC.average_time_mdt)
	text += "The average objective result for OSACO:\n"
	text += fmt.Sprintf("1000ms: O1: %f O2: %f O3: pass O4: %f \n", OC.average_objs_osaco[4][0], OC.average_objs_osaco[4][1], OC.average_objs_osaco[4][3])
	text += fmt.Sprintf("Computering time: %v\n", OC.average_time_osaco[4])
	text += fmt.Sprintf("800ms: O1: %f O2: %f O3: pass O4: %f \n", OC.average_objs_osaco[3][0], OC.average_objs_osaco[3][1], OC.average_objs_osaco[3][3])
	text += fmt.Sprintf("Computering time: %v\n", OC.average_time_osaco[3])
	text += fmt.Sprintf("600ms: O1: %f O2: %f O3: pass O4: %f \n", OC.average_objs_osaco[2][0], OC.average_objs_osaco[2][1], OC.average_objs_osaco[2][3])
	text += fmt.Sprintf("Computering time: %v\n", OC.average_time_osaco[2])
	text += fmt.Sprintf("400ms: O1: %f O2: %f O3: pass O4: %f \n", OC.average_objs_osaco[1][0], OC.average_objs_osaco[1][1], OC.average_objs_osaco[1][3])
	text += fmt.Sprintf("Computering time: %v\n", OC.average_time_osaco[1])
	text += fmt.Sprintf("200ms: O1: %f O2: %f O3: pass O4: %f \n", OC.average_objs_osaco[0][0], OC.average_objs_osaco[0][1], OC.average_objs_osaco[0][3])
	text += fmt.Sprintf("Computering time: %v\n", OC.average_time_osaco[0])

	// Try to open the file in append mode first
	name := fmt.Sprintf("%s_testcase%d_tsn%d_avb%d.txt", topology_name, test_case, tsn, avb)
	log.Printf("Opening file %s in append mode\n", name)
	file, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		// If it fails, create the file
		log.Printf("File not found, creating file %s\n", name)
		file, err = os.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to create or open file %s: %v\n", name, err)
		}
	}

	defer file.Close() // Ensure file is closed even if errors occur

	log.Printf("Writing text to file %s\n", name)
	fmt.Fprintln(file, text)

	// Change the working directory back
	log.Printf("Switching back to parent directory\n")
	err = os.Chdir("..")
	if err != nil {
		log.Fatalf("Failed to switch directories: %v\n", err)
	}
}

//func (mm2 *Memorizer2) M_Store_Files(test_case int) {
//
//}

//func (mm3 *Memorizer2) M_Store_Files(test_case int) {
//
//}
