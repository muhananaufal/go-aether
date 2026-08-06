package port

import (
	"context"
	"io"

	"github.com/muhananaufal/go-aether/internal/core/domain"
)

// TemplateEngine abstracts the rendering process of embedded or customized text/templates.
type TemplateEngine interface {
	// Render executes the specified template file from embed.FS using the supplied TemplateData.
	Render(ctx context.Context, templatePath string, data *domain.TemplateData) ([]byte, error)

	// ListTemplates returns a comprehensive list of available template filenames within a pattern directory.
	ListTemplates(pattern string) ([]string, error)
}

// FileWriter abstracts filesystem operations to support real disk IO, in-memory testing (afero),
// dry-run diff display, and atomic transactional write buffers.
type FileWriter interface {
	// WriteFile writes content to targetAbsolutePath, applying backup or skip behavior based on flags.
	WriteFile(ctx context.Context, targetPath string, content []byte, overwrite, dryRun bool) error

	// Exists returns true if a target file or directory currently exists on the filesystem.
	Exists(path string) (bool, error)

	// ReadFile retrieves the binary contents of an existing file on the filesystem.
	ReadFile(path string) ([]byte, error)

	// DeleteFile removes a file from disk during transactional atomic rollback procedures.
	DeleteFile(path string) error

	// MkdirAll ensures the full directory tree exists for a target path before writing.
	MkdirAll(dirPath string) error
}

// ManifestResolver abstracts locating, parsing, and persisting the aether.yaml manifest file.
type ManifestResolver interface {
	// Resolve locates aether.yaml by walking up the directory tree from current working directory.
	Resolve(ctx context.Context, startDir string) (*domain.AetherManifest, string, error)

	// Save writes the validated AetherManifest out to the specified file system destination.
	Save(ctx context.Context, destPath string, manifest *domain.AetherManifest, writer FileWriter) error
}

// ScaffoldService defines the core application orchestration contract for project CLI operations.
type ScaffoldService interface {
	// InitProject bootstraps a new greenfield Go project architecture in the specified directory.
	InitProject(ctx context.Context, destDir, projectName, moduleName, arch, dbDriver, router string, dryRun bool) error

	// AdoptProject performs interactive anomaly mapping to integrate go-aether into a legacy repository.
	AdoptProject(ctx context.Context, destDir string, scan bool, dryRun bool) error

	// MakeModule generates all vertical slice components (domain, port, service, handler, repo) for a feature.
	MakeModule(ctx context.Context, startDir, moduleName string, transports []string, hasCache, hasWorker, dryRun, force bool) error

	// RunDoctor audits project architecture health and checks consistency against aether.yaml.
	RunDoctor(ctx context.Context, startDir string, fix bool, out io.Writer) error

	// MakeService generates only the service layer component for a specific module.
	MakeService(ctx context.Context, startDir, moduleName string, dryRun, force bool) error

	// MakeHandler generates only the transport handler component for a specific module.
	MakeHandler(ctx context.Context, startDir, moduleName, transport string, dryRun, force bool) error

	// MakeDomain generates only the domain layer entity for a specific module.
	MakeDomain(ctx context.Context, startDir, moduleName string, dryRun, force bool) error

	// MakePort generates only the port interface contract for a specific module.
	MakePort(ctx context.Context, startDir, moduleName string, dryRun, force bool) error

	// MakeRepository generates only the infrastructure repository for a specific module.
	MakeRepository(ctx context.Context, startDir, moduleName string, dryRun, force bool) error

	// AddMiddleware injects middleware components into a target module's transport handler.
	AddMiddleware(ctx context.Context, startDir, moduleName, middlewareType string, dryRun, force bool) error

	// AddCache sets up the global cache layer configuration and generates the cache provider infrastructure.
	AddCache(ctx context.Context, startDir, cacheType string, dryRun, force bool) error

	// AddTransport registers a new global transport protocol.
	AddTransport(ctx context.Context, startDir, transport string, dryRun, force bool) error

	// AddWorker generates an asynchronous background processor in the worker directory.
	AddWorker(ctx context.Context, startDir, workerName, broker, pattern string, dryRun, force bool) error

	// AddEventing sets up the global Publisher and Subscriber interfaces for event-driven architecture.
	AddEventing(ctx context.Context, startDir, broker string, dryRun, force bool) error

	// AddMetrics sets up the Prometheus metrics middleware and endpoint.
	AddMetrics(ctx context.Context, startDir, provider string, dryRun, force bool) error

	// AddTracing sets up the OpenTelemetry tracing infrastructure.
	AddTracing(ctx context.Context, startDir, exporter string, dryRun, force bool) error

	// AddDeploy sets up the deployment manifests (e.g. Kubernetes, Helm).
	AddDeploy(ctx context.Context, startDir, target string, dryRun, force bool) error

	// AddCICD sets up the CI/CD pipelines (e.g. GitHub Actions).
	AddCICD(ctx context.Context, startDir, provider string, dryRun, force bool) error

	// AddAI sets up the LLM proxy interface and stub.
	AddAI(ctx context.Context, startDir, provider string, dryRun, force bool) error

	// AddDI sets up a dependency injection container (e.g. fx, wire).
	AddDI(ctx context.Context, startDir, diType string, dryRun, force bool) error

	// AddConfig sets up a centralized configuration manager (e.g. viper, koanf).
	AddConfig(ctx context.Context, startDir, configType string, dryRun, force bool) error

	// AddError sets up a standardized centralized error handler.
	AddError(ctx context.Context, startDir string, dryRun, force bool) error

	// AddValidator sets up the struct validation wrapper (e.g. go-playground/validator).
	AddValidator(ctx context.Context, startDir, validatorType string, dryRun, force bool) error
}
