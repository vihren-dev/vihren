package codegen

import (
	"fmt"
	"go/token"
	"strings"
)

// parseMarkerComment parses a single line comment if it is a vihren marker.
func parseMarkerComment(
	fileSet *token.FileSet,
	commentText string,
	commentPos token.Pos,
) (MarkerKind, MarkerOptions, bool, []Diagnostic) {
	text := strings.TrimSpace(strings.TrimPrefix(commentText, "//"))
	if !strings.HasPrefix(text, "vihren:") {
		return "", MarkerOptions{}, false, nil
	}
	fields := strings.Fields(text)
	if len(fields) == 0 {
		return "", MarkerOptions{}, false, nil
	}
	kindText := strings.TrimPrefix(fields[0], "vihren:")
	kind := MarkerKind(kindText)
	position := fileSet.Position(commentPos).String()
	if kind != ActivityMarker && kind != WorkflowMarker {
		return "", MarkerOptions{}, false, []Diagnostic{{
			Position: position,
			Message:  fmt.Sprintf("unknown vihren marker %q", fields[0]),
		}}
	}
	options, diagnostics := parseMarkerOptions(kind, fields[1:], position)
	return kind, options, true, diagnostics
}

// parseMarkerOptions parses space-separated key=value pairs and boolean flags.
func parseMarkerOptions(kind MarkerKind, fields []string, position string) (MarkerOptions, []Diagnostic) {
	var options MarkerOptions
	var diagnostics []Diagnostic
	for _, field := range fields {
		key, value, hasValue := strings.Cut(field, "=")
		switch key {
		case "name":
			options.Name = value
			if !hasValue {
				diagnostics = append(diagnostics, missingValue(position, key))
			}
		case "proxy":
			options.Proxy = value
			if kind != ActivityMarker {
				diagnostics = append(diagnostics, unknownOption(position, key))
			}
			if !hasValue {
				diagnostics = append(diagnostics, missingValue(position, key))
			}
		case "versioningBehavior":
			options.VersioningBehavior = value
			if kind != WorkflowMarker {
				diagnostics = append(diagnostics, unknownOption(position, key))
			}
			if !hasValue {
				diagnostics = append(diagnostics, missingValue(position, key))
			}
		case "disableAlreadyRegisteredCheck":
			options.DisableAlreadyRegisteredCheck = true
			if hasValue {
				diagnostics = append(diagnostics, flagHasValue(position, key))
			}
		case "skipInvalidStructFunctions":
			options.SkipInvalidStructFunctions = true
			if kind != ActivityMarker {
				diagnostics = append(diagnostics, unknownOption(position, key))
			}
			if hasValue {
				diagnostics = append(diagnostics, flagHasValue(position, key))
			}
		default:
			diagnostics = append(diagnostics, unknownOption(position, key))
		}
	}
	return options, diagnostics
}

// missingValue reports an option that requires a value.
func missingValue(position string, key string) Diagnostic {
	return Diagnostic{Position: position, Message: fmt.Sprintf("option %q requires a value", key)}
}

// flagHasValue reports a boolean flag that was written as key=value.
func flagHasValue(position string, key string) Diagnostic {
	return Diagnostic{Position: position, Message: fmt.Sprintf("option %q is a boolean flag", key)}
}

// unknownOption reports an option outside the marker grammar.
func unknownOption(position string, key string) Diagnostic {
	return Diagnostic{Position: position, Message: fmt.Sprintf("unknown option %q", key)}
}
