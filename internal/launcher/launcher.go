package launcher

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
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
	ConfigDir string            // Optional: Sets CLAUDE_CONFIG_DIR environment variable
	OtelEnv   map[string]string // Optional: OpenTelemetry environment variables
}

// Launch executes Claude Code with the specified options
func (l *Launcher) Launch(opts LaunchOptions) error {
	args := make([]string, 0)

	if opts.Continue {
		args = append(args, "--continue")
	}

	args = append(args, opts.Args...)

	// #nosec G204 -- ClaudePath defaults to "claude" and args are user-provided CLI arguments
	cmd := exec.Command(l.ClaudePath, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = buildOtelEnv(os.Environ(), opts.OtelEnv)

	if opts.ConfigDir != "" {
		cmd.Env = append(cmd.Env, "CLAUDE_CONFIG_DIR="+opts.ConfigDir)
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run claude: %w", err)
	}

	return nil
}

// buildOtelEnv merges otelEnv into base, skipping keys already present in base.
// Shell env vars (base) take highest priority.
func buildOtelEnv(base []string, otelEnv map[string]string) []string {
	if len(otelEnv) == 0 {
		return base
	}

	existing := make(map[string]bool, len(base))
	for _, e := range base {
		key := strings.SplitN(e, "=", 2)[0]
		existing[key] = true
	}

	result := make([]string, len(base), len(base)+len(otelEnv))
	copy(result, base)

	for k, v := range otelEnv {
		if !existing[k] {
			result = append(result, k+"="+v)
		}
	}

	return result
}
