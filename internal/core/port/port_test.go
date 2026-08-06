package port_test

import (
	"context"
	"testing"

	"github.com/muhananaufal/go-aether/internal/core/domain"
	"github.com/muhananaufal/go-aether/internal/core/port"
)

// mockTemplateEngine confirms that any external adapter implementing TemplateEngine satisfies the port contract.
type mockTemplateEngine struct{}

func (m *mockTemplateEngine) Render(ctx context.Context, templatePath string, data *domain.TemplateData) ([]byte, error) {
	if templatePath == "hexagonal/service.go.tmpl" {
		return []byte("package service\n// generated service for " + data.ModuleNameTitle), nil
	}
	return nil, domain.ErrTemplateMissing
}

func (m *mockTemplateEngine) ListTemplates(pattern string) ([]string, error) {
	return []string{"domain.go.tmpl", "service.go.tmpl"}, nil
}

func TestTemplateEngineInterface_ContractCompliance(t *testing.T) {
	var engine port.TemplateEngine = &mockTemplateEngine{}

	m := domain.NewDefaultManifest("app", "github.com/test/app", "1.23", "hexagonal", "postgres", "chi")
	data, err := domain.NewTemplateData("order", m, []string{"http"}, false, false)
	if err != nil {
		t.Fatalf("expected template data creation to succeed, got %v", err)
	}

	content, err := engine.Render(context.Background(), "hexagonal/service.go.tmpl", data)
	if err != nil {
		t.Fatalf("expected render to succeed, got %v", err)
	}

	expectedSnippet := "package service\n// generated service for Order"
	if string(content) != expectedSnippet {
		t.Errorf("expected content %q, got %q", expectedSnippet, string(content))
	}
}
