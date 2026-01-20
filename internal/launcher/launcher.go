package launcher

import (
	"fmt"
	"os"
	"os/exec"
)

// Launcher handles launching Claude Code
type Launcher struct {
	ClaudePath string
}

// NewLauncher creates a new Launcher
func NewLauncher() *Launcher {
	return &Launcher{
		ClaudePath: "claude",
	}
}

// LaunchOptions contains options for launching Claude
type LaunchOptions struct {
	Continue  bool
	Args      []string
	ConfigDir string // Optional: Sets CLAUDE_CONFIG_DIR environment variable
}

// Launch executes Claude Code with the specified options
func (l *Launcher) Launch(opts LaunchOptions) error {
	args := make([]string, 0)

	if opts.Continue {
		args = append(args, "--continue")
	}

	args = append(args, opts.Args...)

	cmd := exec.Command(l.ClaudePath, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	// Set CLAUDE_CONFIG_DIR if specified
	if opts.ConfigDir != "" {
		cmd.Env = append(cmd.Env, "CLAUDE_CONFIG_DIR="+opts.ConfigDir)
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run claude: %w", err)
	}

	return nil
}
