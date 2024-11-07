package main

import (
	"bytes"
	"flag"
	"log"
	"os"

	"github.com/sadovnikoff/GoConcurrencyCourse/homework_3/internal"
	"github.com/sadovnikoff/GoConcurrencyCourse/homework_3/internal/common"
)

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

	server, err := internal.Setup(cfg, logger)
	if err != nil {
		log.Fatal("error setting up a server:", err.Error())
	}

	server.Run()
}
