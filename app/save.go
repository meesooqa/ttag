package main

import (
	"sync"

	"go.uber.org/zap"

	"github.com/meesooqa/ttag/app/config"
	"github.com/meesooqa/ttag/app/db"
	"github.com/meesooqa/ttag/app/fs"
	"github.com/meesooqa/ttag/app/proc"
	"github.com/meesooqa/ttag/app/repositories"
	"github.com/meesooqa/ttag/app/tg"
)

func main() {
	var wg sync.WaitGroup
	initLogger()
	logger.Info("start saving")
	conf, err := config.Load("etc/config.yml")
	if err != nil {
		logger.Error("can't load config", zap.Error(err))
	}

	wg.Add(1)
	filesChan := make(chan string, 2)
	finder := fs.NewFinder(logger)
	go finder.FindFiles(conf.System.DataPath, filesChan, &wg)

	wg.Add(1)
	mongoDB := db.NewMongoDB(logger, conf.Mongo)
	err = mongoDB.Init()
	if err != nil {
		logger.Error("db connection failed", zap.Error(err))
	}
	defer mongoDB.Close()

	repo := repositories.NewMessageRepository(logger, mongoDB.GetCollectionMessages())
	tgService := tg.NewService(logger, conf.System)
	processor := proc.NewProcessor(logger, tgService, repo)
	go processor.ProcessFile(filesChan, &wg)

	wg.Wait()
	logger.Info("all goroutines are done")
}
