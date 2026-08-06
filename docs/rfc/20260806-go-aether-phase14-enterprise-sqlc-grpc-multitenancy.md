# 🏛️ RFC 20260806-16: Phase 14 — Enterprise SQLC, Advanced gRPC Streaming, Gateway & Multi-Tenancy Engine

> **Inisiatif**: Core Enterprise Data & High-Performance RPC Engine  
> **Versi Target**: `v0.8.8`  
> **Status**: `PROPOSED`  
> **Auditor/Author**: Aetheris L8 Principal Engineer  
> **Target Branch**: `feature/v0.8.8-sqlc-grpc-multitenancy`  
> **Dependencies**: `v0.8.7` (Native QA, Testing & Chaos Engine)

---

## 🧭 1. First-Principles Rooting (5W1H)

### Why (Mengapa Arsitektur Ini Mutlak Diperlukan?)
Pada skala enterprise, tiga tantangan data & komunikasi backend selalu muncul:
1. **Runtime Query Fragility (ORM/Raw SQL):** ORM sering menyebabkan N+1 queries, sedangkan raw SQL rentan typo yang baru meledak saat runtime. **SQLC** menyelesaikan masalah ini dengan meng-compile query SQL murni menjadi kode Go bertipe aman (*zero-reflection, zero runtime overhead, 100% type-safe*).
2. **High-Frequency RPC & Real-Time Sync:** REST HTTP/1.1 memiliki overhead header dan multiplexing buruk. **gRPC Bidirectional Streaming** memungkinkan komunikasi dupleks kontinu berkecepatan tinggi antar microservices.
3. **Dual-Protocol REST/gRPC Maintenance:** Menulis handler HTTP terpisah dari handler gRPC membuang waktu. **gRPC-Gateway** membangkitkan proxy HTTP/REST JSON otomatis langsung dari kontrak `.proto`.
4. **Tenant Data Leakage:** Dalam SaaS multi-tenant, kebocoran data antar tenant adalah risiko hukum fatal. **Tenant Isolation Context** menyuntikkan tenant identifier secara ketat ke seluruh DB queries dan context metadata.

### What (Apa Saja yang Dibangun?)
Empat (4) generator enterprise core:
1. `add:sqlc`: Men-generate konfigurasi `sqlc.yaml`, skema dasar `schema.sql`, query dasar `query.sql`, dan Makefile runner directive.
2. `add:grpc-stream`: Men-generate gRPC bidirectional & client/server streaming handler serta kontrak `.proto` komplit.
3. `add:grpc-gateway`: Men-generate reverse-proxy gRPC-Gateway untuk mengekspos endpoint gRPC sebagai REST/JSON HTTP endpoint.
4. `add:tenant-context`: Men-generate multi-tenancy middleware, tenant context extractor, dan database scoping helper.

### Who & Where (Siapa & Di Mana?)
- **CLI Adapter**: `internal/adapter/cli/add.go` & `internal/adapter/cli/root.go`.
- **Core Service**: `internal/core/service/add_service.go` & `generator.go`.
- **Templates**:
  - `templates/plugins/sqlc_yaml.tmpl`
  - `templates/plugins/sqlc_schema.sql.tmpl`
  - `templates/plugins/sqlc_query.sql.tmpl`
  - `templates/plugins/grpc_stream.go.tmpl`
  - `templates/plugins/grpc_stream.proto.tmpl`
  - `templates/plugins/grpc_gateway.go.tmpl`
  - `templates/plugins/tenant_context.go.tmpl`
- **Generated Artifacts**: `sqlc.yaml`, `db/schema/`, `db/queries/`, `internal/adapter/handler/grpc/`, `pkg/middleware/tenant.go`, `pkg/tenant/`.

---

## ⚔️ 2. Adversarial Anticipation & Trade-Off Matrix

| Opsi Keputusan | Trade-Off Positif (+) | Trade-Off Negatif (-) | Mitigasi / Alasan Terpilih |
| :--- | :--- | :--- | :--- |
| **SQLC vs GORM/Ent** | Zero runtime reflection, 100% compile-time SQL validation, eksekusi secepat raw SQL. | Membutuhkan binary SQLC terinstall untuk regenerasi kode Go. | **Menyediakan fallback generator template dan Makefile task `make sqlc`** agar mudah dieksekusi di local dan CI/CD. |
| **gRPC-Gateway vs Manual Dual Handlers** | Single Source of Truth (Protobuf `.proto`), otomatis sinkron antara REST dan gRPC. | Butuh protoc plugin `protoc-gen-grpc-gateway`. | **Menyediakan standalone HTTP handler wrapper** yang langsung mengalirkan request ke gRPC server in-process atau via TCP. |
| **Tenant by Column (RLS) vs Tenant by Database/Schema** | Paling efisien secara sumber daya, mudah di-scale ke ribuan tenant tanpa overhead multi-pool. | Memerlukan klausul `WHERE tenant_id = ?` di setiap query. | **Tenant Context Middleware + Helper otomatis menyuntikkan tenant_id** dan memvalidasi keberadaan tenant di context. |

---

## 📊 3. Exhaustive Taxonomy & Matrix Scaffolding

