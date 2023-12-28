package main

import (
	"flag"
	"fmt"
	"os"
	"src/network"
	"src/plan"
)

func main() {
	// Define Parameters
	test_case := flag.Int("test_case", 100, "Conducting 50 experiments.")
	topology_name := flag.String("topology_name", "typical_complex", "Topology architecture has typical_complex, typical_simple, ring and layered_ring")
	tsn := flag.Int("tsn", 70, "Number of TSN flows.")
	avb := flag.Int("avb", 30, "Number of AVB flows.")
	hyperperiod := flag.Int("hyperperiod", 6000, "Greatest Common Divisor of Simulated Time LCM.")
	bandwidth := flag.Float64("bandwidth", 1e9, "1 Gbps.")

	show_network := flag.Bool("show_network", false, "Present all network information comprehensively.")
	show_topology := flag.Bool("show_topology", false, "Display all topology information.")
	show_flows := flag.Bool("show_flows", false, "Display all flows information.")
	show_graphs := flag.Bool("show_graphs", false, "Display all show_graphs information.")

	show_plan := flag.Bool("show_plan", false, "Provide a comprehensive display of all plan information.")
	show_smt := flag.Bool("show_smt", false, "Display all Steiner Tree compute information comprehensively.")
	show_mdt := flag.Bool("show_mdt", false, "Display all Distance Tree compute information comprehensively.")
	show_osaco := flag.Bool("show_osaco", false, "Display all Osaco compute information comprehensively.")
	flag.Parse()

	fmt.Println("The experimental parameters are as follows.")
	fmt.Println("****************************************")
	fmt.Printf("Topology file name: %s\n", *topology_name)
	fmt.Printf("Test Case: %d\n", *test_case)
	fmt.Printf("TSN flow: %d, AVB flow: %d\n", *tsn, *avb)
	fmt.Printf("HyperPeriod: %d us\n", *hyperperiod)
	fmt.Printf("Bandwidth: 1Gbps ==> %f bytes\n", *bandwidth)
	fmt.Println("****************************************")

	// Test-Case
	var (
		average_obj_smt    [4]float64    // {o1, o2, o3, o4}
		average_obj_mdt    [4]float64    // {o1, o2, o3, o4}
		average_objs_osaco [5][4]float64 // 1000ms{o1, o2, o3, o4} 800ms{o1, o2, o3, o4} 600ms{o1, o2, o3, o4}, 400ms{o1, o2, o3, o4}, 200ms{o1, o2, o3, o4}
	)

	for ts := 0; ts < *test_case; ts++ {
		fmt.Printf("\nTestCase%d \n", ts+1)
		fmt.Println("****************************************")
		// Network (1.Topology 2.Flows 3.Graphs)
		Network := network.Generate_Network(*topology_name, *tsn, *avb, *hyperperiod, *bandwidth, *show_topology, *show_flows, *show_graphs)
		if *show_network {
			Network.Show_Network()
		}

		// Plan (1.SteinerTree  2.DistanceTree  3.OSACO)
		Plan := plan.NewPlan(Network)
		obj_smt, obj_mdt, objs_osaco := Plan.InitiatePlan(*show_plan, *show_smt, *show_mdt, *show_osaco)

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

	// Result
	for i := 0; i < 5; i++ {
		for j := 0; j < 4; j++ {
			if i == 1 {
				average_obj_smt[j] = average_obj_smt[j] / float64(*test_case)
				average_obj_mdt[j] = average_obj_mdt[j] / float64(*test_case)
			}
			average_objs_osaco[i][j] = average_objs_osaco[i][j] / float64(*test_case)
		}
	}
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

	// Store files
	dirName := "result"
	_, err := os.Stat(dirName)
	if err != nil && os.IsNotExist(err) {
		err := os.MkdirAll(dirName, os.ModePerm)
		if err != nil {
			fmt.Println("Folder creation failed：", err)
			return
		}
	}
	err = os.Chdir(dirName)
	if err != nil {
		fmt.Println("Folder switching failed：", err)
		return
	}

	var text string
	text += "--- The experimental results are as follows --- \n"
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

	name := fmt.Sprintf("%s_testcase%d_tsn%d_avb%d_output.txt", *topology_name, *test_case, *tsn, *avb)
	file, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		file, err = os.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Use the file failed：", err)
			return
		}
		fmt.Fprintln(file, text)
		file.Close()

	} else {
		fmt.Fprintln(file, text)
		file.Close()
	}

	err = os.Chdir("..")
	if err != nil {
		fmt.Println("Folder switching failed：", err)
		return
	}
}
