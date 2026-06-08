package matcher

import (
	"testing"

	"github.com/ysuke25/hariko/internal/config"
)

func TestStaticMatcher_Match(t *testing.T) {
	tests := []struct {
		name          string
		routes        []config.Route
		method        string
		path          string
		expectMatch   bool
		expectedRoute *config.Route
	}{
		{
			name: "exact match",
			routes: []config.Route{
				{
					Method: "GET",
					Path:   "/api/users",
					Response: config.Response{
						Status: 200,
						Body:   `{"users": []}`,
					},
				},
			},
			method:      "GET",
			path:        "/api/users",
			expectMatch: true,
			expectedRoute: &config.Route{
				Method: "GET",
				Path:   "/api/users",
				Response: config.Response{
					Status: 200,
					Body:   `{"users": []}`,
				},
			},
		},
		{
			name: "case-insensitive method matching",
			routes: []config.Route{
				{
					Method: "GET",
					Path:   "/api/users",
					Response: config.Response{
						Status: 200,
					},
				},
			},
			method:      "get",
			path:        "/api/users",
			expectMatch: true,
			expectedRoute: &config.Route{
				Method: "GET",
				Path:   "/api/users",
				Response: config.Response{
					Status: 200,
				},
			},
		},
		{
			name: "different method - no match",
			routes: []config.Route{
				{
					Method: "GET",
					Path:   "/api/users",
					Response: config.Response{
						Status: 200,
					},
				},
			},
			method:        "POST",
			path:          "/api/users",
			expectMatch:   false,
			expectedRoute: nil,
		},
		{
			name: "different path - no match",
			routes: []config.Route{
				{
					Method: "GET",
					Path:   "/api/users",
					Response: config.Response{
						Status: 200,
					},
				},
			},
			method:        "GET",
			path:          "/api/posts",
			expectMatch:   false,
			expectedRoute: nil,
		},
		{
			name: "multiple routes - match correct one",
			routes: []config.Route{
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
						Body:   `{"created": true}`,
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
			method:      "POST",
			path:        "/api/users",
			expectMatch: true,
			expectedRoute: &config.Route{
				Method: "POST",
				Path:   "/api/users",
				Response: config.Response{
					Status: 201,
					Body:   `{"created": true}`,
				},
			},
		},
		{
			name:          "empty routes - no match",
			routes:        []config.Route{},
			method:        "GET",
			path:          "/api/users",
			expectMatch:   false,
			expectedRoute: nil,
		},
		{
			name: "nil routes - no match",
			routes: nil,
			method: "GET",
			path:   "/api/users",
			expectMatch: false,
			expectedRoute: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher := NewStaticMatcher(tt.routes)
			route, matched := matcher.Match(tt.method, tt.path)

			if matched != tt.expectMatch {
				t.Errorf("expected match=%v, got %v", tt.expectMatch, matched)
			}

			if tt.expectMatch {
				if route == nil {
					t.Fatal("expected route to be non-nil")
				}
				if route.Method != tt.expectedRoute.Method {
					t.Errorf("expected method %q, got %q", tt.expectedRoute.Method, route.Method)
				}
				if route.Path != tt.expectedRoute.Path {
					t.Errorf("expected path %q, got %q", tt.expectedRoute.Path, route.Path)
				}
				if route.Response.Status != tt.expectedRoute.Response.Status {
					t.Errorf("expected status %d, got %d", tt.expectedRoute.Response.Status, route.Response.Status)
				}
				if route.Response.Body != tt.expectedRoute.Response.Body {
					t.Errorf("expected body %q, got %q", tt.expectedRoute.Response.Body, route.Response.Body)
				}
			} else {
				if route != nil {
					t.Errorf("expected route to be nil, got %+v", route)
				}
			}
		})
	}
}
