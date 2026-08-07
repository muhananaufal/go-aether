package service

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/muhananaufal/go-aether/internal/core/domain"
	"github.com/muhananaufal/go-aether/internal/core/port"
)

// AetherScaffoldService implements port.ScaffoldService to orchestrate project architecture generation.
type AetherScaffoldService struct {
	engine   port.TemplateEngine
	resolver port.ManifestResolver
	fs       port.FileWriter

	// scanner and byoDetector are only consulted by brownfield adoption. They are
	// injected rather than constructed here so the core keeps depending on ports
	// alone, and a missing wiring is caught at the call site instead of surfacing
	// as a nil dereference during someone's first `adopt`.
	scanner     port.LayoutScanner
	byoDetector port.BYODetector
}

// NewAetherScaffoldService constructs the main orchestrator engine.
//
// scanner and byoDetector may be nil for callers that never adopt an existing
// repository; AdoptProject reports the omission rather than panicking.
func NewAetherScaffoldService(
	engine port.TemplateEngine,
	resolver port.ManifestResolver,
	fs port.FileWriter,
	scanner port.LayoutScanner,
	byoDetector port.BYODetector,
) *AetherScaffoldService {
	return &AetherScaffoldService{
		engine:      engine,
		resolver:    resolver,
		fs:          fs,
		scanner:     scanner,
		byoDetector: byoDetector,
	}
}

// layerTemplate resolves the embedded template for one layer of a vertical slice.
//
// Two layers vary with the project's stack. The HTTP handler is router-specific
// because chi, gin, echo, fiber and net/http share no handler signature; the
// repository is driver-specific because upsert syntax and the "no rows" sentinel
// differ per database. Expressing that as separate templates rather than
// branches inside one file keeps each readable and stops a syntax error in one
// dialect from breaking the others.
func layerTemplate(archPrefix, layer string, manifest *domain.AetherManifest) string {
	switch layer {
	case "handler_http":
		router := manifest.Stack.Router
		if router == "" {
			router = domain.DefaultRouter
		}
		return fmt.Sprintf("%s/handler_http_%s.go.tmpl", archPrefix, strings.ToLower(router))
	case "repository":
		driver := manifest.Stack.Database.Driver
		if driver == "" {
			driver = domain.DefaultDBDriver
		}
		return fmt.Sprintf("%s/repository_%s.go.tmpl", archPrefix, strings.ToLower(driver))
	default:
		return fmt.Sprintf("%s/%s.go.tmpl", archPrefix, layer)
	}
}

// MakeModule generates all vertical slice components for a new module feature.
func (s *AetherScaffoldService) MakeModule(ctx context.Context, startDir, moduleName string, transports []string, hasCache, hasWorker, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	if manifest.HasModule(moduleName) && !force {
		return fmt.Errorf("%w: %s", domain.ErrModuleAlreadyExists, moduleName)
	}

	data, err := domain.NewTemplateData(moduleName, manifest, transports, hasCache, hasWorker)
	if err != nil {
		return err
	}

	// Initialize the transactional buffer
	tx := s.fs.BeginTransaction()
	projectRoot := filepath.Dir(manifestPath)
	archPrefix := manifest.Architecture.Pattern

	// A slice, not a map: map iteration order is randomised, so a failure in one
	// layer would be reported before or after the others at random, and partial
	// output would differ between runs.
	layers := []struct{ name, relPath string }{
		{"domain", manifest.Architecture.Paths.Domain},
		{"port", manifest.Architecture.Paths.Port},
		{"service", manifest.Architecture.Paths.Service},
		{"handler_http", manifest.Architecture.Paths.HandlerHTTP},
		{"repository", manifest.Architecture.Paths.Repository},
	}

	for _, layer := range layers {
		content, err := s.engine.Render(ctx, layerTemplate(archPrefix, layer.name, manifest), data)
		if err != nil {
			return err
		}

		// e.g. "internal/core/service/order_service.go"
		fileName := fmt.Sprintf("%s_%s.go", strings.ToLower(moduleName), strings.Split(layer.name, "_")[0])
		if layer.name == "domain" {
			fileName = fmt.Sprintf("%s.go", strings.ToLower(moduleName))
		}

		destFile := filepath.Join(projectRoot, layer.relPath, fileName)
		if err := tx.WriteFile(ctx, destFile, content, force, dryRun); err != nil {
			return err
		}
	}

	// Update manifest
	if !manifest.HasModule(moduleName) {
		reg := domain.ModuleRegistry{
			Name:       moduleName,
			Transports: transports,
			HasCache:   hasCache,
			HasWorker:  hasWorker,
		}
		_ = manifest.AddModule(reg)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("module generation failed during disk write, transaction rolled back: %w", err)
	}

	if !dryRun {
		_ = s.resolver.Save(ctx, manifestPath, manifest, s.fs)
	}

	return nil
}

