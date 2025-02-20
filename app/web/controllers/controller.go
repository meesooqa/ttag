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
	GetRoute() string
	GetTitle() string
	GetChildren() []Controller
	AddChildren(cc ...Controller)
}

type TemplateFiller interface {
	fillTemplateData(r *http.Request)
}

type BaseController struct {
	self         TemplateFiller
	log          *slog.Logger
	tpl          Template
	method       string
	route        string
	title        string
	contentTpl   string
	templateData TemplateData
	templates    *template.Template
	children     []Controller
}

func (c *BaseController) Router(mux *http.ServeMux) {
	c.initTemplates()
	mux.HandleFunc(c.route, c.handler)
}

func (c *BaseController) GetRoute() string {
	return c.route
}

func (c *BaseController) GetTitle() string {
	return c.title
}

func (c *BaseController) GetChildren() []Controller {
	return c.children
}

func (c *BaseController) AddChildren(cc ...Controller) {
	c.children = append(c.children, cc...)
}

func (c *BaseController) handler(w http.ResponseWriter, r *http.Request) {
	c.self.fillTemplateData(r)
	if r.Method != c.method {
		c.log.Error("method is not allowed", "method", c.method, "route", c.route)
		http.Error(w, "method is not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := c.templates.ExecuteTemplate(w, c.tpl.GetMainTpl(), &c.templateData); err != nil {
		c.log.Error("executing template", "err", err, "contentTpl", c.contentTpl, "route", c.route)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *BaseController) initTemplates() {
	if c.contentTpl == "" {
		c.contentTpl = c.tpl.GetDefaultContentTpl()
	}
	tl := c.tpl.GetTemplatesLocation()

	var files []string

	topLevel, err := fs.Glob(os.DirFS(tl), "*.html")
	if err != nil {
		c.log.Error("find tpls - topLevel", "err", err)
	}
	files = append(files, topLevel...)

	subDir, err := fs.Glob(os.DirFS(tl), c.contentTpl)
	if err != nil {
		c.log.Error("find tpls - subDir", "err", err)
		log.Fatal(err)
	}
	files = append(files, subDir...)

	for i, f := range files {
		files[i] = tl + "/" + f
	}

	c.templates, err = template.ParseFiles(files...)
	if err != nil {
		c.log.Error("parse tpls", "err", err)
	}
}
