package memorizer

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func (OC *OmacoMemorizer) MStoreFile(fileName string) {
	text := averageDataToResult(fileName)

	dirName := "result"
	createFolder(dirName)
	switchWorkingPath(dirName)

	// Try to open the file in append mode first
	txtName := fileName + ".txt"
	log.Printf("Opening file %s in append mode\n", txtName)
	file, err := os.OpenFile(txtName, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		// If it fails, create the file
		log.Printf("File not found, creating file %s\n", txtName)
		file, err = os.OpenFile(txtName, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to create or open file %s: %v\n", txtName, err)
		}
	}
	defer file.Close()
	log.Printf("Writing text to file %s\n", txtName)
	fmt.Fprintln(file, text)

	switchWorkingPath("..")
}

func (OS *OsroMemorizer) MStoreFile(fileName string) {
	text := averageDataToResult_OSRO(fileName)

	dirName := "result"
	createFolder(dirName)
	switchWorkingPath(dirName)

	// Try to open the file in append mode first
	txtName := fileName + ".txt"
	log.Printf("Opening file %s in append mode\n", txtName)
	file, err := os.OpenFile(txtName, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		// If it fails, create the file
		log.Printf("File not found, creating file %s\n", txtName)
		file, err = os.OpenFile(txtName, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to create or open file %s: %v\n", txtName, err)
		}
	}
	defer file.Close()
	log.Printf("Writing text to file %s\n", txtName)
	fmt.Fprintln(file, text)

	switchWorkingPath("..")
}

func averageDataToResult(fileName string) string {
	data, testcaseNumbers := getAverageData(fileName)

	text := fmt.Sprintf("( testcase numbers: %d ) ", testcaseNumbers)
	text += "--- The experimental results are as follows --- \n"
	text += "The average objective result for the Steiner Tree:\n"
	text += fmt.Sprintf("O1: %f O2: %f O3: pass O4: %f \n", data["average_obj_smt_o1"], data["average_obj_smt_o2"], data["average_obj_smt_o4"])
	text += "The average objective result for the MDTC:\n"
	text += fmt.Sprintf("O1: %f O2: %f O3: pass O4: %f \n", data["average_obj_mdt_o1"], data["average_obj_mdt_o2"], data["average_obj_mdt_o4"])
	text += fmt.Sprintf("Computering time: %v ms\n", data["average_time_mdt"])
	text += "The average objective result for OSACO:\n"
	text += fmt.Sprintf("timeout_X5: O1: %f O2: %f O3: pass O4: %f \n", data["average_objs_osaco_X5_O1"], data["average_objs_osaco_X5_O2"], data["average_objs_osaco_X5_O4"])
	text += fmt.Sprintf("timeout_X4: O1: %f O2: %f O3: pass O4: %f \n", data["average_objs_osaco_X4_O1"], data["average_objs_osaco_X4_O2"], data["average_objs_osaco_X4_O4"])
	text += fmt.Sprintf("timeout_X3: O1: %f O2: %f O3: pass O4: %f \n", data["average_objs_osaco_X3_O1"], data["average_objs_osaco_X3_O2"], data["average_objs_osaco_X3_O4"])
	text += fmt.Sprintf("timeout_X2: O1: %f O2: %f O3: pass O4: %f \n", data["average_objs_osaco_X2_O1"], data["average_objs_osaco_X2_O2"], data["average_objs_osaco_X2_O4"])
	text += fmt.Sprintf("timeout_X1: O1: %f O2: %f O3: pass O4: %f \n", data["average_objs_osaco_X1_O1"], data["average_objs_osaco_X1_O2"], data["average_objs_osaco_X1_O4"])
	text += fmt.Sprintf("Computering time: %v ms\n", data["average_time_osaco"])
	text += "The average objective result for OSACO_APTED:\n"
	text += fmt.Sprintf("timeout_X5: O1: %f O2: %f O3: pass O4: %f \n", data["average_objs_osaco_apted_X5_O1"], data["average_objs_osaco_apted_X5_O2"], data["average_objs_osaco_apted_X5_O4"])
	text += fmt.Sprintf("timeout_X4: O1: %f O2: %f O3: pass O4: %f \n", data["average_objs_osaco_apted_X4_O1"], data["average_objs_osaco_apted_X4_O2"], data["average_objs_osaco_apted_X4_O4"])
	text += fmt.Sprintf("timeout_X3: O1: %f O2: %f O3: pass O4: %f \n", data["average_objs_osaco_apted_X3_O1"], data["average_objs_osaco_apted_X3_O2"], data["average_objs_osaco_apted_X3_O4"])
	text += fmt.Sprintf("timeout_X2: O1: %f O2: %f O3: pass O4: %f \n", data["average_objs_osaco_apted_X2_O1"], data["average_objs_osaco_apted_X2_O2"], data["average_objs_osaco_apted_X2_O4"])
	text += fmt.Sprintf("timeout_X1: O1: %f O2: %f O3: pass O4: %f \n", data["average_objs_osaco_apted_X1_O1"], data["average_objs_osaco_apted_X1_O2"], data["average_objs_osaco_apted_X1_O4"])
	text += fmt.Sprintf("Computering time: %v ms\n", data["average_time_osaco_apted"])

	return text
}

