package domain

import "errors"

// Sentinel errors representing all operational anomalies and edge cases in go-aether.
// These errors are leveraged across the entire engine for structured and predictive failure handling.
var (
	ErrManifestNotFound     = errors.New("aether: manifest file (aether.yaml) not found in directory tree")
	ErrManifestCorrupt      = errors.New("aether: manifest file is corrupted or contains syntax errors")
	ErrManifestUnsupported  = errors.New("aether: unsupported manifest schema version")
	ErrInvalidIdentifier    = errors.New("aether: module name is not a valid Go identifier")
	ErrStdlibCollision      = errors.New("aether: module name conflicts with Go standard library package name")
	ErrModuleAlreadyExists  = errors.New("aether: module already exists in manifest or directory")
	ErrModuleNotFound       = errors.New("aether: module not registered in manifest")
	ErrFileConflict         = errors.New("aether: target generation file already exists on filesystem")
	ErrArchitectureMismatch = errors.New("aether: requested pattern is incompatible with configured architecture")
	ErrGoVersionUnsupported = errors.New("aether: Go compiler version older than minimum supported requirement (1.21)")
	ErrWriteFailed          = errors.New("aether: transactional disk write failed; rolling back changes")
	ErrTemplateMissing      = errors.New("aether: requested scaffold template missing from embedded filesystem")
	ErrEnvBindingMissing    = errors.New("aether: required infrastructure connection variable missing in .env.example")
	ErrLegacyLayoutUnmapped = errors.New("aether: brownfield folder architecture requires explicit anomaly mapping")

	// ErrGeneratedCodeInvalid guards the render pipeline: a template that produces
	// syntactically broken Go must never reach the user's disk. Emitting code that
	// fails to parse turns a template bug into a debugging session for the user.
	ErrGeneratedCodeInvalid = errors.New("aether: rendered template produced invalid Go source")

	// ErrPathEscape is returned when a resolved destination would land outside the
	// project root, which is the CLI's hard security boundary against path traversal
	// supplied through command arguments.
	ErrPathEscape = errors.New("aether: resolved destination escapes the project root")

	// ErrReservedName covers identifiers that are legal Go but illegal or hazardous
	// as filenames on a supported platform, most notably the Win32 reserved device
	// names (CON, PRN, AUX, NUL, COM1-9, LPT1-9) which silently swallow file writes.
	ErrReservedName = errors.New("aether: name collides with a reserved operating system device name")

	// ErrGoModMissing signals that scaffolding was requested in a directory tree that
	// is not a Go module, so any generated import path would be unresolvable.
	ErrGoModMissing = errors.New("aether: no go.mod found; run 'go mod init <module-path>' first")
)
