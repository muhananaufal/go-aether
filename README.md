<p align="center">
  <h1 align="center">🚀 go-aether</h1>
  <p align="center">
    <strong>A Zero-Runtime, Domain-First Scaffolding CLI Engine for Production-Grade Go Backends.</strong>
  </p>
  <p align="center">
    <a href="https://golang.org/doc/devel/release.html"><img src="https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat-square&logo=go" alt="Go Version"></a>
    <a href="https://github.com/muhananaufal/go-aether/releases"><img src="https://img.shields.io/github/v/release/muhananaufal/go-aether?style=flat-square" alt="Release"></a>
    <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square" alt="License"></a>
  </p>
</p>

---

**`go-aether`** is a developer-first CLI scaffolding engine designed for Go Backend Engineers and Enterprise Architects. It eliminates the friction of boilerplate construction by generating production-grade, strictly typed code across multiple architectural styles — including **Hexagonal (Ports & Adapters)**, **Domain-Driven Design (DDD)**, **CQRS**, **Distributed Systems**, **Fintech**, and **Realtime Protocols** — in milliseconds.

Unlike opinionated frameworks, `go-aether` is **strictly a Dev-Time Scaffolding Tool**. It embeds clean standard `text/template` files into a single self-contained binary, generating pure, unencumbered Go code with **zero runtime lock-in**.

---

## ✨ Core Highlights

- 🏗️ **Domain-First Scaffolding**: 90 commands organized into semantic groups (`arch:*`, `db:*`, `api:*`, `security:*`, `fintech:*`, etc.)
- ⚡ **Zero-Runtime Overhead**: All generated code is 100% native Go — no external runtime frameworks, no hidden dependencies.
- 🛡️ **Transactional Disk Buffer**: Atomic multi-file generation with automatic rollback on error. Your filesystem is never left in a partial state.
- 🩺 **Structural Health Check**: Built-in `doctor` command for AST-based diagnostics against the Single Source of Truth (`aether.yaml`).
- 🔄 **Enterprise Patterns Built-In**: CQRS, Transactional Outbox, Distributed Saga, Unit of Work, and more — generated with a single command.
- 🔒 **Security & Observability First**: Scaffold OpenTelemetry tracing, Prometheus metrics, Argon2id hashing, Casbin RBAC, and Zero-Trust patterns out of the box.
- 🔍 **Brownfield Adoption**: Scan existing Go codebases and seamlessly integrate them into the Aether manifest.

---

## 🚀 Installation

```bash
go install github.com/muhananaufal/go-aether@latest
```

*Ensure `$(go env GOPATH)/bin` is in your system `$PATH`.*

---

## 📖 Quick Start

### 1. Initialize a New Project

Bootstrap a clean architecture workspace and generate the master `aether.yaml` manifest:

```bash
mkdir my-service && cd my-service
go-aether init my-service --arch hexagonal --router chi --db postgres
```

### 2. Scaffold a Domain Module

Generate a complete vertical slice (`order`) — Domain Entity, Port Interface, Service, HTTP Handler, and Postgres Repository:

```bash
go-aether arch:module order --transports http
```

### 3. Add Enterprise Patterns

```bash
# Add CQRS Command & Query Bus for the order domain
go-aether platform:cqrs order

# Add Transactional Outbox + SQL migrations
go-aether async:outbox

# Add Distributed Saga with compensation workflow
go-aether async:saga checkout

# Add Unit of Work transactional orchestrator
go-aether db:uow

# Add Read/Write database replica splitter
go-aether db:readreplica
```

### 4. Add Security & Observability

```bash
go-aether security:auth oauth2
go-aether security:argon2
go-aether o11y:tracing otel
go-aether o11y:metrics prometheus
go-aether infra:healthcheck
```

### 5. Run Structural Diagnostics

```bash
go-aether doctor
```

---

## 🚩 Important Command Flags

While `go-aether` has many commands, most core scaffolding commands share a common set of flags. 