func getAverageData(fileName string) (map[string]float64, int) {
	currentDir, _ := os.Getwd()
	dir := filepath.Join(currentDir + "/data/" + fileName)
	data := make(map[string]float64)

	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("Error reading directory: %v", err)
	}

	testcaseNumbers := 0
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(dir, file.Name())
		fmt.Printf("Reading file: %s\n", filePath)
		csvFile, err := os.Open(filePath)
		if err != nil {
			log.Printf("Error opening file %s: %v", filePath, err)
			continue
		}
		defer csvFile.Close()

		reader := csv.NewReader(csvFile)
		columns := make([][]string, 0)

		// Read columns instead of rows
		for {
			row, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("Error reading CSV from file %s: %v", filePath, err)
				continue
			}
			if len(columns) == 0 {
				columns = make([][]string, len(row))
			}
			for i, value := range row {
				columns[i] = append(columns[i], value)
			}
		}

		convertToFloat := func(row []string) (float64, error) {
			sum := 0.0
			for _, value := range row {
				num, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return 0, fmt.Errorf("error converting value %s to float64: %v", value, err)
				}
				sum += num
			}
			return sum / float64(len(row)), err
		}

		switch file.Name() {
		case "computering_time.csv":
			if val, err := convertToFloat(columns[0]); err == nil {
				data["average_time_mdt"] = val
			}
			if val, err := convertToFloat(columns[1]); err == nil {
				data["average_time_osaco"] = val
			}
			if val, err := convertToFloat(columns[2]); err == nil {
				data["average_time_osaco_apted"] = val
			}
			if val, err := convertToFloat(columns[3]); err == nil {
				testcaseNumbers = int(val)
			}

		case "MTDC.csv":
			if val, err := convertToFloat(columns[0]); err == nil {
				data["average_obj_mdt_o1"] = val
			}
			if val, err := convertToFloat(columns[1]); err == nil {
				data["average_obj_mdt_o2"] = val
			}
			if val, err := convertToFloat(columns[3]); err == nil {
				data["average_obj_mdt_o4"] = val
			}

		case "SMT.csv":
			if val, err := convertToFloat(columns[0]); err == nil {
				data["average_obj_smt_o1"] = val
			}
			if val, err := convertToFloat(columns[1]); err == nil {
				data["average_obj_smt_o2"] = val
			}
			if val, err := convertToFloat(columns[3]); err == nil {
				data["average_obj_smt_o4"] = val
			}
		}

		// OSACO timeout cases
		for i := 1; i <= 5; i++ {
			if file.Name() == fmt.Sprintf("OSACO_timeout_X%d.csv", i) {
				if val, err := convertToFloat(columns[0]); err == nil {
					data[fmt.Sprintf("average_objs_osaco_X%d_O1", i)] = val
				}
				if val, err := convertToFloat(columns[1]); err == nil {
					data[fmt.Sprintf("average_objs_osaco_X%d_O2", i)] = val
				}
				if val, err := convertToFloat(columns[3]); err == nil {
					data[fmt.Sprintf("average_objs_osaco_X%d_O4", i)] = val
				}
			}
		}

		// OSACO APTED timeout cases
		for i := 1; i <= 5; i++ {
			if file.Name() == fmt.Sprintf("OSACO_APTED_timeout_X%d.csv", i) {
				if val, err := convertToFloat(columns[0]); err == nil {
					data[fmt.Sprintf("average_objs_osaco_apted_X%d_O1", i)] = val
				}
				if val, err := convertToFloat(columns[1]); err == nil {
					data[fmt.Sprintf("average_objs_osaco_apted_X%d_O2", i)] = val
				}
				if val, err := convertToFloat(columns[3]); err == nil {
					data[fmt.Sprintf("average_objs_osaco_apted_X%d_O4", i)] = val
				}
			}
		}
	}

	return data, testcaseNumbers
}

