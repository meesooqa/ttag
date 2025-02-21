package api

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/meesooqa/ttag/app/analysis"
	"github.com/meesooqa/ttag/app/repositories"
	"github.com/meesooqa/ttag/app/web/adapter"
)

type CooccClustersD3ApiController struct {
	BaseApiController
	provider analysis.AnalyzedDataProvider
	adapter  adapter.D3DataAdapter
}

func NewCooccClustersD3ApiController(log *slog.Logger, repo repositories.Repository) *CooccClustersD3ApiController {
	c := &CooccClustersD3ApiController{
		BaseApiController: BaseApiController{
			log:    log,
			method: http.MethodGet,
			route:  "/api/coocc_clusters_d3.json",
		},
		provider: analysis.NewCooccClustersDataProvider(log, repo),
		adapter:  adapter.NewCooccClustersD3DataAdapter(log),
	}
	c.self = c
	return c
}

func (c *CooccClustersD3ApiController) fillData(r *http.Request) {
	group := r.URL.Query().Get("group")
	analysedData := c.provider.GetData(context.TODO(), group)
	c.data = c.adapter.PrepareData(analysedData)
}
