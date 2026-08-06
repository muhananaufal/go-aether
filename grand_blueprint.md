# `go-aether` — Grand Blueprint
### *Opinionated Architecture Scaffold CLI Engine for Go Backend Engineers*

> **Status:** Pre-RFC (Brainstorming Finalized) — Menunggu "Gasskan"  
> **Versi Dokumen:** 2.0.0 *(Updated: Audit Komprehensif 221 File)*  
> **Ditulis oleh:** AETHERIS + muhananaufal  
> **Target:** Eskalasi Karir → Principal Backend Golang Engineer (L7/L8)

---

## BAB 0 — TL;DR (Executive Summary)

`go-aether` adalah sebuah **CLI code generator** yang ditulis dalam bahasa Go itu sendiri. Saat dijalankan, ia menghasilkan struktur proyek dan file `.go` berstandar arsitektur enterprise (Hexagonal, Clean, DDD) tanpa menambahkan satu baris kode runtime ke dalam binary produksi. Dengan kata lain: hasil generatenya adalah **pure idiomatic Go** — tidak ada dependency tersembunyi, tidak ada magic, tidak ada overhead.

---

## BAB 1 — 5W + 1H (Fondasi Filosofis)

### WHY (Mengapa Ini Harus Dibangun?)

**Problem yang menyebabkan rasa sakit:**
1. Junior Go developer tahu syntax Go, tapi tidak tahu bagaimana menyusun project yang scalable. Mereka googling, mendapat 50 tutorial berbeda dengan 50 struktur folder berbeda. Hasilnya: bingung, lalu copy-paste spaghetti.
2. Official CLI untuk Go tidak ada. `go mod init` hanya menginisiasi `go.mod` — bukan struktur arsitektur.
3. Brownfield (proyek legacy) adalah kenyataan 90% pekerja backend — tidak ada tool yang membantu menyisipkan kode berstandardisasi ke dalamnya tanpa merusak sistem lama.
4. Best practice Go (Constructor Injection, Sentinel Errors, Table-driven Tests, slog OTel) tersebar di ribuan artikel. Tool ini mengkonsolidasikannya dalam satu perintah.

**Apa yang terjadi jika tool ini tidak ada:**  
Developer junior menulis kode Go yang berantakan → dapat penolakan di code review berulang → frustrasi → kabur ke framework yang lebih "magic" → hilangnya talenta Go berkualitas di industri.

---

### WHO (Untuk Siapa?)

**Target Primer: Junior-to-Mid Backend Go Developer (0–3 tahun pengalaman Go)**
- Baru beralih dari PHP/Laravel atau TypeScript/NestJS ke Go
- Ingin belajar best practice tanpa membaca 10 buku terlebih dahulu
- Membutuhkan "guardrail arsitektur" agar kode mereka tidak liar

**Target Sekunder: Mid-to-Senior Engineer yang Memimpin Tim**
- Ingin menegakkan standar arsitektur konsisten di seluruh tim
- Butuh cara cepat untuk onboard anggota tim baru ke konvensi yang sudah disepakati
- Bekerja di proyek brownfield yang butuh modernisasi bertahap

**Yang BUKAN target:**
- Engineer Go berpengalaman >5 tahun yang sudah punya opini arsitektur sendiri
- Tim yang butuh solusi yang sangat spesifik untuk industri tertentu (HFT, K8s Operator) — ini cakupan Niche yang di-support via flags, bukan default

---

### WHAT (Apa Persisnya?)

`go-aether` adalah **CLI binary** yang di-install sekali di mesin developer:

```
go install github.com/muhananaufal/go-aether@latest
```

Setelah terinstall, ia menjadi "arsitek virtual" yang selalu standby:

```bash
# Greenfield: Proyek baru dari nol
go-aether init payment-service --arch hexagonal --db postgres --router chi

# Brownfield: Adopsi proyek legacy
go-aether adopt --scan

# Generate komponen (make:*)
go-aether make:module order
go-aether make:module payment --pattern cqrs
go-aether make:module shipment --pattern outbox
go-aether make:service invoice
go-aether make:handler user --transport grpc --stream server
go-aether make:repository product --cache redis
go-aether make:migration add_orders_table

# Tambah kapabilitas (add:*)
go-aether add:middleware jwt-auth
go-aether add:middleware rate-limit
go-aether add:middleware bulkhead
go-aether add:middleware idempotency
go-aether add:middleware audit-log
go-aether add:transport grpc
go-aether add:transport graphql
go-aether add:transport mqtt
go-aether add:transport connect-rpc
go-aether add:worker email-sender --broker redis --pattern asynq
go-aether add:worker invoice-generator --broker kafka --pattern watermill
go-aether add:cron daily-report
go-aether add:workflow payment-saga --engine temporal
go-aether add:config
go-aether add:di wire
go-aether add:eventing
go-aether add:tracing
go-aether add:profiling
go-aether add:metrics prometheus
go-aether add:healthcheck
go-aether add:versioning
go-aether add:webhook
go-aether add:feature-flags
go-aether add:db read-replica
go-aether add:db pgvector
go-aether add:cache l2
go-aether add:auth argon2
go-aether add:authz casbin
go-aether add:crypto aes-gcm
go-aether add:security mtls
go-aether add:secrets vault
go-aether add:lock redis
go-aether add:discovery etcd
go-aether add:cloud aws
go-aether add:cloud gcp
go-aether add:storage s3
go-aether add:deploy k8s
go-aether add:deploy lambda
go-aether add:deploy helm
go-aether add:multitenancy --strategy rls
go-aether add:ai embedding
go-aether add:ai llm-proxy
go-aether add:test e2e
go-aether add:test fuzz
go-aether add:test property
go-aether add:test synctest
go-aether add:validation
go-aether add:lint custom
go-aether add:middleware cdn-cache

# Utilitas
go-aether doctor           # Validasi konsistensi aether.yaml
go-aether ls               # List semua modul yang sudah di-generate
go-aether upgrade          # Update aether.yaml ke versi spec terbaru
```

---

### WHEN (Kapan Digunakan?)

| Skenario | Command |
|:---|:---|
| Hari pertama proyek baru | `go-aether init` |
| Memulai fitur baru (user story baru) | `go-aether make:module` |
| Butuh tambah lapisan transport (gRPC, WS) | `go-aether add:transport` |
| Bergabung ke proyek legacy baru | `go-aether adopt` |
| Butuh tambah worker/consumer event | `go-aether add:worker` |
| Setiap kali curiga ada inkonsistensi | `go-aether doctor` |

---

### WHERE (Di Mana Tool Ini Hidup?)

