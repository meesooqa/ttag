package api

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/meesooqa/ttag/app/analysis"
	"github.com/meesooqa/ttag/app/repositories"
	"github.com/meesooqa/ttag/app/web/adapter"
)

type CooccPmiD3ApiController struct {
	BaseApiController
	provider analysis.AnalyzedDataProvider
	adapter  adapter.D3DataAdapter
}

func NewCooccPmiD3ApiController(log *slog.Logger, repo repositories.Repository) *CooccPmiD3ApiController {
	c := &CooccPmiD3ApiController{
		BaseApiController: BaseApiController{
			log:    log,
			method: http.MethodGet,
			route:  "/api/coocc_pmi_d3.json",
		},
		provider: analysis.NewCooccPmiDataProvider(log, repo),
		adapter:  adapter.NewCooccPmiD3DataAdapter(log),
	}
	c.self = c
	return c
}

func (c *CooccPmiD3ApiController) fillData(r *http.Request) {
	group := r.URL.Query().Get("group")
	analysedData := c.provider.GetData(context.TODO(), group)
	c.data = c.adapter.PrepareData(analysedData)
}
