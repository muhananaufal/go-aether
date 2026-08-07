package e2e_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeLegacyRepo materialises a repository shaped like the ones adoption
// actually meets: conventional directory names borrowed from another ecosystem,
// an entrypoint that already builds its own clients, and no hexagonal layout
// anywhere in sight.
func writeLegacyRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()

	files := map[string]string{
		"go.mod": "module github.com/acme/legacy-billing\n\ngo 1.25\n",
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
		"web/controllers/invoice_controller.go": `package controllers

import "net/http"

func Show(w http.ResponseWriter, r *http.Request) {}
`,
		"logic/invoice_service.go": `package logic

func Recalculate() error { return nil }
`,
		"data/invoice_dao.go": `package data

import "database/sql"

func Load(db *sql.DB) error { return nil }
`,
		"models/invoice.go": `package models

type Invoice struct{ ID string }
`,
	}

	for rel, content := range files {
		full := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", rel, err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatalf("write %s: %v", rel, err)
		}
	}
	return root
}

// TestAdopt_PreviewWritesNothing pins the safety property that matters most for
// this command: it runs against repositories somebody already depends on, so the
// default must be to propose rather than to act.
func TestAdopt_PreviewWritesNothing(t *testing.T) {
	root := writeLegacyRepo(t)
	svc := newRealFSService()

	if err := svc.AdoptProject(context.Background(), root, true, true); err != nil {
		t.Fatalf("AdoptProject preview failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(root, "aether.yaml")); err == nil {
		t.Error("preview wrote aether.yaml; adoption must not touch a repository until --apply")
	}
}

// TestAdopt_MapsLegacyLayoutIntoManifest is the test the previous implementation
// could never have passed: it wrote the default hexagonal paths no matter what
// the repository looked like, so every later generator aimed at directories that
// do not exist.
func TestAdopt_MapsLegacyLayoutIntoManifest(t *testing.T) {
	root := writeLegacyRepo(t)
	svc := newRealFSService()

	if err := svc.AdoptProject(context.Background(), root, true, false); err != nil {
		t.Fatalf("AdoptProject failed: %v", err)
	}

	manifest := readFileOrFail(t, filepath.Join(root, "aether.yaml"))

	// The real layout, not the default one.
	for _, want := range []string{
		"handler_http: web/controllers",
		"service: logic",
		"repository: data",
		"domain: models",
	} {
		if !strings.Contains(manifest, want) {
			t.Errorf("manifest is missing %q; the scan result was not applied\n--- aether.yaml ---\n%s", want, manifest)
		}
	}

	// The module path comes from go.mod. Adoption used to invent
	// "github.com/adopted/<dirname>", which is never the real import path, so
	// every generated import was wrong from the first file.
	if !strings.Contains(manifest, "github.com/acme/legacy-billing") {
		t.Errorf("module path was not read from go.mod\n--- aether.yaml ---\n%s", manifest)
	}
	if strings.Contains(manifest, "github.com/adopted/") {
		t.Errorf("a module path was invented instead of read\n--- aether.yaml ---\n%s", manifest)
	}

	// Clients the project already builds, so generated constructors accept them
	// rather than opening a second pool alongside.
	if !strings.Contains(manifest, "existing_db_var: pool") {
		t.Errorf("the existing database pool was not recorded\n--- aether.yaml ---\n%s", manifest)
	}
	if !strings.Contains(manifest, "existing_redis_var: cache") {
		t.Errorf("the existing redis client was not recorded\n--- aether.yaml ---\n%s", manifest)
	}

	if !strings.Contains(manifest, "mode: brownfield") {
		t.Errorf("adopted project was not marked brownfield\n--- aether.yaml ---\n%s", manifest)
	}
}

// TestAdopt_UnrecognisedLayoutRaisesAnomalyMode proves anomaly_mode is a
// measurement rather than a label somebody types. A layout the scanner could not
// read is exactly the case where later commands must not assume the standard
// tree.
func TestAdopt_UnrecognisedLayoutRaisesAnomalyMode(t *testing.T) {
	root := t.TempDir()
	write := func(rel, content string) {
		full := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", rel, err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatalf("write %s: %v", rel, err)
		}
	}

	write("go.mod", "module example.com/opaque\n\ngo 1.25\n")
	write("zz/qq/thing.go", "package qq\n\ntype Thing struct{}\n")

	svc := newRealFSService()
	if err := svc.AdoptProject(context.Background(), root, true, false); err != nil {
		t.Fatalf("AdoptProject failed: %v", err)
	}

	manifest := readFileOrFail(t, filepath.Join(root, "aether.yaml"))
	if !strings.Contains(manifest, "anomaly_mode: true") {
		t.Errorf("an unrecognisable layout did not raise anomaly_mode\n--- aether.yaml ---\n%s", manifest)
	}
}

// TestAdopt_ReadsGoModWithByteOrderMark banks a defect found by running the
// built binary rather than by reading the code.
//
// Editors on Windows routinely save UTF-8 with a BOM. A BOM on the first line of
// go.mod turns "module" into something that no longer matches, so a repository
// the user can plainly see is a Go module was reported as not being one.
func TestAdopt_ReadsGoModWithByteOrderMark(t *testing.T) {
	root := t.TempDir()
	bom := "\ufeff"
	if err := os.WriteFile(filepath.Join(root, "go.mod"),
		[]byte(bom+"module github.com/acme/bom-project\n\ngo 1.25\n"), 0o644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}

	svc := newRealFSService()
	if err := svc.AdoptProject(context.Background(), root, true, false); err != nil {
		t.Fatalf("a go.mod with a BOM must still be readable: %v", err)
	}

	manifest := readFileOrFail(t, filepath.Join(root, "aether.yaml"))
	if !strings.Contains(manifest, "github.com/acme/bom-project") {
		t.Errorf("module path was not recovered from a BOM-prefixed go.mod\n--- aether.yaml ---\n%s", manifest)
	}
}

// TestAdopt_RefusesRepositoryWithoutGoMod stops adoption from producing a
// manifest whose import paths could never resolve.
func TestAdopt_RefusesRepositoryWithoutGoMod(t *testing.T) {
	root := t.TempDir()
	svc := newRealFSService()

	err := svc.AdoptProject(context.Background(), root, true, false)
	if err == nil {
		t.Fatal("adoption succeeded in a directory that is not a Go module")
	}
	if !strings.Contains(err.Error(), "go.mod") {
		t.Errorf("error should name the missing go.mod, got: %v", err)
	}
	if _, statErr := os.Stat(filepath.Join(root, "aether.yaml")); statErr == nil {
		t.Error("a rejected adoption still wrote aether.yaml")
	}
}
