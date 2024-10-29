package wal

import (
	"errors"
	"github.com/sadovnikoff/GoConcurrencyCourse/homework_3/internal/database/compute"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/sadovnikoff/GoConcurrencyCourse/homework_3/internal/common"
)

func TestNewWAL(t *testing.T) {
	tests := []struct {
		name                 string
		config               *common.WalConfig
		logger               *common.Logger
		noConfig             bool
		wrongFlushingTimeout bool
		wrongSegmentSize     bool
		validConfig          bool
		expectedError        error
		expectedNilObj       bool
	}{
		{
			name:           "New WAL without config",
			expectedNilObj: true,
			noConfig:       true,
		},
		{
			name:           "New WAL without logger",
			config:         &common.WalConfig{},
			expectedNilObj: true,
			expectedError:  errors.New("logger is invalid"),
		},
		{
			name: "New WAL with wrong flushing timeout config",
			logger: func() *common.Logger {
				logger, _ := common.NewLogger("", "")
				return logger
			}(),
			wrongFlushingTimeout: true,
			config: &common.WalConfig{
				FlushingTimeout: "abc",
			},
		},
		{
			name: "New WAL with wrong segment size config",
			logger: func() *common.Logger {
				logger, _ := common.NewLogger("", "")
				return logger
			}(),
			wrongSegmentSize: true,
			config: &common.WalConfig{
				SegmentSize: "abc",
			},
		},
		{
			name: "New valid WAL",
			logger: func() *common.Logger {
				logger, _ := common.NewLogger("", "")
				return logger
			}(),
			validConfig: true,
			config: &common.WalConfig{
				SegmentSize:     "10MB",
				FlushingTimeout: "10ms",
				BatchSize:       100,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wal, err := NewWAL(tt.config, tt.logger)

			if tt.expectedNilObj {
				if tt.noConfig {
					if err != nil {
						t.Errorf("want nil; got %+v", err)
					}
				} else if err == nil {
					t.Errorf("want non-nil error; got nil")
				} else if err.Error() != tt.expectedError.Error() {
					t.Errorf("want %v; got %v", tt.expectedError, err)
				}

				if wal != nil {
					t.Errorf("want nil; got %+v", wal)
				}
			} else {
				if err != nil {
					t.Errorf("want nil; got %+v", err)
				}

				if wal == nil {
					t.Errorf("want not nil object; got nil")
				}

				if tt.wrongFlushingTimeout && wal.flushingTimeout != defaultFlushingTimeout {
					t.Errorf("wrong flushing timeout: want %s; got %s", defaultFlushingTimeout, wal.flushingTimeout)
				}

				if tt.wrongSegmentSize && wal.segmentSize != defaultSegmentSize {
					t.Errorf("wrong segment size: want %d; got %d", defaultSegmentSize, wal.flushingTimeout)
				}

				if tt.validConfig {
					if wal.logger == nil {
						t.Errorf("want not nil logger; got nil")
					}

					if wal.logsManager == nil {
						t.Errorf("want not nil logsManager; got nil")
					}

					if wal.segmentSize != 10485760 {
						t.Errorf("wrong segment size: want %d; got %d", 10485760, wal.segmentSize)
					}

					if wal.batchSize != tt.config.BatchSize {
						t.Errorf("wrong batch size: want %d; got %d", tt.config.BatchSize, wal.batchSize)
					}

					if wal.flushingTimeout != 10000000 {
						t.Errorf("wrong flashing timeout: want %d; got %d", 10000000, wal.flushingTimeout)
					}

					if wal.batches == nil {
						t.Errorf("want not nil batches; got nil")
					}
				}
			}
		})
	}
}

func TestWAL_Start_WriteByTimeout(t *testing.T) {
	logger, _ := common.NewLogger("", "")
	config := &common.WalConfig{
		BatchSize:       5,
		FlushingTimeout: "200ms",
		SegmentSize:     "1KB",
	}

	wal, err := NewWAL(config, logger)
	if err != nil {
		t.Errorf("wal cannot be created: %s", err)
	}

	start := time.Now()
	wal.Start()
	err = wal.Set("key1", "value1")
	duration := time.Since(start)

	if err != nil {
		t.Errorf("wal write error: %s", err)
	}

	if duration < 200*time.Millisecond || duration > 220*time.Millisecond {
		t.Errorf("wromg wal write timeout: expected 200-220ms, got %d", duration)
	}
}

func TestWAL_Start_WriteByBatchSize(t *testing.T) {
	logger, _ := common.NewLogger("", "")
	config := &common.WalConfig{
		BatchSize:       3,
		FlushingTimeout: "200ms",
		SegmentSize:     "1KB",
	}

	wal, err := NewWAL(config, logger)
	if err != nil {
		t.Errorf("wal cannot be created: %s", err)
	}

	var wg sync.WaitGroup
	wg.Add(3)

	start := time.Now()
	wal.Start()

	go func() {
		defer wg.Done()

		err = wal.Set("key1", "value1")
		if err != nil {
			t.Errorf("wal write error: %s", err)
		}
	}()

	go func() {
		defer wg.Done()

		err = wal.Del("key1")
		if err != nil {
			t.Errorf("wal write error: %s", err)
		}
	}()

	go func() {
		defer wg.Done()

		err = wal.Set("key1", "value1")
		if err != nil {
			t.Errorf("wal write error: %s", err)
		}
	}()

	wg.Wait()
	duration := time.Since(start)

	if duration > 20*time.Millisecond {
		t.Errorf("wromg wal write timeout: expected less than 20ms, got %d", duration)
	}
}

func TestWAL_Recover(t *testing.T) {
	logger, _ := common.NewLogger("", "")
	config := &common.WalConfig{
		BatchSize:       5,
		FlushingTimeout: "200ms",
		SegmentSize:     "1KB",
		DirPath:         "test_data",
	}

	wal, err := NewWAL(config, logger)
	if err != nil {
		t.Errorf("wal cannot be created: %s", err)
	}

	requests, err := wal.Recover()
	if err != nil {
		t.Errorf("wal recover error: %s", err)
	}

	if len(requests) != 3 {
		t.Errorf("wal recover got %d requests, expected 3", len(requests))
	}

	if requests[0].Command != compute.SetCommand {
		t.Errorf("wal recover got command %s, expected %s", requests[0].Command, compute.SetCommand)
	}

	if !reflect.DeepEqual(requests[0].Arguments, []string{"1", "a"}) {
		t.Errorf("wal recover got argumetns %+v, expected %+v", requests[0].Arguments, []string{"1", "a"})
	}

	if requests[1].Command != compute.SetCommand {
		t.Errorf("wal recover got command %s, expected %s", requests[0].Command, compute.SetCommand)
	}

	if !reflect.DeepEqual(requests[1].Arguments, []string{"2", "b"}) {
		t.Errorf("wal recover got argumetns %+v, expected %+v", requests[0].Arguments, []string{"2", "b"})
	}

	if requests[2].Command != compute.DelCommand {
		t.Errorf("wal recover got command %s, expected %s", requests[0].Command, compute.SetCommand)
	}

	if !reflect.DeepEqual(requests[2].Arguments, []string{"1"}) {
		t.Errorf("wal recover got argumetns %+v, expected %+v", requests[0].Arguments, []string{"1"})
	}
}
