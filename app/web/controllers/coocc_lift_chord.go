package controllers

import (
	"log/slog"
	"net/http"
)

type CooccLiftChordController struct {
	BaseController
}

// Chord Diagram (диаграмма хорд)
func NewCooccLiftChordController(log *slog.Logger, tpl Template) *CooccLiftChordController {
	c := &CooccLiftChordController{BaseController{
		log:        log,
		tpl:        tpl,
		method:     http.MethodGet,
		route:      "/co-occ/lift-chord/",
		title:      "Lift Chord",
		contentTpl: "content/co-occ-lift-chord.html",
	}}
	c.self = c
	return c
}

func (c *CooccLiftChordController) fillTemplateData(r *http.Request) {
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
