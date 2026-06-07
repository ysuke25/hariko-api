package server

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoggingMiddleware_PassesThrough(t *testing.T) {
	handlerCalled := false
	expectedBody := "test response"

	innerHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(expectedBody))
	})

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	middleware := loggingMiddleware(logger, innerHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	middleware.ServeHTTP(rec, req)

	if !handlerCalled {
		t.Error("inner handler was not called")
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Body.String() != expectedBody {
		t.Errorf("expected body %q, got %q", expectedBody, rec.Body.String())
	}
}

func TestLoggingMiddleware_LogsRequest(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		statusCode     int
		expectedFields []string
	}{
		{
			name:       "GET request",
			method:     "GET",
			path:       "/test",
			statusCode: 200,
			expectedFields: []string{
				"method=GET",
				"path=/test",
				"status=200",
				"duration=",
			},
		},
		{
			name:       "POST request",
			method:     "POST",
			path:       "/api/users",
			statusCode: 201,
			expectedFields: []string{
				"method=POST",
				"path=/api/users",
				"status=201",
				"duration=",
			},
		},
		{
			name:       "404 response",
			method:     "GET",
			path:       "/notfound",
			statusCode: 404,
			expectedFields: []string{
				"method=GET",
				"path=/notfound",
				"status=404",
				"duration=",
			},
		},
		{
			name:       "500 error",
			method:     "POST",
			path:       "/error",
			statusCode: 500,
			expectedFields: []string{
				"method=POST",
				"path=/error",
				"status=500",
				"duration=",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var logBuf bytes.Buffer

			innerHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			})

			logger := slog.New(slog.NewTextHandler(&logBuf, nil))
			middleware := loggingMiddleware(logger, innerHandler)

			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()

			middleware.ServeHTTP(rec, req)

			logOutput := logBuf.String()

			for _, field := range tt.expectedFields {
				if !strings.Contains(logOutput, field) {
					t.Errorf("log output missing field %q\nLog output: %s", field, logOutput)
				}
			}
		})
	}
}

func TestResponseWriter_CapturesStatus(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		expectedStatus int
	}{
		{
			name:           "200 OK",
			statusCode:     200,
			expectedStatus: 200,
		},
		{
			name:           "201 Created",
			statusCode:     201,
			expectedStatus: 201,
		},
		{
			name:           "400 Bad Request",
			statusCode:     400,
			expectedStatus: 400,
		},
		{
			name:           "404 Not Found",
			statusCode:     404,
			expectedStatus: 404,
		},
		{
			name:           "500 Internal Server Error",
			statusCode:     500,
			expectedStatus: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			rw := newResponseWriter(rec)

			rw.WriteHeader(tt.statusCode)

			if rw.statusCode != tt.expectedStatus {
				t.Errorf("expected status code %d, got %d", tt.expectedStatus, rw.statusCode)
			}

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected recorder status code %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestResponseWriter_DefaultStatus(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := newResponseWriter(rec)

	if rw.statusCode != http.StatusOK {
		t.Errorf("expected default status code %d, got %d", http.StatusOK, rw.statusCode)
	}
}

func TestResponseWriter_WriteBody(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := newResponseWriter(rec)

	expectedBody := "test body"
	rw.Write([]byte(expectedBody))

	if rec.Body.String() != expectedBody {
		t.Errorf("expected body %q, got %q", expectedBody, rec.Body.String())
	}
}
