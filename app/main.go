package main

import (
	"github.com/meesooqa/ttag/app/fs"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
)

var logger *zap.Logger

func main() {
	var wg sync.WaitGroup

	initLogger()

	path := "var/data/test"
	filesChan := make(chan string, 2)

	logger.Info("Start", zap.String("path", path))

	wg.Add(1)
	finder := fs.NewFinder(logger)
	go finder.FindFiles(path, filesChan, &wg)

	wg.Add(1)
	go processFile(filesChan, &wg)

	wg.Wait()
	logger.Info("All goroutines are done")
}

func processFile(filesChan <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for filePath := range filesChan {
		logger.Info("Got file:", zap.String("filePath", filePath))
		Worker(filePath)
	}
}

func Worker(filename string) {
	// process File Worker
	time.Sleep(1 * time.Second)
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
