package main

import (
	"context"
	"log/slog"

	"github.com/meesooqa/ttag/app/config"
	"github.com/meesooqa/ttag/app/db"
	"github.com/meesooqa/ttag/app/repositories"
	"github.com/meesooqa/ttag/app/web"
	"github.com/meesooqa/ttag/app/web/api"
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
	cca := provideApiControllers(logger, repo)
	tplDefault := controllers.NewDefaultTemplate(logger, repo)
	cc := provideControllers(logger, tplDefault)
	tplDefault.SetMenuControllers(cc)
	server := web.NewServer(logger, cca, cc)
	server.Run(context.Background(), 8080, tplDefault.GetStaticLocation())
}

func provideApiControllers(log *slog.Logger, repo repositories.Repository) []api.ApiController {
	return []api.ApiController{
		api.NewGroupsApiController(log, repo),
		api.NewSearchTagsApiController(log, repo),
		api.NewCooccPairsD3ApiController(log, repo),
		api.NewCooccClustersD3ApiController(log, repo),
		api.NewCooccMatrixD3ApiController(log, repo),
		api.NewCooccLiftD3ApiController(log, repo),
		api.NewCooccSupportD3ApiController(log, repo),
		api.NewCooccJaccardD3ApiController(log, repo),
		api.NewCooccPmiD3ApiController(log, repo),
	}
}

func provideControllers(log *slog.Logger, tpl controllers.Template) []controllers.Controller {
	return []controllers.Controller{
		controllers.NewIndexController(log, tpl),
		controllers.NewCooccController(log, tpl),
	}
}
