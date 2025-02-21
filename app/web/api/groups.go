package api

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/meesooqa/ttag/app/repositories"
)

type GroupsApiController struct {
	BaseApiController
	repo repositories.Repository
}

func NewGroupsApiController(log *slog.Logger, repo repositories.Repository) *GroupsApiController {
	c := &GroupsApiController{
		BaseApiController: BaseApiController{
			log:    log,
			method: http.MethodGet,
			route:  "/api/groups.json",
		},
		repo: repo,
	}
	c.self = c
	return c
}

func (c *GroupsApiController) fillData(r *http.Request) {
	items, err := c.repo.GetUniqueValues(context.TODO(), "group")
	if err != nil {
		c.log.Error("DefaultTemplate group getting", "err", err)
	}
	c.data = items
}
