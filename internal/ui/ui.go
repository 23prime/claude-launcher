package ui

import (
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
)

// Printer handles formatted output with colors
type Printer struct {
	Writer io.Writer
}

// NewPrinter creates a new Printer
func NewPrinter(writer io.Writer) *Printer {
	if writer == nil {
		writer = os.Stderr
	}
	return &Printer{Writer: writer}
}

// Success prints a success message in green
func (p *Printer) Success(format string, args ...interface{}) {
	green := color.New(color.FgGreen)
	_, _ = green.Fprintf(p.Writer, format, args...) //nolint:errcheck // UI output errors are not critical
}

// Error prints an error message in red
func (p *Printer) Error(format string, args ...interface{}) {
	red := color.New(color.FgRed)
	_, _ = red.Fprintf(p.Writer, format, args...) //nolint:errcheck // UI output errors are not critical
}

// Warning prints a warning message in yellow
func (p *Printer) Warning(format string, args ...interface{}) {
	yellow := color.New(color.FgYellow, color.Bold)
	_, _ = yellow.Fprintf(p.Writer, format, args...) //nolint:errcheck // UI output errors are not critical
}

// Print prints a normal message
func (p *Printer) Print(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(p.Writer, format, args...) //nolint:errcheck // UI output errors are not critical
}

// ShowAllowedDirs displays the list of allowed directories
func (p *Printer) ShowAllowedDirs(dirs []string) {
	p.Print("Allowed directories:\n")
	for _, dir := range dirs {
		p.Print("  - %s\n", dir)
	}
}

// ShowAccessDenied shows an access denied message with details
func (p *Printer) ShowAccessDenied(currentDir string, allowedDirs []string) {
	p.Error("✗ Access denied\n")
	p.Print("\n")
	p.Print("Current directory: %s\n", currentDir)
	p.Print("\n")
	p.Print("Claude Code is not allowed to run in this directory.\n")
	p.Print("Allowed directories:\n")
	for _, dir := range allowedDirs {
		p.Print("  - %s\n", dir)
	}
	p.Print("\n")
}

// ShowConfigError shows a configuration error message
func (p *Printer) ShowConfigError() {
	p.Error("Error: No allowed directories configured\n")
	p.Print("\n")
	p.Print("Please set allowed directories using one of these methods:\n")
	p.Print("\n")
	p.Print("1. Environment variable (colon-separated):\n")
	p.Print("   export CLAUDE_SAFE_DIRS=\"$HOME/projects:$HOME/work\"\n")
	p.Print("\n")
	p.Print("2. Create ~/.config/claude-launcher/config.json:\n")
	p.Print("   {\"allowedDirs\": [\"/home/user/projects\"]}\n")
	p.Print("\n")
}

// ShowDirectoryAllowed shows that the directory check passed
func (p *Printer) ShowDirectoryAllowed() {
	p.Success("✓")
	p.Print(" Directory allowed\n")
	p.Print("\n")
}

// ShowContinuingSession shows that we're continuing the previous session
func (p *Printer) ShowContinuingSession() {
	p.Success("→")
	p.Print(" Continuing previous session...\n")
}

// ShowStartingNewSession shows that we're starting a new session
func (p *Printer) ShowStartingNewSession() {
	p.Success("→")
	p.Print(" Starting new session...\n")
}

// ShowAccountSelected shows that an account was selected
func (p *Printer) ShowAccountSelected(name string, configDir string) {
	p.Success("✓")
	p.Print(" Account: %s (%s)\n", name, configDir)
	p.Print("\n")
}

// ShowNoAccountsConfigured shows that no accounts are configured (using default)
func (p *Printer) ShowNoAccountsConfigured() {
	p.Print("Using default Claude configuration\n")
	p.Print("\n")
}
