package controllers

import (
	"log/slog"
	"net/http"
)

type CooccController struct {
	BaseController
}

func NewCooccController(log *slog.Logger, tpl Template) *CooccController {
	c := &CooccController{BaseController{
		log:        log,
		tpl:        tpl,
		method:     http.MethodGet,
		route:      "/co-occ/",
		title:      "Co-occurrence Analysis",
		contentTpl: "content/co-occ.html",
	}}
	c.self = c
	c.AddChildren(
		NewCooccPairsController(log, tpl),
		NewCooccClustersController(log, tpl),
		NewCooccMatrixController(log, tpl),
		NewCooccLiftController(log, tpl),
		NewCooccAmController(log, tpl),
	)
	return c
}

func (c *CooccController) fillTemplateData(r *http.Request) {
	td, ok := c.tpl.GetData(r).(*DefaultTemplateData)
	if !ok {
		c.log.Error("template data invalid")
		return
	}
	c.templateData = struct {
		Group string
		Menu  []MenuItem
		// Controller Vars
		Title    string
		CooccVar string
	}{
		Group:    td.Group,
		Menu:     td.Menu,
		Title:    c.title,
		CooccVar: "CooccController::CooccVar",
	}
}
