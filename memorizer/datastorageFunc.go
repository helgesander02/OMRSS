package memorizer

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"
)

func (OC *OMACO_Memorizer) M_Store_Data(fileName string, testcase int) {
	dirName := "data"
	createFolder(dirName)
	switchWorkingPath(dirName)
	createFolder(fileName)
	switchWorkingPath(fileName)

	// Write results to CSV file
	log.Printf("Storing SMT data to CSV...\n")
	StoreCSV("SMT.csv", OC.average_obj_smt, testcase)
	log.Printf("Storing MTDC data to CSV...\n")
	StoreCSV("MTDC.csv", OC.average_obj_mdt, testcase)
	log.Printf("Storing OSACO data to CSV...\n")
	StoreCSV("OSACO_timeout_X5.csv", OC.average_objs_osaco[4], testcase)
	StoreCSV("OSACO_timeout_X4.csv", OC.average_objs_osaco[3], testcase)
	StoreCSV("OSACO_timeout_X3.csv", OC.average_objs_osaco[2], testcase)
	StoreCSV("OSACO_timeout_X2.csv", OC.average_objs_osaco[1], testcase)
	StoreCSV("OSACO_timeout_X1.csv", OC.average_objs_osaco[0], testcase)
	log.Printf("Storing OSACO_IAS data to CSV...\n")
	StoreCSV("OSACO_IAS_timeout_X5.csv", OC.average_objs_osaco_ias[4], testcase)
	StoreCSV("OSACO_IAS_timeout_X4.csv", OC.average_objs_osaco_ias[3], testcase)
	StoreCSV("OSACO_IAS_timeout_X3.csv", OC.average_objs_osaco_ias[2], testcase)
	StoreCSV("OSACO_IAS_timeout_X2.csv", OC.average_objs_osaco_ias[1], testcase)
	StoreCSV("OSACO_IAS_timeout_X1.csv", OC.average_objs_osaco_ias[0], testcase)
	StoreComputeringTimeCSV("computering_time.csv", OC.average_time_mdt, OC.average_time_osaco[4], OC.average_time_osaco_ias[4], testcase)

	switchWorkingPath("../..")
}

func StoreCSV(name string, data [4]float64, testcase int) {
	log.Printf("Opening file %s in append mode\n", name)
	csvFile, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("File not found, creating file %s\n", name)
		csvFile, err = os.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to create or open file data.csv: %v\n", err)
		}
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	err = writer.Write(Fary2Sary(data, testcase))
	if err != nil {
		log.Fatalf("Failed to writer %s.csv: %v\n", name, err)
		return
	}

	writer.Flush()
	csvFile.Close()
}

func StoreComputeringTimeCSV(name string, average_time_mdt time.Duration, average_time_osaco time.Duration, average_time_osaco_ias time.Duration, test_case int) {
	log.Printf("Opening file %s in append mode\n", name)
	csvFile, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("File not found, creating file %s\n", name)
		csvFile, err = os.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to create or open file data.csv: %v\n", err)
		}
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	err = writer.Write(Tary2Sary(average_time_mdt, average_time_osaco, average_time_osaco_ias, test_case))
	if err != nil {
		log.Fatalf("Failed to writer %s.csv: %v\n", name, err)
		return
	}

	writer.Flush()
	csvFile.Close()
}

func Fary2Sary(Fary [4]float64, testcase int) []string {
	var Sary []string
	for _, f := range Fary {
		str := fmt.Sprintf("%.6f", f)
		Sary = append(Sary, str)
	}
	str := fmt.Sprintf("%d", testcase)
	Sary = append(Sary, str)

	return Sary
}

func Tary2Sary(average_time_mdt time.Duration, average_time_osaco time.Duration, average_time_osaco_ias time.Duration, test_case int) []string {
	var (
		Sary []string
		Tary = []time.Duration{average_time_mdt, average_time_osaco, average_time_osaco_ias}
	)

	for _, t := range Tary {
		nanoseconds := float64(t.Milliseconds())
		Sary = append(Sary, fmt.Sprintf("%.6f", nanoseconds))
	}
	Sary = append(Sary, fmt.Sprintf("%d", test_case))

	return Sary
}

//func (mm2 *Memorizer2) M_Store_Data(test_case int) {
//
//}

//func (mm3 *Memorizer2) M_Store_Data(test_case int) {
//
//}
