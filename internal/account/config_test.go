package account

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseAccountsString(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedLen int
		wantErr     bool
	}{
		{
			name:        "single account",
			input:       "Personal:/home/user/.claude-personal",
			expectedLen: 1,
			wantErr:     false,
		},
		{
			name:        "multiple accounts",
			input:       "Personal:/home/user/.claude-personal,Work:/home/user/.claude-work",
			expectedLen: 2,
			wantErr:     false,
		},
		{
			name:        "with spaces",
			input:       " Personal : /home/user/.claude-personal , Work : /home/user/.claude-work ",
			expectedLen: 2,
			wantErr:     false,
		},
		{
			name:        "with tilde",
			input:       "Personal:~/.claude-personal",
			expectedLen: 1,
			wantErr:     false,
		},
		{
			name:    "invalid format - no colon",
			input:   "Personal",
			wantErr: true,
		},
		{
			name:    "invalid format - empty name",
			input:   ":/home/user/.claude",
			wantErr: true,
		},
		{
			name:    "invalid format - empty configDir",
			input:   "Personal:",
			wantErr: true,
		},
		{
			name:        "empty entries are skipped",
			input:       "Personal:/home/user/.claude-personal,,Work:/home/user/.claude-work",
			expectedLen: 2,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accounts, err := parseAccountsString(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("parseAccountsString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(accounts) != tt.expectedLen {
				t.Errorf("parseAccountsString() returned %d accounts, expected %d", len(accounts), tt.expectedLen)
			}
		})
	}
}

func TestEnvLoader(t *testing.T) {
	tests := []struct {
		name        string
		envValue    string
		wantErr     bool
		expectedLen int
	}{
		{
			name:        "single account",
			envValue:    "Personal:/home/user/.claude-personal",
			wantErr:     false,
			expectedLen: 1,
		},
		{
			name:        "multiple accounts",
			envValue:    "Personal:/home/user/.claude-personal,Work:/home/user/.claude-work",
			wantErr:     false,
			expectedLen: 2,
		},
		{
			name:        "with tilde",
			envValue:    "Personal:~/.claude-personal,Work:~/.claude-work",
			wantErr:     false,
			expectedLen: 2,
		},
		{
			name:     "empty value",
			envValue: "",
			wantErr:  true,
		},
		{
			name:     "invalid format",
			envValue: "InvalidEntry",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			oldEnv := os.Getenv("CLAUDE_ACCOUNTS")
			defer os.Setenv("CLAUDE_ACCOUNTS", oldEnv)

			if tt.envValue != "" {
				os.Setenv("CLAUDE_ACCOUNTS", tt.envValue)
			} else {
				os.Unsetenv("CLAUDE_ACCOUNTS")
			}

			// Test
			loader := &EnvLoader{}
			cfg, err := loader.Load()

			if (err != nil) != tt.wantErr {
				t.Errorf("EnvLoader.Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if cfg == nil {
					t.Error("EnvLoader.Load() returned nil config")
					return
				}
				if len(cfg.Accounts) != tt.expectedLen {
					t.Errorf("EnvLoader.Load() returned %d accounts, expected %d", len(cfg.Accounts), tt.expectedLen)
				}
			}
		})
	}
}

