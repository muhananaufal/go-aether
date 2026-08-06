package service

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/muhananaufal/go-aether/internal/adapter/writer"
	"github.com/muhananaufal/go-aether/internal/core/domain"
	"github.com/muhananaufal/go-aether/internal/core/port"
)

// AetherScaffoldService implements port.ScaffoldService to orchestrate project architecture generation.
type AetherScaffoldService struct {
	engine   port.TemplateEngine
	resolver port.ManifestResolver
	fs       port.FileWriter
}

// NewAetherScaffoldService constructs the main orchestrator engine.
func NewAetherScaffoldService(engine port.TemplateEngine, resolver port.ManifestResolver, fs port.FileWriter) *AetherScaffoldService {
	return &AetherScaffoldService{
		engine:   engine,
		resolver: resolver,
		fs:       fs,
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
	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	archPrefix := manifest.Architecture.Pattern

	// Define mapping of target layers to templates
	layers := map[string]string{
		"domain":       manifest.Architecture.Paths.Domain,
		"port":         manifest.Architecture.Paths.Port,
		"service":      manifest.Architecture.Paths.Service,
		"handler_http": manifest.Architecture.Paths.HandlerHTTP,
		"repository":   manifest.Architecture.Paths.Repository,
	}

	for layer, relPath := range layers {
		// e.g. "hexagonal/service.go.tmpl"
		tmplName := fmt.Sprintf("%s/%s", archPrefix, layer)
		if layer == "repository" {
			tmplName += fmt.Sprintf("_%s.go.tmpl", manifest.Stack.Database.Driver)
		} else {
			tmplName += ".go.tmpl"
		}

		content, err := s.engine.Render(ctx, tmplName, data)
		if err != nil {
			return err
		}
		
		// e.g. "internal/core/service/order_service.go"
		fileName := fmt.Sprintf("%s_%s.go", strings.ToLower(moduleName), strings.Split(layer, "_")[0])
		if layer == "domain" {
			fileName = fmt.Sprintf("%s.go", strings.ToLower(moduleName))
		}
		
		destFile := filepath.Join(projectRoot, relPath, fileName)
		tx.Stage(destFile, content, force, dryRun)
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

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	archPrefix := manifest.Architecture.Pattern

	tmplName := fmt.Sprintf("%s/service.go.tmpl", archPrefix)
	content, err := s.engine.Render(ctx, tmplName, data)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s_service.go", strings.ToLower(moduleName))
	destFile := filepath.Join(projectRoot, manifest.Architecture.Paths.Service, fileName)
	tx.Stage(destFile, content, force, dryRun)

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

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	archPrefix := manifest.Architecture.Pattern

	tmplName := fmt.Sprintf("%s/handler_%s.go.tmpl", archPrefix, transport)
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
	tx.Stage(destFile, content, force, dryRun)

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

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	archPrefix := manifest.Architecture.Pattern

	tmplName := fmt.Sprintf("%s/domain_only.go.tmpl", archPrefix)
	content, err := s.engine.Render(ctx, tmplName, data)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s.go", strings.ToLower(moduleName))
	destFile := filepath.Join(projectRoot, manifest.Architecture.Paths.Domain, fileName)
	tx.Stage(destFile, content, force, dryRun)

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

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	archPrefix := manifest.Architecture.Pattern

	tmplName := fmt.Sprintf("%s/port_only.go.tmpl", archPrefix)
	content, err := s.engine.Render(ctx, tmplName, data)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s_port.go", strings.ToLower(moduleName))
	destFile := filepath.Join(projectRoot, manifest.Architecture.Paths.Port, fileName)
	tx.Stage(destFile, content, force, dryRun)

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

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	archPrefix := manifest.Architecture.Pattern

	tmplName := fmt.Sprintf("%s/repository_only.go.tmpl", archPrefix)
	content, err := s.engine.Render(ctx, tmplName, data)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s_repository.go", strings.ToLower(moduleName))
	destFile := filepath.Join(projectRoot, manifest.Architecture.Paths.Repository, fileName)
	tx.Stage(destFile, content, force, dryRun)

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

	tx := writer.NewTransactionalBuffer(s.fs)
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

	tx.Stage(upDest, upContent, force, dryRun)
	tx.Stage(downDest, downContent, force, dryRun)

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

	tx := writer.NewTransactionalBuffer(s.fs)
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
	tx.Stage(filepath.Join(destDir, fileName), content, force, dryRun)

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

	tx := writer.NewTransactionalBuffer(s.fs)
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
	tx.Stage(filepath.Join(destDir, fileName), content, force, dryRun)

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

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	
	content, err := s.engine.Render(ctx, "common/seeder.go.tmpl", data)
	if err != nil {
		return err
	}

	destFile := filepath.Join(projectRoot, "cmd", "seeder", fmt.Sprintf("%s.go", strings.ToLower(name)))
	tx.Stage(destFile, content, force, dryRun)

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

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)

	content, err := s.engine.Render(ctx, "hexagonal/valueobject.go.tmpl", data)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s_vo.go", strings.ToLower(name))
	destFile := filepath.Join(projectRoot, manifest.Architecture.Paths.Domain, fileName)
	tx.Stage(destFile, content, force, dryRun)

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

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)

	content, err := s.engine.Render(ctx, "hexagonal/aggregate.go.tmpl", data)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s_aggregate.go", strings.ToLower(name))
	destFile := filepath.Join(projectRoot, manifest.Architecture.Paths.Domain, fileName)
	tx.Stage(destFile, content, force, dryRun)

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

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)

	content, err := s.engine.Render(ctx, "hexagonal/event.go.tmpl", data)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s_event.go", strings.ToLower(name))
	destFile := filepath.Join(projectRoot, manifest.Architecture.Paths.Domain, fileName)
	tx.Stage(destFile, content, force, dryRun)

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

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)

	content, err := s.engine.Render(ctx, "hexagonal/command.go.tmpl", data)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s_command.go", strings.ToLower(name))
	destFile := filepath.Join(projectRoot, manifest.Architecture.Paths.Service, fileName)
	tx.Stage(destFile, content, force, dryRun)

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

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)

	content, err := s.engine.Render(ctx, "hexagonal/query.go.tmpl", data)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s_query.go", strings.ToLower(name))
	destFile := filepath.Join(projectRoot, manifest.Architecture.Paths.Service, fileName)
	tx.Stage(destFile, content, force, dryRun)

	return tx.Commit(ctx)
}
