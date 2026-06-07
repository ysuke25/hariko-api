package server

import (
	"net/http"

	"github.com/ysuke25/hariko/internal/config"
)

func newRouteHandler(route config.Route) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for key, value := range route.Response.Headers {
			w.Header().Set(key, value)
		}

		if _, hasContentType := route.Response.Headers["Content-Type"]; !hasContentType {
			w.Header().Set("Content-Type", "application/json")
		}

		w.WriteHeader(route.Response.Status)

		if route.Response.Body != "" {
			w.Write([]byte(route.Response.Body))
		}
	}
}

func notFoundHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "route not found"}`))
	}
}