```
========================================================================================================
Command              Target Path Output                       Dependencies / Stack       Tujuan Enterprise
========================================================================================================
add:sqlc             sqlc.yaml, db/schema/, db/queries/       sqlc-dev / sqlc            Compile-Time Type-Safe SQL
add:grpc-stream      internal/adapter/handler/grpc/stream.go  google.golang.org/grpc     Duplex Streaming RPC
add:grpc-gateway     internal/adapter/handler/http/gateway.go grpc-ecosystem/grpc-gateway REST/JSON HTTP-to-gRPC Bridge
add:tenant-context   pkg/middleware/tenant.go, pkg/tenant/    Go Context / Chi           Tenant Scoping & Data Isolation
========================================================================================================
```

---

## ⚠️ 4. Zero-Exception Anomaly Hunting & Edge Cases

1. **SQLC Tidak Terinstall di Komputer Developer:**
   - *Mitigasi:* Template `sqlc.yaml` dan query di-scaffold lengkap dengan instruksi cara install (`go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`) dan Makefile command.
2. **gRPC Stream Goroutine Leak saat Client Disconnect:**
   - *Mitigasi:* Handler gRPC stream menggunakan `ctx.Done()` listener di goroutine terpisah untuk memastikan seluruh background channels segera ditutup saat client terputus.
3. **Cross-Tenant Context Pollution:**
   - *Mitigasi:* Middleware `pkg/middleware/tenant.go` secara ketat meng-abort request (HTTP 401 / 403) jika header `X-Tenant-ID` hilang atau formatnya tidak valid (anti empty-tenant bypass).
4. **gRPC-Gateway Error Mapping ke HTTP Status Codes:**
   - *Mitigasi:* Gateway handler menyertakan custom error handler yang memetakan gRPC status codes (`codes.NotFound` -> `404`, `codes.InvalidArgument` -> `400`, `codes.Unauthenticated` -> `401`) secara konsisten.

---

## 🗺️ 5. Visual Architecture Blueprint (Mermaid Grounding)

```mermaid
flowchart TD
    subgraph Clients ["🌐 External Clients & Microservices"]
        RESTClient["HTTP/REST Frontend Clients"]
        GRPCClient["Internal gRPC Microservices"]
    end

    subgraph TransportLayer ["🚪 Transport & Gateway Layer"]
        TenantMW["pkg/middleware/tenant.go<br/>(Extracts & Validates X-Tenant-ID)"]
        Gateway["internal/adapter/handler/http/gateway.go<br/>(gRPC-Gateway Reverse Proxy)"]
        StreamHandler["internal/adapter/handler/grpc/stream.go<br/>(Bi-Directional gRPC Stream)"]
    end

    subgraph ServiceLayer ["⚙️ Core Domain & Service Layer"]
        TenantCtx["pkg/tenant/context.go<br/>(Tenant Context Binding)"]
        DomainSvc["Core Business Logic Services"]
    end

    subgraph DataLayer ["🗄️ Persistence & Type-Safe SQLC Layer"]
        SQLCQueries["db/queries/ & db/schema/<br/>(Type-Safe Generated Queries)"]
        PostgresDB["Multi-Tenant PostgreSQL (Scoped by TenantID)"]
    end

    RESTClient -->|HTTP/JSON Request| TenantMW
    TenantMW --> Gateway
    Gateway -->|In-Process gRPC Call| StreamHandler
    GRPCClient -->|Protobuf RPC Stream| StreamHandler

    StreamHandler --> TenantCtx
    TenantCtx --> DomainSvc
    DomainSvc --> SQLCQueries
    SQLCQueries -->|Scoped SQL (WHERE tenant_id = $1)| PostgresDB
```

---

## 📋 6. Batch Execution Plan (DAG Dependencies)

### Batch 1: Enterprise Port Contracts Definition
- `DependsOn: []`
- Update `internal/core/port/generator.go` dengan interface `AddSQLC`, `AddGRPCStream`, `AddGRPCGateway`, `AddTenantContext`.

### Batch 2: Enterprise Templates Implementation
- `DependsOn: [Batch 1]`
- Buat templates:
  - `templates/plugins/sqlc_yaml.tmpl`
  - `templates/plugins/sqlc_schema.sql.tmpl`
  - `templates/plugins/sqlc_query.sql.tmpl`
  - `templates/plugins/grpc_stream.go.tmpl`
  - `templates/plugins/grpc_stream.proto.tmpl`
  - `templates/plugins/grpc_gateway.go.tmpl`
  - `templates/plugins/tenant_context.go.tmpl`

### Batch 3: Service Layer & CLI Registration
- `DependsOn: [Batch 2]`
- Implementasikan 4 methods di `internal/core/service/add_service.go`.
- Registrasikan subcommands di `internal/adapter/cli/add.go` dan kaitkan di `root.go`.

### Batch 4: Quality Gate & E2E Validation
- `DependsOn: [Batch 3]`
- Jalankan unit tests `go test -v ./...`
- Jalankan E2E validation.
- Lakukan Git commit atomic, merge ke `main`, tag `v0.8.8`, dan push.
