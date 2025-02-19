package main

import (
	"context"
	"log/slog"

	"github.com/meesooqa/ttag/app/config"
	"github.com/meesooqa/ttag/app/web"
	"github.com/meesooqa/ttag/app/web/controllers"
)

func main() {
	logger, cleanup := config.InitLogger("var/log/main.log", slog.LevelDebug) // LevelDebug
	defer cleanup()

	//conf, err := config.Load("etc/config.yml")
	//if err != nil {
	//	logger.Error("can't load config", "err", err)
	//}

	server := web.NewServer(logger, []controllers.Controller{
		controllers.NewIndexController(logger),
	})
	server.Run(context.Background(), 8080)
}
