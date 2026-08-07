package scanner_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/muhananaufal/go-aether/internal/adapter/scanner"
	"github.com/muhananaufal/go-aether/internal/core/domain"
)

// writeTree materialises a synthetic repository. Fixtures are built at runtime
// rather than committed because committed .go fixtures become part of this
// module and would be compiled, vetted and gofmt-checked alongside real code.
func writeTree(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()

	for rel, content := range files {
		full := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatalf("mkdir for %s: %v", rel, err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatalf("write %s: %v", rel, err)
		}
	}
	return root
}

// legacyMVC models the shape go-aether previously could not adopt at all: a
// repository whose directories carry conventional names from another ecosystem.
func legacyMVC() map[string]string {
	return map[string]string{
		"web/controllers/order_controller.go": `package controllers

import "net/http"

func Handle(w http.ResponseWriter, r *http.Request) {}
`,
		"logic/order_service.go": `package logic

func Process() error { return nil }
`,
		"data/order_dao.go": `package data

import "database/sql"

func Fetch(db *sql.DB) error { return nil }
`,
		"models/order.go": `package models

type Order struct {
	ID string
}
`,
	}
}

func TestGoLayoutScanner_RecognisesNonStandardLegacyLayout(t *testing.T) {
	root := writeTree(t, legacyMVC())

	report, err := scanner.NewGoLayoutScanner().Scan(context.Background(), root)
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	expected := map[domain.LayoutKind]string{
		domain.LayerHandler:    "web/controllers",
		domain.LayerService:    "logic",
		domain.LayerRepository: "data",
		domain.LayerDomain:     "models",
	}

	for kind, wantDir := range expected {
		best, ok := report.Best(kind)
		if !ok {
			t.Errorf("no confident %s candidate; the scanner would fall back to a default "+
				"and write files into a directory this project does not use", kind)
			continue
		}
		if best.Dir != wantDir {
			t.Errorf("%s: expected %q, got %q (score %.2f)", kind, wantDir, best.Dir, best.Score)
		}
		if len(best.Evidence) == 0 {
			t.Errorf("%s: candidate carries no evidence, so the user cannot judge the proposal", kind)
		}
	}

	if report.FilesSeen != 4 {
		t.Errorf("expected 4 Go files seen, got %d", report.FilesSeen)
	}
	if report.Confidence() != 1.0 {
		t.Errorf("expected full confidence on a recognisable layout, got %.2f", report.Confidence())
	}
}

// TestGoLayoutScanner_NameAloneIsNotEnough is the guard against confident
// nonsense. Brownfield repositories are full of directories called "api" or
// "data" that hold something else entirely, and acting on the name alone is how
// a generator writes into the wrong place.
func TestGoLayoutScanner_NameAloneIsNotEnough(t *testing.T) {
	root := writeTree(t, map[string]string{
		// Named like a handler package, but it is plainly a config loader.
		"api/settings.go": `package api

type Settings struct {
	Region string
}
`,
	})

	report, err := scanner.NewGoLayoutScanner().Scan(context.Background(), root)
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	if best, ok := report.Best(domain.LayerHandler); ok {
		t.Errorf("directory name alone produced a confident handler match on %q (score %.2f); "+
			"corroborating evidence should have been required", best.Dir, best.Score)
	}
}

// TestGoLayoutScanner_DomainRejectsInfrastructureImports proves the domain rule
// uses absence as evidence: a package full of entities that also opens database
// connections is not a domain layer, whatever it is called.
func TestGoLayoutScanner_DomainRejectsInfrastructureImports(t *testing.T) {
	root := writeTree(t, map[string]string{
		"domain/order.go": `package domain

import "database/sql"

type Order struct{ DB *sql.DB }
`,
	})

	report, err := scanner.NewGoLayoutScanner().Scan(context.Background(), root)
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	if _, ok := report.Best(domain.LayerDomain); ok {
		t.Error("a package importing database/sql was accepted as a pure domain layer")
	}
}

func TestGoLayoutScanner_SkipsVendorAndCaches(t *testing.T) {
	files := legacyMVC()
	files["vendor/github.com/x/y/lib.go"] = "package y\n\nimport \"net/http\"\n\nvar _ = http.StatusOK\n"
	files["node_modules/pkg/index.go"] = "package pkg\n"
	files[".git/hooks/hook.go"] = "package hooks\n"

	root := writeTree(t, files)

	report, err := scanner.NewGoLayoutScanner().Scan(context.Background(), root)
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	if report.FilesSeen != 4 {
		t.Errorf("vendored and cached trees were walked: expected 4 files, got %d", report.FilesSeen)
	}
	for _, c := range report.Candidates {
		if strings.HasPrefix(c.Dir, "vendor/") || strings.HasPrefix(c.Dir, "node_modules/") {
			t.Errorf("candidate produced from an excluded tree: %s", c.Dir)
		}
	}
	if len(report.SkippedDirs) == 0 {
		t.Error("skipped directories were not reported, so their absence is unexplainable")
	}
}

