package session

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/23prime/claude-launcher/internal/ui"
)

// Prompter is an interface for asking user about session continuation
type Prompter interface {
	AskContinue() (bool, error)
}

// InteractivePrompter prompts the user interactively
type InteractivePrompter struct {
	Reader  io.Reader
	Printer *ui.Printer
}

// NewInteractivePrompter creates a new InteractivePrompter
func NewInteractivePrompter(reader io.Reader, printer *ui.Printer) *InteractivePrompter {
	return &InteractivePrompter{
		Reader:  reader,
		Printer: printer,
	}
}

// AskContinue asks the user if they want to continue the previous session
func (p *InteractivePrompter) AskContinue() (bool, error) {
	p.Printer.Warning("Continue previous Claude session?\n")
	p.Printer.Print("  [Y/n] (default: y): ")

	scanner := bufio.NewScanner(p.Reader)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return false, fmt.Errorf("failed to read input: %w", err)
		}
		// EOF or no input, use default (yes)
		return true, nil
	}

	response := strings.TrimSpace(scanner.Text())
	response = strings.ToLower(response)

	switch response {
	case "n", "no":
		return false, nil
	case "", "y", "yes":
		return true, nil
	default:
		// For any other input, default to yes
		return true, nil
	}
}
