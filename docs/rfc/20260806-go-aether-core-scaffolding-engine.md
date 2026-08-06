# RFC: `go-aether` — Core Scaffolding Engine, CLI Framework & Manifest SSOT

- **Status:** `ACCEPTED` *(Disetujui via "Gasskan" - Sedang Eksekusi Batch 1)*
- **Tanggal:** 2026-08-06
- **Target Branch:** `feature/core-engine`
- **Penulis:** AETHERIS + muhananaufal
- **Standar Eksekusi:** All-Out Principal Engineering Standard (AETHERIS Master Key Formula / L8)

---

## 1. Konteks & Problem Statement (PRD Core & 5W1H First-Principles)

### 1.1 Latar Belakang & Urgensi (5W1H Rooting)
Dalam ekosistem rekayasa perangkat lunak backend masa kini, bahasa **Golang** dihormati karena kesederhanaan sintaks (*simplicity*), konkurensi native (*goroutines/channels*), dan performa komputasi yang mendekati bahasa C. Namun, ekosistem Golang sengaja dibangun berdasarkan filosofi **"library over framework"**. Tidak ada kerangka kerja monolitik resmi (seperti Laravel di PHP atau NestJS di TypeScript) yang mendikte struktur direktori atau pola desain enterprise.

Akibat dari kebebasan absolut ini menciptakan permasalahan fatal di lapangan:
1. **Fragmentasi Arsitektur & Spaghetti Code:** Engineer junior hingga menengah (0–3 tahun pengalaman Go) kerap merasa kebingungan menghadapi kanvas kosong. Mereka mengandalkan tutorial parsial yang berujung pada inkonsistensi proyek, pencampuran logika bisnis di lapisan HTTP handler, serta kebergantungan global (*global variable state*) yang merusak testability.
2. **Kutukan Brownfield (Legacy Sprawl):** Lebih dari 90% engineer berhadapan dengan proyek eksisting (*Brownfield*) yang berantakan. Mengadopsi standar baru biasanya menuntut pengodean ulang dari nol (*Big Bang Rewrite*), sebuah risiko finansial dan teknis yang hampir selalu ditolak oleh pemangku kepentingan.
3. **Overhead Runtime dari Framework Alternatif:** Upaya meniru kerangka kerja berbasis refleksi (*runtime reflection*) atau *magic dependency injection* justru merusak filosofi Go: waktu komputasi melonjak, binary menjadi gempal, dan alur eksekusi kehilangan transparansi.

**`go-aether`** lahir sebagai solusi atas kebutuhan tersebut. Ia adalah sebuah **Opinionated Architecture Scaffold CLI Engine**: perangkat bantu berbasis CLI berkecepatan tinggi yang berjalan **hanya pada waktu pengembang (*dev-time*)**. Ia mematuhi pola *Hexagonal / Ports & Adapters*, menyertakan *Constructor Injection* murni, dan memanfaatkan `embed.FS` dari standard library Go. Dengan demikian, `go-aether` menghasilkan struktur kode yang 100% *idiomatic*, *strongly-typed*, dan menolak menambahkan satu baris pun *runtime overhead* ke dalam binary produksi.

### 1.2 Target Pengguna (Who)
- **Target Primer:** Junior-to-Mid Backend Go Developer yang baru berpindah dari stack monolitik (Laravel/NestJS) dan memerlukan sistem pengaman arsitekural (*architectural guardrails*) tanpa harus langsung mengkaji puluhan buku literarur sistem terdistribusi.
- **Target Sekunder:** Tech Lead, Staff, atau Principal Engineer yang bertanggung jawab menstandardisasi konvensi penulisan kode (*lint, test harness, folder layers*) melintasi banyak microservices di dalam organisasi.
- **Dampak Karir:** Membangun CLI generator di dalam bahasa Go itu sendiri (*dogfooding*) adalah proyeksional keterampilan tingkat L7/L8 (Principal Backend Engineer), memperagakan penguasaan mendalam atas Go AST, filesystem transactional writing, cross-platform path resolution, serta strategi *Brownfield adoption*.

### 1.3 Scope Boundaries
- **IN-SCOPE (MVP v0.1.0 Fokus Foundation):**
  - Inisiasi Proyek Greenfield (`go-aether init`) untuk pola arsitektur **Hexagonal Architecture** dengan router `chi` dan database `postgres` (SQLC ready).
  - Generator Modul Vertikal (`go-aether make:module`, `make:service`, `make:handler`, `make:repository`) yang secara otomatis menyambungkan lapisan *Domain*, *Ports*, dan *Adapters*.
  - Pemrosesan spesifikasi registri sistem abadi dalam file manifest **`aether.yaml`** (SSOT) secara reflektif berbasis `viper`.
  - Sistem validasi keselarasan proyek otomatis (`go-aether doctor`).
  - Mesin rendering berbasis `text/template` dengan bundling aset menggunakan `embed.FS` serta isolasi sistem berkas untuk pengujian mutasi via `afero`.
- **EXPLICIT NON-GOALS (Batas Keras / Anti-Scope):**
  - DILARANG membangun runtime HTTP framework atau menyuntikkan *magic runtime package*. Semua keluaran wajib murni memakai stdlib dan pustaka OSS mapan.
  - DILARANG membuat ORM kustom. Interaksi database secara default mengandalkan SQL generator bertipe statis (SQLC / pgxpool).
  - DILARANG menulis antarmuka Web UI / SaaS pengelolaan proyek. Alat ini sepenuhnya hidup di lingkungan terminal dan CI/CD pipeline.

