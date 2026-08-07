package domain

import "sort"

// LayoutKind names an architectural layer the scanner can recognise in a
// repository that was never organised with go-aether in mind.
type LayoutKind string

const (
	LayerHandler    LayoutKind = "handler"
	LayerService    LayoutKind = "service"
	LayerRepository LayoutKind = "repository"
	LayerDomain     LayoutKind = "domain"
)

// minConfidentScore is the point below which a guess stops being useful.
//
// Chosen so that a directory matched only by its name (0.5) is never applied
// silently: brownfield repositories are full of directories called "api" or
// "data" that mean something else entirely. Corroborating evidence from the
// import graph is required before the scanner will speak with any confidence.
const minConfidentScore = 0.75

// LayoutCandidate is one directory the scanner believes holds a given layer,
// together with the reasons it thinks so.
//
// Evidence is carried rather than discarded because the user has to approve the
// mapping, and "we found handlers in web/controllers" is not something anyone
// can sensibly agree to without being told why.
type LayoutCandidate struct {
	// Dir is relative to the scan root and always slash-separated, so it can be
	// written into aether.yaml unchanged on any operating system.
	Dir      string
	Kind     LayoutKind
	Score    float64
	Evidence []string
	GoFiles  int
}

// LayoutReport is the complete result of scanning a legacy repository.
type LayoutReport struct {
	Candidates []LayoutCandidate

	// FilesSeen and Truncated make the scan's own limits visible. A report that
	// silently stopped at 5000 files would otherwise look identical to one that
	// genuinely found nothing further, and the user would trust it equally.
	FilesSeen int
	Truncated bool

	// SkippedDirs records directories deliberately not descended into, so the
	// absence of a result there is explainable rather than mysterious.
	SkippedDirs []string
}

// BYODetection holds the identifiers of infrastructure clients an existing
// project already constructs.
//
// This is the difference between a generator that fits into a brownfield code
// base and one that fights it: without these, generated repositories open their
// own connection pool alongside the one main.go already built, and the service
// quietly runs at double the connection count until the database refuses.
type BYODetection struct {
	DBVar     string
	RedisVar  string
	LoggerVar string

	// Sources maps each identifier to the file it was found in, so the proposal
	// shown to the user can cite its evidence.
	Sources map[string]string
}

// IsEmpty reports whether anything reusable was detected at all.
func (b *BYODetection) IsEmpty() bool {
	return b == nil || (b.DBVar == "" && b.RedisVar == "" && b.LoggerVar == "")
}

// Best returns the strongest candidate for a layer.
//
// The boolean result is false when nothing crossed minConfidentScore, which the
// caller must treat as "ask the user" rather than "use the default": applying a
// default silently is how a generator ends up writing files into a directory
// that means something completely different in that code base.
func (r *LayoutReport) Best(kind LayoutKind) (LayoutCandidate, bool) {
	var best LayoutCandidate
	found := false

	for _, c := range r.Candidates {
		if c.Kind != kind {
			continue
		}
		if !found || c.Score > best.Score || (c.Score == best.Score && c.GoFiles > best.GoFiles) {
			best, found = c, true
		}
	}

	if !found || best.Score < minConfidentScore {
		return LayoutCandidate{}, false
	}
	return best, true
}

// Confidence summarises how much of the expected architecture was recognised,
// as a fraction of the four layers go-aether generates into.
func (r *LayoutReport) Confidence() float64 {
	kinds := []LayoutKind{LayerHandler, LayerService, LayerRepository, LayerDomain}
	matched := 0
	for _, k := range kinds {
		if _, ok := r.Best(k); ok {
			matched++
		}
	}
	return float64(matched) / float64(len(kinds))
}

// Sorted returns candidates ordered strongest first, then by directory, so the
// proposal presented to the user is stable between runs. Map iteration and
// filesystem walk order are both unstable enough to reorder an otherwise
// identical report, which makes it impossible to review a diff of one.
func (r *LayoutReport) Sorted() []LayoutCandidate {
	out := make([]LayoutCandidate, len(r.Candidates))
	copy(out, r.Candidates)

	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Score != out[j].Score {
			return out[i].Score > out[j].Score
		}
		if out[i].Kind != out[j].Kind {
			return out[i].Kind < out[j].Kind
		}
		return out[i].Dir < out[j].Dir
	})
	return out
}
