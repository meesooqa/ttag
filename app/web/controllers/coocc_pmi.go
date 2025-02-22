package controllers

import (
	"log/slog"
	"net/http"
)

type CooccPmiController struct {
	BaseController
}

func NewCooccPmiController(log *slog.Logger, tpl Template) *CooccPmiController {
	c := &CooccPmiController{BaseController{
		log:        log,
		tpl:        tpl,
		method:     http.MethodGet,
		route:      "/co-occ/pmi/",
		title:      "Pointwise Mutual Information",
		contentTpl: "content/co-occ-pmi.html",
	}}
	c.self = c
	return c
}

func (c *CooccPmiController) fillTemplateData(r *http.Request) {
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
