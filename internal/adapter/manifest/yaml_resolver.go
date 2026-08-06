package manifest

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/muhananaufal/go-aether/internal/core/domain"
	"github.com/muhananaufal/go-aether/internal/core/port"
	"gopkg.in/yaml.v3"
)

const manifestFilename = "aether.yaml"

// YamlResolver implements port.ManifestResolver using yaml.v3 with recursive walk-up capabilities.
type YamlResolver struct {
	writer port.FileWriter
}

// NewYamlResolver constructs a resolver backed by an explicit filesystem writer.
func NewYamlResolver(writer port.FileWriter) *YamlResolver {
	return &YamlResolver{writer: writer}
}

// Resolve locates aether.yaml by walking up from startDir until hitting the root or repository boundary.
func (r *YamlResolver) Resolve(ctx context.Context, startDir string) (*domain.AetherManifest, string, error) {
	curr := startDir
	maxDepth := 12 // Guard against infinite filesystem recursion cycles

	for i := 0; i < maxDepth; i++ {
		candidate := filepath.Join(curr, manifestFilename)
		exists, err := r.writer.Exists(candidate)
		if err == nil && exists {
			data, readErr := r.writer.ReadFile(candidate)
			if readErr != nil {
				return nil, "", fmt.Errorf("%w: failed to read %s: %v", domain.ErrManifestCorrupt, candidate, readErr)
			}

			var m domain.AetherManifest
			if err := yaml.Unmarshal(data, &m); err != nil {
				return nil, "", fmt.Errorf("%w: yaml unmarshal failure: %v", domain.ErrManifestCorrupt, err)
			}

			if valErr := m.Validate(); valErr != nil {
				return nil, "", valErr
			}
			return &m, filepath.ToSlash(candidate), nil
		}

		// Stop traversal if we hit repository top level boundary (.git) without finding aether.yaml
		gitDir := filepath.Join(curr, ".git")
		gitExists, _ := r.writer.Exists(gitDir)
		if gitExists && i > 0 {
			break
		}

		parent := filepath.Dir(curr)
		if parent == curr || strings.TrimSpace(parent) == "" || parent == "/" || parent == "." {
			break
		}
		curr = parent
	}

	return nil, "", domain.ErrManifestNotFound
}

// Save validates and persists an AetherManifest out to disk in strict YAML structure.
func (r *YamlResolver) Save(ctx context.Context, destPath string, manifest *domain.AetherManifest, writer port.FileWriter) error {
	if err := manifest.Validate(); err != nil {
		return err
	}

	data, err := yaml.Marshal(manifest)
	if err != nil {
		return fmt.Errorf("failed to marshal manifest to yaml: %w", err)
	}

	target := destPath
	if !strings.HasSuffix(destPath, manifestFilename) {
		target = filepath.Join(destPath, manifestFilename)
	}

	if err := writer.WriteFile(ctx, target, data, true, false); err != nil {
		return fmt.Errorf("failed to write manifest file: %w", err)
	}

	return nil
}
