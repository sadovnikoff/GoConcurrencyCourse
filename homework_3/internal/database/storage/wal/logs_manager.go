package wal

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/sadovnikoff/GoConcurrencyCourse/homework_3/internal/common"
)

type segment interface {
	Write([]byte) error
	ReadAll() ([][]byte, error)
}

type LogsManager struct {
	segment segment
	logger  *common.Logger
}

func NewLogsManager(segment segment, logger *common.Logger) (*LogsManager, error) {
	if segment == nil {
		return nil, errors.New("segment is invalid")
	}

	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	return &LogsManager{segment: segment, logger: logger}, nil
}

func (l *LogsManager) Write(requests []Request) {
	var buffer bytes.Buffer
	for _, req := range requests {
		if err := req.Encode(&buffer); err != nil {
			l.acknowledge(requests, err)
			return
		}
	}

	err := l.segment.Write(buffer.Bytes())
	if err != nil {
		l.logger.Error("failed to write request data: %s", err)
	}

	l.acknowledge(requests, err)
}

func (l *LogsManager) Read() ([]Request, error) {
	segmentsData, err := l.segment.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read segments: %w", err)
	}

	var requests []Request
	for _, data := range segmentsData {
		requests, err = l.readSegment(requests, data)
		if err != nil {
			return nil, fmt.Errorf("failed to read segments: %w", err)
		}
	}

	return requests, nil
}

func (l *LogsManager) acknowledge(requests []Request, err error) {
	for _, req := range requests {
		req.doneStatus <- err
		close(req.doneStatus)
	}
}

func (l *LogsManager) readSegment(requests []Request, data []byte) ([]Request, error) {
	buffer := bytes.NewBuffer(data)
	for buffer.Len() > 0 {
		var request Request
		if err := request.Decode(buffer); err != nil {
			return nil, fmt.Errorf("failed to parse logs data: %w", err)
		}

		requests = append(requests, request)
	}

	return requests, nil
}
