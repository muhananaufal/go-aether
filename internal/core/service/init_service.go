package service

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/muhananaufal/go-aether/internal/core/domain"
)

// detectGoVersion reports the major.minor version of the toolchain running this
// binary, e.g. "1.25".
//
// Hardcoding a version meant the manifest, the generated Dockerfile and the
// go.mod written by `go mod init` could disagree with each other. A container
// image pinned below the language version in go.mod fails to build, and the
// error names a dependency rather than the toolchain, which sends the reader
// looking in the wrong place.
func detectGoVersion() string {
	const fallback = "1.23"

	raw := strings.TrimPrefix(runtime.Version(), "go")
	parts := strings.Split(raw, ".")
	if len(parts) < 2 {
		return fallback
	}
	// Guard against development toolchains such as "devel go1.26-abc1234".
	if _, err := strconv.Atoi(parts[0]); err != nil {
		return fallback
	}
	if _, err := strconv.Atoi(parts[1]); err != nil {
		return fallback
	}
	return parts[0] + "." + parts[1]
}

// InitProject bootstraps a new greenfield project architecture and saves the SSOT manifest.
func (s *AetherScaffoldService) InitProject(ctx context.Context, destDir, projectName, moduleName, arch, dbDriver, router string, dryRun bool) error {
	// Validated before anything is created. An unsupported selection used to pass
	// straight through and yield a chi + Postgres project no matter what the user
	// asked for; failing here leaves the target directory untouched.
	arch, dbDriver, router = domain.NormalizeStack(arch, dbDriver, router)
	if err := domain.ValidateStackSelection(arch, dbDriver, router); err != nil {
		return err
	}

	manifest := domain.NewDefaultManifest(projectName, moduleName, detectGoVersion(), arch, dbDriver, router)

	if err := manifest.Validate(); err != nil {
		return fmt.Errorf("invalid default manifest parameters: %w", err)
	}

	manifestPath := filepath.Join(destDir, "aether.yaml")
	exists, err := s.fs.Exists(manifestPath)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("%w: %s", domain.ErrFileConflict, manifestPath)
	}

	if dryRun {
		return nil
	}

	if err := s.resolver.Save(ctx, destDir, manifest, s.fs); err != nil {
		return fmt.Errorf("failed writing initial manifest: %w", err)
	}

	// go.mod must exist before the skeleton is rendered, because the templates
	// embed the module path in their import statements. Dependency resolution,
	// however, must run *after* the files exist — see runModuleTidy below.
	onRealDisk := false
	if _, statErr := os.Stat(destDir); statErr == nil {
		onRealDisk = true
		if err := s.ensureGoModule(ctx, destDir, moduleName); err != nil {
			return err
		}
	}

	templateData := &domain.TemplateData{
		ModuleName:      moduleName,
		ModuleNameTitle: projectName,
		PackagePath:     moduleName,
		GoVersion:       manifest.Project.GoVersion,
		ArchPattern:     manifest.Architecture.Pattern,
		Router:          manifest.Stack.Router,
		DBDriver:        manifest.Stack.Database.Driver,
		Paths:           manifest.Architecture.Paths,
		Timestamp:       manifest.Project.CreatedAt,
		AetherVersion:   manifest.Project.AetherVersion,
	}

	// Rendered through a transaction so a failure midway leaves no half-written
	// skeleton behind, matching the guarantee every other generator already gives.
	tx := s.fs.BeginTransaction()
	for _, spec := range skeletonFiles(router, dbDriver) {
		outPath := filepath.Join(destDir, spec.dest)
		content, err := s.engine.Render(ctx, spec.template, templateData)
		if err != nil {
			return fmt.Errorf("failed to render %s: %w", spec.dest, err)
		}
		if err := tx.WriteFile(ctx, outPath, content, false, false); err != nil {
			return fmt.Errorf("failed to stage %s: %w", spec.dest, err)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("skeleton generation failed, transaction rolled back: %w", err)
	}

	// Only now can `go mod tidy` see the imports it is supposed to resolve.
	if onRealDisk {
		s.runModuleTidy(ctx, destDir)
	}

	return nil
}

// skeletonSpec pairs an embedded template with its destination inside the new
// project. A slice rather than a map because map iteration order is randomised,
// and a generator that reports its failures in a different order on every run is
// needlessly hard to debug.
type skeletonSpec struct {
	template string
	dest     string
}

// skeletonFiles resolves the concrete template set for a stack selection.
//
// The router and driver are part of the template *name* rather than a branch
// inside one giant template. Five routers expressed as nested conditionals in a
// single file becomes unreadable long before the fifth is added, and a syntax
// error in one branch breaks every other branch with it.
func skeletonFiles(router, dbDriver string) []skeletonSpec {
	specs := []skeletonSpec{
		{fmt.Sprintf("common/main_%s.go.tmpl", router), filepath.Join("cmd", "server", "main.go")},
		{"common/config_viper.go.tmpl", filepath.Join("pkg", "config", "config.go")},
	}

	// "none" means the service genuinely has no datastore, so emitting a
	// connection pool that nothing calls would be noise the reader has to
	// understand and then delete.
	if domain.HasDatabase(dbDriver) {
		specs = append(specs, skeletonSpec{
			template: fmt.Sprintf("common/db_%s.go.tmpl", dbDriver),
			dest:     filepath.Join("pkg", "database", "db.go"),
		})
	}

	return append(specs,
		skeletonSpec{"common/Makefile.tmpl", "Makefile"},
		skeletonSpec{"common/Dockerfile.tmpl", "Dockerfile"},
		skeletonSpec{"common/dockerignore.tmpl", ".dockerignore"},
		skeletonSpec{"common/env_example.tmpl", ".env.example"},
	)
}

// ensureGoModule creates go.mod unless the directory already is a module.
// Re-running `go mod init` in an existing module is a hard error from the
// toolchain, which would otherwise abort adoption of a directory that merely
// happens to already be initialised.
func (s *AetherScaffoldService) ensureGoModule(ctx context.Context, destDir, moduleName string) error {
	if _, err := os.Stat(filepath.Join(destDir, "go.mod")); err == nil {
		return nil
	}

	cmd := exec.CommandContext(ctx, "go", "mod", "init", moduleName)
	cmd.Dir = destDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to run go mod init: %s (%w)", strings.TrimSpace(string(out)), err)
	}
	return nil
}

// runModuleTidy resolves the imports the freshly rendered skeleton introduced.
//
// Deliberately non-fatal: `go mod tidy` needs the network, and a developer
// bootstrapping a project on a plane should still end up with a complete,
// correct source tree. Failing hard here would delete a perfectly good scaffold
// over a transient DNS error. The user is told exactly what to run instead.
func (s *AetherScaffoldService) runModuleTidy(ctx context.Context, destDir string) {
	tidyCtx, cancel := context.WithTimeout(ctx, moduleTidyTimeout)
	defer cancel()

	cmd := exec.CommandContext(tidyCtx, "go", "mod", "tidy")
	cmd.Dir = destDir
	out, err := cmd.CombinedOutput()
	if err == nil {
		return
	}

	fmt.Fprintf(os.Stderr,
		"\n⚠️  Could not resolve dependencies automatically (%v).\n"+
			"   The generated source is complete and correct; only go.mod is missing entries.\n"+
			"   Run this once you have network access:\n\n       cd %s && go mod tidy\n\n%s\n",
		err, destDir, strings.TrimSpace(string(out)))
}

// moduleTidyTimeout bounds the only network-dependent step in the CLI so a
// stalled proxy cannot hang `init` indefinitely.
const moduleTidyTimeout = 120 * time.Second
