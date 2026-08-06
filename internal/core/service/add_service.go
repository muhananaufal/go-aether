package service

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/muhananaufal/go-aether/internal/adapter/writer"
	"github.com/muhananaufal/go-aether/internal/core/domain"
)

// AddMiddleware injects middleware logic into existing adapters.
func (s *AetherScaffoldService) AddMiddleware(ctx context.Context, startDir, moduleName, middlewareType string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}
	
	projectRoot := filepath.Dir(manifestPath)
	handlerFile := filepath.Join(projectRoot, manifest.Architecture.Paths.HandlerHTTP, fmt.Sprintf("%s_handler.go", strings.ToLower(moduleName)))
	
	exists, err := s.fs.Exists(handlerFile)
	if err != nil || !exists {
		return fmt.Errorf("handler file not found for module %s at %s", moduleName, handlerFile)
	}

	contentBytes, err := s.fs.ReadFile(handlerFile)
	if err != nil {
		return err
	}
	content := string(contentBytes)

	marker := "// @aether:inject:middleware"
	if !strings.Contains(content, marker) {
		return fmt.Errorf("marker %q not found in %s. Please add it inside your chi.Route definition", marker, handlerFile)
	}

	var injection string
	var tmplName string
	var destPkgFile string
	switch middlewareType {
	case "jwt-auth":
		injection = "\t\tr.Use(middleware.RequireAuth())\n"
		tmplName = "common/middleware_jwt.go.tmpl"
		destPkgFile = "jwt_auth.go"
	case "rate-limit":
		injection = "\t\tr.Use(middleware.RateLimit())\n"
		tmplName = "common/middleware_ratelimit.go.tmpl"
		destPkgFile = "rate_limit.go"
	default:
		return fmt.Errorf("unknown middleware type: %s", middlewareType)
	}

	// Prevent duplicate injection
	if strings.Contains(content, strings.TrimSpace(injection)) {
		return nil // already injected
	}

	newContent := strings.Replace(content, marker, marker+"\n"+injection, 1)

	// Additionally, generate the middleware logic itself into pkg/middleware/
	tx := writer.NewTransactionalBuffer(s.fs)
	
	// Write the modified handler (always overwrite because it already exists)
	tx.Stage(handlerFile, []byte(newContent), true, dryRun)
	
	// Write the middleware pkg file
	data, _ := domain.NewTemplateData("middleware", manifest, nil, false, false)
	mwContent, err := s.engine.Render(ctx, tmplName, data)
	if err == nil {
		mwDest := filepath.Join(projectRoot, manifest.Architecture.Paths.Pkg, "middleware", destPkgFile)
		tx.Stage(mwDest, mwContent, force, dryRun)
	}

	return tx.Commit(ctx)
}

// AddCache sets up the global cache layer configuration and generates the cache provider infrastructure.
func (s *AetherScaffoldService) AddCache(ctx context.Context, startDir, cacheType string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	manifest.Stack.Cache = cacheType
	
	data, _ := domain.NewTemplateData("cache", manifest, nil, false, false)
	tmplName := fmt.Sprintf("common/cache_%s.go.tmpl", cacheType)
	content, err := s.engine.Render(ctx, tmplName, data)
	if err == nil {
		tx := writer.NewTransactionalBuffer(s.fs)
		destFile := filepath.Join(filepath.Dir(manifestPath), manifest.Architecture.Paths.Pkg, "cache", fmt.Sprintf("%s.go", cacheType))
		tx.Stage(destFile, content, force, dryRun)
		if err := tx.Commit(ctx); err != nil {
			return err
		}
	}

	if !dryRun {
		return s.resolver.Save(ctx, manifestPath, manifest, s.fs)
	}
	return nil
}

// AddTransport registers a new global transport protocol.
func (s *AetherScaffoldService) AddTransport(ctx context.Context, startDir, transport string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	found := false
	for _, t := range manifest.Stack.Transport {
		if t == transport {
			found = true
			break
		}
	}
	if !found {
		manifest.Stack.Transport = append(manifest.Stack.Transport, transport)
	}
	
	if !dryRun {
		return s.resolver.Save(ctx, manifestPath, manifest, s.fs)
	}
	return nil
}
