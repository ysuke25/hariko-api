package server

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ysuke25/hariko/internal/config"
)

func TestNew(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "localhost",
			Port: 8080,
		},
		Routes: []config.Route{
			{
				Method: "GET",
				Path:   "/api/users",
				Response: config.Response{
					Status: 200,
					Body:   `{"message": "success"}`,
				},
			},
		},
	}

	s := New(cfg, logger)

	if s == nil {
		t.Fatal("expected New to return a non-nil server")
	}

	if s.cfg != cfg {
		t.Error("expected server config to match provided config")
	}

	if s.mux == nil {
		t.Error("expected mux to be initialized")
	}

	if s.matcher == nil {
		t.Error("expected matcher to be initialized")
	}

	if s.logger != logger {
		t.Error("expected logger to match provided logger")
	}
}

func TestBuildPattern(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		path     string
		expected string
	}{
		{
			name:     "GET with simple path",
			method:   "GET",
			path:     "/api/users",
			expected: "GET /api/users",
		},
		{
			name:     "POST with simple path",
			method:   "POST",
			path:     "/api/users",
			expected: "POST /api/users",
		},
		{
			name:     "lowercase method gets uppercased",
			method:   "get",
			path:     "/api/users",
			expected: "GET /api/users",
		},
		{
			name:     "mixed case method gets uppercased",
			method:   "pOsT",
			path:     "/api/posts",
			expected: "POST /api/posts",
		},
		{
			name:     "DELETE with complex path",
			method:   "DELETE",
			path:     "/api/users/123",
			expected: "DELETE /api/users/123",
		},
		{
			name:     "PUT with root path",
			method:   "PUT",
			path:     "/",
			expected: "PUT /",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildPattern(tt.method, tt.path)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestContainsParam(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "path with single parameter",
			path:     "/api/users/:id",
			expected: true,
		},
		{
			name:     "path without parameters",
			path:     "/api/users",
			expected: false,
		},
		{
			name:     "path with multiple parameters",
			path:     "/api/:org/repos/:id",
			expected: true,
		},
		{
			name:     "path with parameter at start",
			path:     "/:resource/items",
			expected: true,
		},
		{
			name:     "path with parameter at end",
			path:     "/api/users/:id",
			expected: true,
		},
		{
			name:     "root path without parameters",
			path:     "/",
			expected: false,
		},
		{
			name:     "complex path without parameters",
			path:     "/api/v1/users/all",
			expected: false,
		},
		{
			name:     "path with colon not as parameter separator",
			path:     "/api/users",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsParam(tt.path)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestMatcherHandler_Found(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "localhost",
			Port: 8080,
		},
		Routes: []config.Route{
			{
				Method: "GET",
				Path:   "/api/users/:id",
				Response: config.Response{
					Status: 200,
					Body:   `{"id": 123, "name": "test"}`,
				},
			},
		},
	}

	s := New(cfg, logger)
	handler := s.matcherHandler()

	req := httptest.NewRequest("GET", "/api/users/123", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	expectedBody := `{"id": 123, "name": "test"}`
	if rec.Body.String() != expectedBody {
		t.Errorf("expected body %q, got %q", expectedBody, rec.Body.String())
	}
}

func TestMatcherHandler_NotFound(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "localhost",
			Port: 8080,
		},
		Routes: []config.Route{
			{
				Method: "GET",
				Path:   "/api/users/:id",
				Response: config.Response{
					Status: 200,
					Body:   `{"id": 123}`,
				},
			},
		},
	}

	s := New(cfg, logger)
	handler := s.matcherHandler()

	req := httptest.NewRequest("GET", "/api/posts/123", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}

	expectedBody := `{"error": "route not found"}`
	if rec.Body.String() != expectedBody {
		t.Errorf("expected body %q, got %q", expectedBody, rec.Body.String())
	}

	expectedContentType := "application/json"
	if contentType := rec.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("expected Content-Type %q, got %q", expectedContentType, contentType)
	}
}

