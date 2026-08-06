<p align="center">
  <h1 align="center">🚀 go-aether</h1>
  <p align="center">
    <strong>Lightning-fast, Zero-Runtime Dependency, Opinionated Architecture Scaffold CLI Engine for Go.</strong>
  </p>
  <p align="center">
    <a href="https://github.com/muhananaufal/go-aether/actions"><img src="https://img.shields.io/github/actions/workflow/status/muhananaufal/go-aether/ci.yml?branch=main&style=flat-square" alt="Build Status"></a>
    <a href="https://golang.org/doc/devel/release.html"><img src="https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat-square&logo=go" alt="Go Version"></a>
    <a href="https://github.com/muhananaufal/go-aether/releases"><img src="https://img.shields.io/github/v/release/muhananaufal/go-aether?style=flat-square" alt="Release"></a>
    <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square" alt="License"></a>
  </p>
</p>

---

**`go-aether`** is an advanced developer tool designed for Go Backend Engineers. It eliminates the repetitive boilerplate of setting up Hexagonal (Ports & Adapters) architectures by scaffolding clean, structurally sound, and production-ready modules in milliseconds.

Unlike heavy web frameworks, `go-aether` is **strictly a Dev-Time CLI**. It embeds standard `text/template` files into a single binary, meaning the generated code uses pure Go standard libraries and your chosen dependencies (like `chi` or `pgx`) with **absolute zero runtime lock-in** to `go-aether` itself.

## ✨ Features

- 🏗️ **Hexagonal Architecture Native:** Enforces Domain, Port, Service, and Adapter layer isolation automatically.
- ⚡ **Zero Runtime Bloat:** Generated code belongs to you. No hidden framework dependencies.
- 📦 **Vertical Slice Scaffolding:** Generate a full feature (Domain entity, interface contracts, service use-cases, HTTP handlers, and Postgres repos) with a single command.
- 🛡️ **Transactional Disk Buffer:** If a file generation fails midway, the CLI rolls back all writes automatically to prevent a corrupted Git tree.
- 🩺 **Aether Doctor:** Built-in structural diagnostics to ensure your project's integrity hasn't drifted from the manifest.

## 🚀 Installation

Install the latest version seamlessly via the Go toolchain:

```bash
go install github.com/muhananaufal/go-aether@latest
```
*(Ensure `$(go env GOPATH)/bin` is appended to your `$PATH`)*

## 📖 Quick Start

### 1. Initialize a New Greenfield Project
Run `init` to bootstrap your base architecture directory and generate the Single Source of Truth manifest (`aether.yaml`).

```bash
mkdir my-enterprise-api
cd my-enterprise-api

# Bootstrap Hexagonal Architecture (defaults to chi router & postgres)
go-aether init my-enterprise-api
```

### 2. Scaffold a Vertical Module
Generate a complete feature slice (e.g., `user` or `order`). `go-aether` handles the wiring of Domain logic, Interface Ports, Service Orchestrators, HTTP Handlers, and DB Repositories.

```bash
go-aether make:module order --transports http
```

### 3. Verify Structural Health
Run the built-in diagnostic tool to validate your project's compliance against the `aether.yaml` schema.

```bash
go-aether doctor
```

## 📂 Generated Architecture Layout

When you scaffold a module (e.g., `order`), `go-aether` injects the code into the strict Hexagonal layers:

```text
my-enterprise-api/
├── aether.yaml                     <-- Project Single Source of Truth
└── internal/
    ├── core/                       <-- Inner Layer (Independent)
    │   ├── domain/
    │   │   └── order.go            <-- Pure Domain Entities & Validation
    │   ├── port/
    │   │   └── order_port.go       <-- Inbound & Outbound Interfaces
    │   └── service/
    │       └── order_service.go    <-- Business Logic & Orchestration
    └── adapter/                    <-- Outer Layer (Dependencies)
        ├── handler/http/
        │   └── order_handler.go    <-- HTTP Transport Controller
        └── repository/
            └── order_repository.go <-- DB Persistence Implementation
```

## 🛠️ Commands Reference

| Command | Description | Example |
|---|---|---|
| `init` | Bootstraps a new repository layout and manifests. | `go-aether init api` |
| `make:module` | Generates a full domain slice across all layers. | `go-aether make:module payment` |
| `doctor` | Runs structural diagnostics and schema validation. | `go-aether doctor` |
| `adopt` | Adopts an existing legacy project into aether. | `go-aether adopt --scan` |

### Global Flags
- `--dry-run`: Simulates the command without actually writing files to the disk.
- `-v, --verbose`: Enables verbose debug logging.

## 🤝 Contributing

We welcome contributions! Please feel free to submit a Pull Request or open an Issue to discuss potential changes.
All tests must pass in the CI/CD pipeline before merging.

## 📜 License

This project is licensed under the [MIT License](LICENSE) © 2026 muhananaufal.

---
<p align="center">
  <i>Built with precision for Go Artisans by AETHERIS.</i>
</p>
