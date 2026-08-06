<p align="center">
  <h1 align="center">🚀 go-aether</h1>
  <p align="center">
    <strong>Lightning-fast, Zero-Runtime Dependency, Opinionated Enterprise Architecture CLI Engine for Go.</strong>
  </p>
  <p align="center">
    <a href="https://github.com/muhananaufal/go-aether/actions"><img src="https://img.shields.io/github/actions/workflow/status/muhananaufal/go-aether/ci.yml?branch=main&style=flat-square" alt="Build Status"></a>
    <a href="https://golang.org/doc/devel/release.html"><img src="https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat-square&logo=go" alt="Go Version"></a>
    <a href="https://github.com/muhananaufal/go-aether/releases"><img src="https://img.shields.io/github/v/release/muhananaufal/go-aether?style=flat-square" alt="Release"></a>
    <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square" alt="License"></a>
  </p>
</p>

---

**`go-aether`** is a high-performance, developer-first CLI engine designed for modern Go Backend Engineers and Enterprise Architects. It eliminates the friction of boilerplate construction by scaffolding production-grade, strictly typed **Hexagonal (Ports & Adapters)** and **Domain-Driven Design (DDD)** architectures in milliseconds.

Unlike heavy opinionated frameworks, `go-aether` is **strictly a Dev-Time Scaffolding Tool**. It embeds clean standard `text/template` files into a single binary, generating pure, unencumbered Go code with **zero runtime lock-in**.

---

## ✨ Core Highlights & Capabilities

- 🏗️ **Hexagonal Architecture (Ports & Adapters)**: Automatic separation of concerns across Domain, Port, Service, and Adapter layers.
- ⚡ **Zero-Runtime Overhead**: The generated code is 100% native Go. No external runtime frameworks required.
- 🛡️ **Transactional Disk Buffer**: Atomic multi-file writes with automatic rollback on generation error.
- 🩺 **Aether Doctor**: Built-in AST structural diagnostics against the Single Source of Truth (`aether.yaml`).
- 🔄 **Distributed Systems Ready**: Out-of-the-box CQRS, Transactional Outbox, Distributed Saga Orchestrator, and HMAC Signed Webhooks.
- ☁️ **Cloud Native & Plugins**: Native AWS S3, OAuth2 / API Key, In-Process Cron Scheduler, Mailer, Firebase Auth/FCM, and OpenTelemetry Tracing.
- 🔍 **Brownfield AST Adoption**: Scan legacy Go codebases and seamlessly adopt them into the Aether manifest.

---

## 🚀 Installation

Install the latest stable release via the standard Go toolchain:

```bash
go install github.com/muhananaufal/go-aether@latest
```

*(Ensure `$(go env GOPATH)/bin` is in your system `$PATH`)*

---

## 📖 Quick Start

### 1. Initialize a New Project
Bootstrap a clean architecture workspace and generate the master `aether.yaml` manifest:

```bash
mkdir my-service && cd my-service
go-aether init my-service --arch hexagonal --router chi --db postgres
```

### 2. Scaffold Vertical Domain Modules
Generate a complete vertical slice (`order`) including Domain Entities, Ports, Service Orchestrator, HTTP Handler, and Postgres Repository:

```bash
go-aether make:module order --transports http
```

### 3. Inject Distributed Patterns & Middlewares
Effortlessly add production-grade patterns into your codebase:

```bash
# Add CQRS Handlers & Command Bus for the Order domain
go-aether add:cqrs order

# Add Transactional Outbox pattern & SQL migrations
go-aether add:outbox

# Add Distributed Saga Workflow with compensation
go-aether add:saga checkout

# Add HMAC-SHA256 Signed Webhook Dispatcher & Receiver
go-aether add:webhook

# Add Cloud Storage (AWS S3 / MinIO) and In-process Cron Scheduler
go-aether add:storage s3
go-aether add:cron cleanup_orders
```

### 4. Run Structural Diagnostics
Verify your project's architectural compliance and detect drift:

```bash
go-aether doctor
```

---

## 🗺️ Complete CLI Command Taxonomy (v0.8.3)

### 🔨 Core Generators & Inspectors (`make:*`, `init`, `ls`)
| Command | Description | Example |
| :--- | :--- | :--- |
| `ls` | List active domain modules & installed plugins | `go-aether ls` |
| `init [name]` | Initialize project layout & `aether.yaml` | `go-aether init e-commerce` |
| `make:module [name]` | Scaffold complete vertical domain slice | `go-aether make:module payment` |
| `make:service [name]` | Scaffold standalone business use-case service | `go-aether make:service order` |
| `make:handler [name]` | Scaffold HTTP / gRPC transport handler | `go-aether make:handler user` |
| `make:domain [name]` | Scaffold pure domain entity & value objects | `go-aether make:domain user` |
| `make:port [name]` | Scaffold interface contracts (Inbound/Outbound) | `go-aether make:port invoice` |
| `make:repository [name]` | Scaffold standalone persistence repository | `go-aether make:repository order` |
| `make:migration [name]` | Generate SQL migration pair (Goose / Golang-Migrate) | `go-aether make:migration add_users` |
| `make:seeder [name]` | Generate database dummy data seeder | `go-aether make:seeder initial_users` |

### 🏦 Fintech & Financial Reliability Engine (`add:*`)
| Command | Description | Example |
| :--- | :--- | :--- |
| `add:idempotency [provider]` | Idempotency-Key validation middleware with atomic lock | `go-aether add:idempotency redis` |
| `add:ledger` | Double-Entry bookkeeping ledger engine (Zero-Sum Invariant) | `go-aether add:ledger` |
| `add:decimal` | High-precision decimal money arithmetic helpers | `go-aether add:decimal` |
| `add:reconciliation` | Automated settlement & transaction reconciliation matcher | `go-aether add:reconciliation` |
| `add:pricing-engine` | Rule-based tiered pricing & dynamic fee calculator | `go-aether add:pricing-engine` |