func TestServer_StaticRoutes(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "localhost",
			Port: 8080,
		},
		Routes: []config.Route{
			{
				Method: "GET",
				Path:   "/api/users",
				Response: config.Response{
					Status: 200,
					Body:   `{"users": []}`,
				},
			},
			{
				Method: "POST",
				Path:   "/api/users",
				Response: config.Response{
					Status: 201,
					Body:   `{"id": 1}`,
				},
			},
			{
				Method: "GET",
				Path:   "/api/posts",
				Response: config.Response{
					Status: 200,
					Body:   `{"posts": []}`,
				},
			},
		},
	}

	s := New(cfg, logger)

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "GET /api/users",
			method:         "GET",
			path:           "/api/users",
			expectedStatus: 200,
			expectedBody:   `{"users": []}`,
		},
		{
			name:           "POST /api/users",
			method:         "POST",
			path:           "/api/users",
			expectedStatus: 201,
			expectedBody:   `{"id": 1}`,
		},
		{
			name:           "GET /api/posts",
			method:         "GET",
			path:           "/api/posts",
			expectedStatus: 200,
			expectedBody:   `{"posts": []}`,
		},
		{
			name:           "GET /api/nonexistent",
			method:         "GET",
			path:           "/api/nonexistent",
			expectedStatus: 404,
			expectedBody:   "404 page not found\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()

			s.mux.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if rec.Body.String() != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, rec.Body.String())
			}
		})
	}
}

func TestServer_ParamRoutes(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "localhost",
			Port: 8080,
		},
		Routes: []config.Route{
			{
				Method: "GET",
				Path:   "/api/users/:id",
				Response: config.Response{
					Status: 200,
					Body:   `{"id": "user_id"}`,
				},
			},
			{
				Method: "GET",
				Path:   "/api/:org/repos/:repo",
				Response: config.Response{
					Status: 200,
					Body:   `{"org": "org_name", "repo": "repo_name"}`,
				},
			},
			{
				Method: "DELETE",
				Path:   "/api/posts/:id",
				Response: config.Response{
					Status: 204,
					Body:   "",
				},
			},
		},
	}

	s := New(cfg, logger)

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "GET /api/users/:id with id=123",
			method:         "GET",
			path:           "/api/users/123",
			expectedStatus: 200,
			expectedBody:   `{"id": "user_id"}`,
		},
		{
			name:           "GET /api/users/:id with id=abc",
			method:         "GET",
			path:           "/api/users/abc",
			expectedStatus: 200,
			expectedBody:   `{"id": "user_id"}`,
		},
		{
			name:           "GET /api/:org/repos/:repo",
			method:         "GET",
			path:           "/api/myorg/repos/myrepo",
			expectedStatus: 200,
			expectedBody:   `{"org": "org_name", "repo": "repo_name"}`,
		},
		{
			name:           "DELETE /api/posts/:id",
			method:         "DELETE",
			path:           "/api/posts/456",
			expectedStatus: 204,
			expectedBody:   "",
		},
		{
			name:           "GET /api/users without id should 404",
			method:         "GET",
			path:           "/api/users",
			expectedStatus: 404,
			expectedBody:   `{"error": "route not found"}`,
		},
		{
			name:           "POST /api/users/:id method mismatch should 404",
			method:         "POST",
			path:           "/api/users/123",
			expectedStatus: 404,
			expectedBody:   `{"error": "route not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()

			s.mux.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if rec.Body.String() != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, rec.Body.String())
			}
		})
	}
}

func TestServer_MixedStaticAndParamRoutes(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "localhost",
			Port: 8080,
		},
		Routes: []config.Route{
			{
				Method: "GET",
				Path:   "/api/users",
				Response: config.Response{
					Status: 200,
					Body:   `{"users": "all"}`,
				},
			},
			{
				Method: "GET",
				Path:   "/api/users/:id",
				Response: config.Response{
					Status: 200,
					Body:   `{"user": "specific"}`,
				},
			},
			{
				Method: "GET",
				Path:   "/api/posts",
				Response: config.Response{
					Status: 200,
					Body:   `{"posts": "all"}`,
				},
			},
		},
	}

	s := New(cfg, logger)

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "static route /api/users",
			method:         "GET",
			path:           "/api/users",
			expectedStatus: 200,
			expectedBody:   `{"users": "all"}`,
		},
		{
			name:           "param route /api/users/:id",
			method:         "GET",
			path:           "/api/users/123",
			expectedStatus: 200,
			expectedBody:   `{"user": "specific"}`,
		},
		{
			name:           "static route /api/posts",
			method:         "GET",
			path:           "/api/posts",
			expectedStatus: 200,
			expectedBody:   `{"posts": "all"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()

			s.mux.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if rec.Body.String() != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, rec.Body.String())
			}
		})
	}
}
