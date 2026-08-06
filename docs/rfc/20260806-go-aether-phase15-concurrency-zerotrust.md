# ⚡ RFC 20260806-17: Phase 15 — Concurrency Pipelines, Zero-Trust Security & Observability

> **Inisiatif**: Final Polish & Grand Parity Completion (The Final Sprint)  
> **Versi Target**: `v0.8.9`  
> **Status**: `PROPOSED`  
> **Auditor/Author**: Aetheris L8 Principal Engineer  
> **Target Branch**: `feature/v0.8.9-concurrency-zerotrust`  
> **Dependencies**: `v0.8.8` (Enterprise SQLC & gRPC Engine)

---

## 🧭 1. First-Principles Rooting (5W1H)

### Why (Mengapa Arsitektur Ini Mutlak Diperlukan?)
Menuju rilis stabil `v1.0.0`, `go-aether` harus menutup *gap* terakhir di aspek performa, keamanan ekstrem, dan keandalan tingkat lanjut:
1. **Concurrency Chaos:** Pemrosesan batch rentan terhadap *goroutine leaks* dan OOM. **Fan-Out/Fan-In Pipeline** membatasi konkurensi (bounded concurrency) secara aman.
2. **Thundering Herd / Cache Stampede:** Lonjakan trafik ke satu *key* bersamaan dapat merobohkan database. **Singleflight** (request deduplication) menyerap ribuan request menjadi 1 query fisik.
3. **Graceful Degradation:** Mematikan server secara brutal saat *deployment* akan memotong request in-flight pengguna. **Drain Manager** memastikan *zero-downtime shutdown*.
4. **Data Tampering & Password Breach:** Regulasi mewajibkan logging yang bebas dari *tampering* dan hashing password anti-GPU. **AuditLog** dan **Argon2id** menjawab syarat compliance ketat (PCI-DSS / HIPAA).

### What (Apa Saja yang Dibangun?)
Delapan (8) generator paripurna:
1. `make:pipeline [name]`: Fan-Out/Fan-In bounded concurrency helper.
2. `add:singleflight`: Deduplikasi request via `golang.org/x/sync/singleflight`.
3. `add:metrics [prometheus]`: Prometheus RED metrics (Rate, Errors, Duration) middleware.
4. `add:drain`: Graceful shutdown & connection drain manager.
5. `add:oauth2 [provider]`: OIDC/OAuth2 login client (Google/GitHub) via `golang.org/x/oauth2` dengan PKCE.
6. `add:auditlog`: Middleware audit log immutability & PII scrubbing otomatis.
7. `add:argon2`: `golang.org/x/crypto/argon2` hasher helper (anti ASIC/GPU cracking).
8. `make:specification`: DDD Specification pattern (Domain Query Rules).

### Who & Where (Siapa & Di Mana?)
- **CLI Adapter**: `internal/adapter/cli/make_pipeline.go`, `add_security.go`, `add_metrics.go`, dll.
- **Core Service**: 8 *methods* tambahan di `internal/core/service/`.
- **Templates**: `templates/plugins/...`

---

## ⚔️ 2. Adversarial Anticipation & Trade-Off Matrix

| Opsi Keputusan | Trade-Off Positif (+) | Trade-Off Negatif (-) | Mitigasi / Alasan Terpilih |
| :--- | :--- | :--- | :--- |
| **Pipeline Channels vs WaitGroup** | Backpressure support, mudah dikontrol limit pool-nya, fail-fast on context cancel. | Potensi *deadlock* jika sender tidak menutup channel (unbuffered block). | **Pipeline generator** selalu di-*scaffold* dengan pola *generator -> worker -> multiplexer* yang secara deterministik akan men-defer `close(ch)`. |
| **Argon2id vs Bcrypt** | Standar keamanan terbaru (OWASP 2026), memori-hard anti GPU/ASIC crackers. | CPU & Memory cost lebih tinggi saat hashing (time-consuming). | Parameter `time`, `memory`, dan `threads` akan disetel moderat (64MB, 1 iterasi, 4 threads) sebagai *baseline* yang aman namun performant. |
| **Singleflight vs Distributed Lock** | Sangat cepat (in-memory) dan tidak butuh Redis. | Hanya mem-blok request di *node* yang sama. Multi-node tetap tembus N query. | Ideal untuk mencegah *local cache stampede* pada DB. Digunakan komplementer dengan L1/L2 Redis Cache yang di-scaffold di Phase 12. |