### ⚡ Distributed Patterns, Locks, Secrets & Authz (`add:*`)
| Command | Description | Example |
| :--- | :--- | :--- |
| `add:cqrs [module]` | In-module Command & Query Handlers + Bus | `go-aether add:cqrs order` |
| `add:outbox` | Transactional Outbox engine & SQL migration | `go-aether add:outbox` |
| `add:saga [workflow]` | Distributed Saga orchestrator & rollback | `go-aether add:saga checkout` |
| `add:webhook` | HMAC-SHA256 signed webhook engine | `go-aether add:webhook` |
| `add:lock [provider]` | Distributed mutex lock engine (Redlock) | `go-aether add:lock redis` |
| `add:secrets [provider]` | Secrets manager client (Vault / AWS) | `go-aether add:secrets vault` |
| `add:authz [provider]` | RBAC / ABAC authorization engine (Casbin) | `go-aether add:authz casbin` |
| `add:crypto [algo]` | Envelope encryption helper (AES-256-GCM) | `go-aether add:crypto aes-gcm` |
| `add:featureflags [provider]` | Feature flags client (Flipt) | `go-aether add:featureflags flipt` |
| `add:discovery [provider]` | Service discovery client (Consul / etcd) | `go-aether add:discovery consul` |
| `add:cache [provider]` | High-performance caching layer (Redis / Valkey) | `go-aether add:cache redis` |
| `add:worker [provider]` | Asynchronous task queue (Asynq / River) | `go-aether add:worker asynq` |
| `add:eventing [provider]` | Event streaming broker (Kafka / RabbitMQ) | `go-aether add:eventing kafka` |

### ☁️ Cloud, Auth, Health, Profiling & Observability
| Command | Description | Example |
| :--- | :--- | :--- |
| `add:healthcheck` | K8s `/livez` & `/readyz` probe handlers | `go-aether add:healthcheck` |
| `add:profiling [provider]` | Protected runtime profiling endpoints (pprof) | `go-aether add:profiling pprof` |
| `add:auth [type]` | Authentication provider (OAuth2 / API Key) | `go-aether add:auth oauth2` |
| `add:storage [provider]` | Cloud blob storage abstraction (S3 / GCS / Local) | `go-aether add:storage s3` |
| `add:cron [job-name]` | In-process recurring job scheduler | `go-aether add:cron report_job` |
| `add:mailer [provider]` | Transactional email client (SMTP / Resend) | `go-aether add:mailer smtp` |
| `add:firebase` | Firebase Auth & FCM push notifications | `go-aether add:firebase` |
| `add:logger [provider]` | Structured JSON logger with trace correlation | `go-aether add:logger slog` |
| `add:tracing [provider]` | OpenTelemetry tracing exporter | `go-aether add:tracing otel` |
| `add:metrics [provider]` | Prometheus metric collectors | `go-aether add:metrics prometheus` |
| `add:deploy [target]` | Cloud deployment manifests (K8s, Helm, Lambda) | `go-aether add:deploy k8s` |
| `add:cicd [provider]` | Automated CI/CD workflows (GitHub Actions / GitLab) | `go-aether add:cicd github` |
| `add:ai [provider]` | AI / LLM gateway client (OpenAI, Anthropic, Ollama) | `go-aether add:ai openai` |
| `add:multitenancy [module]` | Row Level Security (RLS) tenant isolation | `go-aether add:multitenancy customer` |
| `add:test` | Integration test suite & mocking harness | `go-aether add:test` |
| `adopt [--scan]` | Scan legacy project and create `aether.yaml` | `go-aether adopt --scan` |
| `doctor` | Structural health and integrity check | `go-aether doctor` |

---

## 📂 Architecture Layout

```text
my-service/
├── aether.yaml                     <-- Single Source of Truth Manifest
├── cmd/
│   └── api/
│       └── main.go                 <-- Application Entrypoint
├── internal/
│   ├── core/
│   │   ├── domain/                 <-- Pure Business Logic & Structs
│   │   ├── port/                   <-- Inbound & Outbound Interfaces
│   │   └── service/                <-- Use Cases & Application Services
│   ├── adapter/
│   │   ├── handler/http/           <-- HTTP Controllers
│   │   └── repository/             <-- Database Implementations
│   ├── jobs/                       <-- Background Cron Workloads
│   └── workflows/                  <-- Saga Workflows
├── pkg/
│   ├── auth/                       <-- OAuth2 & API Key Validators
│   ├── cron/                       <-- Job Schedulers
│   ├── discovery/                  <-- Consul & etcd Clients
│   ├── firebase/                   <-- Firebase Auth & Push
│   ├── logger/                     <-- Slog Correlation Handler
│   ├── mailer/                     <-- Transactional Mailer
│   ├── outbox/                     <-- Transactional Outbox Engine
│   ├── saga/                       <-- Saga Orchestrator Core
│   ├── storage/                    <-- S3 Cloud Storage
│   └── webhook/                    <-- HMAC Webhook Engine
└── migrations/                     <-- SQL Migrations
```

---

## 📜 License

This project is open-sourced under the [MIT License](LICENSE) © 2026 muhananaufal.

---
<p align="center">
  <i>Engineered with extreme precision for Go Artisans by AETHERIS.</i>
</p>
