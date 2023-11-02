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
		average_objs_smt   [3][4]float64 // 100ms{o1, o2, o3, o4} 50ms{o1, o2, o3, o4} 10ms{o1, o2, o3, o4}
		average_objs_osaco [3][4]float64 // 100ms{o1, o2, o3, o4} 50ms{o1, o2, o3, o4} 10ms{o1, o2, o3, o4}
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
		objs_smt, objs_osaco := Plan.InitiatePlan(*K, *show_plan, *show_smt, *show_osaco, *timeout)

		for i := 0; i < 3; i++ {
			for j := 0; j < 4; j++ {
				average_objs_smt[i][j] += objs_smt[i][j]
				average_objs_osaco[i][j] += objs_osaco[i][j]
			}
		}

		fmt.Println("****************************************")
	}

	// Result
	for i := 0; i < 3; i++ {
		for j := 0; j < 4; j++ {
			average_objs_smt[i][j] = average_objs_smt[i][j] / float64(*test_case)
			average_objs_osaco[i][j] = average_objs_osaco[i][j] / float64(*test_case)
		}
	}

	fmt.Println()
	fmt.Println("--- The experimental results are as follows ---")
	fmt.Println("The average objective result for the Steiner Tree:")
	fmt.Printf("100ms: O1: %f O2: %f O3: pass O4: %f \n", average_objs_smt[0][0], average_objs_smt[0][1], average_objs_smt[0][3])
	fmt.Printf("50ms: O1: %f O2: %f O3: pass O4: %f \n", average_objs_smt[1][0], average_objs_smt[1][1], average_objs_smt[1][3])
	fmt.Printf("10ms: O1: %f O2: %f O3: pass O4: %f \n", average_objs_smt[2][0], average_objs_smt[2][1], average_objs_smt[2][3])
	fmt.Println("The average objective result for OSACO:")
	fmt.Printf("100ms: O1: %f O2: %f O3: pass O4: %f \n", average_objs_osaco[0][0], average_objs_osaco[0][1], average_objs_osaco[0][3])
	fmt.Printf("50ms: O1: %f O2: %f O3: pass O4: %f \n", average_objs_osaco[1][0], average_objs_osaco[1][1], average_objs_osaco[1][3])
	fmt.Printf("10ms: O1: %f O2: %f O3: pass O4: %f \n", average_objs_osaco[2][0], average_objs_osaco[2][1], average_objs_osaco[2][3])
}
