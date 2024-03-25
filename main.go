package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"src/network"
	"src/plan"
)

var (
	// network parameters
	test_case     int
	topology_name string
	tsn           int
	avb           int
	hyperperiod   int
	bandwidth     float64

	// algo parameters
	osaco_timeout int
	osaco_K       int
	osaco_P       float64

	// show parameters
	show_network bool
	show_plan    bool
	store_data   bool
)

func init() {
	// Define Parameters
	flag.IntVar(&test_case, "test_case", 100, "Conducting 50 experiments.")
	flag.StringVar(&topology_name, "topology_name", "typical_complex", "Topology architecture has typical_complex, typical_simple, ring and layered_ring.")
	flag.IntVar(&tsn, "tsn", 70, "Number of TSN flows.")
	flag.IntVar(&avb, "avb", 30, "Number of AVB flows.")
	flag.IntVar(&hyperperiod, "hyperperiod", 6000, "Greatest Common Divisor of Simulated Time LCM.")
	flag.Float64Var(&bandwidth, "bandwidth", 1e9, "1 Gbps.")
	flag.IntVar(&osaco_timeout, "osaco_timeout", 200, "Timeout in milliseconds")
	flag.IntVar(&osaco_K, "osaco_K", 5, "Select K trees with different weights.")
	flag.Float64Var(&osaco_P, "osaco_P", 0.6, "The pheromone value of each routing path starts to evaporate, where P is the given evaporation coefficient. (0 <= p <= 1)")

	flag.BoolVar(&show_network, "show_network", false, "Present all network information comprehensively.")
	flag.BoolVar(&show_plan, "show_plan", false, "Provide a comprehensive display of all plan information.")
	flag.BoolVar(&store_data, "store_data", false, "Store all statistical data")
}

func main() {
	// OMRSS
	// --------------------------------------------------------------------------------
	// Process the arguments on the command line
	flag.Parse()

	// Run Test-Case
	var (
		average_obj_smt    [4]float64    // {o1, o2, o3, o4}
		average_obj_mdt    [4]float64    // {o1, o2, o3, o4}
		average_objs_osaco [5][4]float64 // 200ms{o1, o2, o3, o4} 400ms{o1, o2, o3, o4} 600ms{o1, o2, o3, o4}, 800ms{o1, o2, o3, o4}, 1000ms{o1, o2, o3, o4}
	)

	for ts := 0; ts < test_case; ts++ {
		fmt.Printf("\nTestCase%d\n", ts+1)
		fmt.Println("****************************************")
		// 1. Network (1.Topology 2.Flows 3.Graphs)
		Network := network.Generate_Network(topology_name, tsn, avb, hyperperiod, bandwidth)
		if show_network {
			Network.Show_Network()
		}

		// 2. Create new plan
		Plan := plan.NewPlan(Network, osaco_timeout, osaco_K, osaco_P) // TODO: To process parameters with a dictionary structure

		// 3. Initiate plan
		Plan.InitiatePlan()
		if show_plan {
			Plan.Show_Plan()
		}

		for i := 0; i < 5; i++ {
			for j := 0; j < 4; j++ {
				if i == 0 {
					average_obj_smt[j] += Plan.SMT.Objs_smt[j]
					average_obj_mdt[j] += Plan.MDTC.Objs_mdtc[j]
				}
				average_objs_osaco[i][j] += Plan.OSACO.Objs_osaco[i][j]
			}
		}
		fmt.Println("****************************************")
	}

	// 4. Average statistical results
	for i := 0; i < 5; i++ {
		for j := 0; j < 4; j++ {
			if i == 0 {
				average_obj_smt[j] = average_obj_smt[j] / float64(test_case)
				average_obj_mdt[j] = average_obj_mdt[j] / float64(test_case)
			}
			average_objs_osaco[i][j] = average_objs_osaco[i][j] / float64(test_case)
		}
	}
	// 5. Statistical results
	if store_data {
		// Save as CSV
		Store_Data(average_obj_smt, average_obj_mdt, average_objs_osaco)
	}
	// --------------------------------------------------------------------------------

	// Output results
	Output_Results(average_obj_smt, average_obj_mdt, average_objs_osaco)
	// Save as TXT
	Store_Files(average_obj_smt, average_obj_mdt, average_objs_osaco)
}

