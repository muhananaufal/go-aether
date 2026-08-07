package service_test

import (
	"bytes"
	"context"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/muhananaufal/go-aether/internal/adapter/manifest"
	"github.com/muhananaufal/go-aether/internal/adapter/template"
	"github.com/muhananaufal/go-aether/internal/adapter/writer"
	"github.com/muhananaufal/go-aether/internal/core/service"
	"github.com/spf13/afero"
)

func TestAetherScaffoldService_E2E_InitAndMakeModule(t *testing.T) {
	memFS := afero.NewMemMapFs()
	w := writer.NewAferoWriter(memFS)
	resolver := manifest.NewYamlResolver(w)

	// Mock embedded templates matching real layout
	mockEmbed := fstest.MapFS{
		"common/aether_yaml.tmpl":               &fstest.MapFile{Data: []byte("version: \"1\"")},
		"common/main.go.tmpl":                   &fstest.MapFile{Data: []byte("package main")},
		"common/config_viper.go.tmpl":           &fstest.MapFile{Data: []byte("package config")},
		"common/postgres.go.tmpl":               &fstest.MapFile{Data: []byte("package database")},
		"common/Makefile.tmpl":                  &fstest.MapFile{Data: []byte("build:")},
		"common/Dockerfile.tmpl":                &fstest.MapFile{Data: []byte("FROM alpine")},
		"common/dockerignore.tmpl":              &fstest.MapFile{Data: []byte("bin/")},
		"common/env_example.tmpl":               &fstest.MapFile{Data: []byte("APP_ENV=")},
		"hexagonal/domain.go.tmpl":              &fstest.MapFile{Data: []byte("package domain\n// domain for {{ .ModuleNameTitle }}")},
		"hexagonal/port.go.tmpl":                &fstest.MapFile{Data: []byte("package port\n// port for {{ .ModuleNameTitle }}")},
		"hexagonal/service.go.tmpl":             &fstest.MapFile{Data: []byte("package service\n// service for {{ .ModuleNameTitle }}")},
		"hexagonal/handler_http.go.tmpl":        &fstest.MapFile{Data: []byte("package http\n// http handler for {{ .ModuleNameTitle }}")},
		"hexagonal/repository_postgres.go.tmpl": &fstest.MapFile{Data: []byte("package repository\n// repo for {{ .ModuleNameTitle }}")},
	}

	engine := template.NewStdEngine(mockEmbed)
	svc := service.NewAetherScaffoldService(engine, resolver, w)
	ctx := context.Background()
	projectDir := "/projects/e2e-app"
	_ = w.MkdirAll(projectDir)

	// 1. Initialize project
	if err := svc.InitProject(ctx, projectDir, "e2e-app", "github.com/e2e/app", "hexagonal", "postgres", "chi", false); err != nil {
		t.Fatalf("InitProject failed: %v", err)
	}

	manifestFile := filepath.Join(projectDir, "aether.yaml")
	exists, _ := w.Exists(manifestFile)
	if !exists {
		t.Fatalf("expected manifest file %s to be created during init", manifestFile)
	}

	// 2. Make module 'order'
	if err := svc.MakeModule(ctx, projectDir, "order", []string{"http"}, false, false, false, false); err != nil {
		t.Fatalf("MakeModule failed: %v", err)
	}

	// Verify generated layer files exist in memory filesystem
	expectedFiles := []string{
		"/projects/e2e-app/internal/core/domain/order.go",
		"/projects/e2e-app/internal/core/port/order_port.go",
		"/projects/e2e-app/internal/core/service/order_service.go",
		"/projects/e2e-app/internal/adapter/handler/http/order_handler.go",
		"/projects/e2e-app/internal/adapter/repository/order_repository.go",
	}

	for _, ef := range expectedFiles {
		found, err := w.Exists(ef)
		if err != nil || !found {
			t.Errorf("expected generated file %q to exist in filesystem, but it was absent", ef)
		}
	}

	// 3. Run Doctor diagnostic checks
	var docBuf bytes.Buffer
	if err := svc.RunDoctor(ctx, projectDir, false, &docBuf); err != nil {
		t.Fatalf("RunDoctor failed: %v", err)
	}
	output := docBuf.String()
	if !bytes.Contains([]byte(output), []byte("Doctor check complete")) {
		t.Errorf("expected doctor check success message, got:\n%s", output)
	}
}

