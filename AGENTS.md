# AI Agents Guide for go-aether

Welcome, AI Agent! If you are reading this, you are working on the `go-aether` repository. 
This is a high-performance, developer-first CLI engine designed for modern Go Backend Engineers. 
It generates idiomatic Go code (Hexagonal Architecture, DDD, CQRS) strictly as a Dev-Time Scaffolding Tool.

**CRITICAL RULE: DO NOT break the architecture and taxonomy.**

## 1. Project Overview & Architecture
- **Language**: Go (v1.23+)
- **Paradigm**: The CLI is built on top of `cobra`. It uses `text/template` files embedded in the binary to scaffold code.
- **Root Entrypoint**: `main.go` -> `internal/adapter/cli/root.go`
- **Generators**: Located in `internal/core/service/` (e.g., `make_service.go`, `add_service.go`).
- **Templates**: Located in `templates/` (must use embedded `//go:embed`).

## 2. Command Taxonomy (Domain-First)
We use a **Domain-First** taxonomy for all commands. **DO NOT** use verb-first (`add:*`, `make:*`) commands.

| Domain Prefix | Purpose | Example |
|---|---|---|
| `arch:*` | Architecture scaffolding (Domain, Port, Service, Handler, Repo) | `arch:module` |
| `db:*` | Database & persistence | `db:migration`, `db:uow` |
| `api:*` | API layer (HTTP, gRPC, GraphQL) | `api:graphql`, `api:transport` |
| `realtime:*` | Realtime protocols | `realtime:ws`, `realtime:sse` |
| `async:*` | Async messaging & events | `async:worker`, `async:saga` |
| `security:*`| Auth, zero-trust, cryptography | `security:auth`, `security:argon2` |
| `cache:*` | Caching layer | `cache:redis`, `cache:multilevel` |
| `infra:*` | DevOps, deployments, linting | `infra:deploy`, `infra:healthcheck` |
| `o11y:*` | Observability (Tracing, Metrics, Logs) | `o11y:tracing`, `o11y:metrics` |
| `fintech:*` | Financial reliability | `fintech:ledger` |
| `notif:*` | Notification & communication | `notif:mail`, `notif:sms` |
| `cloud:*` | Cloud storage | `cloud:s3` |
| `platform:*`| Distributed platform components | `platform:cqrs`, `platform:discovery` |
| `test:*` | QA engine | `test:integration`, `test:fuzz` |
| (Flat) | Core lifecycle commands | `init`, `adopt`, `doctor`, `ls` |

## 3. Mandatory Rules for AI Agents
1. **Never use `add:*` or `make:*`**: Always use the correct domain prefix (`arch:*`, `api:*`, etc.).
2. **Never add external runtime dependencies to the generated code**: `go-aether` must generate pure Go code without lock-in.
3. **Always run tests after changes**: Use `go test -v ./...` and ensure all tests pass. E2E tests are required for new commands.
4. **Zero Magic**: Code must be explicit, strictly typed, and idiomatic. No reflect magic unless absolutely necessary.
5. **Update README.md**: If you add a new command or flag, you MUST update `README.md` to reflect the changes.
6. **Maintain SSOT**: Ensure that changes to commands are reflected in `internal/adapter/cli/root.go` and the corresponding service interface in `internal/core/port/generator.go`.

## 4. Development Workflow
- Run `go build -o bin/aether.exe .` to compile the binary.
- Use `git add .` and `git commit -m "feat/fix(scope): message"` following Conventional Commits.
- Always include the Co-Authored-By trailer for AETHERIS: `Co-authored-by: aetheris <agents.aetheris@gmail.com>`.
