package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoad_ValidConfig(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		wantPort int
		wantHost string
		wantRoutes int
	}{
		{
			name: "complete config with all fields",
			yaml: `server:
  port: 3000
  host: "localhost"
routes:
  - method: GET
    path: /api/users
    response:
      status: 200
      headers:
        Content-Type: application/json
      body: '{"users":[]}'
  - method: POST
    path: /api/users
    response:
      status: 201
      body: '{"created":true}'`,
			wantPort:   3000,
			wantHost:   "localhost",
			wantRoutes: 2,
		},
		{
			name: "minimal valid config",
			yaml: `server:
  port: 8080
  host: "0.0.0.0"
routes:
  - method: GET
    path: /health
    response:
      status: 200
      body: OK`,
			wantPort:   8080,
			wantHost:   "0.0.0.0",
			wantRoutes: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "config.yaml")
			if err := os.WriteFile(configPath, []byte(tt.yaml), 0644); err != nil {
				t.Fatalf("failed to write test config: %v", err)
			}

			cfg, err := Load(configPath)
			if err != nil {
				t.Fatalf("Load() error = %v, want nil", err)
			}

			if cfg.Server.Port != tt.wantPort {
				t.Errorf("Server.Port = %d, want %d", cfg.Server.Port, tt.wantPort)
			}
			if cfg.Server.Host != tt.wantHost {
				t.Errorf("Server.Host = %q, want %q", cfg.Server.Host, tt.wantHost)
			}
			if len(cfg.Routes) != tt.wantRoutes {
				t.Errorf("len(Routes) = %d, want %d", len(cfg.Routes), tt.wantRoutes)
			}
		})
	}
}

func TestLoad_Defaults(t *testing.T) {
	tests := []struct {
		name         string
		yaml         string
		wantPort     int
		wantHost     string
		wantStatus   int
		description  string
	}{
		{
			name: "default port 8080",
			yaml: `server:
  host: "localhost"
routes:
  - method: GET
    path: /test
    response:
      body: test`,
			wantPort:    8080,
			wantHost:    "localhost",
			wantStatus:  200,
			description: "port defaults to 8080 when not specified",
		},
		{
			name: "default host 0.0.0.0",
			yaml: `server:
  port: 3000
routes:
  - method: GET
    path: /test
    response:
      body: test`,
			wantPort:    3000,
			wantHost:    "0.0.0.0",
			wantStatus:  200,
			description: "host defaults to 0.0.0.0 when not specified",
		},
		{
			name: "default status 200",
			yaml: `server:
  port: 8080
  host: "0.0.0.0"
routes:
  - method: GET
    path: /test
    response:
      body: test`,
			wantPort:    8080,
			wantHost:    "0.0.0.0",
			wantStatus:  200,
			description: "response status defaults to 200 when not specified",
		},
		{
			name: "all defaults applied",
			yaml: `routes:
  - method: GET
    path: /test
    response:
      body: test`,
			wantPort:    8080,
			wantHost:    "0.0.0.0",
			wantStatus:  200,
			description: "all defaults applied when server config omitted",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "config.yaml")
			if err := os.WriteFile(configPath, []byte(tt.yaml), 0644); err != nil {
				t.Fatalf("failed to write test config: %v", err)
			}

			cfg, err := Load(configPath)
			if err != nil {
				t.Fatalf("Load() error = %v, want nil", err)
			}

			if cfg.Server.Port != tt.wantPort {
				t.Errorf("Server.Port = %d, want %d (%s)", cfg.Server.Port, tt.wantPort, tt.description)
			}
			if cfg.Server.Host != tt.wantHost {
				t.Errorf("Server.Host = %q, want %q (%s)", cfg.Server.Host, tt.wantHost, tt.description)
			}
			if len(cfg.Routes) > 0 && cfg.Routes[0].Response.Status != tt.wantStatus {
				t.Errorf("Routes[0].Response.Status = %d, want %d (%s)",
					cfg.Routes[0].Response.Status, tt.wantStatus, tt.description)
			}
		})
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	nonExistentPath := filepath.Join(tmpDir, "does-not-exist.yaml")

	cfg, err := Load(nonExistentPath)
	if err == nil {
		t.Fatal("Load() error = nil, want error for non-existent file")
	}
	if cfg != nil {
		t.Errorf("Load() returned config = %v, want nil", cfg)
	}
	if !strings.Contains(err.Error(), "reading config file") {
		t.Errorf("error message = %q, want to contain 'reading config file'", err.Error())
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	tests := []struct {
		name        string
		yaml        string
		wantErrMsg  string
	}{
		{
			name: "malformed YAML",
			yaml: `server:
  port: not_a_number
  host: "localhost"
routes:
  - method: GET`,
			wantErrMsg: "parsing config file",
		},
		{
			name: "invalid indentation",
			yaml: `server:
port: 8080
  host: "localhost"
routes:
  - method: GET
    path: /test`,
			wantErrMsg: "parsing config file",
		},
		{
			name: "unclosed quote",
			yaml: `server:
  port: 8080
  host: "localhost
routes:
  - method: GET
    path: /test`,
			wantErrMsg: "parsing config file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "config.yaml")
			if err := os.WriteFile(configPath, []byte(tt.yaml), 0644); err != nil {
				t.Fatalf("failed to write test config: %v", err)
			}

			cfg, err := Load(configPath)
			if err == nil {
				t.Fatal("Load() error = nil, want error for invalid YAML")
			}
			if cfg != nil {
				t.Errorf("Load() returned config = %v, want nil", cfg)
			}
			if !strings.Contains(err.Error(), tt.wantErrMsg) {
				t.Errorf("error message = %q, want to contain %q", err.Error(), tt.wantErrMsg)
			}
		})
	}
}