// TestGoLayoutScanner_RespectsFileBudget proves the scan is genuinely bounded
// and, just as importantly, admits when it stopped early.
func TestGoLayoutScanner_RespectsFileBudget(t *testing.T) {
	files := map[string]string{}
	for i := 0; i < 40; i++ {
		files[filepath.ToSlash(filepath.Join("pkg", "gen", "f"+string(rune('a'+i%26))+string(rune('a'+i/26))+".go"))] =
			"package gen\n"
	}
	root := writeTree(t, files)

	s := scanner.NewGoLayoutScanner()
	s.MaxFiles = 10

	report, err := s.Scan(context.Background(), root)
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	if report.FilesSeen > 10 {
		t.Errorf("file budget exceeded: %d files read with MaxFiles=10", report.FilesSeen)
	}
	if !report.Truncated {
		t.Error("a truncated scan reported itself as complete; the user would trust a partial result")
	}
}

func TestGoLayoutScanner_HonoursCancellation(t *testing.T) {
	root := writeTree(t, legacyMVC())

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := scanner.NewGoLayoutScanner().Scan(ctx, root); err == nil {
		t.Error("a cancelled scan completed anyway; a monorepo walk would be uninterruptible")
	}
}

// TestGoLayoutScanner_SortedIsStable pins deterministic output. Filesystem walk
// order and map iteration are both unstable enough to reorder an otherwise
// identical report, which makes reviewing a diff of a proposal impossible.
func TestGoLayoutScanner_SortedIsStable(t *testing.T) {
	root := writeTree(t, legacyMVC())
	s := scanner.NewGoLayoutScanner()

	first, err := s.Scan(context.Background(), root)
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}
	second, err := s.Scan(context.Background(), root)
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	a, b := first.Sorted(), second.Sorted()
	if len(a) != len(b) {
		t.Fatalf("candidate count differs between runs: %d vs %d", len(a), len(b))
	}
	for i := range a {
		if a[i].Dir != b[i].Dir || a[i].Kind != b[i].Kind {
			t.Fatalf("ordering differs at %d: %s/%s vs %s/%s", i, a[i].Kind, a[i].Dir, b[i].Kind, b[i].Dir)
		}
	}
}

func TestGoBYODetector_FindsExistingClients(t *testing.T) {
	root := writeTree(t, map[string]string{
		"cmd/api/main.go": `package main

import (
	"database/sql"

	"github.com/redis/go-redis/v9"
)

func main() {
	pool, err := sql.Open("postgres", "dsn")
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	cache := redis.NewClient(&redis.Options{})
	_ = cache
}
`,
	})

	detection, err := scanner.NewGoBYODetector().Detect(context.Background(), root)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if detection.DBVar != "pool" {
		t.Errorf("expected the existing pool to be named %q, got %q", "pool", detection.DBVar)
	}
	if detection.RedisVar != "cache" {
		t.Errorf("expected the existing redis client to be named %q, got %q", "cache", detection.RedisVar)
	}
	if src := detection.Sources["pool"]; !strings.Contains(src, "cmd/api/main.go") {
		t.Errorf("detection must cite its source file, got %q", src)
	}
	if detection.IsEmpty() {
		t.Error("IsEmpty reported true despite two clients being found")
	}
}

// TestGoBYODetector_IgnoresBlankAndErrorIdentifiers guards the obvious mistake:
// reporting a client named "_" or "err" and then generating a constructor that
// asks the user to pass it.
func TestGoBYODetector_IgnoresBlankAndErrorIdentifiers(t *testing.T) {
	root := writeTree(t, map[string]string{
		"main.go": `package main

import "database/sql"

func main() {
	_, err := sql.Open("postgres", "dsn")
	if err != nil {
		panic(err)
	}
}
`,
	})

	detection, err := scanner.NewGoBYODetector().Detect(context.Background(), root)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if detection.DBVar == "_" || detection.DBVar == "err" {
		t.Errorf("detector reported an unusable identifier: %q", detection.DBVar)
	}
}

func TestGoBYODetector_ReadsEntrypointsOnly(t *testing.T) {
	root := writeTree(t, map[string]string{
		// Not package main: a pool constructed here cannot be handed to generated
		// code, so proposing it would be advice the user cannot follow.
		"internal/store/store.go": `package store

import "database/sql"

func New() (*sql.DB, error) { return sql.Open("postgres", "dsn") }
`,
	})

	detection, err := scanner.NewGoBYODetector().Detect(context.Background(), root)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if !detection.IsEmpty() {
		t.Errorf("a non-entrypoint package was mined for clients: %+v", detection)
	}
}

func TestGoBYODetector_HandlesUnparseableEntrypoint(t *testing.T) {
	root := writeTree(t, map[string]string{
		"main.go": "package main\n\nfunc main() { this is not go",
	})

	detection, err := scanner.NewGoBYODetector().Detect(context.Background(), root)
	if err != nil {
		t.Fatalf("a legacy repository containing broken Go must not fail detection: %v", err)
	}
	if !detection.IsEmpty() {
		t.Error("clients were reported from a file that does not parse")
	}
}
