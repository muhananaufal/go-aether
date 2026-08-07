package e2e_test

import (
	"context"
	"testing"
)

// TestCompileGate_GranularLayerCommandsCompile extends the compile gate to the
// per-layer generators.
//
// arch:module is covered by the matrix, but arch:domain, arch:port and
// arch:repository render entirely different templates (domain_only, port_only,
// repository_only) that no test has ever compiled. They are the commands a
// newcomer reaches for when adding one piece to an existing slice, so shipping
// them unverified is exactly backwards.
func TestCompileGate_GranularLayerCommandsCompile(t *testing.T) {
	requireGoToolchain(t)
	if testing.Short() {
		t.Skip("compile gate invokes the Go toolchain; skipped under -short")
	}

	dir := t.TempDir()
	svc := newRealFSService()
	ctx := context.Background()

	if err := svc.InitProject(ctx, dir, "probeapp", "example.com/probeapp", "hexagonal", "postgres", "chi", false); err != nil {
		t.Fatalf("InitProject failed: %v", err)
	}

	// The order a developer would actually use: entity first, then its contract,
	// then the adapter that satisfies it.
	steps := []struct {
		name string
		run  func() error
	}{
		{"arch:domain", func() error { return svc.MakeDomain(ctx, dir, "invoice", false, false) }},
		{"arch:port", func() error { return svc.MakePort(ctx, dir, "invoice", false, false) }},
		{"arch:repository", func() error { return svc.MakeRepository(ctx, dir, "invoice", false, false) }},
	}
	for _, step := range steps {
		if err := step.run(); err != nil {
			t.Fatalf("%s failed: %v", step.name, err)
		}
	}

	out, err := runInDir(t, dir, "go", "build", "./...")
	if err != nil {
		t.Fatalf("the granular layer commands produce a project that does not compile.\n"+
			"These are the commands used to extend an existing slice one piece at a time.\n"+
			"%v\n--- output ---\n%s", err, out)
	}
}

// TestCompileGate_GranularSliceBuiltPieceByPiece walks the path a developer
// takes when they build an entity up one command at a time rather than asking
// for the whole slice.
//
// This combination crosses two template families: arch:domain renders
// domain_only, which carries no behaviour, while arch:service renders the same
// service template arch:module uses, which calls methods on the entity. Nothing
// guarantees the two agree, and no test has ever put them in the same package.
func TestCompileGate_GranularSliceBuiltPieceByPiece(t *testing.T) {
	requireGoToolchain(t)
	if testing.Short() {
		t.Skip("compile gate invokes the Go toolchain; skipped under -short")
	}

	dir := t.TempDir()
	svc := newRealFSService()
	ctx := context.Background()

	if err := svc.InitProject(ctx, dir, "probeapp", "example.com/probeapp", "hexagonal", "postgres", "chi", false); err != nil {
		t.Fatalf("InitProject failed: %v", err)
	}

	steps := []struct {
		name string
		run  func() error
	}{
		{"arch:domain", func() error { return svc.MakeDomain(ctx, dir, "invoice", false, false) }},
		{"arch:port", func() error { return svc.MakePort(ctx, dir, "invoice", false, false) }},
		{"arch:service", func() error { return svc.MakeService(ctx, dir, "invoice", false, false) }},
		{"arch:repository", func() error { return svc.MakeRepository(ctx, dir, "invoice", false, false) }},
		{"arch:handler", func() error { return svc.MakeHandler(ctx, dir, "invoice", "http", false, false) }},
	}
	for _, step := range steps {
		if err := step.run(); err != nil {
			t.Fatalf("%s failed: %v", step.name, err)
		}
	}

	out, err := runInDir(t, dir, "go", "build", "./...")
	if err != nil {
		t.Fatalf("a slice assembled command by command does not compile.\n"+
			"arch:domain and arch:service render from different template families "+
			"that were never verified against each other.\n%v\n--- output ---\n%s", err, out)
	}
}

// TestCompileGate_GranularServiceAndHandlerCompile covers the other two per-layer
// commands, which reuse the templates arch:module already compiles but have
// never been exercised through their own code path.
func TestCompileGate_GranularServiceAndHandlerCompile(t *testing.T) {
	requireGoToolchain(t)
	if testing.Short() {
		t.Skip("compile gate invokes the Go toolchain; skipped under -short")
	}

	dir := t.TempDir()
	svc := newRealFSService()
	ctx := context.Background()

	if err := svc.InitProject(ctx, dir, "probeapp", "example.com/probeapp", "hexagonal", "postgres", "chi", false); err != nil {
		t.Fatalf("InitProject failed: %v", err)
	}

	// A full slice first, so the service and handler have a domain and port to
	// compile against, then the granular commands for a second entity.
	if err := svc.MakeModule(ctx, dir, "order", []string{"http"}, false, false, false, false); err != nil {
		t.Fatalf("MakeModule failed: %v", err)
	}

	if err := svc.MakeService(ctx, dir, "order", false, true); err != nil {
		t.Fatalf("arch:service --force failed: %v", err)
	}
	if err := svc.MakeHandler(ctx, dir, "order", "http", false, true); err != nil {
		t.Fatalf("arch:handler --force failed: %v", err)
	}

	out, err := runInDir(t, dir, "go", "build", "./...")
	if err != nil {
		t.Fatalf("regenerating the service and handler broke the build: %v\n--- output ---\n%s", err, out)
	}
}
