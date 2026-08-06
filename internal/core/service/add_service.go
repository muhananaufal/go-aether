package service

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

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

// AddWorker generates an asynchronous background processor in the worker directory.
func (s *AetherScaffoldService) AddWorker(ctx context.Context, startDir, workerName, broker, pattern string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData(workerName, manifest, nil, false, false)
	if err != nil {
		return err
	}

	tmplName := fmt.Sprintf("distributed/worker_%s.go.tmpl", pattern)
	content, err := s.engine.Render(ctx, tmplName, data)
	if err != nil {
		return err
	}

	found := false
	for _, w := range manifest.Workers {
		if w.Name == workerName {
			found = true
			break
		}
	}
	if !found {
		manifest.Workers = append(manifest.Workers, domain.WorkerRegistry{
			Name:    workerName,
			Broker:  broker,
			Pattern: pattern,
		})
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	destFile := filepath.Join(projectRoot, manifest.Architecture.Paths.Pkg, "worker", fmt.Sprintf("%s_worker.go", strings.ToLower(workerName)))
	
	tx.Stage(destFile, content, force, dryRun)
	
	if err := tx.Commit(ctx); err != nil {
		return err
	}

	if !dryRun {
		return s.resolver.Save(ctx, manifestPath, manifest, s.fs)
	}
	return nil
}

// AddEventing sets up the global Publisher and Subscriber interfaces for event-driven architecture.
func (s *AetherScaffoldService) AddEventing(ctx context.Context, startDir, broker string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, _ := domain.NewTemplateData("event", manifest, nil, false, false)
	tmplName := "distributed/event_bus.go.tmpl"
	content, err := s.engine.Render(ctx, tmplName, data)
	if err != nil {
		return err
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	destFile := filepath.Join(projectRoot, manifest.Architecture.Paths.Pkg, "event", "event_bus.go")
	
	tx.Stage(destFile, content, force, dryRun)
	
	return tx.Commit(ctx)
}

// AddMetrics sets up the Prometheus metrics middleware and endpoint.
func (s *AetherScaffoldService) AddMetrics(ctx context.Context, startDir, provider string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, _ := domain.NewTemplateData("middleware", manifest, nil, false, false)
	
	tmplName := "common/metrics_prometheus.go.tmpl"
	content, err := s.engine.Render(ctx, tmplName, data)
	if err != nil {
		return err
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	destFile := filepath.Join(projectRoot, manifest.Architecture.Paths.Pkg, "middleware", "metrics_prometheus.go")
	
	tx.Stage(destFile, content, force, dryRun)
	
	return tx.Commit(ctx)
}

// AddTracing sets up the OpenTelemetry tracing infrastructure.
func (s *AetherScaffoldService) AddTracing(ctx context.Context, startDir, exporter string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, _ := domain.NewTemplateData("telemetry", manifest, nil, false, false)
	tmplName := "common/tracing_otel.go.tmpl"
	content, err := s.engine.Render(ctx, tmplName, data)
	if err != nil {
		return err
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	destFile := filepath.Join(projectRoot, manifest.Architecture.Paths.Pkg, "telemetry", "tracing_otel.go")
	
	tx.Stage(destFile, content, force, dryRun)
	
	return tx.Commit(ctx)
}

// AddDeploy sets up the deployment manifests (e.g. Kubernetes, Helm).
func (s *AetherScaffoldService) AddDeploy(ctx context.Context, startDir, target string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, _ := domain.NewTemplateData("deploy", manifest, nil, false, false)
	tmplName := fmt.Sprintf("cloud/%s_deployment.yaml.tmpl", target)
	content, err := s.engine.Render(ctx, tmplName, data)
	if err != nil {
		return err
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	destFile := filepath.Join(projectRoot, "deploy", fmt.Sprintf("%s.yaml", target))
	
	tx.Stage(destFile, content, force, dryRun)
	return tx.Commit(ctx)
}

// AddCICD sets up the CI/CD pipelines (e.g. GitHub Actions).
func (s *AetherScaffoldService) AddCICD(ctx context.Context, startDir, provider string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, _ := domain.NewTemplateData("cicd", manifest, nil, false, false)
	var tmplName, destFile string

	projectRoot := filepath.Dir(manifestPath)
	if provider == "github" {
		tmplName = "cloud/github_actions.yml.tmpl"
		destFile = filepath.Join(projectRoot, ".github", "workflows", "ci.yml")
	} else {
		return fmt.Errorf("unsupported CI/CD provider: %s", provider)
	}

	content, err := s.engine.Render(ctx, tmplName, data)
	if err != nil {
		return err
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	tx.Stage(destFile, content, force, dryRun)
	return tx.Commit(ctx)
}

// AddAI sets up the LLM proxy interface and stub.
func (s *AetherScaffoldService) AddAI(ctx context.Context, startDir, provider string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, _ := domain.NewTemplateData("ai", manifest, nil, false, false)
	tmplName := "ai/llm_proxy.go.tmpl"
	content, err := s.engine.Render(ctx, tmplName, data)
	if err != nil {
		return err
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	destFile := filepath.Join(projectRoot, manifest.Architecture.Paths.Pkg, "ai", "llm_proxy.go")
	
	tx.Stage(destFile, content, force, dryRun)
	return tx.Commit(ctx)
}

// AddDI sets up a dependency injection container (e.g. fx, wire).
func (s *AetherScaffoldService) AddDI(ctx context.Context, startDir, diType string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, _ := domain.NewTemplateData("di", manifest, nil, false, false)
	tmplName := fmt.Sprintf("common/di_%s.go.tmpl", diType)
	content, err := s.engine.Render(ctx, tmplName, data)
	if err != nil {
		return err
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	destFile := filepath.Join(projectRoot, manifest.Architecture.Paths.Pkg, "config", "di.go")
	
	tx.Stage(destFile, content, force, dryRun)
	return tx.Commit(ctx)
}

// AddConfig sets up a centralized configuration manager (e.g. viper, koanf).
func (s *AetherScaffoldService) AddConfig(ctx context.Context, startDir, configType string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, _ := domain.NewTemplateData("config", manifest, nil, false, false)
	tmplName := fmt.Sprintf("common/config_%s.go.tmpl", configType)
	content, err := s.engine.Render(ctx, tmplName, data)
	if err != nil {
		return err
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	destFile := filepath.Join(projectRoot, manifest.Architecture.Paths.Pkg, "config", "config.go")
	
	tx.Stage(destFile, content, force, dryRun)
	return tx.Commit(ctx)
}

// AddError sets up a standardized centralized error handler.
func (s *AetherScaffoldService) AddError(ctx context.Context, startDir string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, _ := domain.NewTemplateData("error", manifest, nil, false, false)
	tmplName := "common/error_handler.go.tmpl"
	content, err := s.engine.Render(ctx, tmplName, data)
	if err != nil {
		return err
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	destFile := filepath.Join(projectRoot, manifest.Architecture.Paths.Pkg, "common", "error.go")
	
	tx.Stage(destFile, content, force, dryRun)
	return tx.Commit(ctx)
}

// AddValidator sets up the struct validation wrapper.
func (s *AetherScaffoldService) AddValidator(ctx context.Context, startDir, validatorType string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, _ := domain.NewTemplateData("validator", manifest, nil, false, false)
	tmplName := fmt.Sprintf("common/validator_%s.go.tmpl", validatorType)
	content, err := s.engine.Render(ctx, tmplName, data)
	if err != nil {
		return err
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	destFile := filepath.Join(projectRoot, manifest.Architecture.Paths.Pkg, "common", "validator.go")
	
	tx.Stage(destFile, content, force, dryRun)
	tx.Stage(destFile, content, force, dryRun)
	return tx.Commit(ctx)
}

// AddTest scaffolds standard integration test setup and helpers.
func (s *AetherScaffoldService) AddTest(ctx context.Context, startDir string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData("tests", manifest, nil, false, false)
	if err != nil {
		return err
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	
	content, err := s.engine.Render(ctx, "common/test_setup.go.tmpl", data)
	if err != nil {
		return err
	}

	destFile := filepath.Join(projectRoot, "tests", "setup.go")
	tx.Stage(destFile, content, force, dryRun)

	return tx.Commit(ctx)
}

// AddMultitenancy scaffolds Row Level Security (RLS) SQL policies for tenant isolation.
func (s *AetherScaffoldService) AddMultitenancy(ctx context.Context, startDir, moduleName string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData(moduleName, manifest, nil, false, false)
	if err != nil {
		return err
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	
	content, err := s.engine.Render(ctx, "common/multitenancy_rls.sql.tmpl", data)
	if err != nil {
		return err
	}

	destFile := filepath.Join(projectRoot, "migrations", fmt.Sprintf("rls_%s.sql", strings.ToLower(moduleName)))
	tx.Stage(destFile, content, force, dryRun)

	return tx.Commit(ctx)
}

// AddCQRS scaffolds Command and Query handlers separated within module scope.
func (s *AetherScaffoldService) AddCQRS(ctx context.Context, startDir, moduleName string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData(moduleName, manifest, nil, false, false)
	if err != nil {
		return err
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	moduleDir := filepath.Join(projectRoot, "internal", "modules", strings.ToLower(moduleName), "cqrs")

	// Render command.go
	cmdContent, err := s.engine.Render(ctx, "distributed/cqrs_command.go.tmpl", data)
	if err != nil {
		return err
	}
	tx.Stage(filepath.Join(moduleDir, "command.go"), cmdContent, force, dryRun)

	// Render query.go
	queryContent, err := s.engine.Render(ctx, "distributed/cqrs_query.go.tmpl", data)
	if err != nil {
		return err
	}
	tx.Stage(filepath.Join(moduleDir, "query.go"), queryContent, force, dryRun)

	// Render bus.go
	busContent, err := s.engine.Render(ctx, "distributed/cqrs_bus.go.tmpl", data)
	if err != nil {
		return err
	}
	tx.Stage(filepath.Join(moduleDir, "bus.go"), busContent, force, dryRun)

	return tx.Commit(ctx)
}

// AddOutbox scaffolds Transactional Outbox pattern infrastructure and SQL migrations.
func (s *AetherScaffoldService) AddOutbox(ctx context.Context, startDir string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData("outbox", manifest, nil, false, false)
	if err != nil {
		return err
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)

	// Render pkg/outbox/outbox.go
	outboxContent, err := s.engine.Render(ctx, "distributed/outbox.go.tmpl", data)
	if err != nil {
		return err
	}
	tx.Stage(filepath.Join(projectRoot, manifest.Architecture.Paths.Pkg, "outbox", "outbox.go"), outboxContent, force, dryRun)

	// Render migrations/timestamp_create_outbox.up.sql & down.sql
	timestamp := time.Now().Format("20060102150405")
	upContent, err := s.engine.Render(ctx, "distributed/outbox_up.sql.tmpl", data)
	if err != nil {
		return err
	}
	tx.Stage(filepath.Join(projectRoot, "migrations", fmt.Sprintf("%s_create_outbox_events.up.sql", timestamp)), upContent, force, dryRun)

	downContent, err := s.engine.Render(ctx, "distributed/outbox_down.sql.tmpl", data)
	if err != nil {
		return err
	}
	tx.Stage(filepath.Join(projectRoot, "migrations", fmt.Sprintf("%s_create_outbox_events.down.sql", timestamp)), downContent, force, dryRun)

	return tx.Commit(ctx)
}

// AddSaga scaffolds Distributed Saga workflow engine and step orchestrators.
func (s *AetherScaffoldService) AddSaga(ctx context.Context, startDir, workflowName string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData(workflowName, manifest, nil, false, false)
	if err != nil {
		return err
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)

	// Render pkg/saga/saga.go (Core engine)
	engineContent, err := s.engine.Render(ctx, "distributed/saga_engine.go.tmpl", data)
	if err != nil {
		return err
	}
	tx.Stage(filepath.Join(projectRoot, manifest.Architecture.Paths.Pkg, "saga", "saga.go"), engineContent, force, dryRun)

	// Render internal/workflows/<workflow>_saga.go
	workflowContent, err := s.engine.Render(ctx, "distributed/saga_workflow.go.tmpl", data)
	if err != nil {
		return err
	}
	tx.Stage(filepath.Join(projectRoot, "internal", "workflows", fmt.Sprintf("%s_saga.go", strings.ToLower(workflowName))), workflowContent, force, dryRun)

	return tx.Commit(ctx)
}

// AddWebhook scaffolds secure signed webhook dispatcher and signature verifier middleware.
func (s *AetherScaffoldService) AddWebhook(ctx context.Context, startDir string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData("webhook", manifest, nil, false, false)
	if err != nil {
		return err
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	webhookDir := filepath.Join(projectRoot, manifest.Architecture.Paths.Pkg, "webhook")

	// Render pkg/webhook/dispatcher.go
	dispContent, err := s.engine.Render(ctx, "distributed/webhook_dispatcher.go.tmpl", data)
	if err != nil {
		return err
	}
	tx.Stage(filepath.Join(webhookDir, "dispatcher.go"), dispContent, force, dryRun)

	// Render pkg/webhook/receiver.go
	recvContent, err := s.engine.Render(ctx, "distributed/webhook_receiver.go.tmpl", data)
	if err != nil {
		return err
	}
	tx.Stage(filepath.Join(webhookDir, "receiver.go"), recvContent, force, dryRun)

	return tx.Commit(ctx)
}

// AddDiscovery scaffolds service discovery registration client (Consul or etcd).
func (s *AetherScaffoldService) AddDiscovery(ctx context.Context, startDir, provider string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	providerLower := strings.ToLower(provider)
	if providerLower != "consul" && providerLower != "etcd" {
		return fmt.Errorf("unsupported discovery provider '%s', must be 'consul' or 'etcd'", provider)
	}

	data, err := domain.NewTemplateData(providerLower, manifest, nil, false, false)
	if err != nil {
		return err
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	discoveryDir := filepath.Join(projectRoot, manifest.Architecture.Paths.Pkg, "discovery")

	templateName := fmt.Sprintf("distributed/discovery_%s.go.tmpl", providerLower)
	content, err := s.engine.Render(ctx, templateName, data)
	if err != nil {
		return err
	}

	destFile := filepath.Join(discoveryDir, fmt.Sprintf("%s.go", providerLower))
	tx.Stage(destFile, content, force, dryRun)

	return tx.Commit(ctx)
}

// AddAuth scaffolds authentication handlers and middleware (OAuth2 or APIKey).
func (s *AetherScaffoldService) AddAuth(ctx context.Context, startDir, authType string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	authTypeLower := strings.ToLower(authType)
	if authTypeLower != "oauth2" && authTypeLower != "apikey" {
		return fmt.Errorf("unsupported auth type '%s', must be 'oauth2' or 'apikey'", authType)
	}

	data, err := domain.NewTemplateData(authTypeLower, manifest, nil, false, false)
	if err != nil {
		return err
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	authDir := filepath.Join(projectRoot, manifest.Architecture.Paths.Pkg, "auth")

	templateName := fmt.Sprintf("plugins/auth_%s.go.tmpl", authTypeLower)
	content, err := s.engine.Render(ctx, templateName, data)
	if err != nil {
		return err
	}

	destFile := filepath.Join(authDir, fmt.Sprintf("%s.go", authTypeLower))
	tx.Stage(destFile, content, force, dryRun)

	return tx.Commit(ctx)
}

// AddStorage scaffolds cloud blob storage interface and S3/MinIO adapter.
func (s *AetherScaffoldService) AddStorage(ctx context.Context, startDir, provider string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	providerLower := strings.ToLower(provider)
	if providerLower != "s3" && providerLower != "gcs" && providerLower != "local" {
		providerLower = "s3"
	}

	data, err := domain.NewTemplateData(providerLower, manifest, nil, false, false)
	if err != nil {
		return err
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	storageDir := filepath.Join(projectRoot, manifest.Architecture.Paths.Pkg, "storage")

	// Render pkg/storage/storage.go
	ifaceContent, err := s.engine.Render(ctx, "plugins/storage_interface.go.tmpl", data)
	if err != nil {
		return err
	}
	tx.Stage(filepath.Join(storageDir, "storage.go"), ifaceContent, force, dryRun)

	// Render pkg/storage/s3.go
	s3Content, err := s.engine.Render(ctx, "plugins/storage_s3.go.tmpl", data)
	if err != nil {
		return err
	}
	tx.Stage(filepath.Join(storageDir, fmt.Sprintf("%s.go", providerLower)), s3Content, force, dryRun)

	return tx.Commit(ctx)
}

// AddCron scaffolds in-process recurring background job scheduler.
func (s *AetherScaffoldService) AddCron(ctx context.Context, startDir, jobName string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData(jobName, manifest, nil, false, false)
	if err != nil {
		return err
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	cronDir := filepath.Join(projectRoot, manifest.Architecture.Paths.Pkg, "cron")
	jobsDir := filepath.Join(projectRoot, "internal", "jobs")

	// Render pkg/cron/scheduler.go
	schedContent, err := s.engine.Render(ctx, "plugins/cron_scheduler.go.tmpl", data)
	if err != nil {
		return err
	}
	tx.Stage(filepath.Join(cronDir, "scheduler.go"), schedContent, force, dryRun)

	// Render internal/jobs/<job>_job.go
	jobContent, err := s.engine.Render(ctx, "plugins/cron_job.go.tmpl", data)
	if err != nil {
		return err
	}
	tx.Stage(filepath.Join(jobsDir, fmt.Sprintf("%s_job.go", data.ModuleNamePkg)), jobContent, force, dryRun)

	return tx.Commit(ctx)
}

// AddMailer scaffolds transactional email delivery client.
func (s *AetherScaffoldService) AddMailer(ctx context.Context, startDir, provider string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData(provider, manifest, nil, false, false)
	if err != nil {
		return err
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	mailerDir := filepath.Join(projectRoot, manifest.Architecture.Paths.Pkg, "mailer")

	content, err := s.engine.Render(ctx, "plugins/mailer.go.tmpl", data)
	if err != nil {
		return err
	}
	tx.Stage(filepath.Join(mailerDir, "mailer.go"), content, force, dryRun)

	return tx.Commit(ctx)
}

// AddFirebase scaffolds Firebase Auth token verifier and FCM push messaging client.
func (s *AetherScaffoldService) AddFirebase(ctx context.Context, startDir string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData("firebase", manifest, nil, false, false)
	if err != nil {
		return err
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	fbDir := filepath.Join(projectRoot, manifest.Architecture.Paths.Pkg, "firebase")

	content, err := s.engine.Render(ctx, "plugins/firebase.go.tmpl", data)
	if err != nil {
		return err
	}
	tx.Stage(filepath.Join(fbDir, "firebase.go"), content, force, dryRun)

	return tx.Commit(ctx)
}

// AddLogger scaffolds structured JSON logger with context correlation tracking.
func (s *AetherScaffoldService) AddLogger(ctx context.Context, startDir, provider string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData(provider, manifest, nil, false, false)
	if err != nil {
		return err
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	loggerDir := filepath.Join(projectRoot, manifest.Architecture.Paths.Pkg, "logger")

	content, err := s.engine.Render(ctx, "plugins/logger.go.tmpl", data)
	if err != nil {
		return err
	}
	tx.Stage(filepath.Join(loggerDir, "logger.go"), content, force, dryRun)

	return tx.Commit(ctx)
}

// AddHealthcheck scaffolds Kubernetes-compatible Liveness & Readiness probe handlers.
func (s *AetherScaffoldService) AddHealthcheck(ctx context.Context, startDir string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData("health", manifest, nil, false, false)
	if err != nil {
		return err
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	healthDir := filepath.Join(projectRoot, manifest.Architecture.Paths.Pkg, "health")

	content, err := s.engine.Render(ctx, "plugins/healthcheck.go.tmpl", data)
	if err != nil {
		return err
	}
	tx.Stage(filepath.Join(healthDir, "healthcheck.go"), content, force, dryRun)

	return tx.Commit(ctx)
}

// AddSecrets scaffolds secret manager client (HashiCorp Vault or AWS Secrets Manager).
func (s *AetherScaffoldService) AddSecrets(ctx context.Context, startDir, provider string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	providerLower := strings.ToLower(provider)
	if providerLower != "vault" && providerLower != "aws" {
		providerLower = "vault"
	}

	data, err := domain.NewTemplateData(providerLower, manifest, nil, false, false)
	if err != nil {
		return err
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	secretsDir := filepath.Join(projectRoot, manifest.Architecture.Paths.Pkg, "secrets")

	templateName := fmt.Sprintf("plugins/secrets_%s.go.tmpl", providerLower)
	content, err := s.engine.Render(ctx, templateName, data)
	if err != nil {
		return err
	}
	destFile := filepath.Join(secretsDir, fmt.Sprintf("%s.go", providerLower))
	tx.Stage(destFile, content, force, dryRun)

	return tx.Commit(ctx)
}

// AddLock scaffolds distributed lock (Redlock / Distributed Mutex) engine.
func (s *AetherScaffoldService) AddLock(ctx context.Context, startDir, provider string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	providerLower := strings.ToLower(provider)
	if providerLower != "redis" {
		providerLower = "redis"
	}

	data, err := domain.NewTemplateData(providerLower, manifest, nil, false, false)
	if err != nil {
		return err
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	lockDir := filepath.Join(projectRoot, manifest.Architecture.Paths.Pkg, "lock")

	content, err := s.engine.Render(ctx, "plugins/lock_redis.go.tmpl", data)
	if err != nil {
		return err
	}
	tx.Stage(filepath.Join(lockDir, "lock.go"), content, force, dryRun)

	return tx.Commit(ctx)
}

// AddAuthz scaffolds RBAC / ABAC authorization engine and middleware (Casbin).
func (s *AetherScaffoldService) AddAuthz(ctx context.Context, startDir, provider string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData("authz", manifest, nil, false, false)
	if err != nil {
		return err
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	authzDir := filepath.Join(projectRoot, manifest.Architecture.Paths.Pkg, "authz")

	content, err := s.engine.Render(ctx, "plugins/authz_casbin.go.tmpl", data)
	if err != nil {
		return err
	}
	tx.Stage(filepath.Join(authzDir, "casbin.go"), content, force, dryRun)

	return tx.Commit(ctx)
}

// AddCrypto scaffolds symmetric envelope encryption helper (AES-GCM).
func (s *AetherScaffoldService) AddCrypto(ctx context.Context, startDir, algorithm string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData("encryption", manifest, nil, false, false)
	if err != nil {
		return err
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	cryptoDir := filepath.Join(projectRoot, manifest.Architecture.Paths.Pkg, "crypto")

	content, err := s.engine.Render(ctx, "plugins/crypto_aes.go.tmpl", data)
	if err != nil {
		return err
	}
	tx.Stage(filepath.Join(cryptoDir, "aes.go"), content, force, dryRun)

	return tx.Commit(ctx)
}

// AddProfiling scaffolds protected runtime pprof / Pyroscope profiling endpoints.
func (s *AetherScaffoldService) AddProfiling(ctx context.Context, startDir, provider string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData("profiling", manifest, nil, false, false)
	if err != nil {
		return err
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	profDir := filepath.Join(projectRoot, manifest.Architecture.Paths.Pkg, "profiling")

	content, err := s.engine.Render(ctx, "plugins/profiling_pprof.go.tmpl", data)
	if err != nil {
		return err
	}
	tx.Stage(filepath.Join(profDir, "pprof.go"), content, force, dryRun)

	return tx.Commit(ctx)
}

// AddFeatureFlags scaffolds feature flag & canary release client (Flipt).
func (s *AetherScaffoldService) AddFeatureFlags(ctx context.Context, startDir, provider string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData("featureflags", manifest, nil, false, false)
	if err != nil {
		return err
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	ffDir := filepath.Join(projectRoot, manifest.Architecture.Paths.Pkg, "featureflags")

	content, err := s.engine.Render(ctx, "plugins/featureflags_flipt.go.tmpl", data)
	if err != nil {
		return err
	}
	tx.Stage(filepath.Join(ffDir, "flipt.go"), content, force, dryRun)

	return tx.Commit(ctx)
}

// AddIdempotency scaffolds Idempotency-Key validation middleware with distributed lock.
func (s *AetherScaffoldService) AddIdempotency(ctx context.Context, startDir, provider string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData("idempotency", manifest, nil, false, false)
	if err != nil {
		return err
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	destDir := filepath.Join(projectRoot, manifest.Architecture.Paths.Pkg, "idempotency")

	content, err := s.engine.Render(ctx, "plugins/idempotency.go.tmpl", data)
	if err != nil {
		return err
	}
	tx.Stage(filepath.Join(destDir, "idempotency.go"), content, force, dryRun)

	return tx.Commit(ctx)
}

// AddLedger scaffolds Double-Entry Bookkeeping Ledger engine.
func (s *AetherScaffoldService) AddLedger(ctx context.Context, startDir string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData("ledger", manifest, nil, false, false)
	if err != nil {
		return err
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	destDir := filepath.Join(projectRoot, manifest.Architecture.Paths.Pkg, "ledger")

	content, err := s.engine.Render(ctx, "plugins/ledger.go.tmpl", data)
	if err != nil {
		return err
	}
	tx.Stage(filepath.Join(destDir, "ledger.go"), content, force, dryRun)

	return tx.Commit(ctx)
}

// AddDecimal scaffolds high-precision financial decimal money arithmetic helpers.
func (s *AetherScaffoldService) AddDecimal(ctx context.Context, startDir string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData("money", manifest, nil, false, false)
	if err != nil {
		return err
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	destDir := filepath.Join(projectRoot, manifest.Architecture.Paths.Pkg, "money")

	content, err := s.engine.Render(ctx, "plugins/decimal.go.tmpl", data)
	if err != nil {
		return err
	}
	tx.Stage(filepath.Join(destDir, "money.go"), content, force, dryRun)

	return tx.Commit(ctx)
}

// AddReconciliation scaffolds automated transaction settlement reconciliation matching engine.
func (s *AetherScaffoldService) AddReconciliation(ctx context.Context, startDir string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData("reconciliation", manifest, nil, false, false)
	if err != nil {
		return err
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	destDir := filepath.Join(projectRoot, manifest.Architecture.Paths.Pkg, "reconciliation")

	content, err := s.engine.Render(ctx, "plugins/reconciliation.go.tmpl", data)
	if err != nil {
		return err
	}
	tx.Stage(filepath.Join(destDir, "reconciliation.go"), content, force, dryRun)

	return tx.Commit(ctx)
}

// AddPricingEngine scaffolds rule-based tiered pricing and fee calculation engine.
func (s *AetherScaffoldService) AddPricingEngine(ctx context.Context, startDir string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData("pricing", manifest, nil, false, false)
	if err != nil {
		return err
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	destDir := filepath.Join(projectRoot, manifest.Architecture.Paths.Pkg, "pricing")

	content, err := s.engine.Render(ctx, "plugins/pricing_engine.go.tmpl", data)
	if err != nil {
		return err
	}
	tx.Stage(filepath.Join(destDir, "pricing.go"), content, force, dryRun)

	return tx.Commit(ctx)
}

// AddWebSocket scaffolds WebSocket hub and multi-client connection pool (Gorilla / Nhooyr).
func (s *AetherScaffoldService) AddWebSocket(ctx context.Context, startDir, provider string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData("websocket", manifest, nil, false, false)
	if err != nil {
		return err
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	destDir := filepath.Join(projectRoot, manifest.Architecture.Paths.Pkg, "websocket")

	content, err := s.engine.Render(ctx, "plugins/websocket_gorilla.go.tmpl", data)
	if err != nil {
		return err
	}
	tx.Stage(filepath.Join(destDir, "websocket.go"), content, force, dryRun)

	return tx.Commit(ctx)
}

// AddSSE scaffolds Server-Sent Events unidirectional streaming broker.
func (s *AetherScaffoldService) AddSSE(ctx context.Context, startDir string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData("sse", manifest, nil, false, false)
	if err != nil {
		return err
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	destDir := filepath.Join(projectRoot, manifest.Architecture.Paths.Pkg, "sse")

	content, err := s.engine.Render(ctx, "plugins/sse.go.tmpl", data)
	if err != nil {
		return err
	}
	tx.Stage(filepath.Join(destDir, "sse.go"), content, force, dryRun)

	return tx.Commit(ctx)
}

// AddWebRTC scaffolds Pion WebRTC peer connection and data channel signaling hub.
func (s *AetherScaffoldService) AddWebRTC(ctx context.Context, startDir, provider string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData("webrtc", manifest, nil, false, false)
	if err != nil {
		return err
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	destDir := filepath.Join(projectRoot, manifest.Architecture.Paths.Pkg, "webrtc")

	content, err := s.engine.Render(ctx, "plugins/webrtc_pion.go.tmpl", data)
	if err != nil {
		return err
	}
	tx.Stage(filepath.Join(destDir, "webrtc.go"), content, force, dryRun)

	return tx.Commit(ctx)
}

// AddMQTT scaffolds Paho MQTT client for IoT telemetry pub/sub messaging.
func (s *AetherScaffoldService) AddMQTT(ctx context.Context, startDir, provider string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData("mqtt", manifest, nil, false, false)
	if err != nil {
		return err
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	destDir := filepath.Join(projectRoot, manifest.Architecture.Paths.Pkg, "mqtt")

	content, err := s.engine.Render(ctx, "plugins/mqtt_paho.go.tmpl", data)
	if err != nil {
		return err
	}
	tx.Stage(filepath.Join(destDir, "mqtt.go"), content, force, dryRun)

	return tx.Commit(ctx)
}

// AddTwilio scaffolds Twilio SMS and WhatsApp omni-channel delivery client.
func (s *AetherScaffoldService) AddTwilio(ctx context.Context, startDir string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData("twilio", manifest, nil, false, false)
	if err != nil {
		return err
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	destDir := filepath.Join(projectRoot, manifest.Architecture.Paths.Pkg, "twilio")

	content, err := s.engine.Render(ctx, "plugins/twilio.go.tmpl", data)
	if err != nil {
		return err
	}
	tx.Stage(filepath.Join(destDir, "twilio.go"), content, force, dryRun)

	return tx.Commit(ctx)
}

// AddMultiLevelCache scaffolds synchronized L1 memory + L2 Redis cache.
func (s *AetherScaffoldService) AddMultiLevelCache(ctx context.Context, startDir string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData("multilevelcache", manifest, nil, false, false)
	if err != nil {
		return err
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	destDir := filepath.Join(projectRoot, manifest.Architecture.Paths.Pkg, "cache")

	content, err := s.engine.Render(ctx, "plugins/multilevelcache.go.tmpl", data)
	if err != nil {
		return err
	}
	tx.Stage(filepath.Join(destDir, "multilevel.go"), content, force, dryRun)

	return tx.Commit(ctx)
}

// AddBloomFilter scaffolds probabilistic cache penetration guard.
func (s *AetherScaffoldService) AddBloomFilter(ctx context.Context, startDir string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData("bloom", manifest, nil, false, false)
	if err != nil {
		return err
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	destDir := filepath.Join(projectRoot, manifest.Architecture.Paths.Pkg, "bloom")

	content, err := s.engine.Render(ctx, "plugins/bloomfilter.go.tmpl", data)
	if err != nil {
		return err
	}
	tx.Stage(filepath.Join(destDir, "bloom.go"), content, force, dryRun)

	return tx.Commit(ctx)
}

// AddS3 scaffolds S3 object storage client with pre-signed URL generator (MinIO / AWS).
func (s *AetherScaffoldService) AddS3(ctx context.Context, startDir, provider string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData("s3", manifest, nil, false, false)
	if err != nil {
		return err
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	destDir := filepath.Join(projectRoot, manifest.Architecture.Paths.Pkg, "storage")

	content, err := s.engine.Render(ctx, "plugins/s3_storage.go.tmpl", data)
	if err != nil {
		return err
	}
	tx.Stage(filepath.Join(destDir, "s3.go"), content, force, dryRun)

	return tx.Commit(ctx)
}

// AddResilience scaffolds Circuit Breaker and Bulkhead resilience engine.
func (s *AetherScaffoldService) AddResilience(ctx context.Context, startDir, provider string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData("resilience", manifest, nil, false, false)
	if err != nil {
		return err
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	destDir := filepath.Join(projectRoot, manifest.Architecture.Paths.Pkg, "resilience")

	content, err := s.engine.Render(ctx, "plugins/resilience.go.tmpl", data)
	if err != nil {
		return err
	}
	tx.Stage(filepath.Join(destDir, "resilience.go"), content, force, dryRun)

	return tx.Commit(ctx)
}

// AddSearch scaffolds full-text typo-tolerant search engine client (Meilisearch / Elasticsearch).
func (s *AetherScaffoldService) AddSearch(ctx context.Context, startDir, provider string, dryRun, force bool) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		return fmt.Errorf("failed to locate aether.yaml: %w", err)
	}

	data, err := domain.NewTemplateData("search", manifest, nil, false, false)
	if err != nil {
		return err
	}

	tx := writer.NewTransactionalBuffer(s.fs)
	projectRoot := filepath.Dir(manifestPath)
	destDir := filepath.Join(projectRoot, manifest.Architecture.Paths.Pkg, "search")

	content, err := s.engine.Render(ctx, "plugins/search_meilisearch.go.tmpl", data)
	if err != nil {
		return err
	}
	tx.Stage(filepath.Join(destDir, "search.go"), content, force, dryRun)

	return tx.Commit(ctx)
}
