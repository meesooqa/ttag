package controllers

import (
	"log/slog"
	"net/http"
)

type IndexController struct {
	BaseController
}

func NewIndexController(log *slog.Logger) *IndexController {
	ic := &IndexController{BaseController{
		log:      log,
		method:   http.MethodGet,
		route:    "/",
		template: "index.html",
	}}
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
	return map[string]string{
		"1": "Group Name 100",
		"2": "Links Group Name",
		"3": "Group 300",
	}
}
