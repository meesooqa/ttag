package controllers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/meesooqa/ttag/app/repositories"
)

type Template interface {
	GetData(r *http.Request, title, route string) TemplateData
}

type TemplateData interface{}

type DefaultTemplate struct {
	log  *slog.Logger
	repo repositories.Repository
	data *DefaultTemplateData
}

type DefaultTemplateData struct {
	Title  string
	Group  string
	Groups []GroupItem
	Menu   []MenuItem
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

func NewDefaultTemplate(log *slog.Logger, repo repositories.Repository) *DefaultTemplate {
	return &DefaultTemplate{
		log:  log,
		repo: repo,
	}
}

func (c DefaultTemplate) GetData(r *http.Request, title, route string) TemplateData {
	queryParams := r.URL.Query()
	group := queryParams.Get("group")
	c.data = &DefaultTemplateData{
		Title:  title,
		Groups: c.getGroups(group),
		Menu:   c.getMenu(route),
		Group:  group,
	}
	return c.data
}

func (c DefaultTemplate) getGroups(group string) []GroupItem {
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

func (c DefaultTemplate) getMenu(menuLink string) []MenuItem {
	return []MenuItem{
		{
			Title:    "Home",
			Link:     "/",
			IsActive: false,
		},
		{
			Title:    "Co-occurrence Analysis",
			Link:     "/co-occ/",
			IsActive: true,
			SubItems: []MenuItem{
				{
					Title:    "Frequency analysis of pairs",
					Link:     "/co-occ/pairs/",
					IsActive: true,
				},
				{
					Title:    "Association measures",
					Link:     "/co-occ/association-measures/",
					IsActive: false,
				},
			},
		},
	}
}
