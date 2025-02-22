package controllers

import (
	"log/slog"
	"net/http"
)

type CooccSupportController struct {
	BaseController
}

func NewCooccSupportController(log *slog.Logger, tpl Template) *CooccSupportController {
	c := &CooccSupportController{BaseController{
		log:        log,
		tpl:        tpl,
		method:     http.MethodGet,
		route:      "/co-occ/support/",
		title:      "Support Measure",
		contentTpl: "content/co-occ-lift-support.html",
	}}
	c.self = c
	return c
}

func (c *CooccSupportController) fillTemplateData(r *http.Request) {
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
