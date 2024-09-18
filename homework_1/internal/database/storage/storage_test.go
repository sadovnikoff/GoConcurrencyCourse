package storage

import (
	"errors"
	"testing"

	"sadovnikoff/go_concurrency_cource/homework_1/internal/common"
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
			name:           "New storage without engine",
			engine:         nil,
			logger:         common.NewLogger(),
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
			name:           "New valid storage",
			engine:         NewMockEngine(),
			logger:         common.NewLogger(),
			expectedError:  errors.New("engine is invalid"),
			expectedNilObj: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage, err := NewStorage(tt.engine, tt.logger)

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

	storage, err := NewStorage(NewMockEngine(), common.NewLogger())
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

	storage, err := NewStorage(NewMockEngine(), common.NewLogger())
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