- **Diinstall di mesin developer lokal** — bukan di server produksi
- **Kodenya di GitHub** (open source, MIT License): `github.com/muhananaufal/go-aether`
- **Distribusinya via** `go install`, GitHub Releases (binary pre-built), dan (tahap lanjut) Homebrew + Scoop
- **aether.yaml** tersimpan di root direktori project (masuk ke version control, di-commit bersama kode)

---

### HOW (Bagaimana Cara Kerjanya Secara Teknis?)

Mekanisme internal go-aether bekerja dalam 5 langkah berurutan saat setiap command dieksekusi:

```
[1. Command Parser]
     cobra parses subcommand + flags
          ↓
[2. Context Resolver]
     Baca aether.yaml (walk up directory tree)
     Validasi skema YAML
     Resolve architecture pattern yang aktif
          ↓
[3. Template Engine]
     Load template .tmpl dari embed.FS (bundled ke binary)
     Inject template variables (ModuleName, PackageName, dll)
     Render ke string buffer
          ↓
[4. Conflict Resolver]
     Cek apakah file tujuan sudah ada
     Pilih strategi: skip / overwrite / merge (interactive)
     Dry-run mode: tampilkan diff tanpa menulis
          ↓
[5. Writer & Formatter]
     Tulis file ke disk
     Jalankan `goimports` / `gofmt` otomatis
     Tampilkan ringkasan: file created, file skipped, next steps
```

---

## BAB 2 — Arsitektur Internal `go-aether` Sendiri

### Tech Stack CLI:

| Komponen | Library | Alasan |
|:---|:---|:---|
| **CLI Framework** | `cobra` + `viper` | Standar industri Go CLI, dipakai oleh kubectl, docker CLI, git-extras |
| **Interactive TUI** | `bubbletea` + `lipgloss` | Kuesioner interaktif yang cantik (untuk `adopt --scan`) |
| **Config** | `viper` (yaml) | Baca `aether.yaml` dengan type-safe binding |
| **Template Engine** | `text/template` (stdlib) | Zero dependency eksternal, sudah sangat powerful |
| **File Embedding** | `embed.FS` (stdlib Go 1.16+) | Template di-bundle ke dalam satu binary, offline 100% |
| **Filesystem Abstraction** | `afero` | Memungkinkan testing tanpa menyentuh disk nyata (in-memory FS) |
| **Testing** | `testify` + `afero` | Assertion yang jelas + mock filesystem |
| **Logging (internal)** | `slog` (stdlib Go 1.21+) | Logging terstruktur untuk debug mode (`--verbose`) |
| **Spinner/Progress** | `briandowns/spinner` | UX loading yang bersih saat scanning |
| **Diff Display** | `pmezard/go-difflib` | Tampilkan diff sebelum overwrite (--dry-run) |

### Struktur Folder `go-aether` Itu Sendiri (Makan masakan sendiri / "Dogfooding"):

```
go-aether/                        ← Proyek ini sendiri pakai Hexagonal!
├── cmd/
│   └── go-aether/
│       └── main.go              ← Entrypoint
├── internal/
│   ├── core/
│   │   ├── domain/
│   │   │   ├── manifest.go      ← Struct AetherManifest (aether.yaml)
│   │   │   ├── module.go        ← Struct Module, Arch, Transport
│   │   │   └── errors.go        ← Sentinel errors
│   │   ├── port/
│   │   │   ├── generator.go     ← Interface Generator
│   │   │   ├── resolver.go      ← Interface ManifestResolver
│   │   │   └── writer.go        ← Interface FileWriter
│   │   └── service/
│   │       ├── init_service.go
│   │       ├── adopt_service.go
│   │       ├── make_service.go
│   │       └── doctor_service.go
│   └── adapter/
│       ├── template/
│       │   └── engine.go        ← text/template implementation
│       ├── manifest/
│       │   └── yaml_resolver.go ← viper-based aether.yaml reader
│       └── writer/
│           └── afero_writer.go  ← Filesystem writer (real + mock)
├── templates/                   ← Semua .tmpl files (di-embed ke binary)
│   ├── hexagonal/
│   │   ├── domain.go.tmpl
│   │   ├── port.go.tmpl
│   │   ├── service.go.tmpl
│   │   ├── handler_http.go.tmpl
│   │   ├── repository_postgres.go.tmpl
│   │   └── repository_redis.go.tmpl
│   ├── clean/
│   │   ├── entity.go.tmpl
│   │   ├── usecase.go.tmpl
│   │   └── ...
│   ├── ddd/
│   │   └── ...
│   └── common/
│       ├── main.go.tmpl
│       ├── config.go.tmpl
│       ├── dockerfile.tmpl
│       ├── makefile.tmpl
│       └── aether_yaml.tmpl
├── aether.yaml                   ← go-aether sendiri pakai go-aether!
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

---

## BAB 3 — Spesifikasi `aether.yaml` (Schema Lengkap)

`aether.yaml` adalah **Single Source of Truth** (SSOT) untuk identitas arsitektur proyek. Ia adalah "otak eksternal" yang memungkinkan semua subcommand bekerja dengan sadar konteks.

```yaml
# aether.yaml — go-aether Project Manifest
# VERSION: 1.0.0
# JANGAN EDIT MANUAL kecuali kamu tahu apa yang kamu lakukan
# Selalu commit file ini bersama kode

version: "1"

project:
  name: "payment-service"
  module: "github.com/company/payment-service"  # dari go.mod
  go_version: "1.23"                             # dari go.mod
  created_at: "2026-08-06T19:00:00+07:00"
  aether_version: "0.1.0"                         # versi go-aether yang membuat ini

architecture:
  pattern: "hexagonal"           # hexagonal | clean | ddd | vertical | mvc
  mode: "greenfield"             # greenfield | brownfield
  
  # Pemetaan path - bisa kustom untuk brownfield
  paths:
    domain:        "internal/core/domain"
    port:          "internal/core/port"
    service:       "internal/core/service"
    handler_http:  "internal/adapter/handler/http"
    handler_grpc:  "internal/adapter/handler/grpc"     # opsional
    repository:    "internal/adapter/repository"
    cmd:           "cmd/server"
    config:        "internal/config"
    pkg:           "pkg"

stack:
  router:    "chi"               # chi | gin | echo | stdlib
  database:
    driver:  "postgres"          # postgres | mysql | sqlite | mongodb | none
    orm:     "sqlc"              # sqlc | gorm | sqlx | raw | none
    pool:    "pgxpool"           # pgxpool | database/sql
  cache:     "redis"             # redis | memcached | none
  transport:
    - "http"                     # http | grpc | websocket | nats | kafka
  auth:      "jwt-rs256"         # jwt-hs256 | jwt-rs256 | oauth2 | none
  logger:    "slog-otel"         # slog-otel | zerolog | zap | slog
  
