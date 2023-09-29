package main

import (
	"flag"
	"fmt"

	"src/algo"
	"src/flow"
	"src/network"
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
	show_input_flows := flag.Bool("show_input_flows", false, "Display all input Flows information.")
	show_bg_flows := flag.Bool("show_bg_flows", false, "Display all background Flows information.")
	show_input_routes := flag.Bool("show_input_routes", false, "Display all inputflows routes information.")
	show_bg_routes := flag.Bool("show_bg_routes", false, "Display all backgroundflows routes information.")
	flag.Parse()

	fmt.Println("The experimental parameters are as follows.")
	fmt.Println("****************************************")
	fmt.Printf("Test Case: %d\n", *test_case)
	fmt.Printf("TSN flow: %d, AVB flow: %d\n", *tsn, *avb)
	fmt.Printf("HyperPeriod: %d ns\n", *hyperperiod)
	fmt.Printf("Bandwidth: %f bytes\n", *bandwidth)
	fmt.Println("****************************************\n")

	// Test-Case
	for ts := 0; ts < *test_case; ts++ {
		fmt.Printf("TestCase%d \n", ts+1)
		fmt.Println("****************************************")
		// 1. Topology Network "generate_network.go"
		fmt.Println("Topology Network")
		fmt.Println("----------------------------------------")
		Topology := network.Generate_Network(*bandwidth)
		fmt.Println("Topology generation completed.")

		if *show_topology {
			Topology.Show_Topology()
		}

		// 2. Generate Flows "generate_flows.go"
		fmt.Println("\nGenerate Flows")
		fmt.Println("----------------------------------------")
		Input_flow_set, BG_flow_set := flow.Generate_Flows(len(Topology.Nodes), *tsn, *avb, *hyperperiod)

		if *show_input_flows {
			Input_flow_set.Show_Flows()
			Input_flow_set.Show_Flow()
			Input_flow_set.Show_Stream()
		}

		if *show_bg_flows {
			BG_flow_set.Show_Flows()
			BG_flow_set.Show_Flow()
			BG_flow_set.Show_Stream()
		}

		// 3. Steiner Tree, K Spanning Tree "generate_trees.go"
		fmt.Printf("\nSteiner Tree and %dth Spanning Tree \n", *K)
		fmt.Println("----------------------------------------")
		ktrees_set, Input_tree_set, BG_tree_set := routes.GetRoutes(Topology, Input_flow_set, BG_flow_set, *bandwidth, *K)

		if *show_input_routes {
			Input_tree_set.Show_Trees_Set()
		}

		if *show_bg_routes {
			BG_tree_set.Show_Trees_Set()
		}

		// 4. OSACO "OSACO.go"
		fmt.Printf("\nOSACO \n")
		fmt.Println("----------------------------------------")
		bytes_rate := 1. / (*bandwidth / float64(*hyperperiod))
		visibility := algo.CompVB(ktrees_set, Input_flow_set, BG_flow_set, bytes_rate)
		algo.OSACO(visibility)

		fmt.Println("****************************************")
	}
}