func Output_Results(average_obj_smt [4]float64, average_obj_mdt [4]float64, average_objs_osaco [5][4]float64) {
	fmt.Println()
	fmt.Println("--- The experimental results are as follows ---")
	fmt.Println("The average objective result for the Steiner Tree:")
	fmt.Printf("O1: %f O2: %f O3: pass O4: %f \n", average_obj_smt[0], average_obj_smt[1], average_obj_smt[3])
	fmt.Println("The average objective result for the MDTC:")
	fmt.Printf("O1: %f O2: %f O3: pass O4: %f \n", average_obj_mdt[0], average_obj_mdt[1], average_obj_mdt[3])
	fmt.Println("The average objective result for OSACO:")
	fmt.Printf("1000ms: O1: %f O2: %f O3: pass O4: %f \n", average_objs_osaco[4][0], average_objs_osaco[4][1], average_objs_osaco[4][3])
	fmt.Printf("800ms: O1: %f O2: %f O3: pass O4: %f \n", average_objs_osaco[3][0], average_objs_osaco[3][1], average_objs_osaco[3][3])
	fmt.Printf("600ms: O1: %f O2: %f O3: pass O4: %f \n", average_objs_osaco[2][0], average_objs_osaco[2][1], average_objs_osaco[2][3])
	fmt.Printf("400ms: O1: %f O2: %f O3: pass O4: %f \n", average_objs_osaco[1][0], average_objs_osaco[1][1], average_objs_osaco[1][3])
	fmt.Printf("200ms: O1: %f O2: %f O3: pass O4: %f \n", average_objs_osaco[0][0], average_objs_osaco[0][1], average_objs_osaco[0][3])
	fmt.Println()
}

func Store_Files(average_obj_smt [4]float64, average_obj_mdt [4]float64, average_objs_osaco [5][4]float64) {
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
	text += fmt.Sprintf("O1: %f O2: %f O3: pass O4: %f \n", average_obj_smt[0], average_obj_smt[1], average_obj_smt[3])
	text += "The average objective result for the MDTC:\n"
	text += fmt.Sprintf("O1: %f O2: %f O3: pass O4: %f \n", average_obj_mdt[0], average_obj_mdt[1], average_obj_mdt[3])
	text += "The average objective result for OSACO:\n"
	text += fmt.Sprintf("1000ms: O1: %f O2: %f O3: pass O4: %f \n", average_objs_osaco[4][0], average_objs_osaco[4][1], average_objs_osaco[4][3])
	text += fmt.Sprintf("800ms: O1: %f O2: %f O3: pass O4: %f \n", average_objs_osaco[3][0], average_objs_osaco[3][1], average_objs_osaco[3][3])
	text += fmt.Sprintf("600ms: O1: %f O2: %f O3: pass O4: %f \n", average_objs_osaco[2][0], average_objs_osaco[2][1], average_objs_osaco[2][3])
	text += fmt.Sprintf("400ms: O1: %f O2: %f O3: pass O4: %f \n", average_objs_osaco[1][0], average_objs_osaco[1][1], average_objs_osaco[1][3])
	text += fmt.Sprintf("200ms: O1: %f O2: %f O3: pass O4: %f \n", average_objs_osaco[0][0], average_objs_osaco[0][1], average_objs_osaco[0][3])

	// Try to open the file in append mode first
	name := fmt.Sprintf("%s_testcase%d_tsn%d_avb%d_osacok%d_p%f_output.txt", topology_name, test_case, tsn, avb, osaco_K, osaco_P)
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

func Store_Data(average_obj_smt [4]float64, average_obj_mdt [4]float64, average_objs_osaco [5][4]float64) {
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
	StoreCSV("SMT.csv", average_obj_smt)
	log.Printf("Storing MTDC data to CSV...\n")
	StoreCSV("MTDC.csv", average_obj_mdt)
	log.Printf("Storing OSACO data to CSV...\n")
	StoreCSV("OSACO1000ms.csv", average_objs_osaco[4])
	StoreCSV("OSACO800ms.csv", average_objs_osaco[3])
	StoreCSV("OSACO600ms.csv", average_objs_osaco[2])
	StoreCSV("OSACO400ms.csv", average_objs_osaco[1])
	StoreCSV("OSACO200ms.csv", average_objs_osaco[0])

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