// MakeService generates only the service layer component for a specific module.
func (s *AetherScaffoldService) MakeService(ctx context.Context, startDir, moduleName string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData(moduleName, manifest, []string{"http"}, false, false)
	if err != nil {
		return err
	}

	tx := s.fs.BeginTransaction()
	projectRoot := filepath.Dir(manifestPath)
	archPrefix := manifest.Architecture.Pattern

	tmplName := fmt.Sprintf("%s/service.go.tmpl", archPrefix)
	content, err := s.engine.Render(ctx, tmplName, data)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s_service.go", strings.ToLower(moduleName))
	destFile := filepath.Join(projectRoot, manifest.Architecture.Paths.Service, fileName)
	tx.WriteFile(ctx, destFile, content, force, dryRun)

	return tx.Commit(ctx)
}

// MakeHandler generates only the transport handler component for a specific module.
func (s *AetherScaffoldService) MakeHandler(ctx context.Context, startDir, moduleName, transport string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData(moduleName, manifest, []string{transport}, false, false)
	if err != nil {
		return err
	}

	tx := s.fs.BeginTransaction()
	projectRoot := filepath.Dir(manifestPath)
	archPrefix := manifest.Architecture.Pattern

	tmplName := fmt.Sprintf("%s/handler_%s.go.tmpl", archPrefix, transport)
	if transport == "http" {
		tmplName = layerTemplate(archPrefix, "handler_http", manifest)
	}
	content, err := s.engine.Render(ctx, tmplName, data)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s_handler.go", strings.ToLower(moduleName))

	destDir := manifest.Architecture.Paths.HandlerHTTP
	if transport == "grpc" {
		destDir = manifest.Architecture.Paths.HandlerGRPC
		if destDir == "" {
			destDir = "internal/adapter/handler/grpc"
		}
	} else if transport != "http" {
		// fallback for other transports
		destDir = fmt.Sprintf("internal/adapter/handler/%s", transport)
	}

	destFile := filepath.Join(projectRoot, destDir, fileName)
	tx.WriteFile(ctx, destFile, content, force, dryRun)

	return tx.Commit(ctx)
}

// MakeDomain generates only the domain layer entity for a specific module.
func (s *AetherScaffoldService) MakeDomain(ctx context.Context, startDir, moduleName string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData(moduleName, manifest, []string{"http"}, false, false)
	if err != nil {
		return err
	}

	tx := s.fs.BeginTransaction()
	projectRoot := filepath.Dir(manifestPath)
	archPrefix := manifest.Architecture.Pattern

	// The canonical domain template, the same one arch:module uses. The former
	// domain_only variant carried neither Validate nor the sentinel errors, so a
	// slice assembled command by command failed to compile the moment the service
	// and handler referenced them.
	content, err := s.engine.Render(ctx, layerTemplate(archPrefix, "domain", manifest), data)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s.go", strings.ToLower(moduleName))
	destFile := filepath.Join(projectRoot, manifest.Architecture.Paths.Domain, fileName)
	tx.WriteFile(ctx, destFile, content, force, dryRun)

	return tx.Commit(ctx)
}

// MakePort generates only the port interface contract for a specific module.
func (s *AetherScaffoldService) MakePort(ctx context.Context, startDir, moduleName string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData(moduleName, manifest, []string{"http"}, false, false)
	if err != nil {
		return err
	}

	tx := s.fs.BeginTransaction()
	projectRoot := filepath.Dir(manifestPath)
	archPrefix := manifest.Architecture.Pattern

	// The canonical port template. port_only declared the repository contract
	// alone, so every service and handler generated afterwards referenced a
	// port.XService that did not exist.
	content, err := s.engine.Render(ctx, layerTemplate(archPrefix, "port", manifest), data)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s_port.go", strings.ToLower(moduleName))
	destFile := filepath.Join(projectRoot, manifest.Architecture.Paths.Port, fileName)
	tx.WriteFile(ctx, destFile, content, force, dryRun)

	return tx.Commit(ctx)
}