---

## 📊 3. Exhaustive Taxonomy & Matrix Scaffolding

```
========================================================================================================
Command              Target Path Output                       Dependencies / Stack       Tujuan Enterprise
========================================================================================================
make:pipeline        pkg/concurrency/[name]_pipeline.go       Native Channels / Sync     Bounded Batch Processing
add:singleflight     pkg/concurrency/singleflight.go          golang.org/x/sync          Anti Thundering Herd
add:metrics          pkg/middleware/prometheus.go             github.com/prometheus/..   RED Monitoring Metrics
add:drain            pkg/server/drain.go                      os/signal & context        Zero-Downtime Shutdown
add:oauth2           internal/adapter/handler/http/oauth2.go  golang.org/x/oauth2        Social SSO + PKCE
add:auditlog         pkg/middleware/auditlog.go               slog & hashing             Tamper-Evident Logs
add:argon2           pkg/security/argon2.go                   golang.org/x/crypto/argon2 GPU-resistant Password Hash
make:specification   internal/core/domain/[spec].go           Native Interfaces          Reusable DDD Rules
========================================================================================================
```

---

## ⚠️ 4. Zero-Exception Anomaly Hunting & Edge Cases

1. **Pipeline Deadlocks pada Error/Panic:**
   - *Mitigasi:* Pipeline template secara eksplisit akan menyisipkan `defer func() { if r := recover(); ... }()` pada setiap worker goroutine dan akan melakukan *fail-fast* melalui errgroup context cancellation agar sender berhenti memblok channel.
2. **PII Data Leakage di Audit Log:**
   - *Mitigasi:* `add:auditlog` akan membangkitkan middleware yang menggunakan regex/masking rules (`x-api-key`, `password`, `ssn`, `credit_card`) sebelum data body/header di-*dump* ke logs.
3. **Shutdown Timeout (Zombie Processes):**
   - *Mitigasi:* `add:drain` menetapkan `context.WithTimeout(..., 30 * time.Second)` yang ketat. Jika koneksi lama (misal: websocket) tidak selesai dalam 30 detik, `os.Exit(1)` akan dieksekusi brutal demi mencegah proses tersangkut selamanya.

---

## 📋 5. Batch Execution Plan (DAG Dependencies)

### Batch 1: Final Interface Contracts Definition
- `DependsOn: []`
- Update `internal/core/port/generator.go` dengan 8 *methods* baru.

### Batch 2: The 8 Grand Templates Implementation
- `DependsOn: [Batch 1]`
- Buat 8 template di `templates/plugins/` (pipeline, singleflight, metrics, drain, oauth2, auditlog, argon2, specification).

### Batch 3: Service Layer & CLI Registration
- `DependsOn: [Batch 2]`
- Implementasikan logic scaffolding di `internal/core/service/add_service.go` dan `make_service.go`.
- Buat file perintah CLI baru `internal/adapter/cli/add_phase15.go` (untuk `add:*`) dan modifikasi `make.go` (untuk `make:*`).
- Daftarkan semuanya di `root.go`.

### Batch 4: QA, Verification & The Grand V1 GA Prep
- `DependsOn: [Batch 3]`
- `go test -v ./...` (Wajib 100% Pass).
- *End-to-End Dry-Run* semua 8 commands.
- Merge ke `main`, tag `v0.8.9`.