> [!NOTE]
> **Design Philosophy: Commands vs Flags**
> - **Flags** (e.g., `arch:module --transports grpc`): Used during **Creation-Time**. They inject specific capabilities into a module *as it is being generated*.
> - **Commands** (e.g., `api:transport grpc`): Used for **Global/Standalone Setup**. They scaffold global infrastructure at the project root or inject new capabilities into an *existing* codebase retrospectively (Brownfield).

### Project Initialization (`init`)
- `--module` : Go module identifier for `go.mod` (default: `github.com/example/app`).
- `--arch` : Architecture blueprint (options: `hexagonal`, `clean`, `ddd`; default: `hexagonal`).
- `--db` : Database engine driver (options: `postgres`, `mysql`, `sqlite`; default: `postgres`).
- `--router` : HTTP routing framework (options: `chi`, `gin`, `echo`; default: `chi`).

### Module Scaffolding (`arch:module`)
- `--transports` : Comma-separated transport targets (e.g., `http,grpc,nats`; default: `http`).
- `--cache` : Inject L1/L2 Redis caching decorators into the repository layer.
- `--worker` : Scaffold async worker job processor for domain events.

### Global Flags (Available on all commands)
- `--dry-run` : Simulate generation and print to stdout without writing to disk.
- `--force` or `-f` : Force overwrite existing generated files.

---

## 🗺️ Complete Command Reference (v0.9.0 — 90 Commands)

> [!TIP]
> **Extensibility Note**: Arguments like `[provider]`, `[type]`, or `[algo]` are designed for forward-compatibility. In `v0.9.0`, we support industry standards out-of-the-box (e.g., `--db postgres`, `security:oauth2 google`, `o11y:tracing otel`). More specialized providers will be unlocked in v1.0.0.

### ⚙️ Core Lifecycle
| Command | Description |
| :--- | :--- |
| `init [name]` | Bootstrap project layout and `aether.yaml` manifest |
| `adopt` | Scan and adopt a legacy brownfield Go repository |
| `doctor` | Run structural health diagnostics against `aether.yaml` |
| `ls` | List active modules, architectural paths, and installed plugins |

---

### 🏗️ Architecture Scaffolding (`arch:*`)
| Command | Description |
| :--- | :--- |
| `arch:module [name]` | Scaffold a complete vertical domain slice (Domain, Port, Service, Handler, Repo) |
| `arch:domain [name]` | Scaffold the domain layer entity only |
| `arch:port [name]` | Scaffold the port interface contract only |
| `arch:service [name]` | Scaffold the service layer only |
| `arch:handler [name]` | Scaffold the transport handler only |
| `arch:repository [name]` | Scaffold the infrastructure repository only |
| `arch:aggregate [name]` | Generate a DDD Aggregate Root with event recording |
| `arch:event [name]` | Generate a Domain Event struct and serializer |
| `arch:valueobject [name]` | Generate an immutable DDD Value Object |
| `arch:command [name]` | Generate a CQRS Command DTO and execution handler |
| `arch:query [name]` | Generate a CQRS Query DTO and read-model handler |
| `arch:mock [interface]` | Scaffold an interface mock via Mockery for isolated unit tests |
| `arch:pipeline [name]` | Generate a Fan-Out / Fan-In bounded concurrency pipeline |
| `arch:specification [name]` | Generate a DDD Specification pattern for dynamic query rules |
| `arch:di [type]` | Set up a Dependency Injection container (fx, wire) |
| `arch:error` | Set up a standardized centralized error handler |

---

### 🗄️ Database & Persistence (`db:*`)
| Command | Description |
| :--- | :--- |
| `db:migration [name]` | Generate a SQL migration file pair (up/down) |
| `db:seeder [name]` | Generate a database seeder file |
| `db:sqlc` | Set up SQLC type-safe query code generator |
| `db:uow` | Set up Unit of Work transactional orchestrator |
| `db:readreplica` | Set up Primary-Write / Replica-Read connection pool splitter |
| `db:paginator` | Generate a cursor-based opaque base64 pagination helper |

---

