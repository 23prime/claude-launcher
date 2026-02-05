package account

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

// Selector is an interface for selecting an account
type Selector interface {
	Select(_ []Account) (*Account, error)
}

// InteractiveSelector provides arrow-key based account selection
type InteractiveSelector struct{}

// NewInteractiveSelector creates a new InteractiveSelector
func NewInteractiveSelector() *InteractiveSelector {
	return &InteractiveSelector{}
}

// Select prompts the user to select an account using arrow keys
// If there's only one account, it's automatically selected
func (s *InteractiveSelector) Select(accounts []Account) (*Account, error) {
	if len(accounts) == 0 {
		return nil, fmt.Errorf("no accounts to select from")
	}

	// Auto-select if only one account
	if len(accounts) == 1 {
		return &accounts[0], nil
	}

	// Create items for the prompt
	items := make([]string, len(accounts))
	for i, acc := range accounts {
		items[i] = fmt.Sprintf("%s (%s)", acc.Name, acc.ConfigDir)
	}

	prompt := promptui.Select{
		Label: "Select Claude account",
		Items: items,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}",
			Active:   "\U0001F449 {{ . | cyan }}",
			Inactive: "  {{ . }}",
			Selected: "\U00002714 {{ . | green }}",
		},
	}

	idx, _, err := prompt.Run()
	if err != nil {
		return nil, fmt.Errorf("account selection failed: %w", err)
	}

	return &accounts[idx], nil
}

// SelectAccount loads account configuration and prompts for selection if needed
// Returns nil if no accounts are configured (uses default)
func SelectAccount() (*Account, error) {
	return SelectAccountInteractively()
}

// FindAccountByName looks up an account by name from config
// Returns (account, found) where found indicates if the name was matched
// Returns (nil, false) if no accounts are configured or name not found
func FindAccountByName(accountName string) (*Account, bool, error) {
	cfg, err := LoadAccountConfig()
	if err != nil {
		return nil, false, fmt.Errorf("failed to load account config: %w", err)
	}

	// No accounts configured - use default
	if cfg == nil || len(cfg.Accounts) == 0 {
		return nil, false, nil
	}

	// If account name is specified, try to find it
	if accountName != "" {
		for i := range cfg.Accounts {
			if cfg.Accounts[i].Name == accountName {
				return &cfg.Accounts[i], true, nil
			}
		}
		// Account name not found
		return nil, false, nil
	}

	// No account name specified, need interactive selection
	return nil, false, nil
}

// SelectAccountInteractively prompts the user to select an account
// Returns nil if no accounts are configured
func SelectAccountInteractively() (*Account, error) {
	cfg, err := LoadAccountConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load account config: %w", err)
	}

	// No accounts configured - use default
	if cfg == nil || len(cfg.Accounts) == 0 {
		return nil, nil
	}

	selector := NewInteractiveSelector()
	return selector.Select(cfg.Accounts)
}
