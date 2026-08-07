package domain

// Tier states how much evidence stands behind a command.
//
// Ninety commands with no way to tell them apart is worse for a newcomer than
// twenty that are all trustworthy: they cannot know which ones have been proven
// and which merely exist. This is that distinction, published rather than
// implied.
type Tier int

const (
	// TierVerified means the code this command generates is compiled by CI.
	// Not "a test asserts a file appeared" — the Go toolchain builds the output.
	TierVerified Tier = iota

	// TierTested means generation is covered by an automated test, but the
	// generated code is never compiled. Template drift can still slip through:
	// exactly how the per-layer commands shipped broken for several releases.
	TierTested

	// TierExperimental means no automated coverage at all.
	//
	// It does not mean "known broken". It means nobody has proven otherwise,
	// which for a tool that writes code into your repository is the thing worth
	// saying out loud.
	TierExperimental
)

// Badge is the marker appended to a command's one-line description.
// Verified commands carry none: meeting the baseline is not news.
func (t Tier) Badge() string {
	switch t {
	case TierTested:
		return "[tested]"
	case TierExperimental:
		return "[experimental]"
	default:
		return ""
	}
}

func (t Tier) String() string {
	switch t {
	case TierVerified:
		return "verified"
	case TierTested:
		return "tested"
	default:
		return "experimental"
	}
}

// verifiedCommands generate code that CI compiles end to end.
var verifiedCommands = map[string]struct{}{
	"init":            {},
	"adopt":           {},
	"arch:module":     {},
	"arch:domain":     {},
	"arch:port":       {},
	"arch:service":    {},
	"arch:repository": {},
	"arch:handler":    {},
}

// testedCommands have an automated test covering generation, without a compile
// step over the result.
var testedCommands = map[string]struct{}{
	"doctor":             {},
	"api:validator":      {},
	"api:middleware":     {},
	"api:graphql":        {},
	"api:grpc":           {},
	"api:grpc-gateway":   {},
	"api:openapi":        {},
	"cache:redis":        {},
	"cache:dedup":        {},
	"db:sqlc":            {},
	"db:uow":             {},
	"db:readreplica":     {},
	"db:paginator":       {},
	"infra:deploy":       {},
	"infra:cicd":         {},
	"infra:drain":        {},
	"infra:lint":         {},
	"o11y:metrics":       {},
	"platform:tenant":    {},
	"arch:mock":          {},
	"arch:pipeline":      {},
	"arch:specification": {},
	"security:oauth2":    {},
	"security:auditlog":  {},
	"security:argon2":    {},
	"test:stress":        {},
	"test:chaos":         {},
	"test:fuzz":          {},
	"test:benchmark":     {},
	"test:container":     {},
	"test:mutation":      {},
}

// TierOf classifies a command by name.
//
// Anything absent from both tables is experimental by default, so adding a
// command without adding coverage produces an honest label rather than a
// flattering one.
func TierOf(command string) Tier {
	if _, ok := verifiedCommands[command]; ok {
		return TierVerified
	}
	if _, ok := testedCommands[command]; ok {
		return TierTested
	}
	return TierExperimental
}

// VerifiedCommands and TestedCommands expose the tables for reporting and for
// the fitness test that keeps them honest.
func VerifiedCommands() []string { return sortedSlice(verifiedCommands) }
func TestedCommands() []string   { return sortedSlice(testedCommands) }

// TierLegend is shown under the command list so the badges mean something to a
// reader who has never seen them before.
func TierLegend() string {
	return "Command maturity:\n" +
		"  (unmarked)      generated code is compiled by CI\n" +
		"  [tested]        generation is tested; the generated code is not compiled\n" +
		"  [experimental]  no automated coverage yet — not known broken, just unproven"
}

// TierCounts summarises the distribution, used by doctor and by the release
// notes so the number cannot quietly drift.
func TierCounts(commands []string) map[Tier]int {
	counts := map[Tier]int{}
	for _, c := range commands {
		counts[TierOf(c)]++
	}
	return counts
}

func init() {
	// A command listed in both tables would report the stronger tier and hide a
	// gap. Cheap to check once at startup, impossible to get wrong later.
	for name := range verifiedCommands {
		if _, dup := testedCommands[name]; dup {
			panic("aether: command " + name + " is classified in two tiers at once")
		}
	}
}
