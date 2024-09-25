package common

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	serverAddress  = "127.0.0.1:8080"
	maxConnections = 100
	maxMessageSize = "4KB"
	idleTimeout    = "5m"
	inMemoryEngine = "in_memory"
	loggingLevel   = "debug"
)

// EngineConfig - engine config
type EngineConfig struct {
	Type string `yaml:"type"`
}

// NetworkConfig - network config
type NetworkConfig struct {
	Address        string `yaml:"address"`
	MaxConnections int    `yaml:"max_connections"`
	MaxMsgSize     string `yaml:"max_message_size"`
	IdleTimeout    string `yaml:"idle_timeout"`
}

// LoggingConfig - logging config
type LoggingConfig struct {
	Level  string `yaml:"level"`
	Output string `yaml:"output"`
}

// Config - server config
type Config struct {
	Engine  *EngineConfig  `yaml:"engine"`
	Network *NetworkConfig `yaml:"network"`
	Logging *LoggingConfig `yaml:"logging"`
}

func ParseConfig(reader io.Reader) (*Config, error) {
	if reader == nil {
		return nil, errors.New("nil reader provided")
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	if config.Engine == nil {
		config.Engine = &EngineConfig{}
	}

	if config.Network == nil {
		config.Network = &NetworkConfig{}
	}

	if config.Logging == nil {
		config.Logging = &LoggingConfig{}
	}

	if config.Network.Address == "" {
		config.Network.Address = serverAddress
	}

	if config.Network.MaxConnections == 0 {
		config.Network.MaxConnections = maxConnections
	}

	if config.Network.MaxMsgSize == "" {
		config.Network.MaxMsgSize = maxMessageSize
	}

	if config.Network.IdleTimeout == "" {
		config.Network.IdleTimeout = idleTimeout
	}

	if config.Engine.Type == "" {
		config.Engine.Type = inMemoryEngine
	}

	if config.Logging.Level == "" {
		config.Logging.Level = loggingLevel
	}

	return &config, nil
}

func ParseBufSize(cfgSize string) (int, error) {
	cfgSize = strings.TrimSpace(strings.ToUpper(cfgSize))
	units := map[string]int{
		"B":  1,
		"KB": 1024,
		"MB": 1024 * 1024,
		"GB": 1024 * 1024 * 1024,
	}

	i := 0
	for i < len(cfgSize) && (cfgSize[i] >= '0' && cfgSize[i] <= '9') {
		i++
	}

	numPart := cfgSize[:i]
	unitPart := cfgSize[i:]
	num, err := strconv.Atoi(numPart)
	if err != nil {
		return 0, fmt.Errorf("invalid buffer size provided: %s", cfgSize)
	}

	multiplier, found := units[unitPart]
	if !found {
		return 0, fmt.Errorf("unknown buffer size unit provided: %s", unitPart)
	}

	return num * multiplier, nil
}
