package e2e_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// supportedRouters and supportedDrivers are the stack combinations go-aether
// advertises. Every entry here is a promise made to the user by the --router and
// --db flags, and this file is what turns each promise into something measured.
var (
	supportedRouters = []string{"chi", "gin", "echo", "fiber", "stdlib"}
	supportedDrivers = []string{"postgres", "mysql", "sqlite"}
)

// routerImportPath maps a router selection to the module its generated code must
// actually import. Asserting on this is what separates real support from the
// previous behaviour, where every selection silently produced a chi project.
var routerImportPath = map[string]string{
	"chi":    "github.com/go-chi/chi/v5",
	"gin":    "github.com/gin-gonic/gin",
	"echo":   "github.com/labstack/echo/v4",
	"fiber":  "github.com/gofiber/fiber/v2",
	"stdlib": "", // net/http only; no third-party module is expected.
}

// driverImportPath maps a database selection to its driver module.
var driverImportPath = map[string]string{
	"postgres": "github.com/jackc/pgx/v5",
	"mysql":    "github.com/go-sql-driver/mysql",
	"sqlite":   "modernc.org/sqlite",
}

// TestFrameworkMatrix_EveryCombinationCompiles is the honesty gate for the
// --router and --db flags.
//
// Before this existed, `init --router gin` produced a project importing chi.
// The flag was accepted, the output was wrong, and nothing failed: the worst
// possible outcome for a newcomer, who has no way to know the tool ignored them.
func TestFrameworkMatrix_EveryCombinationCompiles(t *testing.T) {
	requireGoToolchain(t)
	if testing.Short() {
		t.Skip("matrix invokes the Go toolchain per combination; skipped under -short")
	}

	for _, router := range supportedRouters {
		for _, driver := range supportedDrivers {
			t.Run(router+"_"+driver, func(t *testing.T) {
				dir := scaffoldProject(t, router, driver)

				if out, err := runInDir(t, dir, "go", "build", "./..."); err != nil {
					t.Fatalf("router=%s db=%s produces a project that does not compile: %v\n--- output ---\n%s",
						router, driver, err, out)
				}

				if out, err := runInDir(t, dir, "gofmt", "-l", "."); err != nil {
					t.Fatalf("gofmt failed to inspect router=%s db=%s: %v\n%s", router, driver, err, out)
				} else if offenders := nonEmptyLines(out); len(offenders) > 0 {
					t.Errorf("router=%s db=%s emitted unformatted source:\n  %s",
						router, driver, strings.Join(offenders, "\n  "))
				}
			})
		}
	}
}

// TestFrameworkMatrix_SelectionReachesGeneratedCode proves the selection is not
// merely accepted but actually honoured, by reading the dependency graph the
// generated project declares.
//
// A project that compiles is necessary but not sufficient: a chi project
// compiles perfectly well when the user asked for echo.
func TestFrameworkMatrix_SelectionReachesGeneratedCode(t *testing.T) {
	requireGoToolchain(t)
	if testing.Short() {
		t.Skip("matrix invokes the Go toolchain per combination; skipped under -short")
	}

	for _, router := range supportedRouters {
		t.Run("router_"+router, func(t *testing.T) {
			dir := scaffoldProject(t, router, "postgres")
			goMod := readFileOrFail(t, filepath.Join(dir, "go.mod"))

			if want := routerImportPath[router]; want != "" && !strings.Contains(goMod, want) {
				t.Errorf("--router %s did not reach the generated code: go.mod lacks %s\n--- go.mod ---\n%s",
					router, want, goMod)
			}

			// Every non-chi selection must be free of chi entirely. This is the
			// assertion that would have caught the original defect.
			if router != "chi" && strings.Contains(goMod, "go-chi/chi") {
				t.Errorf("--router %s still pulled in chi; the selection was ignored\n--- go.mod ---\n%s",
					router, goMod)
			}
		})
	}

	for _, driver := range supportedDrivers {
		t.Run("db_"+driver, func(t *testing.T) {
			dir := scaffoldProject(t, "chi", driver)
			goMod := readFileOrFail(t, filepath.Join(dir, "go.mod"))

			if want := driverImportPath[driver]; !strings.Contains(goMod, want) {
				t.Errorf("--db %s did not reach the generated code: go.mod lacks %s\n--- go.mod ---\n%s",
					driver, want, goMod)
			}
		})
	}
}

