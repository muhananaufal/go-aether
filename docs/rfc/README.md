# 📚 RFC Architecture Catalog & Master Index (`go-aether`)

Dokumentasi rancangan arsitektur, keputusan sistem (ADR), dan riwayat implementasi teknis untuk **`go-aether` — Opinionated Architecture Scaffold CLI Engine for Go Backend Engineers**.

---

## 🏛️ Master RFC & Release Registry

| RFC ID | Inisiatif Fitur / Arsitektur | Domain / Modul | Target Branch | Versi Release | Status | Tanggal Rilis / Target |
| :--- | :--- | :--- | :--- | :---: | :---: | :---: |
| `20260806-01` | [Phase 1: Hexagonal DDD Core Engine](20260806-go-aether-core-scaffolding-engine.md) | Core CLI / Scaffold Engine | `feature/core-engine` | `v0.1.0` | `IMPLEMENTED` | 2026-08-06 |
| `20260806-02` | [Phase 2: Granular Ecosystem & Middlewares](20260806-go-aether-phase2-ecosystem.md) | Ecosystem Addons | `feature/phase2-ecosystem` | `v0.2.0` | `IMPLEMENTED` | 2026-08-06 |
| `20260806-03` | [Phase 3: Distributed Systems & Observability](20260806-go-aether-phase3-distributed.md) | Distributed & OTel | `feature/phase3-distributed` | `v0.3.0` | `IMPLEMENTED` | 2026-08-06 |
| `20260806-04` | [Phase 4: Cloud-Native, CI/CD & AI](20260806-go-aether-phase4-cloudnative.md) | Cloud & CLI TUI | `feature/phase4-cloudnative` | `v0.4.0` | `IMPLEMENTED` | 2026-08-06 |
| `20260806-05` | [Phase 5: Core Enterprise Patterns](20260806-go-aether-phase5-core-patterns.md) | Architectural Patterns | `feature/phase5-core-patterns` | `v0.5.0` | `IMPLEMENTED` | 2026-08-06 |
| `20260806-06` | [Phase 6: Data, Migration & QA](20260806-go-aether-phase6-data-qa.md) | Database & Test | `feature/phase6-data-qa` | `v0.6.0` | `IMPLEMENTED` | 2026-08-06 |
| `20260806-07` | [Phase 7: Advanced Distributed Patterns](20260806-go-aether-phase7-distributed-patterns.md) | Distributed Patterns | `feature/phase7-distributed-patterns` | `v0.7.0` | `IMPLEMENTED` | 2026-08-06 |
| `20260806-08` | [Phase 8: Cloud, Auth & 3rd-Party Plugins](20260806-go-aether-phase8-cloud-plugins.md) | Plugins & Ecosystem | `feature/phase8-cloud-plugins` | `v0.8.0` | `IMPLEMENTED` | 2026-08-06 |
| `20260806-09` | Utility & Security Extensions (`ls`, `healthcheck`, `secrets`, `lock`) | Core Extensions | `feature/v0.8.1-utility-plugins` | `v0.8.1` | `IMPLEMENTED` | 2026-08-06 |
| `20260806-10` | Advanced Security & Ops (`authz`, `crypto`, `profiling`, `featureflags`) | Security & Ops | `feature/v0.8.2-advanced-security-ops` | `v0.8.2` | `IMPLEMENTED` | 2026-08-06 |
| `20260806-11` | [Phase 9: Fintech & Financial Reliability Engine](20260806-go-aether-phase9-fintech-engine.md) | Fintech Reliability | `feature/v0.8.3-fintech-engine` | `v0.8.3` | `IMPLEMENTED` | 2026-08-06 |
| `20260806-12` | [Phase 10: Tactical Domain-Driven Design (DDD) Code-Gen](20260806-go-aether-phase10-tactical-ddd.md) | Tactical DDD | `feature/v0.8.4-tactical-ddd` | `v0.8.4` | `IMPLEMENTED` | 2026-08-06 |
| `20260806-13` | [Phase 11: Realtime Protocols & Telecom](20260806-go-aether-phase11-concurrency-resilience.md) | Realtime Networking | `feature/v0.8.5-concurrency-resilience` | `v0.8.5` | `IMPLEMENTED` | 2026-08-06 |
| `20260806-14` | [Phase 12: Caching & Cloud Storage Engine](20260806-go-aether-phase12-database-storage.md) | Cache & S3 Storage | `feature/v0.8.6-database-storage` | `v0.8.6` | `IMPLEMENTED` | 2026-08-06 |
| `20260806-15` | [Phase 13: Native QA, Testing, Stress & Chaos Engine](20260806-go-aether-phase13-testing-qa-engine.md) | Testing & Resilience | `feature/v0.8.7-testing-qa-engine` | `v0.8.7` | `IMPLEMENTED` | 2026-08-06 |
| `20260806-16` | [Phase 14: Enterprise SQLC, Advanced gRPC Streaming & Multi-Tenancy](20260806-go-aether-phase14-enterprise-sqlc-grpc-multitenancy.md) | SQLC, gRPC & Tenant | `feature/v0.8.8-sqlc-grpc-multitenancy` | `v0.8.8` | `IMPLEMENTED` | 2026-08-06 |
| `20260806-17` | [Phase 15: Concurrency Pipelines, Zero-Trust & Observability](20260806-go-aether-phase15-concurrency-zerotrust.md) | Concurrency & Security | `feature/v0.8.9-concurrency-zerotrust` | `v0.8.9` | `IMPLEMENTED` | 2026-08-06 |
| `20260807-01` | [Phase 16: The Final Gap to v1.0.0](20260807-go-aether-v0.9.0-final-gap.md) | CLI Gap Closure | `feature/v0.9.0-final-gap` | `v0.9.0` | `IMPLEMENTED` | 2026-08-07 |
| `20260807-02` | [Advanced UX & DX Deep Dive Overhaul](20260807-ux-dx-deep-dive.md) | Core CLI / UX | `feature/ux-dx-deep-dive` | `v0.2.0` | `PROPOSED` | - |
| `20260807-03` | [v0.4.0: Production Hardening, Framework Matrix & Real Brownfield Engine](20260807-v0.4.0-production-hardening.md) | Core Quality / Brownfield | `feature/v0.3.0-production-hardening` | `v0.4.0` | `ACCEPTED` | - |