func TestFileLoader(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		jsonContent string
		wantErr     bool
		expectedLen int
	}{
		{
			name: "valid accounts",
			jsonContent: `{
				"customConfig": {
					"accounts": [
						{"name": "Personal", "configDir": "/home/user/.claude-personal"},
						{"name": "Work", "configDir": "/home/user/.claude-work"}
					]
				}
			}`,
			wantErr:     false,
			expectedLen: 2,
		},
		{
			name: "valid accounts with tilde",
			jsonContent: `{
				"customConfig": {
					"accounts": [
						{"name": "Personal", "configDir": "~/.claude-personal"}
					]
				}
			}`,
			wantErr:     false,
			expectedLen: 1,
		},
		{
			name: "empty accounts",
			jsonContent: `{
				"customConfig": {
					"accounts": []
				}
			}`,
			wantErr: true,
		},
		{
			name:        "invalid JSON",
			jsonContent: `{invalid json`,
			wantErr:     true,
		},
		{
			name: "missing customConfig",
			jsonContent: `{
				"otherConfig": {}
			}`,
			wantErr: true,
		},
		{
			name: "account with empty name",
			jsonContent: `{
				"customConfig": {
					"accounts": [
						{"name": "", "configDir": "/home/user/.claude"}
					]
				}
			}`,
			wantErr: true,
		},
		{
			name: "account with empty configDir",
			jsonContent: `{
				"customConfig": {
					"accounts": [
						{"name": "Personal", "configDir": ""}
					]
				}
			}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tmpDir, "settings.json")
			if err := os.WriteFile(testFile, []byte(tt.jsonContent), 0644); err != nil {
				t.Fatalf("failed to create test file: %v", err)
			}

			loader := &FileLoader{Path: testFile}
			cfg, err := loader.Load()

			if (err != nil) != tt.wantErr {
				t.Errorf("FileLoader.Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if cfg == nil {
					t.Error("FileLoader.Load() returned nil config")
					return
				}
				if len(cfg.Accounts) != tt.expectedLen {
					t.Errorf("FileLoader.Load() returned %d accounts, expected %d", len(cfg.Accounts), tt.expectedLen)
				}
			}
		})
	}
}

func TestFileLoaderNonExistentFile(t *testing.T) {
	loader := &FileLoader{Path: "/non/existent/path/settings.json"}
	_, err := loader.Load()
	if err == nil {
		t.Error("FileLoader.Load() should return error for non-existent file")
	}
}

func TestChainLoader(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "settings.json")
	jsonContent := `{
		"customConfig": {
			"accounts": [
				{"name": "FromFile", "configDir": "/from/file"}
			]
		}
	}`
	if err := os.WriteFile(testFile, []byte(jsonContent), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	tests := []struct {
		name         string
		setupEnv     func()
		expectedName string
		expectNil    bool
	}{
		{
			name: "env takes priority",
			setupEnv: func() {
				os.Setenv("CLAUDE_ACCOUNTS", "FromEnv:/from/env")
			},
			expectedName: "FromEnv",
			expectNil:    false,
		},
		{
			name: "fallback to file",
			setupEnv: func() {
				os.Unsetenv("CLAUDE_ACCOUNTS")
			},
			expectedName: "FromFile",
			expectNil:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldEnv := os.Getenv("CLAUDE_ACCOUNTS")
			defer os.Setenv("CLAUDE_ACCOUNTS", oldEnv)

			tt.setupEnv()

			loader := &ChainLoader{
				Loaders: []Loader{
					&EnvLoader{},
					&FileLoader{Path: testFile},
				},
			}

			cfg, err := loader.Load()
			if err != nil {
				t.Errorf("ChainLoader.Load() error = %v", err)
				return
			}

			if tt.expectNil {
				if cfg != nil {
					t.Error("ChainLoader.Load() should return nil")
				}
				return
			}

			if cfg == nil {
				t.Error("ChainLoader.Load() returned nil config")
				return
			}

			if len(cfg.Accounts) == 0 {
				t.Error("ChainLoader.Load() returned empty Accounts")
				return
			}

			if cfg.Accounts[0].Name != tt.expectedName {
				t.Errorf("ChainLoader.Load() returned account %v, expected %v", cfg.Accounts[0].Name, tt.expectedName)
			}
		})
	}
}

func TestChainLoaderAllFail(t *testing.T) {
	oldEnv := os.Getenv("CLAUDE_ACCOUNTS")
	defer os.Setenv("CLAUDE_ACCOUNTS", oldEnv)
	os.Unsetenv("CLAUDE_ACCOUNTS")

	loader := &ChainLoader{
		Loaders: []Loader{
			&EnvLoader{},
			&FileLoader{Path: "/non/existent/file.json"},
		},
	}

	cfg, err := loader.Load()
	if err != nil {
		t.Errorf("ChainLoader.Load() should not return error when all loaders fail, got %v", err)
	}
	if cfg != nil {
		t.Error("ChainLoader.Load() should return nil when all loaders fail")
	}
}

func TestAccountConfigExpansion(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home directory: %v", err)
	}

	oldEnv := os.Getenv("CLAUDE_ACCOUNTS")
	defer os.Setenv("CLAUDE_ACCOUNTS", oldEnv)

	os.Setenv("CLAUDE_ACCOUNTS", "Personal:~/.claude-personal")

	loader := &EnvLoader{}
	cfg, err := loader.Load()
	if err != nil {
		t.Fatalf("EnvLoader.Load() error = %v", err)
	}

	expectedDir := filepath.Join(homeDir, ".claude-personal")
	if cfg.Accounts[0].ConfigDir != expectedDir {
		t.Errorf("ConfigDir = %v, expected %v", cfg.Accounts[0].ConfigDir, expectedDir)
	}
}
