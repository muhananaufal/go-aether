package cli_test

import (
	"sort"
	"strings"
	"testing"

	"github.com/muhananaufal/go-aether/internal/adapter/cli"
	"github.com/muhananaufal/go-aether/internal/core/domain"
)

// registeredCommands builds the real command tree and returns every subcommand
// name. The service is nil on purpose: the constructors only capture it, and
// nothing here executes a command.
func registeredCommands(t *testing.T) []string {
	t.Helper()

	root := cli.NewRootCommand(nil, "test", "none", "unknown")
	var names []string
	for _, cmd := range root.Commands() {
		// Cobra injects help and completion; neither is ours to classify.
		if cmd.Name() == "help" || cmd.Name() == "completion" {
			continue
		}
		names = append(names, cmd.Name())
	}
	sort.Strings(names)
	return names
}

// TestTierTable_HasNoStaleEntries is the fitness function that keeps the
// maturity claims honest.
//
// A typo in the table, or a command that was renamed or removed, would silently
// downgrade a real command to experimental while an entry pointing at nothing
// kept claiming coverage that does not exist. Both are invisible without this.
func TestTierTable_HasNoStaleEntries(t *testing.T) {
	registered := map[string]struct{}{}
	for _, name := range registeredCommands(t) {
		registered[name] = struct{}{}
	}

	for _, claimed := range append(domain.VerifiedCommands(), domain.TestedCommands()...) {
		if _, ok := registered[claimed]; !ok {
			t.Errorf("tier table claims coverage for %q, which is not a registered command; "+
				"either the name is misspelled or the command was renamed", claimed)
		}
	}
}

// TestTierAnnotation_ReachesEveryCommand proves the badge and the machine
// readable annotation are applied to the whole tree, not just the ones that
// happened to be checked by hand.
func TestTierAnnotation_ReachesEveryCommand(t *testing.T) {
	root := cli.NewRootCommand(nil, "test", "none", "unknown")

	for _, cmd := range root.Commands() {
		if cmd.Name() == "help" || cmd.Name() == "completion" {
			continue
		}

		tier, ok := cmd.Annotations["aether.tier"]
		if !ok {
			t.Errorf("%s carries no aether.tier annotation", cmd.Name())
			continue
		}

		expected := domain.TierOf(cmd.Name())
		if tier != expected.String() {
			t.Errorf("%s is annotated %q but classified %q", cmd.Name(), tier, expected)
		}

		// The badge must be visible in the one-line description, since that is
		// what a user actually reads in `--help`.
		if badge := expected.Badge(); badge != "" && !strings.HasSuffix(cmd.Short, badge) {
			t.Errorf("%s is %s but its description does not end with %s: %q",
				cmd.Name(), expected, badge, cmd.Short)
		}
	}
}

// TestTierDistribution_IsReported records the split at the time of writing.
//
// It is deliberately an exact assertion. The number moving is the entire point
// of publishing it: a change here means either coverage was added, which should
// be celebrated in the release notes, or a command was added without coverage,
// which should be argued for rather than slipped in.
func TestTierDistribution_IsReported(t *testing.T) {
	commands := registeredCommands(t)
	counts := domain.TierCounts(commands)

	const (
		wantVerified     = 8
		wantTested       = 31
		wantExperimental = 51
	)

	if got := counts[domain.TierVerified]; got != wantVerified {
		t.Errorf("verified commands: want %d, got %d", wantVerified, got)
	}
	if got := counts[domain.TierTested]; got != wantTested {
		t.Errorf("tested commands: want %d, got %d", wantTested, got)
	}
	if got := counts[domain.TierExperimental]; got != wantExperimental {
		t.Errorf("experimental commands: want %d, got %d "+
			"(a new command defaults to experimental, which is correct but should be deliberate)",
			wantExperimental, got)
	}

	if total := len(commands); total != wantVerified+wantTested+wantExperimental {
		t.Errorf("command count changed: %d registered, tiers account for %d",
			total, wantVerified+wantTested+wantExperimental)
	}
}