// MakeRepository generates only the infrastructure repository for a specific module.
func (s *AetherScaffoldService) MakeRepository(ctx context.Context, startDir, moduleName string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData(moduleName, manifest, []string{"http"}, false, false)
	if err != nil {
		return err
	}

	tx := s.fs.BeginTransaction()
	projectRoot := filepath.Dir(manifestPath)
	archPrefix := manifest.Architecture.Pattern

	// The canonical, driver-specific repository. repository_only produced an
	// adapter with no connection at all, which compiled in isolation and then
	// diverged from whatever arch:module had already generated.
	content, err := s.engine.Render(ctx, layerTemplate(archPrefix, "repository", manifest), data)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s_repository.go", strings.ToLower(moduleName))
	destFile := filepath.Join(projectRoot, manifest.Architecture.Paths.Repository, fileName)
	tx.WriteFile(ctx, destFile, content, force, dryRun)

	return tx.Commit(ctx)
}

// MakeMigration generates a SQL migration file pair (up/down).
func (s *AetherScaffoldService) MakeMigration(ctx context.Context, startDir, name string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData(name, manifest, nil, false, false)
	if err != nil {
		return err
	}

	tx := s.fs.BeginTransaction()
	projectRoot := filepath.Dir(manifestPath)

	upContent, err := s.engine.Render(ctx, "common/migration_up.sql.tmpl", data)
	if err != nil {
		return err
	}
	downContent, err := s.engine.Render(ctx, "common/migration_down.sql.tmpl", data)
	if err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102150405")
	baseName := fmt.Sprintf("%s_%s", timestamp, strings.ToLower(name))

	upDest := filepath.Join(projectRoot, "migrations", fmt.Sprintf("%s.up.sql", baseName))
	downDest := filepath.Join(projectRoot, "migrations", fmt.Sprintf("%s.down.sql", baseName))

	tx.WriteFile(ctx, upDest, upContent, force, dryRun)
	tx.WriteFile(ctx, downDest, downContent, force, dryRun)

	return tx.Commit(ctx)
}

// MakePipeline scaffolds Fan-Out / Fan-In bounded concurrency pipeline helper.
func (s *AetherScaffoldService) MakePipeline(ctx context.Context, startDir, name string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData(name+"_pipeline", manifest, nil, false, false)
	if err != nil {
		return err
	}

	tx := s.fs.BeginTransaction()
	projectRoot := filepath.Dir(manifestPath)
	pkgDir := manifest.Architecture.Paths.Pkg
	if pkgDir == "" {
		pkgDir = "pkg"
	}
	destDir := filepath.Join(projectRoot, pkgDir, "concurrency")

	content, err := s.engine.Render(ctx, "plugins/pipeline.go.tmpl", data)
	if err != nil {
		return err
	}
	// Prefix with name to allow multiple pipelines
	fileName := fmt.Sprintf("%s_pipeline.go", strings.ToLower(name))
	tx.WriteFile(ctx, filepath.Join(destDir, fileName), content, force, dryRun)

	return tx.Commit(ctx)
}

// MakeSpecification scaffolds reusable DDD Specification pattern for dynamic query rules.
func (s *AetherScaffoldService) MakeSpecification(ctx context.Context, startDir, name string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData(name, manifest, nil, false, false)
	if err != nil {
		return err
	}

	tx := s.fs.BeginTransaction()
	projectRoot := filepath.Dir(manifestPath)
	domainDir := manifest.Architecture.Paths.Domain
	if domainDir == "" {
		domainDir = filepath.Join("internal", "core", "domain")
	}
	destDir := filepath.Join(projectRoot, domainDir)

	content, err := s.engine.Render(ctx, "plugins/specification.go.tmpl", data)
	if err != nil {
		return err
	}
	fileName := fmt.Sprintf("%s_specification.go", strings.ToLower(name))
	tx.WriteFile(ctx, filepath.Join(destDir, fileName), content, force, dryRun)

	return tx.Commit(ctx)
}

// MakeSeeder generates a database seeder file.
func (s *AetherScaffoldService) MakeSeeder(ctx context.Context, startDir, name string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData(name, manifest, nil, false, false)
	if err != nil {
		return err
	}

	tx := s.fs.BeginTransaction()
	projectRoot := filepath.Dir(manifestPath)

	content, err := s.engine.Render(ctx, "common/seeder.go.tmpl", data)
	if err != nil {
		return err
	}

	destFile := filepath.Join(projectRoot, "cmd", "seeder", fmt.Sprintf("%s.go", strings.ToLower(name)))
	tx.WriteFile(ctx, destFile, content, force, dryRun)

	return tx.Commit(ctx)
}

