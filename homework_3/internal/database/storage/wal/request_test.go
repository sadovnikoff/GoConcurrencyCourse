package wal

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/sadovnikoff/GoConcurrencyCourse/homework_3/internal/database/compute"
)

func TestSingleSerialization(t *testing.T) {
	expectedRequest := Request{
		Command:   compute.SetCommand,
		Arguments: []string{"key", "value"},
	}

	var writeBuffer bytes.Buffer
	err := expectedRequest.Encode(&writeBuffer)
	if err != nil {
		t.Errorf("cannot encode request: %s", err)
	}

	data := writeBuffer.Bytes()
	readBuffer := bytes.NewBuffer(data)

	var request Request
	err = request.Decode(readBuffer)
	if err != nil {
		t.Errorf("cannot decode request: %s", err)
	}

	if request.Command != expectedRequest.Command {
		t.Errorf("comparison issue: got command %s, expected %s", request.Command, expectedRequest.Command)
	}

	if len(request.Arguments) != len(expectedRequest.Arguments) {
		t.Errorf("comparison issue: got %d arguments, expected %d", len(request.Arguments), len(expectedRequest.Arguments))
	}

	if request.Arguments[0] != expectedRequest.Arguments[0] || request.Arguments[1] != expectedRequest.Arguments[1] {
		t.Errorf("comparison issue: got arguments %+v, expected %+v", request.Arguments, expectedRequest.Arguments)
	}
}

func TestMultipleSerialization(t *testing.T) {
	t.Parallel()

	expectedRequests := []Request{
		{
			Command:   compute.SetCommand,
			Arguments: []string{"key", "value"},
		},
		{
			Command:   compute.DelCommand,
			Arguments: []string{"key"},
		},
	}

	var writeBuffer bytes.Buffer
	for _, request := range expectedRequests {
		err := request.Encode(&writeBuffer)
		if err != nil {
			t.Errorf("cannot encode request: %s", err)
		}
	}

	data := writeBuffer.Bytes()
	readBuffer := bytes.NewBuffer(data)

	var requests []Request
	for readBuffer.Len() > 0 {
		var request Request
		err := request.Decode(readBuffer)
		if err != nil {
			t.Errorf("cannot decode request: %s", err)
		}

		requests = append(requests, request)
	}

	if !reflect.DeepEqual(requests, expectedRequests) {
		t.Errorf("comparison issue: got %+v requests, expected %+v", requests, expectedRequests)
	}
}
