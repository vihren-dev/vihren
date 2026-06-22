package main

import (
	"strings"
	"testing"
)

// TestRunRequiresPackagePattern keeps the CLI from silently generating nothing.
func TestRunRequiresPackagePattern(t *testing.T) {
	t.Parallel()
	err := run(nil)
	if err == nil || !strings.Contains(err.Error(), "at least one package pattern is required") {
		t.Fatalf("run error = %v, want package pattern error", err)
	}
}

// TestRunDryRunRendersValidFixture proves the CLI can exercise discovery and
// rendering without mutating fixture packages.
func TestRunDryRunRendersValidFixture(t *testing.T) {
	t.Parallel()
	if err := run([]string{"--dry-run", "./internal/codegen/testdata/fixtures/basic"}); err != nil {
		t.Fatalf("dry run: %v", err)
	}
}

// TestRunFormatsDiagnostics proves generator diagnostics are surfaced through
// the command-line boundary.
func TestRunFormatsDiagnostics(t *testing.T) {
	t.Parallel()
	err := run([]string{"--dry-run", "./internal/codegen/testdata/fixtures/invalid"})
	if err == nil {
		t.Fatal("dry run invalid fixture error = nil, want diagnostics")
	}
	for _, want := range []string{
		"codegen diagnostics:",
		"activity proxy name \"Duplicate\" collides",
		"unknown option \"timeout\"",
	} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("error %q missing %q", err.Error(), want)
		}
	}
}

// TestRunRejectsManifestNameCollision proves worker-wide Temporal identity
// collisions are checked at the command boundary.
func TestRunRejectsManifestNameCollision(t *testing.T) {
	t.Parallel()
	err := run([]string{
		"--dry-run",
		"./internal/codegen/testdata/fixtures/names/explicitone",
		"./internal/codegen/testdata/fixtures/names/explicittwo",
	})
	if err == nil {
		t.Fatal("collision run error = nil, want diagnostics")
	}
	if !strings.Contains(err.Error(), `Temporal name "shared.operation" collides`) {
		t.Fatalf("error %q missing collision diagnostic", err.Error())
	}
}
