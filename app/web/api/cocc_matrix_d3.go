package api

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/meesooqa/ttag/app/analysis"
	"github.com/meesooqa/ttag/app/repositories"
)

type CooccMatrixD3ApiController struct {
	BaseApiController
	//provider analysis.AnalyzedDataProvider
	provider *analysis.CooccMatrixDataProvider
}

func NewCooccMatrixD3ApiController(log *slog.Logger, repo repositories.Repository) *CooccMatrixD3ApiController {
	c := &CooccMatrixD3ApiController{
		BaseApiController: BaseApiController{
			log:    log,
			method: http.MethodGet,
			route:  "/api/coocc_matrix_d3.json",
		},
		provider: analysis.NewCooccMatrixDataProvider(log, repo),
	}
	c.self = c
	return c
}

func (c *CooccMatrixD3ApiController) fillData(r *http.Request) {
	minFrequency, _ := strconv.Atoi(r.URL.Query().Get("min_frequency"))
	c.provider.SetMinFrequency(minFrequency)
	c.data = c.provider.GetData(context.TODO(), r.URL.Query().Get("group"))
}
