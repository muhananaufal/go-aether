package domain_test

import (
	"testing"

	"github.com/muhananaufal/go-aether/internal/core/domain"
)

func TestValidateGoIdentifier_ValidNames(t *testing.T) {
	validNames := []string{
		"order",
		"payment_gateway",
		"invoice",
		"user2026",
		"cart_item",
	}
	for _, name := range validNames {
		t.Run(name, func(t *testing.T) {
			if err := domain.ValidateGoIdentifier(name); err != nil {
				t.Errorf("expected valid identifier for %q, got err: %v", name, err)
			}
		})
	}
}

func TestValidateGoIdentifier_InvalidNames(t *testing.T) {
	invalidCases := []struct {
		name        string
		expectedErr error
	}{
		{"my-order", domain.ErrInvalidIdentifier},     // hyphen not allowed in Go package names
		{"212fast", domain.ErrInvalidIdentifier},      // cannot start with number
		{"order space", domain.ErrInvalidIdentifier},  // spaces not allowed
		{"context", domain.ErrStdlibCollision},        // standard library package name collision
		{"sync", domain.ErrStdlibCollision},           // standard library package name collision
		{"utils", domain.ErrStdlibCollision},          // forbidden anti-pattern package name
		{"", domain.ErrInvalidIdentifier},             // empty string
	}
	for _, tc := range invalidCases {
		t.Run(tc.name, func(t *testing.T) {
			err := domain.ValidateGoIdentifier(tc.name)
			if err == nil {
				t.Errorf("expected error for %q, but got nil", tc.name)
			}
		})
	}
}

func TestToTitleCase(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"order", "Order"},
		{"order_item", "OrderItem"},
		{"payment_gateway_client", "PaymentGatewayClient"},
	}
	for _, c := range cases {
		if result := domain.ToTitleCase(c.input); result != c.expected {
			t.Errorf("expected %q, got %q", c.expected, result)
		}
	}
}

func TestAetherManifest_ValidationAndRegistry(t *testing.T) {
	m := domain.NewDefaultManifest("payment", "github.com/company/payment", "1.23", "hexagonal", "postgres", "chi")
	if err := m.Validate(); err != nil {
		t.Fatalf("expected default manifest to be valid, got err: %v", err)
	}

	if m.HasModule("order") {
		t.Errorf("expected module 'order' to not exist initially")
	}

	reg := domain.ModuleRegistry{
		Name:       "order",
		Transports: []string{"http"},
		HasCache:   true,
	}
	if err := m.AddModule(reg); err != nil {
		t.Fatalf("expected AddModule to succeed, got err: %v", err)
	}

	if !m.HasModule("order") {
		t.Errorf("expected module 'order' to exist after registration")
	}

	// Adding duplicate should fail
	if err := m.AddModule(reg); err == nil {
		t.Errorf("expected duplicate module addition to fail, but it succeeded")
	}
}

func TestNewTemplateData(t *testing.T) {
	m := domain.NewDefaultManifest("service", "github.com/company/service", "1.23", "hexagonal", "postgres", "chi")
	m.Adapters.ExistingDBVar = "app.DBPool" // Simulate brownfield BYO injection

	td, err := domain.NewTemplateData("order_item", m, []string{"http", "grpc"}, true, false)
	if err != nil {
		t.Fatalf("expected NewTemplateData to succeed, got err: %v", err)
	}
	if td.ModuleNameTitle != "OrderItem" {
		t.Errorf("expected title 'OrderItem', got %q", td.ModuleNameTitle)
	}
	if td.ExistingDBVar != "app.DBPool" {
		t.Errorf("expected BYO DB injection var to be preserved, got %q", td.ExistingDBVar)
	}
}
