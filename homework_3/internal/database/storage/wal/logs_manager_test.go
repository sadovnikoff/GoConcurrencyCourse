package wal

import (
	"errors"
	"reflect"
	"testing"

	"github.com/sadovnikoff/GoConcurrencyCourse/homework_3/internal/common"
	"github.com/sadovnikoff/GoConcurrencyCourse/homework_3/internal/database/filesystem"
)

func TestNewLogsManager(t *testing.T) {
	tests := []struct {
		name           string
		segment        segment
		logger         *common.Logger
		expectedError  error
		expectedNilObj bool
	}{
		{
			name:           "New LogsManager without segment",
			expectedNilObj: true,
			expectedError:  errors.New("segment is invalid"),
		},
		{
			name:           "New LogsManager without logger",
			segment:        filesystem.NewSegment("tmp_test_data", 10),
			expectedNilObj: true,
			expectedError:  errors.New("logger is invalid"),
		},
		{
			name:    "New LogsManager",
			segment: filesystem.NewSegment("tmp_test_data", 10),
			logger: func() *common.Logger {
				logger, _ := common.NewLogger("", "")
				return logger
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logsManager, err := NewLogsManager(tt.segment, tt.logger)

			if tt.expectedNilObj {
				if tt.expectedError.Error() != err.Error() {
					t.Errorf("want %+v; got %+v", tt.expectedError, err)
				}

				if logsManager != nil {
					t.Errorf("want %+v; got %+v", nil, logsManager)
				}
			} else {
				if err != nil {
					t.Errorf("want %+v; got %+v", tt.expectedError, err)
				}

				if logsManager == nil {
					t.Errorf("want %+v; got %+v", "not nil object", nil)
				}

				if logsManager.segment == nil {
					t.Errorf("want %+v; got %+v", "not nil segment", nil)
				}
			}
		})
	}
}

func TestLogsManager_Write(t *testing.T) {
	MockRequest1.doneStatus = make(chan error, 1)
	MockRequest2.doneStatus = make(chan error, 1)
	requests := []Request{
		MockRequest1,
		MockRequest2,
	}

	segment := NewMockWalSegment(false)
	logger, err := common.NewLogger("", "")
	if err != nil {
		t.Errorf("logger creation issue: %s", err)
	}

	logsManager, err := NewLogsManager(segment, logger)
	if err != nil {
		t.Errorf("logs manager creation issue: %s", err)
	}

	logsManager.Write(requests)

	for _, request := range requests {
		err := <-request.doneStatus
		if err != nil {
			t.Errorf("logs manager write issue: %s", err)
		}

		_, ok := <-request.doneStatus
		if ok {
			t.Errorf("logs manager write issue: channel was not closed [request: %+v]", request)
		}
	}
}

func TestLogsManager_Write_WithSegmentError(t *testing.T) {

	MockRequest1.doneStatus = make(chan error, 1)
	MockRequest2.doneStatus = make(chan error, 1)
	requests := []Request{
		MockRequest1,
		MockRequest2,
	}

	segment := NewMockWalSegment(true)
	logger, err := common.NewLogger("", "")
	if err != nil {
		t.Errorf("logger creation issue: %s", err)
	}

	logsManager, err := NewLogsManager(segment, logger)
	if err != nil {
		t.Errorf("logs manager creation issue: %s", err)
	}

	logsManager.Write(requests)

	for _, request := range requests {
		err := <-request.doneStatus
		if err.Error() != TestWriteSegmentError {
			t.Errorf("logs manager write issue: expected error %s, got %s", TestWriteSegmentError, err)
		}

		_, ok := <-request.doneStatus
		if ok {
			t.Errorf("logs manager write issue: channel was not closed [request: %+v]", request)
		}
	}
}

func TestLogsManager_Read(t *testing.T) {
	segment := NewMockWalSegment(false)
	logger, err := common.NewLogger("", "")
	if err != nil {
		t.Errorf("logger creation issue: %s", err)
	}

	logsManager, err := NewLogsManager(segment, logger)
	if err != nil {
		t.Errorf("logs manager creation issue: %s", err)
	}

	requests, err := logsManager.Read()
	if err != nil {
		t.Errorf("logs manager read issue: %s", err)
	}

	if len(requests) != 2 {
		t.Errorf("logs manager read issue: expected 2 requests, got %d", len(requests))
	}

	MockRequest1.doneStatus = nil
	MockRequest2.doneStatus = nil
	expectedRequests := []Request{
		MockRequest1,
		MockRequest2,
	}
	if !reflect.DeepEqual(requests, expectedRequests) {
		t.Errorf("logs manager read issue: got %+v requests, expected %+v", requests, expectedRequests)
	}
}

func TestLogsManager_Read_WithSegmentError(t *testing.T) {
	segment := NewMockWalSegment(true)
	logger, err := common.NewLogger("", "")
	if err != nil {
		t.Errorf("logger creation issue: %s", err)
	}

	logsManager, err := NewLogsManager(segment, logger)
	if err != nil {
		t.Errorf("logs manager creation issue: %s", err)
	}

	requests, err := logsManager.Read()
	if err == nil || err.Error() != TestReadAllSegmentError {
		t.Errorf("logs manager read issue: expected error %s, got %s", TestReadAllSegmentError, err)
	}

	if len(requests) != 0 {
		t.Errorf("logs manager read issue: expected 0 requests, got %d", len(requests))
	}
}
