# Changelog

All notable changes to `go-aether` are recorded here.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [0.4.0] — 2026-08-07

The release that made the tool's own claims verifiable.

Every defect below was found by building the binary and running it, then locked
behind a test so it cannot return. See
[`docs/rfc/20260807-v0.4.0-production-hardening.md`](docs/rfc/20260807-v0.4.0-production-hardening.md)
for the full design and the anomaly matrix.

### ⚠️ Breaking changes

- **`adopt` now previews by default.** It prints the proposed layer mapping and
  writes nothing. Pass `--apply` to save `aether.yaml`. Scripts that relied on
  `adopt` writing immediately must be updated. The command runs against
  repositories somebody already depends on, where a wrong mapping sends every
  later generator into a directory that means something else.
- **Unsupported stack selections are rejected instead of silently substituted.**
  `--arch clean`, `--arch ddd`, `--router mux`, and deploy targets other than
  `k8s` now fail with a message listing what is supported. They previously
  produced a chi + Postgres project regardless of what was asked for.
- **Exit codes are differentiated.** `2` for input the user can fix, `1` for
  operational failure. Previously everything was `1`, so a typo in a flag was
  indistinguishable from a full disk.

### Added

- **Compile gate.** Generated projects are built with the real Go toolchain in
  CI: `go build`, `go vet` and `gofmt -l` over the output. This is the layer that
  had never existed, and its absence is why the defects below shipped.
- **Real multi-framework support.** `--router` accepts `chi`, `gin`, `echo`,
  `fiber` and `stdlib`; `--db` accepts `postgres`, `mysql`, `sqlite` and `none`.
  All 15 combinations are compiled in CI, and `go.mod` is asserted to prove the
  selection actually reached the generated code.
- **A brownfield engine that reads the repository.** `adopt --scan` walks the
  tree (bounded: depth 4, 5000 files, vendor and caches skipped), classifies
  layers from directory naming *and* the import graph, reads the module path from
  `go.mod`, and detects connection pools the entrypoint already builds so
  generated constructors accept them rather than opening a second one.
- **Command maturity tiers.** `--help` marks every command: unmarked means the
  generated code is compiled by CI (8), `[tested]` means generation is tested but
  the output is not compiled (31), `[experimental]` means no automated coverage
  yet (51). Experimental is the default for anything unclassified.
- **`doctor` practice audit.** Ten rules covering exposed `.env` files, unbounded
  connection pools, servers without timeouts, missing probes, panics in library
  packages, detached contexts, unstructured logging, absent tests and missing
  linter config. Every finding names the fix, and where a `go-aether` command
  already provides it, the hint names that command.
- **Regression suite.** Permanent coverage for every defect listed below.
- **`govulncheck` and `golangci-lint`** in CI, and `-race` promoted to required.

### Fixed

- **`init` produced a project that would not build.** `go mod tidy` ran before
  any template had been rendered, so it had no imports to resolve and `go.mod`
  was left empty. Four modules were unresolved on the user's first command.
- **`api:validator` could never succeed.** It staged the same destination twice
  in one transaction, so the unit of work rolled back every time.
- **`cache:redis <unknown>` reported success while writing nothing** and
  persisted a manifest declaring a driver whose implementation did not exist.
- **`api:middleware` mutated the handler even when the middleware could not be
  generated**, leaving a project that referenced a package that was never
  written.
- **Per-layer commands produced code that did not compile.** `arch:domain` and
  `arch:port` rendered `domain_only` and `port_only` templates that had drifted
  from the canonical ones: no `Validate`, no sentinel errors, no service
  contract. Assembling a slice one command at a time failed with six errors. The
  duplicates are deleted; all layer commands now render the canonical templates.
- **Generated Go was never formatted.** `gofmt -l` flagged seven files from a
  fresh `init`. Output now passes through `go/format`, and a template that
  renders invalid syntax fails with the template name and line instead of
  reaching disk.
- **Windows reserved device names were accepted.** `arch:domain con` exited 0
  and printed success while the content went to the `CON` console device and
  vanished. `CON`, `PRN`, `AUX`, `NUL`, `COM0-9` and `LPT0-9` are now rejected on
  every OS, so a module authored on Linux stays checkoutable on Windows.
