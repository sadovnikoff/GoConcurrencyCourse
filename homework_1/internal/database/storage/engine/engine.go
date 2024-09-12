package engine

import (
	"errors"
	"sadovnikoff/go_concurrency_cource/homework_1/internal/common"
	"sadovnikoff/go_concurrency_cource/homework_1/internal/database/storage"
)

type Engine struct {
	DB     map[string]string
	logger *common.Logger
}

func NewEngine(logger *common.Logger) (*Engine, error) {
	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	return &Engine{DB: make(map[string]string), logger: logger}, nil
}

func (e *Engine) Set(key, value string) {
	e.DB[key] = value
	e.logger.ILog.Printf("successfull SET query [key %s, value %s]", key, value)
}

func (e *Engine) Get(key string) (string, error) {
	value, ok := e.DB[key]
	if !ok {
		return "", storage.ErrNotFound
	}

	e.logger.ILog.Printf("successfull GET query [key %s, value %s]", key, value)
	return value, nil
}

func (e *Engine) Del(key string) {
	delete(e.DB, key)
	e.logger.ILog.Printf("successfull DEL query [key %s]", key)
}
