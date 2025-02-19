package controllers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/meesooqa/ttag/app/repositories"
)

type IndexController struct {
	BaseController
	repo repositories.Repository
}

func NewIndexController(log *slog.Logger, repo repositories.Repository) *IndexController {
	ic := &IndexController{
		BaseController: BaseController{
			log:      log,
			method:   http.MethodGet,
			route:    "/",
			template: "index.html",
		},
		repo: repo,
	}
	ic.self = ic
	return ic
}

func (c *IndexController) fillTemplateData(r *http.Request) {
	groups := c.getGroups()
	queryParams := r.URL.Query()
	groupId := queryParams.Get("group")
	group := groups[groupId]

	c.templateData = struct {
		Title   string
		Groups  map[string]string
		Group   string
		GroupId string
	}{
		Title:   "NewIndexController",
		Groups:  groups,
		Group:   group,
		GroupId: groupId,
	}
}

func (c *IndexController) getGroups() map[string]string {
	items, err := c.repo.GetUniqueValues(context.TODO(), "group")
	if err != nil {
		c.log.Error(err.Error(), "err", err)
		return nil
	}

	result := make(map[string]string, len(items))
	for _, item := range items {
		result[item] = item
	}

	return result
}
