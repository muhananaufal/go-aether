package domain

import "sort"

// AuditSeverity ranks a finding by how much attention it deserves.
type AuditSeverity int

const (
	// AuditCritical is reserved for findings that leak secrets or lose data.
	AuditCritical AuditSeverity = iota
	// AuditWarn marks a practice that holds up in development and fails under
	// production load: the class of defect a newcomer cannot yet foresee.
	AuditWarn
	// AuditInfo is a suggestion, not a problem.
	AuditInfo
)

func (s AuditSeverity) String() string {
	switch s {
	case AuditCritical:
		return "CRITICAL"
	case AuditWarn:
		return "WARN"
	default:
		return "INFO"
	}
}

// Icon renders the severity for terminal output.
func (s AuditSeverity) Icon() string {
	switch s {
	case AuditCritical:
		return "🔴"
	case AuditWarn:
		return "🟡"
	default:
		return "🔵"
	}
}

// AuditFinding is one observation about a user's project.
//
// Hint is not optional in spirit. A finding that only names the problem tells a
// newcomer they are wrong without telling them what right looks like, which is
// the least useful thing a mentor can do.
type AuditFinding struct {
	Severity AuditSeverity
	Rule     string
	Message  string
	Hint     string
	File     string
}

// AuditReport collects findings from a project inspection.
type AuditReport struct {
	Findings  []AuditFinding
	FilesSeen int
}

// Add appends a finding.
func (r *AuditReport) Add(f AuditFinding) {
	r.Findings = append(r.Findings, f)
}

// Sorted returns findings worst-first, then by rule, so the output order is
// stable between runs and the most urgent line is the one nearest the prompt.
func (r *AuditReport) Sorted() []AuditFinding {
	out := make([]AuditFinding, len(r.Findings))
	copy(out, r.Findings)

	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Severity != out[j].Severity {
			return out[i].Severity < out[j].Severity
		}
		if out[i].Rule != out[j].Rule {
			return out[i].Rule < out[j].Rule
		}
		return out[i].File < out[j].File
	})
	return out
}

// Count returns how many findings carry a given severity.
func (r *AuditReport) Count(severity AuditSeverity) int {
	n := 0
	for _, f := range r.Findings {
		if f.Severity == severity {
			n++
		}
	}
	return n
}

// HasBlocking reports whether anything found would be irresponsible to deploy.
func (r *AuditReport) HasBlocking() bool {
	return r.Count(AuditCritical) > 0
}
