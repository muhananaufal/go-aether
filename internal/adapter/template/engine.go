package template

import (
	"bytes"
	"context"
	"fmt"
	"go/format"
	"io/fs"
	"path"
	"strings"
	"text/template"

	"github.com/muhananaufal/go-aether/internal/core/domain"
	"github.com/muhananaufal/go-aether/internal/core/port"
)

// goTemplateSuffix marks templates whose rendered output is Go source and must
// therefore pass through the formatter before it is allowed near a disk.
const goTemplateSuffix = ".go.tmpl"

// Compile-time interface compliance check
var _ port.TemplateEngine = (*StdEngine)(nil)

// StdEngine implements port.TemplateEngine using Go standard library text/template backed by fs.FS.
type StdEngine struct {
	embedFS fs.FS
	funcMap template.FuncMap
}

// NewStdEngine constructs a template generator engine wrapped around an embedded filesystem.
func NewStdEngine(fileSystem fs.FS) *StdEngine {
	funcMap := template.FuncMap{
		"toLower": strings.ToLower,
		"toUpper": strings.ToUpper,
		"toTitle": domain.ToTitleCase,
		"join":    strings.Join,
	}
	return &StdEngine{
		embedFS: fileSystem,
		funcMap: funcMap,
	}
}

// Render compiles and interpolates a designated template path with provided domain TemplateData.
func (e *StdEngine) Render(ctx context.Context, templatePath string, data *domain.TemplateData) ([]byte, error) {
	normalizedPath := path.Clean(filepathToSlash(templatePath))

	rawBytes, err := fs.ReadFile(e.embedFS, normalizedPath)
	if err != nil {
		return nil, fmt.Errorf("%w: unable to read embedded template %q: %v", domain.ErrTemplateMissing, normalizedPath, err)
	}

	tmpl, err := template.New(path.Base(normalizedPath)).Funcs(e.funcMap).Parse(string(rawBytes))
	if err != nil {
		return nil, fmt.Errorf("aether: template syntax parsing failure on %q: %w", normalizedPath, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("aether: template execution failed on %q: %w", normalizedPath, err)
	}

	// Go output passes through gofmt before anyone sees it. Two reasons, both
	// user-facing: unformatted source is rewritten by every editor on save and
	// pollutes the project's first commit, and a template that renders broken
	// syntax must fail here — with the template name and offending line — rather
	// than land on disk and surface later as an opaque compiler error.
	if strings.HasSuffix(normalizedPath, goTemplateSuffix) {
		formatted, fmtErr := format.Source(buf.Bytes())
		if fmtErr != nil {
			return nil, fmt.Errorf("%w: template %q: %v", domain.ErrGeneratedCodeInvalid, normalizedPath, fmtErr)
		}
		return formatted, nil
	}

	return buf.Bytes(), nil
}

// ListTemplates scans the filesystem within a specified pattern directory and returns valid template filenames.
func (e *StdEngine) ListTemplates(pattern string) ([]string, error) {
	var templates []string
	cleanPattern := filepathToSlash(pattern)

	err := fs.WalkDir(e.embedFS, cleanPattern, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".tmpl") {
			templates = append(templates, path)
		}
		return nil
	})
	if err != nil && err != fs.ErrNotExist {
		return nil, fmt.Errorf("failed walking template fs for pattern %s: %w", cleanPattern, err)
	}
	return templates, nil
}

// filepathToSlash converts Windows backslash paths to standard Unix slashes required by embed.FS.
func filepathToSlash(p string) string {
	return strings.ReplaceAll(p, "\\", "/")
}
