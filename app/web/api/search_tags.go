package api

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/meesooqa/ttag/app/repositories"
)

type SearchTagsApiController struct {
	BaseApiController
	repo repositories.Repository
}

func NewSearchTagsApiController(log *slog.Logger, repo repositories.Repository) *SearchTagsApiController {
	c := &SearchTagsApiController{
		BaseApiController: BaseApiController{
			log:    log,
			method: http.MethodGet,
			route:  "/api/search_tags.json",
		},
		repo: repo,
	}
	c.self = c
	return c
}

func (c *SearchTagsApiController) fillData(r *http.Request) {
	query := strings.ToLower(r.URL.Query().Get("q"))
	items, err := c.repo.GetTags(context.TODO(), query)
	if err != nil {
		c.log.Error("DefaultTemplate SearchTags getting", "err", err)
	}
	c.data = items
}