func TestAetherScaffoldService_Phase13_TestingAndQASuite(t *testing.T) {
	memFS := afero.NewMemMapFs()
	w := writer.NewAferoWriter(memFS)
	resolver := manifest.NewYamlResolver(w)

	mockEmbed := fstest.MapFS{
		"common/aether_yaml.tmpl":             &fstest.MapFile{Data: []byte("version: \"1\"")},
		"common/main.go.tmpl":                 &fstest.MapFile{Data: []byte("package main")},
		"common/config_viper.go.tmpl":         &fstest.MapFile{Data: []byte("package config")},
		"common/postgres.go.tmpl":             &fstest.MapFile{Data: []byte("package database")},
		"common/Makefile.tmpl":                &fstest.MapFile{Data: []byte("build:")},
		"common/Dockerfile.tmpl":              &fstest.MapFile{Data: []byte("FROM alpine")},
		"common/dockerignore.tmpl":            &fstest.MapFile{Data: []byte("bin/")},
		"common/env_example.tmpl":             &fstest.MapFile{Data: []byte("APP_ENV=")},
		"plugins/stress_k6.js.tmpl":           &fstest.MapFile{Data: []byte("// k6 stress test for {{ .ModuleName }}")},
		"plugins/stress_vegeta.sh.tmpl":       &fstest.MapFile{Data: []byte("# vegeta stress test for {{ .ModuleNameTitle }}")},
		"plugins/chaos.go.tmpl":               &fstest.MapFile{Data: []byte("package middleware\n// chaos middleware")},
		"plugins/fuzz_test.go.tmpl":           &fstest.MapFile{Data: []byte("package fuzz_test\n// fuzz test")},
		"plugins/benchmark_test.go.tmpl":      &fstest.MapFile{Data: []byte("package bench_test\n// benchmark test")},
		"plugins/testcontainers_test.go.tmpl": &fstest.MapFile{Data: []byte("package integration_test\n// testcontainers")},
		"plugins/mutation_test.ps1.tmpl":      &fstest.MapFile{Data: []byte("# mutation test script")},
		"plugins/mock.go.tmpl":                &fstest.MapFile{Data: []byte("package mocks\n// mock for {{ .ModuleNameTitle }}")},
	}

	engine := template.NewStdEngine(mockEmbed)
	svc := service.NewAetherScaffoldService(engine, resolver, w)
	ctx := context.Background()
	projectDir := "/projects/qa-app"
	_ = w.MkdirAll(projectDir)

	if err := svc.InitProject(ctx, projectDir, "qa-app", "github.com/qa/app", "hexagonal", "postgres", "chi", false); err != nil {
		t.Fatalf("InitProject failed: %v", err)
	}

	// 1. TestStress k6 & vegeta
	if err := svc.TestStress(ctx, projectDir, "k6", false, false); err != nil {
		t.Fatalf("TestStress k6 failed: %v", err)
	}
	exists, _ := w.Exists("/projects/qa-app/tests/stress/load_test.js")
	if !exists {
		t.Error("expected load_test.js to exist")
	}

	if err := svc.TestStress(ctx, projectDir, "vegeta", false, false); err != nil {
		t.Fatalf("TestStress vegeta failed: %v", err)
	}
	exists, _ = w.Exists("/projects/qa-app/tests/stress/vegeta_attack.sh")
	if !exists {
		t.Error("expected vegeta_attack.sh to exist")
	}

	// 2. TestChaos
	if err := svc.TestChaos(ctx, projectDir, false, false); err != nil {
		t.Fatalf("TestChaos failed: %v", err)
	}
	exists, _ = w.Exists("/projects/qa-app/pkg/middleware/chaos.go")
	if !exists {
		t.Error("expected pkg/middleware/chaos.go to exist")
	}

	// 3. TestFuzz
	if err := svc.TestFuzz(ctx, projectDir, "json", false, false); err != nil {
		t.Fatalf("TestFuzz failed: %v", err)
	}
	exists, _ = w.Exists("/projects/qa-app/tests/fuzz/fuzz_test.go")
	if !exists {
		t.Error("expected tests/fuzz/fuzz_test.go to exist")
	}

	// 4. TestBenchmark
	if err := svc.TestBenchmark(ctx, projectDir, false, false); err != nil {
		t.Fatalf("TestBenchmark failed: %v", err)
	}
	exists, _ = w.Exists("/projects/qa-app/tests/bench/alloc_benchmark_test.go")
	if !exists {
		t.Error("expected tests/bench/alloc_benchmark_test.go to exist")
	}

	// 5. TestContainer
	if err := svc.TestContainer(ctx, projectDir, "postgres-redis", false, false); err != nil {
		t.Fatalf("TestContainer failed: %v", err)
	}
	exists, _ = w.Exists("/projects/qa-app/tests/integration/postgres_redis_test.go")
	if !exists {
		t.Error("expected tests/integration/postgres_redis_test.go to exist")
	}

	// 6. TestMutation
	if err := svc.TestMutation(ctx, projectDir, false, false); err != nil {
		t.Fatalf("TestMutation failed: %v", err)
	}
	exists, _ = w.Exists("/projects/qa-app/scripts/mutation_test.ps1")
	if !exists {
		t.Error("expected scripts/mutation_test.ps1 to exist")
	}

	// 7. MakeMock
	if err := svc.MakeMock(ctx, projectDir, "PaymentGateway", false, false); err != nil {
		t.Fatalf("MakeMock failed: %v", err)
	}
	exists, _ = w.Exists("/projects/qa-app/mocks/mock_paymentgateway.go")
	if !exists {
		t.Error("expected mocks/mock_paymentgateway.go to exist")
	}
}

