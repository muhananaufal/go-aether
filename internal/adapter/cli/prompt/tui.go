package prompt

import (
	"errors"
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/mattn/go-isatty"
)

// IsInteractive checks if the current execution environment supports TTY interaction.
func IsInteractive() bool {
	// If CI is set, assume non-interactive
	if os.Getenv("CI") != "" {
		return false
	}
	return isatty.IsTerminal(os.Stdout.Fd()) && isatty.IsTerminal(os.Stdin.Fd())
}

// EnsureInteractive guarantees the environment supports TTY, returning an error otherwise.
func EnsureInteractive() error {
	if !IsInteractive() {
		return errors.New("terminal is not interactive and required flags are missing (CI mode detected)")
	}
	return nil
}

// AskString prompts the user for a simple string input.
func AskString(title, description string, required bool) (string, error) {
	if err := EnsureInteractive(); err != nil {
		return "", err
	}

	var result string
	input := huh.NewInput().
		Title(title).
		Description(description).
		Value(&result)

	if required {
		input.Validate(func(v string) error {
			if v == "" {
				return fmt.Errorf("this field is required")
			}
			return nil
		})
	}

	err := huh.NewForm(huh.NewGroup(input)).Run()
	return result, err
}

// AskSelect prompts the user to select from a list of options.
func AskSelect(title string, options []string) (string, error) {
	if err := EnsureInteractive(); err != nil {
		return "", err
	}

	if len(options) == 0 {
		return "", errors.New("no options provided for selection")
	}

	var result string
	huhOptions := make([]huh.Option[string], len(options))
	for i, opt := range options {
		huhOptions[i] = huh.NewOption(opt, opt)
	}

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(title).
				Options(huhOptions...).
				Value(&result),
		),
	).Run()
	
	return result, err
}

// PrintDryRunDiff prints a visually distinct message indicating what would have happened.
func PrintDryRunDiff(message string) {
	// Simple fallback to ANSI for dry-run since we removed lipgloss from go.mod requirement to avoid bloat,
	// but wait! lipgloss is typically installed with huh.
	fmt.Printf("\033[36m[DRY RUN] %s\033[0m\n", message)
}
