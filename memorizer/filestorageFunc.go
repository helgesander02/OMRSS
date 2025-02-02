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

func (OC *OMACO_Memorizer) M_Store_File(file_name string) {
	text := average_data_to_result(file_name)

	dirName := "result"
	createFolder(dirName)
	switchWorkingPath(dirName)

	// Try to open the file in append mode first
	txt_name := file_name + ".txt"
	log.Printf("Opening file %s in append mode\n", txt_name)
	file, err := os.OpenFile(txt_name, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		// If it fails, create the file
		log.Printf("File not found, creating file %s\n", txt_name)
		file, err = os.OpenFile(txt_name, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to create or open file %s: %v\n", txt_name, err)
		}
	}
	defer file.Close()
	log.Printf("Writing text to file %s\n", txt_name)
	fmt.Fprintln(file, text)

	switchWorkingPath("..")
}

func average_data_to_result(file_name string) string {
	data, testcase_numbers := get_average_data(file_name)

	text := fmt.Sprintf("( testcase numbers: %d ) ", testcase_numbers)
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
	text += "The average objective result for OSACO_IAS:\n"
	text += fmt.Sprintf("timeout_X5: O1: %f O2: %f O3: pass O4: %f \n", data["average_objs_osaco_ias_X5_O1"], data["average_objs_osaco_ias_X5_O2"], data["average_objs_osaco_ias_X5_O4"])
	text += fmt.Sprintf("timeout_X4: O1: %f O2: %f O3: pass O4: %f \n", data["average_objs_osaco_ias_X4_O1"], data["average_objs_osaco_ias_X4_O2"], data["average_objs_osaco_ias_X4_O4"])
	text += fmt.Sprintf("timeout_X3: O1: %f O2: %f O3: pass O4: %f \n", data["average_objs_osaco_ias_X3_O1"], data["average_objs_osaco_ias_X3_O2"], data["average_objs_osaco_ias_X3_O4"])
	text += fmt.Sprintf("timeout_X2: O1: %f O2: %f O3: pass O4: %f \n", data["average_objs_osaco_ias_X2_O1"], data["average_objs_osaco_ias_X2_O2"], data["average_objs_osaco_ias_X2_O4"])
	text += fmt.Sprintf("timeout_X1: O1: %f O2: %f O3: pass O4: %f \n", data["average_objs_osaco_ias_X1_O1"], data["average_objs_osaco_ias_X1_O2"], data["average_objs_osaco_ias_X1_O4"])
	text += fmt.Sprintf("Computering time: %v ms\n", data["average_time_osaco_ias"])

	return text
}

func get_average_data(file_name string) (map[string]float64, int) {
	currentDir, _ := os.Getwd()
	dir := filepath.Join(currentDir + "/data/" + file_name)
	data := make(map[string]float64)

	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("Error reading directory: %v", err)
	}

	testcase_numbers := 0
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
				data["average_time_osaco_ias"] = val
			}
			if val, err := convertToFloat(columns[3]); err == nil {
				testcase_numbers = int(val)
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

		// OSACO IAS timeout cases
		for i := 1; i <= 5; i++ {
			if file.Name() == fmt.Sprintf("OSACO_IAS_timeout_X%d.csv", i) {
				if val, err := convertToFloat(columns[0]); err == nil {
					data[fmt.Sprintf("average_objs_osaco_ias_X%d_O1", i)] = val
				}
				if val, err := convertToFloat(columns[1]); err == nil {
					data[fmt.Sprintf("average_objs_osaco_ias_X%d_O2", i)] = val
				}
				if val, err := convertToFloat(columns[3]); err == nil {
					data[fmt.Sprintf("average_objs_osaco_ias_X%d_O4", i)] = val
				}
			}
		}
	}

	return data, testcase_numbers
}

//func (mm2 *Memorizer2) M_Store_Files(test_case int) {
//
//}

//func (mm3 *Memorizer2) M_Store_Files(test_case int) {
//
//}
