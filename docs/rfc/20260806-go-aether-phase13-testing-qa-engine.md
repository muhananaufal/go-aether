# 🧪 RFC 20260806-15: Phase 13 — Native QA, Testing, Stress & Chaos Engine

> **Inisiatif**: Core QA & Reliability Parity Engine  
> **Versi Target**: `v0.8.7`  
> **Status**: `PROPOSED`  
> **Auditor/Author**: Aetheris L8 Principal Engineer  
> **Target Branch**: `feature/v0.8.7-testing-qa-engine`  
> **Dependencies**: `v0.8.6` (Caching & Cloud Storage Engine)

---

## 🧭 1. First-Principles Rooting (5W1H)

### Why (Mengapa Arsitektur Ini Mutlak Diperlukan?)
Dalam arsitektur backend Go berkinerja tinggi, **klaim "kode kami cepat dan andal" tidak memiliki nilai tanpa bukti empiris (Proof-of-Performance & Proof-of-Resilience)**.
Mayoritas framework scaffolding hanya menyediakan unit test dasar (`testing.T`) yang rentan terhadap:
1. **Tautological Assertions:** Test yang selalu lulus karena hanya menguji hal dangkal tanpa memvalidasi *state mutation* dan edge cases.
2. **Hidden Concurrency Race Conditions:** Bug yang hanya muncul di bawah beban 10.000 RPS atau saat downstream service mengalami degradasi latensi tinggi.
3. **Unfuzzed Boundary Crashes:** Panic runtime yang disebabkan oleh payload biner cacat atau encoding yang lolos dari validasi standar.
4. **Mock Brittleness:** Mock yang dibuat manual rentan usang (*stale*) saat interface berubah.

Phase 13 menghadirkan **Quality & Resilience Suite kelas Tier-1 SRE** langsung ke dalam CLI `go-aether`.

### What (Apa Saja yang Dibangun?)
Tujuh (7) generator pengujian dan ketahanan sistem:
1. `test:stress [k6|vegeta]`: Load testing harness & latency histogram benchmark.
2. `test:chaos`: Failure injection middleware (simulasi latensi, error rate acak, connection abort).
3. `test:fuzz`: Continuous Go native fuzzing harness (`go test -fuzz`) untuk memvalidasi parser/domain invariants.
4. `test:mutation`: Mutation testing runner script untuk mengukur kualitas unit test suite.
5. `test:benchmark`: Micro-benchmark suite + memory allocation & escape analysis profiler (`testing.B`).
6. `test:container`: Testcontainers integration testing harness (PostgreSQL & Redis asli di Docker otomatis).
7. `make:mock [interface]`: Interface mock generator berbasis Mockery untuk isolated unit tests.

### Who & Where (Siapa & Di Mana?)
- **CLI Adapter**: `internal/adapter/cli/test.go` & `internal/adapter/cli/make_mock.go`.
- **Core Service**: `internal/core/service/test_service.go` & `scaffold_service.go`.
- **Templates**: `templates/tests/*.tmpl` dan `templates/plugins/chaos.go.tmpl`.
- **Generated Artifacts**: Folder `tests/stress/`, `tests/fuzz/`, `tests/bench/`, `tests/integration/`, `pkg/middleware/chaos.go`, dan `mocks/`.

---

## ⚔️ 2. Adversarial Anticipation & Trade-Off Matrix

| Opsi Keputusan | Trade-Off Positif (+) | Trade-Off Negatif (-) | Mitigasi / Alasan Terpilih |
| :--- | :--- | :--- | :--- |
| **K6 vs Vegeta untuk `test:stress`** | K6 mendukung scripting JS modern & cloud metrics; Vegeta adalah biner Go murni tanpa dependency. | K6 butuh runtime K6 terinstall; Vegeta skripnya kurang fleksibel untuk skenario bertahap. | **Disediakan flag `--engine=k6` (default) dan `--engine=vegeta`** dengan fallback script mandiri. |
| **Chaos Injection in Middleware vs Service Mesh** | Tidak butuh Istio/Linkerd, bisa jalan di local dev dan integration test. | Sedikit overhead evaluasi header di production jika tidak di-disable. | **Wajib di-guarded oleh environment check `CHAOS_ENABLED=true`** dan security token header. |
| **Testcontainers vs In-Memory SQLite/Miniredis** | 100% kompatibel dengan fitur native Postgres (JSONB, pgvector, CTE, RLS) dan Redis Cluster. | Membutuhkan Docker daemon aktif saat pengujian integrasi dijalankan. | **Auto-detect Docker daemon:** Jika Docker tidak aktif, fallback gracefully dengan pesan instruksi yang jelas. |

---

## 📊 3. Exhaustive Taxonomy & Matrix Scaffolding

