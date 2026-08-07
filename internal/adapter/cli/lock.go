package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gofrs/flock"
)

var globalFileLock *flock.Flock

// acquireProjectLock acquires a global file lock for the project to prevent concurrent UOW race conditions.
func acquireProjectLock(ctx context.Context, startDir string) error {
	// Look for aether.yaml to place the lock next to it
	curr := startDir
	manifestPath := ""
	for i := 0; i < 12; i++ {
		candidate := filepath.Join(curr, "aether.yaml")
		if _, err := os.Stat(candidate); err == nil {
			manifestPath = candidate
			break
		}
		if _, err := os.Stat(filepath.Join(curr, ".git")); err == nil && i > 0 {
			break
		}
		parent := filepath.Dir(curr)
		if parent == curr {
			break
		}
		curr = parent
	}

	var lockPath string
	if manifestPath != "" {
		lockPath = filepath.Join(filepath.Dir(manifestPath), ".aether.lock")
	} else {
		lockPath = filepath.Join(startDir, ".aether.lock")
	}

	globalFileLock = flock.New(lockPath)

	locked, err := globalFileLock.TryLockContext(ctx, 10*time.Millisecond)
	if err != nil {
		return fmt.Errorf("failed to attempt project lock: %w", err)
	}

	if !locked {
		fmt.Println("⏳ Another go-aether process is currently running. Waiting for lock...")
		// Wait up to 10 seconds for the lock
		waitCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		
		locked, err = globalFileLock.TryLockContext(waitCtx, 500*time.Millisecond)
		if err != nil || !locked {
			return fmt.Errorf("timeout waiting for project lock at %s. Ensure no other go-aether process is stuck", lockPath)
		}
	}

	return nil
}

// releaseProjectLock unlocks the global file lock
func releaseProjectLock() {
	if globalFileLock != nil {
		_ = globalFileLock.Unlock()
	}
}
