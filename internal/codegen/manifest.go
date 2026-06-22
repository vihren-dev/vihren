package codegen

import (
	"fmt"
	"sort"
)

// ManifestEntry is the durable identity record for one generated activity or
// workflow.
type ManifestEntry struct {
	Kind         MarkerKind `json:"kind"`
	ImportPath   string     `json:"importPath"`
	FunctionName string     `json:"functionName"`
	ProxyName    string     `json:"proxyName,omitempty"`
	TemporalName string     `json:"temporalName"`
	ExplicitName bool       `json:"explicitName"`
	Position     string     `json:"position"`
}

// BuildManifest flattens discovered packages into stable identity records.
func BuildManifest(packages []DiscoveredPackage) []ManifestEntry {
	var entries []ManifestEntry
	for _, pkg := range packages {
		for _, marker := range pkg.Markers {
			entries = append(entries, ManifestEntry{
				Kind:         marker.Kind,
				ImportPath:   marker.ImportPath,
				FunctionName: marker.FunctionName,
				ProxyName:    marker.ProxyName,
				TemporalName: marker.TemporalName,
				ExplicitName: marker.ExplicitName,
				Position:     marker.Position,
			})
		}
	}
	sort.Slice(entries, func(i int, j int) bool {
		if entries[i].TemporalName == entries[j].TemporalName {
			return entries[i].Position < entries[j].Position
		}
		return entries[i].TemporalName < entries[j].TemporalName
	})
	return entries
}

// ValidateManifestNames reports worker-wide Temporal type-name collisions.
func ValidateManifestNames(entries []ManifestEntry) []Diagnostic {
	var diagnostics []Diagnostic
	seen := map[string]ManifestEntry{}
	for _, entry := range entries {
		existing, ok := seen[entry.TemporalName]
		if !ok {
			seen[entry.TemporalName] = entry
			continue
		}
		diagnostics = append(diagnostics, Diagnostic{
			Position: entry.Position,
			Message: fmt.Sprintf(
				"Temporal name %q collides with %s at %s",
				entry.TemporalName,
				existing.FunctionName,
				existing.Position,
			),
		})
	}
	return diagnostics
}
