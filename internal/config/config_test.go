package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExpandPath(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home directory: %v", err)
	}

	tests := []struct {
		name     string
		path     string
		expected string
		wantErr  bool
	}{
		{
			name:     "tilde only",
			path:     "~",
			expected: homeDir,
			wantErr:  false,
		},
		{
			name:     "tilde with path",
			path:     "~/projects",
			expected: filepath.Join(homeDir, "projects"),
			wantErr:  false,
		},
		{
			name:     "absolute path",
			path:     "/home/user/projects",
			expected: "/home/user/projects",
			wantErr:  false,
		},
		{
			name:     "relative path",
			path:     "projects",
			expected: "projects",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ExpandPath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExpandPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.expected {
				t.Errorf("ExpandPath() = %v, expected %v", result, tt.expected)
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
			name:        "single directory",
			envValue:    "/home/user/projects",
			wantErr:     false,
			expectedLen: 1,
		},
		{
			name:        "multiple directories",
			envValue:    "/home/user/projects:/home/user/work",
			wantErr:     false,
			expectedLen: 2,
		},
		{
			name:        "with tilde",
			envValue:    "~/projects:~/work",
			wantErr:     false,
			expectedLen: 2,
		},
		{
			name:     "empty value",
			envValue: "",
			wantErr:  true,
		},
		{
			name:        "with empty entries",
			envValue:    "/home/user/projects::/home/user/work",
			wantErr:     false,
			expectedLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			oldEnv := os.Getenv("CLAUDE_SAFE_DIRS")
			defer os.Setenv("CLAUDE_SAFE_DIRS", oldEnv)

			if tt.envValue != "" {
				os.Setenv("CLAUDE_SAFE_DIRS", tt.envValue)
			} else {
				os.Unsetenv("CLAUDE_SAFE_DIRS")
			}

			// Test
			loader := &EnvLoader{}
			config, err := loader.Load()

			if (err != nil) != tt.wantErr {
				t.Errorf("EnvLoader.Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if config == nil {
					t.Error("EnvLoader.Load() returned nil config")
					return
				}
				if len(config.AllowedDirs) != tt.expectedLen {
					t.Errorf("EnvLoader.Load() returned %d dirs, expected %d", len(config.AllowedDirs), tt.expectedLen)
				}
			}
		})
	}
}

func TestFileLoader(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		jsonContent string
		wantErr     bool
		expectedLen int
	}{
		{
			name: "valid config",
			jsonContent: `{
				"allowedDirs": ["/home/user/projects", "/home/user/work"]
			}`,
			wantErr:     false,
			expectedLen: 2,
		},
		{
			name: "valid config with tilde",
			jsonContent: `{
				"allowedDirs": ["~/projects"]
			}`,
			wantErr:     false,
			expectedLen: 1,
		},
		{
			name: "empty allowedDirs",
			jsonContent: `{
				"allowedDirs": []
			}`,
			wantErr: true,
		},
		{
			name:        "invalid JSON",
			jsonContent: `{invalid json`,
			wantErr:     true,
		},
		{
			name: "missing allowedDirs",
			jsonContent: `{
				"otherConfig": {}
			}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test file
			testFile := filepath.Join(tmpDir, "settings.json")
			if err := os.WriteFile(testFile, []byte(tt.jsonContent), 0o644); err != nil {
				t.Fatalf("failed to create test file: %v", err)
			}

			// Test
			loader := &FileLoader{Path: testFile}
			config, err := loader.Load()

			if (err != nil) != tt.wantErr {
				t.Errorf("FileLoader.Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if config == nil {
					t.Error("FileLoader.Load() returned nil config")
					return
				}
				if len(config.AllowedDirs) != tt.expectedLen {
					t.Errorf("FileLoader.Load() returned %d dirs, expected %d", len(config.AllowedDirs), tt.expectedLen)
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
	testFile := filepath.Join(tmpDir, "config.json")
	jsonContent := `{
		"allowedDirs": ["/from/file"]
	}`
	if err := os.WriteFile(testFile, []byte(jsonContent), 0o644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	tests := []struct {
		name        string
		setupEnv    func()
		expectedDir string
	}{
		{
			name: "env takes priority",
			setupEnv: func() {
				os.Setenv("CLAUDE_SAFE_DIRS", "/from/env")
			},
			expectedDir: "/from/env",
		},
		{
			name: "fallback to file",
			setupEnv: func() {
				os.Unsetenv("CLAUDE_SAFE_DIRS")
			},
			expectedDir: "/from/file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			oldEnv := os.Getenv("CLAUDE_SAFE_DIRS")
			defer os.Setenv("CLAUDE_SAFE_DIRS", oldEnv)

			tt.setupEnv()

			// Test
			loader := &ChainLoader{
				Loaders: []Loader{
					&EnvLoader{},
					&FileLoader{Path: testFile},
				},
			}

			config, err := loader.Load()
			if err != nil {
				t.Errorf("ChainLoader.Load() error = %v", err)
				return
			}

			if config == nil {
				t.Error("ChainLoader.Load() returned nil config")
				return
			}

			if len(config.AllowedDirs) == 0 {
				t.Error("ChainLoader.Load() returned empty AllowedDirs")
				return
			}

			if config.AllowedDirs[0] != tt.expectedDir {
				t.Errorf("ChainLoader.Load() returned %v, expected %v", config.AllowedDirs[0], tt.expectedDir)
			}
		})
	}
}

func TestChainLoaderAllFail(t *testing.T) {
	// Ensure env is not set
	oldEnv := os.Getenv("CLAUDE_SAFE_DIRS")
	defer os.Setenv("CLAUDE_SAFE_DIRS", oldEnv)
	os.Unsetenv("CLAUDE_SAFE_DIRS")

	loader := &ChainLoader{
		Loaders: []Loader{
			&EnvLoader{},
			&FileLoader{Path: "/non/existent/file.json"},
		},
	}

	_, err := loader.Load()
	if err == nil {
		t.Error("ChainLoader.Load() should return error when all loaders fail")
	}
}
