package scanner_test

import (
	"context"
	"strings"
	"testing"

	"github.com/muhananaufal/go-aether/internal/adapter/scanner"
	"github.com/muhananaufal/go-aether/internal/core/domain"
)

// findingsByRule indexes a report so assertions can name the rule they mean.
func findingsByRule(report *domain.AuditReport) map[string]domain.AuditFinding {
	out := map[string]domain.AuditFinding{}
	for _, f := range report.Findings {
		out[f.Rule] = f
	}
	return out
}

func TestGoCodeAuditor_FindsProductionHazards(t *testing.T) {
	root := writeTree(t, map[string]string{
		".env": "DB_PASS=hunter2\n",
		// .gitignore deliberately absent: the secret is exposed.
		"cmd/api/main.go": `package main

import (
	"database/sql"
	"net/http"
)

func main() {
	db, _ := sql.Open("postgres", "dsn")
	_ = db
	srv := &http.Server{Addr: ":8080"}
	_ = srv.ListenAndServe()
}
`,
		"internal/service/order_service.go": `package service

import (
	"context"
	"fmt"
)

func Process(ctx context.Context, id string) error {
	fmt.Println("processing", id)
	inner := context.Background()
	_ = inner
	if id == "" {
		panic("empty id")
	}
	return nil
}
`,
	})

	report, err := scanner.NewGoCodeAuditor().Audit(context.Background(), root)
	if err != nil {
		t.Fatalf("Audit failed: %v", err)
	}

	found := findingsByRule(report)
	expected := []string{
		"secrets/env-not-ignored",
		"db/unbounded-pool",
		"http/no-server-timeouts",
		"k8s/no-readiness-probe",
		"k8s/no-liveness-probe",
		"test/no-tests",
		"quality/no-linter-config",
		"reliability/panic-in-library",
		"context/detached-context",
		"observability/unstructured-log",
	}
	for _, rule := range expected {
		if _, ok := found[rule]; !ok {
			t.Errorf("rule %s did not fire on a project that plainly violates it", rule)
		}
	}

	// A finding that only names the problem tells a newcomer they are wrong
	// without telling them what right looks like.
	for rule, finding := range found {
		if strings.TrimSpace(finding.Hint) == "" {
			t.Errorf("finding %s carries no hint", rule)
		}
	}

	if !report.HasBlocking() {
		t.Error("an exposed .env did not register as blocking")
	}
}

// TestGoCodeAuditor_DoesNotFlagStartupCode is the false-positive guard, and it
// exists because the first version of this auditor flagged go-aether's own
// generated output.
//
// A pool constructor that builds its own context with a timeout, and a config
// loader that logs before any logger exists, are both correct. Reporting them
// on a project's very first doctor run teaches the reader to ignore the whole
// report, which costs more than the rules are worth.
func TestGoCodeAuditor_DoesNotFlagStartupCode(t *testing.T) {
	root := writeTree(t, map[string]string{
		".gitignore": ".env\n",
		"pkg/database/db.go": `package database

import (
	"context"
	"database/sql"
	"time"
)

func NewDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", "dsn")
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(25)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}
	return db, nil
}
`,
		"pkg/config/config.go": `package config

import "log"

func Load() error {
	log.Println("config: no .env file found; reading from the environment")
	return nil
}
`,
	})

	report, err := scanner.NewGoCodeAuditor().Audit(context.Background(), root)
	if err != nil {
		t.Fatalf("Audit failed: %v", err)
	}

	found := findingsByRule(report)
	for _, rule := range []string{"context/detached-context", "observability/unstructured-log"} {
		if finding, fired := found[rule]; fired {
			t.Errorf("%s fired on correct startup code (%s); "+
				"the rule must require a context.Context parameter to be in scope",
				rule, finding.File)
		}
	}
}

func TestGoCodeAuditor_CleanProjectReportsNothing(t *testing.T) {
	root := writeTree(t, map[string]string{
		".gitignore":    ".env\n",
		".golangci.yml": "linters:\n  enable:\n    - govet\n",
		"internal/core/domain/order.go": `package domain

type Order struct{ ID string }
`,
		"internal/core/domain/order_test.go": `package domain

import "testing"

func TestOrder(t *testing.T) {
	if (Order{ID: "1"}).ID != "1" {
		t.Fatal("unexpected")
	}
}
`,
	})

	report, err := scanner.NewGoCodeAuditor().Audit(context.Background(), root)
	if err != nil {
		t.Fatalf("Audit failed: %v", err)
	}

	if len(report.Findings) != 0 {
		var rules []string
		for _, f := range report.Findings {
			rules = append(rules, f.Rule)
		}
		t.Errorf("a clean project produced %d finding(s): %s",
			len(report.Findings), strings.Join(rules, ", "))
	}
}

func TestGoCodeAuditor_SortedPutsCriticalFirst(t *testing.T) {
	report := &domain.AuditReport{}
	report.Add(domain.AuditFinding{Severity: domain.AuditInfo, Rule: "a/info"})
	report.Add(domain.AuditFinding{Severity: domain.AuditCritical, Rule: "z/critical"})
	report.Add(domain.AuditFinding{Severity: domain.AuditWarn, Rule: "m/warn"})

	sorted := report.Sorted()
	if sorted[0].Severity != domain.AuditCritical {
		t.Errorf("the most urgent finding must be printed nearest the prompt, got %s first", sorted[0].Rule)
	}
	if sorted[len(sorted)-1].Severity != domain.AuditInfo {
		t.Errorf("suggestions must sort last, got %s", sorted[len(sorted)-1].Rule)
	}
}
