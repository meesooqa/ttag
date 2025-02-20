package controllers

import (
	"log/slog"
	"net/http"
)

type IndexController struct {
	BaseController
	tpl Template
}

func NewIndexController(log *slog.Logger, tpl Template) *IndexController {
	c := &IndexController{
		BaseController: BaseController{
			log:      log,
			method:   http.MethodGet,
			route:    "/",
			template: "content/home.html",
		},
		tpl: tpl,
	}
	c.self = c
	return c
}

func (c *IndexController) fillTemplateData(r *http.Request) {
	c.templateData = c.tpl.GetData(r, "Home", c.route)
}
