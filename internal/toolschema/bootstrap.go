package toolschema

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

func writeBootstrapToolAccessors(info PackageInfo, tools []DiscoveredToolSpec) (func(), error) {
	if len(tools) == 0 {
		return func() {}, nil
	}
	existing, err := existingFunctionNames(info.Dir)
	if err != nil {
		return nil, err
	}
	missing := make([]DiscoveredToolSpec, 0, len(tools))
	for _, tool := range tools {
		if _, ok := existing[tool.AccessorName]; !ok {
			missing = append(missing, tool)
		}
	}
	if len(missing) == 0 {
		return func() {}, nil
	}
	data, err := renderBootstrapToolAccessors(info.Name, missing)
	if err != nil {
		return nil, err
	}
	path := filepath.Join(info.Dir, "zz_toolschema_bootstrap_tmp.go")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return nil, err
	}
	return func() {
		_ = os.Remove(path)
	}, nil
}

func existingFunctionNames(dir string) (map[string]struct{}, error) {
	matches, err := filepath.Glob(filepath.Join(dir, "*.go"))
	if err != nil {
		return nil, err
	}
	names := map[string]struct{}{}
	fileSet := token.NewFileSet()
	for _, path := range matches {
		if strings.HasSuffix(path, "_test.go") {
			continue
		}
		file, err := parser.ParseFile(fileSet, path, nil, 0)
		if err != nil {
			return nil, err
		}
		for _, declaration := range file.Decls {
			function, ok := declaration.(*ast.FuncDecl)
			if ok {
				names[function.Name.Name] = struct{}{}
			}
		}
	}
	return names, nil
}

func renderBootstrapToolAccessors(packageName string, tools []DiscoveredToolSpec) ([]byte, error) {
	var buffer bytes.Buffer
	buffer.WriteString("package ")
	buffer.WriteString(packageName)
	buffer.WriteString("\n\n")
	qualifier := "toolcontract."
	if packageName != "toolcontract" {
		buffer.WriteString("import toolcontract \"github.com/vihren-dev/vihren/platform/toolcontract\"\n\n")
	} else {
		qualifier = ""
	}
	for _, tool := range tools {
		buffer.WriteString("func ")
		buffer.WriteString(tool.AccessorName)
		buffer.WriteString("() ")
		buffer.WriteString(qualifier)
		buffer.WriteString("Tool {\n")
		buffer.WriteString("\treturn ")
		buffer.WriteString(qualifier)
		buffer.WriteString("NewSchemaDerivedTool(")
		buffer.WriteString(tool.SpecVariableName)
		buffer.WriteString(".Name, ")
		buffer.WriteString(tool.SpecVariableName)
		buffer.WriteString(".Description, nil, nil)\n")
		buffer.WriteString("}\n\n")
	}
	return format.Source(buffer.Bytes())
}
