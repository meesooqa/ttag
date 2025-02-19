package main

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/meesooqa/ttag/app/config"
	"github.com/meesooqa/ttag/app/db"
	"github.com/meesooqa/ttag/app/fs"
	"github.com/meesooqa/ttag/app/proc"
	"github.com/meesooqa/ttag/app/repositories"
	"github.com/meesooqa/ttag/app/tg"
)

var logger *zap.Logger

func main() {
	var wg sync.WaitGroup
	initLogger()
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

func initLogger() {
	logFile, _ := os.OpenFile("var/log/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	writer := zapcore.AddSync(logFile)
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:      "timestamp",
		LevelKey:     "level",
		MessageKey:   "message",
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		writer,
		zap.DebugLevel, // zap.InfoLevel
	)

	logger = zap.New(core)
	defer logger.Sync()
}
