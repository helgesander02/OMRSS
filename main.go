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
	test_case     int
	topology_name string
	input_tsn     int
	input_avb     int
	bg_tsn        int
	bg_avb        int
	hyperperiod   int
	bandwidth     float64

	// algo parameters
	plan_name     string
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
	flag.StringVar(&topology_name, "topology_name", "typical_complex", "Topology architecture has typical_complex, typical_simple, ring, layered_ring and industrial.")
	flag.IntVar(&input_tsn, "input_tsn", 35, "Number of TSN input flows.")
	flag.IntVar(&input_avb, "input_avb", 15, "Number of AVB input flows.")
	flag.IntVar(&bg_tsn, "bg_tsn", 35, "Number of TSN bg flows.")
	flag.IntVar(&bg_avb, "bg_avb", 15, "Number of AVB bg flows.")
	flag.IntVar(&hyperperiod, "hyperperiod", 6000, "Greatest Common Divisor of Simulated Time LCM.")
	flag.Float64Var(&bandwidth, "bandwidth", 1e9, "1 Gbps.")

	flag.StringVar(&plan_name, "plan_name", "omaco", "The plan comprises OMACO.")
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

	// Data storage architecture
	Memorizers5 := memorizer.New_Memorizers()
	Memorizer5 := Memorizers5[plan_name]
	Memorizers10 := memorizer.New_Memorizers()
	Memorizer10 := Memorizers10[plan_name]
	Memorizers15 := memorizer.New_Memorizers()
	Memorizer15 := Memorizers15[plan_name]
	Memorizers20 := memorizer.New_Memorizers()
	Memorizer20 := Memorizers20[plan_name]

	// Run Test-Case
	for ts := 0; ts < test_case; ts++ {
		fmt.Printf("\nTestCase%d\n", ts+1)
		fmt.Println("****************************************")
		// 1. Network (a.Topology b.Flows c.Graphs)
		Network := network.Generate_Network(topology_name, bg_tsn, bg_avb, input_tsn, input_avb, hyperperiod, bandwidth)
		if show_network {
			Network.Show_Network()
		}

		// 2. Create new plans (a.OMACO ... )
		Plans5 := plan.New_Plans(Network, osaco_timeout, osaco_K, 0.8) // TODO: To process parameters with a dictionary structure
		Plans10 := plan.New_Plans(Network, osaco_timeout, osaco_K, 0.7)
		Plans15 := plan.New_Plans(Network, osaco_timeout, osaco_K, 0.6)
		Plans20 := plan.New_Plans(Network, osaco_timeout, osaco_K, 0.5)

		// 3. Select plan
		Plan5 := Plans5[plan_name]
		Plan10 := Plans10[plan_name]
		Plan15 := Plans15[plan_name]
		Plan20 := Plans20[plan_name]

		// 4. Initiate plan
		Plan5.Initiate_Plan()
		Plan10.Initiate_Plan()
		Plan15.Initiate_Plan()
		Plan20.Initiate_Plan()
		if show_plan {
			Plan5.Show_Plan()
			Plan10.Show_Plan()
			Plan15.Show_Plan()
			Plan20.Show_Plan()
		}

		// 5. Cumulative quantity
		Memorizer5.M_Cumulative(Plan5)
		Memorizer10.M_Cumulative(Plan10)
		Memorizer15.M_Cumulative(Plan15)
		Memorizer20.M_Cumulative(Plan20)

		fmt.Println("****************************************")
	}

	// 6. Average statistical results
	Memorizer5.M_Average(test_case)
	Memorizer10.M_Average(test_case)
	Memorizer15.M_Average(test_case)
	Memorizer20.M_Average(test_case)

	// 7. Output results
	Memorizer5.M_Output_Results()
	Memorizer10.M_Output_Results()
	Memorizer15.M_Output_Results()
	Memorizer20.M_Output_Results()

	// 8. Save as TXT
	Memorizer5.M_Store_Files(topology_name, test_case, input_tsn, input_avb, osaco_K, 0.8)
	Memorizer10.M_Store_Files(topology_name, test_case, input_tsn, input_avb, osaco_K, 0.7)
	Memorizer15.M_Store_Files(topology_name, test_case, input_tsn, input_avb, osaco_K, 0.6)
	Memorizer20.M_Store_Files(topology_name, test_case, input_tsn, input_avb, osaco_K, 0.5)
	// --------------------------------------------------------------------------------

	// 9. Save as CSV
	if store_data {
		Memorizer5.M_Store_Data()
	}
}
