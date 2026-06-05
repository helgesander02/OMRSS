package config

import (
	"fmt"
	"src/pkg/logger"
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
	Flows       FlowsConfig `mapstructure:"flows"`
	Bandwidth   float64     `mapstructure:"bandwidth"`
	ByteRate    float64     // microseconds per byte
}

type FlowsConfig struct {
	// RoutingMode controls the shape of TT flow destinations:
	//   - "tree": flow has multiple destinations (multicast), used by OMACO
	//   - "path": flow has a single destination (unicast),   used by OSRO
	RoutingMode string        `mapstructure:"routing_mode"`
	TSN         TSNFlowConfig `mapstructure:"tsn"`
	AVB         AVBFlowConfig `mapstructure:"avb"`
	CAN         CANFlowConfig `mapstructure:"can"`
}

type TSNFlowConfig struct {
	Input      int                 `mapstructure:"input"`
	Background int                 `mapstructure:"background"`
	Params     TSNFlowParamsConfig `mapstructure:"params"`
}

type AVBFlowConfig struct {
	Input      int                 `mapstructure:"input"`
	Background int                 `mapstructure:"background"`
	Params     AVBFlowParamsConfig `mapstructure:"params"`
}

type TSNFlowParamsConfig struct {
	Periods   []int     `mapstructure:"periods"`   // microseconds (array of possible values)
	DataSizes []float64 `mapstructure:"datasizes"` // bytes (array of possible values)
}

type AVBFlowParamsConfig struct {
	Period    int       `mapstructure:"period"`    // microseconds (fixed value)
	Deadline  int       `mapstructure:"deadline"`  // microseconds (fixed value)
	DataSizes []float64 `mapstructure:"datasizes"` // bytes (array of possible values)
}

type CANFlowConfig struct {
	Nodes       int                 `mapstructure:"nodes"` // Number of CAN nodes randomly selected from topology
	Important   int                 `mapstructure:"important"`
	Unimportant int                 `mapstructure:"unimportant"`
	Params      CANFlowParamsConfig `mapstructure:"params"`
}

type CANFlowParamsConfig struct {
	Important   CANImportantParams   `mapstructure:"important_params"`
	Unimportant CANUnimportantParams `mapstructure:"unimportant_params"`
}

type CANImportantParams struct {
	Period   int     `mapstructure:"period"`   // microseconds
	Deadline int     `mapstructure:"deadline"` // microseconds
	DataSize float64 `mapstructure:"datasize"` // bytes
}

type CANUnimportantParams struct {
	Periods   []int   `mapstructure:"periods"`   // microseconds (array of possible values)
	Deadlines []int   `mapstructure:"deadlines"` // microseconds (array of possible values)
	DataSize  float64 `mapstructure:"datasize"`  // bytes
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
	LogLevel    string `mapstructure:"log_level"`
}

func Load(configName string) (*Config, error) {
	v := viper.New()

	v.SetConfigType("yaml")
	v.AddConfigPath("./pkg/config")

	switch configName {
	case "omaco":
		v.SetConfigName("config.omaco")
	case "osro":
		v.SetConfigName("config.osro")
	default:
		v.SetConfigName("config")
	}

	v.SetEnvPrefix("OMRSS")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		logger.Printf("Config file (%s) not found, using defaults\n", configName)
	} else {
		logger.Printf("Using config file: %s\n", v.ConfigFileUsed())
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	bytesPerUs := (cfg.Network.Bandwidth / 8) * 1e-6                      // bytes per microsecond
	cfg.Network.ByteRate = 1.0 / bytesPerUs                               // microseconds per byte
	cfg.Network.Bandwidth = bytesPerUs * float64(cfg.Network.Hyperperiod) // bytes that can be transmitted in one hyperperiod

	return &cfg, nil
}

func (c *Config) Validate() error {
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

	if c.Algorithm.Name != "omaco" && c.Algorithm.Name != "osro" {
		return fmt.Errorf("invalid algorithm name: %s (valid options: omaco, osro)", c.Algorithm.Name)
	}

	if c.Algorithm.OSACO.PheromoneEvaporation < 0 || c.Algorithm.OSACO.PheromoneEvaporation > 1 {
		return fmt.Errorf("pheromone evaporation must be between 0 and 1, got: %f",
			c.Algorithm.OSACO.PheromoneEvaporation)
	}

	if c.Experiment.TestCases <= 0 {
		return fmt.Errorf("test_cases must be positive, got: %d", c.Experiment.TestCases)
	}

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

	// OSRO is the only algorithm that consumes CAN nodes, so only enforce
	// the positivity check when it is selected.
	if c.Algorithm.Name == "osro" && c.Network.Flows.CAN.Nodes <= 0 {
		return fmt.Errorf("network.flows.can.nodes must be positive, got: %d", c.Network.Flows.CAN.Nodes)
	}

	validRoutingModes := map[string]bool{"tree": true, "path": true}
	if !validRoutingModes[c.Network.Flows.RoutingMode] {
		return fmt.Errorf("invalid network.flows.routing_mode: %q (valid options: tree, path)", c.Network.Flows.RoutingMode)
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
	return fmt.Sprintf("%s_tsn%d_avb%d_%s_K%d_P%.1f_timeout%d",
		c.Network.Topology,
		c.Network.Flows.TSN.Input,
		c.Network.Flows.AVB.Input,
		c.Algorithm.Name,
		c.Algorithm.OSACO.KTrees,
		c.Algorithm.OSACO.PheromoneEvaporation,
		c.Algorithm.OSACO.Timeout,
	)
}
