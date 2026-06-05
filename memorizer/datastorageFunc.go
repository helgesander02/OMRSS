package memorizer

import (
	"encoding/csv"
	"fmt"
	"os"
	"src/pkg/logger"
	"time"
)

func (OC *OmacoMemorizer) MStoreData(fileName string, testcase int) {
	dirName := "data"
	createFolder(dirName)
	switchWorkingPath(dirName)
	createFolder(fileName)
	switchWorkingPath(fileName)

	// Write results to CSV file
	logger.Printf("Storing SMT data to CSV...\n")
	StoreCSV("SMT.csv", OC.average_obj_smt, testcase)
	logger.Printf("Storing MTDC data to CSV...\n")
	StoreCSV("MTDC.csv", OC.average_obj_mdt, testcase)
	logger.Printf("Storing OSACO data to CSV...\n")
	StoreCSV("OSACO_timeout_X5.csv", OC.average_objs_osaco[4], testcase)
	StoreCSV("OSACO_timeout_X4.csv", OC.average_objs_osaco[3], testcase)
	StoreCSV("OSACO_timeout_X3.csv", OC.average_objs_osaco[2], testcase)
	StoreCSV("OSACO_timeout_X2.csv", OC.average_objs_osaco[1], testcase)
	StoreCSV("OSACO_timeout_X1.csv", OC.average_objs_osaco[0], testcase)
	logger.Printf("Storing OSACO_APTED data to CSV...\n")
	StoreCSV("OSACO_APTED_timeout_X5.csv", OC.average_objs_osaco_apted[4], testcase)
	StoreCSV("OSACO_APTED_timeout_X4.csv", OC.average_objs_osaco_apted[3], testcase)
	StoreCSV("OSACO_APTED_timeout_X3.csv", OC.average_objs_osaco_apted[2], testcase)
	StoreCSV("OSACO_APTED_timeout_X2.csv", OC.average_objs_osaco_apted[1], testcase)
	StoreCSV("OSACO_APTED_timeout_X1.csv", OC.average_objs_osaco_apted[0], testcase)
	StoreComputeringTimeCSV("computering_time.csv", OC.average_time_mdt, OC.average_time_osaco[4], OC.average_time_osaco_apted[4], testcase)

	switchWorkingPath("../..")
}

func (OS *OsroMemorizer) MStoreData(fileName string, testcase int) {
	dirName := "data"
	createFolder(dirName)
	switchWorkingPath(dirName)
	createFolder(fileName)
	switchWorkingPath(fileName)

	// Write results to CSV file
	logger.Printf("Storing Shortest Path data to CSV...\n")
	StoreCSV("ShortestPath.csv", OS.average_obj_smt, testcase)
	logger.Printf("Storing OSACO (Path) data to CSV...\n")
	StoreCSV("OSACO_Path_timeout_X5.csv", OS.average_objs_osaco[4], testcase)
	StoreCSV("OSACO_Path_timeout_X4.csv", OS.average_objs_osaco[3], testcase)
	StoreCSV("OSACO_Path_timeout_X3.csv", OS.average_objs_osaco[2], testcase)
	StoreCSV("OSACO_Path_timeout_X2.csv", OS.average_objs_osaco[1], testcase)
	StoreCSV("OSACO_Path_timeout_X1.csv", OS.average_objs_osaco[0], testcase)
	StoreComputeringTimeCSV_OSRO("computering_time.csv", OS.average_time_mdt, OS.average_time_osaco[4], testcase)

	switchWorkingPath("../..")
}

func StoreCSV(name string, data [4]float64, testcase int) {
	logger.Printf("Opening file %s in append mode\n", name)
	csvFile, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		logger.Printf("File not found, creating file %s\n", name)
		csvFile, err = os.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			logger.Fatalf("Failed to create or open file data.csv: %v\n", err)
		}
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	err = writer.Write(Fary2Sary(data, testcase))
	if err != nil {
		logger.Fatalf("Failed to writer %s.csv: %v\n", name, err)
		return
	}

	writer.Flush()
	csvFile.Close()
}

func StoreComputeringTimeCSV(name string, average_time_mdt time.Duration, average_time_osaco time.Duration, average_time_osaco_apted time.Duration, testCase int) {
	logger.Printf("Opening file %s in append mode\n", name)
	csvFile, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		logger.Printf("File not found, creating file %s\n", name)
		csvFile, err = os.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			logger.Fatalf("Failed to create or open file data.csv: %v\n", err)
		}
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	err = writer.Write(Tary2Sary(average_time_mdt, average_time_osaco, average_time_osaco_apted, testCase))
	if err != nil {
		logger.Fatalf("Failed to writer %s.csv: %v\n", name, err)
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

func Tary2Sary(average_time_mdt time.Duration, average_time_osaco time.Duration, average_time_osaco_apted time.Duration, testCase int) []string {
	var (
		Sary []string
		Tary = []time.Duration{average_time_mdt, average_time_osaco, average_time_osaco_apted}
	)

	for _, t := range Tary {
		nanoseconds := float64(t.Milliseconds())
		Sary = append(Sary, fmt.Sprintf("%.6f", nanoseconds))
	}
	Sary = append(Sary, fmt.Sprintf("%d", testCase))

	return Sary
}

func StoreComputeringTimeCSV_OSRO(name string, average_time_sp time.Duration, average_time_osaco_path time.Duration, testCase int) {
	logger.Printf("Opening file %s in append mode\n", name)
	csvFile, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		logger.Printf("File not found, creating file %s\n", name)
		csvFile, err = os.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			logger.Fatalf("Failed to create or open file data.csv: %v\n", err)
		}
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	err = writer.Write(Tary2Sary_OSRO(average_time_sp, average_time_osaco_path, testCase))
	if err != nil {
		logger.Fatalf("Failed to writer %s.csv: %v\n", name, err)
		return
	}

	writer.Flush()
	csvFile.Close()
}

func Tary2Sary_OSRO(average_time_sp time.Duration, average_time_osaco_path time.Duration, testCase int) []string {
	var (
		Sary []string
		Tary = []time.Duration{average_time_sp, average_time_osaco_path}
	)

	for _, t := range Tary {
		nanoseconds := float64(t.Milliseconds())
		Sary = append(Sary, fmt.Sprintf("%.6f", nanoseconds))
	}
	Sary = append(Sary, fmt.Sprintf("%d", testCase))

	return Sary
}
