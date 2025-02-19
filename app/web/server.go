package web

import (
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"github.com/meesooqa/ttag/app/web/controllers"
)

type Server struct {
	log         *slog.Logger
	controllers []controllers.Controller

	tplLocationPattern string
	tplStaticLocation  string

	httpServer *http.Server
	templates  *template.Template
}

func NewServer(log *slog.Logger, controllers []controllers.Controller) *Server {
	return &Server{
		log:                log,
		tplLocationPattern: "app/web/templates/*",
		tplStaticLocation:  "app/web/static",
		controllers:        controllers,
	}
}

func (s *Server) Run(ctx context.Context, port int) {
	s.log.Info("[INFO] starting server", "port", port)

	s.templates = template.Must(template.ParseGlob(s.tplLocationPattern))
	s.httpServer = &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           s.router(),
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second, // HTTPResponseTimeout
		IdleTimeout:       60 * time.Second,
	}

	err := s.httpServer.ListenAndServe()
	if err != nil {
		s.log.Error("failed to start server", "err", err)
	}
}

func (s *Server) router() http.Handler {
	mux := http.NewServeMux()

	// Static
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(s.tplStaticLocation))))

	// Route
	for _, controller := range s.controllers {
		controller.Router(mux, s.templates)
	}

	return mux
}
