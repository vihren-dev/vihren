package toolcontract

import (
	"encoding/json"
	"testing"
)

func TestToolSpecToolCarriesGeneratedSchemas(t *testing.T) {
	t.Parallel()

	inputSchema := json.RawMessage(`{"type":"object","properties":{"value":{"type":"string"}}}`)
	outputSchema := json.RawMessage(`{"type":"object","properties":{"ok":{"type":"boolean"}}}`)
	spec := ToolSpec[struct{}, struct{}]{
		Name:         "test_tool",
		Description:  "Test tool.",
		InputSchema:  inputSchema,
		OutputSchema: outputSchema,
	}

	tool := spec.Tool()
	if tool.Name() != spec.Name || tool.Description() != spec.Description {
		t.Fatalf("tool metadata mismatch: %#v", tool)
	}
	schemaTool, ok := tool.(SchemaBackedTool)
	if !ok {
		t.Fatalf("ToolSpec.Tool returned non-schema-backed tool: %T", tool)
	}
	if string(schemaTool.InputSchemaJSON()) != string(inputSchema) {
		t.Fatalf("input schema = %s, want %s", schemaTool.InputSchemaJSON(), inputSchema)
	}
	if string(schemaTool.OutputSchemaJSON()) != string(outputSchema) {
		t.Fatalf("output schema = %s, want %s", schemaTool.OutputSchemaJSON(), outputSchema)
	}
}
