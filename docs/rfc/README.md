# 📚 RFC Architecture Catalog & Master Index (`go-aether`)

Dokumentasi rancangan arsitektur, keputusan sistem (ADR), dan riwayat implementasi teknis untuk **`go-aether` — Opinionated Architecture Scaffold CLI Engine for Go Backend Engineers**.

| RFC ID | Inisiatif Fitur / Arsitektur | Domain / Modul | Target Branch | Status | Tanggal Rilis |
| :--- | :--- | :--- | :--- | :---: | :---: |
| `20260806-01` | [Core Engine Scaffolder, CLI Framework & Manifest SSOT](20260806-go-aether-core-scaffolding-engine.md) | Core CLI / Scaffold Engine | `feature/core-engine` | `IMPLEMENTED` | 2026-08-06 |
| `20260806-02` | [Phase 2: Granular Ecosystem & Middleware Expansion](20260806-go-aether-phase2-ecosystem.md) | Ecosystem Addons | `feature/phase2-ecosystem` | `IMPLEMENTED` | 2026-08-06 |
| `20260806-03` | [Phase 3: Distributed Systems & Observability](20260806-go-aether-phase3-distributed.md) | Distributed & OTel | `feature/phase3-distributed` | `IMPLEMENTED` | 2026-08-06 |
| `20260806-04` | [Phase 4: Cloud-Native, CI/CD & AI](20260806-go-aether-phase4-cloudnative.md) | Cloud & CLI TUI | `feature/phase4-cloudnative` | `IMPLEMENTED` | 2026-08-06 |
| `20260806-05` | [Phase 5: Core Enterprise Patterns](20260806-go-aether-phase5-core-patterns.md) | Architectural Patterns | `feature/phase5-core-patterns` | `IMPLEMENTED` | 2026-08-06 |
| `20260806-06` | [Phase 6: Data, Migration & QA](20260806-go-aether-phase6-data-qa.md) | Database & Test | `feature/phase6-data-qa` | `IMPLEMENTED` | 2026-08-06 |

### 🗺️ ROADMAP: The 8 Phases of Evolution
Proyek `go-aether` dibagi menjadi 8 tonggak pencapaian besar (Target Production Release di `v1.0.0`):
- **Fase 1: The Core Engine (v0.1.0)** — Mesin CLI dasar, sistem penulisan transaksional, YAML resolver, dan pembentukan struktur murni Hexagonal (Domain ke Repository). *(SELESAI)*
- **Fase 2: Granular Ecosystem & Middleware (v0.2.0)** — Membangun generator `make:*` parsial dan injeksi `add:middleware` (JWT, Rate Limit, Auth, dsb) untuk operasional standar industri. *(SELESAI)*
- **Fase 3: Distributed Systems & Observability (v0.3.0)** — Membangun `add:worker` (Kafka/Redis), Sagas (Temporal), Metrics (Prometheus), dan arsitektur *High-Frequency*. *(SELESAI)*
- **Fase 4: Cloud-Native, CI/CD & AI (v0.4.0)** — Otomasi `add:cicd`, `add:deploy` (K8s/Helm), integrasi `add:ai` (LLM), dan sihir adopsi kode lama (`adopt --scan` dengan AST parsing). *(SELESAI)*
- **Fase 5: Core Enterprise Patterns (v0.5.0)** — Ekspansi ke `make:repository`, `make:domain`, `add:di`, `add:config`, dan standar *error/validation*. *(SELESAI)*
- **Fase 6: Data, Migration & QA (v0.6.0)** — Database migration (`make:migration`), *seeder*, tes komprehensif (`add:test`), dan multitenancy. *(SELESAI)*
- **Fase 7: Advanced Distributed Patterns (v0.7.0)** — Injeksi CQRS, Saga/Temporal, Outbox, Webhook, dan Service Discovery. *(TERJADWAL)*
- **Fase 8: Cloud, Auth & 3rd-Party Plugins (v0.8.0)** — Otentikasi Oauth2, S3 Storage, Cron, Mailer, Firebase, dan Logger. *(TERJADWAL)*

### Panduan Siklus Status RFC:
* `PROPOSED`: Sedang dalam tahap perancangan & review (menunggu "Gasskan").
* `ACCEPTED`: Telah disetujui user ("Gasskan"), siap/sedang dieksekusi.
* `IMPLEMENTED`: Seluruh batch task tuntas dan lolos Quality Gate.
* `SUPERSEDED`: Digantikan oleh dokumen RFC yang lebih baru (sertakan link ke RFC pengganti).