### 🌐 API Layer (`api:*`)
| Command | Description |
| :--- | :--- |
| `api:graphql` | Set up gqlgen GraphQL server with DataLoader boilerplate |
| `api:openapi` | Set up Swagger / OpenAPI 3.0 documentation middleware |
| `api:grpc` | Set up gRPC bi-directional duplex streaming handler |
| `api:grpc-gateway` | Set up gRPC-Gateway REST/JSON reverse-proxy |
| `api:middleware [type]` | Inject a middleware (jwt-auth, rate-limit) into a module handler |
| `api:transport [type]` | Register a new global transport protocol in `aether.yaml` |
| `api:validator [type]` | Set up struct validation wrapper (playground) |
| `api:idempotency [type]` | Set up Idempotency-Key validation middleware |

---

### ⚡ Realtime Protocols (`realtime:*`)
| Command | Description |
| :--- | :--- |
| `realtime:ws [provider]` | Set up WebSocket hub and connection pool |
| `realtime:sse` | Set up Server-Sent Events (SSE) streaming broker |
| `realtime:webrtc [pion]` | Set up Pion WebRTC peer-to-peer data channel hub |
| `realtime:mqtt [paho]` | Set up Paho MQTT client for IoT telemetry |

---

### 📨 Async Messaging & Events (`async:*`)
| Command | Description |
| :--- | :--- |
| `async:worker [name]` | Generate an asynchronous background job processor |
| `async:cron [job-name]` | Set up in-process recurring cron scheduler |
| `async:outbox` | Set up Transactional Outbox pattern with SQL migrations |
| `async:saga [workflow]` | Set up Distributed Saga orchestrator with compensation |
| `async:eventing` | Set up global Pub/Sub event bus interfaces |
| `async:webhook` | Set up HMAC-SHA256 signed webhook dispatcher |

---

### 🔐 Security & Zero-Trust (`security:*`)
| Command | Description |
| :--- | :--- |
| `security:auth [type]` | Set up authentication middleware (oauth2, apikey) |
| `security:authz [casbin]` | Set up RBAC / ABAC authorization engine (Casbin) |
| `security:oauth2 [provider]` | Set up OIDC/OAuth2 login with PKCE state verification |
| `security:argon2` | Set up GPU-resistant Argon2id password hasher |
| `security:crypto [algo]` | Set up symmetric envelope encryption (AES-256-GCM) |
| `security:secrets [provider]` | Set up secret manager client (Vault, AWS) |
| `security:auditlog` | Set up tamper-evident immutable audit log with PII scrubbing |

---

### 💾 Caching Layer (`cache:*`)
| Command | Description |
| :--- | :--- |
| `cache:redis [type]` | Set up global cache layer (Redis, Valkey) |
| `cache:multilevel` | Set up synchronized L1 (memory) + L2 (Redis) cache |
| `cache:bloom` | Set up probabilistic Bloom Filter cache penetration guard |
| `cache:dedup` | Set up Singleflight request deduplication (anti thundering herd) |

---

### 🏭 Infrastructure & DevOps (`infra:*`)
| Command | Description |
| :--- | :--- |
| `infra:cicd [provider]` | Set up CI/CD pipelines (GitHub Actions, GitLab CI) |
| `infra:deploy [target]` | Set up deployment manifests (K8s, Helm, Lambda) |
| `infra:healthcheck` | Set up Kubernetes Liveness `/livez` and Readiness `/readyz` probes |
| `infra:drain` | Set up zero-downtime graceful shutdown and draining manager |
| `infra:profiling [pprof]` | Set up protected runtime profiling endpoints |
| `infra:lint` | Set up golangci-lint configuration and git pre-commit hooks |
| `infra:featureflags [provider]` | Set up feature flags and canary release client (Flipt) |
| `infra:config [type]` | Set up centralized config manager (Viper, Koanf) |
| `infra:search [provider]` | Set up full-text search engine client (Elasticsearch, Meilisearch) |

---

### 📊 Observability (`o11y:*`)
| Command | Description |
| :--- | :--- |
| `o11y:tracing [exporter]` | Set up OpenTelemetry distributed tracing (Jaeger, stdout) |
| `o11y:metrics [provider]` | Set up Prometheus RED metrics middleware and endpoint |
| `o11y:logger [provider]` | Set up structured JSON logger with trace correlation (slog, zap) |