func TestAetherScaffoldService_Phase14_SQLC_GRPC_MultiTenancy(t *testing.T) {
	memFS := afero.NewMemMapFs()
	w := writer.NewAferoWriter(memFS)
	resolver := manifest.NewYamlResolver(w)

	mockEmbed := fstest.MapFS{
		"common/aether_yaml.tmpl":        &fstest.MapFile{Data: []byte("version: \"1\"")},
		"common/main.go.tmpl":            &fstest.MapFile{Data: []byte("package main")},
		"common/config_viper.go.tmpl":    &fstest.MapFile{Data: []byte("package config")},
		"common/postgres.go.tmpl":        &fstest.MapFile{Data: []byte("package database")},
		"common/Makefile.tmpl":           &fstest.MapFile{Data: []byte("build:")},
		"common/Dockerfile.tmpl":         &fstest.MapFile{Data: []byte("FROM alpine")},
		"common/dockerignore.tmpl":       &fstest.MapFile{Data: []byte("bin/")},
		"common/env_example.tmpl":        &fstest.MapFile{Data: []byte("APP_ENV=")},
		"plugins/sqlc_yaml.tmpl":         &fstest.MapFile{Data: []byte("version: \"2\"")},
		"plugins/sqlc_schema.sql.tmpl":   &fstest.MapFile{Data: []byte("CREATE TABLE accounts ();")},
		"plugins/sqlc_query.sql.tmpl":    &fstest.MapFile{Data: []byte("-- name: GetAccountByID :one\nSELECT 1;")},
		"plugins/grpc_stream.proto.tmpl": &fstest.MapFile{Data: []byte("syntax = \"proto3\";")},
		"plugins/grpc_stream.go.tmpl":    &fstest.MapFile{Data: []byte("package grpc\n// stream handler")},
		"plugins/grpc_gateway.go.tmpl":   &fstest.MapFile{Data: []byte("package http\n// grpc gateway")},
		"plugins/tenant_context.go.tmpl": &fstest.MapFile{Data: []byte("package tenant\n// tenant context")},
	}

	engine := template.NewStdEngine(mockEmbed)
	svc := service.NewAetherScaffoldService(engine, resolver, w)
	ctx := context.Background()
	projectDir := "/projects/enterprise-app"
	_ = w.MkdirAll(projectDir)

	if err := svc.InitProject(ctx, projectDir, "enterprise-app", "github.com/enterprise/app", "hexagonal", "postgres", "chi", false); err != nil {
		t.Fatalf("InitProject failed: %v", err)
	}

	// 1. AddSQLC
	if err := svc.AddSQLC(ctx, projectDir, false, false); err != nil {
		t.Fatalf("AddSQLC failed: %v", err)
	}
	exists, _ := w.Exists("/projects/enterprise-app/sqlc.yaml")
	if !exists {
		t.Error("expected sqlc.yaml to exist")
	}
	exists, _ = w.Exists("/projects/enterprise-app/db/schema/001_initial_schema.sql")
	if !exists {
		t.Error("expected db/schema/001_initial_schema.sql to exist")
	}
	exists, _ = w.Exists("/projects/enterprise-app/db/queries/accounts.sql")
	if !exists {
		t.Error("expected db/queries/accounts.sql to exist")
	}

	// 2. AddGRPCStream
	if err := svc.AddGRPCStream(ctx, projectDir, false, false); err != nil {
		t.Fatalf("AddGRPCStream failed: %v", err)
	}
	exists, _ = w.Exists("/projects/enterprise-app/api/proto/stream/v1/stream.proto")
	if !exists {
		t.Error("expected api/proto/stream/v1/stream.proto to exist")
	}
	exists, _ = w.Exists("/projects/enterprise-app/internal/adapter/handler/grpc/stream_handler.go")
	if !exists {
		t.Error("expected internal/adapter/handler/grpc/stream_handler.go to exist")
	}

	// 3. AddGRPCGateway
	if err := svc.AddGRPCGateway(ctx, projectDir, false, false); err != nil {
		t.Fatalf("AddGRPCGateway failed: %v", err)
	}
	exists, _ = w.Exists("/projects/enterprise-app/internal/adapter/handler/http/grpc_gateway.go")
	if !exists {
		t.Error("expected internal/adapter/handler/http/grpc_gateway.go to exist")
	}

	// 4. AddTenantContext
	if err := svc.AddTenantContext(ctx, projectDir, false, false); err != nil {
		t.Fatalf("AddTenantContext failed: %v", err)
	}
	exists, _ = w.Exists("/projects/enterprise-app/pkg/tenant/context.go")
	if !exists {
		t.Error("expected pkg/tenant/context.go to exist")
	}
}