adapters:
  # Koneksi yang sudah ada (untuk brownfield - BYO Dependency)
  # Jika non-empty, go-aether TIDAK generate kode koneksi baru
  # melainkan menerima koneksi dari luar via constructor
  existing_db_var:    ""         # misal: "app.DB" (untuk brownfield)
  existing_redis_var: ""         # misal: "app.RedisClient"
  existing_logger_var: ""        # misal: "app.Logger"

modules:
  # Daftar semua modul yang sudah di-generate (dikelola otomatis)
  - name: "order"
    created_at: "2026-08-06T20:00:00+07:00"
    transports: ["http"]
    has_cache: true
    has_worker: false
  - name: "payment"
    created_at: "2026-08-06T21:00:00+07:00"
    transports: ["http", "grpc"]
    has_cache: false
    has_worker: true

workers:
  - name: "email-sender"
    broker: "redis"              # redis | kafka | nats | rabbitmq
    pattern: "asynq"             # asynq | watermill | machinery

migrations:
  tool: "goose"                  # goose | migrate | atlas | none
  dir:  "migrations/"

lint:
  enabled: true
  config: ".golangci.yml"

test:
  coverage_threshold: 80         # minimum coverage % yang diizinkan

meta:
  # Metadata brownfield anomaly mapping
  anomaly_mode: false
  legacy_notes: ""               # catatan bebas untuk tim
```

---

## BAB 4 — Spesifikasi Command Lengkap

> **Total setelah audit komprehensif 221 file:** ~65 subcommand aktif, ~115 nilai flag, ~15 boolean flags = **~195 pilihan distinkt**

---

### GRUP A: Project Lifecycle

### 4.1 `go-aether init`

**Fungsi:** Inisiasi proyek Greenfield dari nol.

**Syntax:**
```bash
go-aether init <project-name> [flags]

Flags:
  --arch        string   hexagonal|clean|ddd|vertical|mvc|modular-monolith [default: hexagonal]
  --router      string   chi|gin|echo|stdlib [default: chi]
  --db          string   postgres|mysql|sqlite|mongodb|timescaledb|badger|none [default: postgres]
  --orm         string   sqlc|gorm|sqlx|raw|none [default: sqlc]
  --cache       string   redis|memcached|none [default: none]
  --auth        string   jwt-rs256|jwt-hs256|oauth2|session|none [default: none]
  --transport   string   grpc|websocket|nats|kafka|rabbitmq|none [default: none]
  --migrations  string   goose|migrate|atlas|none [default: goose]
  --logger      string   slog-otel|zerolog|zap|slog [default: slog-otel]
  --ci          string   github|gitlab [default: github]
  --module      string   Go module name (default: github.com/<user>/<project>)
  --go          string   Go version (default: deteksi dari go version)
  --monorepo             Init sebagai workspace monorepo (go.work)
  --dry-run              Tampilkan file yang akan di-generate tanpa menulis
  --verbose              Log detail setiap langkah
```

**Output yang di-generate:**
```
payment-service/
├── cmd/server/main.go              ← Entrypoint + Dependency Injection wiring
├── internal/
│   ├── config/config.go            ← Env-based config dengan validation
│   ├── core/
│   │   ├── domain/                 ← Structs, value objects, domain errors
│   │   ├── port/                   ← Interfaces (service + repository)
│   │   └── service/                ← Business logic (kosong, siap diisi)
│   └── adapter/
│       ├── handler/http/           ← Chi router setup + health endpoint
│       └── repository/postgres/    ← DB pool + ping check
├── pkg/
│   ├── middleware/                 ← RequestID, Logger, Recovery
│   └── validator/                  ← Input validation helper
├── migrations/                     ← Goose migration placeholder
├── .env.example                    ← Semua env vars yang dibutuhkan
├── .golangci.yml                   ← Opinionated linter config
├── Dockerfile                      ← Multi-stage build (builder + distroless)
├── docker-compose.yml              ← Postgres + Redis services
├── Makefile                        ← run, test, build, lint, migrate targets
├── aether.yaml                      ← Manifest (auto-generated)
└── go.mod
```

---

### 4.2 `go-aether adopt`

**Fungsi:** Adopsi proyek Go yang sudah ada (Brownfield) tanpa merusak apapun.

**Syntax:**
```bash
go-aether adopt [flags]

Flags:
  --scan     Auto-scan folder structure (mode interaktif)
  --arch     Override arsitektur yang terdeteksi
  --dry-run  Tampilkan aether.yaml yang akan dibuat tanpa menulis
```

**Alur Eksekusi (Mode `--scan`):**

```
1. go-aether adopt --scan
   ↓
2. CLI baca go.mod untuk mendapat module name & go version
   ↓
3. CLI scan struktur folder (2 level deep) untuk mendeteksi pola:
   - Ada folder "handler" / "controller"? → tanya konfirmasi
   - Ada folder "service" / "usecase"?    → tanya konfirmasi
   - Ada folder "repository" / "repo" / "store"? → tanya konfirmasi
   ↓
4. Kuesioner Interaktif (bubbletea):
   ? Arsitektur proyek ini paling mendekati: (pilih)
     ❯ hexagonal
       clean
       mvc
       tidak yakin (biarkan go-aether mapping manual)

   ? Folder handler/controller ada di: [jawab bebas]
   ? Folder service/usecase ada di:    [jawab bebas]
   ? Folder repository ada di:         [jawab bebas]
   ↓
5. Generate aether.yaml dengan anomaly_mode: true
   dan path sesuai jawaban user
   ↓
6. Tampilkan ringkasan dan instruksi langkah selanjutnya
```

---

### 4.3 `go-aether make:module`

**Fungsi:** Generate satu modul lengkap (domain + port + service + handler + repository) dalam satu perintah.

```bash
go-aether make:module <name> [flags]

Flags:
  --transport  string   http|grpc|both [default: http]
  --pattern    string   standard|cqrs|outbox|saga [default: standard]
  --cache               Tambahkan Redis cache layer
  --worker              Generate worker untuk modul ini
  --no-test             Skip generate _test.go
  --dry-run
```

**Contoh Output (hexagonal, `go-aether make:module order`):**
```
internal/core/domain/order.go                           ← Entity + Value Objects + Errors
internal/core/port/order_service_port.go                ← Interface IOrderService
internal/core/port/order_repo_port.go                   ← Interface IOrderRepository
internal/core/service/order_service.go                  ← Business logic stub + constructor
internal/adapter/handler/http/order_handler.go          ← Chi handler + routes
internal/adapter/repository/postgres/order_repository.go ← SQLC-ready stub
internal/core/service/order_service_test.go             ← Table-driven test skeleton
```

---

### 4.4 `go-aether make:service`

**Fungsi:** Generate service layer saja (tanpa handler & repository), siap diinject via constructor.

```bash
go-aether make:service <name> [flags]
  --no-test    Skip generate test file
