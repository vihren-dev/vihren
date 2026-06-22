package toolschema

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestRenderGoFileSortsAndFormatsSchemas(t *testing.T) {
	t.Parallel()
	data, err := RenderGoFile("fixturepkg", map[string]json.RawMessage{
		"ZSchema": json.RawMessage(`{"type":"string"}`),
		"ASchema": json.RawMessage(`{"type":"object"}`),
	})
	if err != nil {
		t.Fatalf("render should succeed: %v", err)
	}
	output := string(data)
	if !strings.Contains(output, "package fixturepkg") {
		t.Fatalf("output missing package: %s", output)
	}
	if strings.Index(output, "ASchema") > strings.Index(output, "ZSchema") {
		t.Fatalf("schemas were not sorted: %s", output)
	}
}

func TestRenderGoFileRejectsInvalidSchemaJSON(t *testing.T) {
	t.Parallel()
	_, err := RenderGoFile("fixturepkg", map[string]json.RawMessage{
		"BrokenSchema": json.RawMessage(`{`),
	})
	if err == nil {
		t.Fatal("expected invalid JSON error")
	}
}

func TestRenderGeneratedGoFileAssignsToolSpecSchemas(t *testing.T) {
	t.Parallel()
	data, err := RenderGeneratedGoFile(
		"fixturepkg",
		map[string]json.RawMessage{
			"ExampleToolInputSchemaJSON":  json.RawMessage(`{"type":"object"}`),
			"ExampleToolOutputSchemaJSON": json.RawMessage(`{"type":"object"}`),
		},
		[]DiscoveredToolSpec{{
			SpecVariableName:         "ExampleToolSpec",
			AccessorName:             "ExampleTool",
			SingletonName:            "exampleTool",
			InputSchemaVariableName:  "ExampleToolInputSchemaJSON",
			OutputSchemaVariableName: "ExampleToolOutputSchemaJSON",
		}},
	)
	if err != nil {
		t.Fatalf("render should succeed: %v", err)
	}
	output := string(data)
	for _, expected := range []string{
		"func init()",
		"ExampleToolSpec.InputSchema = json.RawMessage",
		"ExampleToolSpec.OutputSchema = json.RawMessage",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("generated output missing %s: %s", expected, output)
		}
	}
	for _, removed := range []string{
		"var ExampleToolInputSchemaJSON",
		"var exampleTool",
		"func ExampleTool()",
		`"github.com/vihren-dev/vihren/platform/toolcontract"`,
	} {
		if strings.Contains(output, removed) {
			t.Fatalf("generated output still contains %s: %s", removed, output)
		}
	}
}
