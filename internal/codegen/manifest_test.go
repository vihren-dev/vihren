package codegen

import (
	"slices"
	"strings"
	"testing"
)

// TestManifestDefaultsUseImportPath proves package-name collisions do not
// collide when default Temporal names use the import path.
func TestManifestDefaultsUseImportPath(t *testing.T) {
	t.Parallel()
	packages, diagnostics, err := Discover(
		DiscoverConfig{Dir: moduleRoot(t)},
		"./internal/codegen/testdata/fixtures/names/defaultalpha",
		"./internal/codegen/testdata/fixtures/names/defaultbeta",
	)
	if err != nil {
		t.Fatalf("discover packages: %v", err)
	}
	if len(diagnostics) != 0 {
		t.Fatalf("diagnostics = %#v, want none", diagnostics)
	}
	entries := BuildManifest(packages)
	if len(entries) != 2 {
		t.Fatalf("entry count = %d, want 2: %#v", len(entries), entries)
	}
	names := []string{entries[0].TemporalName, entries[1].TemporalName}
	if names[0] == names[1] {
		t.Fatalf("default names collided: %#v", entries)
	}
	for _, name := range names {
		if !strings.Contains(name, "internal/codegen/testdata/fixtures/names/default") {
			t.Fatalf("default name %q does not contain import path", name)
		}
	}
	if nameDiagnostics := ValidateManifestNames(entries); len(nameDiagnostics) != 0 {
		t.Fatalf("name diagnostics = %#v, want none", nameDiagnostics)
	}
}

// TestManifestRejectsExplicitNameCollisions proves explicit name overrides are
// worker-wide and need generated diagnostics.
func TestManifestRejectsExplicitNameCollisions(t *testing.T) {
	t.Parallel()
	packages, diagnostics, err := Discover(
		DiscoverConfig{Dir: moduleRoot(t)},
		"./internal/codegen/testdata/fixtures/names/explicitone",
		"./internal/codegen/testdata/fixtures/names/explicittwo",
	)
	if err != nil {
		t.Fatalf("discover packages: %v", err)
	}
	if len(diagnostics) != 0 {
		t.Fatalf("diagnostics = %#v, want none", diagnostics)
	}
	entries := BuildManifest(packages)
	if !slices.ContainsFunc(entries, func(entry ManifestEntry) bool {
		return entry.ExplicitName && entry.TemporalName == "shared.operation"
	}) {
		t.Fatalf("entries do not record explicit shared.operation: %#v", entries)
	}
	nameDiagnostics := ValidateManifestNames(entries)
	if len(nameDiagnostics) != 1 {
		t.Fatalf("name diagnostics = %#v, want one", nameDiagnostics)
	}
	if !strings.Contains(nameDiagnostics[0].Message, `Temporal name "shared.operation" collides`) {
		t.Fatalf("diagnostic = %q", nameDiagnostics[0].Message)
	}
}
