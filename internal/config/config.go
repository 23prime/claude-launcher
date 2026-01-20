package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Config represents the configuration for claude-launcher
type Config struct {
	AllowedDirs []string
}

// Loader is an interface for loading configuration
type Loader interface {
	Load() (*Config, error)
}

// EnvLoader loads configuration from environment variables
type EnvLoader struct{}

// Load implements the Loader interface for EnvLoader
func (e *EnvLoader) Load() (*Config, error) {
	envValue := os.Getenv("CLAUDE_SAFE_DIRS")
	if envValue == "" {
		return nil, fmt.Errorf("CLAUDE_SAFE_DIRS environment variable not set")
	}

	dirs := strings.Split(envValue, ":")
	expandedDirs := make([]string, 0, len(dirs))
	for _, dir := range dirs {
		if dir == "" {
			continue
		}
		expanded, err := ExpandPath(dir)
		if err != nil {
			return nil, fmt.Errorf("failed to expand path %s: %w", dir, err)
		}
		expandedDirs = append(expandedDirs, expanded)
	}

	if len(expandedDirs) == 0 {
		return nil, fmt.Errorf("no valid directories in CLAUDE_SAFE_DIRS")
	}

	return &Config{AllowedDirs: expandedDirs}, nil
}

// DefaultConfigPath returns the default configuration file path
func DefaultConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, ".config", "claude-launcher", "config.json"), nil
}

// FileLoader loads configuration from ~/.config/claude-launcher/config.json
type FileLoader struct {
	Path string
}

// configJSON represents the structure of the config file
type configJSON struct {
	AllowedDirs []string `json:"allowedDirs"`
}

// Load implements the Loader interface for FileLoader
func (f *FileLoader) Load() (*Config, error) {
	path := f.Path
	if path == "" {
		var err error
		path, err = DefaultConfigPath()
		if err != nil {
			return nil, err
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg configJSON
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	if len(cfg.AllowedDirs) == 0 {
		return nil, fmt.Errorf("no allowedDirs found in config file")
	}

	expandedDirs := make([]string, 0, len(cfg.AllowedDirs))
	for _, dir := range cfg.AllowedDirs {
		expanded, err := ExpandPath(dir)
		if err != nil {
			return nil, fmt.Errorf("failed to expand path %s: %w", dir, err)
		}
		expandedDirs = append(expandedDirs, expanded)
	}

	return &Config{AllowedDirs: expandedDirs}, nil
}

// ChainLoader tries multiple loaders in order
type ChainLoader struct {
	Loaders []Loader
}

// Load implements the Loader interface for ChainLoader
func (c *ChainLoader) Load() (*Config, error) {
	var errors []error

	for _, loader := range c.Loaders {
		config, err := loader.Load()
		if err == nil {
			return config, nil
		}
		errors = append(errors, err)
	}

	if len(errors) == 0 {
		return nil, fmt.Errorf("no loaders configured")
	}

	return nil, fmt.Errorf("all loaders failed: %v", errors)
}

// LoadConfig loads configuration with priority order:
// 1. CLAUDE_SAFE_DIRS environment variable
// 2. ~/.config/claude-launcher/config.json
func LoadConfig() (*Config, error) {
	loader := &ChainLoader{
		Loaders: []Loader{
			&EnvLoader{},
			&FileLoader{},
		},
	}

	return loader.Load()
}

// ExpandPath expands ~ to home directory
func ExpandPath(path string) (string, error) {
	if !strings.HasPrefix(path, "~") {
		return path, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	if path == "~" {
		return homeDir, nil
	}

	if strings.HasPrefix(path, "~/") {
		return filepath.Join(homeDir, path[2:]), nil
	}

	return path, nil
}
