package adapter_test

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/muhananaufal/go-aether/internal/adapter/manifest"
	"github.com/muhananaufal/go-aether/internal/adapter/template"
	"github.com/muhananaufal/go-aether/internal/adapter/writer"
	"github.com/muhananaufal/go-aether/internal/core/domain"
	"github.com/spf13/afero"
)

func TestAferoWriter_WriteAndConflict(t *testing.T) {
	memFS := afero.NewMemMapFs()
	w := writer.NewAferoWriter(memFS)
	ctx := context.Background()

	path := "internal/core/service/order_service.go"
	content := []byte("package service\n// order service")

	// First write should succeed
	if err := w.WriteFile(ctx, path, content, false, false); err != nil {
		t.Fatalf("expected write to succeed, got %v", err)
	}

	exists, _ := w.Exists(path)
	if !exists {
		t.Errorf("expected file to exist after write")
	}

	// Second write without overwrite should fail with ErrFileConflict
	err := w.WriteFile(ctx, path, []byte("overwrite attempt"), false, false)
	if !errors.Is(err, domain.ErrFileConflict) {
		t.Errorf("expected ErrFileConflict on second write, got %v", err)
	}

	// Third write WITH overwrite should succeed and create .bak
	if err := w.WriteFile(ctx, path, []byte("new content"), true, false); err != nil {
		t.Fatalf("expected overwrite to succeed, got %v", err)
	}

	bakExists, _ := w.Exists(path + ".bak")
	if !bakExists {
		t.Errorf("expected defensive backup .bak file to be generated")
	}
}

func TestTransactionalBuffer_RollbackOnFailure(t *testing.T) {
	memFS := afero.NewReadOnlyFs(afero.NewMemMapFs()) // read-only FS causes writes to fail
	baseWriter := writer.NewAferoWriter(memFS)
	tx := writer.NewTransactionalBuffer(baseWriter)

	tx.Stage("file1.go", []byte("package test"), true, false)
	tx.Stage("file2.go", []byte("package test"), true, false)

	err := tx.Commit(context.Background())
	if !errors.Is(err, domain.ErrWriteFailed) {
		t.Errorf("expected ErrWriteFailed when writing to read-only disk, got %v", err)
	}
}

func TestYamlResolver_WalkUpSearch(t *testing.T) {
	memFS := afero.NewMemMapFs()
	w := writer.NewAferoWriter(memFS)
	resolver := manifest.NewYamlResolver(w)
	ctx := context.Background()

	rootManifest := domain.NewDefaultManifest("payment-service", "github.com/company/payment", "1.23", "hexagonal", "postgres", "chi")
	if err := resolver.Save(ctx, "/projects/payment", rootManifest, w); err != nil {
		t.Fatalf("failed to save manifest: %v", err)
	}

	// Simulate walk-up from 3 directories deep in the structure
	deepDir := "/projects/payment/internal/adapter/handler/http"
	_ = memFS.MkdirAll(deepDir, 0755)

	found, foundPath, err := resolver.Resolve(ctx, deepDir)
	if err != nil {
		t.Fatalf("expected resolver to find aether.yaml via walk-up search, got err: %v", err)
	}

	expectedPath := "/projects/payment/aether.yaml"
	if filepath.ToSlash(foundPath) != expectedPath {
		t.Errorf("expected path %q, got %q", expectedPath, foundPath)
	}
	if found.Project.Name != "payment-service" {
		t.Errorf("expected project name 'payment-service', got %q", found.Project.Name)
	}
}

func TestStdEngine_RenderTemplate(t *testing.T) {
	mockEmbed := fstest.MapFS{
		"hexagonal/service.go.tmpl": &fstest.MapFile{
			Data: []byte("package service\n\n// {{ .ModuleNameTitle }}Service handles core logic for {{ .ModuleName }}.\ntype {{ .ModuleNameTitle }}Service struct {}"),
		},
	}
	eng := template.NewStdEngine(mockEmbed)
	ctx := context.Background()

	m := domain.NewDefaultManifest("app", "github.com/company/app", "1.23", "hexagonal", "postgres", "chi")
	data, err := domain.NewTemplateData("order", m, nil, false, false)
	if err != nil {
		t.Fatalf("failed to construct TemplateData: %v", err)
	}

	rendered, err := eng.Render(ctx, "hexagonal/service.go.tmpl", data)
	if err != nil {
		t.Fatalf("expected render to succeed, got %v", err)
	}

	expected := "package service\n\n// OrderService handles core logic for order.\ntype OrderService struct {}"
	if string(rendered) != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, string(rendered))
	}
}