// TestFrameworkMatrix_SqliteBuildsWithoutCgo guards a trap specific to SQLite.
//
// The popular mattn/go-sqlite3 driver requires cgo. The Dockerfile this tool
// generates builds with CGO_ENABLED=0 to produce a static binary, so choosing a
// cgo driver would yield a project that builds on the developer's laptop and
// fails inside the container with a link error nobody expects.
func TestFrameworkMatrix_SqliteBuildsWithoutCgo(t *testing.T) {
	requireGoToolchain(t)
	if testing.Short() {
		t.Skip("compile gate invokes the Go toolchain; skipped under -short")
	}

	dir := scaffoldProject(t, "chi", "sqlite")

	cmd := []string{"go", "build", "./..."}
	out, err := runInDirWithEnv(t, dir, []string{"CGO_ENABLED=0"}, cmd[0], cmd[1:]...)
	if err != nil {
		t.Fatalf("sqlite project fails to build with CGO_ENABLED=0, which is exactly how "+
			"the generated Dockerfile builds it: %v\n--- output ---\n%s", err, out)
	}
}

// TestFrameworkMatrix_UnsupportedSelectionIsRejected asserts the tool refuses
// what it cannot deliver instead of quietly substituting a default.
//
// Silently downgrading is how the original defect stayed invisible for so long.
func TestFrameworkMatrix_UnsupportedSelectionIsRejected(t *testing.T) {
	svc := newRealFSService()
	ctx := context.Background()

	cases := []struct {
		name           string
		router, driver string
		arch           string
	}{
		{"unknown router", "sinatra", "postgres", "hexagonal"},
		{"unknown driver", "chi", "cassandra", "hexagonal"},
		{"unimplemented architecture", "chi", "postgres", "clean"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			err := svc.InitProject(ctx, dir, "probeapp", "example.com/probeapp", tc.arch, tc.driver, tc.router, false)
			if err == nil {
				t.Fatalf("InitProject accepted an unsupported stack (arch=%s router=%s db=%s)",
					tc.arch, tc.router, tc.driver)
			}

			// The message has to name the valid options, or the user is left guessing.
			if !strings.Contains(err.Error(), "supported") {
				t.Errorf("rejection message must list what is supported, got: %v", err)
			}

			// A rejected selection must leave nothing behind.
			if entries, readErr := os.ReadDir(dir); readErr == nil && len(entries) > 0 {
				var names []string
				for _, e := range entries {
					names = append(names, e.Name())
				}
				t.Errorf("rejected init left %d artefact(s) on disk: %s", len(names), strings.Join(names, ", "))
			}
		})
	}
}

// TestFrameworkMatrix_SecondModuleDoesNotCollide catches a whole class of defect
// that every single-module test is blind to.
//
// All generated handlers land in one package. Any package-level helper in a
// handler template — a writeJSON, a maxBodyBytes const — is therefore redeclared
// the moment a second module is scaffolded, and the project stops compiling on
// the user's second command rather than their first.
//
// The check runs across every router because each handler template is written
// independently and only one of them has to slip.
func TestFrameworkMatrix_SecondModuleDoesNotCollide(t *testing.T) {
	requireGoToolchain(t)
	if testing.Short() {
		t.Skip("compile gate invokes the Go toolchain; skipped under -short")
	}

	svc := newRealFSService()
	ctx := context.Background()

	for _, router := range supportedRouters {
		t.Run(router, func(t *testing.T) {
			dir := t.TempDir()
			if err := svc.InitProject(ctx, dir, "probeapp", "example.com/probeapp", "hexagonal", "postgres", router, false); err != nil {
				t.Fatalf("InitProject failed: %v", err)
			}

			for _, module := range []string{"order", "product", "invoice"} {
				if err := svc.MakeModule(ctx, dir, module, []string{"http"}, false, false, false, false); err != nil {
					t.Fatalf("MakeModule(%s) failed: %v", module, err)
				}
			}

			if out, err := runInDir(t, dir, "go", "build", "./..."); err != nil {
				t.Fatalf("router=%s: three modules in one project do not compile together.\n"+
					"A package-level helper in the handler template is redeclared per module.\n"+
					"%v\n--- output ---\n%s", router, err, out)
			}
		})
	}
}

func readFileOrFail(t *testing.T, path string) string {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read %s: %v", path, err)
	}
	return string(raw)
}