func TestAetherScaffoldService_Phase15_Concurrency_ZeroTrust(t *testing.T) {
	memFS := afero.NewMemMapFs()
	w := writer.NewAferoWriter(memFS)
	resolver := manifest.NewYamlResolver(w)

	mockEmbed := fstest.MapFS{
		"common/aether_yaml.tmpl":       &fstest.MapFile{Data: []byte("version: \"1\"\narchitecture:\n  paths:\n    domain: internal/core/domain\n    port: internal/core/port\n    service: internal/core/service\n    handler_http: internal/adapter/handler/http\n    repository: internal/adapter/repository\n    cmd: cmd/server\n    pkg: pkg")},
		"common/main.go.tmpl":           &fstest.MapFile{Data: []byte("package main")},
		"common/config_viper.go.tmpl":   &fstest.MapFile{Data: []byte("package config")},
		"common/postgres.go.tmpl":       &fstest.MapFile{Data: []byte("package database")},
		"common/Makefile.tmpl":          &fstest.MapFile{Data: []byte("build:")},
		"common/Dockerfile.tmpl":        &fstest.MapFile{Data: []byte("FROM alpine")},
		"common/dockerignore.tmpl":      &fstest.MapFile{Data: []byte("bin/")},
		"common/env_example.tmpl":       &fstest.MapFile{Data: []byte("APP_ENV=")},
		"plugins/pipeline.go.tmpl":      &fstest.MapFile{Data: []byte("package concurrency\n// pipeline")},
		"plugins/singleflight.go.tmpl":  &fstest.MapFile{Data: []byte("package concurrency\n// singleflight")},
		"plugins/metrics.go.tmpl":       &fstest.MapFile{Data: []byte("package middleware\n// metrics")},
		"plugins/drain.go.tmpl":         &fstest.MapFile{Data: []byte("package server\n// drain")},
		"plugins/oauth2.go.tmpl":        &fstest.MapFile{Data: []byte("package http\n// oauth2")},
		"plugins/auditlog.go.tmpl":      &fstest.MapFile{Data: []byte("package middleware\n// auditlog")},
		"plugins/argon2.go.tmpl":        &fstest.MapFile{Data: []byte("package security\n// argon2")},
		"plugins/specification.go.tmpl": &fstest.MapFile{Data: []byte("package domain\n// specification")},
	}

	engine := template.NewStdEngine(mockEmbed)
	svc := service.NewAetherScaffoldService(engine, resolver, w)
	ctx := context.Background()
	projectDir := "/projects/final-app"
	_ = w.MkdirAll(projectDir)

	if err := svc.InitProject(ctx, projectDir, "final-app", "github.com/final/app", "hexagonal", "postgres", "chi", false); err != nil {
		t.Fatalf("InitProject failed: %v", err)
	}

	tests := []struct {
		name     string
		action   func() error
		expected string
	}{
		{"MakePipeline", func() error { return svc.MakePipeline(ctx, projectDir, "batch", false, false) }, "/projects/final-app/pkg/concurrency/batch_pipeline.go"},
		{"MakeSpecification", func() error { return svc.MakeSpecification(ctx, projectDir, "active_user", false, false) }, "/projects/final-app/internal/core/domain/active_user_specification.go"},
		{"AddSingleflight", func() error { return svc.AddSingleflight(ctx, projectDir, false, false) }, "/projects/final-app/pkg/concurrency/singleflight.go"},
		{"AddMetrics", func() error { return svc.AddMetrics(ctx, projectDir, "prometheus", false, false) }, "/projects/final-app/pkg/middleware/prometheus.go"},
		{"AddDrain", func() error { return svc.AddDrain(ctx, projectDir, false, false) }, "/projects/final-app/pkg/server/drain.go"},
		{"AddOAuth2", func() error { return svc.AddOAuth2(ctx, projectDir, "google", false, false) }, "/projects/final-app/internal/adapter/handler/http/oauth2.go"},
		{"AddAuditLog", func() error { return svc.AddAuditLog(ctx, projectDir, false, false) }, "/projects/final-app/pkg/middleware/auditlog.go"},
		{"AddArgon2", func() error { return svc.AddArgon2(ctx, projectDir, false, false) }, "/projects/final-app/pkg/security/argon2.go"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.action(); err != nil {
				t.Fatalf("%s failed: %v", tt.name, err)
			}
			exists, _ := w.Exists(tt.expected)
			if !exists {
				t.Errorf("expected %s to exist", tt.expected)
			}
		})
	}
}

