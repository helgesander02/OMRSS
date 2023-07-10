package main

import (
    "fmt"
    "flag"

    "src/topology"
    "src/flow" 
    "src/SMT"   
)

func main() {    
    // Define Parameters
    test_case :=  flag.Int("test_case", 50, "Conducting 50 experiments")
    tsn := flag.Int("tsn", 70, "Number of TSN flows")
    avb := flag.Int("avb", 30, "Number of AVB flows")
    hyperperiod := flag.Int("hyperperiod", 6000, "Greatest Common Divisor of Simulated Time LCM.")
    bandwidth := flag.Float64("bandwidth", 750000., "1 Gbps ==> bytes/hyperperiod")
    show_topology := flag.Bool("show_topology", false, "Display all topology information.")
    show_flows := flag.Bool("show_flows", false, "Display all Flows information.")
    show_trees := flag.Bool("show_trees", false, "Display all Trees information.")
    flag.Parse()

    fmt.Println("The experimental parameters are as follows.")
    fmt.Println("****************************************")
    fmt.Printf("Test Case: %d\n", *test_case)
    fmt.Printf("TSN flow: %d, AVB flow: %d\n", *tsn, *avb)
    fmt.Printf("HyperPeriod: %d\n", *hyperperiod)
    fmt.Println("****************************************\n")

    // Test-Case
    for ts:=0; ts<*test_case; ts++ { 
        fmt.Printf("TestCase%d \n", ts+1)
        fmt.Println("****************************************")
        // 1. Topology Network "src/topology"
        fmt.Println("Topology Network")
        fmt.Println("----------------------------------------")
        Topology := topology.Generate_Network(*bandwidth) 
        fmt.Println("Topology generation completed.")

        if *show_topology { Topology.Show_Topology() }

        // 2. Generate Flows "src/flow" 
        fmt.Println("\nGenerate Flows")
        fmt.Println("----------------------------------------")
        Flows := stream.Generate_Flows(len(Topology.Nodes), *tsn, *avb, *hyperperiod)

        if *show_flows { 
            Flows.Show_Flows()
            Flows.Show_Flow()
            Flows.Show_Stream()
        }

        // 3. Steiner Tree "src/SMT"
        fmt.Println("\nSteiner Tree")
        fmt.Println("----------------------------------------")
        Trees := SMT.SteninerTree(Topology, Flows, *bandwidth)

        if *show_trees {
            Trees.Show_Trees()
        }

        // 4. K Spanning Tree
        //fmt.Println("\nK Spanning Tree")
        //fmt.Println("----------------------------------------")


        fmt.Println("****************************************")
    }
}