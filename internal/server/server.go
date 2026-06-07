package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/ysuke25/hariko/internal/config"
)

type Server struct {
	cfg    *config.Config
	mux    *http.ServeMux
	logger *slog.Logger
}

func New(cfg *config.Config, logger *slog.Logger) *Server {
	s := &Server{
		cfg:    cfg,
		mux:    http.NewServeMux(),
		logger: logger,
	}
	s.registerRoutes()
	return s
}

func (s *Server) registerRoutes() {
	for _, route := range s.cfg.Routes {
		pattern := buildPattern(route.Method, route.Path)
		s.mux.HandleFunc(pattern, newRouteHandler(route))
		s.logger.Info("registered route", "method", route.Method, "path", route.Path)
	}
}

func buildPattern(method, path string) string {
	return strings.ToUpper(method) + " " + path
}

func (s *Server) Run() error {
	addr := fmt.Sprintf("%s:%d", s.cfg.Server.Host, s.cfg.Server.Port)

	server := &http.Server{
		Addr:    addr,
		Handler: s.mux,
	}

	s.logger.Info("starting hariko", "addr", addr, "routes", len(s.cfg.Routes))
	return server.ListenAndServe()
}
