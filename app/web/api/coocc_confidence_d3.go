package api

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/meesooqa/ttag/app/analysis"
	"github.com/meesooqa/ttag/app/repositories"
	"github.com/meesooqa/ttag/app/web/adapter"
)

type CooccConfidenceD3ApiController struct {
	BaseApiController
	provider analysis.AnalyzedDataProvider
	adapter  adapter.D3DataAdapter
}

func NewCooccConfidenceD3ApiController(log *slog.Logger, repo repositories.Repository) *CooccConfidenceD3ApiController {
	c := &CooccConfidenceD3ApiController{
		BaseApiController: BaseApiController{
			log:    log,
			method: http.MethodGet,
			route:  "/api/coocc_confidence_d3.json",
		},
		provider: analysis.NewCooccConfidenceDataProvider(log, repo),
		adapter:  adapter.NewCooccConfidenceD3DataAdapter(log),
	}
	c.self = c
	return c
}

func (c *CooccConfidenceD3ApiController) fillData(r *http.Request) {
	group := r.URL.Query().Get("group")
	analysedData := c.provider.GetData(context.TODO(), group)
	c.data = c.adapter.PrepareData(analysedData)
}
