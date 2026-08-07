// Package regression, continued: filesystem boundaries and transaction
// atomicity.
//
// Reference: docs/rfc/20260807-v0.4.0-production-hardening.md §4 anomalies 6, 7.
package regression_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/muhananaufal/go-aether/internal/adapter/manifest"
	"github.com/muhananaufal/go-aether/internal/adapter/template"
	"github.com/muhananaufal/go-aether/internal/adapter/writer"
	"github.com/muhananaufal/go-aether/internal/core/domain"
	"github.com/muhananaufal/go-aether/internal/core/port"
	"github.com/muhananaufal/go-aether/internal/core/service"
	"github.com/spf13/afero"
)

// TestRegression_Issue20260807_DeployTargetCannotEscapeProject proves the
// generator refuses to write outside the project root.
//
// `infra:deploy ../../../pwned` was previously stopped only because the template
// lookup failed first, and the resulting error quoted the traversal path back to
// the terminal. Ordering is not a security boundary.
func TestRegression_Issue20260807_DeployTargetCannotEscapeProject(t *testing.T) {
	svc, w, _ := newFixture(t, fstest.MapFS{
		"cloud/k8s_deployment.yaml.tmpl": &fstest.MapFile{Data: []byte("kind: Deployment\n")},
	})
	ctx := context.Background()

	for _, target := range []string{"../../../pwned", `..\..\pwned`, "helm"} {
		t.Run(target, func(t *testing.T) {
			err := svc.AddDeploy(ctx, projectDir, target, false, false)
			if err == nil {
				t.Fatalf("AddDeploy accepted %q", target)
			}
			if !errors.Is(err, domain.ErrUnsupportedStack) && !errors.Is(err, domain.ErrPathEscape) {
				t.Errorf("expected a typed rejection, got %v", err)
			}
			if strings.Contains(err.Error(), "pwned_deployment.yaml.tmpl") {
				t.Errorf("error echoes the attacker-controlled template path: %v", err)
			}
		})
	}

	// The supported target must still work, or the test above would be satisfied
	// by a generator that refuses everything.
	if err := svc.AddDeploy(ctx, projectDir, "k8s", false, false); err != nil {
		t.Fatalf("AddDeploy(k8s) must succeed, got %v", err)
	}
	if exists, _ := w.Exists(projectDir + "/deploy/k8s.yaml"); !exists {
		t.Error("the supported target produced no manifest")
	}
}

// failingWriter fails the Nth write, modelling a disk that fills up or a
// permission that is revoked partway through a multi-file generation.
type failingWriter struct {
	port.FileWriter
	failOn int
	seen   int
}

func (f *failingWriter) WriteFile(ctx context.Context, path string, content []byte, overwrite, dryRun bool) error {
	f.seen++
	if f.seen == f.failOn {
		return fmt.Errorf("simulated ENOSPC on write %d", f.seen)
	}
	return f.FileWriter.WriteFile(ctx, path, content, overwrite, dryRun)
}

func (f *failingWriter) BeginTransaction() port.TransactionalWriter {
	return writer.NewUOWWriter(f)
}

