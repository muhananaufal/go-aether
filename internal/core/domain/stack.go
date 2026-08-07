package domain

import (
	"fmt"
	"sort"
	"strings"
)

// The sets below are the single source of truth for what --arch, --router and
// --db actually deliver. They exist because the CLI used to accept any value and
// then generate a chi + lib/pq project regardless, which is the most expensive
// kind of bug: the flag is accepted, the output is wrong, and nothing fails.
//
// Adding a value here is a promise. It must be backed by a template and covered
// by the matrix gate in tests/e2e, or the promise is empty again.
var (
	// supportedArchitectures is deliberately a single entry. Clean and DDD layouts
	// were offered by the interactive prompt for several releases without any
	// templates behind them; listing only what exists is more useful than listing
	// aspirations.
	supportedArchitectures = map[string]struct{}{
		"hexagonal": {},
	}

	supportedRouters = map[string]struct{}{
		"chi":    {},
		"gin":    {},
		"echo":   {},
		"fiber":  {},
		"stdlib": {},
	}

	// "none" is supported and means the generated project wires no database at
	// all, which is the correct choice for a pure API gateway or a CLI service.
	supportedDBDrivers = map[string]struct{}{
		"postgres": {},
		"mysql":    {},
		"sqlite":   {},
		"none":     {},
	}
)

// Default stack selections, applied when a caller supplies an empty value.
const (
	DefaultArchitecture = "hexagonal"
	DefaultRouter       = "chi"
	DefaultDBDriver     = "postgres"
)

// NormalizeStack fills in defaults and lowercases the selection so that
// "Postgres" and "postgres" are the same choice. It performs no validation;
// call ValidateStackSelection on the result.
func NormalizeStack(arch, dbDriver, router string) (string, string, string) {
	normalize := func(value, fallback string) string {
		trimmed := strings.ToLower(strings.TrimSpace(value))
		if trimmed == "" {
			return fallback
		}
		return trimmed
	}
	return normalize(arch, DefaultArchitecture),
		normalize(dbDriver, DefaultDBDriver),
		normalize(router, DefaultRouter)
}

// ValidateStackSelection rejects a combination the generator cannot honour.
//
// Rejection happens before a single byte is written, so an unsupported choice
// leaves the target directory exactly as it was found rather than half
// populated with a stack the user did not ask for.
func ValidateStackSelection(arch, dbDriver, router string) error {
	if _, ok := supportedArchitectures[arch]; !ok {
		return fmt.Errorf("%w: architecture %q is not implemented; supported: %s",
			ErrUnsupportedStack, arch, sortedKeys(supportedArchitectures))
	}
	if _, ok := supportedRouters[router]; !ok {
		return fmt.Errorf("%w: router %q has no template; supported: %s",
			ErrUnsupportedStack, router, sortedKeys(supportedRouters))
	}
	if _, ok := supportedDBDrivers[dbDriver]; !ok {
		return fmt.Errorf("%w: database driver %q has no template; supported: %s",
			ErrUnsupportedStack, dbDriver, sortedKeys(supportedDBDrivers))
	}
	return nil
}

// SupportedRouters reports the router identifiers the CLI may offer, sorted so
// the interactive prompt and the help text stay stable between runs.
func SupportedRouters() []string { return sortedSlice(supportedRouters) }

// SupportedDBDrivers reports the database identifiers the CLI may offer.
func SupportedDBDrivers() []string { return sortedSlice(supportedDBDrivers) }

// SupportedArchitectures reports the architecture identifiers the CLI may offer.
func SupportedArchitectures() []string { return sortedSlice(supportedArchitectures) }

// HasDatabase reports whether a driver selection implies generating database
// wiring at all.
func HasDatabase(dbDriver string) bool {
	return strings.ToLower(strings.TrimSpace(dbDriver)) != "none"
}

func sortedSlice(set map[string]struct{}) []string {
	out := make([]string, 0, len(set))
	for k := range set {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// sortedKeys renders a set for an error message. Map iteration order is
// randomised in Go, and an error whose wording changes between runs is harder to
// search for and impossible to golden-test.
func sortedKeys(set map[string]struct{}) string {
	return strings.Join(sortedSlice(set), ", ")
}
