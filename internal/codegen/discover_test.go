package codegen

import (
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"
)

// TestDiscoverValidMarkers proves go/packages exposes the marker, type,
// receiver, import-path, and source-position data the generator depends on.
func TestDiscoverValidMarkers(t *testing.T) {
	t.Parallel()
	packages, diagnostics, err := Discover(
		DiscoverConfig{Dir: moduleRoot(t)},
		"./internal/codegen/testdata/fixtures/basic",
	)
	if err != nil {
		t.Fatalf("discover packages: %v", err)
	}
	if len(diagnostics) != 0 {
		t.Fatalf("diagnostics = %#v, want none", diagnostics)
	}
	if len(packages) != 1 {
		t.Fatalf("package count = %d, want 1", len(packages))
	}
	markers := packages[0].Markers
	if len(markers) != 5 {
		t.Fatalf("marker count = %d, want 5: %#v", len(markers), markers)
	}
	charge := markerByFunction(t, markers, "ChargeCard")
	if charge.Kind != ActivityMarker || charge.ReceiverName != "Activities" {
		t.Fatalf("charge marker = %#v, want activity on Activities", charge)
	}
	if charge.ProxyName != "Charge" {
		t.Fatalf("charge proxy = %q, want Charge", charge.ProxyName)
	}
	wantChargeName := "github.com/vihren-dev/vihren/internal/codegen/testdata/fixtures/basic.Charge"
	if charge.TemporalName != wantChargeName {
		t.Fatalf("charge temporal name = %q, want %q", charge.TemporalName, wantChargeName)
	}
	normalize := markerByFunction(t, markers, "Normalize")
	if normalize.TemporalName != "external.normalize" {
		t.Fatalf("normalize temporal name = %q", normalize.TemporalName)
	}
	price := markerByFunction(t, markers, "Price")
	if price.InputCount != 2 {
		t.Fatalf("price input count = %d, want 2", price.InputCount)
	}
	ping := markerByFunction(t, markers, "Ping")
	if ping.InputCount != 0 || ping.HasOutput {
		t.Fatalf("ping marker = %#v, want zero inputs and no output", ping)
	}
	checkout := markerByFunction(t, markers, "Checkout")
	if checkout.Kind != WorkflowMarker || checkout.ReceiverName != "" {
		t.Fatalf("checkout marker = %#v, want workflow function", checkout)
	}
	if !strings.Contains(checkout.Position, "basic.go") {
		t.Fatalf("checkout position = %q, want source file", checkout.Position)
	}
}

// TestDiscoverInvalidMarkers proves bad signatures, duplicate proxies, and
// marker options become diagnostics rather than loader failures.
func TestDiscoverInvalidMarkers(t *testing.T) {
	t.Parallel()
	_, diagnostics, err := Discover(
		DiscoverConfig{Dir: moduleRoot(t)},
		"./internal/codegen/testdata/fixtures/invalid",
	)
	if err != nil {
		t.Fatalf("discover packages: %v", err)
	}
	messages := diagnosticMessages(diagnostics)
	for _, want := range []string{
		"activity must not use workflow.Context",
		"input 1 type contains func()",
		"input 1 type contains chan int",
		"input 1 type contains unsafe.Pointer",
		"input 1 type contains a non-string map key",
		"activity proxy name \"Duplicate\" collides",
		"workflow first parameter has the wrong context type",
		"unknown option \"timeout\"",
	} {
		if !slices.ContainsFunc(messages, func(message string) bool { return strings.Contains(message, want) }) {
			t.Fatalf("diagnostics %q do not contain %q", strings.Join(messages, "\n"), want)
		}
	}
}

// markerByFunction finds one marker by Go function name.
func markerByFunction(t *testing.T, markers []DiscoveredMarker, name string) DiscoveredMarker {
	t.Helper()
	for _, marker := range markers {
		if marker.FunctionName == name {
			return marker
		}
	}
	t.Fatalf("marker %q not found in %#v", name, markers)
	return DiscoveredMarker{}
}

// diagnosticMessages returns diagnostics without positions for substring checks.
func diagnosticMessages(diagnostics []Diagnostic) []string {
	messages := make([]string, 0, len(diagnostics))
	for _, diagnostic := range diagnostics {
		messages = append(messages, diagnostic.Message)
	}
	return messages
}

// moduleRoot returns the repository root for go/packages patterns.
func moduleRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime caller unavailable")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "../.."))
}
