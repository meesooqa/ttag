package api

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/meesooqa/ttag/app/analysis"
	"github.com/meesooqa/ttag/app/repositories"
)

type CooccClustersD3ApiController struct {
	BaseApiController
	provider analysis.AnalyzedDataProvider
}

func NewCooccClustersD3ApiController(log *slog.Logger, repo repositories.Repository) *CooccClustersD3ApiController {
	c := &CooccClustersD3ApiController{
		BaseApiController: BaseApiController{
			log:    log,
			method: http.MethodGet,
			route:  "/api/coocc_clusters_d3.json",
		},
		provider: analysis.NewCooccClustersDataProvider(log, repo),
	}
	c.self = c
	return c
}

func (c *CooccClustersD3ApiController) fillData(r *http.Request) {
	c.data = c.provider.GetData(context.TODO(), r.URL.Query().Get("group"))
}
