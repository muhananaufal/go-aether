package service

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strings"
)

// RunDoctor performs structural health diagnostics and configuration checks against aether.yaml.
func (s *AetherScaffoldService) RunDoctor(ctx context.Context, startDir string, fix bool, out io.Writer) error {
	manifest, manifestPath, err := s.resolver.Resolve(ctx, startDir)
	if err != nil {
		_, _ = fmt.Fprintf(out, "❌ [CRITICAL] Manifest resolution failed: %v\n", err)
		return err
	}

	_, _ = fmt.Fprintf(out, "🩺 go-aether Doctor Diagnostics for project: %s\n", manifest.Project.Name)
	_, _ = fmt.Fprintf(out, "   Manifest location: %s\n", manifestPath)
	_, _ = fmt.Fprintf(out, "   Architecture mode: %s (%s)\n", manifest.Architecture.Pattern, manifest.Architecture.Mode)

	// Validate manifest structural compliance
	if err := manifest.Validate(); err != nil {
		_, _ = fmt.Fprintf(out, "❌ [ERROR] Manifest Schema Invalid: %v\n", err)
		return err
	}
	_, _ = fmt.Fprintf(out, "✅ Manifest Schema Validation: PASS\n")

	// Audit registered feature modules
	_, _ = fmt.Fprintf(out, "ℹ️ Registered Modules count: %d\n", len(manifest.Modules))
	projectRoot := filepath.Dir(manifestPath)
	var corruptCount int

	for _, mod := range manifest.Modules {
		domainPath := filepath.Join(projectRoot, manifest.Architecture.Paths.Domain, fmt.Sprintf("%s.go", strings.ToLower(mod.Name)))
		exists, err := s.fs.Exists(domainPath)

		status := "✅"
		if err != nil || !exists {
			status = "❌ (MISSING FILE)"
			corruptCount++
		}

		_, _ = fmt.Fprintf(out, "   %s Module [%s]: Transports=%v, Cache=%t\n", status, mod.Name, mod.Transports, mod.HasCache)
		if status != "✅" {
			_, _ = fmt.Fprintf(out, "      -> Expected physical file not found: %s\n", domainPath)
		}
	}

	if corruptCount > 0 {
		_, _ = fmt.Fprintf(out, "⚠️ [CRITICAL] %d module(s) are registered in aether.yaml but missing physical domain files!\n", corruptCount)
	}

	if manifest.Meta.AnomalyMode {
		_, _ = fmt.Fprintf(out, "⚠️ [NOTE] Anomaly mode active for legacy brownfield structure.\n")
	}

	_, _ = fmt.Fprintf(out, "\n🔍 Environment Dependencies Check:\n")

	dependencies := []string{"go", "sqlc", "mockgen", "docker"}
	allDepsMet := true
	for _, dep := range dependencies {
		if path, err := exec.LookPath(dep); err == nil {
			_, _ = fmt.Fprintf(out, "   ✅ %-10s found at %s\n", dep, path)
		} else {
			_, _ = fmt.Fprintf(out, "   ❌ %-10s NOT FOUND in PATH\n", dep)
			allDepsMet = false
		}
	}

	if !allDepsMet {
		_, _ = fmt.Fprintf(out, "\n⚠️ Warning: Some external dependencies are missing. Some generator features might not work.\n")
	}

	if corruptCount > 0 {
		_, _ = fmt.Fprintf(out, "\n❌ Doctor check complete. Project structure is CORRUPT. Please fix the missing files.\n")
	} else {
		_, _ = fmt.Fprintf(out, "\n✨ Doctor check complete. Project structure is sound.\n")
	}
	return nil
}
