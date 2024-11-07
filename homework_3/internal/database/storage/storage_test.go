package storage

import (
	"errors"
	"testing"

	"github.com/sadovnikoff/GoConcurrencyCourse/homework_3/internal/common"
)

func TestNewStorage(t *testing.T) {
	tests := []struct {
		name           string
		logger         *common.Logger
		engine         Engine
		expectedError  error
		expectedNilObj bool
	}{
		{
			name:   "New storage without engine",
			engine: nil,
			logger: func() *common.Logger {
				logger, _ := common.NewLogger("", "")
				return logger
			}(),
			expectedError:  errors.New("engine is invalid"),
			expectedNilObj: true,
		},
		{
			name:           "New storage without logger",
			engine:         NewMockEngine(),
			expectedError:  errors.New("logger is invalid"),
			expectedNilObj: true,
		},
		{
			name:   "New valid storage",
			engine: NewMockEngine(),
			logger: func() *common.Logger {
				logger, _ := common.NewLogger("", "")
				return logger
			}(),
			expectedError:  errors.New("engine is invalid"),
			expectedNilObj: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage, err := NewStorage(tt.engine, nil, tt.logger)

			if tt.expectedNilObj {
				if tt.expectedError.Error() != err.Error() {
					t.Errorf("want %+v; got %+v", tt.expectedError, err)
				}

				if storage != nil {
					t.Errorf("want %+v; got %+v", nil, storage)
				}
			} else {
				if err != nil {
					t.Errorf("want %+v; got %+v", tt.expectedError, err)
				}

				if storage == nil {
					t.Errorf("want %+v; got %+v", "not nil storage", nil)
				}
			}
		})
	}
}

func TestStorage_Set(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value string
	}{
		{
			name:  "Storage SET command",
			key:   "some_key",
			value: "some_value",
		},
	}

	logger, _ := common.NewLogger("", "")
	storage, err := NewStorage(NewMockEngine(), nil, logger)
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage.Set(tt.key, tt.value)

			value, err := storage.Get(tt.key)
			if errors.Is(err, ErrNotFound) {
				t.Errorf("want %+v; got %+v", nil, err)
			}

			if value != tt.value {
				t.Errorf("want %+v; got %+v", tt.value, value)
			}
		})
	}
}

func TestStorage_Del(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value string
	}{
		{
			name:  "Storage DEL command",
			key:   "some_key",
			value: "",
		},
	}

	logger, _ := common.NewLogger("", "")
	storage, err := NewStorage(NewMockEngine(), nil, logger)
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage.Set(tt.key, "some_value")
			storage.Del(tt.key)

			value, err := storage.Get(tt.key)
			if !errors.Is(err, ErrNotFound) {
				t.Errorf("want %+v; got %+v", ErrNotFound, err)
			}

			if value != tt.value {
				t.Errorf("want %+v; got %+v", tt.value, value)
			}
		})
	}
}
