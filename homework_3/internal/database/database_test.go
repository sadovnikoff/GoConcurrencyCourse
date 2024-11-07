package database

import (
	"errors"
	"testing"

	"github.com/sadovnikoff/GoConcurrencyCourse/homework_3/internal/common"
	"github.com/sadovnikoff/GoConcurrencyCourse/homework_3/internal/database/compute"
)

func TestNewDatabase(t *testing.T) {
	tests := []struct {
		name           string
		logger         *common.Logger
		computeLayer   computeLayer
		storageLayer   storageLayer
		expectedError  error
		expectedNilObj bool
	}{
		{
			name:         "New database without compute layer",
			computeLayer: nil,
			logger: func() *common.Logger {
				logger, _ := common.NewLogger("", "")
				return logger
			}(),
			storageLayer:   NewMockStorageLayer(),
			expectedError:  errors.New("compute is invalid"),
			expectedNilObj: true,
		},
		{
			name:         "New database without storage layer",
			computeLayer: NewMockComputeLayer(),
			storageLayer: nil,
			logger: func() *common.Logger {
				logger, _ := common.NewLogger("", "")
				return logger
			}(),
			expectedError:  errors.New("storage is invalid"),
			expectedNilObj: true,
		},
		{
			name:           "New database without logger",
			computeLayer:   NewMockComputeLayer(),
			storageLayer:   NewMockStorageLayer(),
			logger:         nil,
			expectedError:  errors.New("logger is invalid"),
			expectedNilObj: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			database, err := NewDatabase(tt.computeLayer, tt.storageLayer, tt.logger)

			if tt.expectedNilObj {
				if tt.expectedError.Error() != err.Error() {
					t.Errorf("want %+v; got %+v", tt.expectedError, err)
				}

				if database != nil {
					t.Errorf("want %+v; got %+v", nil, database)
				}
			} else {
				if err != nil {
					t.Errorf("want %+v; got %+v", tt.expectedError, err)
				}

				if database == nil {
					t.Errorf("want %+v; got %+v", "not nil storage", nil)
				}
			}
		})
	}
}

func TestDatabase_HandleQuery(t *testing.T) {
	tests := []struct {
		name          string
		cmd           string
		response      string
		isValid       bool
		expectedError error
	}{
		{
			name:          "Database HandleQuery SET command",
			cmd:           compute.SetCommand,
			response:      "[ok]",
			isValid:       true,
			expectedError: nil,
		},
		{
			name:          "Database HandleQuery GET command",
			cmd:           compute.GetCommand,
			response:      "[ok] value",
			isValid:       true,
			expectedError: nil,
		},
		{
			name:          "Database HandleQuery DEL command",
			cmd:           compute.DelCommand,
			response:      "[ok]",
			isValid:       true,
			expectedError: nil,
		},
		{
			name:          "Database HandleQuery invalid command",
			cmd:           "",
			response:      "",
			isValid:       false,
			expectedError: errors.New("compute layer is incorrect"),
		},
	}

	logger, _ := common.NewLogger("", "")
	database, err := NewDatabase(NewMockComputeLayer(), NewMockStorageLayer(), logger)
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := database.HandleQuery(tt.cmd)
			if tt.isValid {
				if err != nil {
					t.Errorf("want %+v; got %+v", tt.expectedError, err)
				}
			} else {
				if err == nil {
					t.Errorf("want %+v; got %+v", tt.expectedError, err)
				}
			}

			if result != tt.response {
				t.Errorf("want %+v; got %+v", tt.response, result)
			}

		})
	}
}