```

---

### 4.5 `go-aether make:handler`

**Fungsi:** Generate transport handler saja untuk antarmuka tertentu.

```bash
go-aether make:handler <name> [flags]
  --transport  string  http|grpc [default: sesuai aether.yaml]
  --methods    string  GET,POST,PUT,DELETE,PATCH (kombinasi bebas)
  --stream     string  server|client|bidi (khusus gRPC streaming)
```

---

### 4.6 `go-aether make:repository`

**Fungsi:** Generate repository layer saja yang terhubung ke database driver tertentu.

```bash
go-aether make:repository <name> [flags]
  --db     string  postgres|mysql|sqlite|mongodb|timescaledb|badger [default: sesuai aether.yaml]
  --cache          Tambahkan cache layer
```

---

### 4.7 `go-aether make:migration`

**Fungsi:** Generate file SQL migrasi kosongan sesuai alat migrasi (goose, migrate, atlas).

```bash
go-aether make:migration <name> [flags]
  --tool   string  goose|migrate|atlas [default: sesuai aether.yaml]
# Contoh: go-aether make:migration add_orders_table
# → migrations/20260806192000_add_orders_table.sql
```

---

### GRUP B: Middleware & Security

### 4.8 `go-aether add:middleware`

**Fungsi:** Inject middleware siap pakai ke project yang sudah ada (zero magic, murni handler chaining).

```bash
go-aether add:middleware <type>

  # Auth & Session
  jwt-auth          JWT RS256/HS256 validation (OWASP-compliant)
  zero-trust        Full OWASP security headers bundle
  csrf              CSRF token protection
  session           Cookie-based session middleware

  # Traffic Control
  rate-limit      Token bucket rate limiter per-IP
  bulkhead        Concurrency limiter per endpoint (bulkhead pattern)
  circuit-breaker gobreaker integration
  timeout         Per-request context deadline
  idempotency     Idempotency-Key header enforcer

  # Observability
  audit-log         Structured access + mutation audit log
  cdn-cache         Cache-Control header policies

  # Infra
  cors              CORS whitelist
  recovery          Panic recovery + structured log
```

---



---

### 4.9 `go-aether add:auth`

```bash
go-aether add:auth <type>
  jwt-rs256   JWT dengan RS256 (asymmetric — production grade)
  jwt-hs256   JWT dengan HS256 (symmetric — simpler)
  oauth2      OAuth2 + OIDC flow (PKCE)
  argon2      Password hashing dengan Argon2id
  session     Session-based auth (Redis store)
```

---

### 4.10 `go-aether add:authz`

```bash
go-aether add:authz <type>
  casbin    RBAC/ABAC dengan Casbin
  opa       Open Policy Agent integration
```

---

### 4.11 `go-aether add:crypto`

```bash
go-aether add:crypto <type>
  aes-gcm   AES-256-GCM encryption helper
  ed25519   Ed25519 signing + verification
```

---

### 4.12 `go-aether add:security`

```bash
go-aether add:security <type>
  mtls          mTLS setup + certificate rotation scaffold
  secrets-scan  .golangci.yml update + pre-commit hook
```

---

### 4.13 `go-aether add:secrets`

```bash
go-aether add:secrets <type>
  vault       HashiCorp Vault dynamic secrets
  aws-ssm     AWS Parameter Store + Secrets Manager
  env         dotenv + envconfig validation
  rotation    Automated secret rotation scaffold
```

---

### GRUP C: Transport Layer

### 4.14 `go-aether add:transport`

```bash
go-aether add:transport <type>
  grpc          gRPC server + proto scaffold + interceptors
  websocket     WebSocket handler (gorilla/websocket)
  graphql       GraphQL server (gqlgen + DataLoader)
  nats          NATS JetStream subscriber
  kafka         Kafka consumer/producer (segmentio/kafka-go)
  rabbitmq      RabbitMQ consumer (amqp091-go)
  sqs           AWS SQS consumer
  mqtt          MQTT broker client (paho)
  connect-rpc   Connect-RPC (gRPC-Web compatible)
  quic          QUIC/HTTP3 server (quic-go)
```

---

### GRUP D: Data Layer

### 4.15 `go-aether add:db`

```bash
go-aether add:db <type>
  read-replica    Read replica routing setup
  pgvector        pgvector extension + similarity search repository
  timescale       TimescaleDB hypertable scaffold
```

---

### 4.16 `go-aether add:cache`

```bash
go-aether add:cache <type>
  l2     Multilevel cache (in-memory L1 + Redis L2 + singleflight)
  cdn    HTTP CDN Cache-Control header policies
```

---

### 4.17 `go-aether add:lock`

```bash
go-aether add:lock <type>
  redis    Distributed lock via Redlock (go-redis/redislock)
  pg       Postgres advisory lock helper
```

---

### GRUP E: Background Processing

### 4.18 `go-aether add:worker`

```bash
go-aether add:worker <name> [flags]
  --broker   string  redis|kafka|nats|rabbitmq [default: redis]
  --pattern  string  asynq|watermill|machinery|river [default: asynq]
```

---

### 4.19 `go-aether add:cron`

```bash
go-aether add:cron <name>
# Generate cron job scaffold menggunakan robfig/cron v3
# → internal/cron/daily_report.go
```

---

### 4.20 `go-aether add:workflow`

```bash
go-aether add:workflow <name> [flags]
  --engine  string  temporal|cadence [default: temporal]
```

---

### 4.21 `go-aether add:eventing`

```bash
go-aether add:eventing [flags]
  --broker  string  nats|kafka|rabbitmq [default: nats]
# Generate: EventBus interface + Publisher + Subscriber + Event registry
```

---

### GRUP F: Observability

### 4.22 `go-aether add:tracing`

```bash
go-aether add:tracing [flags]
  --exporter  string  jaeger|tempo|datadog|stdout [default: jaeger]
```

---

### 4.23 `go-aether add:metrics`

```bash
go-aether add:metrics <type>
  prometheus       /metrics endpoint + RED + USE metric helpers
  datadog          Datadog DogStatsD client
  otel-collector   OTel Collector gRPC exporter
```

---

### 4.24 `go-aether add:profiling`

```bash
go-aether add:profiling [flags]
  --pgo    Enable Profile-Guided Optimization scaffold
# Generate /debug/pprof endpoint + Makefile pgo target
```

---

### 4.25 `go-aether add:healthcheck`

```bash
go-aether add:healthcheck
# Generate /health/readiness + /health/liveness (Kubernetes probes)
# Checks: DB ping, Redis ping, uptime
```

---

### GRUP G: API Design

### 4.26 `go-aether add:docs`

```bash
go-aether add:docs <type>
  swagger    Swagger UI + ogen codegen dari OpenAPI spec
  scalar     Scalar modern API reference UI
