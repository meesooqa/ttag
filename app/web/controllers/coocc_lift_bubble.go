package controllers

import (
	"log/slog"
	"net/http"
)

type CooccLiftBubbleController struct {
	BaseController
}

// Bubble Chart (пузырьковая диаграмма)
func NewCooccLiftBubbleController(log *slog.Logger, tpl Template) *CooccLiftBubbleController {
	c := &CooccLiftBubbleController{BaseController{
		log:        log,
		tpl:        tpl,
		method:     http.MethodGet,
		route:      "/co-occ/lift-bubble/",
		title:      "Lift Bubble",
		contentTpl: "content/co-occ-lift-bubble.html",
	}}
	c.self = c
	return c
}

func (c *CooccLiftBubbleController) fillTemplateData(r *http.Request) {
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
