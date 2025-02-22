package api

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/meesooqa/ttag/app/analysis"
	"github.com/meesooqa/ttag/app/repositories"
	"github.com/meesooqa/ttag/app/web/adapter"
)

type CooccMatrixD3ApiController struct {
	BaseApiController
	//provider analysis.AnalyzedDataProvider
	provider *analysis.CooccMatrixDataProvider
	adapter  adapter.D3DataAdapter
}

func NewCooccMatrixD3ApiController(log *slog.Logger, repo repositories.Repository) *CooccMatrixD3ApiController {
	c := &CooccMatrixD3ApiController{
		BaseApiController: BaseApiController{
			log:    log,
			method: http.MethodGet,
			route:  "/api/coocc_matrix_d3.json",
		},
		provider: analysis.NewCooccMatrixDataProvider(log, repo),
		adapter:  adapter.NewCooccMatrixD3DataAdapter(log),
	}
	c.self = c
	return c
}

func (c *CooccMatrixD3ApiController) fillData(r *http.Request) {
	group := r.URL.Query().Get("group")
	minFrequency, _ := strconv.Atoi(r.URL.Query().Get("min_frequency"))
	c.provider.SetMinFrequency(minFrequency)
	analysedData := c.provider.GetData(context.TODO(), group)
	c.data = c.adapter.PrepareData(analysedData)
}
