package domain_test

import (
	"errors"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/muhananaufal/go-aether/internal/core/domain"
)

func TestSafeJoin_AcceptsPathsInsideRoot(t *testing.T) {
	root := filepath.Join("C:", "projects", "app")
	if runtime.GOOS != "windows" {
		root = filepath.Join("/projects", "app")
	}

	cases := []struct {
		name  string
		elems []string
	}{
		{"single segment", []string{"Makefile"}},
		{"nested segments", []string{"deploy", "k8s.yaml"}},
		{"manifest-supplied layer path", []string{"internal/core/domain", "order.go"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := domain.SafeJoin(root, tc.elems...)
			if err != nil {
				t.Fatalf("SafeJoin(%v) rejected a legitimate destination: %v", tc.elems, err)
			}
			if !strings.HasPrefix(got, root) {
				t.Errorf("result %q is not under root %q", got, root)
			}
		})
	}
}

// TestSafeJoin_RejectsEscapes is the security assertion. Command arguments reach
// filenames throughout the generator, and before this the only thing stopping
// traversal was that the template lookup happened to fail first.
func TestSafeJoin_RejectsEscapes(t *testing.T) {
	root := filepath.Join("C:", "projects", "app")
	if runtime.GOOS != "windows" {
		root = filepath.Join("/projects", "app")
	}

	cases := []struct {
		name  string
		elems []string
	}{
		{"parent traversal", []string{"..", "..", "etc", "passwd"}},
		{"traversal inside a segment", []string{"../../../pwned.yaml"}},
		{"windows separator traversal", []string{`..\..\pwned.yaml`}},
		{"separator smuggled into a filename", []string{"deploy", "../../pwned.yaml"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := domain.SafeJoin(root, tc.elems...)
			if err == nil {
				t.Fatalf("SafeJoin(%v) produced %q instead of refusing to leave the project root", tc.elems, got)
			}
			if !errors.Is(err, domain.ErrPathEscape) {
				t.Errorf("expected ErrPathEscape so callers can react, got %v", err)
			}
			if got != "" {
				t.Errorf("a rejected join must return no path, got %q", got)
			}
		})
	}
}

func TestValidateDeployTarget(t *testing.T) {
	if err := domain.ValidateDeployTarget("k8s"); err != nil {
		t.Errorf("k8s has a template and must be accepted, got %v", err)
	}

	// Advertised in the README but never implemented. Accepting them silently
	// produced a template-not-found error that quoted the user's input back at
	// them without explaining that the target was never supported.
	for _, target := range []string{"helm", "lambda", "docker", "../../etc"} {
		t.Run(target, func(t *testing.T) {
			err := domain.ValidateDeployTarget(target)
			if err == nil {
				t.Fatalf("unsupported target %q was accepted", target)
			}
			if !errors.Is(err, domain.ErrUnsupportedStack) {
				t.Errorf("expected ErrUnsupportedStack, got %v", err)
			}
			if !strings.Contains(err.Error(), "supported") {
				t.Errorf("message must name what is supported, got %v", err)
			}
		})
	}
}
