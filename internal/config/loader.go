package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	cfg.setDefaults()

	if err := validate(&cfg); err != nil {
		return nil, fmt.Errorf("validating config: %w", err)
	}

	return &cfg, nil
}

func validate(cfg *Config) error {
	if len(cfg.Routes) == 0 {
		return fmt.Errorf("no routes defined")
	}

	for i, r := range cfg.Routes {
		if r.Method == "" {
			return fmt.Errorf("route[%d]: method is required", i)
		}
		if r.Path == "" {
			return fmt.Errorf("route[%d]: path is required", i)
		}
		if !strings.HasPrefix(r.Path, "/") {
			return fmt.Errorf("route[%d]: path must start with /", i)
		}
	}

	return nil
}
