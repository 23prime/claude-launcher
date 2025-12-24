package security

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DirectoryChecker checks if a directory is allowed
type DirectoryChecker struct {
	AllowedDirs []string
}

// NewDirectoryChecker creates a new DirectoryChecker
func NewDirectoryChecker(allowedDirs []string) *DirectoryChecker {
	return &DirectoryChecker{
		AllowedDirs: allowedDirs,
	}
}

// IsAllowed checks if the current directory is allowed
func (dc *DirectoryChecker) IsAllowed(currentDir string) (bool, error) {
	// Resolve the current directory path
	resolvedCurrent, err := ResolvePath(currentDir)
	if err != nil {
		return false, fmt.Errorf("failed to resolve current directory: %w", err)
	}

	for _, allowedDir := range dc.AllowedDirs {
		// Skip if the allowed directory doesn't exist
		if _, err := os.Stat(allowedDir); os.IsNotExist(err) {
			continue
		}

		// Resolve the allowed directory path
		resolvedAllowed, err := ResolvePath(allowedDir)
		if err != nil {
			// Skip this allowed directory if we can't resolve it
			continue
		}

		// Check if current directory is the allowed directory or a subdirectory
		if isPathEqual(resolvedCurrent, resolvedAllowed) || isSubdirectory(resolvedCurrent, resolvedAllowed) {
			return true, nil
		}
	}

	return false, nil
}

// ResolvePath resolves symlinks and returns the absolute path
func ResolvePath(path string) (string, error) {
	// Get absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Evaluate symlinks
	resolvedPath, err := filepath.EvalSymlinks(absPath)
	if err != nil {
		// If EvalSymlinks fails, return the absolute path
		// This can happen if the path doesn't exist yet
		return absPath, nil
	}

	return resolvedPath, nil
}

// isPathEqual checks if two paths are equal
func isPathEqual(path1, path2 string) bool {
	// Clean both paths to normalize them
	clean1 := filepath.Clean(path1)
	clean2 := filepath.Clean(path2)
	return clean1 == clean2
}

// isSubdirectory checks if child is a subdirectory of parent
func isSubdirectory(child, parent string) bool {
	// Clean both paths to normalize them
	cleanChild := filepath.Clean(child)
	cleanParent := filepath.Clean(parent)

	// Same directory is not a subdirectory
	if cleanChild == cleanParent {
		return false
	}

	// Add trailing separator to parent to avoid false positives
	// e.g., /home/user/projects should not match /home/user/project
	if !strings.HasSuffix(cleanParent, string(filepath.Separator)) {
		cleanParent += string(filepath.Separator)
	}

	return strings.HasPrefix(cleanChild+string(filepath.Separator), cleanParent)
}
