package controllers

import (
	"log/slog"
	"net/http"
)

type CooccPairsController struct {
	BaseController
}

func NewCooccPairsController(log *slog.Logger, tpl Template) *CooccPairsController {
	c := &CooccPairsController{BaseController{
		log:        log,
		tpl:        tpl,
		method:     http.MethodGet,
		route:      "/co-occ/pairs/",
		title:      "Frequency analysis of pairs",
		contentTpl: "content/co-occ-pairs.html",
	}}
	c.self = c
	return c
}

func (c *CooccPairsController) fillTemplateData(r *http.Request) {
	td, ok := c.tpl.GetData(r).(*DefaultTemplateData)
	if !ok {
		c.log.Error("template data invalid")
		return
	}
	c.templateData = struct {
		Group  string
		Groups []GroupItem
		Menu   []MenuItem
		// Controller Vars
		Title         string
		CooccPairsVar string
	}{
		Group:         td.Group,
		Groups:        td.Groups,
		Menu:          td.Menu,
		Title:         c.title,
		CooccPairsVar: "CooccPairsController::CooccPairsVar",
	}
}
