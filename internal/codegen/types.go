package codegen

import "go/types"

// MarkerKind identifies which vihren marker was attached to a function.
type MarkerKind string

const (
	// ActivityMarker marks a Temporal activity implementation.
	ActivityMarker MarkerKind = "activity"

	// WorkflowMarker marks a Temporal workflow implementation.
	WorkflowMarker MarkerKind = "workflow"
)

// MarkerOptions carries parsed marker options from one comment.
type MarkerOptions struct {
	Name                          string
	Proxy                         string
	VersioningBehavior            string
	DisableAlreadyRegisteredCheck bool
	SkipInvalidStructFunctions    bool
}

// DiscoveredMarker is the normalized view of one annotated function.
type DiscoveredMarker struct {
	Kind         MarkerKind
	PackageName  string
	PackageDir   string
	ImportPath   string
	FunctionName string
	ReceiverName string
	ProxyName    string
	TemporalName string
	ExplicitName bool
	Position     string
	InputCount   int
	HasOutput    bool
	Options      MarkerOptions
	signature    *types.Signature
}

// DiscoveredPackage groups markers discovered from one loaded package.
type DiscoveredPackage struct {
	Name       string
	Dir        string
	ImportPath string
	Markers    []DiscoveredMarker
	types      *types.Package
}

// Diagnostic records a recoverable generator diagnostic with a source position.
type Diagnostic struct {
	Position string
	Message  string
}
