package toolschema

import (
	"fmt"
	"strings"
)

// ParseSchemaSpec parses a CLI schema spec formatted as VariableName=TypeName.
func ParseSchemaSpec(raw string) (SchemaSpec, error) {
	left, right, ok := strings.Cut(raw, "=")
	if !ok {
		return SchemaSpec{}, fmt.Errorf("schema spec %q must be VariableName=TypeName", raw)
	}
	spec := SchemaSpec{VariableName: strings.TrimSpace(left), TypeName: strings.TrimSpace(right)}
	if spec.VariableName == "" || spec.TypeName == "" {
		return SchemaSpec{}, fmt.Errorf("schema spec %q must include variable and type names", raw)
	}
	return spec, nil
}

// FormatSchemaSpec formats one schema spec for diagnostics and tests.
func FormatSchemaSpec(spec SchemaSpec) string {
	return spec.VariableName + "=" + spec.TypeName
}