func averageDataToResult_OSRO(fileName string) string {
	data, testcaseNumbers := getAverageData_OSRO(fileName)

	text := fmt.Sprintf("( testcase numbers: %d ) ", testcaseNumbers)
	text += "--- The experimental results are as follows --- \n"
	text += "The average objective result for the Shortest Path:\n"
	text += fmt.Sprintf("O1: %f O2: %f O3: pass O4: %f \n", data["average_obj_sp_o1"], data["average_obj_sp_o2"], data["average_obj_sp_o4"])
	text += fmt.Sprintf("Computering time: %v ms\n", data["average_time_sp"])
	text += "The average objective result for OSACO (Path-based):\n"
	text += fmt.Sprintf("timeout_X5: O1: %f O2: %f O3: pass O4: %f \n", data["average_objs_osaco_path_X5_O1"], data["average_objs_osaco_path_X5_O2"], data["average_objs_osaco_path_X5_O4"])
	text += fmt.Sprintf("timeout_X4: O1: %f O2: %f O3: pass O4: %f \n", data["average_objs_osaco_path_X4_O1"], data["average_objs_osaco_path_X4_O2"], data["average_objs_osaco_path_X4_O4"])
	text += fmt.Sprintf("timeout_X3: O1: %f O2: %f O3: pass O4: %f \n", data["average_objs_osaco_path_X3_O1"], data["average_objs_osaco_path_X3_O2"], data["average_objs_osaco_path_X3_O4"])
	text += fmt.Sprintf("timeout_X2: O1: %f O2: %f O3: pass O4: %f \n", data["average_objs_osaco_path_X2_O1"], data["average_objs_osaco_path_X2_O2"], data["average_objs_osaco_path_X2_O4"])
	text += fmt.Sprintf("timeout_X1: O1: %f O2: %f O3: pass O4: %f \n", data["average_objs_osaco_path_X1_O1"], data["average_objs_osaco_path_X1_O2"], data["average_objs_osaco_path_X1_O4"])
	text += fmt.Sprintf("Computering time: %v ms\n", data["average_time_osaco_path"])

	return text
}

func getAverageData_OSRO(fileName string) (map[string]float64, int) {
	currentDir, _ := os.Getwd()
	dir := filepath.Join(currentDir + "/data/" + fileName)
	data := make(map[string]float64)

	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("Error reading directory: %v", err)
	}

	testcaseNumbers := 0
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(dir, file.Name())
		fmt.Printf("Reading file: %s\n", filePath)
		csvFile, err := os.Open(filePath)
		if err != nil {
			log.Printf("Error opening file %s: %v", filePath, err)
			continue
		}
		defer csvFile.Close()

		reader := csv.NewReader(csvFile)
		columns := make([][]string, 0)

		// Read columns instead of rows
		for {
			row, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("Error reading CSV from file %s: %v", filePath, err)
				continue
			}
			if len(columns) == 0 {
				columns = make([][]string, len(row))
			}
			for i, value := range row {
				columns[i] = append(columns[i], value)
			}
		}

		convertToFloat := func(row []string) (float64, error) {
			sum := 0.0
			for _, value := range row {
				num, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return 0, fmt.Errorf("error converting value %s to float64: %v", value, err)
				}
				sum += num
			}
			return sum / float64(len(row)), err
		}

		switch file.Name() {
		case "computering_time.csv":
			if val, err := convertToFloat(columns[0]); err == nil {
				data["average_time_sp"] = val
			}
			if val, err := convertToFloat(columns[1]); err == nil {
				data["average_time_osaco_path"] = val
			}
			if val, err := convertToFloat(columns[2]); err == nil {
				testcaseNumbers = int(val)
			}

		case "ShortestPath.csv":
			if val, err := convertToFloat(columns[0]); err == nil {
				data["average_obj_sp_o1"] = val
			}
			if val, err := convertToFloat(columns[1]); err == nil {
				data["average_obj_sp_o2"] = val
			}
			if val, err := convertToFloat(columns[3]); err == nil {
				data["average_obj_sp_o4"] = val
			}
		}

		// OSACO Path timeout cases
		for i := 1; i <= 5; i++ {
			if file.Name() == fmt.Sprintf("OSACO_Path_timeout_X%d.csv", i) {
				if val, err := convertToFloat(columns[0]); err == nil {
					data[fmt.Sprintf("average_objs_osaco_path_X%d_O1", i)] = val
				}
				if val, err := convertToFloat(columns[1]); err == nil {
					data[fmt.Sprintf("average_objs_osaco_path_X%d_O2", i)] = val
				}
				if val, err := convertToFloat(columns[3]); err == nil {
					data[fmt.Sprintf("average_objs_osaco_path_X%d_O4", i)] = val
				}
			}
		}
	}

	return data, testcaseNumbers
}
