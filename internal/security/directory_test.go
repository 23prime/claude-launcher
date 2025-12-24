package security

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolvePath(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test directory
	testDir := filepath.Join(tmpDir, "test")
	if err := os.Mkdir(testDir, 0755); err != nil {
		t.Fatalf("failed to create test directory: %v", err)
	}

	// Create a symlink
	symlinkPath := filepath.Join(tmpDir, "symlink")
	if err := os.Symlink(testDir, symlinkPath); err != nil {
		t.Fatalf("failed to create symlink: %v", err)
	}

	tests := []struct {
		name     string
		path     string
		wantSame bool // whether resolved path should be same as input (after cleaning)
	}{
		{
			name:     "regular directory",
			path:     testDir,
			wantSame: true,
		},
		{
			name:     "symlink",
			path:     symlinkPath,
			wantSame: false, // should resolve to target
		},
		{
			name:     "relative path",
			path:     ".",
			wantSame: false, // should become absolute
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolved, err := ResolvePath(tt.path)
			if err != nil {
				t.Errorf("ResolvePath() error = %v", err)
				return
			}

			if resolved == "" {
				t.Error("ResolvePath() returned empty string")
				return
			}

			// Check if path is absolute
			if !filepath.IsAbs(resolved) {
				t.Errorf("ResolvePath() returned relative path: %v", resolved)
			}

			// For symlink test, verify it resolves to the target
			if tt.path == symlinkPath {
				if resolved != testDir {
					t.Errorf("ResolvePath() symlink = %v, expected %v", resolved, testDir)
				}
			}
		})
	}
}

