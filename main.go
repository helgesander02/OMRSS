package main

import (
	"flag"
	"fmt"
	"src/memorizer"
	"src/network"
	"src/plan"
)

var (
	// network parameters
	test_case       int
	topology_name   string
	input_tsn       int
	input_avb       int
	bg_tsn          int
	bg_avb          int
	important_can   int
	unimportant_can int
	hyperperiod     int
	bandwidth       float64

	// algo parameters
	plan_name     string
	osaco_timeout int
	osaco_K       int
	osaco_P       float64

	// schedule parameters
	o1_cost int
	o2_cost int
	o3_cost int
	o4_cost int

	// show parameters
	show_network bool
	show_plan    bool
)

func init() {
	// Define Parameters
	flag.IntVar(&test_case, "test_case", 100, "Conducting 50 experiments.")
	flag.StringVar(&topology_name, "topology_name", "typical_complex", "Topology architecture has typical_complex, typical_simple, ring, layered_ring and industrial.")
	flag.IntVar(&input_tsn, "input_tsn", 35, "Number of TSN input flows.")
	flag.IntVar(&input_avb, "input_avb", 15, "Number of AVB input flows.")
	flag.IntVar(&bg_tsn, "bg_tsn", 35, "Number of TSN bg flows.")
	flag.IntVar(&bg_avb, "bg_avb", 15, "Number of AVB bg flows.")
	flag.IntVar(&important_can, "important_can", 5, "Number of CAN important flows.")
	flag.IntVar(&unimportant_can, "unimportant_can", 25, "Number of CAN unimportant flows.")
	flag.IntVar(&hyperperiod, "hyperperiod", 6000, "Greatest Common Divisor of Simulated Time LCM.")
	flag.Float64Var(&bandwidth, "bandwidth", 1e9, "1 Gbps.")

	flag.StringVar(&plan_name, "plan_name", "omaco", "The plan comprises omaco and osro.")
	flag.IntVar(&osaco_timeout, "osaco_timeout", 200, "Timeout in milliseconds")
	flag.IntVar(&osaco_K, "osaco_K", 5, "Select K trees with different weights.")
	flag.Float64Var(&osaco_P, "osaco_P", 0.7, "The pheromone value of each routing path starts to evaporate, where P is the given evaporation coefficient. (0 <= p <= 1)")

	flag.IntVar(&o1_cost, "o1_cost", 100000000, "O1 cost")
	flag.IntVar(&o2_cost, "o2_cost", 100000, "O2 cost")
	flag.IntVar(&o3_cost, "o3_cost", 0, "O3 cost")
	flag.IntVar(&o4_cost, "o4_cost", 1, "O4 cost")

	flag.BoolVar(&show_network, "show_network", false, "Present all network information comprehensively.")
	flag.BoolVar(&show_plan, "show_plan", false, "Provide a comprehensive display of all plan information.")
}

func main() {
	// OMRSS
	// --------------------------------------------------------------------------------
	// Process the arguments on the command line
	flag.Parse()

	// Data storage architecture
	// ----------------------------------------------
	Memorizers := memorizer.New_Memorizers()
	Memorizer := Memorizers[plan_name]

	// Run Test-Case
	for ts := 0; ts < test_case; ts++ {
		fmt.Printf("\nTestCase%d\n", ts+1)
		fmt.Println("****************************************")
		// 1. Network (a.OMACO ... )
		Networks := network.New_Networks(topology_name, bg_tsn, bg_avb, input_tsn, input_avb, important_can, unimportant_can, hyperperiod, bandwidth)
		Network := Networks[plan_name]
		Network.Generate_Network() // Generate a.Topology b.Flows c.Graphs
		if show_network {
			Network.Show_Network()
		}

		// 2. Create new plans (a.OMACO ... )
		// -----------------------------------------------------------
		Plans := plan.New_Plans(Network, osaco_timeout, osaco_K, osaco_P)
		Plan := Plans[plan_name]

		// 3. Initiate plan
		// ----------------------------------------------------------
		cost_setting1 := [4]int{o1_cost, o2_cost, o3_cost, o4_cost}
		Plan.Initiate_Plan(cost_setting1)
		if show_plan {
			Plan.Show_Plan()
		}

		// 4. Cumulative quantity
		// ------------------------------------------
		Memorizer.M_Cumulative(Plan)
		fmt.Println("****************************************")
	}

	// 5. Average statistical results
	// ---------------------------------
	Memorizer.M_Average(test_case)

	// 6. Output results
	// -----------------------------------
	Memorizer.M_Output_Results()

	// 7. Save as CSV
	// --------------------------------------------------------------------------------------------------------------------------------------------------
	name := fmt.Sprintf("%s_tsn%d_avb%d_K%d_P%.1f_timeout%d_O1%d_O2%d", topology_name, input_tsn, input_avb, osaco_K, osaco_P, osaco_timeout, o1_cost, o2_cost)
	Memorizer.M_Store_Data(name, test_case)

	// 8. Save as TXT
	// ------------------------------------
	Memorizer.M_Store_File(name)
	// --------------------------------------------------------------------------------
}
