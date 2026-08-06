# Contributing to go-aether

First off, thank you for considering contributing to go-aether! We welcome contributions from the community.

## 1. Local Development Setup

To run go-aether locally:

1. Clone the repository: git clone https://github.com/muhananaufal/go-aether.git
2. Navigate into the directory: cd go-aether
3. Build the binary: go build -o bin/aether.exe .
4. Run tests: go test -v ./...

## 2. Architecture Guidelines

go-aether is built with a strict separation of concerns following Hexagonal Architecture:
- internal/adapter/cli: Handles the command-line interface via Cobra.
- internal/adapter/template: Reads and compiles embedded 	ext/template files.
- internal/core/port: Interfaces (Contracts) connecting adapters and services.
- internal/core/service: Business logic for scaffolding modules and components.
- 	emplates/: Contains all the raw templates used for generation.

**Important**: We use a Domain-First Command Taxonomy (e.g., rch:module, pi:middleware). Please read AGENTS.md before making any structural changes to the commands.

## 3. How to Submit a Pull Request

1. **Fork the repository** and create your branch from main.
2. **Write tests** for your new generator or bug fix.
3. Ensure your code passes the build and all tests pass (go test -v ./...).
4. Make sure your commit messages follow the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) standard.
5. Create a Pull Request describing your changes in detail.

## 4. Code of Conduct

Please note that this project is released with a Contributor Code of Conduct. By participating in this project you agree to abide by its terms. Be respectful and constructive.
