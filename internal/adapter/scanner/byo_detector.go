package scanner

import (
	"context"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/muhananaufal/go-aether/internal/core/domain"
	"github.com/muhananaufal/go-aether/internal/core/port"
)

// constructorSignal maps a package-qualified constructor call to the kind of
// client it produces.
//
// Matching the call rather than the variable's declared type is deliberate:
// most real entrypoints write `db, err := sql.Open(...)`, where the type is
// never spelled out and only inference knows it.
type constructorSignal struct {
	pkg   string
	funcs []string
	kind  clientKind
}

type clientKind int

const (
	clientDB clientKind = iota
	clientRedis
	clientLogger
)

var constructorSignals = []constructorSignal{
	{pkg: "sql", funcs: []string{"Open", "OpenDB"}, kind: clientDB},
	{pkg: "pgxpool", funcs: []string{"New", "NewWithConfig", "Connect"}, kind: clientDB},
	{pkg: "pgx", funcs: []string{"Connect", "ConnectConfig"}, kind: clientDB},
	{pkg: "sqlx", funcs: []string{"Open", "Connect", "MustConnect"}, kind: clientDB},
	{pkg: "gorm", funcs: []string{"Open"}, kind: clientDB},
	{pkg: "redis", funcs: []string{"NewClient", "NewFailoverClient", "NewClusterClient"}, kind: clientRedis},
	{pkg: "zap", funcs: []string{"NewProduction", "NewDevelopment", "New"}, kind: clientLogger},
	{pkg: "zerolog", funcs: []string{"New"}, kind: clientLogger},
	{pkg: "logrus", funcs: []string{"New"}, kind: clientLogger},
	{pkg: "slog", funcs: []string{"New"}, kind: clientLogger},
}

// typeSignals recognise explicit declarations such as `var db *sql.DB`, which
// appear in projects that keep their clients in package-level state.
var typeSignals = map[string]clientKind{
	"sql.DB":         clientDB,
	"pgxpool.Pool":   clientDB,
	"sqlx.DB":        clientDB,
	"gorm.DB":        clientDB,
	"redis.Client":   clientRedis,
	"zap.Logger":     clientLogger,
	"slog.Logger":    clientLogger,
	"zerolog.Logger": clientLogger,
}

// GoBYODetector implements port.BYODetector by reading a project's entrypoints.
type GoBYODetector struct {
	MaxDepth int
}

var _ port.BYODetector = (*GoBYODetector)(nil)

// NewGoBYODetector constructs a detector with the default depth bound.
func NewGoBYODetector() *GoBYODetector {
	return &GoBYODetector{MaxDepth: DefaultMaxDepth}
}

// Detect inspects package main files under root for infrastructure clients the
// project already builds.
//
// Only entrypoints are read. A connection pool constructed deep inside a package
// is not something generated code can be handed, so reporting it would produce a
// suggestion the user cannot act on.
func (d *GoBYODetector) Detect(ctx context.Context, root string) (*domain.BYODetection, error) {
	maxDepth := d.MaxDepth
	if maxDepth <= 0 {
		maxDepth = DefaultMaxDepth
	}

	result := &domain.BYODetection{Sources: map[string]string{}}
	fset := token.NewFileSet()

	walkErr := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			if entry != nil && entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}

		if entry.IsDir() {
			if path == root {
				return nil
			}
			name := strings.ToLower(entry.Name())
			if _, skip := skippedDirNames[name]; skip || strings.HasPrefix(name, ".") {
				return filepath.SkipDir
			}
			if depthOf(root, path) > maxDepth {
				return filepath.SkipDir
			}
			return nil
		}

		if !strings.HasSuffix(entry.Name(), ".go") || strings.HasSuffix(entry.Name(), "_test.go") {
			return nil
		}

		src, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil
		}

		// Cheap gate first: parse the package clause alone and move on unless this
		// is an entrypoint.
		head, headErr := parser.ParseFile(fset, path, src, parser.PackageClauseOnly|parser.SkipObjectResolution)
		if headErr != nil || head.Name == nil || head.Name.Name != "main" {
			return nil
		}

		file, parseErr := parser.ParseFile(fset, path, src, parser.SkipObjectResolution)
		if parseErr != nil {
			// An entrypoint that does not parse cannot be reasoned about, but it is
			// also not a reason to fail the whole detection.
			return nil
		}

		d.inspectFile(file, relSlash(root, path), result)
		return nil
	})

	if walkErr != nil && !isBenignWalkStop(walkErr) {
		return result, walkErr
	}
	return result, nil
}

