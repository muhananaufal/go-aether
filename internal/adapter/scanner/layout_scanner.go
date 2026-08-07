// Package scanner reads an existing repository and reports what it appears to
// contain, so brownfield adoption can propose a mapping instead of imposing one.
package scanner

import (
	"context"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/muhananaufal/go-aether/internal/core/domain"
	"github.com/muhananaufal/go-aether/internal/core/port"
)

// Scan limits. A brownfield repository can be a monorepo with millions of
// files; without a ceiling `adopt --scan` becomes a command that appears to
// hang, and the user kills it before it ever produces a proposal.
//
// Four levels is enough to reach internal/app/order/handler, which is about as
// deep as Go layouts usefully nest.
const (
	DefaultMaxDepth = 4
	DefaultMaxFiles = 5000
)

// skippedDirNames are never descended into. Each would contribute thousands of
// files that describe someone else's code, not this project's architecture.
var skippedDirNames = map[string]struct{}{
	".git": {}, ".idea": {}, ".vscode": {}, ".claude": {},
	"vendor": {}, "node_modules": {}, "testdata": {},
	"dist": {}, "build": {}, "bin": {}, "tmp": {}, "coverage": {},
}

// layerSignals describes how one architectural layer announces itself.
type layerSignals struct {
	kind LayoutKind
	// dirNames are directory names that suggest the layer by convention.
	dirNames map[string]struct{}
	// importPrefixes are import paths that corroborate the guess.
	importPrefixes []string
	// fileSuffixes are filename endings that corroborate the guess.
	fileSuffixes []string
	// negativeImports disqualify the layer when present.
	negativeImports []string
}

// LayoutKind is re-exported locally purely to keep the table below readable.
type LayoutKind = domain.LayoutKind

// Scoring weights. A directory matched only by name never reaches the
// confidence threshold on its own, because brownfield repositories are full of
// directories called "api" or "data" that mean something else entirely.
const (
	weightDirName     = 0.50
	weightImport      = 0.40
	weightFileSuffix  = 0.25
	weightCleanDomain = 0.35
)

var signalTable = []layerSignals{
	{
		kind: domain.LayerHandler,
		dirNames: set("handler", "handlers", "controller", "controllers",
			"api", "http", "rest", "transport", "web", "delivery", "endpoint", "endpoints",
			"router", "routers", "routes"),
		importPrefixes: []string{
			"net/http",
			"github.com/gin-gonic/gin",
			"github.com/labstack/echo",
			"github.com/gofiber/fiber",
			"github.com/go-chi/chi",
			"github.com/gorilla/mux",
			"google.golang.org/grpc",
		},
		fileSuffixes: []string{"_handler.go", "_controller.go", "_api.go"},
	},
	{
		kind: domain.LayerRepository,
		dirNames: set("repository", "repositories", "repo", "repos",
			"store", "stores", "dao", "persistence", "data", "database", "db", "storage"),
		importPrefixes: []string{
			"database/sql",
			"github.com/jackc/pgx",
			"gorm.io/gorm",
			"github.com/jmoiron/sqlx",
			"go.mongodb.org/mongo-driver",
			"github.com/redis/go-redis",
			"entgo.io/ent",
			"modernc.org/sqlite",
			"github.com/go-sql-driver/mysql",
		},
		fileSuffixes: []string{"_repository.go", "_repo.go", "_dao.go", "_store.go"},
	},
	{
		kind: domain.LayerService,
		dirNames: set("service", "services", "usecase", "usecases", "uc",
			"business", "logic", "application", "app", "core"),
		fileSuffixes: []string{"_service.go", "_usecase.go", "_uc.go"},
	},
	{
		kind: domain.LayerDomain,
		dirNames: set("domain", "entity", "entities", "model", "models",
			"aggregate", "aggregates", "valueobject"),
		fileSuffixes: []string{"_entity.go", "_model.go"},
		// Business rules that import a web framework or a database driver are not
		// a domain layer, whatever the directory is called.
		negativeImports: []string{"net/http", "database/sql", "github.com/gin-gonic/gin"},
	},
}

// GoLayoutScanner implements port.LayoutScanner over the real filesystem.
type GoLayoutScanner struct {
	MaxDepth int
	MaxFiles int
}

var _ port.LayoutScanner = (*GoLayoutScanner)(nil)

// NewGoLayoutScanner constructs a scanner with the default bounds.
func NewGoLayoutScanner() *GoLayoutScanner {
	return &GoLayoutScanner{MaxDepth: DefaultMaxDepth, MaxFiles: DefaultMaxFiles}
}

// dirFacts accumulates what was observed about a single directory.
type dirFacts struct {
	goFiles     int
	imports     map[string]struct{}
	fileSuffix  map[string]struct{}
	unparseable int
}

