// Package e2e holds the compile gate: the only test layer that proves the code
// go-aether emits is code the Go toolchain actually accepts.
//
// Every other test in this repository verifies that the generator wrote *a* file.
// None of them verify that the file compiles. That gap is why a released version
// could ship an `init` command whose output fails `go build` on the first try.
//
// This suite closes it by generating into a real temporary directory and then
// running the real toolchain over the result.
package e2e_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/muhananaufal/go-aether/internal/adapter/manifest"
	"github.com/muhananaufal/go-aether/internal/adapter/scanner"
	"github.com/muhananaufal/go-aether/internal/adapter/template"
	"github.com/muhananaufal/go-aether/internal/adapter/writer"
	"github.com/muhananaufal/go-aether/internal/core/port"
	"github.com/muhananaufal/go-aether/internal/core/service"
	"github.com/muhananaufal/go-aether/templates"
	"github.com/spf13/afero"
)

// newRealFSService wires the production object graph against the real OS
// filesystem. Using afero.MemMapFs here would defeat the purpose: the defects
// this suite hunts (go.mod population, gofmt, OS-reserved filenames) only exist
// on a real disk with a real toolchain.
func newRealFSService() port.ScaffoldService {
	osFS := afero.NewOsFs()
	fileWriter := writer.NewAferoWriter(osFS)
	resolver := manifest.NewYamlResolver(fileWriter)
	engine := template.NewStdEngine(templates.FS)
	return service.NewAetherScaffoldService(
		engine, resolver, fileWriter,
		scanner.NewGoLayoutScanner(),
		scanner.NewGoBYODetector(),
	)
}

// requireGoToolchain skips rather than fails when the toolchain is unavailable,
// because an absent `go` binary is an environment problem, not a defect in the
// code under test. Everything else is a hard failure on purpose.
func requireGoToolchain(t *testing.T) string {
	t.Helper()
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skipf("go toolchain not present in PATH: %v", err)
	}
	return goBin
}

// runInDir executes a toolchain command inside dir and returns its combined
// output. The output is returned even on failure because the assertion messages
// are worthless without it.
func runInDir(t *testing.T, dir, name string, args ...string) (string, error) {
	t.Helper()
	return runInDirWithEnv(t, dir, nil, name, args...)
}