func (d *GoBYODetector) inspectFile(file *ast.File, relPath string, result *domain.BYODetection) {
	record := func(kind clientKind, name, evidence string) {
		if name == "" || name == "_" {
			return
		}
		switch kind {
		case clientDB:
			if result.DBVar == "" {
				result.DBVar = name
				result.Sources[name] = relPath + " (" + evidence + ")"
			}
		case clientRedis:
			if result.RedisVar == "" {
				result.RedisVar = name
				result.Sources[name] = relPath + " (" + evidence + ")"
			}
		case clientLogger:
			if result.LoggerVar == "" {
				result.LoggerVar = name
				result.Sources[name] = relPath + " (" + evidence + ")"
			}
		}
	}

	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.AssignStmt:
			kind, evidence, ok := matchConstructor(node.Rhs)
			if !ok {
				return true
			}
			if name, found := firstIdent(node.Lhs); found {
				record(kind, name, evidence)
			}

		case *ast.ValueSpec:
			// var db *sql.DB
			if node.Type != nil {
				if kind, evidence, ok := matchType(node.Type); ok && len(node.Names) > 0 {
					record(kind, node.Names[0].Name, evidence)
					return true
				}
			}
			// var db = sql.Open(...)
			if kind, evidence, ok := matchConstructor(node.Values); ok && len(node.Names) > 0 {
				record(kind, node.Names[0].Name, evidence)
			}
		}
		return true
	})
}

// matchConstructor reports whether any expression is a recognised constructor
// call such as sql.Open or redis.NewClient.
func matchConstructor(exprs []ast.Expr) (clientKind, string, bool) {
	for _, expr := range exprs {
		call, ok := expr.(*ast.CallExpr)
		if !ok {
			continue
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			continue
		}
		pkgIdent, ok := sel.X.(*ast.Ident)
		if !ok {
			continue
		}

		for _, sig := range constructorSignals {
			if pkgIdent.Name != sig.pkg {
				continue
			}
			for _, fn := range sig.funcs {
				if sel.Sel.Name == fn {
					return sig.kind, pkgIdent.Name + "." + fn, true
				}
			}
		}
	}
	return 0, "", false
}

// matchType reports whether a declared type is a recognised client, handling
// both *sql.DB and sql.DB.
func matchType(expr ast.Expr) (clientKind, string, bool) {
	if star, ok := expr.(*ast.StarExpr); ok {
		expr = star.X
	}
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return 0, "", false
	}
	pkgIdent, ok := sel.X.(*ast.Ident)
	if !ok {
		return 0, "", false
	}

	qualified := pkgIdent.Name + "." + sel.Sel.Name
	if kind, found := typeSignals[qualified]; found {
		return kind, "declared as *" + qualified, true
	}
	return 0, "", false
}

// firstIdent returns the first usable identifier on the left-hand side, skipping
// the blank identifier so `_, err := sql.Open(...)` is not reported as a client
// named "_".
func firstIdent(exprs []ast.Expr) (string, bool) {
	for _, expr := range exprs {
		ident, ok := expr.(*ast.Ident)
		if !ok || ident.Name == "_" || ident.Name == "err" {
			continue
		}
		return ident.Name, true
	}
	return "", false
}