---

### 💰 Fintech & Financial Reliability (`fintech:*`)
| Command | Description |
| :--- | :--- |
| `fintech:ledger` | Set up Double-Entry bookkeeping ledger engine |
| `fintech:decimal` | Set up high-precision decimal money arithmetic helpers |
| `fintech:reconcile` | Set up automated settlement and transaction reconciliation |
| `fintech:pricing` | Set up rule-based tiered pricing and fee calculation engine |

---

### 📣 Notifications & Communication (`notif:*`)
| Command | Description |
| :--- | :--- |
| `notif:sms` | Set up Twilio SMS and WhatsApp omni-channel messaging |
| `notif:push` | Set up Firebase Auth and FCM push notifications |
| `notif:mail [provider]` | Set up transactional email client (SMTP, Resend) |

---

### ☁️ Cloud Storage (`cloud:*`)
| Command | Description |
| :--- | :--- |
| `cloud:s3 [provider]` | Set up S3 object storage client with pre-signed URL generator |
| `cloud:storage [provider]` | Set up cloud blob storage adapter (S3, GCS, local) |

---

### 🌐 Distributed Platform (`platform:*`)
| Command | Description |
| :--- | :--- |
| `platform:discovery [provider]` | Set up service discovery client (Consul, etcd) |
| `platform:lock [redis]` | Set up distributed Redlock mutex engine |
| `platform:resilience [provider]` | Set up Circuit Breaker and Bulkhead resilience engine |
| `platform:cqrs [module]` | Set up CQRS Command and Query handlers for a module |
| `platform:multitenancy [module]` | Set up Row Level Security (RLS) SQL tenant isolation |
| `platform:tenant` | Set up tenant context middleware and isolation helper |
| `platform:ai [provider]` | Set up LLM proxy infrastructure (OpenAI, Ollama) |

---

### 🧪 Testing & QA Engine (`test:*`)
| Command | Description |
| :--- | :--- |
| `test:stress` | Scaffold high-concurrency load testing suite (k6 / Vegeta) |
| `test:chaos` | Scaffold chaos engineering fault injection middleware |
| `test:fuzz` | Scaffold Go native continuous fuzz testing harness |
| `test:mutation` | Scaffold mutation testing verification runner |
| `test:benchmark` | Scaffold micro-benchmark suite and memory allocation profiler |
| `test:container` | Scaffold Testcontainers integration testing (PostgreSQL + Redis) |
| `test:integration` | Set up integration test helpers and mocking base |

---

## 📂 Generated Project Layout

```text
my-service/
├── aether.yaml                     <- Single Source of Truth Manifest
├── cmd/
│   └── api/
│       └── main.go                 <- Application Entrypoint
├── internal/
│   ├── core/
│   │   ├── domain/                 <- Pure Business Logic & Entities
│   │   ├── port/                   <- Inbound & Outbound Interface Contracts
│   │   └── service/                <- Use Cases & Application Orchestration
│   ├── adapter/
│   │   ├── handler/http/           <- HTTP Controllers
│   │   └── repository/             <- Database Implementations
│   ├── jobs/                       <- Background Cron Workloads
│   └── workflows/                  <- Saga Orchestration Workflows
├── pkg/
│   ├── auth/                       <- OAuth2 & API Key Validators
│   ├── cache/                      <- L1/L2 Multi-Level Cache
│   ├── concurrency/                <- Pipeline & Singleflight Helpers
│   ├── middleware/                  <- Metrics, Audit Log, Rate Limit
│   ├── security/                   <- Argon2id, Crypto, Secrets
│   ├── server/                     <- Drain & Health Check Handlers
│   └── tenant/                     <- Multi-Tenant Context Isolation
└── migrations/                     <- SQL Migrations (Up/Down)
```

---

## 📜 License

This project is open-sourced under the [MIT License](LICENSE) © 2026 muhananaufal.