- **Path traversal was stopped only by accident.** `infra:deploy ../../../x` was
  blocked because the template lookup failed first, and the error echoed the
  path back. `SafeJoin` now enforces the project root explicitly.
- **The project lock was never released on failure.** `PersistentPostRunE` is
  skipped by cobra when a command errors, so every failed command left the flock
  held.
- **`ls` searched upward without limit,** and could describe an unrelated project
  several directories up as if it were this one.
- **`errors.Is` in the file writer.** A `==` comparison against `fs.ErrNotExist`
  never matched, since both `os` and `afero` wrap it, so a missing manifest
  surfaced as a generic read failure.
- **`adopt` invented a module path** (`github.com/adopted/<dirname>`) instead of
  reading `go.mod`, making every generated import wrong from the first file.
- **A `go.mod` saved with a UTF-8 BOM** was reported as not being a Go module.
- **The generated Dockerfile could not build.** It copied `.env`, which
  `.dockerignore` excludes, and would have baked local credentials into an image
  layer had it succeeded. The image now runs as a non-root user, carries
  `ca-certificates` and a `HEALTHCHECK`, and uses an exec-form `ENTRYPOINT` so
  `SIGTERM` reaches the process.
- **CI pinned Go 1.23 while `go.mod` required 1.25,** so every workflow run
  failed before reaching a test, with an error that read like a dependency
  problem.
- **`viper.Unmarshal` could not see environment-only values** for keys without a
  registered default, so `DB_USER` and `DB_PASS` were silently empty in any
  deployment without a `.env` file.
- **`aether.yaml` was read in full before its size was checked.** Now bounded at
  256 KiB, checked before the read.
- **Windows `MAX_PATH`** is detected with a message naming both the limit and the
  registry key that lifts it.

### Changed

- Repositories are generated with a typed pool (`*pgxpool.Pool` or `*sql.DB`)
  and working SQL per dialect, replacing a `db interface{}` field and `// TODO`
  stubs that taught exactly the untyped pattern this tool exists to discourage.
- SQLite uses `modernc.org/sqlite` rather than `mattn/go-sqlite3`: the latter
  requires cgo, which the generated Dockerfile disables, so the project would
  build on a laptop and fail inside the container.
- Connection pools set explicit limits. The defaults are unbounded, and the
  ceiling is otherwise discovered when the database starts refusing connections.
- Handler helpers are methods rather than package-level functions, so generating
  a second module into the same package no longer collides.
- Skeleton generation is transactional; a failure midway leaves nothing behind.
- `go mod tidy` is non-fatal and bounded by a timeout, so a developer without
  network still gets a complete source tree plus the command to run later.
- `--arch`, `--db` and `--router` help text is generated from the same sets the
  validator enforces, so it cannot advertise a value that would be rejected.
- Releases are reproducible (`-trimpath`, fixed `mod_timestamp`) and ship an
  SBOM per archive.

### Removed

- `templates/hexagonal/repository_redis.go.tmpl` — unreachable, since `redis` was
  never a supported driver, and it implemented `Create`/`Get` against a port
  requiring `Save`/`FindByID`, so it could never have compiled.
- `templates/hexagonal/{domain,port,repository}_only.go.tmpl` — superseded by the
  canonical layer templates.
- `templates/common/{main,postgres}.go.tmpl` — superseded by the per-router and
  per-driver variants.

### Known gaps

Stated rather than rounded up:

- `api:middleware` refuses non-chi projects. Its templates carry chi's
  middleware signature; router-aware middleware is not implemented.
- The `--vertical` escape hatch for unrecognisable legacy layouts does not exist.
- `SafeJoin` guards the three call sites where free-form input reaches a
  filename. Enforcement at the writer layer is still open.
- 51 of 90 commands have no automated coverage. They are labelled.

---

## [0.3.0] — 2026-08-07

Executable skeleton generation and type-safe config scaffolding.

## [0.2.0] — 2026-08-06

Interactive TUI fallback, atomic unit-of-work disk buffer, context memory.

## [0.1.0] — 2026-08-06

Initial release: hexagonal DDD core scaffolding engine.
