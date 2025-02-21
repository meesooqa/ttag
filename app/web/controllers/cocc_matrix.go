package controllers

import (
	"log/slog"
	"net/http"
)

type CooccMatrixController struct {
	BaseController
}

func NewCooccMatrixController(log *slog.Logger, tpl Template) *CooccMatrixController {
	c := &CooccMatrixController{BaseController{
		log:        log,
		tpl:        tpl,
		method:     http.MethodGet,
		route:      "/co-occ/matrix/",
		title:      "Co-occurrence Matrix",
		contentTpl: "content/co-occ-matrix.html",
	}}
	c.self = c
	return c
}

func (c *CooccMatrixController) fillTemplateData(r *http.Request) {
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
