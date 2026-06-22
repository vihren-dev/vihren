// Package toolcontract owns logical tool contracts and shared tool-call result
// types used by tool catalogs and provider adapters.
package toolcontract

import (
	"encoding/json"
	"strings"

	"github.com/vihren-dev/vihren/platform/blobref"
)

// ToolName is the stable identity for a logical model-callable tool.
type ToolName string

// Tool is the broad logical contract sessions use to allow tools.
type Tool interface {
	Name() ToolName
	Description() string
}

// SchemaBackedTool exposes schemas to provider rendering and metadata paths.
type SchemaBackedTool interface {
	Tool
	InputSchemaJSON() json.RawMessage
	OutputSchemaJSON() json.RawMessage
}

// SchemaDerivedTool is the default logical tool contract backed by generated
// JSON Schema artifacts.
type SchemaDerivedTool struct {
	ToolName        ToolName
	ToolDescription string
	InputSchema     json.RawMessage
	OutputSchema    json.RawMessage
}

// ToolSpec is the author-written static contract that source generation turns
// into a schema-backed logical tool.
type ToolSpec[I any, O any] struct {
	Name         ToolName
	Description  string
	InputSchema  json.RawMessage
	OutputSchema json.RawMessage
}

// Tool returns the schema-backed logical tool described by the spec. Generated
// schema files populate InputSchema and OutputSchema during package
// initialization, so callers should read this at runtime rather than from other
// package-level initializers.
func (spec ToolSpec[I, O]) Tool() Tool {
	return NewSchemaDerivedTool(spec.Name, spec.Description, spec.InputSchema, spec.OutputSchema)
}

// NewSchemaDerivedTool constructs a schema-backed logical tool.
func NewSchemaDerivedTool(
	name ToolName,
	description string,
	inputSchema json.RawMessage,
	outputSchema json.RawMessage,
) SchemaDerivedTool {
	return SchemaDerivedTool{
		ToolName:        name,
		ToolDescription: description,
		InputSchema:     inputSchema,
		OutputSchema:    outputSchema,
	}
}

// Name returns the stable logical tool identity.
func (tool SchemaDerivedTool) Name() ToolName {
	return tool.ToolName
}

// Description returns the model-facing tool description.
func (tool SchemaDerivedTool) Description() string {
	return tool.ToolDescription
}

// InputSchemaJSON returns the generated input schema for provider rendering.
func (tool SchemaDerivedTool) InputSchemaJSON() json.RawMessage {
	return tool.InputSchema
}

// OutputSchemaJSON returns the generated output schema metadata.
func (tool SchemaDerivedTool) OutputSchemaJSON() json.RawMessage {
	return tool.OutputSchema
}

// ModelToolCall is the provider-neutral model tool-call payload the catalog
// accepts after a provider adapter decodes a concrete function-call response.
type ModelToolCall struct {
	CallID        string `json:"call_id"`
	Name          string `json:"name"`
	ArgumentsJSON string `json:"arguments_json"`
}

// FailureKind classifies why a model-shaped tool call was rejected or failed.
type FailureKind string

const (
	// FailureUnknownTool means no reviewed tool registration exists.
	FailureUnknownTool FailureKind = "unknown_tool"

	// FailureDisallowedTool means the session did not allow the tool.
	FailureDisallowedTool FailureKind = "disallowed_tool"

	// FailureInvalidArguments means tool arguments did not satisfy the contract.
	FailureInvalidArguments FailureKind = "invalid_arguments"

	// FailureInvalidOutput means output could not be mechanically encoded for the
	// next model turn.
	FailureInvalidOutput FailureKind = "invalid_output"

	// FailureExecution means Temporal activity execution failed.
	FailureExecution FailureKind = "execution_failed"
)

// ToolFailure explains a rejected or failed tool call for review output.
type ToolFailure struct {
	Kind    FailureKind `json:"kind"`
	Message string      `json:"message"`
}

// ToolResult carries control data and provider continuation output for one model
// tool call.
type ToolResult struct {
	CallID  string                         `json:"call_id"`
	Output  blobref.Value[json.RawMessage] `json:"output,omitempty"`
	Failure *ToolFailure                   `json:"failure,omitempty"`
}

// ValidationProblem explains one schema or contract validation failure.
type ValidationProblem struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// FormatValidationProblems formats validation failures for tool execution
// messages.
func FormatValidationProblems(problems []ValidationProblem) string {
	messages := make([]string, 0, len(problems))
	for _, problem := range problems {
		messages = append(messages, problem.Field+": "+problem.Message)
	}
	return strings.Join(messages, "; ")
}
