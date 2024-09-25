package internal

import (
	"errors"
	"fmt"

	"github.com/sadovnikoff/GoConcurrencyCourse/homework_2/internal/common"
	"github.com/sadovnikoff/GoConcurrencyCourse/homework_2/internal/database"
	"github.com/sadovnikoff/GoConcurrencyCourse/homework_2/internal/database/compute"
	"github.com/sadovnikoff/GoConcurrencyCourse/homework_2/internal/database/storage"
	"github.com/sadovnikoff/GoConcurrencyCourse/homework_2/internal/database/storage/engine"
	"github.com/sadovnikoff/GoConcurrencyCourse/homework_2/internal/network/tcp"
)

type NetworkLayer interface {
	Run()
}

func Setup(cfg *common.Config, logger *common.Logger) (NetworkLayer, error) {
	if cfg == nil {
		return nil, errors.New("config is nil")
	}

	computeLayer, err := compute.NewParser(logger)
	if err != nil {
		return nil, err
	}

	var dbEngine storage.Engine
	if cfg.Engine.Type == engine.InMemoryEngine {
		logger.Debug("setup server: in-memory engine has been chosen")
		dbEngine, err = engine.NewEngine(logger)
		if err != nil {
			logger.Debug("setup server: in-memory engine cannot be set up")
			return nil, err
		}
	} else {
		logger.Debug("setup server: engine [%s] not supported", cfg.Engine.Type)
		return nil, fmt.Errorf("engine type '%s' not supported", cfg.Engine.Type)
	}

	storageLayer, err := storage.NewStorage(dbEngine, logger)
	if err != nil {
		logger.Debug("setup server: storage cannot be set up")
		return nil, err
	}

	db, err := database.NewDatabase(computeLayer, storageLayer, logger)
	if err != nil {
		logger.Debug("setup server: database cannot be set up")
		return nil, err
	}

	server, err := tcp.NewServer(cfg, db, logger)
	if err != nil {
		logger.Debug("setup server: tcp server cannot be set up")
		return nil, err
	}

	return server, nil
}
