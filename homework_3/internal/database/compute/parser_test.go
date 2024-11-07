package compute

import (
	"errors"
	"testing"

	"github.com/sadovnikoff/GoConcurrencyCourse/homework_3/internal/common"
)

func TestNewParser(t *testing.T) {
	tests := []struct {
		name              string
		logger            *common.Logger
		expectedError     error
		expectedNilParser bool
	}{
		{
			name:              "New parser without logger",
			expectedNilParser: true,
			expectedError:     errors.New("logger is invalid"),
		},
		{
			name: "New parser with logger",
			logger: func() *common.Logger {
				logger, _ := common.NewLogger("", "")
				return logger
			}(),
			expectedNilParser: false,
			expectedError:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser, err := NewParser(tt.logger)

			if tt.expectedNilParser {
				if tt.expectedError.Error() != err.Error() {
					t.Errorf("want %+v; got %+v", tt.expectedError, err)
				}

				if parser != nil {
					t.Errorf("want %+v; got %+v", nil, parser)
				}
			} else {
				if err != nil {
					t.Errorf("want %+v; got %+v", tt.expectedError, err)
				}

				if parser == nil {
					t.Errorf("want %+v; got %+v", "not nil parser", nil)
				}
			}
		})
	}
}

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name          string
		request       string
		expectedQuery Query
		expectedErr   error
	}{
		{
			name:          "Invalid request - less than 2 tokens",
			request:       "GET",
			expectedQuery: NewQuery("", "", ""),
			expectedErr:   errInvalidRequest,
		},
		{
			name:          "Invalid GET request - not enough args",
			request:       "GET ",
			expectedQuery: NewQuery("", "", ""),
			expectedErr:   errInvalidRequest,
		},
		{
			name:          "Valid GET request",
			request:       "GET key_1",
			expectedQuery: NewQuery("GET", "key_1", ""),
			expectedErr:   nil,
		},
		{
			name:          "Valid GET request with any characters",
			request:       "GET key_1/qw**as,agh.y(#)",
			expectedQuery: NewQuery("GET", "key_1/qw**as,agh.y(#)", ""),
			expectedErr:   nil,
		},
		{
			name:          "Valid DEL request",
			request:       "DEL some_key",
			expectedQuery: NewQuery("DEL", "some_key", ""),
			expectedErr:   nil,
		},
		{
			name:          "Valid SET request",
			request:       "SET some_key some_value",
			expectedQuery: NewQuery("SET", "some_key", "some_value"),
			expectedErr:   nil,
		},
		{
			name:          "Valid SET request with trailing spaces",
			request:       "	SET some_key some_value  ",
			expectedQuery: NewQuery("SET", "some_key", "some_value"),
			expectedErr:   nil,
		},
		{
			name:          "Invalid GET request - too many args",
			request:       "GET some_key qwe",
			expectedQuery: NewQuery("", "", ""),
			expectedErr:   errInvalidArguments,
		},
		{
			name:          "Invalid DEL request - too many args",
			request:       "DEL some_key qwe",
			expectedQuery: NewQuery("", "", ""),
			expectedErr:   errInvalidArguments,
		},
		{
			name:          "Invalid SET request - too many args",
			request:       "SET some_key qwe 123",
			expectedQuery: NewQuery("", "", ""),
			expectedErr:   errInvalidArguments,
		},
		{
			name:          "Invalid SET request - not enough args",
			request:       "SET some_key",
			expectedQuery: NewQuery("", "", ""),
			expectedErr:   errInvalidArguments,
		},
		{
			name:          "Invalid command request - lowercase SET",
			request:       "set some_key some_value",
			expectedQuery: NewQuery("", "", ""),
			expectedErr:   errInvalidCommand,
		},
		{
			name:          "Invalid command request - unknown command",
			request:       "qwerty some_key ",
			expectedQuery: NewQuery("", "", ""),
			expectedErr:   errInvalidCommand,
		},
	}

	logger, _ := common.NewLogger("", "")
	parser, err := NewParser(logger)
	if err != nil {
		t.Errorf("want %+v; got %+v", nil, err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, err := parser.Parse(tt.request)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("want %q; got %q", tt.expectedErr, err)
			}

			if query.Command() != tt.expectedQuery.Command() {
				t.Errorf("want %q; got %q", tt.expectedQuery.Command(), query.Command())
			}

			if query.KeyArgument() != tt.expectedQuery.KeyArgument() {
				t.Errorf("want %q; got %q", tt.expectedQuery.KeyArgument(), query.KeyArgument())
			}

			if query.ValueArgument() != tt.expectedQuery.ValueArgument() {
				t.Errorf("want %q; got %q", tt.expectedQuery.ValueArgument(), query.ValueArgument())
			}
		})
	}
}
