package controllers

import (
	"html/template"
	"log/slog"
	"net/http"
)

type Controller interface {
	Router(mux *http.ServeMux, templates *template.Template)
}

type BaseController struct {
	log    *slog.Logger
	method string
	route  string

	template     string
	templateData any

	templates *template.Template
}

func (c *BaseController) Router(mux *http.ServeMux, templates *template.Template) {
	c.templates = templates
	//mux.HandleFunc(fmt.Sprintf("%s %s", c.method, c.route), c.handler)
	mux.HandleFunc(c.route, c.handler)
}

func (c *BaseController) handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != c.method {
		c.log.Error("method is not allowed", "method", c.method, "route", c.route)
		http.Error(w, "method is not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := c.templates.ExecuteTemplate(w, c.template, &c.templateData); err != nil {
		c.log.Error("executing template", "err", err, "template", c.template, "route", c.route)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type IndexController struct {
	BaseController
}

func NewIndexController(log *slog.Logger) *IndexController {
	groups := map[string]string{
		"1": "Group Name 100",
		"2": "Links Group Name",
		"3": "Group 300",
	}
	return &IndexController{BaseController{
		log:      log,
		method:   http.MethodGet,
		route:    "/",
		template: "index.html",
		templateData: struct {
			Title   string
			Groups  map[string]string
			Group   string
			GroupId string
		}{
			Title:   "NewIndexController",
			Groups:  groups,
			Group:   "",
			GroupId: "",
		},
	}}
}
