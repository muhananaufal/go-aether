package writer

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"runtime"

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

// maxWindowsPath is the classic MAX_PATH limit minus room for the temporary
// suffixes the toolchain appends. Long-path support exists but is off by default
// on most machines, and the error the OS returns when it is off names neither the
// limit nor the fix.
const maxWindowsPath = 240

// WriteFile writes data to disk, managing conflicts, backup files, and dry-run notifications.
func (w *AferoWriter) WriteFile(ctx context.Context, targetPath string, content []byte, overwrite, dryRun bool) error {
	if dryRun {
		// In dry-run mode, we abort execution before touching physical media or allocating OS descriptors.
		return nil
	}

	if runtime.GOOS == "windows" && len(targetPath) > maxWindowsPath {
		return fmt.Errorf(
			"%w: destination path is %d characters, over the %d Windows supports by default: %s\n"+
				"Move the project closer to the drive root, shorten the module name, or enable long paths "+
				"(Computer\\HKEY_LOCAL_MACHINE\\SYSTEM\\CurrentControlSet\\Control\\FileSystem\\LongPathsEnabled = 1)",
			domain.ErrWriteFailed, len(targetPath), maxWindowsPath, targetPath)
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

// BeginTransaction returns a TransactionalWriter decorator.
func (w *AferoWriter) BeginTransaction() port.TransactionalWriter {
	return NewUOWWriter(w)
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

// Size reports a file's length without reading it.
func (w *AferoWriter) Size(path string) (int64, error) {
	info, err := w.fs.Stat(path)
	if err != nil {
		return 0, fmt.Errorf("failed to stat %s: %w", path, err)
	}
	return info.Size(), nil
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
//
// A direct == comparison never matched: both os and afero return *fs.PathError,
// which wraps the sentinel rather than being it. The consequence was that a
// genuinely missing manifest surfaced as a generic read failure instead of
// ErrManifestNotFound, so callers could not tell "no project here" apart from
// "disk is broken".
func errorsIsNotExist(err error) bool {
	return errors.Is(err, fs.ErrNotExist)
}

// UOWWriter implements port.FileWriter as a Unit of Work decorator that buffers writes in memory.
// It commits them atomically and rolls back on failure to maintain a clean git tree.
type UOWWriter struct {
	baseWriter port.FileWriter
	pending    []bufferedFile
	written    []writtenFile
}

type bufferedFile struct {
	path      string
	content   []byte
	overwrite bool
	dryRun    bool
}

type writtenFile struct {
	path           string
	wasOverwritten bool
}

// NewUOWWriter constructs a transactional wrapper around an existing FileWriter.
func NewUOWWriter(base port.FileWriter) *UOWWriter {
	return &UOWWriter{
		baseWriter: base,
		pending:    make([]bufferedFile, 0),
		written:    make([]writtenFile, 0),
	}
}

// WriteFile intercepts file writes and stages them into memory. It returns nil.
func (t *UOWWriter) WriteFile(ctx context.Context, targetPath string, content []byte, overwrite, dryRun bool) error {
	t.pending = append(t.pending, bufferedFile{
		path:      targetPath,
		content:   content,
		overwrite: overwrite,
		dryRun:    dryRun,
	})
	return nil
}

// Commit executes all staged write operations sequentially. If any step fails, an atomic rollback is performed.
func (t *UOWWriter) Commit(ctx context.Context) error {
	for _, bf := range t.pending {
		exists, _ := t.baseWriter.Exists(bf.path)
		wasOverwritten := exists && bf.overwrite

		if err := t.baseWriter.WriteFile(ctx, bf.path, bf.content, bf.overwrite, bf.dryRun); err != nil {
			t.Rollback()
			return fmt.Errorf("%w: error writing %s: %v (transaction rolled back)", domain.ErrWriteFailed, bf.path, err)
		}
		if !bf.dryRun {
			t.written = append(t.written, writtenFile{path: bf.path, wasOverwritten: wasOverwritten})
		}
	}
	t.pending = nil
	return nil
}

// Rollback iterates backwards over successfully written files and removes them to leave a clean git tree.
func (t *UOWWriter) Rollback() {
	for i := len(t.written) - 1; i >= 0; i-- {
		wf := t.written[i]
		_ = t.baseWriter.DeleteFile(wf.path)

		if wf.wasOverwritten {
			if backupData, err := t.baseWriter.ReadFile(wf.path + ".bak"); err == nil {
				// Restore original content from .bak. Since the file was just deleted, overwrite=false is fine.
				_ = t.baseWriter.WriteFile(context.Background(), wf.path, backupData, false, false)
				_ = t.baseWriter.DeleteFile(wf.path + ".bak")
			}
		}
	}
	t.written = nil
	t.pending = nil
}

// Exists delegates to the base writer.
func (t *UOWWriter) Exists(path string) (bool, error) {
	return t.baseWriter.Exists(path)
}

// ReadFile delegates to the base writer.
func (t *UOWWriter) ReadFile(path string) ([]byte, error) {
	return t.baseWriter.ReadFile(path)
}

// Size delegates to the base writer.
func (t *UOWWriter) Size(path string) (int64, error) {
	return t.baseWriter.Size(path)
}

// DeleteFile delegates to the base writer.
func (t *UOWWriter) DeleteFile(path string) error {
	return t.baseWriter.DeleteFile(path)
}

// MkdirAll delegates to the base writer.
func (t *UOWWriter) MkdirAll(dirPath string) error {
	return t.baseWriter.MkdirAll(dirPath)
}

// BeginTransaction returns the current UOW instance (no nested transactions).
func (t *UOWWriter) BeginTransaction() port.TransactionalWriter {
	return t
}
