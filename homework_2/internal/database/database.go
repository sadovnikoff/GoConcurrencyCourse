package database

import (
	"errors"
	"fmt"

	"github.com/sadovnikoff/GoConcurrencyCourse/homework_2/internal/common"
	"github.com/sadovnikoff/GoConcurrencyCourse/homework_2/internal/database/compute"
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
	d.logger.Info("handling request [%s]", request)

	query, err := d.computeLayer.Parse(request)
	if err != nil {
		d.logger.Debug("compute layer is incorrect")
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
	default:
		return "", errors.New("unknown command")
	}

	return response, nil
}