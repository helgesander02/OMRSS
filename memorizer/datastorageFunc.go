package memorizer

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func (OC *OMACO_Memorizer) M_Store_Data() {
	// Check if the directory exists
	dirName := "data"
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

	log.Printf("Storing SMT data to CSV...\n")
	StoreCSV("SMT.csv", OC.average_obj_smt)
	log.Printf("Storing MTDC data to CSV...\n")
	StoreCSV("MTDC.csv", OC.average_obj_mdt)
	log.Printf("Storing OSACO data to CSV...\n")
	StoreCSV("OSACO1000ms.csv", OC.average_objs_osaco[4])
	StoreCSV("OSACO800ms.csv", OC.average_objs_osaco[3])
	StoreCSV("OSACO600ms.csv", OC.average_objs_osaco[2])
	StoreCSV("OSACO400ms.csv", OC.average_objs_osaco[1])
	StoreCSV("OSACO200ms.csv", OC.average_objs_osaco[0])
	log.Printf("Storing OSACO_IAS data to CSV...\n")
	StoreCSV("OSACO1000ms.csv", OC.average_objs_osaco_ias[4])
	StoreCSV("OSACO800ms.csv", OC.average_objs_osaco_ias[3])
	StoreCSV("OSACO600ms.csv", OC.average_objs_osaco_ias[2])
	StoreCSV("OSACO400ms.csv", OC.average_objs_osaco_ias[1])
	StoreCSV("OSACO200ms.csv", OC.average_objs_osaco_ias[0])

	// Change the working directory back
	log.Printf("Switching back to parent directory\n")
	err = os.Chdir("..")
	if err != nil {
		log.Fatalf("Failed to switch directories: %v\n", err)
	}
}

func StoreCSV(name string, data [4]float64) {
	log.Printf("Opening file %s in append mode\n", name)
	csvFile, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		// If it fails, create the file
		log.Printf("File not found, creating file %s\n", name)
		csvFile, err = os.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to create or open file data.csv: %v\n", err)
		}
	}
	defer csvFile.Close() // Ensure file is closed even if errors occur

	writer := csv.NewWriter(csvFile)
	err = writer.Write(Fary2Sary(data))
	if err != nil {
		log.Fatalf("Failed to writer %s.csv: %v\n", name, err)
		return
	}

	writer.Flush()
	csvFile.Close()
}

func Fary2Sary(Fary [4]float64) []string {
	var Sary []string
	for _, f := range Fary {
		str := fmt.Sprintf("%.6f", f)
		Sary = append(Sary, str)
	}

	return Sary
}

//func (mm2 *Memorizer2) M_Store_Data(test_case int) {
//
//}

//func (mm3 *Memorizer2) M_Store_Data(test_case int) {
//
//}
