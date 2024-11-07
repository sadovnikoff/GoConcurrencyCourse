package storage

import (
	"errors"

	"github.com/sadovnikoff/GoConcurrencyCourse/homework_3/internal/common"
	"github.com/sadovnikoff/GoConcurrencyCourse/homework_3/internal/database/compute"
	"github.com/sadovnikoff/GoConcurrencyCourse/homework_3/internal/database/storage/wal"
)

var ErrNotFound = errors.New("storage: requested data not found")

type Engine interface {
	Set(string, string)
	Get(string) (string, error)
	Del(string)
}

type WAL interface {
	Set(string, string) error
	Del(string) error
	Recover() ([]wal.Request, error)
}

type Storage struct {
	engine Engine
	wal    WAL
	logger *common.Logger
}

func NewStorage(engine Engine, wal WAL, logger *common.Logger) (*Storage, error) {
	if engine == nil {
		return nil, errors.New("engine is invalid")
	}

	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	storage := &Storage{
		engine: engine,
		wal:    wal,
		logger: logger,
	}

	if storage.wal != nil {
		requests, err := storage.wal.Recover()
		if err != nil {
			logger.Error("failed to recover data from WAL: %s", err)
		} else {
			storage.restore(requests)
		}
	}

	return storage, nil
}

func (s *Storage) Set(key, value string) error {
	if s.wal != nil {
		if err := s.wal.Set(key, value); err != nil {
			return err
		}
	}

	s.engine.Set(key, value)
	return nil
}

func (s *Storage) Get(key string) (string, error) {
	return s.engine.Get(key)
}

func (s *Storage) Del(key string) error {
	if s.wal != nil {
		if err := s.wal.Del(key); err != nil {
			return err
		}
	}

	s.engine.Del(key)
	return nil
}

func (s *Storage) restore(requests []wal.Request) {
	for _, request := range requests {
		switch request.Command {
		case compute.SetCommand:
			s.engine.Set(request.Arguments[0], request.Arguments[1])
		case compute.DelCommand:
			s.engine.Del(request.Arguments[0])
		}
	}
}
