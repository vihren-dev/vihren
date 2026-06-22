package codegen

import (
	"fmt"
	"go/types"
)

// SignatureShape records the generated call shape implied by a function.
type SignatureShape struct {
	InputCount int
	HasOutput  bool
}

// validateSignature applies the ADR's activity/workflow signature rules.
func validateSignature(kind MarkerKind, signature *types.Signature, position string) (SignatureShape, []Diagnostic) {
	var diagnostics []Diagnostic
	results := signature.Results()
	if signature.Variadic() {
		diagnostics = append(diagnostics, Diagnostic{Position: position, Message: fmt.Sprintf("%s functions must not be variadic", kind)})
	}
	shape, shapeDiagnostics := validateParameterShape(kind, signature, position)
	diagnostics = append(diagnostics, shapeDiagnostics...)
	if kind == WorkflowMarker && signature.Recv() != nil {
		diagnostics = append(diagnostics, Diagnostic{Position: position, Message: "workflow marker requires a function, not a method"})
	}
	switch results.Len() {
	case 1:
		if !isErrorType(results.At(0).Type()) {
			diagnostics = append(diagnostics, Diagnostic{Position: position, Message: fmt.Sprintf("%s single return must be error", kind)})
		}
	case 2:
		shape.HasOutput = true
		diagnostics = append(diagnostics, validateEndpointType(position, "output", results.At(0).Type())...)
		if !isErrorType(results.At(1).Type()) {
			diagnostics = append(diagnostics, Diagnostic{Position: position, Message: fmt.Sprintf("%s second return must be error", kind)})
		}
	default:
		diagnostics = append(diagnostics, Diagnostic{Position: position, Message: fmt.Sprintf("%s returns must be (Out, error) or (error)", kind)})
	}
	return shape, diagnostics
}

// validateParameterShape validates the activity/workflow parameter list.
func validateParameterShape(kind MarkerKind, signature *types.Signature, position string) (SignatureShape, []Diagnostic) {
	params := signature.Params()
	switch kind {
	case ActivityMarker:
		return validateActivityParameters(params, position)
	case WorkflowMarker:
		return validateWorkflowParameters(params, position)
	default:
		return SignatureShape{}, []Diagnostic{{Position: position, Message: fmt.Sprintf("unknown marker kind %q", kind)}}
	}
}

// validateActivityParameters accepts every non-dynamic activity shape the SDK
// supports: optional context.Context followed by zero or more business args.
func validateActivityParameters(params *types.Tuple, position string) (SignatureShape, []Diagnostic) {
	var diagnostics []Diagnostic
	start := 0
	if params.Len() > 0 && isActivityContext(params.At(0).Type()) {
		start = 1
	}
	for index := 0; index < params.Len(); index++ {
		if isWorkflowContext(params.At(index).Type()) {
			diagnostics = append(diagnostics, Diagnostic{Position: position, Message: "activity must not use workflow.Context"})
		}
	}
	for index := start; index < params.Len(); index++ {
		label := fmt.Sprintf("input %d", index-start+1)
		diagnostics = append(diagnostics, validateEndpointType(position, label, params.At(index).Type())...)
	}
	return SignatureShape{InputCount: params.Len() - start}, diagnostics
}

// validateWorkflowParameters keeps the ADR's workflow shape intentionally
// narrower than the SDK: workflow.Context plus zero or one business input.
func validateWorkflowParameters(params *types.Tuple, position string) (SignatureShape, []Diagnostic) {
	var diagnostics []Diagnostic
	if params.Len() == 0 || !isWorkflowContext(params.At(0).Type()) {
		diagnostics = append(diagnostics, Diagnostic{Position: position, Message: "workflow first parameter has the wrong context type"})
	}
	if params.Len() > 2 {
		diagnostics = append(diagnostics, Diagnostic{Position: position, Message: "workflow accepts at most one business input"})
	}
	if params.Len() == 2 {
		diagnostics = append(diagnostics, validateEndpointType(position, "input", params.At(1).Type())...)
	}
	return SignatureShape{InputCount: max(params.Len()-1, 0)}, diagnostics
}

// isActivityContext reports whether typ is context.Context.
func isActivityContext(typ types.Type) bool {
	return typeText(typ) == "context.Context"
}

// isWorkflowContext reports whether typ is workflow.Context.
func isWorkflowContext(typ types.Type) bool {
	text := typeText(typ)
	return text == "go.temporal.io/sdk/workflow.Context" ||
		text == "go.temporal.io/sdk/internal.Context"
}

// isErrorType reports whether typ is the builtin error interface.
func isErrorType(typ types.Type) bool {
	return typeText(typ) == "error"
}

// validateEndpointType applies the generated serializable subset.
func validateEndpointType(position string, label string, typ types.Type) []Diagnostic {
	var problems []string
	validateSerializableType(typ, map[types.Type]bool{}, &problems)
	diagnostics := make([]Diagnostic, 0, len(problems))
	for _, problem := range problems {
		diagnostics = append(diagnostics, Diagnostic{Position: position, Message: fmt.Sprintf("%s type %s", label, problem)})
	}
	return diagnostics
}

// validateSerializableType recursively rejects shapes outside the conservative subset.
func validateSerializableType(typ types.Type, seen map[types.Type]bool, problems *[]string) {
	typ = types.Unalias(typ)
	if seen[typ] {
		return
	}
	seen[typ] = true
	switch concrete := typ.(type) {
	case *types.Basic:
		if concrete.Kind() == types.UnsafePointer {
			*problems = append(*problems, "contains unsafe.Pointer")
		}
	case *types.Named:
		validateSerializableType(concrete.Underlying(), seen, problems)
	case *types.Pointer:
		validateSerializableType(concrete.Elem(), seen, problems)
	case *types.Slice:
		validateSerializableType(concrete.Elem(), seen, problems)
	case *types.Array:
		validateSerializableType(concrete.Elem(), seen, problems)
	case *types.Map:
		if basic, ok := types.Unalias(concrete.Key()).(*types.Basic); !ok || basic.Kind() != types.String {
			*problems = append(*problems, "contains a non-string map key")
		}
		validateSerializableType(concrete.Elem(), seen, problems)
	case *types.Chan:
		*problems = append(*problems, "contains "+typeText(typ))
	case *types.Signature:
		*problems = append(*problems, "contains "+typeText(typ))
	case *types.Struct:
		validateStructType(concrete, seen, problems)
	default:
		*problems = append(*problems, "contains "+typeText(typ))
	}
}

// validateStructType validates field types; unexported fields are left to the
// Temporal data converter, matching the SDK's permissive registration behavior.
func validateStructType(structType *types.Struct, seen map[types.Type]bool, problems *[]string) {
	for index := 0; index < structType.NumFields(); index++ {
		field := structType.Field(index)
		validateSerializableType(field.Type(), seen, problems)
	}
}

// typeText renders a type with package paths, which keeps aliases unambiguous.
func typeText(typ types.Type) string {
	return types.TypeString(typ, func(pkg *types.Package) string {
		return pkg.Path()
	})
}