func TestIsPathEqual(t *testing.T) {
	tests := []struct {
		name     string
		path1    string
		path2    string
		expected bool
	}{
		{
			name:     "identical paths",
			path1:    "/home/user/projects",
			path2:    "/home/user/projects",
			expected: true,
		},
		{
			name:     "with trailing slash",
			path1:    "/home/user/projects/",
			path2:    "/home/user/projects",
			expected: true,
		},
		{
			name:     "with dots",
			path1:    "/home/user/./projects",
			path2:    "/home/user/projects",
			expected: true,
		},
		{
			name:     "different paths",
			path1:    "/home/user/projects",
			path2:    "/home/user/work",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isPathEqual(tt.path1, tt.path2)
			if result != tt.expected {
				t.Errorf("isPathEqual() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestIsSubdirectory(t *testing.T) {
	tests := []struct {
		name     string
		child    string
		parent   string
		expected bool
	}{
		{
			name:     "direct subdirectory",
			child:    "/home/user/projects/myproject",
			parent:   "/home/user/projects",
			expected: true,
		},
		{
			name:     "nested subdirectory",
			child:    "/home/user/projects/myproject/src",
			parent:   "/home/user/projects",
			expected: true,
		},
		{
			name:     "not a subdirectory",
			child:    "/home/user/work",
			parent:   "/home/user/projects",
			expected: false,
		},
		{
			name:     "same directory",
			child:    "/home/user/projects",
			parent:   "/home/user/projects",
			expected: false,
		},
		{
			name:     "parent is longer",
			child:    "/home/user/proj",
			parent:   "/home/user/projects",
			expected: false,
		},
		{
			name:     "similar prefix but not subdirectory",
			child:    "/home/user/project",
			parent:   "/home/user/proj",
			expected: false,
		},
		{
			name:     "with trailing slash in parent",
			child:    "/home/user/projects/myproject",
			parent:   "/home/user/projects/",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isSubdirectory(tt.child, tt.parent)
			if result != tt.expected {
				t.Errorf("isSubdirectory(%q, %q) = %v, expected %v", tt.child, tt.parent, result, tt.expected)
			}
		})
	}
}

func TestDirectoryChecker_IsAllowed(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test directory structure
	projectsDir := filepath.Join(tmpDir, "projects")
	myProject := filepath.Join(projectsDir, "myproject")
	workDir := filepath.Join(tmpDir, "work")

	for _, dir := range []string{projectsDir, myProject, workDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("failed to create test directory %s: %v", dir, err)
		}
	}

	tests := []struct {
		name        string
		allowedDirs []string
		currentDir  string
		expected    bool
		wantErr     bool
	}{
		{
			name:        "allowed directory",
			allowedDirs: []string{projectsDir},
			currentDir:  projectsDir,
			expected:    true,
			wantErr:     false,
		},
		{
			name:        "subdirectory of allowed",
			allowedDirs: []string{projectsDir},
			currentDir:  myProject,
			expected:    true,
			wantErr:     false,
		},
		{
			name:        "not allowed",
			allowedDirs: []string{projectsDir},
			currentDir:  workDir,
			expected:    false,
			wantErr:     false,
		},
		{
			name:        "multiple allowed dirs - first match",
			allowedDirs: []string{projectsDir, workDir},
			currentDir:  projectsDir,
			expected:    true,
			wantErr:     false,
		},
		{
			name:        "multiple allowed dirs - second match",
			allowedDirs: []string{projectsDir, workDir},
			currentDir:  workDir,
			expected:    true,
			wantErr:     false,
		},
		{
			name:        "parent of allowed dir not allowed",
			allowedDirs: []string{myProject},
			currentDir:  projectsDir,
			expected:    false,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checker := NewDirectoryChecker(tt.allowedDirs)
			result, err := checker.IsAllowed(tt.currentDir)

			if (err != nil) != tt.wantErr {
				t.Errorf("DirectoryChecker.IsAllowed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if result != tt.expected {
				t.Errorf("DirectoryChecker.IsAllowed() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestDirectoryChecker_IsAllowed_WithSymlink(t *testing.T) {
	tmpDir := t.TempDir()

	// Create real directory
	realDir := filepath.Join(tmpDir, "real")
	if err := os.Mkdir(realDir, 0755); err != nil {
		t.Fatalf("failed to create real directory: %v", err)
	}

	// Create symlink to real directory
	symlinkDir := filepath.Join(tmpDir, "symlink")
	if err := os.Symlink(realDir, symlinkDir); err != nil {
		t.Fatalf("failed to create symlink: %v", err)
	}

	tests := []struct {
		name        string
		allowedDirs []string
		currentDir  string
		expected    bool
	}{
		{
			name:        "real dir allowed, access via real path",
			allowedDirs: []string{realDir},
			currentDir:  realDir,
			expected:    true,
		},
		{
			name:        "real dir allowed, access via symlink",
			allowedDirs: []string{realDir},
			currentDir:  symlinkDir,
			expected:    true,
		},
		{
			name:        "symlink allowed, access via real path",
			allowedDirs: []string{symlinkDir},
			currentDir:  realDir,
			expected:    true,
		},
		{
			name:        "symlink allowed, access via symlink",
			allowedDirs: []string{symlinkDir},
			currentDir:  symlinkDir,
			expected:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checker := NewDirectoryChecker(tt.allowedDirs)
			result, err := checker.IsAllowed(tt.currentDir)

			if err != nil {
				t.Errorf("DirectoryChecker.IsAllowed() error = %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("DirectoryChecker.IsAllowed() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestDirectoryChecker_IsAllowed_NonExistentAllowedDir(t *testing.T) {
	tmpDir := t.TempDir()
	existingDir := filepath.Join(tmpDir, "existing")
	if err := os.Mkdir(existingDir, 0755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}

	checker := NewDirectoryChecker([]string{
		"/non/existent/path",
		existingDir,
	})

	result, err := checker.IsAllowed(existingDir)
	if err != nil {
		t.Errorf("DirectoryChecker.IsAllowed() error = %v", err)
		return
	}

	if !result {
		t.Error("DirectoryChecker.IsAllowed() should return true for existing allowed dir")
	}
}