func TestLoad_Validation(t *testing.T) {
	tests := []struct {
		name        string
		yaml        string
		wantErr     bool
		wantErrMsg  string
	}{
		{
			name: "no routes defined",
			yaml: `server:
  port: 8080
  host: "0.0.0.0"
routes: []`,
			wantErr:    true,
			wantErrMsg: "no routes defined",
		},
		{
			name: "missing method",
			yaml: `server:
  port: 8080
routes:
  - path: /test
    response:
      body: test`,
			wantErr:    true,
			wantErrMsg: "method is required",
		},
		{
			name: "missing path",
			yaml: `server:
  port: 8080
routes:
  - method: GET
    response:
      body: test`,
			wantErr:    true,
			wantErrMsg: "path is required",
		},
		{
			name: "path without leading slash",
			yaml: `server:
  port: 8080
routes:
  - method: GET
    path: api/test
    response:
      body: test`,
			wantErr:    true,
			wantErrMsg: "path must start with /",
		},
		{
			name: "valid minimal config",
			yaml: `routes:
  - method: GET
    path: /
    response:
      body: OK`,
			wantErr: false,
		},
		{
			name: "multiple routes with one invalid",
			yaml: `routes:
  - method: GET
    path: /valid
    response:
      body: OK
  - method: POST
    path: invalid
    response:
      body: error`,
			wantErr:    true,
			wantErrMsg: "path must start with /",
		},
		{
			name: "empty method string",
			yaml: `routes:
  - method: ""
    path: /test
    response:
      body: test`,
			wantErr:    true,
			wantErrMsg: "method is required",
		},
		{
			name: "empty path string",
			yaml: `routes:
  - method: GET
    path: ""
    response:
      body: test`,
			wantErr:    true,
			wantErrMsg: "path is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "config.yaml")
			if err := os.WriteFile(configPath, []byte(tt.yaml), 0644); err != nil {
				t.Fatalf("failed to write test config: %v", err)
			}

			cfg, err := Load(configPath)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("Load() error = nil, want error containing %q", tt.wantErrMsg)
				}
				if !strings.Contains(err.Error(), tt.wantErrMsg) {
					t.Errorf("error message = %q, want to contain %q", err.Error(), tt.wantErrMsg)
				}
				if cfg != nil {
					t.Errorf("Load() returned config = %v, want nil on error", cfg)
				}
			} else {
				if err != nil {
					t.Fatalf("Load() error = %v, want nil", err)
				}
				if cfg == nil {
					t.Error("Load() returned nil config, want non-nil")
				}
			}
		})
	}
}
