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

type GroupItem struct {
	Title    string
	Value    string
	IsActive bool
}

type MenuItem struct {
	Title    string
	Link     string
	IsActive bool
	SubItems []MenuItem
}

func NewIndexController(log *slog.Logger, repo repositories.Repository) *IndexController {
	c := &IndexController{
		BaseController: BaseController{
			log:      log,
			method:   http.MethodGet,
			route:    "/",
			template: "content/home.html",
		},
		repo: repo,
	}
	c.self = c
	return c
}

func (c *IndexController) fillTemplateData(r *http.Request) {
	queryParams := r.URL.Query()
	group := queryParams.Get("group")
	c.templateData = struct {
		Title  string
		Groups []GroupItem
		Group  string
		Menu   []MenuItem
	}{
		Title:  "NewIndexController",
		Groups: c.getGroups(group),
		Group:  group,
		Menu:   c.getMenu(),
	}
}

func (c *IndexController) getGroups(group string) []GroupItem {
	items, err := c.repo.GetUniqueValues(context.TODO(), "group")
	if err != nil {
		c.log.Error(err.Error(), "err", err)
		return nil
	}

	result := make([]GroupItem, len(items))
	for i, item := range items {
		result[i] = GroupItem{
			Title:    item,
			Value:    item,
			IsActive: item == group,
		}
	}

	return result
}

func (c *IndexController) getMenu() []MenuItem {
	return []MenuItem{
		{
			Title: "Home",
			Link:  "/",
		},
		{
			Title: "Co-occurrence Analysis",
			Link:  "/co-occ/",
			SubItems: []MenuItem{
				{
					Title: "Frequency analysis of pairs",
					Link:  "/co-occ/pairs/",
				},
				{
					Title: "Association measures",
					Link:  "/co-occ/association-measures/",
				},
			},
		},
	}
}
