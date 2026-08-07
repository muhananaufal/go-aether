package service

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/muhananaufal/go-aether/internal/core/domain"
)

// adoptScanTimeout bounds the filesystem walk. Adoption is normally pointed at
// somebody else's repository, whose size is unknown until it is walked.
const adoptScanTimeout = 30 * time.Second

// AdoptProject inspects an existing repository and proposes an aether.yaml that
// matches the structure that is actually there.
//
// The previous implementation asked five questions and then wrote the default
// hexagonal paths regardless of the answers, so adopting a legacy repository
// produced a manifest describing a layout that project does not have. Every
// generator then wrote into directories that did not exist.
//
// Two rules govern this command, both because it operates on code somebody
// already depends on:
//
//   - Nothing is written unless the caller explicitly opts in. dryRun is the
//     default at the CLI layer.
//   - Every proposed path is printed with the evidence behind it, so the mapping
//     can be judged rather than trusted.
func (s *AetherScaffoldService) AdoptProject(ctx context.Context, destDir string, scan, dryRun bool) error {
	manifestPath := filepath.Join(destDir, "aether.yaml")
	exists, err := s.fs.Exists(manifestPath)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("%w: repository already contains aether.yaml", domain.ErrFileConflict)
	}

	out := s.adoptOutput()

	moduleName := detectModulePath(destDir)
	if moduleName == "" {
		// Without go.mod there is no import path, so every generated import would
		// be unresolvable. Better to stop than to invent a module name.
		return fmt.Errorf("%w: cannot adopt %s", domain.ErrGoModMissing, destDir)
	}

	manifest := domain.NewDefaultManifest(filepath.Base(destDir), moduleName, detectGoVersion(),
		domain.DefaultArchitecture, domain.DefaultDBDriver, domain.DefaultRouter)
	manifest.Architecture.Mode = "brownfield"
	manifest.Project.CreatedAt = time.Now().Format(time.RFC3339)

	var report *domain.LayoutReport
	if scan {
		report, err = s.scanLayout(ctx, destDir, manifest, out)
		if err != nil {
			return err
		}
		s.applyBYODetection(ctx, destDir, manifest, out)
	}

	s.printAdoptionPlan(out, destDir, manifest, report, dryRun)

	if dryRun {
		return nil
	}

	if err := s.resolver.Save(ctx, destDir, manifest, s.fs); err != nil {
		return fmt.Errorf("failed saving adopted manifest: %w", err)
	}
	_, _ = fmt.Fprintf(out, "\n✅ Wrote %s\n", filepath.ToSlash(manifestPath))
	return nil
}

// scanLayout runs the bounded scan and folds confident findings into the
// manifest, leaving the defaults in place where the scanner is unsure.
func (s *AetherScaffoldService) scanLayout(
	ctx context.Context,
	destDir string,
	manifest *domain.AetherManifest,
	out io.Writer,
) (*domain.LayoutReport, error) {
	if s.scanner == nil {
		return nil, fmt.Errorf("aether: adoption requires a layout scanner; none was wired")
	}

	scanCtx, cancel := context.WithTimeout(ctx, adoptScanTimeout)
	defer cancel()

	report, err := s.scanner.Scan(scanCtx, destDir)
	if err != nil {
		return nil, fmt.Errorf("layout scan failed: %w", err)
	}

	assign := map[domain.LayoutKind]*string{
		domain.LayerHandler:    &manifest.Architecture.Paths.HandlerHTTP,
		domain.LayerService:    &manifest.Architecture.Paths.Service,
		domain.LayerRepository: &manifest.Architecture.Paths.Repository,
		domain.LayerDomain:     &manifest.Architecture.Paths.Domain,
	}
	for kind, target := range assign {
		if best, ok := report.Best(kind); ok {
			*target = best.Dir
		}
	}

	// Anomaly mode is not a label the user picks; it is a measurement. A layout
	// the scanner could not recognise is precisely the case where later commands
	// must not assume the standard hexagonal tree.
	confidence := report.Confidence()
	manifest.Meta.AnomalyMode = confidence < 1.0
	manifest.Meta.LegacyNotes = fmt.Sprintf(
		"Adopted by go-aether scan on %s: %d Go files inspected, layout confidence %.0f%%.",
		time.Now().Format(time.RFC3339), report.FilesSeen, confidence*100)

	if report.Truncated {
		_, _ = fmt.Fprintf(out, "\n⚠️  Scan stopped at the %d file limit; deeper directories were not inspected.\n",
			report.FilesSeen)
	}

	return report, nil
}

