package domain

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// validGoPackageRegex defines strict structural rules for Golang package names.
var validGoPackageRegex = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

// stdlibReservedWords contains standard library package names and core language identifiers
// that must never be used as module or feature names to prevent namespace collisions.
var stdlibReservedWords = map[string]struct{}{
	"archive": {}, "bufio": {}, "builtin": {}, "bytes": {}, "compress": {},
	"container": {}, "context": {}, "crypto": {}, "database": {}, "debug": {},
	"embed": {}, "encoding": {}, "errors": {}, "expvar": {}, "flag": {},
	"fmt": {}, "go": {}, "hash": {}, "html": {}, "image": {},
	"index": {}, "io": {}, "log": {}, "math": {}, "mime": {},
	"net": {}, "os": {}, "path": {}, "reflect": {}, "regexp": {},
	"runtime": {}, "sort": {}, "strconv": {}, "strings": {}, "sync": {},
	"syscall": {}, "testing": {}, "text": {}, "time": {}, "unicode": {},
	"unsafe": {}, "main": {}, "init": {}, "common": {}, "util": {}, "utils": {},
}

// windowsReservedDeviceNames lists the Win32 device names that cannot be used as
// a filename stem, even with an extension: "con.go" resolves to the CON console
// device, not to a file.
//
// The failure mode this prevents is the worst kind. Writing succeeds, the CLI
// reports success, and the content vanishes into the device — leaving a project
// that references a source file which does not exist.
//
// Enforced on every operating system rather than only on Windows, so a module
// generated on Linux never becomes un-checkoutable for a teammate on Windows.
var windowsReservedDeviceNames = map[string]struct{}{
	"con": {}, "prn": {}, "aux": {}, "nul": {},
	"com0": {}, "com1": {}, "com2": {}, "com3": {}, "com4": {},
	"com5": {}, "com6": {}, "com7": {}, "com8": {}, "com9": {},
	"lpt0": {}, "lpt1": {}, "lpt2": {}, "lpt3": {}, "lpt4": {},
	"lpt5": {}, "lpt6": {}, "lpt7": {}, "lpt8": {}, "lpt9": {},
}

// TemplateData represents the unified context passed into every text/template during generation.
type TemplateData struct {
	// Module and package branding
	ModuleName      string `json:"module_name"`       // e.g., "order"
	ModuleNamePkg   string `json:"module_name_pkg"`   // e.g., "order" (sanitized lower)
	ModuleNameTitle string `json:"module_name_title"` // e.g., "Order" (TitleCase)
	PackagePath     string `json:"package_path"`      // e.g., "github.com/company/service"
	GoVersion       string `json:"go_version"`        // e.g., "1.23"

	// Architecture context
	ArchPattern string  `json:"arch_pattern"` // e.g., "hexagonal"
	Paths       PathMap `json:"paths"`        // Resolved from AetherManifest

	// Stack configurations
	Router      string   `json:"router"`       // e.g., "chi"
	DBDriver    string   `json:"db_driver"`    // e.g., "postgres"
	CacheDriver string   `json:"cache_driver"` // e.g., "redis"
	HasCache    bool     `json:"has_cache"`
	HasWorker   bool     `json:"has_worker"`
	Transports  []string `json:"transports"`

	// Security & Observability
	AuthPattern   string `json:"auth_pattern"`   // e.g., "jwt-rs256"
	LoggerPattern string `json:"logger_pattern"` // e.g., "slog-otel"

	// Brownfield Bring-Your-Own (BYO) constructor injection
	ExistingDBVar     string `json:"existing_db_var"` // Non-empty triggers constructor injection instead of global connection setup
	ExistingRedisVar  string `json:"existing_redis_var"`
	ExistingLoggerVar string `json:"existing_logger_var"`

	// Generation provenance metadata
	Timestamp     string `json:"timestamp"`      // ISO8601 generation time
	AetherVersion string `json:"aether_version"` // e.g., "0.1.0"
}

// ValidateGoIdentifier inspects a proposed feature or module name to guarantee it compiles
// clean in Go without breaking naming rules or colliding with the standard library.
func ValidateGoIdentifier(name string) error {
	cleaned := strings.TrimSpace(name)
	if len(cleaned) == 0 {
		return fmt.Errorf("%w: name cannot be empty", ErrInvalidIdentifier)
	}
	if !validGoPackageRegex.MatchString(cleaned) {
		return fmt.Errorf("%w: %q (must start with letter and contain only lowercase letters, digits, or underscores)", ErrInvalidIdentifier, cleaned)
	}
	lowered := strings.ToLower(cleaned)
	if _, isReserved := stdlibReservedWords[lowered]; isReserved {
		return fmt.Errorf("%w: %q is a reserved Go standard library package or forbidden word", ErrStdlibCollision, cleaned)
	}
	if _, isDevice := windowsReservedDeviceNames[lowered]; isDevice {
		return fmt.Errorf("%w: %q maps to a Windows device (writes are silently discarded); try %q instead",
			ErrReservedName, cleaned, cleaned+"_module")
	}
	return nil
}

// ToTitleCase transforms a package string into an exported struct title (e.g., "order_item" -> "OrderItem").
func ToTitleCase(input string) string {
	words := strings.Split(input, "_")
	var result strings.Builder
	for _, word := range words {
		if len(word) > 0 {
			r := []rune(word)
			r[0] = unicode.ToUpper(r[0])
			result.WriteString(string(r))
		}
	}
	return result.String()
}

// NewTemplateData constructs a complete template interpolation context from an existing manifest.
func NewTemplateData(moduleName string, manifest *AetherManifest, transports []string, hasCache, hasWorker bool) (*TemplateData, error) {
	if err := ValidateGoIdentifier(moduleName); err != nil {
		return nil, err
	}
	if len(transports) == 0 {
		transports = manifest.Stack.Transport
	}

	return &TemplateData{
		ModuleName:        moduleName,
		ModuleNamePkg:     strings.ToLower(moduleName),
		ModuleNameTitle:   ToTitleCase(moduleName),
		PackagePath:       manifest.Project.Module,
		GoVersion:         manifest.Project.GoVersion,
		ArchPattern:       manifest.Architecture.Pattern,
		Paths:             manifest.Architecture.Paths,
		Router:            manifest.Stack.Router,
		DBDriver:          manifest.Stack.Database.Driver,
		CacheDriver:       manifest.Stack.Cache,
		HasCache:          hasCache,
		HasWorker:         hasWorker,
		Transports:        transports,
		AuthPattern:       manifest.Stack.Auth,
		LoggerPattern:     manifest.Stack.Logger,
		ExistingDBVar:     manifest.Adapters.ExistingDBVar,
		ExistingRedisVar:  manifest.Adapters.ExistingRedisVar,
		ExistingLoggerVar: manifest.Adapters.ExistingLoggerVar,
		Timestamp:         manifest.Project.CreatedAt,
		AetherVersion:     manifest.Project.AetherVersion,
	}, nil
}
