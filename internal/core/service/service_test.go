package service_test

import (
	"bytes"
	"context"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/muhananaufal/go-aether/internal/adapter/manifest"
	"github.com/muhananaufal/go-aether/internal/adapter/template"
	"github.com/muhananaufal/go-aether/internal/adapter/writer"
	"github.com/muhananaufal/go-aether/internal/core/service"
	"github.com/spf13/afero"
)

func TestAetherScaffoldService_E2E_InitAndMakeModule(t *testing.T) {
	memFS := afero.NewMemMapFs()
	w := writer.NewAferoWriter(memFS)
	resolver := manifest.NewYamlResolver(w)

	// Mock embedded templates matching real layout
	mockEmbed := fstest.MapFS{
		"common/aether_yaml.tmpl":                      &fstest.MapFile{Data: []byte("version: \"1\"")},
		"hexagonal/domain.go.tmpl":                     &fstest.MapFile{Data: []byte("package domain\n// domain for {{ .ModuleNameTitle }}")},
		"hexagonal/port.go.tmpl":                       &fstest.MapFile{Data: []byte("package port\n// port for {{ .ModuleNameTitle }}")},
		"hexagonal/service.go.tmpl":                    &fstest.MapFile{Data: []byte("package service\n// service for {{ .ModuleNameTitle }}")},
		"hexagonal/handler_http.go.tmpl":               &fstest.MapFile{Data: []byte("package http\n// http handler for {{ .ModuleNameTitle }}")},
		"hexagonal/repository_postgres.go.tmpl":        &fstest.MapFile{Data: []byte("package repository\n// repo for {{ .ModuleNameTitle }}")},
	}

	engine := template.NewStdEngine(mockEmbed)
	svc := service.NewAetherScaffoldService(engine, resolver, w)
	ctx := context.Background()
	projectDir := "/projects/e2e-app"
	_ = w.MkdirAll(projectDir)

	// 1. Initialize project
	if err := svc.InitProject(ctx, projectDir, "e2e-app", "github.com/e2e/app", "hexagonal", "postgres", "chi", false); err != nil {
		t.Fatalf("InitProject failed: %v", err)
	}

	manifestFile := filepath.Join(projectDir, "aether.yaml")
	exists, _ := w.Exists(manifestFile)
	if !exists {
		t.Fatalf("expected manifest file %s to be created during init", manifestFile)
	}

	// 2. Make module 'order'
	if err := svc.MakeModule(ctx, projectDir, "order", []string{"http"}, false, false, false, false); err != nil {
		t.Fatalf("MakeModule failed: %v", err)
	}

	// Verify generated layer files exist in memory filesystem
	expectedFiles := []string{
		"/projects/e2e-app/internal/core/domain/order.go",
		"/projects/e2e-app/internal/core/port/order_port.go",
		"/projects/e2e-app/internal/core/service/order_service.go",
		"/projects/e2e-app/internal/adapter/handler/http/order_handler.go",
		"/projects/e2e-app/internal/adapter/repository/order_repository.go",
	}

	for _, ef := range expectedFiles {
		found, err := w.Exists(ef)
		if err != nil || !found {
			t.Errorf("expected generated file %q to exist in filesystem, but it was absent", ef)
		}
	}

	// 3. Run Doctor diagnostic checks
	var docBuf bytes.Buffer
	if err := svc.RunDoctor(ctx, projectDir, false, &docBuf); err != nil {
		t.Fatalf("RunDoctor failed: %v", err)
	}
	output := docBuf.String()
	if !bytes.Contains([]byte(output), []byte("Doctor check complete")) {
		t.Errorf("expected doctor check success message, got:\n%s", output)
	}
}
