package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/23prime/claude-launcher/internal/account"
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

var (
	// Version information - set by ldflags during build
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
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

	showVersion := flag.Bool("version", false, "Show version information")
	flag.BoolVar(showVersion, "v", false, "Show version information (shorthand)")

	showConfig := flag.Bool("show-config", false, "Show configuration file path and contents")
	flag.BoolVar(showConfig, "c", false, "Show configuration file path and contents (shorthand)")

	accountName := flag.String("account", "", "Account name to use (must exist in config)")
	flag.StringVar(accountName, "a", "", "Account name to use (shorthand)")

	flag.Parse()

	printer := ui.NewPrinter(os.Stderr)

	// Show help if requested
	if *showHelp {
		showHelpMessage()
		return exitSuccess
	}

	if *showVersion {
		showVersionInformation()
		return exitSuccess
	}

	if *showConfig {
		showConfigFile()
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

	// Select account (if configured)
	var selectedAccount *account.Account
	if *accountName != "" {
		// Try to find the specified account
		found, foundOk, err := account.FindAccountByName(*accountName)
		if err != nil {
			printer.Error("Failed to find account: %v\n", err)
			return exitError
		}

		if foundOk {
			selectedAccount = found
		} else {
			// Account not found - show warning before interactive selection
			printer.ShowAccountNotFound(*accountName)
			selectedAccount, err = account.SelectAccountInteractively()
			if err != nil {
				printer.Error("Failed to select account: %v\n", err)
				return exitError
			}
		}
	} else {
		// No account name specified - use interactive selection
		var err error
		selectedAccount, err = account.SelectAccountInteractively()
		if err != nil {
			printer.Error("Failed to select account: %v\n", err)
			return exitError
		}
	}

	var configDir string
	if selectedAccount != nil {
		printer.ShowAccountSelected(selectedAccount.Name, selectedAccount.ConfigDir)
		configDir = selectedAccount.ConfigDir
	}

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
		Continue:  shouldContinue,
		Args:      flag.Args(),
		ConfigDir: configDir,
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
    -h, --help         Show this help message
    -l, --show-dirs    Show configured allowed directories
    -c, --show-config  Show configuration file path and contents
    -v, --version      Show version information
    -a, --account      Account name to use (skips interactive selection)

DESCRIPTION:
    Combines directory security, account selection, and session management
    for Claude Code.

    1. Checks if current directory is in allowed list
    2. Prompts to select account (if multiple accounts configured)
    3. Prompts to continue previous session or start fresh
    4. Launches Claude Code with appropriate flags

CONFIGURATION (priority order):
    Allowed Directories:
    1. CLAUDE_SAFE_DIRS (highest priority)
        Colon-separated list of allowed directory paths
        Example: export CLAUDE_SAFE_DIRS="$HOME/projects:$HOME/work"

    2. ~/.config/claude-launcher/config.json (fallback)
        Read from allowedDirs array
        Example: {"allowedDirs": ["/home/user/projects"]}

    Multiple Accounts (optional):
    1. CLAUDE_ACCOUNTS environment variable (highest priority)
        Comma-separated list of Name:ConfigDir pairs
        Example: export CLAUDE_ACCOUNTS="Personal:~/.claude-personal,Work:~/.claude-work"

    2. ~/.config/claude-launcher/config.json (fallback)
        Read from accounts array
        Example: {"accounts": [
            {"name": "Personal", "configDir": "~/.claude-personal"},
            {"name": "Work", "configDir": "~/.claude-work"}
        ]}

EXAMPLES:
    # Configure allowed directories via environment variable
    export CLAUDE_SAFE_DIRS="$HOME/develop:$HOME/projects"

    # Configure multiple accounts via environment variable
    export CLAUDE_ACCOUNTS="Personal:~/.claude-personal,Work:~/.claude-work"

    # Or configure via config file
    # Create ~/.config/claude-launcher/config.json:
    # {
    #   "allowedDirs": ["/home/user/develop", "/home/user/projects"],
    #   "accounts": [
    #     {"name": "Personal", "configDir": "~/.claude-personal"},
    #     {"name": "Work", "configDir": "~/.claude-work"}
    #   ]
    # }

    # Launch Claude Code (interactive account selection)
    claude-launcher

    # Launch with specific account (skips interactive selection)
    claude-launcher --account Personal

    # Show allowed directories
    claude-launcher --show-dirs
`
	fmt.Print(help)
}

func showVersionInformation() {
	fmt.Printf("claude-launcher %s\n", Version)
	fmt.Printf("  commit: %s\n", GitCommit)
	fmt.Printf("  built:  %s\n", BuildDate)
}

func showConfigFile() {
	configPath, err := config.DefaultConfigPath()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	configPath = filepath.Clean(configPath)
	fmt.Printf("Config file: %s\n\n", configPath)

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("(file does not exist)")
			fmt.Println("\nCreate it with:")
			fmt.Printf("  mkdir -p %s\n", filepath.Dir(configPath))
			fmt.Printf("  echo '{\"allowedDirs\": []}' > %s\n", configPath)
		} else {
			fmt.Printf("Error reading file: %v\n", err)
		}
		return
	}

	fmt.Println("Contents:")
	fmt.Println(string(data))
}
