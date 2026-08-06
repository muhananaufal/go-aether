package service

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/muhananaufal/go-aether/internal/core/domain"
)

// InitProject bootstraps a new greenfield project architecture and saves the SSOT manifest.
func (s *AetherScaffoldService) InitProject(ctx context.Context, destDir, projectName, moduleName, arch, dbDriver, router string, dryRun bool) error {
	manifest := domain.NewDefaultManifest(projectName, moduleName, "1.23", arch, dbDriver, router)
	
	if err := manifest.Validate(); err != nil {
		return fmt.Errorf("invalid default manifest parameters: %w", err)
	}

	manifestPath := filepath.Join(destDir, "aether.yaml")
	exists, err := s.fs.Exists(manifestPath)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("%w: %s", domain.ErrFileConflict, manifestPath)
	}

	if !dryRun {
		if err := s.resolver.Save(ctx, destDir, manifest, s.fs); err != nil {
			return fmt.Errorf("failed writing initial manifest: %w", err)
		}
	}

	return nil
}

// AdoptProject scans an existing repository to construct an initial aether.yaml for brownfield adoption.
func (s *AetherScaffoldService) AdoptProject(ctx context.Context, destDir string, scan, dryRun bool) error {
	manifestPath := filepath.Join(destDir, "aether.yaml")
	exists, err := s.fs.Exists(manifestPath)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("%w: repository already contains aether.yaml", domain.ErrFileConflict)
	}

	manifest := domain.NewDefaultManifest(filepath.Base(destDir), "github.com/adopted/"+filepath.Base(destDir), "1.23", "hexagonal", "postgres", "chi")
	manifest.Architecture.Mode = "brownfield"
	manifest.Project.CreatedAt = time.Now().Format(time.RFC3339)

	if scan {
		manifest.Meta.AnomalyMode = true
		manifest.Meta.LegacyNotes = "Scanned and adopted via go-aether CLI; BYO dependencies enabled."
	}

	if !dryRun {
		if err := s.resolver.Save(ctx, destDir, manifest, s.fs); err != nil {
			return fmt.Errorf("failed saving adopted manifest: %w", err)
		}
	}

	return nil
}
