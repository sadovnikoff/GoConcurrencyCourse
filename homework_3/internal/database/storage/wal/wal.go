package wal

import (
	"errors"
	"sync"
	"time"

	"github.com/sadovnikoff/GoConcurrencyCourse/homework_3/internal/common"
	"github.com/sadovnikoff/GoConcurrencyCourse/homework_3/internal/database/compute"
	"github.com/sadovnikoff/GoConcurrencyCourse/homework_3/internal/database/filesystem"
)

const (
	defaultFlushingTimeout = 10 * time.Millisecond
	defaultSegmentSize     = 10_485_760
)

type logManager interface {
	Write([]Request)
	Read() ([]Request, error)
}

type WAL struct {
	logsManager     logManager
	batchSize       int
	segmentSize     int
	batches         chan []Request
	writeStatus     <-chan error
	flushingTimeout time.Duration
	logger          *common.Logger

	mutex sync.Mutex
	batch []Request
}

func NewWAL(cfg *common.WalConfig, logger *common.Logger) (*WAL, error) {
	if cfg == nil {
		return nil, nil
	}

	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	timeout, err := time.ParseDuration(cfg.FlushingTimeout)
	if err != nil {
		timeout = defaultFlushingTimeout
	}

	segmentSize, err := common.ParseSize(cfg.SegmentSize)
	if err != nil {
		segmentSize = defaultSegmentSize
	}

	segment := filesystem.NewSegment(cfg.DirPath, segmentSize)
	logsManager, err := NewLogsManager(segment, logger)
	if err != nil {
		return nil, err
	}

	wal := &WAL{
		logsManager:     logsManager,
		batchSize:       cfg.BatchSize,
		batches:         make(chan []Request, 1),
		flushingTimeout: timeout,
		segmentSize:     segmentSize,
		logger:          logger,
	}

	return wal, nil
}

func (w *WAL) Start() {
	go func() {
		ticker := time.NewTicker(w.flushingTimeout)
		defer ticker.Stop()

		for {
			select {
			case batch := <-w.batches:
				w.logsManager.Write(batch)
				ticker.Reset(w.flushingTimeout)
			case <-ticker.C:
				w.flushBatch()
			}
		}
	}()
}

func (w *WAL) Set(key, value string) error {
	w.push(compute.SetCommand, []string{key, value})
	return <-w.writeStatus
}

func (w *WAL) Del(key string) error {
	w.push(compute.DelCommand, []string{key})
	return <-w.writeStatus
}

func (w *WAL) Recover() ([]Request, error) {
	return w.logsManager.Read()
}

func (w *WAL) push(cmd string, args []string) {
	request := NewRequest(cmd, args)

	w.mutex.Lock()
	w.batch = append(w.batch, request)
	if len(w.batch) == w.batchSize {
		w.batches <- w.batch
		w.batch = nil
	}
	w.mutex.Unlock()

	w.writeStatus = request.doneStatus
}

func (w *WAL) flushBatch() {
	var batch []Request

	w.mutex.Lock()
	batch = w.batch
	w.batch = nil
	w.mutex.Unlock()

	if len(batch) != 0 {
		w.logsManager.Write(batch)
	}
}
