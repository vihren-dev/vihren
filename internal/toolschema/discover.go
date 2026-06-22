package toolschema

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
	"unicode"
)

// DiscoverToolSpecs finds package-level toolcontract.ToolSpec[I, O] values in
// the target package.
func DiscoverToolSpecs(info PackageInfo) ([]DiscoveredToolSpec, error) {
	files, err := parsePackageFiles(info.Dir)
	if err != nil {
		return nil, err
	}
	var discovered []DiscoveredToolSpec
	seen := map[string]struct{}{}
	for _, file := range files {
		aliases := toolcontractAliases(file, info.ModulePath)
		for _, declaration := range file.Decls {
			general, ok := declaration.(*ast.GenDecl)
			if !ok || general.Tok != token.VAR {
				continue
			}
			for _, spec := range general.Specs {
				values, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}
				found, err := discoverValueSpec(values, aliases, info.Name == "toolcontract")
				if err != nil {
					return nil, err
				}
				for _, tool := range found {
					for _, generatedName := range []string{
						tool.AccessorName,
						tool.SingletonName,
						tool.InputSchemaVariableName,
						tool.OutputSchemaVariableName,
					} {
						if _, ok := seen[generatedName]; ok {
							return nil, fmt.Errorf("duplicate generated name %q", generatedName)
						}
						seen[generatedName] = struct{}{}
					}
					discovered = append(discovered, tool)
				}
			}
		}
	}
	return discovered, nil
}

func parsePackageFiles(dir string) ([]*ast.File, error) {
	matches, err := filepath.Glob(filepath.Join(dir, "*.go"))
	if err != nil {
		return nil, err
	}
	files := make([]*ast.File, 0, len(matches))
	fileSet := token.NewFileSet()
	for _, path := range matches {
		if strings.HasSuffix(path, "_test.go") || strings.HasSuffix(path, ".gen.go") {
			continue
		}
		file, err := parser.ParseFile(fileSet, path, nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	return files, nil
}

func toolcontractAliases(file *ast.File, modulePath string) map[string]struct{} {
	aliases := map[string]struct{}{}
	toolcontractPath := modulePath + "/platform/toolcontract"
	for _, importSpec := range file.Imports {
		if strings.Trim(importSpec.Path.Value, `"`) != toolcontractPath {
			continue
		}
		if importSpec.Name != nil {
			aliases[importSpec.Name.Name] = struct{}{}
			continue
		}
		aliases["toolcontract"] = struct{}{}
	}
	return aliases
}

func discoverValueSpec(
	spec *ast.ValueSpec,
	aliases map[string]struct{},
	samePackage bool,
) ([]DiscoveredToolSpec, error) {
	discovered := make([]DiscoveredToolSpec, 0, len(spec.Names))
	for index, name := range spec.Names {
		typeExpression := spec.Type
		if typeExpression == nil && index < len(spec.Values) {
			if composite, ok := spec.Values[index].(*ast.CompositeLit); ok {
				typeExpression = composite.Type
			}
		}
		if typeExpression == nil {
			continue
		}
		inputType, outputType, ok, err := toolSpecTypeArguments(typeExpression, aliases, samePackage)
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}
		tool, err := newDiscoveredToolSpec(name.Name, inputType, outputType)
		if err != nil {
			return nil, err
		}
		discovered = append(discovered, tool)
	}
	return discovered, nil
}

func toolSpecTypeArguments(
	expression ast.Expr,
	aliases map[string]struct{},
	samePackage bool,
) (string, string, bool, error) {
	indexed, ok := expression.(*ast.IndexListExpr)
	if !ok || len(indexed.Indices) != 2 {
		return "", "", false, nil
	}
	if !isToolSpecExpression(indexed.X, aliases, samePackage) {
		return "", "", false, nil
	}
	input, ok := indexed.Indices[0].(*ast.Ident)
	if !ok {
		return "", "", true, fmt.Errorf("ToolSpec input type must be a package-local named type")
	}
	output, ok := indexed.Indices[1].(*ast.Ident)
	if !ok {
		return "", "", true, fmt.Errorf("ToolSpec output type must be a package-local named type")
	}
	if !ast.IsExported(input.Name) || !ast.IsExported(output.Name) {
		return "", "", true, fmt.Errorf("ToolSpec input and output types must be exported")
	}
	return input.Name, output.Name, true, nil
}

func isToolSpecExpression(expression ast.Expr, aliases map[string]struct{}, samePackage bool) bool {
	if samePackage {
		if ident, ok := expression.(*ast.Ident); ok && ident.Name == "ToolSpec" {
			return true
		}
	}
	selector, ok := expression.(*ast.SelectorExpr)
	if !ok || selector.Sel.Name != "ToolSpec" {
		return false
	}
	qualifier, ok := selector.X.(*ast.Ident)
	if !ok {
		return false
	}
	_, ok = aliases[qualifier.Name]
	return ok
}

func newDiscoveredToolSpec(specVariableName string, inputType string, outputType string) (DiscoveredToolSpec, error) {
	baseName, ok := strings.CutSuffix(specVariableName, "Spec")
	if !ok || baseName == "" {
		return DiscoveredToolSpec{}, fmt.Errorf("ToolSpec variable %q must end with Spec", specVariableName)
	}
	return DiscoveredToolSpec{
		SpecVariableName:         specVariableName,
		AccessorName:             baseName,
		SingletonName:            lowerFirst(baseName),
		InputTypeName:            inputType,
		OutputTypeName:           outputType,
		InputSchemaVariableName:  baseName + "InputSchemaJSON",
		OutputSchemaVariableName: baseName + "OutputSchemaJSON",
	}, nil
}

func lowerFirst(value string) string {
	if value == "" {
		return ""
	}
	runes := []rune(value)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}
