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

// FileLoader loads configuration from ~/.claude/settings.json
type FileLoader struct {
	Path string
}

// settingsJSON represents the structure of ~/.claude/settings.json
type settingsJSON struct {
	CustomConfig struct {
		AllowedDirs []string `json:"allowedDirs"`
	} `json:"customConfig"`
}

// Load implements the Loader interface for FileLoader
func (f *FileLoader) Load() (*Config, error) {
	path := f.Path
	if path == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		path = filepath.Join(homeDir, ".claude", "settings.json")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read settings file: %w", err)
	}

	var settings settingsJSON
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, fmt.Errorf("failed to parse settings JSON: %w", err)
	}

	if len(settings.CustomConfig.AllowedDirs) == 0 {
		return nil, fmt.Errorf("no allowedDirs found in settings.json")
	}

	expandedDirs := make([]string, 0, len(settings.CustomConfig.AllowedDirs))
	for _, dir := range settings.CustomConfig.AllowedDirs {
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
// 2. ~/.claude/settings.json
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
