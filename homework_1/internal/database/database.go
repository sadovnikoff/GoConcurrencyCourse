package database

import (
	"errors"
	"fmt"

	"sadovnikoff/go_concurrency_cource/homework_1/internal/common"
	"sadovnikoff/go_concurrency_cource/homework_1/internal/database/compute"
)

type computeLayer interface {
	Parse(string) (compute.Query, error)
}

type storageLayer interface {
	Set(string, string)
	Get(string) (string, error)
	Del(string)
}

type Database struct {
	computeLayer computeLayer
	storageLayer storageLayer
	logger       *common.Logger
}

func NewDatabase(computeLayer computeLayer, storageLayer storageLayer, logger *common.Logger) (*Database, error) {
	if computeLayer == nil {
		return nil, errors.New("compute is invalid")
	}

	if storageLayer == nil {
		return nil, errors.New("storage is invalid")
	}

	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	return &Database{
		computeLayer: computeLayer,
		storageLayer: storageLayer,
		logger:       logger,
	}, nil
}

func (d *Database) HandleQuery(request string) (string, error) {
	d.logger.ILog.Printf("handling request [%s]\n", request)

	query, err := d.computeLayer.Parse(request)
	if err != nil {
		d.logger.DLog.Printf("compute layer is incorrect")
		return "", err
	}

	var response string
	switch query.Command() {
	case compute.SetCommand:
		d.storageLayer.Set(query.KeyArgument(), query.ValueArgument())
		response = "[ok]"
	case compute.GetCommand:
		val, err := d.storageLayer.Get(query.KeyArgument())
		if err != nil {
			return "", err
		}
		response = fmt.Sprintf("[ok] %s", val)
	case compute.DelCommand:
		d.storageLayer.Del(query.KeyArgument())
		response = fmt.Sprintf("[ok]")
	}

	return response, nil
}
