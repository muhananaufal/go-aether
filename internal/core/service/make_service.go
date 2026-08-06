package service

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

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
