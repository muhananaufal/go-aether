package port

import (
	"context"

	"github.com/muhananaufal/go-aether/internal/core/domain"
)

// LayoutScanner inspects an existing repository and reports where its
// architectural layers appear to live.
//
// Implementations must be bounded. A brownfield repository can be a monorepo
// with millions of files, and a scanner that walks all of them turns `adopt`
// into a command nobody dares run twice.
type LayoutScanner interface {
	// Scan walks root and classifies directories. It returns a report even when
	// nothing is recognised: an empty report with FilesSeen populated is a useful
	// answer, whereas an error is not.
	//
	// The context must be honoured; a scan interrupted by cancellation returns
	// what it has along with ctx.Err().
	Scan(ctx context.Context, root string) (*domain.LayoutReport, error)
}

// BYODetector finds infrastructure clients an existing project already builds,
// so generated code can accept them rather than construct duplicates.
type BYODetector interface {
	// Detect inspects the entrypoints under root. Finding nothing is not an
	// error; it simply means the project has no reusable clients to inherit.
	Detect(ctx context.Context, root string) (*domain.BYODetection, error)
}
