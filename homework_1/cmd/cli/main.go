package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/sadovnikoff/GoConcurrencyCourse/homework_1/internal/common"
	"github.com/sadovnikoff/GoConcurrencyCourse/homework_1/internal/database"
	"github.com/sadovnikoff/GoConcurrencyCourse/homework_1/internal/database/compute"
	"github.com/sadovnikoff/GoConcurrencyCourse/homework_1/internal/database/storage"
	"github.com/sadovnikoff/GoConcurrencyCourse/homework_1/internal/database/storage/engine"
)

func setup(logger *common.Logger) (*database.Database, error) {
	computeLayer, err := compute.NewParser(logger)
	if err != nil {
		return nil, err
	}

	dbEngine, err := engine.NewEngine(logger)
	if err != nil {
		return nil, err
	}

	storageLayer, err := storage.NewStorage(dbEngine, logger)
	if err != nil {
		return nil, err
	}

	db, err := database.NewDatabase(computeLayer, storageLayer, logger)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	fmt.Println(" In-memory key-value DB server is running.\nType a command with arguments and press Enter.")
	fmt.Println("Available commands: SET, GET, DEL")

	logger := common.NewLogger()

	db, err := setup(logger)
	if err != nil {
		logger.ELog.Printf("error setting up database: %s", err.Error())
		return
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			logger.ELog.Println(err.Error())
			return
		}

		request := strings.TrimSpace(input)
		response, err := db.HandleQuery(request)
		if err != nil {
			logger.ELog.Println(err.Error())
			continue
		}

		logger.DLog.Println(response)
	}
}
