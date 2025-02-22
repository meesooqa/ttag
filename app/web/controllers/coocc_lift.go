package controllers

import (
	"log/slog"
	"net/http"
)

type CooccLiftController struct {
	BaseController
}

// Меры ассоциации: Lift
// Lift: Отношение наблюдаемой совместной частоты к ожидаемой при независимости появления тегов.
// Lift: отношение совместной вероятности появления тегов к произведению их индивидуальных вероятностей.
// Граф взаимосвязей тегов (метрика Lift)
func NewCooccLiftController(log *slog.Logger, tpl Template) *CooccLiftController {
	c := &CooccLiftController{BaseController{
		log:        log,
		tpl:        tpl,
		method:     http.MethodGet,
		route:      "/co-occ/lift/",
		title:      "Lift Measure",
		contentTpl: "content/co-occ-lift.html",
	}}
	c.self = c
	return c
}

func (c *CooccLiftController) fillTemplateData(r *http.Request) {
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
