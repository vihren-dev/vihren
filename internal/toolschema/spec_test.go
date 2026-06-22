package toolschema

import "testing"

func TestParseSchemaSpec(t *testing.T) {
	t.Parallel()
	spec, err := ParseSchemaSpec("ExampleInputSchemaJSON=ExampleInput")
	if err != nil {
		t.Fatalf("schema spec should parse: %v", err)
	}
	if spec.VariableName != "ExampleInputSchemaJSON" || spec.TypeName != "ExampleInput" {
		t.Fatalf("unexpected spec: %#v", spec)
	}
}

func TestParseSchemaSpecRejectsInvalidShape(t *testing.T) {
	t.Parallel()
	if _, err := ParseSchemaSpec("ExampleInput"); err == nil {
		t.Fatal("expected invalid schema spec error")
	}
}