// runInDirWithEnv is runInDir with additional environment entries appended, used
// to reproduce the exact conditions of the generated Dockerfile build.
func runInDirWithEnv(t *testing.T, dir string, extraEnv []string, name string, args ...string) (string, error) {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	// GOFLAGS is cleared so a developer's ambient -mod=vendor does not silently
	// change what this gate measures.
	cmd.Env = append(os.Environ(), "GOFLAGS=")
	cmd.Env = append(cmd.Env, extraEnv...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// scaffoldProject runs the exact sequence a newcomer runs on their first day:
// initialise a project, then generate one vertical slice.
func scaffoldProject(t *testing.T, router, dbDriver string) string {
	t.Helper()
	dir := t.TempDir()
	svc := newRealFSService()
	ctx := context.Background()

	if err := svc.InitProject(ctx, dir, "probeapp", "example.com/probeapp", "hexagonal", dbDriver, router, false); err != nil {
		t.Fatalf("InitProject(router=%s, db=%s) failed: %v", router, dbDriver, err)
	}
	if err := svc.MakeModule(ctx, dir, "order", []string{"http"}, false, false, false, false); err != nil {
		t.Fatalf("MakeModule(order) failed: %v", err)
	}
	return dir
}

// TestCompileGate_InitAndModuleProduceBuildableProject is the headline gate.
//
// It asserts the single promise that matters most to the target user: the code
// go-aether generates compiles without any manual intervention.
func TestCompileGate_InitAndModuleProduceBuildableProject(t *testing.T) {
	requireGoToolchain(t)
	if testing.Short() {
		t.Skip("compile gate invokes the Go toolchain; skipped under -short")
	}

	dir := scaffoldProject(t, "chi", "postgres")

	out, err := runInDir(t, dir, "go", "build", "./...")
	if err != nil {
		t.Fatalf("generated project does not compile.\n"+
			"This is the defect a newcomer hits on their very first command.\n"+
			"go build ./... exited with %v\n--- output ---\n%s", err, out)
	}
	if strings.TrimSpace(out) != "" {
		t.Errorf("go build emitted diagnostics on a freshly generated project:\n%s", out)
	}
}

// TestCompileGate_GeneratedProjectPassesVet catches semantically valid but
// suspicious constructs (unreachable code, bad printf verbs, shadowed errors)
// that compile cleanly yet teach bad habits to whoever reads the output.
func TestCompileGate_GeneratedProjectPassesVet(t *testing.T) {
	requireGoToolchain(t)
	if testing.Short() {
		t.Skip("compile gate invokes the Go toolchain; skipped under -short")
	}

	dir := scaffoldProject(t, "chi", "postgres")

	// vet requires a resolvable dependency graph, so a build failure would surface
	// here as a confusing vet error. Establish the precondition explicitly.
	if out, err := runInDir(t, dir, "go", "build", "./..."); err != nil {
		t.Skipf("skipping vet because the project does not build yet:\n%s", out)
	}

	out, err := runInDir(t, dir, "go", "vet", "./...")
	if err != nil {
		t.Fatalf("go vet rejected the generated project: %v\n--- output ---\n%s", err, out)
	}
}

// TestCompileGate_GeneratedCodeIsGofmtClean asserts that emitted Go is already
// formatted. Unformatted output is not cosmetic here: the generator's entire
// purpose is to demonstrate idiomatic Go, and every reader's editor will
// reformat it on save, producing noise in the very first commit of a project.
func TestCompileGate_GeneratedCodeIsGofmtClean(t *testing.T) {
	requireGoToolchain(t)
	if testing.Short() {
		t.Skip("compile gate invokes the Go toolchain; skipped under -short")
	}

	dir := scaffoldProject(t, "chi", "postgres")

	out, err := runInDir(t, dir, "gofmt", "-l", ".")
	if err != nil {
		t.Fatalf("gofmt could not inspect the generated tree: %v\n%s", err, out)
	}

	offenders := nonEmptyLines(out)
	if len(offenders) > 0 {
		t.Errorf("generated Go source is not gofmt-clean; %d file(s) would be rewritten on save:\n  %s",
			len(offenders), strings.Join(offenders, "\n  "))
	}
}

// TestCompileGate_ModuleGraphIsSelfSufficient pins the specific failure observed
// on v0.3.0: go.mod was written before any template had been rendered, so the
// dependency resolution step had nothing to discover.
//
// Asserting on go.mod content rather than on build success gives a failure
// message that names the actual root cause instead of a wall of compiler errors.
func TestCompileGate_ModuleGraphIsSelfSufficient(t *testing.T) {
	requireGoToolchain(t)
	if testing.Short() {
		t.Skip("compile gate invokes the Go toolchain; skipped under -short")
	}

	dir := scaffoldProject(t, "chi", "postgres")

	raw, err := os.ReadFile(filepath.Join(dir, "go.mod"))
	if err != nil {
		t.Fatalf("generated project has no readable go.mod: %v", err)
	}
	goMod := string(raw)

	// Every module below is imported by a template that init unconditionally writes.
	// If init claims success, these must already be requirements.
	required := []string{
		"github.com/go-chi/chi/v5",
		"github.com/spf13/viper",
	}
	var missing []string
	for _, mod := range required {
		if !strings.Contains(goMod, mod) {
			missing = append(missing, mod)
		}
	}
	if len(missing) > 0 {
		t.Errorf("go.mod omits %d module(s) that generated code imports: %s\n"+
			"The user must run 'go mod tidy' by hand before the project builds.\n--- go.mod ---\n%s",
			len(missing), strings.Join(missing, ", "), goMod)
	}
}

// TestCompileGate_ReservedDeviceNameIsRejected locks the Win32 anomaly proven on
// v0.3.0: `arch:domain con` reported success while the content was routed to the
// CON console device and discarded, leaving no file behind.
//
// The rule is enforced on every OS, not just Windows, so that a project authored
// on Linux stays portable.
func TestCompileGate_ReservedDeviceNameIsRejected(t *testing.T) {
	dir := t.TempDir()
	svc := newRealFSService()
	ctx := context.Background()

	if err := svc.InitProject(ctx, dir, "probeapp", "example.com/probeapp", "hexagonal", "postgres", "chi", false); err != nil {
		t.Fatalf("InitProject failed: %v", err)
	}

	err := svc.MakeDomain(ctx, dir, "con", false, false)
	if err == nil {
		t.Fatalf("MakeDomain(%q) reported success for a reserved device name; "+
			"on Windows the content is written to the CON device and silently lost", "con")
	}

	// A generic failure is not enough: the caller must be able to distinguish this
	// from an I/O error in order to print actionable guidance.
	if !strings.Contains(strings.ToLower(err.Error()), "reserved") {
		t.Errorf("expected a reserved-name error, got a different failure: %v", err)
	}

	if runtime.GOOS == "windows" {
		return // Test-Path cannot meaningfully probe a device name.
	}
	if _, statErr := os.Stat(filepath.Join(dir, "internal", "core", "domain", "con.go")); statErr == nil {
		t.Errorf("con.go was created despite the identifier being rejected")
	}
}

// nonEmptyLines splits toolchain output into meaningful lines, tolerating the
// CRLF that Windows toolchains emit.
func nonEmptyLines(s string) []string {
	var out []string
	for _, line := range strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n") {
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}