```
========================================================================================================
Command              Target Path Output                       Dependencies / Stack       Tujuan Pengujian
========================================================================================================
test:stress          tests/stress/load_test.js                k6 / vegeta                SLA 10k RPS & p99 < 50ms
test:chaos           pkg/middleware/chaos.go                  Go Stdlib / Chi / Gin      Fault Tolerance & Resiliency
test:fuzz            tests/fuzz/fuzz_test.go                  Go 1.18+ Native Fuzz       Parser & Security Crash Detection
test:mutation        scripts/mutation_test.ps1                go-mutesting / harness     Anti-Tautological Test Validation
test:benchmark       tests/bench/alloc_benchmark_test.go      testing.B (0 allocs/op)    Nanosecond & Heap Alloc Profiling
test:container       tests/integration/postgres_redis_test.go testcontainers-go          Real Postgres/Redis E2E Testing
make:mock            mocks/[interface]_mock.go                vektra/mockery             Decoupled Unit Testing
========================================================================================================
```

---

## ⚠️ 4. Zero-Exception Anomaly Hunting & Edge Cases

1. **Docker Daemon Mati saat `test:container` dijalankan:**
   - *Mitigasi:* Test suite menyertakan pre-flight check `testcontainers.CheckDocker()`. Jika Docker offline, test diskip dengan status `t.Skip("Docker daemon not running, skipping container integration tests")` tanpa membuat CI merah secara palsu.
2. **Chaos Injection Bocor ke Production Traffic:**
   - *Mitigasi:* Chaos middleware secara default nonaktif (`CHAOS_ENABLED=false`). Hanya aktif jika header rahasia `X-Chaos-Secret` cocok dengan konfigurasi server.
3. **Fuzz Testing Kehabisan Memori (OOM) di CI Server:**
   - *Mitigasi:* Template `fuzz_test.go` membatasi ukuran byte input buffer (`if len(data) > 65536 { return }`) untuk mencegah memory exhaustion.
4. **Stale Mock Files setelah Interface Dimutasi:**
   - *Mitigasi:* `make:mock` menyertakan build tag dan instruksi generation directive `//go:generate mockery` agar sinkronisasi dapat diperbarui satu perintah `go generate ./...`.

---

## 🗺️ 5. Visual Architecture Blueprint (Mermaid Grounding)

```mermaid
flowchart TD
    subgraph DevEnvironment ["💻 Developer & CI/CD Pipeline"]
        CLI["go-aether CLI Engine"]
    end

    subgraph GeneratedTestingSuites ["🧪 Generated Quality & Resilience Suites"]
        TS["test:stress<br/>(k6 / Vegeta Load Tests)"]
        TC["test:chaos<br/>(pkg/middleware/chaos.go)"]
        TF["test:fuzz<br/>(Continuous Fuzz Harness)"]
        TB["test:benchmark<br/>(testing.B Memory Profiler)"]
        TI["test:container<br/>(Testcontainers E2E in Docker)"]
        TM["make:mock<br/>(Mockery Isolated Interfaces)"]
        TU["test:mutation<br/>(Mutation Assertion Verifier)"]
    end

    subgraph AppLayers ["🏢 Application Layers Under Test"]
        HTTP["HTTP / gRPC Handlers"]
        SVC["Core Domain Services"]
        DB["Real PostgreSQL & Redis Containers"]
    end

    CLI -->|Scaffolds| TS
    CLI -->|Scaffolds| TC
    CLI -->|Scaffolds| TF
    CLI -->|Scaffolds| TB
    CLI -->|Scaffolds| TI
    CLI -->|Scaffolds| TM
    CLI -->|Scaffolds| TU

    TS -->|Load Injection| HTTP
    TC -->|Fault Injection (Latency/Abort)| HTTP
    TF -->|Corrupted Bytes Fuzzing| SVC
    TM -->|Injected Mock Contracts| SVC
    TI -->|Real Integration Testing| DB
```

---

## 📋 6. Batch Execution Plan (DAG Dependencies)

### Batch 1: Testing & Mocking Contracts Definition
- `DependsOn: []`
- Update `internal/core/port/generator.go` & `scaffold_service.go` dengan interface baru untuk 7 generator pengujian.

### Batch 2: High-Performance Test Templates Implementation
- `DependsOn: [Batch 1]`
- Buat templates:
  - `templates/tests/k6_stress.js.tmpl` & `vegeta_stress.tmpl`
  - `templates/plugins/chaos.go.tmpl`
  - `templates/tests/fuzz.go.tmpl`
  - `templates/tests/benchmark.go.tmpl`
  - `templates/tests/testcontainers.go.tmpl`
  - `templates/tests/mutation.tmpl`
  - `templates/core/mock.go.tmpl`

### Batch 3: Core Service & CLI Commands Registration
- `DependsOn: [Batch 2]`
- Implementasikan logic scaffolding di `internal/core/service/test_service.go` dan `add_service.go`.
- Registrasikan subcommands di `internal/adapter/cli/test.go`, `make_mock.go`, dan kaitkan di `root.go`.

### Batch 4: Quality Gate & E2E Validation
- `DependsOn: [Batch 3]`
- Jalankan suite test `go test -race ./...`
- Jalankan E2E scaffolding verification untuk seluruh 7 command baru.
- Lakukan Git commit atomic dan tag release `v0.8.7`.
