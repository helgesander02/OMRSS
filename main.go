package main

import (
	"flag"
	"fmt"
	"src/network"
	"src/plan"
)

func main() {
	// Define Parameters
	test_case := flag.Int("test_case", 50, "Conducting 50 experiments.")
	tsn := flag.Int("tsn", 70, "Number of TSN flows.")
	avb := flag.Int("avb", 30, "Number of AVB flows.")
	hyperperiod := flag.Int("hyperperiod", 6000, "Greatest Common Divisor of Simulated Time LCM.")
	bandwidth := flag.Float64("bandwidth", 1e9, "1 Gbps.")
	K := flag.Int("K", 5, "Finds kth minimum spanning tree.")
	timeout := flag.Int("timeout", 100, "The timeout of each run is set as 10~100 ms.")

	show_network := flag.Bool("show_network", false, "Present all network information comprehensively.")
	show_topology := flag.Bool("show_topology", false, "Display all topology information.")
	show_flows := flag.Bool("show_flows", false, "Display all flows information.")
	show_graphs := flag.Bool("show_graphs", false, "Display all show_graphs information.")

	show_plan := flag.Bool("show_plan", false, "Provide a comprehensive display of all plan information.")
	show_smt := flag.Bool("show_smt", false, "Display all Steiner Tree compute information comprehensively.")
	show_osaco := flag.Bool("show_osaco", false, "Display all Osaco compute information comprehensively.")
	flag.Parse()

	fmt.Println("The experimental parameters are as follows.")
	fmt.Println("****************************************")
	fmt.Printf("Test Case: %d\n", *test_case)
	fmt.Printf("TSN flow: %d, AVB flow: %d\n", *tsn, *avb)
	fmt.Printf("HyperPeriod: %d us\n", *hyperperiod)
	fmt.Printf("Bandwidth: 1Gbps ==> %f bytes\n", *bandwidth)
	fmt.Println("****************************************")

	// Test-Case
	var (
		average_obj_SMT   [4]float64 = [4]float64{0.0, 0.0, 0.0, 0.0}
		average_obj_osaco [4]float64 = [4]float64{0.0, 0.0, 0.0, 0.0}
	)

	for ts := 0; ts < *test_case; ts++ {
		fmt.Printf("\nTestCase%d \n", ts+1)
		fmt.Println("****************************************")
		// Network (1.Topology 2.Flows 3.Graphs)
		Network := network.Generate_Network(*tsn, *avb, *hyperperiod, *bandwidth, *show_topology, *show_flows, *show_graphs)

		if *show_network {
			Network.Show_Network()
		}

		// Plan (1.SteinerTree 2.OSACO 3. ...)
		Plan := plan.NewPlan(Network)
		obj_SMT, obj_osaco := Plan.InitiatePlan(*K, *show_plan, *show_smt, *show_osaco, *timeout)
		for i := 0; i < 4; i++ {
			average_obj_SMT[i] += obj_SMT[i]
			average_obj_osaco[i] += obj_osaco[i]
		}

		fmt.Println("****************************************")
	}

	// Result
	for i := 0; i < 4; i++ {
		average_obj_SMT[i] = average_obj_SMT[i] / float64(*test_case)
		average_obj_osaco[i] = average_obj_osaco[i] / float64(*test_case)
	}
	fmt.Println()
	fmt.Println("--- The experimental results are as follows ---")
	fmt.Printf("The average objective result for the Steiner Tree: O1: %f O2: %f O3: pass O4: %f \n", average_obj_SMT[0], average_obj_SMT[1], average_obj_SMT[3])
	fmt.Printf("The average objective result for OSACO: O1: %f O2: %f O3: pass O4: %f \n", average_obj_osaco[0], average_obj_osaco[1], average_obj_osaco[3])
}
