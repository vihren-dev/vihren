package toolschema

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestGenerateGoFileReflectsFixturePackage(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	moduleRoot, err := FindModuleRoot(".")
	if err != nil {
		t.Fatalf("module root should resolve: %v", err)
	}
	data, err := GenerateGoFile(ctx, GenerateRequest{
		ModuleRoot: moduleRoot,
		Package:    "./internal/toolschema/testdata/fixturepkg",
		Schemas: []SchemaSpec{
			{VariableName: "ExampleInputSchemaJSON", TypeName: "ExampleInput"},
			{VariableName: "ExampleOutputSchemaJSON", TypeName: "ExampleOutput"},
		},
	})
	if err != nil {
		t.Fatalf("generate should succeed: %v", err)
	}
	output := string(data)
	for _, expected := range []string{
		"package fixturepkg",
		"ExampleInputSchemaJSON",
		"ExampleOutputSchemaJSON",
		`"name"`,
		`"description"`,
		`"minimum":1`,
		`"maximum":5`,
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("generated output missing %s: %s", expected, output)
		}
	}
}

func TestGenerateGoFileDiscoversToolSpecs(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	moduleRoot, err := FindModuleRoot(".")
	if err != nil {
		t.Fatalf("module root should resolve: %v", err)
	}
	data, err := GenerateGoFile(ctx, GenerateRequest{
		ModuleRoot: moduleRoot,
		Package:    "./internal/toolschema/testdata/fixturepkg",
	})
	if err != nil {
		t.Fatalf("generate should succeed: %v", err)
	}
	output := string(data)
	for _, expected := range []string{
		"func init()",
		"ExampleToolSpec.InputSchema = json.RawMessage",
		"ExampleToolSpec.OutputSchema = json.RawMessage",
		`"name"`,
		`"summary"`,
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("generated output missing %s: %s", expected, output)
		}
	}
	for _, removed := range []string{
		"func ExampleTool()",
		"var exampleTool",
		"var ExampleToolInputSchemaJSON",
	} {
		if strings.Contains(output, removed) {
			t.Fatalf("generated output still contains %s: %s", removed, output)
		}
	}
}