// Scan walks root and classifies each directory that contains Go files.
func (s *GoLayoutScanner) Scan(ctx context.Context, root string) (*domain.LayoutReport, error) {
	maxDepth, maxFiles := s.MaxDepth, s.MaxFiles
	if maxDepth <= 0 {
		maxDepth = DefaultMaxDepth
	}
	if maxFiles <= 0 {
		maxFiles = DefaultMaxFiles
	}

	report := &domain.LayoutReport{}
	facts := map[string]*dirFacts{}
	fset := token.NewFileSet()

	walkErr := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// An unreadable directory is a fact about the repository, not a reason
			// to abandon the scan. Permission-denied on one subtree is common in
			// checkouts shared between users.
			if d != nil && d.IsDir() {
				report.SkippedDirs = append(report.SkippedDirs, relSlash(root, path))
				return filepath.SkipDir
			}
			return nil
		}

		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}

		if d.IsDir() {
			if path == root {
				return nil
			}
			if _, skip := skippedDirNames[strings.ToLower(d.Name())]; skip {
				report.SkippedDirs = append(report.SkippedDirs, relSlash(root, path))
				return filepath.SkipDir
			}
			// Hidden directories rarely hold application code and frequently hold
			// large tool caches.
			if strings.HasPrefix(d.Name(), ".") {
				report.SkippedDirs = append(report.SkippedDirs, relSlash(root, path))
				return filepath.SkipDir
			}
			if depthOf(root, path) > maxDepth {
				report.SkippedDirs = append(report.SkippedDirs, relSlash(root, path))
				return filepath.SkipDir
			}
			return nil
		}

		if !strings.HasSuffix(d.Name(), ".go") {
			return nil
		}
		if report.FilesSeen >= maxFiles {
			report.Truncated = true
			return filepath.SkipAll
		}
		report.FilesSeen++

		dir := relSlash(root, filepath.Dir(path))
		f, ok := facts[dir]
		if !ok {
			f = &dirFacts{imports: map[string]struct{}{}, fileSuffix: map[string]struct{}{}}
			facts[dir] = f
		}
		f.goFiles++
		f.fileSuffix[strings.ToLower(d.Name())] = struct{}{}

		// Imports only. Full parsing of a large repository costs seconds per
		// thousand files and buys nothing: the import graph already separates a
		// web handler from a domain entity.
		src, readErr := os.ReadFile(path)
		if readErr != nil {
			f.unparseable++
			return nil
		}
		parsed, parseErr := parser.ParseFile(fset, path, src, parser.ImportsOnly|parser.SkipObjectResolution)
		if parseErr != nil {
			// Legacy repositories contain files that do not compile. Counting them
			// is more honest than pretending the scan saw everything.
			f.unparseable++
			return nil
		}
		for _, imp := range parsed.Imports {
			f.imports[strings.Trim(imp.Path.Value, `"`)] = struct{}{}
		}
		return nil
	})

	if walkErr != nil && !isBenignWalkStop(walkErr) {
		return report, walkErr
	}

	for dir, f := range facts {
		for _, sig := range signalTable {
			if c, ok := score(dir, f, sig); ok {
				report.Candidates = append(report.Candidates, c)
			}
		}
	}

	return report, nil
}

// score turns the observations about one directory into a candidate for one
// layer, or reports that there is no case to make.
func score(dir string, f *dirFacts, sig layerSignals) (domain.LayoutCandidate, bool) {
	var total float64
	var evidence []string

	base := strings.ToLower(pathBase(dir))
	if _, hit := sig.dirNames[base]; hit {
		total += weightDirName
		evidence = append(evidence, "directory named "+base)
	}

	for _, prefix := range sig.importPrefixes {
		if hasImportPrefix(f.imports, prefix) {
			total += weightImport
			evidence = append(evidence, "imports "+prefix)
			break
		}
	}

	for _, suffix := range sig.fileSuffixes {
		if hasFileSuffix(f.fileSuffix, suffix) {
			total += weightFileSuffix
			evidence = append(evidence, "contains *"+suffix)
			break
		}
	}

	// A domain layer is recognised partly by what it does not depend on. Nothing
	// else in the table can use absence as evidence, because absence is only
	// meaningful when the layer is defined by its purity.
	if sig.kind == domain.LayerDomain && total > 0 {
		clean := true
		for _, bad := range sig.negativeImports {
			if hasImportPrefix(f.imports, bad) {
				clean = false
				break
			}
		}
		if clean {
			total += weightCleanDomain
			evidence = append(evidence, "no web or database imports")
		} else {
			return domain.LayoutCandidate{}, false
		}
	}

	if total <= 0 {
		return domain.LayoutCandidate{}, false
	}
	if total > 1 {
		total = 1
	}

	return domain.LayoutCandidate{
		Dir:      dir,
		Kind:     sig.kind,
		Score:    total,
		Evidence: evidence,
		GoFiles:  f.goFiles,
	}, true
}

func hasImportPrefix(imports map[string]struct{}, prefix string) bool {
	for imp := range imports {
		if imp == prefix || strings.HasPrefix(imp, prefix+"/") {
			return true
		}
	}
	return false
}

func hasFileSuffix(names map[string]struct{}, suffix string) bool {
	for name := range names {
		if strings.HasSuffix(name, suffix) {
			return true
		}
	}
	return false
}

// relSlash renders a path relative to root using forward slashes, because the
// result is written into aether.yaml and a backslash there would make the
// manifest unusable on any other operating system.
func relSlash(root, path string) string {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return filepath.ToSlash(path)
	}
	if rel == "." {
		return ""
	}
	return filepath.ToSlash(rel)
}

func depthOf(root, path string) int {
	rel := relSlash(root, path)
	if rel == "" {
		return 0
	}
	return strings.Count(rel, "/") + 1
}

func pathBase(dir string) string {
	if dir == "" {
		return ""
	}
	if idx := strings.LastIndex(dir, "/"); idx >= 0 {
		return dir[idx+1:]
	}
	return dir
}

// isBenignWalkStop reports whether the walk ended because a limit was reached
// rather than because something went wrong.
func isBenignWalkStop(err error) bool {
	return err == filepath.SkipAll || err == filepath.SkipDir
}

func set(values ...string) map[string]struct{} {
	out := make(map[string]struct{}, len(values))
	for _, v := range values {
		out[v] = struct{}{}
	}
	return out
}