// MakeValueObject generates an immutable DDD Value Object struct.
func (s *AetherScaffoldService) MakeValueObject(ctx context.Context, startDir, name string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData(name, manifest, nil, false, false)
	if err != nil {
		return err
	}

	tx := s.fs.BeginTransaction()
	projectRoot := filepath.Dir(manifestPath)

	content, err := s.engine.Render(ctx, "hexagonal/valueobject.go.tmpl", data)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s_vo.go", strings.ToLower(name))
	destFile := filepath.Join(projectRoot, manifest.Architecture.Paths.Domain, fileName)
	tx.WriteFile(ctx, destFile, content, force, dryRun)

	return tx.Commit(ctx)
}

// MakeAggregate generates a DDD Aggregate Root entity with event recording.
func (s *AetherScaffoldService) MakeAggregate(ctx context.Context, startDir, name string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData(name, manifest, nil, false, false)
	if err != nil {
		return err
	}

	tx := s.fs.BeginTransaction()
	projectRoot := filepath.Dir(manifestPath)

	content, err := s.engine.Render(ctx, "hexagonal/aggregate.go.tmpl", data)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s_aggregate.go", strings.ToLower(name))
	destFile := filepath.Join(projectRoot, manifest.Architecture.Paths.Domain, fileName)
	tx.WriteFile(ctx, destFile, content, force, dryRun)

	return tx.Commit(ctx)
}

// MakeEvent generates a Domain Event struct and serializer.
func (s *AetherScaffoldService) MakeEvent(ctx context.Context, startDir, name string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData(name, manifest, nil, false, false)
	if err != nil {
		return err
	}

	tx := s.fs.BeginTransaction()
	projectRoot := filepath.Dir(manifestPath)

	content, err := s.engine.Render(ctx, "hexagonal/event.go.tmpl", data)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s_event.go", strings.ToLower(name))
	destFile := filepath.Join(projectRoot, manifest.Architecture.Paths.Domain, fileName)
	tx.WriteFile(ctx, destFile, content, force, dryRun)

	return tx.Commit(ctx)
}

// MakeCommand generates a CQRS Command DTO and execution handler.
func (s *AetherScaffoldService) MakeCommand(ctx context.Context, startDir, name string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData(name, manifest, nil, false, false)
	if err != nil {
		return err
	}

	tx := s.fs.BeginTransaction()
	projectRoot := filepath.Dir(manifestPath)

	content, err := s.engine.Render(ctx, "hexagonal/command.go.tmpl", data)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s_command.go", strings.ToLower(name))
	destFile := filepath.Join(projectRoot, manifest.Architecture.Paths.Service, fileName)
	tx.WriteFile(ctx, destFile, content, force, dryRun)

	return tx.Commit(ctx)
}

// MakeQuery generates a CQRS Query DTO and read-model handler.
func (s *AetherScaffoldService) MakeQuery(ctx context.Context, startDir, name string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData(name, manifest, nil, false, false)
	if err != nil {
		return err
	}

	tx := s.fs.BeginTransaction()
	projectRoot := filepath.Dir(manifestPath)

	content, err := s.engine.Render(ctx, "hexagonal/query.go.tmpl", data)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s_query.go", strings.ToLower(name))
	destFile := filepath.Join(projectRoot, manifest.Architecture.Paths.Service, fileName)
	tx.WriteFile(ctx, destFile, content, force, dryRun)

	return tx.Commit(ctx)
}

// MakeCursorPaginator scaffolds cursor-based opaque base64 pagination helper.
func (s *AetherScaffoldService) MakeCursorPaginator(ctx context.Context, startDir string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData("cursor_paginator", manifest, nil, false, false)
	if err != nil {
		return err
	}

	tx := s.fs.BeginTransaction()
	projectRoot := filepath.Dir(manifestPath)

	pkgDir := manifest.Architecture.Paths.Pkg
	if pkgDir == "" {
		pkgDir = "pkg"
	}
	destDir := filepath.Join(projectRoot, pkgDir, "pagination")

	content, err := s.engine.Render(ctx, "plugins/cursor_paginator.go.tmpl", data)
	if err != nil {
		return err
	}
	tx.WriteFile(ctx, filepath.Join(destDir, "cursor.go"), content, force, dryRun)

	return tx.Commit(ctx)
}