```

---

### 4.27 `go-aether add:versioning`

```bash
go-aether add:versioning [flags]
  --strategy  string  url-prefix|header|content-type [default: url-prefix]
# Generate API versioning middleware + router group (/v1, /v2)
```

---

### 4.28 `go-aether add:webhook`

```bash
go-aether add:webhook <name>
# Generate webhook delivery engine:
# outbox table + delivery worker + retry + HMAC signature verification
```

---

### GRUP H: Architecture Patterns

### 4.29 `go-aether add:di`

```bash
go-aether add:di <type>
  wire    Google Wire DI codegen setup
```

---

### 4.30 `go-aether add:config`

```bash
go-aether add:config
# viper setup + env binding + validation + .env.example sync
```

---

### 4.31 `go-aether add:feature-flags`

```bash
go-aether add:feature-flags <type>
  flipt          Flipt self-hosted feature flags
  launchdarkly   LaunchDarkly SDK client
  env            Simple env-based feature flags (zero dependency)
```

---

### GRUP I: Deploy & Cloud

### 4.32 `go-aether add:deploy`

```bash
go-aether add:deploy <type>
  k8s       Kubernetes manifests (Deployment + Service + HPA + ConfigMap)
  lambda    AWS Lambda handler wrapper + SAM template
  helm      Helm chart scaffold (values.yaml + templates/)
```

---

### 4.33 `go-aether add:cloud`

```bash
go-aether add:cloud <type>
  aws    AWS SDK v2 setup + credentials chain
  gcp    Google Cloud client setup + ADC
```

---

### 4.34 `go-aether add:storage`

```bash
go-aether add:storage <type>
  s3       AWS S3 (aws-sdk-go-v2)
  gcs      Google Cloud Storage
  minio    MinIO self-hosted
  local    Local filesystem adapter
```

---

### 4.35 `go-aether add:discovery`

```bash
go-aether add:discovery <type>
  etcd       etcd service registration + discovery
  consul     Consul service mesh client
```

---

### GRUP J: Communication

### 4.36 `go-aether add:mailer`

```bash
go-aether add:mailer <type>
  smtp       SMTP via net/smtp + template
  sendgrid   SendGrid API v3 client
  mailgun    Mailgun API client
  ses        AWS SES client
```

---

### 4.37 `go-aether add:notification`

```bash
go-aether add:notification <type>
  fcm        Firebase Cloud Messaging (push)
  apns       Apple Push Notification Service
  web-push   Web Push API (VAPID)
```

---

### 4.38 `go-aether add:httpclient`

```bash
go-aether add:httpclient <name> [flags]
  --retry              Retry + exponential backoff + jitter
  --circuit-breaker    gobreaker integration
  --timeout  string    Default request timeout
  --tracing            Inject OTel trace headers
```

---

### GRUP K: AI & Data

### 4.39 `go-aether add:ai`

```bash
go-aether add:ai <type>
  embedding    Embedding pipeline scaffold
  llm-proxy    LLM inference proxy (OpenAI/Ollama compatible)
```

---

### GRUP L: Testing

### 4.40 `go-aether add:test`

```bash
go-aether add:test <type>
  integration   testcontainers-go scaffold
  e2e           End-to-end test scaffold (httptest)
  fuzz          Go native fuzz test scaffold per modul
  property      Property-based test (pgregory.net/rapid)
  synctest      Deterministic concurrency test (Go 1.24)
```

---

### GRUP M: Code Quality

### 4.41 `go-aether add:validation`

```bash
go-aether add:validation <type>
  playground   go-playground/validator v10
  ozzo         go-ozzo/ozzo-validation
```

---

### 4.42 `go-aether add:lint`

```bash
go-aether add:lint <type>
  custom    Custom linter menggunakan go/analysis framework
```

---

### GRUP N: Utilities

### 4.43 `go-aether add:multitenancy`

```bash
go-aether add:multitenancy [flags]
  --strategy  string  schema-per-tenant|rls|subdomain [default: rls]
```

---

### 4.44 `go-aether doctor`

**Fungsi:** Validasi kesehatan proyek dan konsistensi seluruh entitas di dalam `aether.yaml` dengan state filesystem & environment nyata.

```bash
go-aether doctor [flags]
  --fix   Auto-fix inkonsistensi yang bisa diperbaiki otomatis

Output:
  ✅ aether.yaml schema: valid (v1.0.0)
  ✅ go.mod module matches aether.yaml: github.com/company/payment-service
  ✅ Go version: 1.23 (minimum 1.21 satisfied)
  ✅ All 3 modules have matching files (order, payment, invoice)
  ⚠️  Module "user" declared in aether.yaml but handler file missing
  ⚠️  redis configured in aether.yaml but no REDIS_URL in .env.example
  ❌ Architecture path "internal/core/domain" does not exist on disk
```

---

### 4.45 `go-aether ls`

**Fungsi:** Menampilkan inventori lengkap seluruh modul, worker, middleware, dan kapabilitas yang telah di-generate di dalam proyek.

```bash
go-aether ls
# Tampilkan semua: modules, workers, middlewares, add-ons ter-register

Output:
  Modules (3):
  ┌─────────────┬──────────────────────┬──────────┬───────────┐
  │ Name        │ Created              │ Transport│ Cache     │
  ├─────────────┼──────────────────────┼──────────┼───────────┤
  │ order       │ 2026-08-06 20:00     │ http     │ redis ✓   │
  │ payment     │ 2026-08-06 21:00     │ http,grpc│ -         │
  │ invoice     │ 2026-08-07 09:00     │ http     │ -         │
  └─────────────┴──────────────────────┴──────────┴───────────┘
```

---

### 4.46 `go-aether upgrade`

**Fungsi:** Memperbarui skema `aether.yaml` ke versi spesifikasi terbaru saat binary `go-aether` di-update ke versi lebih tinggi.

```bash
go-aether upgrade [flags]
  --dry-run  Tampilkan perubahan yang akan dilakukan tanpa apply
```

---

## BAB 5 — Template Engine: Desain & Isi

### 5.1 Filosofi Template

Semua template di-embed ke dalam binary menggunakan `//go:embed templates/**`. Ini berarti:
- Tool bisa dipakai **offline sepenuhnya** (tidak perlu internet saat generate)
- Tidak ada CDN template yang bisa down
- Distribusi = 1 binary saja

### 5.2 Template Variables

Setiap template menerima struct `TemplateData`:

