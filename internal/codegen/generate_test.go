package codegen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestGenerateRendersFilesAndManifest proves the orchestration layer returns
// package-local generated files and a stable invocation manifest.
func TestGenerateRendersFilesAndManifest(t *testing.T) {
	t.Parallel()
	result, diagnostics, err := Generate(GenerateRequest{
		ModuleRoot: moduleRoot(t),
		Patterns:   []string{"./internal/codegen/testdata/fixtures/basic"},
	})
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if len(diagnostics) != 0 {
		t.Fatalf("diagnostics = %#v, want none", diagnostics)
	}
	if len(result.Files) != 1 {
		t.Fatalf("generated file count = %d, want 1", len(result.Files))
	}
	if !strings.HasSuffix(result.Files[0].Path, "internal/codegen/testdata/fixtures/basic/vihren.gen.go") {
		t.Fatalf("generated path = %q", result.Files[0].Path)
	}
	if !strings.Contains(string(result.Files[0].Source), "func Register(r worker.Registry, activities *Activities)") {
		t.Fatalf("generated source missing Register:\n%s", result.Files[0].Source)
	}
	if len(result.Manifest.Entries) != 5 {
		t.Fatalf("manifest entries = %#v, want 5", result.Manifest.Entries)
	}
	manifest, err := RenderManifest(result.Manifest)
	if err != nil {
		t.Fatalf("render manifest: %v", err)
	}
	if !strings.Contains(string(manifest), `"temporalName": "external.normalize"`) {
		t.Fatalf("manifest missing explicit name:\n%s", manifest)
	}
}

// TestGenerateStopsOnDiagnostics proves generation does not render files when
// package diagnostics need user attention first.
func TestGenerateStopsOnDiagnostics(t *testing.T) {
	t.Parallel()
	result, diagnostics, err := Generate(GenerateRequest{
		ModuleRoot: moduleRoot(t),
		Patterns:   []string{"./internal/codegen/testdata/fixtures/invalid"},
	})
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if len(diagnostics) == 0 {
		t.Fatal("diagnostics = none, want invalid fixture diagnostics")
	}
	if len(result.Files) != 0 {
		t.Fatalf("generated files = %#v, want none", result.Files)
	}
	if !strings.Contains(FormatDiagnostics(diagnostics), "activity proxy name \"Duplicate\" collides") {
		t.Fatalf("formatted diagnostics missing duplicate proxy:\n%s", FormatDiagnostics(diagnostics))
	}
}

// TestWriteGeneratedWritesFilesAndManifest proves the disk-writing boundary
// persists both package-local generated files and the root manifest.
func TestWriteGeneratedWritesFilesAndManifest(t *testing.T) {
	t.Parallel()
	tempDir := t.TempDir()
	packageDir := filepath.Join(tempDir, "fixture")
	if err := os.Mkdir(packageDir, 0o755); err != nil {
		t.Fatalf("make package dir: %v", err)
	}
	result := GenerateResult{
		Files: []GeneratedFile{{
			ImportPath: "example.com/fixture",
			PackageDir: packageDir,
			Path:       filepath.Join(packageDir, "vihren.gen.go"),
			Source:     []byte("package fixture\n"),
		}},
		Manifest: Manifest{Entries: []ManifestEntry{{
			Kind:         ActivityMarker,
			ImportPath:   "example.com/fixture",
			FunctionName: "Touch",
			ProxyName:    "Touch",
			TemporalName: "example.com/fixture.Touch",
			Position:     "fixture.go:1:1",
		}}},
	}
	if err := WriteGenerated(GenerateRequest{ModuleRoot: tempDir}, result); err != nil {
		t.Fatalf("write generated: %v", err)
	}
	source, err := os.ReadFile(filepath.Join(packageDir, "vihren.gen.go"))
	if err != nil {
		t.Fatalf("read generated file: %v", err)
	}
	if string(source) != "package fixture\n" {
		t.Fatalf("generated source = %q", source)
	}
	manifest, err := os.ReadFile(filepath.Join(tempDir, "vihren.manifest.json"))
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	if !strings.Contains(string(manifest), `"temporalName": "example.com/fixture.Touch"`) {
		t.Fatalf("manifest missing entry:\n%s", manifest)
	}
}

// TestWriteGeneratedCreatesManifestParentDirectory proves example go:generate
// directives can write manifests under .cache in a fresh checkout.
func TestWriteGeneratedCreatesManifestParentDirectory(t *testing.T) {
	t.Parallel()
	tempDir := t.TempDir()
	result := GenerateResult{
		Manifest: Manifest{Entries: []ManifestEntry{{
			Kind:         ActivityMarker,
			ImportPath:   "example.com/fixture",
			FunctionName: "Touch",
			ProxyName:    "Touch",
			TemporalName: "example.com/fixture.Touch",
			Position:     "fixture.go:1:1",
		}}},
	}
	if err := WriteGenerated(GenerateRequest{
		ModuleRoot:       tempDir,
		ManifestFileName: filepath.Join(".cache", "vihren.manifest.json"),
	}, result); err != nil {
		t.Fatalf("write generated: %v", err)
	}
	manifest, err := os.ReadFile(filepath.Join(tempDir, ".cache", "vihren.manifest.json"))
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	if !strings.Contains(string(manifest), `"temporalName": "example.com/fixture.Touch"`) {
		t.Fatalf("manifest missing entry:\n%s", manifest)
	}
}

// TestGenerateAllowsActivityProxyBootstrap proves first-time generation works
// even when workflow source already calls the generated Activity proxy.
func TestGenerateAllowsActivityProxyBootstrap(t *testing.T) {
	t.Parallel()
	result, diagnostics, err := Generate(GenerateRequest{
		ModuleRoot: moduleRoot(t),
		Patterns:   []string{"./internal/codegen/testdata/fixtures/bootstrap"},
	})
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if len(diagnostics) != 0 {
		t.Fatalf("diagnostics = %#v, want none", diagnostics)
	}
	if len(result.Files) != 1 {
		t.Fatalf("generated file count = %d, want 1", len(result.Files))
	}
	if !strings.Contains(string(result.Files[0].Source), "var Activity activityProxy") {
		t.Fatalf("generated source missing Activity proxy:\n%s", result.Files[0].Source)
	}
}
