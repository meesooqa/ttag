package main

import (
	"context"
	"log/slog"

	"github.com/meesooqa/ttag/app/config"
	"github.com/meesooqa/ttag/app/db"
	"github.com/meesooqa/ttag/app/repositories"
	"github.com/meesooqa/ttag/app/web"
	"github.com/meesooqa/ttag/app/web/controllers"
)

func main() {
	logger, cleanup := config.InitLogger("var/log/main.log", slog.LevelDebug)
	defer cleanup()

	conf, err := config.Load("etc/config.yml")
	if err != nil {
		logger.Error("can't load config", "err", err)
	}

	mongoDB := db.NewMongoDB(logger, conf.Mongo)
	err = mongoDB.Init()
	if err != nil {
		logger.Error("db connection failed", "err", err)
	}
	defer mongoDB.Close()

	repo := repositories.NewMessageRepository(logger, mongoDB)
	server := web.NewServer(logger, []controllers.Controller{
		controllers.NewIndexController(logger, repo),
	})
	server.Run(context.Background(), 8080)
}
