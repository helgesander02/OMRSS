package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"src/internal/config"
	"src/internal/random"
	"src/memorizer"
	"src/network"
	"src/plan"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error: %v\n", err)
	}
}

func run() error {
	// Load configuration
	configFile := flag.String("config", "omaco", "Path to configuration file (omaco, osro)")
	flag.Parse()

	cfg, err := config.Load(*configFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Display configuration summary
	fmt.Println("========================================")
	fmt.Println("Configuration Summary")
	fmt.Println("========================================")
	fmt.Printf("Topology: %s\n", cfg.Network.Topology)
	fmt.Printf("Algorithm: %s\n", cfg.Algorithm.Name)
	fmt.Printf("Test Cases: %d\n", cfg.Experiment.TestCases)
	fmt.Println("========================================")
	fmt.Println()

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	go func() {
		<-ctx.Done()
		log.Println("\nReceived shutdown signal, cleaning up...")
	}()

	// Initialize random number generator
	rng := random.New(cfg.Experiment.RandomSeed)
	log.Printf("Initialized RNG with seed: %d\n", cfg.Experiment.RandomSeed)

	network.FillRNG(rng)
	network.FillTTParams(cfg)
	network.FillCANParams(cfg)
	plan.FillRNG(rng)

	// Run experiments
	if err := runExperiments(ctx, cfg); err != nil {
		return fmt.Errorf("experiment execution failed: %w", err)
	}

	return nil
}

func runExperiments(ctx context.Context, cfg *config.Config) error {
	memorizers := memorizer.NewMemorizers()
	memorizer, ok := memorizers[cfg.Algorithm.Name]
	if !ok {
		return fmt.Errorf("unknown algorithm: %s", cfg.Algorithm.Name)
	}

	log.Printf(
		"Starting %d test cases with algorithm: %s\n",
		cfg.Experiment.TestCases,
		cfg.Algorithm.Name,
	)

	startTime := time.Now()
	for ts := 0; ts < cfg.Experiment.TestCases; ts++ {

		select {
		case <-ctx.Done():
			log.Println("\nExperiment interrupted by user")
			return ctx.Err()
		default:
		}

		fmt.Printf("\nTestCase %d/%d\n", ts+1, cfg.Experiment.TestCases)
		fmt.Println("****************************************")

		// 1. Generate Network
		networkInstance := network.GenerateNetwork(cfg)
		if cfg.Output.ShowNetwork {
			networkInstance.ShowNetwork()
		}

		// 2. Create Plan
		planInstance := plan.NewPlans(networkInstance, cfg)

		// 3. Initiate Plan
		costSetting := cfg.GetCostArray()
		planInstance.InitiatePlan(costSetting, cfg)
		if cfg.Output.ShowPlan {
			planInstance.ShowPlan()
		}

		// 4. Accumulate Results
		memorizer.MCumulative(planInstance)

		fmt.Println("****************************************")
	}

	elapsed := time.Since(startTime)

	log.Printf(
		"\nCompleted %d test cases in %v\n",
		cfg.Experiment.TestCases,
		elapsed,
	)

	// 5. Calculate Average
	memorizer.MAverage(cfg.Experiment.TestCases)

	// 6. Output Results
	memorizer.MOutputResults()

	// 7. Save Results
	experimentName := cfg.GetExperimentName()

	memorizer.MStoreData(experimentName, cfg.Experiment.TestCases)
	memorizer.MStoreFile(experimentName)

	log.Printf("Results saved as: %s\n", experimentName)

	return nil
}
