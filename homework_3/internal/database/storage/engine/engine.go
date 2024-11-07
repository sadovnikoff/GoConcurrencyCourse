package engine

import (
	"errors"
	"sync"

	"github.com/sadovnikoff/GoConcurrencyCourse/homework_3/internal/common"
	"github.com/sadovnikoff/GoConcurrencyCourse/homework_3/internal/database/storage"
)

const InMemoryEngine = "in_memory"

type Engine struct {
	logger *common.Logger

	m  sync.Mutex
	DB map[string]string
}

func NewEngine(logger *common.Logger) (*Engine, error) {
	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	return &Engine{DB: make(map[string]string), logger: logger}, nil
}

func (e *Engine) Set(key, value string) {
	e.m.Lock()
	e.DB[key] = value
	e.m.Unlock()

	e.logger.Debug("successful SET query [key %s, value %s]", key, value)
}

func (e *Engine) Get(key string) (string, error) {
	e.m.Lock()
	value, ok := e.DB[key]
	e.m.Unlock()
	if !ok {
		e.logger.Debug("GET query [key %s, value %s]: key not found", key, value)
		return "", storage.ErrNotFound
	}

	e.logger.Info("successful GET query [key %s, value %s]", key, value)
	return value, nil
}

func (e *Engine) Del(key string) {
	e.m.Lock()
	delete(e.DB, key)
	e.m.Unlock()

	e.logger.Debug("successful DEL query [key %s]", key)
}
