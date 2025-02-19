package controllers

import (
	"html/template"
	"log/slog"
	"net/http"
)

type Controller interface {
	Router(mux *http.ServeMux, templates *template.Template)
}

type TemplateFiller interface {
	fillTemplateData(r *http.Request)
}

type BaseController struct {
	self   TemplateFiller
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
	c.self.fillTemplateData(r)
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