```go
type TemplateData struct {
    // Project info
    ModuleName    string  // "order"
    ModuleNamePkg string  // "order" (lowercase, sanitized)
    ModuleNameTitle string // "Order" (TitleCase)
    PackagePath   string  // "github.com/company/payment-service"
    GoVersion     string  // "1.23"
    
    // Architecture
    ArchPattern   string  // "hexagonal"
    Paths         PathMap // dari aether.yaml
    
    // Stack
    Router        string  // "chi"
    DBDriver      string  // "postgres"
    CacheDriver   string  // "redis"
    HasCache      bool
    HasWorker     bool
    Transports    []string
    
    // Auth
    AuthPattern   string  // "jwt-rs256"
    
    // BYO Dependency (brownfield)
    ExistingDBVar    string  // kosong = generate baru, non-kosong = terima dari luar
    ExistingRedisVar string
    
    // Meta
    Timestamp     string  // ISO8601
    AetherVersion  string  // "0.1.0"
}
```

### 5.3 Contoh Template Hexagonal: Domain

```go
// templates/hexagonal/domain.go.tmpl
package domain

import "errors"

// {{.ModuleNameTitle}} represents the core {{.ModuleNamePkg}} entity.
type {{.ModuleNameTitle}} struct {
    ID        string
    CreatedAt int64
    UpdatedAt int64
    DeletedAt *int64
}

// Sentinel errors for {{.ModuleNamePkg}} domain.
var (
    Err{{.ModuleNameTitle}}NotFound    = errors.New("{{.ModuleNamePkg}}: not found")
    Err{{.ModuleNameTitle}}InvalidInput = errors.New("{{.ModuleNamePkg}}: invalid input")
    Err{{.ModuleNameTitle}}Duplicate   = errors.New("{{.ModuleNamePkg}}: already exists")
)
```

### 5.4 Contoh Template Hexagonal: Service (dengan BYO injection)

```go
// templates/hexagonal/service.go.tmpl
package service

import (
    "context"
    
    "{{.PackagePath}}/internal/core/domain"
    "{{.PackagePath}}/internal/core/port"
    {{if .HasCache}}"{{.PackagePath}}/internal/core/port"{{end}}
)

type {{.ModuleNameTitle}}Service struct {
    repo  port.I{{.ModuleNameTitle}}Repository
    {{if .HasCache}}cache port.I{{.ModuleNameTitle}}Cache{{end}}
}

// New{{.ModuleNameTitle}}Service creates a new service with constructor injection.
// All dependencies are passed from outside — no global state.
func New{{.ModuleNameTitle}}Service(
    repo port.I{{.ModuleNameTitle}}Repository,
    {{if .HasCache}}cache port.I{{.ModuleNameTitle}}Cache,{{end}}
) *{{.ModuleNameTitle}}Service {
    return &{{.ModuleNameTitle}}Service{
        repo:  repo,
        {{if .HasCache}}cache: cache,{{end}}
    }
}
```

---

## BAB 6 — Edge Cases (Penanganan Kasus Tepi Komprehensif)

### 6.1 aether.yaml Tidak Ditemukan

**Skenario:** User menjalankan `go-aether make:service order` di dalam subdirektori proyek.

**Solusi:** CLI berjalan dengan algoritma `walk-up` — mencari `aether.yaml` mulai dari direktori saat ini ke atas hingga root filesystem atau batas git repo (`.git/` ditemukan):

```
Cari aether.yaml di: ./
  → Tidak ada
Cari aether.yaml di: ../
  → Tidak ada
Cari aether.yaml di: ../../
  → DITEMUKAN: ../../aether.yaml
  → Set working directory ke ../../ untuk semua path resolution
```

Jika sampai root tidak ditemukan:
```
ERROR: aether.yaml not found.
Run `go-aether init` to start a new project, or
Run `go-aether adopt` to use go-aether with an existing project.
```

---

### 6.2 File Tujuan Sudah Ada (Conflict)

**Skenario:** User menjalankan `go-aether make:module order` padahal `order.go` sudah ada dari generate sebelumnya.

**Solusi (3 Mode):**

```
⚠️  Conflict detected: internal/core/domain/order.go already exists.

? What would you like to do?
  ❯ skip       (keep existing file, do not overwrite)
    overwrite  (replace existing file with new template)
    diff       (show diff between existing and new)
    backup     (copy existing to order.go.bak, then overwrite)
    abort      (cancel entire generation)
```

Jika flag `--force` diberikan: langsung overwrite tanpa tanya.
Jika flag `--dry-run` diberikan: tampilkan diff, tidak ada perubahan.

---

### 6.3 Nama Modul Tidak Valid (Go Identifier)

**Skenario:** `go-aether make:module my-order` (mengandung tanda hubung).

**Solusi:**
```
ERROR: "my-order" is not a valid Go package name.
Go identifiers cannot contain hyphens.

Suggestions:
  → go-aether make:module myorder
  → go-aether make:module my_order
  → go-aether make:module MyOrder
```

---

### 6.4 Nama Modul Duplikat

**Skenario:** `aether.yaml` sudah punya modul `order`, tapi user coba `go-aether make:module order` lagi.

**Solusi:**
```
ERROR: Module "order" already exists in aether.yaml.

Existing files:
  internal/core/domain/order.go
  internal/core/service/order_service.go
  ...

If you want to regenerate, use: --force
If you want to add a transport, use: go-aether make:handler order --transport grpc
```

---

### 6.5 Dependency BYO Brownfield: Redis/DB Sudah Ada

**Skenario:** Proyek legacy punya `var redisClient *redis.Client` di `main.go`. User adopsi dengan go-aether dan generate repository baru.

**Solusi:**
Jika `adapters.existing_redis_var` di `aether.yaml` **non-empty**, generator mengetahui bahwa ia harus **menerima koneksi dari luar** via constructor, bukan membuat koneksi baru:

```go
// YANG DI-GENERATE (brownfield mode):
type OrderRepository struct {
    db    *pgxpool.Pool   // diterima dari constructor
    cache *redis.Client   // diterima dari constructor  
}

func NewOrderRepository(db *pgxpool.Pool, cache *redis.Client) *OrderRepository {
    return &OrderRepository{db: db, cache: cache}
}
```

Di `main.go` legacy, user tinggal:
```go
// main.go legacy (tidak diubah go-aether)
orderRepo := repository.NewOrderRepository(existingDB, existingRedisClient)
```

Zero duplikasi koneksi, zero connection pool collision.

---

### 6.6 Deteksi Anomali Brownfield (Folder Tidak Konsisten)

**Skenario:** Struktur folder legacy sangat tidak konvensional:
```
myapp/
  ├── web/controllers/    ← Handler ada di sini
  ├── logic/              ← Service ada di sini
  └── data/               ← Repository ada di sini
```

