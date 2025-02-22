package api

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/meesooqa/ttag/app/analysis"
	"github.com/meesooqa/ttag/app/repositories"
	"github.com/meesooqa/ttag/app/web/adapter"
)

type CooccJaccardD3ApiController struct {
	BaseApiController
	provider analysis.AnalyzedDataProvider
	adapter  adapter.D3DataAdapter
}

func NewCooccJaccardD3ApiController(log *slog.Logger, repo repositories.Repository) *CooccJaccardD3ApiController {
	c := &CooccJaccardD3ApiController{
		BaseApiController: BaseApiController{
			log:    log,
			method: http.MethodGet,
			route:  "/api/coocc_jaccard_d3.json",
		},
		provider: analysis.NewCooccJaccardDataProvider(log, repo),
		adapter:  adapter.NewCooccJaccardD3DataAdapter(log),
	}
	c.self = c
	return c
}

func (c *CooccJaccardD3ApiController) fillData(r *http.Request) {
	group := r.URL.Query().Get("group")
	analysedData := c.provider.GetData(context.TODO(), group)
	c.data = c.adapter.PrepareData(analysedData)
}
