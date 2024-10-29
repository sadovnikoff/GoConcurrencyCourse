package tcp

import (
	"net"
	"sync"
	"testing"
	"time"

	"github.com/sadovnikoff/GoConcurrencyCourse/homework_3/internal/common"
)

func TestNewServer(t *testing.T) {
	tests := []struct {
		name            string
		cfg             *common.Config
		db              database
		logger          *common.Logger
		expectedBufSize int
		expectedAddr    string
		expectedTimeout time.Duration
		expectedNilObj  bool
	}{
		{
			name: "New server without config",
			cfg:  nil,
			db:   NewMockDatabase(),
			logger: func() *common.Logger {
				logger, _ := common.NewLogger("", "")
				return logger
			}(),
			expectedNilObj: true,
		},
		{
			name: "New server without database",
			cfg:  &common.Config{},
			db:   nil,
			logger: func() *common.Logger {
				logger, _ := common.NewLogger("", "")
				return logger
			}(),
			expectedNilObj: true,
		},
		{
			name:           "New server without logger",
			cfg:            &common.Config{},
			db:             NewMockDatabase(),
			logger:         nil,
			expectedNilObj: true,
		},
		{
			name: "New server with invalid address",
			cfg: &common.Config{
				Network: &common.NetworkConfig{
					Address:     "invalid address",
					MaxMsgSize:  "4KB",
					IdleTimeout: "5m",
				},
			},
			db: NewMockDatabase(),
			logger: func() *common.Logger {
				logger, _ := common.NewLogger("", "")
				return logger
			}(),
			expectedNilObj: true,
		},
		{
			name: "New server with invalid idle timeout",
			cfg: &common.Config{
				Network: &common.NetworkConfig{
					Address:     "127.0.0.1:8090",
					MaxMsgSize:  "2KB",
					IdleTimeout: "not valid",
				},
			},
			db: NewMockDatabase(),
			logger: func() *common.Logger {
				logger, _ := common.NewLogger("", "")
				return logger
			}(),
			expectedTimeout: defaultIdleTimeout,
			expectedBufSize: 2048,
			expectedAddr:    "127.0.0.1:8090",
			expectedNilObj:  false,
		},
		{
			name: "New server with invalid buffer size",
			cfg: &common.Config{
				Network: &common.NetworkConfig{
					Address:     "127.0.0.1:8080",
					MaxMsgSize:  "5BB",
					IdleTimeout: "3m",
				},
			},
			db: NewMockDatabase(),
			logger: func() *common.Logger {
				logger, _ := common.NewLogger("", "")
				return logger
			}(),
			expectedBufSize: defaultBufSize,
			expectedTimeout: 180 * time.Second,
			expectedAddr:    "127.0.0.1:8080",
			expectedNilObj:  false,
		},
		{
			name: "New valid server",
			cfg: &common.Config{
				Network: &common.NetworkConfig{
					Address:     "127.0.0.1:8080",
					MaxMsgSize:  "4KB",
					IdleTimeout: "2m",
				},
			},
			db: NewMockDatabase(),
			logger: func() *common.Logger {
				logger, _ := common.NewLogger("", "")
				return logger
			}(),
			expectedBufSize: 4096,
			expectedAddr:    "127.0.0.1:8080",
			expectedTimeout: 120 * time.Second,
			expectedNilObj:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, err := NewServer(tt.cfg, tt.db, tt.logger)

			if tt.expectedNilObj {
				if err == nil {
					t.Errorf("want not nil error; got %+v", err)
				}

				if server != nil {
					t.Errorf("want %+v; got %+v", nil, server)
				}
			} else {
				if err != nil {
					t.Errorf("want nil error; got %+v", err)
				}

				if server == nil {
					t.Errorf("want %+v; got %+v", "not nil storage", nil)
				}

				if server.lis.Addr().String() != tt.expectedAddr {
					t.Errorf("want %+v; got %+v", tt.expectedAddr, server.lis.Addr().String())
				}

				if server.idleTimeout != tt.expectedTimeout {
					t.Errorf("want %+v; got %+v", tt.expectedTimeout, server.idleTimeout)
				}

				if server.bufferSize != tt.expectedBufSize {
					t.Errorf("want %+v; got %+v", tt.expectedBufSize, server.bufferSize)
				}

				if err := server.lis.Close(); err != nil {
					server.logger.Error("failed to close listener %s", err.Error())
				}
			}
		})
	}
}

func TestServer_Run(t *testing.T) {
	t.Parallel()

	addr := "127.0.0.1:8080"
	cfg := &common.Config{
		Network: &common.NetworkConfig{
			Address:        addr,
			MaxConnections: 5,
			MaxMsgSize:     "4KB",
			IdleTimeout:    "2m",
		},
	}

	logger, _ := common.NewLogger("", "")
	server, err := NewServer(cfg, NewMockDatabase(), logger)
	if err != nil {
		t.Errorf("want nil error; got %+v", err)
	}

	go server.Run()

	time.Sleep(100 * time.Millisecond)
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		connection, clientErr := net.Dial("tcp", addr)
		if clientErr != nil {
			t.Errorf("want nil error; got %+v", err)
		}

		_, clientErr = connection.Write([]byte("client-1"))
		if err != nil {
			t.Errorf("want nil error; got %+v", err)
		}

		buffer := make([]byte, 1024)
		size, clientErr := connection.Read(buffer)
		if err != nil {
			t.Errorf("want nil error; got %+v", err)
		}

		clientErr = connection.Close()
		if err != nil {
			t.Errorf("want nil error; got %+v", err)
		}

		if string(buffer[:size]) != "successful response to the [client-1] request" {
			t.Errorf("want: successful response to the [client-1] request; got: %+v", string(buffer[:size]))
		}
	}()

	go func() {
		defer wg.Done()

		connection, clientErr := net.Dial("tcp", addr)
		if clientErr != nil {
			t.Errorf("want nil error; got %+v", err)
		}

		_, clientErr = connection.Write([]byte("client-2 error"))
		if clientErr != nil {
			t.Errorf("want nil error; got %+v", err)
		}

		buffer := make([]byte, 1024)
		size, clientErr := connection.Read(buffer)
		if clientErr != nil {
			t.Errorf("want nil error; got %+v", err)
		}

		clientErr = connection.Close()
		if clientErr != nil {
			t.Errorf("want nil error; got %+v", err)
		}

		if string(buffer[:size]) != "error has been occurred during request handling" {
			t.Errorf("want: error has been occurred during request handling; got: %+v", string(buffer[:size]))
		}
	}()

	wg.Wait()

	if err := server.lis.Close(); err != nil {
		t.Errorf("failed to close listener %s", err.Error())
	}
}
