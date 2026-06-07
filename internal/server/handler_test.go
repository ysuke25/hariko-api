package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ysuke25/hariko/internal/config"
)

func TestNewRouteHandler_Status(t *testing.T) {
	tests := []struct {
		name           string
		status         int
		expectedStatus int
	}{
		{
			name:           "200 OK",
			status:         200,
			expectedStatus: 200,
		},
		{
			name:           "201 Created",
			status:         201,
			expectedStatus: 201,
		},
		{
			name:           "404 Not Found",
			status:         404,
			expectedStatus: 404,
		},
		{
			name:           "500 Internal Server Error",
			status:         500,
			expectedStatus: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			route := config.Route{
				Method: "GET",
				Path:   "/test",
				Response: config.Response{
					Status: tt.status,
				},
			}

			handler := newRouteHandler(route)
			req := httptest.NewRequest("GET", "/test", nil)
			rec := httptest.NewRecorder()

			handler(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestNewRouteHandler_Headers(t *testing.T) {
	tests := []struct {
		name            string
		headers         map[string]string
		expectedHeaders map[string]string
	}{
		{
			name:    "default Content-Type",
			headers: map[string]string{},
			expectedHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			name: "custom Content-Type",
			headers: map[string]string{
				"Content-Type": "text/plain",
			},
			expectedHeaders: map[string]string{
				"Content-Type": "text/plain",
			},
		},
		{
			name: "custom headers with default Content-Type",
			headers: map[string]string{
				"X-Custom-Header": "custom-value",
			},
			expectedHeaders: map[string]string{
				"Content-Type":    "application/json",
				"X-Custom-Header": "custom-value",
			},
		},
		{
			name: "multiple custom headers",
			headers: map[string]string{
				"Content-Type":    "application/xml",
				"X-Custom-Header": "custom-value",
				"X-Another":       "another-value",
			},
			expectedHeaders: map[string]string{
				"Content-Type":    "application/xml",
				"X-Custom-Header": "custom-value",
				"X-Another":       "another-value",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			route := config.Route{
				Method: "GET",
				Path:   "/test",
				Response: config.Response{
					Status:  200,
					Headers: tt.headers,
				},
			}

			handler := newRouteHandler(route)
			req := httptest.NewRequest("GET", "/test", nil)
			rec := httptest.NewRecorder()

			handler(rec, req)

			for key, expectedValue := range tt.expectedHeaders {
				actualValue := rec.Header().Get(key)
				if actualValue != expectedValue {
					t.Errorf("header %s: expected %q, got %q", key, expectedValue, actualValue)
				}
			}
		})
	}
}

func TestNewRouteHandler_Body(t *testing.T) {
	tests := []struct {
		name         string
		body         string
		expectedBody string
	}{
		{
			name:         "JSON body",
			body:         `{"message": "success"}`,
			expectedBody: `{"message": "success"}`,
		},
		{
			name:         "plain text body",
			body:         "success",
			expectedBody: "success",
		},
		{
			name:         "empty body",
			body:         "",
			expectedBody: "",
		},
		{
			name:         "complex JSON body",
			body:         `{"user": {"id": 1, "name": "test"}, "items": [1, 2, 3]}`,
			expectedBody: `{"user": {"id": 1, "name": "test"}, "items": [1, 2, 3]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			route := config.Route{
				Method: "GET",
				Path:   "/test",
				Response: config.Response{
					Status: 200,
					Body:   tt.body,
				},
			}

			handler := newRouteHandler(route)
			req := httptest.NewRequest("GET", "/test", nil)
			rec := httptest.NewRecorder()

			handler(rec, req)

			if rec.Body.String() != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, rec.Body.String())
			}
		})
	}
}

func TestNotFoundHandler(t *testing.T) {
	handler := notFoundHandler()
	req := httptest.NewRequest("GET", "/nonexistent", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}

	expectedContentType := "application/json"
	if contentType := rec.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("expected Content-Type %q, got %q", expectedContentType, contentType)
	}

	expectedBody := `{"error": "route not found"}`
	if rec.Body.String() != expectedBody {
		t.Errorf("expected body %q, got %q", expectedBody, rec.Body.String())
	}
}
