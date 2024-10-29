package wal

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"

	"github.com/sadovnikoff/GoConcurrencyCourse/homework_3/internal/database/compute"
)

const (
	TestWriteSegmentError   = "write segment error"
	TestReadAllSegmentError = "failed to read segments: read all segment error"
)

var MockRequest1 = Request{
	Command:   compute.SetCommand,
	Arguments: []string{"key", "value"},
}

var MockRequest2 = Request{
	Command:   compute.DelCommand,
	Arguments: []string{"key"},
}

// MockWalSegment is a mock of wal Segment interface
type MockWalSegment struct {
	ShouldFailWrite bool
}

func NewMockWalSegment(shouldFailWrite bool) *MockWalSegment {
	return &MockWalSegment{
		ShouldFailWrite: shouldFailWrite,
	}
}

func (mws *MockWalSegment) Write(data []byte) error {
	if mws.ShouldFailWrite {
		return errors.New(TestWriteSegmentError)
	}

	readBuffer := bytes.NewBuffer(data)
	var requests []Request
	for readBuffer.Len() > 0 {
		var request Request
		err := request.Decode(readBuffer)
		if err != nil {
			return err
		}

		requests = append(requests, request)
	}

	if len(requests) != 2 {
		return fmt.Errorf("wrong amount of requests: got %d, expected 2", len(requests))
	}

	expectedRequests := []Request{
		MockRequest1,
		MockRequest2,
	}

	if requests[0].Command != expectedRequests[0].Command {
		return fmt.Errorf("wrong request command: got %s, expected %s", requests[0].Command, expectedRequests[0].Command)
	}

	if !reflect.DeepEqual(requests[0].Arguments, expectedRequests[0].Arguments) {
		return fmt.Errorf("wrong request arguments: got %+v, expected %+v", requests[0].Arguments, expectedRequests[0].Arguments)
	}

	if requests[1].Command != expectedRequests[1].Command {
		return fmt.Errorf("wrong request command: got %s, expected %s", requests[0].Command, expectedRequests[0].Command)
	}

	if !reflect.DeepEqual(requests[1].Arguments, expectedRequests[1].Arguments) {
		return fmt.Errorf("wrong request arguments: got %+v, expected %+v", requests[1].Arguments, expectedRequests[1].Arguments)
	}

	return nil
}

func (mws *MockWalSegment) ReadAll() ([][]byte, error) {
	if mws.ShouldFailWrite {
		return nil, errors.New("read all segment error")
	}

	var writeBuffer bytes.Buffer
	err := MockRequest1.Encode(&writeBuffer)
	if err != nil {
		return nil, err
	}

	err = MockRequest2.Encode(&writeBuffer)
	if err != nil {
		return nil, err
	}

	result := [][]byte{
		writeBuffer.Bytes(),
	}

	return result, nil
}
