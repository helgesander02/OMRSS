package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Network    NetworkConfig    `mapstructure:"network"`
	Algorithm  AlgorithmConfig  `mapstructure:"algorithm"`
	Schedule   ScheduleConfig   `mapstructure:"schedule"`
	Experiment ExperimentConfig `mapstructure:"experiment"`
	Output     OutputConfig     `mapstructure:"output"`
}

type NetworkConfig struct {
	Topology    string      `mapstructure:"topology"`
	Hyperperiod int         `mapstructure:"hyperperiod"`
	Bandwidth   float64     `mapstructure:"bandwidth"`
	Flows       FlowsConfig `mapstructure:"flows"`
}

type FlowsConfig struct {
	TSN TSNFlowConfig `mapstructure:"tsn"`
	AVB AVBFlowConfig `mapstructure:"avb"`
	CAN CANFlowConfig `mapstructure:"can"`
}

type TSNFlowConfig struct {
	Input      int `mapstructure:"input"`
	Background int `mapstructure:"background"`
}

type AVBFlowConfig struct {
	Input      int `mapstructure:"input"`
	Background int `mapstructure:"background"`
}

type CANFlowConfig struct {
	Important   int `mapstructure:"important"`
	Unimportant int `mapstructure:"unimportant"`
}

type AlgorithmConfig struct {
	Name  string      `mapstructure:"name"`
	OSACO OSACOConfig `mapstructure:"osaco"`
}

type OSACOConfig struct {
	Timeout              int     `mapstructure:"timeout"`
	KTrees               int     `mapstructure:"k_trees"`
	PheromoneEvaporation float64 `mapstructure:"pheromone_evaporation"`
	MethodNumber         int     `mapstructure:"method_number"`
}

type ScheduleConfig struct {
	Costs CostObjectives `mapstructure:"costs"`
}

type CostObjectives struct {
	O1 int `mapstructure:"o1"`
	O2 int `mapstructure:"o2"`
	O3 int `mapstructure:"o3"`
	O4 int `mapstructure:"o4"`
}

type ExperimentConfig struct {
	TestCases  int   `mapstructure:"test_cases"`
	RandomSeed int64 `mapstructure:"random_seed"`
}

type OutputConfig struct {
	ShowNetwork bool   `mapstructure:"show_network"`
	ShowPlan    bool   `mapstructure:"show_plan"`
	ResultsDir  string `mapstructure:"results_dir"`
	LogLevel    string `mapstructure:"log_level"`
}

func Load(configName string) (*Config, error) {
	v := viper.New()

	setDefaults(v)

	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.AddConfigPath(".")

	switch configName {
	case "omaco":
		v.SetConfigName("config.omaco")
	case "osro":
		v.SetConfigName("config.osro")
	case "":
		v.SetConfigName("config")
	default:
		v.SetConfigFile(configName)
	}

	v.SetEnvPrefix("OMRSS")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		fmt.Printf("Config file (%s) not found, using defaults\n", configName)
	} else {
		fmt.Printf("Using config file: %s\n", v.ConfigFileUsed())
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

func setDefaults(v *viper.Viper) {
	// Network defaults
	v.SetDefault("network.topology", "typical_complex")
	v.SetDefault("network.hyperperiod", 6000)
	v.SetDefault("network.bandwidth", 1e9)
	v.SetDefault("network.flows.tsn.input", 35)
	v.SetDefault("network.flows.tsn.background", 35)
	v.SetDefault("network.flows.avb.input", 15)
	v.SetDefault("network.flows.avb.background", 15)
	v.SetDefault("network.flows.can.important", 5)
	v.SetDefault("network.flows.can.unimportant", 25)

	// Algorithm defaults
	v.SetDefault("algorithm.name", "omaco")
	v.SetDefault("algorithm.osaco.timeout", 200)
	v.SetDefault("algorithm.osaco.k_trees", 5)
	v.SetDefault("algorithm.osaco.pheromone_evaporation", 0.7)
	v.SetDefault("algorithm.osaco.method_number", 0)

	// Schedule defaults
	v.SetDefault("schedule.costs.o1", 100000000)
	v.SetDefault("schedule.costs.o2", 100000)
	v.SetDefault("schedule.costs.o3", 0)
	v.SetDefault("schedule.costs.o4", 1)

	// Experiment defaults
	v.SetDefault("experiment.test_cases", 100)
	v.SetDefault("experiment.random_seed", 42)

	// Output defaults
	v.SetDefault("output.show_network", false)
	v.SetDefault("output.show_plan", false)
	v.SetDefault("output.results_dir", "./results")
	v.SetDefault("output.log_level", "info")
}

func (c *Config) Validate() error {
	// Validate topology
	validTopologies := map[string]bool{
		"typical_complex": true,
		"typical_simple":  true,
		"ring":            true,
		"layered_ring":    true,
		"industrial":      true,
	}
	if !validTopologies[c.Network.Topology] {
		return fmt.Errorf("invalid topology: %s (valid options: typical_complex, typical_simple, ring, layered_ring, industrial)", c.Network.Topology)
	}

	// Validate algorithm name
	if c.Algorithm.Name != "omaco" && c.Algorithm.Name != "osro" {
		return fmt.Errorf("invalid algorithm name: %s (valid options: omaco, osro)", c.Algorithm.Name)
	}

	// Validate pheromone evaporation
	if c.Algorithm.OSACO.PheromoneEvaporation < 0 || c.Algorithm.OSACO.PheromoneEvaporation > 1 {
		return fmt.Errorf("pheromone evaporation must be between 0 and 1, got: %f",
			c.Algorithm.OSACO.PheromoneEvaporation)
	}

	// Validate test cases
	if c.Experiment.TestCases <= 0 {
		return fmt.Errorf("test_cases must be positive, got: %d", c.Experiment.TestCases)
	}

	// Validate positive values
	if c.Network.Hyperperiod <= 0 {
		return fmt.Errorf("hyperperiod must be positive, got: %d", c.Network.Hyperperiod)
	}

	if c.Network.Bandwidth <= 0 {
		return fmt.Errorf("bandwidth must be positive, got: %f", c.Network.Bandwidth)
	}

	if c.Algorithm.OSACO.Timeout <= 0 {
		return fmt.Errorf("osaco timeout must be positive, got: %d", c.Algorithm.OSACO.Timeout)
	}

	if c.Algorithm.OSACO.KTrees <= 0 {
		return fmt.Errorf("osaco k_trees must be positive, got: %d", c.Algorithm.OSACO.KTrees)
	}

	return nil
}

func (c *Config) GetCostArray() [4]int {
	return [4]int{
		c.Schedule.Costs.O1,
		c.Schedule.Costs.O2,
		c.Schedule.Costs.O3,
		c.Schedule.Costs.O4,
	}
}

func (c *Config) GetTimeoutDuration() time.Duration {
	return time.Duration(c.Algorithm.OSACO.Timeout) * time.Millisecond
}

func (c *Config) GetExperimentName() string {
	return fmt.Sprintf("%s_tsn%d_avb%d_K%d_P%.1f_timeout%d_O1%d_O2%d",
		c.Network.Topology,
		c.Network.Flows.TSN.Input,
		c.Network.Flows.AVB.Input,
		c.Algorithm.OSACO.KTrees,
		c.Algorithm.OSACO.PheromoneEvaporation,
		c.Algorithm.OSACO.Timeout,
		c.Schedule.Costs.O1,
		c.Schedule.Costs.O2,
	)
}