### 1.4 Asumsi Epistemik & Target Performa (Skala Sistem CLI & Output)
- **Eksekusi Dev-Time CLI:** Waktu parsing command, validasi YAML, pembicara template, dan penulisan berkas bertransaksi disk harus selesai di bawah **100 milidetik (<100ms)** pada mode offline absolut.
- **Kondisi Output Produksi:** Kode Go yang dihasilkan oleh `go-aether` bergaransi **0% Overhead Tambahan**. Binary akhir tidak mengimpor modul `go-aether` itu sendiri, menjamin latensi komputasi murni setara rawa stdlib/chi (p95 < 10ms, p99 < 50ms untuk operasi I/O standar).
- **Target Lingkungan Eksekusi:** Berjalan mulus di lintas sistem operasi: Windows (PowerShell/CMD), macOS (Darwin ARM64/AMD64), dan Linux (AMD64/ARM64) dengan normalisasi *path separating separator* (`/` vs `\`) yang konsisten secara atomik.

---

## 2. Eksplorasi Arsitektur & Trade-off Matrix (§4 Anti-Yes-Man Protocol)

Untuk memenuhi kebutuhan standardisasi kode tanpa melompoti norma etis rekayasa software sistem, kami telah membandingkan **3 Pendekatan Arsitektur Netral** secara klinis dan kritis tanpa bias persetujuan cepat.

### Opsi A: Eksternal Script / Generic Template Engine (Yeoman / Cookiecutter / Python)
Menggunakan ekosistem di luar bahasa Go (seperti Cookiecutter berbasis Python atau Yeoman berbasis Node.js) yang memanipulasi direktori string secara universal.
- **Deskripsi Arsitektur:** Script Python/JS diunduh secara eksternal; skrip membaca folder template kasar dan langsung melalap nama direktori serta string penampung (`{{ variable }}`) di filesystem.
- **Kelebihan (≥3):**
  1. Sangat cepat dibangun untuk iterasi awal tanpa mempedulikan struktur kompilasi biner.
  2. Didukung oleh ekosistem tool generator string lintas bahasa yang amat matang di dunia Python/Node.
  3. Memiliki dukungan kuesioner terminal bawaan yang mudah diintegrasikan dengan sintaks konfigurasi JSON sederhana.
- **Kekurangan & Failure Modes (≥3):**
  1. **Dependency Pollution:** Developer Go diwajibkan memasang dan mengelola environment Python/Node.js beserta virtualenv atau `node_modules` raksasa hanya demi membuat folder proyek Go.
  2. **Zero Go Context Awareness:** Engine eksternal buta terhadap hukum bahasa Go. Ia tidak memahami validasi identifikasi package (`package my-order` yang cacat), sintaks AST, ataupun impor go.mod.
  3. **Fragility in Brownfield:** Tidak mampu membaca struktur *package interface* eksisting atau mengurai manifest modular untuk diselaraskan dengan aman ke dalam kode lama.
- **Reversibility:** `Two-Way Door` (Mudah dibuang, namun merugikan dari sisi rekayasa kredibilitas developer Go).

### Opsi B: Runtime Application Framework Berbasis Refleksi (Loco / Fiber Monolith Style)
Membangun framework murni di dalam Go yang merangkum router, kontainers dependensi otomatis (*IoC Container*), dan abstraksi database ke dalam sebuah paket komputasi terpadu yang menyertai running binaries.
- **Deskripsi Arsitektur:** Developer mengimpor `github.com/go-aether/framework` di dalam `main.go`. Inisiasi router, penambatan repository, dan penanganan HTTP request dilayani secara tersembunyi via *struct tags* dan *reflect.Value*.
- **Kelebihan (≥3):**
  1. Pengembang baru merasa sangat familiar seperti bekerja dengan NestJS, Spring Boot, atau Laravel.
  2. Jumlah baris kode yang ditulis di dalam proyek pengguna (*boilerplate reduction*) menjadi sangat renggang dan singkat.
  3. Perubahan logika fundamental pada antarmuka framework dapat di-update serentak hanya dengan menaikkan versi dependency di `go.mod`.
- **Kekurangan & Failure Modes (≥3):**
  1. **Performance Degradation & Memory Bloat:** Refleksi (*runtime reflection*) merusak mekanisme optimasi statis kompilator Go, memangkas kecepatan eksekusi CPU hingga 3–5 kali lipat, dan memicu *Heap Allocations* tak terkendali saat *Garbage Collection (GC)*.
  2. **Obfuscation & Vendor Lock-In:** Alur eksekusi tertidur di balik tirai hitam "magic". Jika terjadi kebocoran memori atau bottleneck konkurensi di production, tracing dengan `pprof` menjadi sangkar labirin yang mustahil dibongkar.
  3. **Violation of Go Idioms:** Secara langsung menentang filosofi utama Golang ("Explicit over Implicit"), membuat adopsinya ditolak mentah-mentah oleh komunitas dan Principal Engineer di perusahaan teknologi kelas atas.
- **Reversibility:** `One-Way Door` (Sangat berbahaya; membongkar framework runtime dari ratusan endpoint menuntut biaya *rewrite* total bernilai mahal).

### Opsi C: Pure Go Opinionated Dev-Time CLI Engine (`embed.FS` + `afero` + BYO Injection)
Membangun perkakas CLI murni berspesikasi komputasi mandiri dalam bahasa Go. Binary berisi cetakan kode berlisensi murni yang menyembahkan output tanpa jejak framework di dalam produksi.
- **Deskripsi Arsitektur:** Binary CLI dipilah atas arsitektur Hexagonal miliknya sendiri. Template disimpan di dalam `//go:embed`, dieksaminasi via `text/template`, dinavigasi oleh `viper`, dan ditulis secara atomik ke filesystem via antarmuka `afero`.
- **Kelebihan (≥3):**
  1. **Zero Runtime Overhead & Idiomatic Purity:** Kode output berwujud murni *plain old Go*. Tiada biaya performa tambahan sepersejuta detik pun pada saat produksi berjalan; komputasi bersinar dalam kecepatan native C-like Go.
  2. **Offline & Self-Contained Distribution:** File eksekusi binernya tunggal (*single binary*). Tidak memerlukan unduhan CDN eksternal, tidak membutuhkan runtime runtime asing, langsung aktif setelah `go install`.
  3. **Brownfield BYO Dependency Resilience:** Dapat dengan elegan melayang ke atas lautan spaghetti kode lawas berkat pemetakan patokan folder anomali dan teknik injeksi konstruktor eksternal (*Bring Your Own Dependency*).
- **Kekurangan & Failure Modes (≥3):**
  1. Membutuhkan penguasaan rekayasa struktural yang rumit di pihak arsitek pembuat CLI (penulisan cetakan bersarang di dalam string Go).
  2. Pemutakhiran standar template di masa depan tidak langsung menyembuhkan kode yang telah di-generate sebelumnya (harus ditangani melalui tool migrasi seperti `go-aether upgrade` atau refactor konvensional).
  3. Menjaga kelayakan template sintaktis menuntut kedisiplinan pengujian ketat (*Golden File & Mutation Testing*) saat siklus pengembangan CLI.
- **Reversibility:** `Two-Way Door` (Sangat fleksibel; jika tim putuskan berhenti memakai `go-aether`, kode yang dihasilkan tetap bernilai sejati, valid, dan dapat dirawat seperti biasa tanpa ada ikatan tali ranting tool sama sekali).

> **Keputusan Arsitektur:** Kita mengeksekusi **Opsi C** secara mutlak dan meyakinkan, karena menjiwai esensi sejati ilmu keahlian Go Backend Engineering berderajat Principal L7/L8.

---

## 3. Spesifikasi Teknis & Desain Sistem Terpilih

### 3.1 Arsitektur Internal CLI ("Dogfooding Pattern")
`go-aether` menolak sikap munafik; ia sendiri disusun menerapkan konvensi **Hexagonal / Ports & Adapters Architecture** berdisiplin baja:

```
go-aether/
├── cmd/
│   └── go-aether/main.go               ← Entrypoint CLI (Koneksi cobra rootCmd)
├── internal/
│   ├── config/                         ← Viper binding & env loading
│   ├── core/
│   │   ├── domain/                     ← Entitas inti (AetherManifest, Module, Arch, Sentinel Errors)
│   │   ├── port/                       ← Antarmuka kontrak (Generator, ManifestResolver, FileWriter, Logger)
│   │   └── service/                    ← Orchestration logic (InitService, AdoptService, MakeService, DoctorService)
│   └── adapter/
│       ├── cli/                        ← Implementasi subcommand cobra & antarmuka TUI (bubbletea)
│       ├── manifest/                   ← Resolusi parser aether.yaml berbasis viper
│       ├── template/                   ← Engine rendering stdlib text/template + embed.FS
│       └── writer/                     ← Filesystem abstraction (os-fs untuk nyata, afero-mem map untuk pengujian)
├── templates/                          ← Seluruh arsip cetakan (*.tmpl) terbundle statis dalam binary
│   ├── hexagonal/
│   ├── clean/
│   ├── ddd/
│   └── common/
└── aether.yaml                         ← go-aether mengonsumsi manifest aether miliknya sendiri!
```

### 3.2 Manifest SSOT (`aether.yaml`) & Skema Kontrak Abadi
File `aether.yaml` diletakkan pada akar (*root*) repositori dan wajib masuk ke version control Git. Ia adalah rekam medis dan komandan strategi sistem:

```yaml
version: "1"
project:
  name: "payment-service"
  module: "github.com/company/payment-service"
  go_version: "1.23"
  created_at: "2026-08-06T19:00:00+07:00"
  aether_version: "0.1.0"
architecture:
  pattern: "hexagonal"           # hexagonal | clean | ddd | vertical | mvc
  mode: "greenfield"             # greenfield | brownfield
  paths:
    domain:        "internal/core/domain"
    port:          "internal/core/port"
    service:       "internal/core/service"
    handler_http:  "internal/adapter/handler/http"
    repository:    "internal/adapter/repository"
    cmd:           "cmd/server"
    config:        "internal/config"
    pkg:           "pkg"
stack:
  router:    "chi"
  database:
    driver:  "postgres"
    orm:     "sqlc"
    pool:    "pgxpool"
  cache:     "redis"
  transport: ["http"]
  auth:      "jwt-rs256"
  logger:    "slog-otel"
adapters:
  existing_db_var: ""            # Kosong untuk greenfield; terisi (misal: app.DBPool) pada brownfield BYO
  existing_redis_var: ""
modules:
  - name: "order"
    created_at: "2026-08-06T20:00:00+07:00"
    transports: ["http"]
    has_cache: true
    has_worker: false
meta:
  anomaly_mode: false
  legacy_notes: ""
```

### 3.3 Exhaustive Command & Capabilities Taxonomy (195-Option Mastery Map)
Sebagai pembuktian pemenuhan pilar *Exhaustive Inventory Mapping*, berikut adalah klasifikasi taksonomik lengkap dari 14 Grup Komandan di dalam blueprint kita yang merekonstruksi 221 file referensi Golang Mastery tanpa ada penyinoniman ataupun pangkas silang (`// dll` dibekukan):

| Grup Komandan | Subcommand / Domain CLI | Daftar Opsi Flag & Kapabilitas Lengkap Tanpa Batas | Target Folder Output / Lapisan Integrasi |
| :--- | :--- | :--- | :--- |
| **Grup A: Lifecycle** | `init`, `adopt`, `make:module`, `make:service`, `make:handler`, `make:repository`, `make:migration` | `--arch` (hexagonal/clean/ddd/vertical/mvc/modular-monolith), `--router` (chi/gin/echo/stdlib), `--db` (postgres/mysql/sqlite/mongodb/timescaledb/badger/none), `--orm` (sqlc/gorm/sqlx/raw/none), `--cache` (redis/memcached/none), `--migrations` (goose/migrate/atlas), `--pattern` (standard/cqrs/outbox/saga), `--stream` (server/client/bidi), `--monorepo`, `--scan`, `--dry-run`, `--verbose` | `cmd/server/main.go`, `internal/core/{domain,port,service}/`, `internal/adapter/{handler,repository}/`, `migrations/`, `.env.example`, `Makefile`, `Dockerfile` |
| **Grup B: Middleware & Security** | `add:middleware`, `add:auth`, `add:authz`, `add:crypto`, `add:security`, `add:secrets` | *Middleware:* jwt-auth, zero-trust, csrf, session, rate-limit (token bucket), bulkhead, circuit-breaker (gobreaker), timeout, idempotency, audit-log, cdn-cache, cors, recovery.<br>*Auth:* jwt-rs256, jwt-hs256, oauth2, argon2, session.<br>*AuthZ:* casbin, opa.<br>*Crypto:* aes-gcm, ed25519.<br>*Security:* mtls, secrets-scan.<br>*Secrets:* vault, aws-ssm, env, rotation | `pkg/middleware/`, `internal/common/security/`, `internal/common/crypto/`, `internal/common/secrets/` |
| **Grup C: Transport** | `add:transport` | grpc, websocket (gorilla), graphql (gqlgen+DataLoader), nats (JetStream), kafka (segmentio), rabbitmq (amqp091), sqs, mqtt (paho), connect-rpc, quic | `internal/adapter/handler/{grpc,ws,graphql,event}/`, `proto/` |
| **Grup D: Data Layer** | `add:db`, `add:cache`, `add:lock` | *DB:* read-replica, pgvector (similarity search), timescale (hypertable).<br>*Cache:* l2 (memory L1 + redis L2 + singleflight), cdn.<br>*Lock:* redis (redlock), pg (advisory) | `internal/adapter/repository/`, `internal/common/cache/`, `internal/common/lock/` |
| **Grup E: Background Jobs** | `add:worker`, `add:cron`, `add:workflow`, `add:eventing` | *Worker:* `--broker` (redis/kafka/nats/rabbitmq), `--pattern` (asynq/watermill/machinery/river).<br>*Cron:* robfig/cron v3.<br>*Workflow:* temporal, cadence.<br>*Eventing:* nats, kafka, rabbitmq | `internal/worker/`, `internal/cron/`, `internal/workflow/`, `internal/common/eventing/` |
| **Grup F: Observance** | `add:tracing`, `add:metrics`, `add:profiling`, `add:healthcheck` | *Tracing:* jaeger, tempo, datadog, stdout.<br>*Metrics:* prometheus (RED/USE), datadog, otel-collector.<br>*Profiling:* `--pgo` (pprof & profile-guided optimization).<br>*Healthcheck:* `/health/liveness`, `/health/readiness` | `pkg/telemetry/`, `pkg/middleware/metrics.go`, `cmd/server/router.go` |
| **Grup G: API Design** | `add:docs`, `add:versioning`, `add:webhook` | *Docs:* swagger (ogen codegen), scalar.<br>*Versioning:* `--strategy` (url-prefix/header/content-type).<br>*Webhook:* outbox delivery worker + HMAC retry verification | `docs/openapi.yaml`, `internal/adapter/handler/http/versioning.go`, `internal/webhook/` |
| **Grup H: Patterns** | `add:di`, `add:config`, `add:feature-flags` | *DI:* wire (Google Wire codegen).<br>*Config:* viper env bind with validation.<br>*Flags:* flipt, launchdarkly, env | `cmd/server/wire.go`, `internal/config/`, `internal/common/featureflags/` |
| **Grup I: Deploy & Cloud**| `add:deploy`, `add:cloud`, `add:storage`, `add:discovery` | *Deploy:* k8s (Deployment+Service+HPA), lambda (SAM), helm.<br>*Cloud:* aws, gcp.<br>*Storage:* s3, gcs, minio, local.<br>*Discovery:* etcd, consul | `deploy/`, `charts/`, `internal/adapter/storage/`, `internal/common/discovery/` |
| **Grup J: Comm / Net** | `add:mailer`, `add:notification`, `add:httpclient` | *Mailer:* smtp, sendgrid, mailgun, ses.<br>*Notif:* fcm, apns, web-push (VAPID).<br>*HTTPClient:* `--retry`, `--circuit-breaker`, `--timeout`, `--tracing` | `internal/adapter/mailer/`, `internal/adapter/notification/`, `pkg/httpclient/` |
| **Grup K: AI & Data Infra**| `add:ai` | embedding (pipeline scaffold), llm-proxy (OpenAI/Ollama inference compatible) | `internal/adapter/ai/` |
| **Grup L: Testing Harness**| `add:test` | integration (testcontainers-go), e2e (httptest), fuzz (Go native fuzz), property (pgregory.net/rapid), synctest (Go 1.24 deterministic) | `tests/integration/`, `tests/e2e/`, `*_test.go` |
| **Grup M: Code Quality** | `add:validation`, `add:lint` | *Validation:* playground (validator v10), ozzo.<br>*Lint:* custom (go/analysis linter framework) | `pkg/validator/`, `tools/linter/` |
| **Grup N: Utilities & Audit**| `add:multitenancy`, `doctor`, `ls`, `upgrade` | *Tenancy:* `--strategy` (rls/schema-per-tenant/subdomain).<br>*Doctor:* `--fix` (validasi dan auto-fix inkonsistensi aether.yaml vs disk).<br>*Ls:* Tampilkan tabel matriks modul terdaftar.<br>*Upgrade:* Migrasikan skema aether.yaml | `internal/common/tenancy/`, `aether.yaml` |

---

### 3.4 Zero-Exception Anomaly Hunting (16 Edge Cases & Operational Resilience)
Sesuai amanat pilar keempat, seluruh skenario kegagalan operasional (*edge cases*) tidak berhak diringkas. Berikut adalah 16 protokol pertahanan penuh terhadap anomali sistem dan lingkungan yang telah dirancang untuk `go-aether`:

| No | Nama Kasus Tepi / Anomali | Skenario Pemicu di Lapangan | Strategi Mitigasi Arsitekural & Code Guardrails |
| :---: | :--- | :--- | :--- |
| **1** | **`aether.yaml` Missing in Subdir** | User berada di `internal/core/service` dan mengeksekusi `go-aether make:handler order`. | CLI mengadopsi algoritma pemecahan **Walk-Up Recursive Directory Search**: menelusuri dari direktori eksekusi `./` naik ke parent `../`, `../../`, hingga mendeteksi file `aether.yaml` atau batas repositori `.git/`. Jika buntu di root OS, lempar Sentinel Error yang membeberkan cara run `init` atau `adopt`. |
| **2** | **Target File Conflict / Already Exists** | Command dijalankan ulang pada nama modul yang sama; file `order.go` sudah berdiam di filesystem. | Bertindak dengan proteksi interaktif: menampilkan menu konfrontasi bersyarat (*Skip / Overwrite / Diff / Backup `.bak` / Abort*). Jika disertai flag `--force`, timpa otomatis. Jika dibubuhi `--dry-run`, hanya cetak perbedaan komparasi delta linier melalui layar konsol tanpa mutasi disk sepotong pun. |
| **3** | **Invalid Go Identifier Naming** | User mengetik nama modul yang melanggar hukum sintaktis Golang, misalnya: `my-order` atau `212fast`. | Validasi Regex fail-fast mencekik eksekusi sebelum memory template dialokasikan. Konsol melemparkan `ErrInvalidIdentifier` beserta 3 saran normalisasi konversi sah seketika: `myorder`, `my_order`, atau `MyOrder`. |
| **4** | **Duplicate Module Declaration** | Entri `order` telah tertidur di dalam array `modules` dalam file `aether.yaml`. | Engine memveto pembuatan modul ulang secara mutlak untuk mencegah penindihan yang berujung amnesing kode. Instruksikan pengguna untuk mengalihkan sasaran kepada `make:handler order --transport grpc` guna menambahkan cabang antarmuka pada modul tersebut, atau tambahkan `--force` untuk regenerasi penuh. |
| **5** | **Brownfield BYO Dependency Collision** | Proyek lawas telah membina koneksi global `DB` atau `Redis` di `main.go`. Generator murni yang awas bisa menciptakan koneksi ganda yang memporak-porandakan batas *connection pool* (max conns exceeded). | Jika parameter `adapters.existing_db_var` / `adapters.existing_redis_var` di dalam `aether.yaml` berisi nilai (bukan string kosong), AST engine segera mengabaikan penulisan inisialiser koneksi baru (`sqlx.Connect` atau `pgxpool.New`). Sebaliknya, ia menyambung jalur parameter **Constructor Injection** pada blok konstruktor (*NewOrderRepository(db *pgxpool.Pool)*), sehingga *main.go* lawas dapat menyulap instance yang sudah ada tanpa pertengkaran sesi koneksi. |
| **6** | **Non-Standard Legacy Folder Layouts** | Arsitektur Brownfield berserakan tak tentu arah (misalnya: `web/controllers/`, `logic/`, `data/`). | Saat user memicu `go-aether adopt --scan`, kuesioner interaktif TUI merekam pemetaan alamat aneh tersebut ke blok `architecture.paths.*` pada manifest dengan flag `anomaly_mode: true`. Generator selanjutnya beroperasi dengan **kepatuhan buta pada manifest kustom ini**, menyemayamkan file-file baru di lokasi direktori nyentrik yang ditentukan tanpa mengubah struktur lama. |
| **7** | **Architecture Style Mismatch** | Manifest `aether.yaml` menyatakan `pattern: hexagonal`, tetapi user mengintervensi dengan `make:handler order --style mvc`. | Mesin penjaga integritas mengunci upaya intervensi inkonsisten ini seketika dengan `ErrArchitectureMismatch`. Menawarkan saran untuk tetap setia pada jalur Hexagonal, atau melakukan migrasi skema perancang via perintah resmi `go-aether upgrade --arch mvc` (disertai peringatan bahwa file eksisting menuntut transliterasi manual). |
| **8** | **Outdated Go Compiler Minimum Version** | Mesin komputasi user menginduki `go.mod` bernomorkan Go lama (misal: `1.18`), padahal pustaka log default yang disandarkan mengadopsi `log/slog` ( Go 1.21+ ). | CLI menganalisis header berkas `go.mod` sesaat pasca-parsing command. Jika versi di bawah `1.21`, lempar panic edukatif berisikan baris perintah remediasi konsol aktual (`go mod edit -go=1.21 && go mod tidy`). |
| **9** | **Partial Disk Write / Permission Failure** | Proses pembicara harus menuliskan 7 berkas beruntung, namun tumbang sehabis menulisberkas ke-5 atas kendala disk penuh (*ENOSPC*) atau kegagalan kepemilikan akses (*EACCES*). | Mengandalkan arsitektur **In-Memory Transactional Write Buffer Pattern**: Seluruh hasil olahan template ditimbun sejenak dalam buffer memori sementara. Hanya saat keseluruhan 7 berkas terkonfirmasi bebas cacat baru dilimpahkan secara *atomic sequence* ke disk nyata. Jika satu penulisan terhadang kendala I/O, mekanisme peniti balikan (*Atomic Rollback*) langsung menghapus berkas bernomor 1 s/d 4 yang baru saja bergeming di filesystem agar tidak meninggalkan sampah setengah matang. |
| **10**| **Missing Go Module Declaration** | Eksekusi dipaksa beroperasi di atas ruang folder hampa yang lalai di-initiate dengan berkas `go.mod`. | CLI menonaktifkan proses generator dan menayangkan panduan langkah *prerequisite*: perintahkan instalasi `go mod init <repo-url>` secara lugas sebelum mengetam gerang eksekusi cadangan scaffolding. |
| **11**| **Uninstalled Formatter binaries in PATH** | Pasca merumahkan file ke hard disk, CLI bermaksud menuntaskan indentasi dengan `goimports` atau `gofmt`, namun utilitas tersebut absen di `ENV:PATH`. | Mengaktifkan pendekatan **Graceful Degradation**: Pembuatan modul dijamin terselesaikan lancar 100% tanpa fatal crash. Namun, terminal menerbitkan notifikasi kuning bercahayakan peringatan `WARN: goimports missing in PATH`, menyerta tautan perintah unduh cepat (`go install golang.org/x/tools/cmd/goimports@latest`). |
| **12**| **Standard Library Package Name Collision** | User bermimpi buruk dengan menamakan modul berselingan kata terlarikan baku: `go-aether make:module context`, `sync`, atau `http`. | Sistem meraba tabel inventaris reserved keyword milik standard library Go. Apabila kata beradu pantul, lempar `ErrStdlibCollision` serta mampat luput luaskan pengetikan baru berawalan konvensi semantik seperti `appcontext` atau `reqcontext`. |
| **13**| **Outdated Schema version in Manifest** | Pengguna meningkatkan binary `go-aether` ke `v1.0.0`, akan tetapi rekam `aether.yaml` proyek masih menganut parameter arkais versi `v0.9`. | Menerjemahkan mode **Backward Compatibility Protection**: CLI `v1.0` senantiasa mempertahankan pembaca arsip lawas secara lancar demi mengelakkan kerusakan fungsi eksisting, sambil mendesak peringatan interaktif untuk menyetujui eksekusi normalisasi melalui perintah autoconvert `go-aether upgrade`. |
| **14**| **Missing Environment Variable Definitions** | Konfigurasi manifest menyanjung pemakaian `cache: redis`, sebaliknya arsip `.env.example` alfa memuat definisi token kunci `REDIS_URL=""`. | Sewaktu pemantapan detektif `go-aether doctor` mengudara, kepincangan ikatan ini dicantumkan sebagai alarm ⚠️ *Incomplete Env Binding*. Pengguna ditawari kemantapan eksekusi perbaikan mandiri seketika berbekalkan opsi `go-aether doctor --fix`. |
| **15**| **Severe Legacy Spaghetti Entanglements** | Kondisi direktori brownfield sebegitu berantakan dan saling melilit bagaikan akar benalu sehingga mapping path konvensional musnah total dan ditolak sistem build. | Mengerahkan jalur evakuasi darurat (**Emergency Escape Hatch / Vertical Slice Mode**): `go-aether make:feature order --vertical --dir ./features/order`. Perintah ini menyekap seluruh entitas *handler, service, repo, domain, dan test harness* dalam satu bunker direktori isolasi yang kedul dan kebal dari badai kekacauan sekitarnya. |
| **16**| **Cross-Platform Path Separator Anomalies** | Pengarsipan lintas sistem operasi: Windows yang mengenyam aksen backslash (`\`) bertentangan dengan Unix/Linux yang memunguti forward slash (`/`). | Seluruh representasi teks internal di dalam file `aether.yaml` **DIHARAMKAN KETAT MENYIMPAN BACKSLASH**. CLI berkomitmen untuk selalu membakukan pembekakan string melewatinya normalisasi mutlak `filepath.ToSlash()`, sembari meneruskan eksekusi baca-tulis sistem file lokal bersendi antarmuka standar kompilator `filepath.Join()`. |

---

### 3.5 STRIDE Threat Model & Security Perimeter (CLI & Generated Code)

Sistem yang hebat bukan cuma indah, melainkan juga harus kebal terhadap eksploitasi cyber. Berikut adalah audit analisis ancaman STRIDE untuk `go-aether`:

| Vektor Ancaman STRIDE | Potensi Celah & Skenario Serangan | Mitigasi Arsitektural & Guardrails `go-aether` |
| :--- | :--- | :--- |
| **Spoofing** *(Pemalsuan identitas)* | Pemasang binary palsu atau eksekusi arsip template jahat dari luar yang memanipulasi *scaffold target*. | Seluruh template dibungkus matang secara statis ke dalam biner resmi via `//go:embed`. Rilis binari diamati oleh enkripsi penandatanganan kriptografik Cosign di pipeline GitHub Actions. |
| **Tampering** *(Manipulasi arsip)* | Path traversal attack pada saat pemanggilan argumen nama modul: `go-aether make:module ../../../etc/passwd` atau penyanderaan skrip pada spesifikasi file `aether.yaml`. | Regex pengukuh nama (*Strict Identifier Sanitization*) membuang seluruh karakter slsh (`/`, `\`, `.`, `..`) sebelum mengolah resolusi string filesystem. Parsing YAML dibatasi skema tipe data kuat dari *Viper* bertutup rintangan eksekusi shell. |
| **Repudiation** *(Penyangklangan aksi)* | Pengguna meretas rancangan struktur, menimpakan kerusakan pada kode lama pasca integrasi, lalu mengklaim ketidaktahuan atas mutasi arsip. | `go-aether` mencubiti penjejaking komentar penanda historis tanpa intervensi logika di akhir tajuk berkas (*Code generation metadata stamp*): `// Code generated by go-aether v0.1.0 on YYYY-MM-DD; DO NOT EDIT EXCEPT IMPLEMENTATION BLOCKS`. |
| **Information Disclosure** *(Kebocoran rahasia)* | Generator menghasilkan file `.env` atau *test files* yang mengandung string token kata sandi autentik internal yang tak sengaja terunggah ke repositori Git. | Generator menertibkan asas **Zero Hardcode Secret (Gerbang §0.3)**: Output yang terbit hanyalah arsip purwa `.env.example` yang dikunci nihil nilai rahasianya, serentak membenamkan pola pelarangan eksklusi `.env*` yang solid pada cetakan file `.dockerignore` dan `.gitignore`. |
| **Denial of Service** *(Beban komputasi)* | Pemaksaan pembuatan struktur bersarang tak terhingga (*Infinite Recursive Folder Loom*) atau manipulasi YAML berspekan besar yang memicu konsumsi RAM meledak (*YAML Bomb*). | Batas kedalaman penjelajahan pembaca direktori (*Directory Scanning Depth Limit*) diborgol kaku maksimal 4 tingkat. Ukuran maksimal muatan yang sanggup dijamar dari berkas `aether.yaml` dibatasi seberat 256 Kilobytes untuk mencegah malapetaka kelebihan kapasitas heap. |
| **Elevation of Privilege** *(Bypass spesifikasi)* | Prosedur generate menciptakan modul *HTTP handler* dengan celah terbuka yang melalaikan autentikasi middleware dan validasi input. | Template penutupan handler senantiasa merangkul kerangka enkapsulasi protektif (*Safe-by-Default Wiring*): Menyisipkan middleware pengaman `RequestID`, `Recoverer`, serta pengait fungsi validasi sematik input terstruktur yang tak dapat diabaikan oleh pengguna tanpa penghapus sengaja. |

---

## 4. Rencana Eksekusi & Living Task Checklist (DAG & Batch Protocol §7)

Untuk meminimalkan risiko dekontekstualisasi dan menjamin determinisme implementasi dari nol hingga MVP v0.1.0, seluruh pekerjaan akan dieksekusi mematuhi hukum **DAG (Directed Acyclic Graph)** dalam pembatalan bertahap (*Batch Writing Protocol*):

### Batch 1: Core Domain Locks & Filesystem Port Abstraction
* **DependsOn:** `[]` *(Fondasi Inti)*
* **Goal:** Mendefinisikan entitas utama, arketipe skema manifest, antarmuka port, dan sistem penguncian error yang absolut sebelum beranjak ke eksekusi masukan CLI.
- [ ] `[NEW]` `go-aether/internal/core/domain/errors.go` *(Daftar konstelasi seluruh Sentinel Errors)*
- [ ] `[NEW]` `go-aether/internal/core/domain/manifest.go` *(Struct type-safe untuk spesifikasi aether.yaml)*
- [ ] `[NEW]` `go-aether/internal/core/domain/module.go` *(Struct arsitektur, parameter stack, dan entitas modul)*
- [ ] `[NEW]` `go-aether/internal/core/port/generator.go` *(Antarmuka kontrak IGenerator, IManifestResolver, dan IFileWriter)*
- [ ] Verifikasi & Unit Test Contract Batch 1 via `go test -race -v ./internal/core/...` **(Target: Pass 100%)**

### Batch 2: Template Rendering Engine & Filesystem Adapters
* **DependsOn:** `[Batch 1]`
* **Goal:** Mewujudkan mesin pengeksekusi template berlandaskan `text/template` dengan bundling `embed.FS` serta penyambungan adaptor sistem berkas terisolasi berprinsipkan *Transactional Rollback Buffer*.
- [ ] `[NEW]` `go-aether/internal/adapter/template/engine.go` *(Implementasi parser text/template + pemetaan variabel TemplateData)*
- [ ] `[NEW]` `go-aether/internal/adapter/manifest/yaml_resolver.go` *(Implementasi Viper reader dengan mekanisme Walk-Up Locator)*
- [ ] `[NEW]` `go-aether/internal/adapter/writer/afero_writer.go` *(Implementasi IFileWriter dengan perlindungan Atomic Transaction Buffer)*
- [ ] `[NEW]` `go-aether/templates/common/aether_yaml.tmpl` *(Template dasar pengisi file manifest aether.yaml)*
- [ ] Verifikasi Adapter & Mock Test FS Batch 2 via `afero.NewMemMapFs()` **(Target: Pass 100%)**

### Batch 3: Core Hexagonal Templates Generation & Service Layer Wiring
* **DependsOn:** `[Batch 2]`
* **Goal:** Menjalur rangkaian cetakan Hexagonal murni (Domain, Port, Service, Chi Handler, Postgres SQLC Stub) dan meramu orchestration service pembimbing eksekusi (`InitService`, `MakeService`).
- [ ] `[NEW]` `go-aether/templates/hexagonal/{domain,port,service,handler_http,repository_postgres}.go.tmpl` *(Cetakan murni Hexagonal)*
- [ ] `[NEW]` `go-aether/internal/core/service/init_service.go` *(Orkestrasi alur penciptaan proyek Greenfield)*
- [ ] `[NEW]` `go-aether/internal/core/service/make_service.go` *(Orkestrasi alur injeksi modul penopang ke dalam repositori)*
- [ ] `[NEW]` `go-aether/internal/core/service/doctor_service.go` *(Mesin diagnosa kesehatan struktur dan keutuhan spesifikasi manifest)*
- [ ] Verifikasi Service & Golden File Integration Tests Batch 3 **(Target: Pass 100%)**

### Batch 4: CLI Entrypoint, Cobra Subcommands & Interactive UX
* **DependsOn:** `[Batch 3]`
* **Goal:** Menyelesaikan pintu gerbang utama binary berwujudkan subcommand `cobra` dan hiasan UX indikator pembaca terminal.
- [ ] `[NEW]` `go-aether/cmd/go-aether/main.go` *(Root entrypoint pembakar aplikasi)*
- [ ] `[NEW]` `go-aether/internal/adapter/cli/root.go` *(Inisialisasi Cobra root command & penetapan flag verbose)*
- [ ] `[NEW]` `go-aether/internal/adapter/cli/{init,make,adopt,doctor,ls,upgrade}.go` *(Penghubung argumen konsol menuju penangan di Service layer)*
- [ ] `[NEW]` `go-aether/Makefile` *(Task runner automasi kompilasi, verifikasi, dan penyembah golden files)*
- [ ] Verifikasi Sanity E2E Binary Generation Test & Linter Full Suite Batch 4 **(Target: Pass 100%)**
- [ ] *(Opsional)* Evaluasi Adversarial QA Subagent atas ketahanan terhadap kasus tepi **(Target: Pass 100%)**

---

## 5. Observabilitas & Verifikasi Mutu (Quality Gate §6)

### 5.1 Rencana Observabilitas & Debug Logging Internal CLI
Mengoperasikan generator tanpa observabilitas adalah undangan menuju keputusasaan saat terjadi salah penempatan path di mesin CI/CD. CLI ini dibekali sistem jejak komersial berbasis **`log/slog`**:
- **Default Mode:** Keluaran konsol bergaya manusiawiwi (*human-readable output*) bertakdir minimalis. Hanya menampilkan indikator sukses ✅, lompat file ⏭️, atau kesalahan fatal ❌.
- **Verbose Mode (`--verbose` / `-v`):** Menjemput peralihan struktur menuju keluaran JSON terstruktur yang mengejutkan setiap langkah resolusi path, durasi alokasi RAM per template, serta rekapan inspeksi *tree-sitter/AST*:
  ```json
  {"time":"2026-08-06T20:00:00Z","level":"DEBUG","msg":"resolving manifest path","cwd":"C:\\Projects\\app","steps_up":2,"found":"C:\\Projects\\aether.yaml","duration_ms":1.4}
  {"time":"2026-08-06T20:00:00Z","level":"INFO","msg":"rendering template","template":"hexagonal/service.go.tmpl","target":"internal/core/service/order_service.go","status":"written"}
  ```

### 5.2 Rencana Pengujian (Testing Harness)
Ketaatan pada Quality Gate §6 berlanjut dengan pengasahan metodologis multi-lapis tanpa belas kasih:
1. **Unit & Race Testing:** `go test -race -coverprofile=coverage.out ./...`. Ambang batas kecakapan pengujian (*Coverage Threshold Gate*) dibekukan di **minimum 80%**. Tidak diperkenankan melepaskan build jika terjungkal di bawahnya.
2. **Deterministic Golden File Testing:** Setiap pemutakhiran pada berkas `.tmpl` diharuskan menjalani verifikasi perbandingan identitas keluaran terhadap arsip "Golden File" acuan di `tests/golden/`. Bila memang sengaja merubah konvensi, pengesahan wajib dipicu via `go test -update-golden`.
3. **In-Memory FS Isolation Test:** Menyematkan pustaka `afero.NewMemMapFs()` untuk mensimulasikan lingkungan disk ekstrem (seperti kuota penuh atau pemetaan izin ilegal 000) tanpa berisiko menampar hard drive asli milik developer.
4. **Mutation Testing Verification:** Pengoperasian perkakas mutasi kode (`go-mutesting`) untuk membakar seluruh asersi *tautologis* gampang (misalkan sekadar menguji `assert.NotNil(err)`). Test harus gagal/merah seketika apabila baris kode internal generator diubah secara asal.

### 5.3 Self-Review & Exit Gate Checklist (Ratchet Enforcement)
Sebelum tanda tangan akhir `[x]` diberikan pada setiap tahapan, komandan sistem wajib mencecar barisan parameter ini:
```text
N+1 query: <nihil — generator bekerja dari memory ke disk tanpa jaringan>
OWASP/STRIDE: <verified — sanitasi nama modul & eliminasi rahasia terpasang>
race condition: <nihil — arsitektur komputasi CLI bersifat sekuensial bebas race data>
input validation: <verified — regex pembendung nama modul ilegal terpasang ketat>
memory leak: <nihil — pembakaran buffer dibatasi <256KB per pemanggilan siklus>
goroutine leak: <nihil — tidak ada goroutine menggantung purna eksekusi CLI>
```

---

## 6. Prosedur Rollback & Canary Strategy (Safe Deployment & Restoration)

### 6.1 Canary Release & Version Adoption Strategy
Mendistribusikan alat bantu CLI ke seluruh engsel organisasi harus berlangsung bertahap agar tidak menghancurkan keutuhan repositori yang ada:
- **Phase 1 (Alpha Dogfooding):** Binary `v0.1.0` hanya diujicobakan secara eksklusif untuk melestarikan dan membangun arsitektur dalam repo `go-aether` itu sendiri (*self-hosting dogfooding*).
- **Phase 2 (Greenfield Pilot Canary):** Dikerahkan pertama kali pada 2 atau 3 layanan baru (*microservices*) berkepencantuman non-kritis (*low-risk blast radius*) dengan pemantauan keluhan dari developer.
- **Phase 3 (Brownfield Organizational Ramp-Up):** Diperkenalkan ke repositori raksasa lawas secara konsolidatif via fitur interaktif `go-aether adopt --scan`.

### 6.2 Pemicu Rollback Darurat (Rollback Triggers & Emergency Circuit Breakers)
Sistem pembakaran *release* dinyatakan gugur dan wajib dimundurkan segera bila ditemui hal berikut:
1. **Kegagalan Kompilasi Output:** Proyek hasil bentukan `go-aether init` atau `make:module` menolak untuk diuji bangun oleh kompilator asli `go build ./...` atau menyorotkan kegagalan validasi statis pada linter `golangci-lint run`.
2. **Koruspsi Manifest Brownfield:** Perintah adopsi meremas arsip eksisting atau menggandakan deklarasi yang mengakibatkan mogoknya alur kerja CI/CD tim.

### 6.3 Prosedur Remediasi Kilat (Instant Recovery Action Plan)
Jika kiamat lokal benar-benar meledak pasca pemanggilan generator di workstation seorang developer:
1. **Pemanfaatan Atomic Git Revert:** Mengandalkan kebersihan status Git sebelum eksekusi. Bila hasil generate tercemarkan, cukup luncurkan `git clean -fd && git checkout -- .` untuk mengembalikan repositori ke kondisi murni seketika.
2. **Rollback Versi Tooling:** Mundurkan versi tool ke stabilan rilis sebelumnya:
   ```bash
   go install github.com/muhananaufal/go-aether@v0.0.9
   ```
3. **Rekam Fakta Insiden:** Buka arsip Post-Mortem baru di bawah `docs/rca/YYYYMMDD-generator-regression.md` dan abadikan perbaikan di barisan bukti `LEARNED.md`.

---
*Dokumen perancangan ini telah disusun secara komprehensif mematuhi 4 Pilar standar Principal L8 tanpa setitik pun pemotongan token.*
**MENUNGGU PERINTAH MUTLAK DARI USER: KETIK "GASSKAN" UNTUK MEMULAI EKSEKUSI BATCH 1!** 🛡️🚀
