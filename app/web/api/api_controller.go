package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type ApiController interface {
	Router(mux *http.ServeMux)
}

type DataFiller interface {
	fillData(r *http.Request)
}

type BaseApiController struct {
	self   DataFiller
	log    *slog.Logger
	method string
	route  string
	data   any
}

func (c *BaseApiController) Router(mux *http.ServeMux) {
	mux.HandleFunc(c.route, c.handler)
}

func (c *BaseApiController) handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != c.method {
		c.log.Error("method is not allowed", "method", c.method, "route", c.route)
		http.Error(w, "method is not allowed", http.StatusMethodNotAllowed)
		return
	}

	c.self.fillData(r)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(c.data)
	if err != nil {
		c.log.Error("executing template", "err", err, "route", c.route)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
