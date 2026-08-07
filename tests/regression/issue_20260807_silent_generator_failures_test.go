// Package regression locks the CLI's immunity to defects that already shipped.
//
// Reference: docs/rfc/20260807-v0.4.0-production-hardening.md §4 anomalies 3, 4.
//
// The defects banked here share one shape and it is the most damaging shape a
// developer tool can have: the command prints success while the filesystem
// disagrees. A hard error teaches; a false success misleads.
package regression_test

import (
	"context"
	"errors"
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

const projectDir = "/proj"

// newFixture builds a service over an in-memory filesystem that already contains
// a valid manifest, so each test can exercise one generator in isolation.
func newFixture(t *testing.T, files fstest.MapFS) (port.ScaffoldService, port.FileWriter, port.ManifestResolver) {
	t.Helper()

	memFS := afero.NewMemMapFs()
	w := writer.NewAferoWriter(memFS)
	resolver := manifest.NewYamlResolver(w)

	m := domain.NewDefaultManifest("proj", "example.com/proj", "1.23", "hexagonal", "postgres", "chi")
	if err := resolver.Save(context.Background(), projectDir, m, w); err != nil {
		t.Fatalf("fixture setup: could not save manifest: %v", err)
	}

	// nil scanner and detector: none of these cases exercise brownfield adoption,
	// and passing real ones would let a scan failure masquerade as a generator bug.
	svc := service.NewAetherScaffoldService(template.NewStdEngine(files), resolver, w, nil, nil)
	return svc, w, resolver
}

// TestRegression_Issue20260807_AddValidatorIsNotSelfConflicting proves that
// `api:validator` can succeed at all.
//
// It staged the identical destination path twice inside one transaction. The
// first write created the file, the second hit "already exists" without
// --force, and the unit of work rolled the whole thing back. The command was
// therefore impossible to complete on a clean project: guaranteed failure,
// dressed as a file conflict the user never caused.
func TestRegression_Issue20260807_AddValidatorIsNotSelfConflicting(t *testing.T) {
	svc, w, _ := newFixture(t, fstest.MapFS{
		"common/validator_playground.go.tmpl": &fstest.MapFile{
			Data: []byte("package common\n\n// ValidateStruct checks struct tags.\nfunc ValidateStruct(s any) error { return nil }\n"),
		},
	})

	err := svc.AddValidator(context.Background(), projectDir, "playground", false, false)
	if err != nil {
		t.Fatalf("AddValidator must succeed on a clean project, got: %v", err)
	}

	dest := projectDir + "/pkg/common/validator.go"
	exists, existErr := w.Exists(dest)
	if existErr != nil {
		t.Fatalf("could not stat %s: %v", dest, existErr)
	}
	if !exists {
		t.Fatalf("AddValidator reported success but %s was never written", dest)
	}

	// Non-tautological: assert the file carries the rendered contract, not merely
	// that some bytes landed on disk.
	raw, readErr := w.ReadFile(dest)
	if readErr != nil {
		t.Fatalf("could not read %s: %v", dest, readErr)
	}
	if !strings.Contains(string(raw), "func ValidateStruct") {
		t.Errorf("validator file does not contain the rendered function:\n%s", raw)
	}
}

// TestRegression_Issue20260807_AddCacheRejectsUnknownDriver proves the CLI stops
// claiming success for work it did not do.
//
// `cache:redis <unknown>` swallowed the template error via `if err == nil`,
// skipped every write, then persisted the mutated manifest and printed
// "⚡ Injected ... cache infrastructure". The project was left declaring a cache
// driver whose implementation does not exist anywhere on disk.
func TestRegression_Issue20260807_AddCacheRejectsUnknownDriver(t *testing.T) {
	svc, w, resolver := newFixture(t, fstest.MapFS{
		"common/cache_redis.go.tmpl": &fstest.MapFile{
			Data: []byte("package cache\n\n// Provider abstracts the cache.\ntype Provider interface{}\n"),
		},
	})
	ctx := context.Background()

	err := svc.AddCache(ctx, projectDir, "memcachedxyz", false, false)
	if err == nil {
		t.Fatal("AddCache accepted an unknown driver and reported success")
	}

	// The manifest is the project's source of truth. A failed command must not
	// leave it describing infrastructure that was never generated.
	reloaded, _, resolveErr := resolver.Resolve(ctx, projectDir)
	if resolveErr != nil {
		t.Fatalf("manifest unreadable after failed command: %v", resolveErr)
	}
	if reloaded.Stack.Cache == "memcachedxyz" {
		t.Errorf("failed AddCache still persisted cache=%q to aether.yaml", reloaded.Stack.Cache)
	}

	if exists, _ := w.Exists(projectDir + "/pkg/cache/memcachedxyz.go"); exists {
		t.Error("a provider file was written for a driver that has no template")
	}
}

// TestRegression_Issue20260807_AddCacheSucceedsForKnownDriver is the companion
// assertion. Rejecting the unknown case is only correct if the known case still
// works; without this the previous test could be satisfied by refusing
// everything.
func TestRegression_Issue20260807_AddCacheSucceedsForKnownDriver(t *testing.T) {
	svc, w, resolver := newFixture(t, fstest.MapFS{
		"common/cache_redis.go.tmpl": &fstest.MapFile{
			Data: []byte("package cache\n\n// Provider abstracts the cache.\ntype Provider interface{}\n"),
		},
	})
	ctx := context.Background()

	if err := svc.AddCache(ctx, projectDir, "redis", false, false); err != nil {
		t.Fatalf("AddCache(redis) must succeed, got: %v", err)
	}

	if exists, _ := w.Exists(projectDir + "/pkg/cache/redis.go"); !exists {
		t.Error("cache provider file was not generated for a supported driver")
	}

	reloaded, _, err := resolver.Resolve(ctx, projectDir)
	if err != nil {
		t.Fatalf("manifest unreadable: %v", err)
	}
	if reloaded.Stack.Cache != "redis" {
		t.Errorf("manifest should record cache=redis after success, got %q", reloaded.Stack.Cache)
	}
}

// TestRegression_Issue20260807_AddMiddlewareIsAllOrNothing proves the injector
// no longer mutates a user's handler when the middleware it depends on cannot be
// produced.
//
// The handler edit was staged unconditionally while the middleware write sat
// behind `if err == nil`. A missing template therefore produced a handler that
// calls middleware.RequireAuth() with no such package on disk: a project that
// no longer compiles, from a command that exited 0.
func TestRegression_Issue20260807_AddMiddlewareIsAllOrNothing(t *testing.T) {
	// Note the absent common/middleware_jwt.go.tmpl: this models a template that
	// fails to render for any reason.
	svc, w, _ := newFixture(t, fstest.MapFS{})
	ctx := context.Background()

	handlerPath := projectDir + "/internal/adapter/handler/http/order_handler.go"
	original := "package http\n\nfunc Route() {\n\t// @aether:inject:middleware\n}\n"
	if err := w.WriteFile(ctx, handlerPath, []byte(original), false, false); err != nil {
		t.Fatalf("fixture setup: could not seed handler: %v", err)
	}

	err := svc.AddMiddleware(ctx, projectDir, "order", "jwt-auth", false, false)
	if err == nil {
		t.Fatal("AddMiddleware reported success although the middleware template is unavailable")
	}

	after, readErr := w.ReadFile(handlerPath)
	if readErr != nil {
		t.Fatalf("handler disappeared after a failed injection: %v", readErr)
	}
	if string(after) != original {
		t.Errorf("handler was modified despite the command failing.\nbefore:\n%s\nafter:\n%s", original, after)
	}
}

// TestRegression_Issue20260807_ReservedDeviceNameRejected pins the Win32 device
// anomaly at the domain layer, where the rule belongs, so it holds for every
// generator rather than only the ones an end-to-end test happens to exercise.
func TestRegression_Issue20260807_ReservedDeviceNameRejected(t *testing.T) {
	for _, name := range []string{"con", "prn", "aux", "nul", "com1", "lpt9"} {
		t.Run(name, func(t *testing.T) {
			err := domain.ValidateGoIdentifier(name)
			if err == nil {
				t.Fatalf("ValidateGoIdentifier(%q) accepted a Win32 device name", name)
			}
			if !errors.Is(err, domain.ErrReservedName) {
				t.Errorf("expected ErrReservedName for %q, got %v", name, err)
			}
		})
	}

	// Guard against an over-broad rule: names that merely start with a reserved
	// prefix are perfectly legal and must keep working.
	for _, name := range []string{"connection", "console", "auxiliary", "company"} {
		t.Run("allowed_"+name, func(t *testing.T) {
			if err := domain.ValidateGoIdentifier(name); err != nil {
				t.Errorf("ValidateGoIdentifier(%q) must be accepted, got %v", name, err)
			}
		})
	}
}
