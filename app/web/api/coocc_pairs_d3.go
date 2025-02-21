package api

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/meesooqa/ttag/app/analysis"
	"github.com/meesooqa/ttag/app/repositories"
	"github.com/meesooqa/ttag/app/web/adapter"
)

type CooccPairsD3ApiController struct {
	BaseApiController
	provider analysis.AnalyzedDataProvider
	adapter  adapter.D3DataAdapter
}

func NewCooccPairsD3ApiController(log *slog.Logger, repo repositories.Repository) *CooccPairsD3ApiController {
	c := &CooccPairsD3ApiController{
		BaseApiController: BaseApiController{
			log:    log,
			method: http.MethodGet,
			route:  "/api/coocc_pairs_d3.json",
		},
		provider: analysis.NewCooccPairsDataProvider(log, repo),
		adapter:  adapter.NewCooccPairsD3DataAdapter(log),
	}
	c.self = c
	return c
}

func (c *CooccPairsD3ApiController) fillData(r *http.Request) {
	group := r.URL.Query().Get("group")
	analysedData := c.provider.GetData(context.TODO(), group)
	c.data = c.adapter.PrepareData(analysedData)
}
