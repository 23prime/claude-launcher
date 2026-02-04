package account

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/23prime/claude-launcher/internal/config"
)

// Account represents a Claude account configuration
type Account struct {
	Name      string
	ConfigDir string
}

// AccountConfig holds the list of configured accounts
type AccountConfig struct {
	Accounts []Account
}

// Loader is an interface for loading account configuration
type Loader interface {
	Load() (*AccountConfig, error)
}

// EnvLoader loads account configuration from CLAUDE_ACCOUNTS environment variable
// Format: "Name1:ConfigDir1,Name2:ConfigDir2"
type EnvLoader struct{}

// Load implements the Loader interface for EnvLoader
func (e *EnvLoader) Load() (*AccountConfig, error) {
	envValue := os.Getenv("CLAUDE_ACCOUNTS")
	if envValue == "" {
		return nil, fmt.Errorf("CLAUDE_ACCOUNTS environment variable not set")
	}

	accounts, err := parseAccountsString(envValue)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CLAUDE_ACCOUNTS: %w", err)
	}

	if len(accounts) == 0 {
		return nil, fmt.Errorf("no valid accounts in CLAUDE_ACCOUNTS")
	}

	return &AccountConfig{Accounts: accounts}, nil
}

// parseAccountsString parses a comma-separated string of "Name:ConfigDir" pairs
func parseAccountsString(s string) ([]Account, error) {
	entries := strings.Split(s, ",")
	accounts := make([]Account, 0, len(entries))

	for _, entry := range entries {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}

		parts := strings.SplitN(entry, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid account entry %q: expected format Name:ConfigDir", entry)
		}

		name := strings.TrimSpace(parts[0])
		configDir := strings.TrimSpace(parts[1])

		if name == "" || configDir == "" {
			return nil, fmt.Errorf("invalid account entry %q: name and configDir cannot be empty", entry)
		}

		expandedDir, err := config.ExpandPath(configDir)
		if err != nil {
			return nil, fmt.Errorf("failed to expand path %s: %w", configDir, err)
		}

		accounts = append(accounts, Account{
			Name:      name,
			ConfigDir: expandedDir,
		})
	}

	return accounts, nil
}

// FileLoader loads account configuration from ~/.config/claude-launcher/config.json
type FileLoader struct {
	Path string
}

// accountJSON represents the account structure in JSON
type accountJSON struct {
	Name      string `json:"name"`
	ConfigDir string `json:"configDir"`
}

// configJSON represents the structure of the config file for accounts
type configJSON struct {
	Accounts []accountJSON `json:"accounts"`
}

// Load implements the Loader interface for FileLoader
func (f *FileLoader) Load() (*AccountConfig, error) {
	path := filepath.Clean(f.Path)
	if path == "" {
		var err error
		path, err = config.DefaultConfigPath()
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

	if len(cfg.Accounts) == 0 {
		return nil, fmt.Errorf("no accounts found in config file")
	}

	accounts := make([]Account, 0, len(cfg.Accounts))
	for _, acc := range cfg.Accounts {
		if acc.Name == "" || acc.ConfigDir == "" {
			return nil, fmt.Errorf("invalid account: name and configDir cannot be empty")
		}

		expandedDir, err := config.ExpandPath(acc.ConfigDir)
		if err != nil {
			return nil, fmt.Errorf("failed to expand path %s: %w", acc.ConfigDir, err)
		}

		accounts = append(accounts, Account{
			Name:      acc.Name,
			ConfigDir: expandedDir,
		})
	}

	return &AccountConfig{Accounts: accounts}, nil
}

// ChainLoader tries multiple loaders in order
type ChainLoader struct {
	Loaders []Loader
}

// Load implements the Loader interface for ChainLoader
// Returns nil config (without error) if no loaders return valid accounts
func (c *ChainLoader) Load() (*AccountConfig, error) {
	for _, loader := range c.Loaders {
		cfg, err := loader.Load()
		if err == nil {
			return cfg, nil
		}
	}

	// No accounts configured - this is not an error, just no accounts
	return nil, nil
}

// LoadAccountConfig loads account configuration with priority order:
// 1. CLAUDE_ACCOUNTS environment variable
// 2. ~/.config/claude-launcher/config.json
// Returns nil if no accounts are configured (not an error)
func LoadAccountConfig() (*AccountConfig, error) {
	loader := &ChainLoader{
		Loaders: []Loader{
			&EnvLoader{},
			&FileLoader{},
		},
	}

	return loader.Load()
}
