package domain

import (
	"fmt"
	"path/filepath"
	"strings"
)

// SafeJoin joins elements onto root and refuses any result that escapes it.
//
// Command arguments become path segments throughout the generator, and most are
// validated as Go identifiers first, which happens to block traversal. A few are
// not: --db, --target and similar free-form values flow into filenames directly.
// Today those are stopped only because the template lookup fails first, which is
// an accident of ordering rather than a boundary.
//
// This is that boundary, stated explicitly.
func SafeJoin(root string, elems ...string) (string, error) {
	for _, elem := range elems {
		if elem == "" {
			continue
		}
		// Reject separators outright. A destination segment is a filename or a
		// directory name, never a path, so anything containing a separator is
		// either a mistake or an attempt.
		if strings.ContainsAny(elem, `/\`) && !isKnownRelativeLayer(elem) {
			return "", fmt.Errorf("%w: segment %q contains a path separator", ErrPathEscape, elem)
		}
	}

	joined := filepath.Join(append([]string{root}, elems...)...)

	rel, err := filepath.Rel(root, joined)
	if err != nil {
		return "", fmt.Errorf("%w: cannot relate %q to %q", ErrPathEscape, joined, root)
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("%w: %q resolves outside %q", ErrPathEscape, joined, root)
	}

	return joined, nil
}

// isKnownRelativeLayer allows the multi-segment paths that come from the
// manifest itself, such as "internal/core/domain".
//
// Those are configuration the user wrote, not arguments an attacker supplies,
// and rejecting them would make SafeJoin unusable at exactly the call sites that
// need it most. They still pass the containment check below, which is the part
// that actually enforces the boundary.
func isKnownRelativeLayer(elem string) bool {
	normalized := filepath.ToSlash(elem)
	if strings.Contains(normalized, "..") {
		return false
	}
	return !filepath.IsAbs(elem) && !strings.HasPrefix(normalized, "/")
}
