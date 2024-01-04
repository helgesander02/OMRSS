package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"src/network"
	"src/plan"
)

var (
	test_case     int
	topology_name string
	tsn           int
	avb           int
	hyperperiod   int
	bandwidth     float64

	show_network  bool
	show_topology bool
	show_flows    bool
	show_graphs   bool

	show_plan  bool
	show_smt   bool
	show_mdt   bool
	show_osaco bool
)

func init() {
	// Define Parameters
	flag.IntVar(&test_case, "test_case", 100, "Conducting 50 experiments.")
	flag.StringVar(&topology_name, "topology_name", "typical_complex", "Topology architecture has typical_complex, typical_simple, ring and layered_ring")
	flag.IntVar(&tsn, "tsn", 70, "Number of TSN flows.")
	flag.IntVar(&avb, "avb", 30, "Number of AVB flows.")
	flag.IntVar(&hyperperiod, "hyperperiod", 6000, "Greatest Common Divisor of Simulated Time LCM.")
	flag.Float64Var(&bandwidth, "bandwidth", 1e9, "1 Gbps.")

	flag.BoolVar(&show_network, "show_network", false, "Present all network information comprehensively.")
	flag.BoolVar(&show_topology, "show_topology", false, "Display all topology information.")
	flag.BoolVar(&show_flows, "show_flows", false, "Display all flows information.")
	flag.BoolVar(&show_graphs, "show_graphs", false, "Display all show_graphs information.")

	flag.BoolVar(&show_plan, "show_plan", false, "Provide a comprehensive display of all plan information.")
	flag.BoolVar(&show_smt, "show_smt", false, "Display all Steiner Tree compute information comprehensively.")
	flag.BoolVar(&show_mdt, "show_mdt", false, "Display all Distance Tree compute information comprehensively.")
	flag.BoolVar(&show_osaco, "show_osaco", false, "Display all Osaco compute information comprehensively.")
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
		average_objs_osaco [5][4]float64 // 1000ms{o1, o2, o3, o4} 800ms{o1, o2, o3, o4} 600ms{o1, o2, o3, o4}, 400ms{o1, o2, o3, o4}, 200ms{o1, o2, o3, o4}
	)

	for ts := 0; ts < test_case; ts++ {
		fmt.Printf("\nTestCase%d\n", ts+1)
		fmt.Println("****************************************")
		// 1. Network (1.Topology 2.Flows 3.Graphs)
		Network := network.Generate_Network(topology_name, tsn, avb, hyperperiod, bandwidth, show_topology, show_flows, show_graphs)
		if show_network {
			Network.Show_Network()
		}

		// 2. Create new plan
		Plan := plan.NewPlan(Network)

		// 3. Initiate plan
		obj_smt, obj_mdt, objs_osaco := Plan.InitiatePlan(show_plan, show_smt, show_mdt, show_osaco)

		// 4. Statistical results
		for i := 0; i < 5; i++ {
			for j := 0; j < 4; j++ {
				if i == 1 {
					average_obj_smt[j] += obj_smt[j]
					average_obj_mdt[j] += obj_mdt[j]
				}
				average_objs_osaco[i][j] += objs_osaco[i][j]
			}
		}
		fmt.Println("****************************************")
	}

	// 5. Average statistical results
	for i := 0; i < 5; i++ {
		for j := 0; j < 4; j++ {
			if i == 1 {
				average_obj_smt[j] = average_obj_smt[j] / float64(test_case)
				average_obj_mdt[j] = average_obj_mdt[j] / float64(test_case)
			}
			average_objs_osaco[i][j] = average_objs_osaco[i][j] / float64(test_case)
		}
	}
	// --------------------------------------------------------------------------------

	// Output results
	outputresults(average_obj_smt, average_obj_mdt, average_objs_osaco)
	// Store files
	storefiles(average_obj_smt, average_obj_mdt, average_objs_osaco)
}

func outputresults(average_obj_smt [4]float64, average_obj_mdt [4]float64, average_objs_osaco [5][4]float64) {
	fmt.Println()
	fmt.Println("--- The experimental results are as follows ---")
	fmt.Println("The average objective result for the Steiner Tree:")
	fmt.Printf("O1: %f O2: %f O3: pass O4: %f \n", average_obj_smt[0], average_obj_smt[1], average_obj_smt[3])
	fmt.Println("The average objective result for the MDTC:")
	fmt.Printf("O1: %f O2: %f O3: pass O4: %f \n", average_obj_mdt[0], average_obj_mdt[1], average_obj_mdt[3])
	fmt.Println("The average objective result for OSACO:")
	fmt.Printf("1000ms: O1: %f O2: %f O3: pass O4: %f \n", average_objs_osaco[0][0], average_objs_osaco[0][1], average_objs_osaco[0][3])
	fmt.Printf("800ms: O1: %f O2: %f O3: pass O4: %f \n", average_objs_osaco[1][0], average_objs_osaco[1][1], average_objs_osaco[1][3])
	fmt.Printf("600ms: O1: %f O2: %f O3: pass O4: %f \n", average_objs_osaco[2][0], average_objs_osaco[2][1], average_objs_osaco[2][3])
	fmt.Printf("400ms: O1: %f O2: %f O3: pass O4: %f \n", average_objs_osaco[3][0], average_objs_osaco[3][1], average_objs_osaco[3][3])
	fmt.Printf("200ms: O1: %f O2: %f O3: pass O4: %f \n", average_objs_osaco[4][0], average_objs_osaco[4][1], average_objs_osaco[4][3])
	fmt.Println()
}

func storefiles(average_obj_smt [4]float64, average_obj_mdt [4]float64, average_objs_osaco [5][4]float64) {
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
	text += fmt.Sprintf("1000ms: O1: %f O2: %f O3: pass O4: %f \n", average_objs_osaco[0][0], average_objs_osaco[0][1], average_objs_osaco[0][3])
	text += fmt.Sprintf("800ms: O1: %f O2: %f O3: pass O4: %f \n", average_objs_osaco[1][0], average_objs_osaco[1][1], average_objs_osaco[1][3])
	text += fmt.Sprintf("600ms: O1: %f O2: %f O3: pass O4: %f \n", average_objs_osaco[2][0], average_objs_osaco[2][1], average_objs_osaco[2][3])
	text += fmt.Sprintf("400ms: O1: %f O2: %f O3: pass O4: %f \n", average_objs_osaco[3][0], average_objs_osaco[3][1], average_objs_osaco[3][3])
	text += fmt.Sprintf("200ms: O1: %f O2: %f O3: pass O4: %f \n", average_objs_osaco[4][0], average_objs_osaco[4][1], average_objs_osaco[4][3])

	// Try to open the file in append mode first
	name := fmt.Sprintf("%s_testcase%d_tsn%d_avb%d_output.txt", topology_name, test_case, tsn, avb)
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
