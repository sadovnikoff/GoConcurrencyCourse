package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/sadovnikoff/GoConcurrencyCourse/homework_2/internal/common"
	"github.com/sadovnikoff/GoConcurrencyCourse/homework_2/internal/database"
	"github.com/sadovnikoff/GoConcurrencyCourse/homework_2/internal/database/compute"
	"github.com/sadovnikoff/GoConcurrencyCourse/homework_2/internal/database/storage"
	"github.com/sadovnikoff/GoConcurrencyCourse/homework_2/internal/database/storage/engine"
	"github.com/sadovnikoff/GoConcurrencyCourse/homework_2/internal/network/tcp"
)

type networkLayer interface {
	Run()
}

func setup(cfg *common.Config, logger *common.Logger) (networkLayer, error) {
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

func main() {
	cfgPath := flag.String("cfg_path", "config.yaml", "Config file path")
	flag.Parse()

	cfg := &common.Config{}
	if *cfgPath != "" {
		data, err := os.ReadFile(*cfgPath)
		if err != nil {
			log.Fatal(err)
		}

		reader := bytes.NewReader(data)
		cfg, err = common.ParseConfig(reader)
		if err != nil {
			log.Fatal(err)
		}
	}

	logger, err := common.NewLogger(cfg.Logging.Level, cfg.Logging.Output)
	if err != nil {
		log.Fatal("error creating logger:", err.Error())
	}
	defer func() { _ = logger.Close() }()

	server, err := setup(cfg, logger)
	if err != nil {
		log.Fatal("error setting up a server:", err.Error())
	}

	server.Run()
}
