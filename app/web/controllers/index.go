package controllers

import (
	"log/slog"
	"net/http"
)

type IndexController struct {
	BaseController
}

func NewIndexController(log *slog.Logger, tpl Template) *IndexController {
	c := &IndexController{BaseController{
		log:        log,
		tpl:        tpl,
		method:     http.MethodGet,
		route:      "/",
		title:      "Home",
		contentTpl: "content/home.html",
	}}
	c.self = c
	return c
}

func (c *IndexController) fillTemplateData(r *http.Request) {
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
		IndexVar string
	}{
		Group:    td.Group,
		Menu:     td.Menu,
		Title:    c.title,
		IndexVar: "IndexController::IndexVar",
	}
}
