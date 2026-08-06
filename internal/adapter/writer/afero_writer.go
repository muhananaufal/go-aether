package writer

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/muhananaufal/go-aether/internal/core/domain"
	"github.com/muhananaufal/go-aether/internal/core/port"
	"github.com/spf13/afero"
)

// AferoWriter implements port.FileWriter using spf13/afero to support OS disks and in-memory testing.
type AferoWriter struct {
	fs afero.Fs
}

// NewAferoWriter initializes a writer backed by the provided afero file system interface.
func NewAferoWriter(fileSystem afero.Fs) *AferoWriter {
	if fileSystem == nil {
		fileSystem = afero.NewOsFs()
	}
	return &AferoWriter{fs: fileSystem}
}

// WriteFile writes data to disk, managing conflicts, backup files, and dry-run notifications.
func (w *AferoWriter) WriteFile(ctx context.Context, targetPath string, content []byte, overwrite, dryRun bool) error {
	if dryRun {
		// In dry-run mode, we abort execution before touching physical media or allocating OS descriptors.
		return nil
	}

	exists, err := w.Exists(targetPath)
	if err != nil {
		return err
	}

	if exists && !overwrite {
		return fmt.Errorf("%w: %s", domain.ErrFileConflict, targetPath)
	}

	if exists && overwrite {
		// Create a defensive backup (.bak) prior to overwriting target file
		existingContent, readErr := afero.ReadFile(w.fs, targetPath)
		if readErr == nil {
			_ = afero.WriteFile(w.fs, targetPath+".bak", existingContent, 0644)
		}
	}

	dir := filepath.Dir(targetPath)
	if err := w.MkdirAll(dir); err != nil {
		return fmt.Errorf("failed to create target directory %s: %w", dir, err)
	}

	if err := afero.WriteFile(w.fs, targetPath, content, 0644); err != nil {
		return fmt.Errorf("%w: %v", domain.ErrWriteFailed, err)
	}

	return nil
}

// Exists checks whether a file or directory path exists within the filesystem.
func (w *AferoWriter) Exists(path string) (bool, error) {
	exists, err := afero.Exists(w.fs, path)
	if err != nil {
		return false, fmt.Errorf("failed to inspect filesystem path %s: %w", path, err)
	}
	return exists, nil
}

// ReadFile retrieves binary contents from an existing filesystem file.
func (w *AferoWriter) ReadFile(path string) ([]byte, error) {
	data, err := afero.ReadFile(w.fs, path)
	if err != nil {
		if errorsIsNotExist(err) {
			return nil, fmt.Errorf("%w: %s", domain.ErrManifestNotFound, path)
		}
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}
	return data, nil
}

// DeleteFile removes a file from disk, typically during atomic transaction rollback.
func (w *AferoWriter) DeleteFile(path string) error {
	exists, err := w.Exists(path)
	if err != nil || !exists {
		return nil
	}
	return w.fs.Remove(path)
}

// MkdirAll creates a directory tree recursively with proper permissions.
func (w *AferoWriter) MkdirAll(dirPath string) error {
	return w.fs.MkdirAll(dirPath, 0755)
}

// errorsIsNotExist evaluates file existence errors across standard and afero boundaries.
func errorsIsNotExist(err error) bool {
	return err == fs.ErrNotExist
}

// TransactionalBuffer implements an in-memory buffer that rolls back written files upon encountering an I/O failure.
type TransactionalBuffer struct {
	baseWriter port.FileWriter
	pending    []bufferedFile
	written    []string
}

type bufferedFile struct {
	path      string
	content   []byte
	overwrite bool
	dryRun    bool
}

// NewTransactionalBuffer constructs a transactional wrapper around an existing FileWriter.
func NewTransactionalBuffer(base port.FileWriter) *TransactionalBuffer {
	return &TransactionalBuffer{
		baseWriter: base,
		pending:    make([]bufferedFile, 0),
		written:    make([]string, 0),
	}
}

// Stage enqueues a file operation into memory without executing immediate disk mutations.
func (t *TransactionalBuffer) Stage(path string, content []byte, overwrite, dryRun bool) {
	t.pending = append(t.pending, bufferedFile{
		path:      path,
		content:   content,
		overwrite: overwrite,
		dryRun:    dryRun,
	})
}

// Commit executes all staged write operations sequentially. If any step fails, an atomic rollback is performed.
func (t *TransactionalBuffer) Commit(ctx context.Context) error {
	for _, bf := range t.pending {
		if err := t.baseWriter.WriteFile(ctx, bf.path, bf.content, bf.overwrite, bf.dryRun); err != nil {
			t.Rollback()
			return fmt.Errorf("%w: error writing %s: %v (transaction rolled back)", domain.ErrWriteFailed, bf.path, err)
		}
		if !bf.dryRun {
			t.written = append(t.written, bf.path)
		}
	}
	t.pending = nil
	return nil
}

// Rollback iterates backwards over successfully written files and removes them to leave a clean git tree.
func (t *TransactionalBuffer) Rollback() {
	for i := len(t.written) - 1; i >= 0; i-- {
		_ = t.baseWriter.DeleteFile(t.written[i])
	}
	t.written = nil
	t.pending = nil
}
