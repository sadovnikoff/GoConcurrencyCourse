package common

import (
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"
)

const testConfig = `
engine:
  type: "in_memory_test"
network:
  address: "127.0.0.1:9999"
  max_connections: 50
  max_message_size: "2KB"
  idle_timeout: 3m
logging:
  level: "info"
  output: "/test/output.log"
`

func TestParseConfig(t *testing.T) {
	tests := []struct {
		name          string
		reader        io.Reader
		expectedCfg   *Config
		isErrExpected bool
		expectedErr   error
	}{
		{
			name:   "load config",
			reader: strings.NewReader(testConfig),
			expectedCfg: &Config{
				Engine: &EngineConfig{
					Type: "in_memory_test",
				},
				Network: &NetworkConfig{
					Address:        "127.0.0.1:9999",
					MaxConnections: 50,
					MaxMsgSize:     "2KB",
					IdleTimeout:    "3m",
				},
				Logging: &LoggingConfig{
					Level:  "info",
					Output: "/test/output.log",
				},
			},
			isErrExpected: false,
			expectedErr:   nil,
		},
		{
			name:   "load empty config",
			reader: strings.NewReader(""),
			expectedCfg: &Config{
				Engine: &EngineConfig{
					Type: inMemoryEngine,
				},
				Network: &NetworkConfig{
					Address:        serverAddress,
					MaxConnections: maxConnections,
					MaxMsgSize:     maxMessageSize,
					IdleTimeout:    idleTimeout,
				},
				Logging: &LoggingConfig{
					Level: loggingLevel,
				},
			},
			isErrExpected: false,
			expectedErr:   nil,
		},
		{
			name:          "load config with nil reader",
			reader:        nil,
			expectedCfg:   nil,
			isErrExpected: true,
			expectedErr:   errors.New("nil reader provided"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := ParseConfig(tt.reader)

			if tt.isErrExpected {
				if tt.expectedErr.Error() != err.Error() {
					t.Errorf("want %+v; got %+v", tt.expectedErr, err)
				}

				if cfg != nil {
					t.Errorf("want %+v; got %+v", tt.expectedCfg, cfg)
				}
			} else {
				if tt.expectedErr != nil {
					t.Errorf("want %+v; got %+v", tt.expectedErr, err)
				}

				if cfg == nil {
					t.Errorf("want %+v; got %+v", tt.expectedCfg, cfg)
				}

				if !reflect.DeepEqual(cfg, tt.expectedCfg) {
					t.Errorf("want %+v; got %+v", tt.expectedCfg, cfg)
				}
			}
		})
	}
}

func TestParseBufSize(t *testing.T) {
	tests := []struct {
		name           string
		cfg            *Config
		ExpectedSize   int
		expectedErr    error
		expectedNilErr bool
	}{
		{
			name: "Parse invalid buffer size",
			cfg: &Config{
				Network: &NetworkConfig{
					MaxMsgSize: "A2KB",
				},
			},
			ExpectedSize: 0,
			expectedErr:  errors.New("invalid buffer size provided: A2KB"),
		},
		{
			name: "Parse unsupported buffer size unit",
			cfg: &Config{
				Network: &NetworkConfig{
					MaxMsgSize: "2PB",
				},
			},
			ExpectedSize: 0,
			expectedErr:  errors.New("unknown buffer size unit provided: PB"),
		},
		{
			name: "Parse valid buffer size",
			cfg: &Config{
				Network: &NetworkConfig{
					MaxMsgSize: "2KB",
				},
			},
			ExpectedSize:   2048,
			expectedErr:    nil,
			expectedNilErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bufSize, err := ParseSize(tt.cfg.Network.MaxMsgSize)

			if tt.ExpectedSize != bufSize {
				t.Errorf("want %+v; got %+v", tt.ExpectedSize, bufSize)
			}

			if tt.expectedNilErr {
				if err != nil {
					t.Errorf("want %+v; got %+v", tt.expectedErr, err)
				}
			} else {
				if tt.expectedErr.Error() != err.Error() {
					t.Errorf("want %+v; got %+v", tt.expectedErr, err)
				}
			}
		})
	}
}