// applyBYODetection records infrastructure clients the project already builds so
// generated constructors can accept them instead of opening duplicates.
func (s *AetherScaffoldService) applyBYODetection(
	ctx context.Context,
	destDir string,
	manifest *domain.AetherManifest,
	out io.Writer,
) {
	if s.byoDetector == nil {
		return
	}

	detection, err := s.byoDetector.Detect(ctx, destDir)
	if err != nil || detection.IsEmpty() {
		// Finding nothing is an ordinary outcome, not a failure: plenty of
		// projects construct their clients somewhere this detector cannot reach.
		return
	}

	manifest.Adapters.ExistingDBVar = detection.DBVar
	manifest.Adapters.ExistingRedisVar = detection.RedisVar
	manifest.Adapters.ExistingLoggerVar = detection.LoggerVar

	_, _ = fmt.Fprintf(out, "\n🔌 Reusable clients found in your entrypoint:\n")
	for _, pair := range []struct{ label, name string }{
		{"database", detection.DBVar},
		{"redis", detection.RedisVar},
		{"logger", detection.LoggerVar},
	} {
		if pair.name == "" {
			continue
		}
		_, _ = fmt.Fprintf(out, "   • %-9s %-14s %s\n", pair.label, pair.name, detection.Sources[pair.name])
	}
	_, _ = fmt.Fprintf(out, "   Generated constructors will ask for these instead of opening their own.\n")
}

// printAdoptionPlan renders the proposal. It always runs, including on a real
// write, so the record of what was decided appears in the terminal either way.
func (s *AetherScaffoldService) printAdoptionPlan(
	out io.Writer,
	destDir string,
	manifest *domain.AetherManifest,
	report *domain.LayoutReport,
	dryRun bool,
) {
	_, _ = fmt.Fprintf(out, "\n🔍 Adoption plan for %s\n", filepath.ToSlash(destDir))
	_, _ = fmt.Fprintf(out, "   module: %s\n", manifest.Project.Module)

	if report != nil {
		_, _ = fmt.Fprintf(out, "   scanned: %d Go files, layout confidence %.0f%%\n",
			report.FilesSeen, report.Confidence()*100)
	}

	_, _ = fmt.Fprintf(out, "\n   Proposed layer mapping:\n")
	rows := []struct {
		label string
		kind  domain.LayoutKind
		path  string
	}{
		{"handler", domain.LayerHandler, manifest.Architecture.Paths.HandlerHTTP},
		{"service", domain.LayerService, manifest.Architecture.Paths.Service},
		{"repository", domain.LayerRepository, manifest.Architecture.Paths.Repository},
		{"domain", domain.LayerDomain, manifest.Architecture.Paths.Domain},
	}
	for _, row := range rows {
		evidence := "default (nothing recognised — verify this)"
		if report != nil {
			if best, ok := report.Best(row.kind); ok {
				evidence = strings.Join(best.Evidence, "; ")
			}
		}
		_, _ = fmt.Fprintf(out, "   %-11s %-34s %s\n", row.label, row.path, evidence)
	}

	if manifest.Meta.AnomalyMode {
		_, _ = fmt.Fprintf(out, "\n   ⚠️  anomaly_mode is on: at least one layer could not be identified.\n")
		_, _ = fmt.Fprintf(out, "      Edit the paths above in aether.yaml before generating anything.\n")
	}

	if dryRun {
		_, _ = fmt.Fprintf(out, "\n   Nothing was written. Re-run with --apply to save this manifest.\n")
	}
}

// adoptOutput centralises where the plan is printed so the destination can be
// redirected in a future revision without touching every call site.
func (s *AetherScaffoldService) adoptOutput() io.Writer { return os.Stdout }

// detectModulePath reads the module line from go.mod.
//
// Adoption previously invented "github.com/adopted/<dirname>", which is never
// the project's real import path, so every generated import was wrong from the
// first file.
func detectModulePath(dir string) string {
	file, err := os.Open(filepath.Join(dir, "go.mod"))
	if err != nil {
		return ""
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	first := true
	for scanner.Scan() {
		line := scanner.Text()
		if first {
			// Editors on Windows routinely save UTF-8 with a byte order mark. A BOM
			// ahead of "module" stops the prefix from matching, and the project is
			// then reported as not being a Go module at all: a confusing failure for
			// a file the user can plainly see is correct.
			line = strings.TrimPrefix(line, "\ufeff")
			first = false
		}

		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "module") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			return strings.Trim(fields[1], `"`)
		}
	}
	return ""
}
