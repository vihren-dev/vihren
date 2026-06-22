package fixturepkg

import "github.com/vihren-dev/vihren/platform/toolcontract"

const (
	// ExampleToolName is the fixture tool identity used by generator tests.
	ExampleToolName toolcontract.ToolName = "example_tool"
)

// ExampleInput describes the model-facing input contract.
type ExampleInput struct {
	// Name selects the record to inspect.
	Name string `json:"name"`

	// Count limits the number of rows returned.
	Count int `json:"count,omitempty" jsonschema:"minimum=1,maximum=5"`
}

// ExampleOutput describes the tool result contract.
type ExampleOutput struct {
	// Summary is the compact reviewed result.
	Summary string `json:"summary"`
}

// ExampleToolSpec is the fixture package-level tool contract.
var ExampleToolSpec = toolcontract.ToolSpec[ExampleInput, ExampleOutput]{
	Name:        ExampleToolName,
	Description: "Read the fixture example input and return a summary.",
}
