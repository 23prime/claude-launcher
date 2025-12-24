package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/23prime/claude-launcher/internal/config"
	"github.com/23prime/claude-launcher/internal/launcher"
	"github.com/23prime/claude-launcher/internal/security"
	"github.com/23prime/claude-launcher/internal/session"
	"github.com/23prime/claude-launcher/internal/ui"
)

const (
	exitSuccess = 0
	exitError   = 1
)

func main() {
	os.Exit(run())
}

func run() int {
	// Parse command-line flags
	showDirs := flag.Bool("show-dirs", false, "Show configured allowed directories")
	flag.BoolVar(showDirs, "l", false, "Show configured allowed directories (shorthand)")
	showHelp := flag.Bool("help", false, "Show help message")
	flag.BoolVar(showHelp, "h", false, "Show help message (shorthand)")
	flag.Parse()

	printer := ui.NewPrinter(os.Stderr)

	// Show help if requested
	if *showHelp {
		showHelpMessage()
		return exitSuccess
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		printer.ShowConfigError()
		return exitError
	}

	// Show allowed directories if requested
	if *showDirs {
		printer.ShowAllowedDirs(cfg.AllowedDirs)
		return exitSuccess
	}

	// Check if current directory is allowed
	currentDir, err := os.Getwd()
	if err != nil {
		printer.Error("Failed to get current directory: %v\n", err)
		return exitError
	}

	checker := security.NewDirectoryChecker(cfg.AllowedDirs)
	allowed, err := checker.IsAllowed(currentDir)
	if err != nil {
		printer.Error("Failed to check directory: %v\n", err)
		return exitError
	}

	if !allowed {
		printer.ShowAccessDenied(currentDir, cfg.AllowedDirs)
		return exitError
	}

	printer.ShowDirectoryAllowed()

	// Ask user about session continuation
	prompter := session.NewInteractivePrompter(os.Stdin, printer)
	shouldContinue, err := prompter.AskContinue()
	if err != nil {
		printer.Error("Failed to read input: %v\n", err)
		return exitError
	}

	// Show what we're doing
	if shouldContinue {
		printer.ShowContinuingSession()
	} else {
		printer.ShowStartingNewSession()
	}

	// Launch Claude
	l := launcher.NewLauncher()
	launchOpts := launcher.LaunchOptions{
		Continue: shouldContinue,
		Args:     flag.Args(),
	}

	if err := l.Launch(launchOpts); err != nil {
		printer.Error("Failed to launch Claude: %v\n", err)
		return exitError
	}

	return exitSuccess
}

func showHelpMessage() {
	help := `claude-launcher - Comprehensive launcher for Claude Code

USAGE:
    claude-launcher [OPTIONS] [CLAUDE_ARGUMENTS...]

OPTIONS:
    -h, --help        Show this help message
    -l, --show-dirs   Show configured allowed directories

DESCRIPTION:
    Combines directory security and session management for Claude Code.

    1. Checks if current directory is in allowed list
    2. Prompts to continue previous session or start fresh
    3. Launches Claude Code with appropriate flags

CONFIGURATION (priority order):
    1. CLAUDE_SAFE_DIRS (highest priority)
        Colon-separated list of allowed directory paths
        Example: export CLAUDE_SAFE_DIRS="$HOME/projects:$HOME/work"

    2. ~/.claude/settings.json (fallback)
        Read from customConfig.allowedDirs array
        Example: {"customConfig": {"allowedDirs": ["/home/user/projects"]}}

EXAMPLES:
    # Configure via environment variable
    export CLAUDE_SAFE_DIRS="$HOME/develop:$HOME/projects"

    # Or configure via settings.json
    # Edit ~/.claude/settings.json and add:
    # {
    #   "customConfig": {
    #     "allowedDirs": ["/home/user/develop", "/home/user/projects"]
    #   }
    # }

    # Launch Claude Code
    claude-launcher

    # Show allowed directories
    claude-launcher --show-dirs
`
	fmt.Print(help)
}
