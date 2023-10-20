package main

import (
	"flag"
	"fmt"
	"src/network"
	"src/plan"
)

func main() {
	// Define Parameters
	test_case := flag.Int("test_case", 50, "Conducting 50 experiments")
	tsn := flag.Int("tsn", 70, "Number of TSN flows")
	avb := flag.Int("avb", 30, "Number of AVB flows")
	hyperperiod := flag.Int("hyperperiod", 6000, "Greatest Common Divisor of Simulated Time LCM.")
	bandwidth := flag.Float64("bandwidth", 1e9, "1 Gbps")
	K := flag.Int("K", 5, "finds kth minimum spanning tree")

	show_network := flag.Bool("show_network", false, "Present all network information comprehensively.")
	show_topology := flag.Bool("show_topology", false, "Display all topology information.")
	show_flows := flag.Bool("show_flows", false, "Display all flows information.")
	show_graphs := flag.Bool("show_graphs", false, "Display all show_graphs information.")

	show_plan := flag.Bool("show_plan", false, "Provide a comprehensive display of all plan information.")
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
	for ts := 0; ts < *test_case; ts++ {
		fmt.Printf("TestCase%d \n", ts+1)
		fmt.Println("****************************************")
		// Network (1. Topology 2. Flows 3. Graphs)
		Network := network.Generate_Network(*tsn, *avb, *hyperperiod, *bandwidth, *show_topology, *show_flows, *show_graphs)

		if *show_network {
			Network.Show_Network()
		}

		Plan := plan.NewPlan(Network)
		// Plan (1. SteinerTree() 2. OSACO(*K) 3. ...)
		Plan.InitiatePlan(*K, *show_plan, *show_osaco)

		fmt.Println("****************************************")
	}
}
