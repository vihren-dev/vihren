// Package toolschema generates checked-in JSON Schema Go artifacts for tool
// input and output structs.
package toolschema

// SchemaSpec maps one generated Go variable to one exported Go type.
type SchemaSpec struct {
	VariableName string
	TypeName     string
}

// DiscoveredToolSpec records one package-level tool spec selected for
// generated schema-backed tool code.
type DiscoveredToolSpec struct {
	SpecVariableName         string
	AccessorName             string
	SingletonName            string
	InputTypeName            string
	OutputTypeName           string
	InputSchemaVariableName  string
	OutputSchemaVariableName string
}

// GeneratedGoFile contains all generated artifacts for one package.
type GeneratedGoFile struct {
	Schemas []SchemaSpec
	Tools   []DiscoveredToolSpec
}

// GenerateRequest configures source-time schema generation for one package.
type GenerateRequest struct {
	ModuleRoot    string
	Package       string
	OutputPackage string
	Schemas       []SchemaSpec
}

// PackageInfo records the target package details needed by the generated
// reflection helper.
type PackageInfo struct {
	ImportPath string
	Dir        string
	Name       string
	ModulePath string
}
