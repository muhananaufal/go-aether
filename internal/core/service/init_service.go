package service

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/charmbracelet/huh"
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

	var moduleName = "github.com/adopted/" + filepath.Base(destDir)
	var goVersion = "1.23"
	var arch = "hexagonal"
	var router = "chi"
	var db = "postgres"

	if scan && !dryRun {
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("What is the Go module name?").
					Description("e.g. github.com/user/project").
					Value(&moduleName),
				
				huh.NewInput().
					Title("What is the Go version?").
					Description("e.g. 1.23").
					Value(&goVersion),
			),
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Which architectural pattern best describes this project?").
					Options(
						huh.NewOption("Hexagonal Architecture", "hexagonal"),
						huh.NewOption("Clean Architecture", "clean"),
						huh.NewOption("Standard MVC", "mvc"),
						huh.NewOption("Unknown / Legacy Anomaly", "anomaly"),
					).
					Value(&arch),
				
				huh.NewSelect[string]().
					Title("Primary HTTP Router").
					Options(
						huh.NewOption("go-chi/chi", "chi"),
						huh.NewOption("gin-gonic", "gin"),
						huh.NewOption("gorilla/mux", "mux"),
						huh.NewOption("fiber", "fiber"),
						huh.NewOption("standard net/http", "stdlib"),
					).
					Value(&router),
					
				huh.NewSelect[string]().
					Title("Primary Database").
					Options(
						huh.NewOption("PostgreSQL", "postgres"),
						huh.NewOption("MySQL", "mysql"),
						huh.NewOption("MongoDB", "mongo"),
						huh.NewOption("None", "none"),
					).
					Value(&db),
			),
		)

		err := form.Run()
		if err != nil {
			return fmt.Errorf("adoption aborted: %w", err)
		}
	}

	manifest := domain.NewDefaultManifest(filepath.Base(destDir), moduleName, goVersion, arch, db, router)
	manifest.Architecture.Mode = "brownfield"
	manifest.Project.CreatedAt = time.Now().Format(time.RFC3339)

	if scan {
		manifest.Meta.AnomalyMode = arch == "anomaly"
		manifest.Meta.LegacyNotes = "Scanned and adopted via go-aether CLI; BYO dependencies enabled."
	}

	if !dryRun {
		if err := s.resolver.Save(ctx, destDir, manifest, s.fs); err != nil {
			return fmt.Errorf("failed saving adopted manifest: %w", err)
		}
		fmt.Printf("\n✨ Successfully generated aether.yaml with anomaly_mode=%v\n", manifest.Meta.AnomalyMode)
	}

	return nil
}
