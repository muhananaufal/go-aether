# 📚 RFC Architecture Catalog & Master Index (`go-aether`)

Dokumentasi rancangan arsitektur, keputusan sistem (ADR), dan riwayat implementasi teknis untuk **`go-aether` — Opinionated Architecture Scaffold CLI Engine for Go Backend Engineers**.

| RFC ID | Inisiatif Fitur / Arsitektur | Domain / Modul | Target Branch | Status | Tanggal Rilis |
| :--- | :--- | :--- | :--- | :---: | :---: |
| `20260806-01` | [Core Engine Scaffolder, CLI Framework & Manifest SSOT](20260806-go-aether-core-scaffolding-engine.md) | Core CLI / Scaffold Engine | `feature/core-engine` | `IMPLEMENTED` | 2026-08-06 |
| `20260806-02` | [Phase 2: Granular Ecosystem & Middleware Expansion](20260806-go-aether-phase2-ecosystem.md) | Ecosystem Addons | `feature/phase2-ecosystem` | `IMPLEMENTED` | 2026-08-06 |
| `20260806-03` | [Phase 3: Distributed Systems & Observability](20260806-go-aether-phase3-distributed.md) | Distributed & OTel | `feature/phase3-distributed` | `IMPLEMENTED` | 2026-08-06 |
| `20260806-04` | [Phase 4: Cloud-Native, CI/CD & AI](20260806-go-aether-phase4-cloudnative.md) | Cloud & CLI TUI | `feature/phase4-cloudnative` | `PROPOSED` | - |

### 🗺️ ROADMAP: The 4 Phases of Evolution
Proyek `go-aether` dibagi menjadi 4 tonggak pencapaian besar:
- **Fase 1: The Core Engine (v1.0.0)** — Mesin CLI dasar, sistem penulisan transaksional, YAML resolver, dan pembentukan struktur murni Hexagonal (Domain ke Repository). *(SELESAI)*
- **Fase 2: Granular Ecosystem & Middleware** — Membangun generator `make:*` parsial dan injeksi `add:middleware` (JWT, Rate Limit, Auth, dsb) untuk operasional standar industri.
- **Fase 3: Distributed Systems & Observability** — Membangun `add:worker` (Kafka/Redis), Sagas (Temporal), Metrics (Prometheus), dan arsitektur *High-Frequency*.
- **Fase 4: Cloud-Native, CI/CD & AI** — Otomasi `add:cicd`, `add:deploy` (K8s/Helm), integrasi `add:ai` (LLM), dan sihir adopsi kode lama (`adopt --scan` dengan AST parsing).

### Panduan Siklus Status RFC:
* `PROPOSED`: Sedang dalam tahap perancangan & review (menunggu "Gasskan").
* `ACCEPTED`: Telah disetujui user ("Gasskan"), siap/sedang dieksekusi.
* `IMPLEMENTED`: Seluruh batch task tuntas dan lolos Quality Gate.
* `SUPERSEDED`: Digantikan oleh dokumen RFC yang lebih baru (sertakan link ke RFC pengganti).