// TestRegression_Issue20260807_PartialWriteRollsBackCompletely proves the unit
// of work leaves no debris when the disk gives out mid-generation.
//
// A vertical slice writes five files. Failing on the fourth must remove the
// three that already landed, because a project containing a domain and a port
// but no service or repository does not compile, and the user has no way to know
// which files were the tool's.
func TestRegression_Issue20260807_PartialWriteRollsBackCompletely(t *testing.T) {
	memFS := afero.NewMemMapFs()
	base := writer.NewAferoWriter(memFS)
	failing := &failingWriter{FileWriter: base, failOn: 4}

	resolver := manifest.NewYamlResolver(base)
	ctx := context.Background()

	m := domain.NewDefaultManifest("proj", "example.com/proj", "1.25", "hexagonal", "postgres", "chi")
	if err := resolver.Save(ctx, projectDir, m, base); err != nil {
		t.Fatalf("fixture setup: %v", err)
	}

	templates := fstest.MapFS{
		"hexagonal/domain.go.tmpl":              &fstest.MapFile{Data: []byte("package domain\n")},
		"hexagonal/port.go.tmpl":                &fstest.MapFile{Data: []byte("package port\n")},
		"hexagonal/service.go.tmpl":             &fstest.MapFile{Data: []byte("package service\n")},
		"hexagonal/handler_http_chi.go.tmpl":    &fstest.MapFile{Data: []byte("package http\n")},
		"hexagonal/repository_postgres.go.tmpl": &fstest.MapFile{Data: []byte("package repository\n")},
	}

	svc := service.NewAetherScaffoldService(template.NewStdEngine(templates), resolver, failing, nil, nil, nil)

	err := svc.MakeModule(ctx, projectDir, "order", []string{"http"}, false, false, false, false)
	if err == nil {
		t.Fatal("MakeModule reported success although a write failed")
	}
	if !errors.Is(err, domain.ErrWriteFailed) {
		t.Errorf("expected ErrWriteFailed so the caller can distinguish disk trouble, got %v", err)
	}

	// Nothing from the slice may survive. Counting files under the layer
	// directories is the non-tautological assertion here: asserting only that
	// err != nil would pass even with three orphaned files on disk.
	var leftovers []string
	for _, dir := range []string{
		"/proj/internal/core/domain",
		"/proj/internal/core/port",
		"/proj/internal/core/service",
		"/proj/internal/adapter/handler/http",
		"/proj/internal/adapter/repository",
	} {
		entries, readErr := afero.ReadDir(memFS, dir)
		if readErr != nil {
			continue
		}
		for _, e := range entries {
			if !e.IsDir() {
				leftovers = append(leftovers, filepath.ToSlash(filepath.Join(dir, e.Name())))
			}
		}
	}

	if len(leftovers) > 0 {
		t.Errorf("rollback left %d orphaned file(s), so the project is half generated:\n  %s",
			len(leftovers), strings.Join(leftovers, "\n  "))
	}
}

// TestRegression_Issue20260807_ManifestSurvivesFailedGeneration confirms the
// SSOT is not mutated by a generation that never completed. A manifest listing a
// module whose files do not exist is exactly the corruption `doctor` reports.
func TestRegression_Issue20260807_ManifestSurvivesFailedGeneration(t *testing.T) {
	memFS := afero.NewMemMapFs()
	base := writer.NewAferoWriter(memFS)
	failing := &failingWriter{FileWriter: base, failOn: 2}

	resolver := manifest.NewYamlResolver(base)
	ctx := context.Background()

	m := domain.NewDefaultManifest("proj", "example.com/proj", "1.25", "hexagonal", "postgres", "chi")
	if err := resolver.Save(ctx, projectDir, m, base); err != nil {
		t.Fatalf("fixture setup: %v", err)
	}

	templates := fstest.MapFS{
		"hexagonal/domain.go.tmpl":              &fstest.MapFile{Data: []byte("package domain\n")},
		"hexagonal/port.go.tmpl":                &fstest.MapFile{Data: []byte("package port\n")},
		"hexagonal/service.go.tmpl":             &fstest.MapFile{Data: []byte("package service\n")},
		"hexagonal/handler_http_chi.go.tmpl":    &fstest.MapFile{Data: []byte("package http\n")},
		"hexagonal/repository_postgres.go.tmpl": &fstest.MapFile{Data: []byte("package repository\n")},
	}

	svc := service.NewAetherScaffoldService(template.NewStdEngine(templates), resolver, failing, nil, nil, nil)

	if err := svc.MakeModule(ctx, projectDir, "order", []string{"http"}, false, false, false, false); err == nil {
		t.Fatal("MakeModule reported success although a write failed")
	}

	reloaded, _, err := resolver.Resolve(ctx, projectDir)
	if err != nil {
		t.Fatalf("manifest unreadable after failed generation: %v", err)
	}
	if reloaded.HasModule("order") {
		t.Error("aether.yaml registers a module whose files were rolled back")
	}
}

// TestRegression_Issue20260807_ReservedNameLeavesNoArtefact closes the Win32
// device anomaly at the filesystem level rather than only at validation.
func TestRegression_Issue20260807_ReservedNameLeavesNoArtefact(t *testing.T) {
	svc, w, _ := newFixture(t, fstest.MapFS{
		"hexagonal/domain.go.tmpl": &fstest.MapFile{Data: []byte("package domain\n")},
	})

	if err := svc.MakeDomain(context.Background(), projectDir, "aux", false, false); err == nil {
		t.Fatal("MakeDomain accepted the reserved device name \"aux\"")
	}

	if exists, _ := w.Exists(projectDir + "/internal/core/domain/aux.go"); exists {
		t.Error("a file was created for a rejected identifier")
	}
	if _, err := os.Stat(filepath.Join(projectDir, "internal", "core", "domain", "aux.go")); err == nil {
		t.Error("a reserved-name file reached the real filesystem")
	}
}
