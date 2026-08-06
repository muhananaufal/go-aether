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

	// MakeMigration generates a SQL migration file pair (up/down).
	MakeMigration(ctx context.Context, startDir, name string, dryRun, force bool) error

	// MakeSeeder generates a database seeder file.
	MakeSeeder(ctx context.Context, startDir, name string, dryRun, force bool) error

	// MakeValueObject generates an immutable DDD Value Object struct.
	MakeValueObject(ctx context.Context, startDir, name string, dryRun, force bool) error

	// MakeAggregate generates a DDD Aggregate Root entity with event recording.
	MakeAggregate(ctx context.Context, startDir, name string, dryRun, force bool) error

	// MakeEvent generates a Domain Event struct and serializer.
	MakeEvent(ctx context.Context, startDir, name string, dryRun, force bool) error

	// MakeCommand generates a CQRS Command DTO and execution handler.
	MakeCommand(ctx context.Context, startDir, name string, dryRun, force bool) error

	// MakeQuery generates a CQRS Query DTO and read-model handler.
	MakeQuery(ctx context.Context, startDir, name string, dryRun, force bool) error

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

	// AddTest scaffolds standard integration test setup and helpers.
	AddTest(ctx context.Context, startDir string, dryRun, force bool) error

	// AddMultitenancy scaffolds Row Level Security (RLS) SQL policies for tenant isolation.
	AddMultitenancy(ctx context.Context, startDir, moduleName string, dryRun, force bool) error

	// AddCQRS scaffolds Command and Query handlers separated within module scope.
	AddCQRS(ctx context.Context, startDir, moduleName string, dryRun, force bool) error

	// AddOutbox scaffolds Transactional Outbox pattern infrastructure and SQL migrations.
	AddOutbox(ctx context.Context, startDir string, dryRun, force bool) error

	// AddSaga scaffolds Distributed Saga workflow engine and step orchestrators.
	AddSaga(ctx context.Context, startDir, workflowName string, dryRun, force bool) error

	// AddWebhook scaffolds secure signed webhook dispatcher and signature verifier middleware.
	AddWebhook(ctx context.Context, startDir string, dryRun, force bool) error

	// AddDiscovery scaffolds service discovery registration client (Consul or etcd).
	AddDiscovery(ctx context.Context, startDir, provider string, dryRun, force bool) error

	// AddAuth scaffolds authentication handlers and middleware (OAuth2 or APIKey).
	AddAuth(ctx context.Context, startDir, authType string, dryRun, force bool) error

	// AddStorage scaffolds cloud blob storage interface and S3/MinIO adapter.
	AddStorage(ctx context.Context, startDir, provider string, dryRun, force bool) error

	// AddCron scaffolds in-process recurring background job scheduler.
	AddCron(ctx context.Context, startDir, jobName string, dryRun, force bool) error

	// AddMailer scaffolds transactional email delivery client.
	AddMailer(ctx context.Context, startDir, provider string, dryRun, force bool) error

	// AddFirebase scaffolds Firebase Auth token verifier and FCM push messaging client.
	AddFirebase(ctx context.Context, startDir string, dryRun, force bool) error

	// AddLogger scaffolds structured JSON logger with context correlation tracking.
	AddLogger(ctx context.Context, startDir, provider string, dryRun, force bool) error

	// AddHealthcheck scaffolds Kubernetes-compatible Liveness & Readiness probe handlers.
	AddHealthcheck(ctx context.Context, startDir string, dryRun, force bool) error

	// AddSecrets scaffolds secret manager client (HashiCorp Vault or AWS Secrets Manager).
	AddSecrets(ctx context.Context, startDir, provider string, dryRun, force bool) error

	// AddLock scaffolds distributed lock (Redlock / Distributed Mutex) engine.
	AddLock(ctx context.Context, startDir, provider string, dryRun, force bool) error

	// AddAuthz scaffolds RBAC / ABAC authorization engine and middleware (Casbin).
	AddAuthz(ctx context.Context, startDir, provider string, dryRun, force bool) error

	// AddCrypto scaffolds symmetric envelope encryption helper (AES-GCM).
	AddCrypto(ctx context.Context, startDir, algorithm string, dryRun, force bool) error

	// AddProfiling scaffolds protected runtime pprof / Pyroscope profiling endpoints.
	AddProfiling(ctx context.Context, startDir, provider string, dryRun, force bool) error

	// AddFeatureFlags scaffolds feature flag & canary release client (Flipt).
	AddFeatureFlags(ctx context.Context, startDir, provider string, dryRun, force bool) error

	// AddIdempotency scaffolds Idempotency-Key validation middleware with distributed lock.
	AddIdempotency(ctx context.Context, startDir, provider string, dryRun, force bool) error

	// AddLedger scaffolds Double-Entry Bookkeeping Ledger engine.
	AddLedger(ctx context.Context, startDir string, dryRun, force bool) error

	// AddDecimal scaffolds high-precision financial decimal money arithmetic helpers.
	AddDecimal(ctx context.Context, startDir string, dryRun, force bool) error

	// AddReconciliation scaffolds automated transaction settlement reconciliation matching engine.
	AddReconciliation(ctx context.Context, startDir string, dryRun, force bool) error

	// AddPricingEngine scaffolds rule-based tiered pricing and fee calculation engine.
	AddPricingEngine(ctx context.Context, startDir string, dryRun, force bool) error

	// AddWebSocket scaffolds WebSocket hub and multi-client connection pool (Gorilla / Nhooyr).
	AddWebSocket(ctx context.Context, startDir, provider string, dryRun, force bool) error

	// AddSSE scaffolds Server-Sent Events unidirectional streaming broker.
	AddSSE(ctx context.Context, startDir string, dryRun, force bool) error

	// AddWebRTC scaffolds Pion WebRTC peer connection and data channel signaling hub.
	AddWebRTC(ctx context.Context, startDir, provider string, dryRun, force bool) error

	// AddMQTT scaffolds Paho MQTT client for IoT telemetry pub/sub messaging.
	AddMQTT(ctx context.Context, startDir, provider string, dryRun, force bool) error

	// AddTwilio scaffolds Twilio SMS and WhatsApp omni-channel delivery client.
	AddTwilio(ctx context.Context, startDir string, dryRun, force bool) error

	// AddMultiLevelCache scaffolds synchronized L1 memory + L2 Redis cache.
	AddMultiLevelCache(ctx context.Context, startDir string, dryRun, force bool) error

	// AddBloomFilter scaffolds probabilistic cache penetration guard.
	AddBloomFilter(ctx context.Context, startDir string, dryRun, force bool) error

	// AddS3 scaffolds S3 object storage client with pre-signed URL generator (MinIO / AWS).
	AddS3(ctx context.Context, startDir, provider string, dryRun, force bool) error

	// AddResilience scaffolds Circuit Breaker and Bulkhead resilience engine.
	AddResilience(ctx context.Context, startDir, provider string, dryRun, force bool) error

	// AddSearch scaffolds full-text typo-tolerant search engine client (Meilisearch / Elasticsearch).
	AddSearch(ctx context.Context, startDir, provider string, dryRun, force bool) error

	// TestStress scaffolds high-concurrency load testing suite (k6 / vegeta).
	TestStress(ctx context.Context, startDir, engine string, dryRun, force bool) error

	// TestChaos scaffolds chaos engineering fault injection middleware.
	TestChaos(ctx context.Context, startDir string, dryRun, force bool) error

	// TestFuzz scaffolds Go native continuous fuzz testing harness.
	TestFuzz(ctx context.Context, startDir, target string, dryRun, force bool) error

	// TestBenchmark scaffolds micro-benchmark and memory allocation profiler.
	TestBenchmark(ctx context.Context, startDir string, dryRun, force bool) error

	// TestContainer scaffolds Testcontainers integration testing harness for PostgreSQL and Redis.
	TestContainer(ctx context.Context, startDir, provider string, dryRun, force bool) error

	// TestMutation scaffolds mutation testing verification harness.
	TestMutation(ctx context.Context, startDir string, dryRun, force bool) error

	// MakeMock scaffolds interface mock implementation using Mockery directives.
	MakeMock(ctx context.Context, startDir, interfaceName string, dryRun, force bool) error

	// AddSQLC scaffolds SQLC configuration, base schema, and type-safe query templates.
	AddSQLC(ctx context.Context, startDir string, dryRun, force bool) error

	// AddGRPCStream scaffolds gRPC bi-directional, client, and server streaming RPC handlers and protobuf definition.
	AddGRPCStream(ctx context.Context, startDir string, dryRun, force bool) error

	// AddGRPCGateway scaffolds gRPC-Gateway reverse-proxy HTTP REST JSON bridge.
	AddGRPCGateway(ctx context.Context, startDir string, dryRun, force bool) error

	// AddTenantContext scaffolds multi-tenancy middleware, tenant context extractor, and DB scoping helper.
	AddTenantContext(ctx context.Context, startDir string, dryRun, force bool) error

	// MakePipeline scaffolds Fan-Out / Fan-In bounded concurrency pipeline helper.
	MakePipeline(ctx context.Context, startDir, name string, dryRun, force bool) error

	// MakeSpecification scaffolds reusable DDD Specification pattern for dynamic query rules.
	MakeSpecification(ctx context.Context, startDir, name string, dryRun, force bool) error

	// AddSingleflight scaffolds request deduplication helper to prevent cache stampede.
	AddSingleflight(ctx context.Context, startDir string, dryRun, force bool) error

	// AddDrain scaffolds zero-downtime graceful shutdown and connection draining manager.
	AddDrain(ctx context.Context, startDir string, dryRun, force bool) error

	// AddOAuth2 scaffolds OIDC/OAuth2 login client with PKCE state verification.
	AddOAuth2(ctx context.Context, startDir, provider string, dryRun, force bool) error

	// AddAuditLog scaffolds tamper-evident immutable audit log with PII scrubbing.
	AddAuditLog(ctx context.Context, startDir string, dryRun, force bool) error

	// AddArgon2 scaffolds GPU-resistant Argon2id password security hasher.
	AddArgon2(ctx context.Context, startDir string, dryRun, force bool) error

	// AddLint scaffolds linter rules (golangci-lint) and git pre-commit hooks.
	AddLint(ctx context.Context, startDir string, dryRun, force bool) error

	// AddUOW scaffolds Unit of Work transactional orchestrator pattern.
	AddUOW(ctx context.Context, startDir string, dryRun, force bool) error

	// AddGraphQL scaffolds gqlgen GraphQL server with dataloader boilerplate.
	AddGraphQL(ctx context.Context, startDir string, dryRun, force bool) error

	// AddReadReplica scaffolds DB connection pool splitter for Primary-Write and Replica-Read.
	AddReadReplica(ctx context.Context, startDir string, dryRun, force bool) error

	// AddOpenAPI scaffolds Swagger OpenAPI docs generator middleware.
	AddOpenAPI(ctx context.Context, startDir string, dryRun, force bool) error

	// MakeCursorPaginator scaffolds cursor-based opaque base64 pagination helper.
	MakeCursorPaginator(ctx context.Context, startDir string, dryRun, force bool) error
}
