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
	c.templateData = c.tpl.GetData(r)
	// TODO IndexVar
}
