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
| `20260806-11` | [Phase 9: Fintech & Financial Reliability Engine](20260806-go-aether-phase9-fintech-engine.md) | Fintech Reliability | `feature/v0.8.3-fintech-engine` | `v0.8.3` | `PROPOSED` | Target Q3 2026 |
| `20260806-12` | [Phase 10: Tactical Domain-Driven Design (DDD) Code-Gen](20260806-go-aether-phase10-tactical-ddd.md) | Tactical DDD | `feature/v0.8.4-tactical-ddd` | `v0.8.4` | `PROPOSED` | Target Q3 2026 |
| `20260806-13` | [Phase 11: Concurrency & High-Load Resilience Engine](20260806-go-aether-phase11-concurrency-resilience.md) | Concurrency & Traffic | `feature/v0.8.5-concurrency-resilience` | `v0.8.5` | `PROPOSED` | Target Q3 2026 |
| `20260806-14` | [Phase 12: Database Extensions & Distributed Storage](20260806-go-aether-phase12-database-storage.md) | Vector, TimeSeries, KV | `feature/v0.8.6-database-storage` | `v0.8.6` | `PROPOSED` | Target Q3 2026 |
| `20260806-15` | [Phase 13: Media, Export & Document Pipeline](20260806-go-aether-phase13-media-pipeline.md) | Media & Processing | `feature/v0.8.7-media-pipeline` | `v0.8.7` | `PROPOSED` | Target Q3 2026 |
| `20260806-16` | [Phase 14: Modern API DX, GraphQL & Zero-Trust Workloads](20260806-go-aether-phase14-api-dx-zerotrust.md) | API & Networking | `feature/v0.8.8-api-dx-zerotrust` | `v0.8.8` | `PROPOSED` | Target Q3 2026 |

---

## 🗺️ Master Evolution Roadmap: 14 Phases to v1.0.0

```mermaid
flowchart TD
    subgraph Selesai [✅ SUDAH DIRILIS (46 Command)]
        F1["v0.1.0: Hexagonal DDD Core"] --> F2["v0.2.0: Granular Middlewares"]
        F2 --> F3["v0.3.0: Distributed Systems"]
        F3 --> F4["v0.4.0: Cloud-Native & AI"]
        F4 --> F5["v0.5.0: Enterprise Patterns"]
        F5 --> F6["v0.6.0: Data & QA"]
        F6 --> F7["v0.7.0: Advanced Patterns"]
        F7 --> F8["v0.8.0: Cloud/Auth Plugins"]
        F8 --> F81["v0.8.1: Utility & Probes"]
        F81 --> F82["v0.8.2: Security & Ops Parity"]
    end

    subgraph Rencana [⏳ ROADMAP BATCH TAHAP BERIKUTNYA (32 Command)]
        F82 --> F9["v0.8.3: Fintech & Financial Engine (5 Cmd)"]
        F9 --> F10["v0.8.4: Tactical DDD Granular (5 Cmd)"]
        F10 --> F11["v0.8.5: Concurrency & Resilience (6 Cmd)"]
        F11 --> F12["v0.8.6: Database & Storage (5 Cmd)"]
        F12 --> F13["v0.8.7: Media & Document Pipeline (4 Cmd)"]
        F13 --> F14["v0.8.8: API DX, GraphQL & SPIFFE (7 Cmd)"]
        F14 --> V1["🚀 v1.0.0: Grand Parity GA (78 Commands)"]
    end
```

---

## 📦 Detail Cakupan Batch Rilis Mendatang:

### 🏦 Batch 1 — Fase 9: Fintech & Financial Reliability (`v0.8.3`)
- `add:idempotency [redis|memory]`: Middleware Idempotency Key dengan atomic redis lock.
- `add:ledger`: Double-entry bookkeeping ledger engine (Debit/Credit balance invariance).
- `add:decimal`: Fixed-point financial math arithmetic helper (`shopspring/decimal`).
- `add:reconciliation`: Automated settlement & transaction batch reconciliation.
- `add:pricing-engine`: Tiered rate-card & dynamic rule-based pricing calculator.

### 🏛️ Batch 2 — Fase 10: Tactical Domain-Driven Design (`v0.8.4`)
- `make:valueobject [name]`: Generator Immutable Value Object dengan validasi & equality comparator.
- `make:aggregate [name]`: Generator DDD Aggregate Root dengan invariant checking & event dispatcher.
- `make:event [name]`: Generator Domain Event struct & serializer JSON/Protobuf.
- `make:command [name]`: Granular CQRS Command DTO & Command Handler generator.
- `make:query [name]`: Granular CQRS Query DTO & Query Handler generator.

### ⚡ Batch 3 — Fase 11: Concurrency & High-Load Resilience (`v0.8.5`)
- `add:singleflight`: Anti-thundering herd caching wrapper (`golang.org/x/sync/singleflight`).
- `add:workerpool`: Bounded concurrency worker pool engine (`golang.org/x/sync/errgroup`).
- `add:ratelimit [tokenbucket|slidingwindow]`: Distributed & memory rate limiting middleware.
- `add:circuitbreaker [gobreaker]`: Circuit breaker state machine untuk outbound HTTP/gRPC client.
- `add:bulkhead`: Resource isolation wrapper untuk isolasi goroutine dan connection pool.
- `add:cache:l2`: Multi-tier hybrid cache (L1 In-Memory Fast Cache + L2 Redis Distributed Cache).

### 🗄️ Batch 4 — Fase 12: Database Extensions & Distributed Storage (`v0.8.6`)
- `add:db:read-replica`: Master-Replica connection pool dengan dynamic read/write query routing.
- `add:db:pgvector`: Vector similarity search schema migration & query helper (HNSW/IVFFlat).
- `add:db:timescaledb`: Time-series hypertable migrations & analytical rollups.
- `add:kv [badger|bbolt]`: Embedded zero-dependency Key-Value database client.
- `add:raft`: In-process Raft consensus state machine adapter (`hashicorp/raft`).

### 📄 Batch 5 — Fase 13: Media, Export & Document Pipeline (`v0.8.7`)
- `add:media:excel`: Streaming million-row Excel generator/parser (`excelize/v2`).
- `add:media:pdf`: Headless Chromium pixel-perfect HTML-to-PDF engine (`chromedp`).
- `add:media:image`: High-speed JIT image processing & WebP converter (`govips`).
- `add:media:video`: FFmpeg wrapper untuk video transcoding & HLS streaming playlist generator.

### 🌐 Batch 6 — Fase 14: API DX, GraphQL & Zero-Trust (`v0.8.8`)
- `add:versioning [header|uri]`: Semantic API Versioning router (`/v1` atau header `X-API-Version`).
- `add:pagination`: Cursor-based, Keyset & Offset pagination DTOs & SQL builders.
- `add:i18n`: Multilingual translation catalog & `Accept-Language` context extractor.
- `add:graphql`: GraphQL schema-first server boilerplate (`gqlgen` + DataLoader).
- `add:sse`: Server-Sent Events (SSE) stream emitter handler.
- `add:openapi`: OpenAPI 3.1 contract generator & Swagger UI handler.
- `add:spiffe`: SPIFFE / SPIRE zero-trust workload identity mTLS attestation client.

---

### Panduan Siklus Status RFC:
* `PROPOSED`: Sedang dalam tahap perancangan & review (menunggu "Gasskan").
* `ACCEPTED`: Telah disetujui user ("Gasskan"), siap/sedang dieksekusi.
* `IMPLEMENTED`: Seluruh batch task tuntas dan lolos Quality Gate.
* `SUPERSEDED`: Digantikan oleh dokumen RFC yang lebih baru.
