package controllers

import (
	"log/slog"
	"net/http"
)

type CooccConfidenceController struct {
	BaseController
}

func NewCooccConfidenceController(log *slog.Logger, tpl Template) *CooccConfidenceController {
	c := &CooccConfidenceController{BaseController{
		log:        log,
		tpl:        tpl,
		method:     http.MethodGet,
		route:      "/co-occ/confidence/",
		title:      "Confidence",
		contentTpl: "content/co-occ-confidence.html",
	}}
	c.self = c
	return c
}

func (c *CooccConfidenceController) fillTemplateData(r *http.Request) {
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
