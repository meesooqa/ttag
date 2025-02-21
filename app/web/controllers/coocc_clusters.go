package controllers

import (
	"log/slog"
	"net/http"
)

type CooccClustersController struct {
	BaseController
}

// Кластеризация тегов
func NewCooccClustersController(log *slog.Logger, tpl Template) *CooccClustersController {
	c := &CooccClustersController{BaseController{
		log:        log,
		tpl:        tpl,
		method:     http.MethodGet,
		route:      "/co-occ/clusters/",
		title:      "Tag Clustering",
		contentTpl: "content/co-occ-clusters.html",
	}}
	c.self = c
	return c
}

func (c *CooccClustersController) fillTemplateData(r *http.Request) {
	td, ok := c.tpl.GetData(r).(*DefaultTemplateData)
	if !ok {
		c.log.Error("template data invalid")
		return
	}
	c.templateData = struct {
		Group string
		Menu  []MenuItem
		// Controller Vars
		Title string
	}{
		Group: td.Group,
		Menu:  td.Menu,
		Title: c.title,
	}
}
