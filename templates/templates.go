package templates

import "embed"

// FS is the global embedded filesystem containing all go-aether scaffolding templates.
// It is strictly injected into the stdlib text/template engine at runtime to guarantee
// a self-contained, single-binary distribution with zero external dependencies.
//
//go:embed common/* hexagonal/* distributed/* cloud/* ai/*
var FS embed.FS
