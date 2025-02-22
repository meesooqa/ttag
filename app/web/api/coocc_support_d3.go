package api

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/meesooqa/ttag/app/analysis"
	"github.com/meesooqa/ttag/app/repositories"
	"github.com/meesooqa/ttag/app/web/adapter"
)

type CooccSupportD3ApiController struct {
	BaseApiController
	provider analysis.AnalyzedDataProvider
	adapter  adapter.D3DataAdapter
}

func NewCooccSupportD3ApiController(log *slog.Logger, repo repositories.Repository) *CooccSupportD3ApiController {
	c := &CooccSupportD3ApiController{
		BaseApiController: BaseApiController{
			log:    log,
			method: http.MethodGet,
			route:  "/api/coocc_support_d3.json",
		},
		provider: analysis.NewCooccSupportDataProvider(log, repo),
		adapter:  adapter.NewCooccSupportD3DataAdapter(log),
	}
	c.self = c
	return c
}

func (c *CooccSupportD3ApiController) fillData(r *http.Request) {
	group := r.URL.Query().Get("group")
	analysedData := c.provider.GetData(context.TODO(), group)
	c.data = c.adapter.PrepareData(analysedData)
}
