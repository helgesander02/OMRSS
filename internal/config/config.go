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
	Flows       FlowsConfig `mapstructure:"flows"`
	Bandwidth   float64     `mapstructure:"bandwidth"`
	ByteRate    float64     // microseconds per byte
}

type FlowsConfig struct {
	TSN TSNFlowConfig `mapstructure:"tsn"`
	AVB AVBFlowConfig `mapstructure:"avb"`
	CAN CANFlowConfig `mapstructure:"can"`
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
	setDefaults(v)

	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.AddConfigPath(".")

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

	bytesPerUs := (cfg.Network.Bandwidth / 8) * 1e-6                      // bytes per microsecond
	cfg.Network.ByteRate = 1.0 / bytesPerUs                               // microseconds per byte
	cfg.Network.Bandwidth = bytesPerUs * float64(cfg.Network.Hyperperiod) // bytes that can be transmitted in one hyperperiod

	return &cfg, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("network.topology", "typical_complex")
	v.SetDefault("network.hyperperiod", 6000)
	v.SetDefault("network.bandwidth", 1e9)
	v.SetDefault("network.flows.tsn.input", 35)
	v.SetDefault("network.flows.tsn.background", 35)
	v.SetDefault("network.flows.avb.input", 15)
	v.SetDefault("network.flows.avb.background", 15)
	v.SetDefault("network.flows.can.important", 5)
	v.SetDefault("network.flows.can.unimportant", 25)

	// TSN flow parameters
	v.SetDefault("network.flows.tsn.params.periods", []int{100, 500, 1000, 1500, 2000})
	v.SetDefault("network.flows.tsn.params.datasizes", []float64{30.0, 40.0, 50.0, 60.0, 70.0, 80.0, 90.0, 100.0})

	// AVB flow parameters
	v.SetDefault("network.flows.avb.params.period", 125)
	v.SetDefault("network.flows.avb.params.deadline", 2000)
	v.SetDefault("network.flows.avb.params.datasizes", []float64{1000.0, 1100.0, 1200.0, 1300.0, 1400.0, 1500.0})

	// CAN important flow parameters
	v.SetDefault("network.flows.can.params.important_params.period", 5000)
	v.SetDefault("network.flows.can.params.important_params.deadline", 5000)
	v.SetDefault("network.flows.can.params.important_params.datasize", 16.0)

	// CAN unimportant flow parameters
	v.SetDefault("network.flows.can.params.unimportant_params.periods", []int{50000, 100000, 150000})
	v.SetDefault("network.flows.can.params.unimportant_params.deadlines", []int{10000, 12000, 14000, 16000, 18000, 20000})
	v.SetDefault("network.flows.can.params.unimportant_params.datasize", 16.0)

	v.SetDefault("algorithm.name", "omaco")
	v.SetDefault("algorithm.osaco.timeout", 200)
	v.SetDefault("algorithm.osaco.k_trees", 5)
	v.SetDefault("algorithm.osaco.pheromone_evaporation", 0.7)
	v.SetDefault("algorithm.osaco.method_number", 0)

	v.SetDefault("schedule.costs.o1", 100000000)
	v.SetDefault("schedule.costs.o2", 100000)
	v.SetDefault("schedule.costs.o3", 0)
	v.SetDefault("schedule.costs.o4", 1)

	v.SetDefault("experiment.test_cases", 100)
	v.SetDefault("experiment.random_seed", 42)

	v.SetDefault("output.show_network", false)
	v.SetDefault("output.show_plan", false)
	v.SetDefault("output.log_level", "info")
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
