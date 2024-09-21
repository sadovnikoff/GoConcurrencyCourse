package storage

import (
	"errors"

	"github.com/sadovnikoff/GoConcurrencyCourse/homework_2/internal/common"
)

var ErrNotFound = errors.New("storage: requested data not found")

type Engine interface {
	Set(string, string)
	Get(string) (string, error)
	Del(string)
}

type Storage struct {
	engine Engine
	logger *common.Logger
}

func NewStorage(engine Engine, logger *common.Logger) (*Storage, error) {
	if engine == nil {
		return nil, errors.New("engine is invalid")
	}

	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	storage := &Storage{
		engine: engine,
		logger: logger,
	}

	return storage, nil
}

func (s *Storage) Set(key, value string) {
	s.engine.Set(key, value)
}

func (s *Storage) Get(key string) (string, error) {
	return s.engine.Get(key)
}

func (s *Storage) Del(key string) {
	s.engine.Del(key)
}
