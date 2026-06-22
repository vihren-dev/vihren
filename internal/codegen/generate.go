package codegen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	defaultGeneratedFileName = "vihren.gen.go"
	defaultManifestFileName  = "vihren.manifest.json"
)

// GenerateRequest configures one generator invocation across package patterns.
type GenerateRequest struct {
	ModuleRoot        string
	Patterns          []string
	GeneratedFileName string
	ManifestFileName  string
}

// GenerateResult contains rendered package files and the root manifest.
type GenerateResult struct {
	Files    []GeneratedFile
	Manifest Manifest
}

// GeneratedFile is one package-local generated Go artifact.
type GeneratedFile struct {
	ImportPath string
	PackageDir string
	Path       string
	Source     []byte
}

// Manifest is the checked-in identity record for one generator invocation.
type Manifest struct {
	Entries []ManifestEntry `json:"entries"`
}

// Generate discovers markers, validates package-wide identities, and renders
// generated files without writing them to disk.
func Generate(request GenerateRequest) (GenerateResult, []Diagnostic, error) {
	request = normalizeGenerateRequest(request)
	packages, diagnostics, err := Discover(DiscoverConfig{Dir: request.ModuleRoot}, request.Patterns...)
	if err != nil {
		return GenerateResult{}, nil, err
	}
	entries := BuildManifest(packages)
	diagnostics = append(diagnostics, ValidateManifestNames(entries)...)
	if len(diagnostics) > 0 {
		return GenerateResult{}, diagnostics, nil
	}
	files := make([]GeneratedFile, 0, len(packages))
	for _, pkg := range packages {
		source, err := RenderPackage(pkg)
		if err != nil {
			return GenerateResult{}, nil, err
		}
		files = append(files, GeneratedFile{
			ImportPath: pkg.ImportPath,
			PackageDir: pkg.Dir,
			Path:       filepath.Join(pkg.Dir, request.GeneratedFileName),
			Source:     source,
		})
	}
	sort.Slice(files, func(i int, j int) bool {
		return files[i].ImportPath < files[j].ImportPath
	})
	return GenerateResult{
		Files:    files,
		Manifest: Manifest{Entries: entries},
	}, nil, nil
}

// WriteGenerated writes generated package files and the invocation manifest.
func WriteGenerated(request GenerateRequest, result GenerateResult) error {
	request = normalizeGenerateRequest(request)
	for _, file := range result.Files {
		if err := writeGeneratedFile(file.Path, file.Source); err != nil {
			return err
		}
	}
	manifest, err := RenderManifest(result.Manifest)
	if err != nil {
		return err
	}
	manifestPath := filepath.Join(request.ModuleRoot, request.ManifestFileName)
	if err := writeGeneratedFile(manifestPath, manifest); err != nil {
		return err
	}
	return nil
}

// writeGeneratedFile creates parent directories before writing generated
// artifacts so first-time go generate invocations work from a clean checkout.
func writeGeneratedFile(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create parent for %s: %w", path, err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

// RenderManifest renders the generated manifest with stable formatting.
func RenderManifest(manifest Manifest) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(manifest); err != nil {
		return nil, fmt.Errorf("render manifest: %w", err)
	}
	return buffer.Bytes(), nil
}

// FormatDiagnostics renders diagnostics for command-line output.
func FormatDiagnostics(diagnostics []Diagnostic) string {
	lines := make([]string, 0, len(diagnostics))
	for _, diagnostic := range diagnostics {
		if diagnostic.Position == "" {
			lines = append(lines, diagnostic.Message)
			continue
		}
		lines = append(lines, diagnostic.Position+": "+diagnostic.Message)
	}
	return strings.Join(lines, "\n")
}

func normalizeGenerateRequest(request GenerateRequest) GenerateRequest {
	if request.GeneratedFileName == "" {
		request.GeneratedFileName = defaultGeneratedFileName
	}
	if request.ManifestFileName == "" {
		request.ManifestFileName = defaultManifestFileName
	}
	return request
}
