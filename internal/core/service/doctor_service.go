package service

import (
	"context"
	"fmt"
	"io"
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
	for _, mod := range manifest.Modules {
		_, _ = fmt.Fprintf(out, "   - Module [%s]: Transports=%v, Cache=%t\n", mod.Name, mod.Transports, mod.HasCache)
	}

	if manifest.Meta.AnomalyMode {
		_, _ = fmt.Fprintf(out, "⚠️ [NOTE] Anomaly mode active for legacy brownfield structure.\n")
	}

	_, _ = fmt.Fprintf(out, "✨ Doctor check complete. Project structure is sound.\n")
	return nil
}
