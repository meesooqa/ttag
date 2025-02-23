package api

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/meesooqa/ttag/app/analysis"
	"github.com/meesooqa/ttag/app/repositories"
)

type CooccSupportD3ApiController struct {
	BaseApiController
	provider analysis.AnalyzedDataProvider
}

func NewCooccSupportD3ApiController(log *slog.Logger, repo repositories.Repository) *CooccSupportD3ApiController {
	c := &CooccSupportD3ApiController{
		BaseApiController: BaseApiController{
			log:    log,
			method: http.MethodGet,
			route:  "/api/coocc_support_d3.json",
		},
		provider: analysis.NewCooccSupportDataProvider(log, repo),
	}
	c.self = c
	return c
}

func (c *CooccSupportD3ApiController) fillData(r *http.Request) {
	c.data = c.provider.GetData(context.TODO(), r.URL.Query().Get("group"))
}

type CooccPmiD3ApiController struct {
	BaseApiController
	provider analysis.AnalyzedDataProvider
}

func NewCooccPmiD3ApiController(log *slog.Logger, repo repositories.Repository) *CooccPmiD3ApiController {
	c := &CooccPmiD3ApiController{
		BaseApiController: BaseApiController{
			log:    log,
			method: http.MethodGet,
			route:  "/api/coocc_pmi_d3.json",
		},
		provider: analysis.NewCooccPmiDataProvider(log, repo),
	}
	c.self = c
	return c
}

func (c *CooccPmiD3ApiController) fillData(r *http.Request) {
	c.data = c.provider.GetData(context.TODO(), r.URL.Query().Get("group"))
}

type CooccLiftD3ApiController struct {
	BaseApiController
	provider analysis.AnalyzedDataProvider
}

func NewCooccLiftD3ApiController(log *slog.Logger, repo repositories.Repository) *CooccLiftD3ApiController {
	c := &CooccLiftD3ApiController{
		BaseApiController: BaseApiController{
			log:    log,
			method: http.MethodGet,
			route:  "/api/coocc_lift_d3.json",
		},
		provider: analysis.NewCooccLiftDataProvider(log, repo),
	}
	c.self = c
	return c
}

func (c *CooccLiftD3ApiController) fillData(r *http.Request) {
	c.data = c.provider.GetData(context.TODO(), r.URL.Query().Get("group"))
}

type CooccJaccardD3ApiController struct {
	BaseApiController
	provider analysis.AnalyzedDataProvider
}

func NewCooccJaccardD3ApiController(log *slog.Logger, repo repositories.Repository) *CooccJaccardD3ApiController {
	c := &CooccJaccardD3ApiController{
		BaseApiController: BaseApiController{
			log:    log,
			method: http.MethodGet,
			route:  "/api/coocc_jaccard_d3.json",
		},
		provider: analysis.NewCooccJaccardDataProvider(log, repo),
	}
	c.self = c
	return c
}

func (c *CooccJaccardD3ApiController) fillData(r *http.Request) {
	c.data = c.provider.GetData(context.TODO(), r.URL.Query().Get("group"))
}

type CooccConfidenceD3ApiController struct {
	BaseApiController
	provider analysis.AnalyzedDataProvider
}

func NewCooccConfidenceD3ApiController(log *slog.Logger, repo repositories.Repository) *CooccConfidenceD3ApiController {
	c := &CooccConfidenceD3ApiController{
		BaseApiController: BaseApiController{
			log:    log,
			method: http.MethodGet,
			route:  "/api/coocc_confidence_d3.json",
		},
		provider: analysis.NewCooccConfidenceDataProvider(log, repo),
	}
	c.self = c
	return c
}

func (c *CooccConfidenceD3ApiController) fillData(r *http.Request) {
	c.data = c.provider.GetData(context.TODO(), r.URL.Query().Get("group"))
}
