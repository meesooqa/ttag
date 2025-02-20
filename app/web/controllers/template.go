package controllers

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/meesooqa/ttag/app/repositories"
)

type Template interface {
	GetTemplatesLocation() string
	GetStaticLocation() string
	GetMainTpl() string
	GetDefaultContentTpl() string
	GetData(r *http.Request) TemplateData
}

type TemplateData interface{}

type DefaultTemplate struct {
	code            string
	log             *slog.Logger
	repo            repositories.Repository
	menuControllers []Controller
	data            *DefaultTemplateData
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
		code: "default",
		log:  log,
		repo: repo,
	}
}

func (t *DefaultTemplate) GetTemplatesLocation() string {
	return "templates/" + t.code
}

func (t *DefaultTemplate) GetStaticLocation() string {
	return t.GetTemplatesLocation() + "/static"
}

func (t *DefaultTemplate) GetMainTpl() string {
	return "layout.html"
}

func (t *DefaultTemplate) GetDefaultContentTpl() string {
	return "content/default.html"
}

func (t *DefaultTemplate) getDefaultTitle() string {
	return "ttag"
}

func (t *DefaultTemplate) SetMenuControllers(menuControllers []Controller) {
	t.menuControllers = menuControllers
}

func (t *DefaultTemplate) GetData(r *http.Request) TemplateData {
	queryParams := r.URL.Query()
	group := queryParams.Get("group")
	t.data = &DefaultTemplateData{
		Title:  t.getTitle(r.URL.Path),
		Groups: t.getGroups(group),
		Menu:   t.getMenu(r.URL.Path),
		Group:  group,
	}
	return t.data
}

func (t *DefaultTemplate) getTitle(current string) string {
	for _, c := range t.menuControllers {
		if t.isMenuLinkCurrent(current, c.GetRoute()) {
			return c.GetTitle()
		}
	}
	return t.getDefaultTitle()
}

func (t *DefaultTemplate) getGroups(group string) []GroupItem {
	items, err := t.repo.GetUniqueValues(context.TODO(), "group")
	if err != nil {
		t.log.Error(err.Error(), "err", err)
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

func (t *DefaultTemplate) getMenu(current string) []MenuItem {
	// TODO SubItems
	var result []MenuItem
	for _, c := range t.menuControllers {
		mi := MenuItem{
			Title:    c.GetTitle(),
			Link:     c.GetRoute(),
			IsActive: t.isMenuLinkCurrent(current, c.GetRoute()),
		}

		result = append(result, mi)
	}
	return result
	/*/
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
	/*/
}

func (t *DefaultTemplate) isMenuLinkCurrent(current, link string) bool {
	if link == "/" {
		return current == "/"
	}
	return strings.HasPrefix(current, link)
}
