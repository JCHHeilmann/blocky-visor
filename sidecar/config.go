package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const defaultBlockyDir = "/opt/blocky"

type Config struct {
	Listen      string   `yaml:"listen"`
	APIKey      string   `yaml:"api_key"`
	CORSOrigins []string `yaml:"cors_origins"`
	DNSResolver string   `yaml:"dns_resolver"`
	Blocky      struct {
		Dir         string `yaml:"dir"`
		ConfigPath  string `yaml:"config_path"`
		LogDir      string `yaml:"log_dir"`
		ServiceName string `yaml:"service_name"`
	} `yaml:"blocky"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	if cfg.Listen == "" {
		cfg.Listen = ":8550"
	}
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("api_key is required")
	}

	// Default blocky dir
	dir := cfg.Blocky.Dir
	if dir == "" {
		dir = defaultBlockyDir
	}

	// Derive paths from dir if not explicitly set
	if cfg.Blocky.ConfigPath == "" {
		cfg.Blocky.ConfigPath = filepath.Join(dir, "config.yml")
	}
	if cfg.Blocky.LogDir == "" {
		cfg.Blocky.LogDir = filepath.Join(dir, "logs")
	}
	if cfg.Blocky.ServiceName == "" {
		cfg.Blocky.ServiceName = "blocky"
	}

	return &cfg, nil
}
