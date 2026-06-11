package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"src/memorizer"
	"src/network"
	"src/pkg/config"
	"src/pkg/logger"
	"src/pkg/random"
	"src/plan"
)

func main() {
	if err := run(); err != nil {
		logger.Fatalf("Error: %v\n", err)
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

	// Initialise the execution-flow logger as soon as we know the
	// experiment name so every subsequent print is captured on disk.
	if err := logger.Init(cfg.GetExperimentName()); err != nil {
		return fmt.Errorf("failed to initialise logger: %w", err)
	}
	defer func() {
		if cerr := logger.Close(); cerr != nil {
			fmt.Fprintf(os.Stderr, "failed to close logger: %v\n", cerr)
		}
	}()

	// Display configuration summary
	logger.Println("========================================")
	logger.Println("Configuration Summary")
	logger.Println("========================================")
	logger.Printf("Topology: %s\n", cfg.Network.Topology)
	logger.Printf("Algorithm: %s\n", cfg.Algorithm.Name)
	logger.Printf("Test Cases: %d\n", cfg.Experiment.TestCases)
	logger.Println("========================================")
	logger.Println()

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	go func() {
		<-ctx.Done()
		logger.Println("\nReceived shutdown signal, cleaning up...")
	}()

	// Initialize random number generator
	rng := random.New(cfg.Experiment.RandomSeed)
	logger.Printf("Initialized RNG with seed: %d\n", cfg.Experiment.RandomSeed)

	network.FillRNG(rng)
	network.FillTTParams(cfg)
	network.FillCANParams(cfg)
	plan.FillRNG(rng)

	// Network-only mode: stop after generating + showing the network so we
	// can inspect the encapsulation result without paying for plan/memorizer.
	if cfg.Experiment.NetworkOnly {
		return runNetworkOnly(ctx, cfg)
	}

	// Run experiments
	if err := runExperiments(ctx, cfg); err != nil {
		return fmt.Errorf("experiment execution failed: %w", err)
	}

	return nil
}

func runNetworkOnly(ctx context.Context, cfg *config.Config) error {
	logger.Printf(
		"Network-only mode: %d test case(s), algorithm: %s\n",
		cfg.Experiment.TestCases,
		cfg.Algorithm.Name,
	)

	startTime := time.Now()
	for ts := 0; ts < cfg.Experiment.TestCases; ts++ {
		select {
		case <-ctx.Done():
			logger.Println("\nNetwork-only run interrupted by user")
			return ctx.Err()
		default:
		}

		logger.Printf("\nTestCase %d/%d\n", ts+1, cfg.Experiment.TestCases)
		logger.Println("****************************************")

		networkInstance := network.GenerateNetwork(cfg)
		networkInstance.ShowNetwork(cfg.Output.ShowNetwork)

		logger.Println("****************************************")
	}

	logger.Printf("\nFinished network-only run in %v\n", time.Since(startTime))
	return nil
}

func runExperiments(ctx context.Context, cfg *config.Config) error {
	memorizers := memorizer.NewMemorizers()
	memorizer, ok := memorizers[cfg.Algorithm.Name]
	if !ok {
		return fmt.Errorf("unknown algorithm: %s", cfg.Algorithm.Name)
	}

	logger.Printf(
		"Starting %d test cases with algorithm: %s\n",
		cfg.Experiment.TestCases,
		cfg.Algorithm.Name,
	)

	startTime := time.Now()
	for ts := 0; ts < cfg.Experiment.TestCases; ts++ {

		select {
		case <-ctx.Done():
			logger.Println("\nExperiment interrupted by user")
			return ctx.Err()
		default:
		}

		logger.Printf("\nTestCase %d/%d\n", ts+1, cfg.Experiment.TestCases)
		logger.Println("****************************************")

		// 1. Generate Network
		networkInstance := network.GenerateNetwork(cfg)
		if cfg.Output.ShowNetwork {
			networkInstance.ShowNetwork(true)
		}

		// 2. Create Plan
		planInstance := plan.NewPlans(networkInstance, cfg)

		// 3. Initiate Plan
		planInstance.InitiatePlan(cfg)
		if cfg.Output.ShowPlan {
			planInstance.ShowPlan()
		}

		// 4. Accumulate Results
		memorizer.MCumulative(planInstance)

		logger.Println("****************************************")
	}

	elapsed := time.Since(startTime)

	logger.Printf(
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

	logger.Printf("Results saved as: %s\n", experimentName)

	return nil
}
