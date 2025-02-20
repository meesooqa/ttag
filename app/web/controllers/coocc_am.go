package controllers

import (
	"log/slog"
	"net/http"
)

type CooccAmController struct {
	BaseController
}

func NewCooccAmController(log *slog.Logger, tpl Template) *CooccAmController {
	c := &CooccAmController{BaseController{
		log:        log,
		tpl:        tpl,
		method:     http.MethodGet,
		route:      "/co-occ/association-measures/",
		title:      "Association measures",
		contentTpl: "content/co-occ-association-measures.html",
	}}
	c.self = c
	return c
}

func (c *CooccAmController) fillTemplateData(r *http.Request) {
	td, ok := c.tpl.GetData(r).(*DefaultTemplateData)
	if !ok {
		c.log.Error("template data invalid")
		return
	}
	// association-measures
	c.templateData = struct {
		Group  string
		Groups []GroupItem
		Menu   []MenuItem
		// Controller Vars
		Title      string
		CooccAmVar string
	}{
		Group:      td.Group,
		Groups:     td.Groups,
		Menu:       td.Menu,
		Title:      c.title,
		CooccAmVar: "CooccAmController::CooccAmVar",
	}
}