**Solusi:** Kuesioner interaktif `go-aether adopt --scan` akan memetakan ini ke `aether.yaml`:

```yaml
architecture:
  pattern: "hexagonal"
  mode: "brownfield"
  paths:
    handler_http:  "web/controllers"   ← Path KUSTOM dari input user
    service:       "logic"
    repository:    "data"
  
meta:
  anomaly_mode: true
  legacy_notes: "Non-standard folder structure from pre-2020 rewrite"
```

Saat `go-aether make:module order` dijalankan di proyek ini, file akan diletakkan di `web/controllers/order_handler.go`, `logic/order_service.go`, `data/order_repository.go` — **menghormati konvensi lama, tidak menciptakan gaya baru!**

---

### 6.7 Arsitektur Mismatch

**Skenario:** aether.yaml menyatakan `pattern: hexagonal`, tapi user mencoba `go-aether make:handler order --style mvc`.

**Solusi:**
```
ERROR: Architecture mismatch.
aether.yaml declares: hexagonal
You requested:        mvc

Hexagonal handlers live in: internal/adapter/handler/http/
MVC-style controllers are not compatible with the hexagonal pattern.

If you want to change architecture, run:
  go-aether upgrade --arch mvc  (WARNING: this will not migrate existing files)
```

---

### 6.8 Go Version Terlalu Lama

**Solusi:** go-aether membutuhkan minimum Go 1.21 (untuk `slog` stdlib):
```
ERROR: go-aether requires Go 1.21 or later.
Your project's go.mod specifies: go 1.18

Please update your go.mod:
  go mod edit -go=1.21
  go mod tidy
```

---

### 6.9 Partial Generation Failure (Rollback)

**Skenario:** go-aether sedang generate 7 file, berhasil 5, lalu terjadi error (disk penuh / permission denied) di file ke-6.

**Solusi:** go-aether menggunakan **transactional write pattern** via buffer:
1. Semua file di-render ke buffer memory terlebih dahulu (tanpa tulis ke disk)
2. Validasi semua buffer sukses → baru batch-write ke disk
3. Jika ada error saat write → rollback semua file yang sudah tertulis

```
ERROR: Failed to write internal/adapter/repository/postgres/order_repository.go
  Reason: permission denied

Rolling back changes:
  ❌ Deleting internal/core/domain/order.go (rolled back)
  ❌ Deleting internal/core/port/order_service_port.go (rolled back)
  ...

No files were permanently modified. Please fix permissions and try again.
```

---

### 6.10 go.mod Tidak Ditemukan

```
ERROR: go.mod not found in current directory or parent directories.

go-aether requires a valid Go module.
Initialize one first: go mod init github.com/yourname/yourproject
Then run: go-aether init OR go-aether adopt
```

---

### 6.11 Goimports / Gofmt Tidak Terinstall

**Skenario:** go-aether mencoba format file setelah generate, tapi `goimports` tidak ada di PATH.

**Solusi:** Graceful degradation — generate tetap berhasil, tapi tampilkan warning:

```
⚠️  WARNING: goimports not found in PATH.
Generated files may have unresolved imports.

Install with: go install golang.org/x/tools/cmd/goimports@latest
Or disable auto-format: go-aether config set autoformat false
```

---

### 6.12 Nama Modul Collision dengan Standard Library

**Skenario:** `go-aether make:module context`

**Solusi:**
```
ERROR: "context" conflicts with Go standard library package name.
Choose a different module name:
  → go-aether make:module appcontext
  → go-aether make:module requestcontext
```

Daftar kata reserved yang diblokir mencakup semua stdlib package name.

---

### 6.13 aether.yaml Versi Lama (Schema Migration)

**Skenario:** User meng-update go-aether ke versi baru tapi `aether.yaml` masih schema lama.

**Solusi:**
```
WARNING: aether.yaml schema v0.9 detected, but go-aether v1.0.0 expects schema v1.0.

Run: go-aether upgrade --dry-run  (lihat perubahan)
     go-aether upgrade             (apply migrasi otomatis)

Backward compatibility: go-aether 1.x tetap bisa baca aether.yaml v0.9,
tapi beberapa fitur baru tidak akan tersedia sampai diupgrade.
```

---

### 6.14 Koneksi yang Dideclare di aether.yaml Tapi Tidak Ada di .env.example

`go-aether doctor` mendeteksi:
```
⚠️  aether.yaml declares cache: redis
    But REDIS_URL is missing from .env.example

Auto-fix available: go-aether doctor --fix will add REDIS_URL="" to .env.example
```

---

### 6.15 Vertical Slice Mode (Emergency Escape Hatch)

Untuk proyek legacy yang terlalu kacau untuk dipetakan:

```bash
go-aether make:feature order --vertical --dir ./features/order
```

Generate semua komponen dalam satu folder terisolasi:
```
features/order/
├── handler.go      ← Self-contained handler
├── service.go      ← Self-contained service
├── repository.go   ← Self-contained repository
├── domain.go       ← Local entity
└── order_test.go   ← Test untuk semuanya
```

Tidak ada pemetaan path, tidak ada konflik, 100% terisolasi.

---

### 6.16 Windows Path Separator

**Skenario:** User di Windows → path `internal\core\domain` vs `internal/core/domain`.

**Solusi:** go-aether selalu menggunakan `filepath.Join()` dan `filepath.ToSlash()` secara internal. `aether.yaml` selalu menyimpan dalam format `/` (Unix-style). Saat dioperasikan di Windows, konversi dilakukan secara transparan.

---

## BAB 7 — Testing Strategy untuk go-aether Sendiri

### 7.1 Unit Tests (via `afero` mock FS)

```go
func TestMakeModule_Hexagonal_WritesCorrectFiles(t *testing.T) {
    fs := afero.NewMemMapFs()
    // Setup aether.yaml di in-memory filesystem
    // Jalankan make:module service
    // Assert file yang dihasilkan ada dan isinya benar
    // TANPA menyentuh disk nyata
}
```

### 7.2 Integration Tests (via temporary real directory)

```go
func TestInitCommand_CreatesValidGoProject(t *testing.T) {
    tmpDir := t.TempDir()
    // Jalankan go-aether init di tmpDir
    // Assert: go.mod valid, go build berhasil, aether.yaml valid
}
```

### 7.3 Golden File Tests

Setiap template punya "golden file" — output yang diharapkan. Saat template berubah, jalankan `go test -update-golden` untuk meregenerasi. Ini mencegah regresi template tanpa sengaja.

### 7.4 Mutation Tests (via `go-mutesting`)

Membuktikan test kita tidak tautologis — test harus merah saat kode dimutasi.

### 7.5 Coverage Gate

`go-aether doctor` sendiri di-test dengan coverage threshold 80% (sesuai yang di-generate ke user).

