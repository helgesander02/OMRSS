package main

import (
	"context"
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
	"src/network/flow/can"
	"src/network/flow/tt"
	"src/network/topology"
	"src/plan"
	"src/plan/algo"
	"src/plan/routes"

	"github.com/spf13/pflag"
)

var (
	configFile string
)

func init() {
	pflag.StringVarP(&configFile, "config", "c", "", "Path to configuration file")

	// Allow overriding config values via CLI flags
	pflag.String("topology", "", "Override topology name")
	pflag.Int("test-cases", 0, "Override number of test cases")
	pflag.Int("tsn-input", 0, "Override number of tsn input streams")
	pflag.Int("avb-input", 0, "Override number of avb input streams")
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error: %v\n", err)
	}
}

func run() error {
	pflag.Parse()

	// Load configuration
	cfg, err := config.Load(configFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	applyCliOverrides(cfg) // pflaf > config

	// Display configuration summary
	fmt.Println("========================================")
	fmt.Println("OMRSS Configuration Summary")
	fmt.Println("========================================")
	fmt.Printf("Topology: %s\n", cfg.Network.Topology)
	fmt.Printf("Algorithm: %s\n", cfg.Algorithm.Name)
	fmt.Printf("Test Cases: %d\n", cfg.Experiment.TestCases)
	fmt.Printf("Random Seed: %d\n", cfg.Experiment.RandomSeed)
	fmt.Println("========================================")
	fmt.Println()

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("\nReceived shutdown signal, cleaning up...")
		cancel()
	}()

	// Initialize random number generator with seed
	rng := random.New(cfg.Experiment.RandomSeed)
	log.Printf("Initialized RNG with seed: %d\n", cfg.Experiment.RandomSeed)

	// Set RNG for all packages that need randomness
	tt.SetRNG(rng)
	can.SetRNG(rng)
	topology.SetRNG(rng)
	routes.SetRNG(rng)
	algo.SetRNG(rng)

	// Create results directory
	if err := os.MkdirAll(cfg.Output.ResultsDir, 0755); err != nil {
		return fmt.Errorf("failed to create results directory: %w", err)
	}

	// Run experiments
	if err := runExperiments(ctx, cfg); err != nil {
		return fmt.Errorf("experiment execution failed: %w", err)
	}

	return nil
}

func applyCliOverrides(cfg *config.Config) {
	if pflag.Lookup("topology").Changed {
		cfg.Network.Topology = pflag.Lookup("topology").Value.String()
	}
	if pflag.Lookup("test-cases").Changed {
		if val, err := pflag.CommandLine.GetInt("test-cases"); err == nil && val > 0 {
			cfg.Experiment.TestCases = val
		}
	}
	if pflag.Lookup("tsn-input").Changed {
		if val, err := pflag.CommandLine.GetInt("tsn-input"); err == nil && val > 0 {
			cfg.Network.Flows.TSN.Input = val
		}
	}
	if pflag.Lookup("avb-input").Changed {
		if val, err := pflag.CommandLine.GetInt("avb-input"); err == nil && val > 0 {
			cfg.Network.Flows.AVB.Input = val
		}
	}
}

func runExperiments(ctx context.Context, cfg *config.Config) error {
	// Initialize memorizer for result tracking
	memorizers := memorizer.NewMemorizers()
	memorizer, ok := memorizers[cfg.Algorithm.Name]
	if !ok {
		return fmt.Errorf("unknown algorithm: %s", cfg.Algorithm.Name)
	}

	log.Printf("Starting %d test cases with algorithm: %s\n",
		cfg.Experiment.TestCases, cfg.Algorithm.Name)

	startTime := time.Now()

	// Run test cases
	for ts := 0; ts < cfg.Experiment.TestCases; ts++ {
		// Check for cancellation
		select {
		case <-ctx.Done():
			log.Println("\nExperiment interrupted by user")
			return ctx.Err()
		default:
		}

		fmt.Printf("\nTestCase %d/%d\n", ts+1, cfg.Experiment.TestCases)
		fmt.Println("****************************************")

		// 1. Generate Network
		Networks := network.NewNetworks(
			cfg.Network.Topology,
			cfg.Network.Flows.TSN.Background,
			cfg.Network.Flows.AVB.Background,
			cfg.Network.Flows.TSN.Input,
			cfg.Network.Flows.AVB.Input,
			cfg.Network.Flows.CAN.Important,
			cfg.Network.Flows.CAN.Unimportant,
			cfg.Network.Hyperperiod,
			cfg.Network.Bandwidth,
		)

		Network := Networks[cfg.Algorithm.Name]
		Network.GenerateNetwork()

		if cfg.Output.ShowNetwork {
			Network.ShowNetwork()
		}

		// 2. Create Plan
		Plan := plan.NewPlans(
			cfg.Algorithm.Name,
			Network,
			cfg.Algorithm.OSACO.Timeout,
			cfg.Algorithm.OSACO.KTrees,
			cfg.Algorithm.OSACO.PheromoneEvaporation,
		)

		// 3. Initiate Plan
		costSetting := cfg.GetCostArray()
		Plan.InitiatePlan(costSetting)

		if cfg.Output.ShowPlan {
			Plan.ShowPlan()
		}

		// 4. Accumulate Results
		memorizer.MCumulative(Plan)
		fmt.Println("****************************************")
	}

	elapsed := time.Since(startTime)
	log.Printf("\nCompleted %d test cases in %v\n", cfg.Experiment.TestCases, elapsed)

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