---

## 🗺️ Master Evolution Roadmap: 16 Phases to v1.0.0

```mermaid
flowchart TD
    subgraph Selesai [✅ SUDAH DIRILIS (v0.1.0 - v0.8.6: 66 Command)]
        F1["v0.1.0 - v0.8.2: Core Hexagonal, Ecosystem, Observability & Security"]
        F1 --> F9["v0.8.3: Fintech & Financial Engine (5 Cmd)"]
        F9 --> F10["v0.8.4: Tactical DDD Granular (5 Cmd)"]
        F10 --> F11["v0.8.5: Realtime Protocols & Telecom (5 Cmd)"]
        F11 --> F12["v0.8.6: Caching & Cloud Storage Engine (5 Cmd)"]
    end

    subgraph CoreEnterpriseGap [🎯 FOKUS PENUTUPAN GAP CORE & ENTERPRISE (29 Command)]
        F12 --> F13["v0.8.7: Native QA, Stress & Chaos Engine (7 Cmd)"]
        F13 --> F14["v0.8.8: Enterprise SQLC, Unit of Work & gRPC/GraphQL (8 Cmd)"]
        F14 --> F15["v0.8.9: Concurrency Pipelines, Zero-Trust & Observability (8 Cmd)"]
        F15 --> F16["v0.9.0: The Final Gap (6 Cmd)"]
        F16 --> V1["🚀 v1.0.0: Grand Parity GA (90 Total Commands)"]
    end
```

---

## 📦 Detail Cakupan 3 Batch Penutup Core & Enterprise:

### 🧪 Batch 5 — Fase 13: Native QA, Testing, Stress & Chaos Engine (`v0.8.7`)
- `test:stress [k6|vegeta]`: High-concurrency load testing harness & RPS throughput profiler.
- `test:chaos`: Chaos engineering failure injector middleware (latency, random abort, packet drop).
- `test:fuzz`: Continuous Go native fuzz testing harness (`go test -fuzz`) untuk menemukan edge-case panics.
- `test:mutation`: Mutation testing runner script untuk memvalidasi ketajaman assertion unit test.
- `test:benchmark`: Micro-benchmark suite + memory allocation & escape analysis profiler (`testing.B`).
- `test:container`: Testcontainers integration testing harness (PostgreSQL & Redis di Docker otomatis).
- `make:mock [interface]`: Interface mock generator berbasis Mockery untuk isolated unit tests.

### 🗄️ Batch 6 — Fase 14: Enterprise SQLC, Unit of Work & Microservices gRPC (`v0.8.8`)
- `add:sqlc`: Scaffolder `sqlc.yaml`, native SQL queries, dan type-safe generator berbasis `pgx/v5`.
- `add:uow`: Unit of Work pattern untuk multi-repository transactional orchestration.
- `make:proto [name]`: Scaffolder `.proto` service definition dengan Protobuf v3, validation rules, dan Buf linter.
- `add:grpc-gateway`: In-process HTTP/REST JSON reverse proxy ke gRPC service.
- `add:graphql [gqlgen]`: GraphQL schema-first server dengan DataLoader anti N+1.
- `add:readreplica`: Connection pool manager pemisah Primary (Write) vs Read-Replica (Read) dengan auto-fallback.
- `make:cursor-paginator`: Reusable opaque base64-encoded cursor-based pagination helper.
- `add:openapi`: Swagger UI / OpenAPI 3.0 documentation server generator.

### ⚡ Batch 7 — Fase 15: Concurrency Pipelines, Zero-Trust Security & Observability (`v0.8.9`)
- `make:pipeline [name]`: Fan-Out / Fan-In concurrency pipeline dengan context cancellation dan bounded worker pools.
- `add:singleflight`: Request deduplication helper (`x/sync/singleflight`) anti thundering herd / cache stampede.
- `add:metrics [prometheus]`: Dedicated Prometheus RED metrics collector & middleware (Rate, Errors, Duration).
- `add:drain`: Zero-downtime graceful shutdown & connection draining manager.
- `add:oauth2 [provider]`: Social Login (Google, GitHub, OIDC) dengan PKCE state verification.
- `add:auditlog`: Tamper-evident immutable audit log dengan automatic PII scrubbing.
- `add:argon2`: Argon2id password security hasher.
- `make:specification`: DDD Specification pattern helper untuk reusable dynamic query rules.
