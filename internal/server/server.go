package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/ysuke25/hariko/internal/config"
	"github.com/ysuke25/hariko/internal/matcher"
)

type Server struct {
	cfg     *config.Config
	mux     *http.ServeMux
	matcher matcher.RouteMatcher
	logger  *slog.Logger
}

func New(cfg *config.Config, logger *slog.Logger) *Server {
	s := &Server{
		cfg:     cfg,
		mux:     http.NewServeMux(),
		matcher: matcher.NewStaticMatcher(cfg.Routes),
		logger:  logger,
	}
	s.registerRoutes()
	return s
}

func (s *Server) registerRoutes() {
	hasParamRoutes := false

	for _, route := range s.cfg.Routes {
		if containsParam(route.Path) {
			hasParamRoutes = true
		} else {
			pattern := buildPattern(route.Method, route.Path)
			s.mux.HandleFunc(pattern, newRouteHandler(route))
		}
		s.logger.Info("registered route", "method", route.Method, "path", route.Path)
	}

	if hasParamRoutes {
		s.mux.HandleFunc("/", s.matcherHandler())
	}
}

func (s *Server) matcherHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route, found := s.matcher.Match(r.Method, r.URL.Path)
		if !found {
			notFoundHandler()(w, r)
			return
		}
		newRouteHandler(*route)(w, r)
	}
}

func containsParam(path string) bool {
	return strings.Contains(path, "/:")
}

func buildPattern(method, path string) string {
	return strings.ToUpper(method) + " " + path
}

func (s *Server) Run() error {
	addr := fmt.Sprintf("%s:%d", s.cfg.Server.Host, s.cfg.Server.Port)

	server := &http.Server{
		Addr:    addr,
		Handler: loggingMiddleware(s.logger, s.mux),
	}

	s.logger.Info("starting hariko", "addr", addr, "routes", len(s.cfg.Routes))
	return server.ListenAndServe()
}