---

## BAB 8 — Distribusi & Open Source Strategy

### 8.1 Release Pipeline

```
Push to main
    → GitHub Actions: go test ./... (all platforms)
    → golangci-lint
    → goreleaser (build binary untuk linux/amd64, darwin/arm64, windows/amd64)
    → Upload ke GitHub Releases
    → Update Homebrew formula (tap: muhananaufal/homebrew-go-aether)
    → Update Scoop bucket (untuk Windows)
```

### 8.2 Install Methods

```bash
# Via go install (paling simpel)
go install github.com/muhananaufal/go-aether@latest

# Via Homebrew (macOS/Linux)
brew install muhananaufal/go-aether/go-aether

# Via Scoop (Windows)
scoop bucket add go-aether https://github.com/muhananaufal/scoop-go-aether
scoop install go-aether

# Via GitHub Releases (manual download)
# Binary pre-built tersedia di github.com/muhananaufal/go-aether/releases
```

### 8.3 Versioning

Ikuti Semantic Versioning (SemVer):
- MAJOR: Breaking change di aether.yaml schema atau command interface
- MINOR: Tambah arsitektur pattern baru, tambah transport baru
- PATCH: Bug fix, perbaikan template

### 8.4 Dokumentasi

- `README.md`: QuickStart 5 menit
- `docs/`: Full documentation site (via MkDocs atau VitePress)
- `CONTRIBUTING.md`: Panduan contribution (template baru, transport baru)
- `CHANGELOG.md`: Semua perubahan per versi

---

## BAB 9 — Roadmap (MVP → v1 → v2)

### MVP (v0.1.0) — Target: 4 Minggu

**Scope ketat — hanya yang benar-benar dibutuhkan:**
- [ ] `go-aether init` dengan `--arch hexagonal` dan `--db postgres` + `--router chi`
- [ ] `go-aether make:module` (domain + port + service + handler + repository)
- [ ] `aether.yaml` schema v1 (greenfield only)
- [ ] `go-aether doctor` (validasi basic)
- [ ] Template: hexagonal only, postgres only, chi only
- [ ] Output: Makefile, Dockerfile multi-stage, .env.example, .golangci.yml
- [ ] Tests: unit + golden files (coverage >80%)
- [ ] Distribusi: `go install` + GitHub Releases (3 platform)

---

### v0.5.0 — Target: 8 Minggu dari MVP

- [ ] `go-aether adopt --scan` (brownfield dengan kuesioner interaktif)
- [ ] `go-aether add:middleware` (jwt-auth, rate-limit, cors, otel-trace)
- [ ] `go-aether make:module` tambahan transport `--transport grpc`
- [ ] Template tambahan: `clean` architecture
- [ ] Tambah router: `gin` dan `echo`

---

### v1.0.0 — Stable Release

- [ ] `go-aether add:worker` (asynq + redis)
- [ ] `go-aether add:transport websocket`
- [ ] Template: `ddd` dan `vertical` architecture
- [ ] `go-aether ls` dengan output table cantik
- [ ] `go-aether upgrade` (schema migration)
- [ ] Dokumentasi lengkap di situs dedicated
- [ ] Homebrew + Scoop
- [ ] Test coverage enforced di CI/CD

---

### v2.0.0 — Ecosystem Play (Kondisional pada Adopsi v1)

- [ ] Plugin system: komunitas bisa tambahkan template custom via `aether.yaml`
- [ ] `go-aether add:observability` (Jaeger + Prometheus + Grafana stack)
- [ ] `go-aether add:auth oauth2` (PKCE flow)
- [ ] Multi-module monorepo support
- [ ] VS Code extension: IntelliSense untuk aether.yaml

---

## BAB 10 — Dampak Karir (Career Trajectory)

### Mengapa membangun `go-aether` mendongkrak karir backend Golang:

| Skill yang Dibuktikan | Cara Dibuktikannya |
|:---|:---|
| **Deep Go Internals** | Pakai `embed.FS`, `text/template`, `afero`, `cobra` — bukan tutorial basic |
| **Architecture Mastery** | Tool kamu *mengajarkan* Hexagonal/Clean/DDD ke orang lain |
| **CLI Engineering** | Domain yang sangat dihargai di Go ecosystem |
| **Brownfield Thinking** | Bisa handle legacy — ini yang 90% perusahaan butuhkan |
| **Open Source Contribution** | Portofolio yang berbicara sendiri tanpa CV |
| **Cross-Platform Expertise** | Windows + macOS + Linux path handling |
| **Testing Discipline** | `afero` mock FS + golden file + mutation testing |
| **DevOps Mindset** | goreleaser + GitHub Actions + Homebrew tap |

### Narasi di Interview (Principal/Staff Level):

> *"Saya membangun go-aether — sebuah architecture scaffold CLI untuk Go yang mendukung Hexagonal, Clean, dan DDD pattern. Yang membedakannya adalah kemampuan Brownfield adoption via Interactive Anomaly Mapping, sehingga tim bisa mengadopsi arsitektur modern tanpa Big Bang Rewrite. Saat ini digunakan oleh [N] developer dan memiliki [K] GitHub stars."*

Narasi ini langsung membuktikan: arsitektur, engineering mindset, dan kemampuan memimpin arah teknis tim.

---

## BAB 11 — Anti-Scope (Yang TIDAK Akan Dibangun)

Batasan eksplisit agar scope tidak membesar tanpa kendali:

| Tidak Dibangun | Alasan |
|:---|:---|
| Runtime framework (laravel/nestjs style) | Melanggar filosofi Go, membunuh performa |
| ORM kustom | SQLC + pgx sudah optimal |
| UI web untuk manage proyek | Ini CLI, bukan SaaS |
| Support PHP/TypeScript/Rust | Out of scope — Go only |
| Auto-migration database | Goose sudah cukup, jangan duplikasi |
| Deployment tool (terraform, k8s) | Bukan domain scaffolding |
| Code review bot | Product berbeda |

---

## PENUTUP

`go-aether` bukan hanya sebuah tool — ini adalah **pernyataan arsitektur** yang tertulis dalam bahasa Go itu sendiri. Setiap perintah yang dihasilkannya adalah kristalisasi dari ratusan jam riset best practice yang sudah kita bangun bersama di 221 referensi domain Golang Mastery.

Saat MVP-nya jalan, setiap junior Go developer yang menjalankan `go-aether init` secara tidak sadar mendapat mentor senior yang mendampinginya dari hari pertama — tanpa harus membayar konsultan atau membaca 10 buku arsitektur.

Dan ketika itulah tool ini benar-benar berhasil.

---

*"Don't build frameworks. Build tools that generate idiomatic code."*  
*— Filosofi Inti go-aether*
