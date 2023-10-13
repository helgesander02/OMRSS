package main

import (
	"flag"
	"fmt"

	"src/flow"
	"src/network"
	algo1 "src/osaco"
	"src/routes"
)

func main() {
	// Define Parameters
	test_case := flag.Int("test_case", 50, "Conducting 50 experiments")
	tsn := flag.Int("tsn", 70, "Number of TSN flows")
	avb := flag.Int("avb", 30, "Number of AVB flows")
	hyperperiod := flag.Int("hyperperiod", 6000, "Greatest Common Divisor of Simulated Time LCM.")
	bandwidth := flag.Float64("bandwidth", 750000., "1 Gbps ==> bytes/hyperperiod")
	K := flag.Int("K", 5, "finds kth minimum spanning tree")

	show_topology := flag.Bool("show_topology", false, "Display all topology information.")
	show_flows := flag.Bool("show_flows", false, "Display all flows information.")
	show_routes := flag.Bool("show_routes", false, "Display all routes information.")
	flag.Parse()

	fmt.Println("The experimental parameters are as follows.")
	fmt.Println("****************************************")
	fmt.Printf("Test Case: %d\n", *test_case)
	fmt.Printf("TSN flow: %d, AVB flow: %d\n", *tsn, *avb)
	fmt.Printf("HyperPeriod: %d ns\n", *hyperperiod)
	fmt.Printf("Bandwidth: %f bytes\n", *bandwidth)
	fmt.Println("****************************************")

	bytes_rate := 1. / (*bandwidth / float64(*hyperperiod))

	// Test-Case
	for ts := 0; ts < *test_case; ts++ {
		fmt.Printf("TestCase%d \n", ts+1)
		fmt.Println("****************************************")
		// 1. Topology Network "generate_network.go"
		fmt.Println("Topology Network")
		fmt.Println("----------------------------------------")
		Topology := network.Generate_Network(bytes_rate)
		fmt.Println("Topology generation completed.")

		if *show_topology {
			Topology.Show_Topology()
		}

		// 2. Generate Flows "generate_flows.go"
		fmt.Println("\nGenerate Flows")
		fmt.Println("----------------------------------------")
		flow_set := flow.Generate_Flows(len(Topology.Nodes), *tsn, *avb, *hyperperiod)

		if *show_flows {
			flow_set.Show_Flows()
			flow_set.Show_Flow()
			flow_set.Show_Stream()
		}

		// 3. Steiner Tree, K Spanning Tree "generate_trees.go"
		fmt.Printf("\nSteiner Tree and %dth Spanning Tree \n", *K)
		fmt.Println("----------------------------------------")
		ktrees_set := routes.GetRoutes(Topology, flow_set, bytes_rate, *K)

		if *show_routes {
			ktrees_set.Show_kTrees_Set()
		}

		// 4. OSACO "OSACO.go"
		fmt.Printf("\nOSACO \n")
		fmt.Println("----------------------------------------")
		visibility := algo1.CompVB(ktrees_set, flow_set)
		algo1.OSACO(visibility)

		fmt.Println("****************************************")
	}
}
