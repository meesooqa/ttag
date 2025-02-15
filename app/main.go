package main

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/meesooqa/ttag/app/db"
	"github.com/meesooqa/ttag/app/fs"
	"github.com/meesooqa/ttag/app/proc"
	"github.com/meesooqa/ttag/app/tg"
)

var logger *zap.Logger

func main() {
	var wg sync.WaitGroup

	initLogger()

	path := "var/data" // @see TgArchivedHTMLParser.ParseFile()
	filesChan := make(chan string, 2)

	logger.Info("Start", zap.String("path", path))

	wg.Add(1)
	finder := fs.NewFinder(logger)
	go finder.FindFiles(path, filesChan, &wg)

	wg.Add(1)
	tgService := tg.NewService(logger)
	mongoDB := db.NewMongoDB(logger, "db_tags", "tags")
	processor := proc.NewProcessor(logger, tgService, mongoDB)
	go processor.ProcessFile(filesChan, &wg)

	wg.Wait()
	logger.Info("All goroutines are done")
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
		zapcore.NewConsoleEncoder(encoderConfig),
		//TODO zapcore.NewJSONEncoder(encoderConfig),
		writer,
		zap.DebugLevel, // zap.InfoLevel
	)

	logger = zap.New(core)
	defer logger.Sync()
}
