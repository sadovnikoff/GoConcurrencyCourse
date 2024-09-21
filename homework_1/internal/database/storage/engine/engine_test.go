package engine

import (
	"errors"
	"testing"

	"github.com/sadovnikoff/GoConcurrencyCourse/homework_1/internal/common"
	"github.com/sadovnikoff/GoConcurrencyCourse/homework_1/internal/database/storage"
)

func TestNewEngine(t *testing.T) {
	tests := []struct {
		name              string
		logger            *common.Logger
		expectedError     error
		expectedNilEngine bool
	}{
		{
			name:              "New engine without logger",
			expectedNilEngine: true,
			expectedError:     errors.New("logger is invalid"),
		},
		{
			name:              "New engine with logger",
			logger:            common.NewLogger(),
			expectedNilEngine: false,
			expectedError:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine, err := NewEngine(tt.logger)

			if tt.expectedNilEngine {
				if tt.expectedError.Error() != err.Error() {
					t.Errorf("want %+v; got %+v", tt.expectedError, err)
				}

				if engine != nil {
					t.Errorf("want %+v; got %+v", nil, engine)
				}
			} else {
				if err != nil {
					t.Errorf("want %+v; got %+v", tt.expectedError, err)
				}

				if engine == nil {
					t.Errorf("want %+v; got %+v", "not nil engine", nil)
				}

				if engine.DB == nil {
					t.Errorf("want %+v; got %+v", "not nil db", engine.DB)
				}
			}
		})
	}
}

func TestEngine_Set(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value string
	}{
		{
			name:  "SET - insert a value by a key",
			key:   "some_key",
			value: "some_value",
		},
	}

	engine, err := NewEngine(common.NewLogger())
	if err != nil {
		t.Errorf("want %+v; got %+v", nil, err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine.Set(tt.key, tt.value)

			value, err := engine.Get(tt.key)
			if errors.Is(err, storage.ErrNotFound) {
				t.Errorf("want %+v; got %+v", nil, err)
			}

			if value != tt.value {
				t.Errorf("want %+v; got %+v", tt.value, value)
			}
		})
	}
}

func TestEngine_Get(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value string
	}{
		{
			name:  "GET - Get a value by a non-exist key",
			key:   "some_key",
			value: "",
		},
	}

	engine, err := NewEngine(common.NewLogger())
	if err != nil {
		t.Errorf("want %+v; got %+v", nil, err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			value, err := engine.Get(tt.key)
			if !errors.Is(err, storage.ErrNotFound) {
				t.Errorf("want %+v; got %+v", storage.ErrNotFound, err)
			}

			if value != tt.value {
				t.Errorf("want %+v; got %+v", tt.value, value)
			}
		})
	}
}

func TestEngine_Del(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value string
	}{
		{
			name:  "Del - delete a value by a key",
			key:   "some_key",
			value: "",
		},
	}

	engine, err := NewEngine(common.NewLogger())
	if err != nil {
		t.Errorf("want %+v; got %+v", nil, err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine.Set(tt.key, "some_value")
			engine.Del(tt.key)

			value, err := engine.Get(tt.key)
			if !errors.Is(err, storage.ErrNotFound) {
				t.Errorf("want %+v; got %+v", storage.ErrNotFound, err)
			}

			if value != tt.value {
				t.Errorf("want %+v; got %+v", tt.value, value)
			}
		})
	}
}
