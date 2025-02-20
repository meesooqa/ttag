package controllers

import (
	"html/template"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"os"
)

type Controller interface {
	Router(mux *http.ServeMux)
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
	templateData TemplateData

	templates *template.Template
}

func (c *BaseController) Router(mux *http.ServeMux) {
	c.initTemplates()
	mux.HandleFunc(c.route, c.handler)
}

func (c *BaseController) handler(w http.ResponseWriter, r *http.Request) {
	c.self.fillTemplateData(r)
	if r.Method != c.method {
		c.log.Error("method is not allowed", "method", c.method, "route", c.route)
		http.Error(w, "method is not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := c.templates.ExecuteTemplate(w, "index.html", &c.templateData); err != nil {
		c.log.Error("executing template", "err", err, "template", c.template, "route", c.route)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *BaseController) initTemplates() {
	if c.template == "" {
		c.template = "content/default.html"
	}

	var files []string

	tplDir := "app/web/templates"
	topLevel, err := fs.Glob(os.DirFS(tplDir), "*.html")
	if err != nil {
		c.log.Error("find tpls - topLevel", "err", err)
	}
	files = append(files, topLevel...)

	subDir, err := fs.Glob(os.DirFS(tplDir), c.template)
	if err != nil {
		c.log.Error("find tpls - subDir", "err", err)
		log.Fatal(err)
	}
	files = append(files, subDir...)

	for i, f := range files {
		files[i] = tplDir + "/" + f
	}

	c.templates, err = template.ParseFiles(files...)
	if err != nil {
		c.log.Error("parse tpls", "err", err)
	}
}
