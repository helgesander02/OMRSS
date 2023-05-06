package main

import (
    "fmt"
    "flag"

    "src/topology"
    "src/stream" 
    "src/steninertree"   
)

func main() {    
    // Define Parameters
    test_case :=  flag.Int("test_case", 50, "Conducting 50 experiments")
    tsn := flag.Int("tsn", 70, "TSN flow")
    avb := flag.Int("avb", 30, "ANB flow")
    HyperPeriod := flag.Int("HyperPeriod", 6000, "Greatest Common Divisor of Simulated Time LCM.")
    topology_print := flag.Bool("topology_print", false, "Display all topology information.")
    flows_print := flag.Bool("flows_print", false, "Display all Flows information.")
    flag.Parse()

    fmt.Println("The experimental parameters are as follows.")
    fmt.Println("****************************************")
    fmt.Printf("Test Case: %d\n", *test_case)
    fmt.Printf("TSN flow: %d, AVB flow: %d\n", *tsn, *avb)
    fmt.Printf("HyperPeriod: %d\n", *HyperPeriod)
    fmt.Println("****************************************\n")

    // Test-Case
    for ts:=0; ts<*test_case; ts++ { 
        fmt.Printf("TestCase%d \n", ts+1)
        fmt.Println("****************************************")
        // 1. Topology Network "src/topology"
        fmt.Println("Topology Network")
        fmt.Println("----------------------------------------")
        Topology := topology.Generate_Network() 
        fmt.Println("Topology generation completed.")

        if *topology_print { Topology.Print_Topology() }

        // 2. Generate Flows "src/stream"
        fmt.Println("\nGenerate Flows")
        fmt.Println("----------------------------------------")
        Flows := stream.Generate_Flows(*tsn, *avb, *HyperPeriod)

        if *flows_print { 
            Flows.Print_Flows()
            Flows.Print_Flow()
            Flows.Print_Stream()
        }

        // 3. Steiner Tree
        steninertree.SteninerTree()
        fmt.Println("****************************************")
    }
}