func TestAetherScaffoldService_Phase16_FinalGap(t *testing.T) {
	memFS := afero.NewMemMapFs()
	w := writer.NewAferoWriter(memFS)
	resolver := manifest.NewYamlResolver(w)

	mockEmbed := fstest.MapFS{
		"common/aether_yaml.tmpl":          &fstest.MapFile{Data: []byte("version: \"1\"\narchitecture:\n  paths:\n    domain: internal/core/domain\n    port: internal/core/port\n    service: internal/core/service\n    handler_http: internal/adapter/handler/http\n    repository: internal/adapter/repository\n    cmd: cmd/server\n    pkg: pkg")},
		"common/main.go.tmpl":              &fstest.MapFile{Data: []byte("package main")},
		"common/config_viper.go.tmpl":      &fstest.MapFile{Data: []byte("package config")},
		"common/postgres.go.tmpl":          &fstest.MapFile{Data: []byte("package database")},
		"common/Makefile.tmpl":             &fstest.MapFile{Data: []byte("build:")},
		"common/Dockerfile.tmpl":           &fstest.MapFile{Data: []byte("FROM alpine")},
		"common/dockerignore.tmpl":         &fstest.MapFile{Data: []byte("bin/")},
		"common/env_example.tmpl":          &fstest.MapFile{Data: []byte("APP_ENV=")},
		"plugins/lint.yml.tmpl":            &fstest.MapFile{Data: []byte("lint")},
		"plugins/uow.go.tmpl":              &fstest.MapFile{Data: []byte("package uow")},
		"plugins/graphql.go.tmpl":          &fstest.MapFile{Data: []byte("package graphql")},
		"plugins/readreplica.go.tmpl":      &fstest.MapFile{Data: []byte("package database")},
		"plugins/openapi.go.tmpl":          &fstest.MapFile{Data: []byte("package http")},
		"plugins/cursor_paginator.go.tmpl": &fstest.MapFile{Data: []byte("package pagination")},
	}

	engine := template.NewStdEngine(mockEmbed)
	svc := service.NewAetherScaffoldService(engine, resolver, w)
	ctx := context.Background()
	projectDir := "/projects/final-gap-app"
	_ = w.MkdirAll(projectDir)

	if err := svc.InitProject(ctx, projectDir, "final-gap-app", "github.com/final/gap/app", "hexagonal", "postgres", "chi", false); err != nil {
		t.Fatalf("InitProject failed: %v", err)
	}

	tests := []struct {
		name     string
		action   func() error
		expected string
	}{
		{"AddLint", func() error { return svc.AddLint(ctx, projectDir, false, false) }, "/projects/final-gap-app/.golangci.yml"},
		{"AddUOW", func() error { return svc.AddUOW(ctx, projectDir, false, false) }, "/projects/final-gap-app/internal/core/port/uow.go"},
		{"AddGraphQL", func() error { return svc.AddGraphQL(ctx, projectDir, false, false) }, "/projects/final-gap-app/internal/adapter/handler/graphql/server.go"},
		{"AddReadReplica", func() error { return svc.AddReadReplica(ctx, projectDir, false, false) }, "/projects/final-gap-app/pkg/database/replica.go"},
		{"AddOpenAPI", func() error { return svc.AddOpenAPI(ctx, projectDir, false, false) }, "/projects/final-gap-app/internal/adapter/handler/http/swagger.go"},
		{"MakeCursorPaginator", func() error { return svc.MakeCursorPaginator(ctx, projectDir, false, false) }, "/projects/final-gap-app/pkg/pagination/cursor.go"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.action(); err != nil {
				t.Fatalf("%s failed: %v", tt.name, err)
			}
			exists, _ := w.Exists(tt.expected)
			if !exists {
				t.Errorf("expected %s to exist", tt.expected)
			}
		})
	}
}
