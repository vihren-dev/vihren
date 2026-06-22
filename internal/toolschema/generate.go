package toolschema

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

// reflectionHelperMu serializes generation sections that temporarily mutate a
// target package directory for bootstrap accessors.
var reflectionHelperMu sync.Mutex

// GenerateGoFile derives schemas for the configured package and renders a Go
// artifact containing json.RawMessage variables.
func GenerateGoFile(ctx context.Context, request GenerateRequest) ([]byte, error) {
	info, err := loadPackageInfo(ctx, request.ModuleRoot, request.Package)
	if err != nil {
		return nil, err
	}
	outputPackage := strings.TrimSpace(request.OutputPackage)
	if outputPackage == "" {
		outputPackage = info.Name
	}
	generatedFile, err := generationFileForRequest(info, request.Schemas)
	if err != nil {
		return nil, err
	}
	schemas, err := runReflectionHelper(ctx, request.ModuleRoot, info, generatedFile.Schemas, generatedFile.Tools)
	if err != nil {
		return nil, err
	}
	return RenderGeneratedGoFile(outputPackage, schemas, generatedFile.Tools)
}

func generationFileForRequest(info PackageInfo, explicitSchemas []SchemaSpec) (GeneratedGoFile, error) {
	if len(explicitSchemas) > 0 {
		return GeneratedGoFile{Schemas: explicitSchemas}, nil
	}
	tools, err := DiscoverToolSpecs(info)
	if err != nil {
		return GeneratedGoFile{}, err
	}
	if len(tools) == 0 {
		return GeneratedGoFile{}, fmt.Errorf("at least one schema spec or ToolSpec value is required")
	}
	schemas := make([]SchemaSpec, 0, len(tools)*2)
	for _, tool := range tools {
		schemas = append(schemas,
			SchemaSpec{VariableName: tool.InputSchemaVariableName, TypeName: tool.InputTypeName},
			SchemaSpec{VariableName: tool.OutputSchemaVariableName, TypeName: tool.OutputTypeName},
		)
	}
	return GeneratedGoFile{Schemas: schemas, Tools: tools}, nil
}

func loadPackageInfo(ctx context.Context, moduleRoot string, packagePattern string) (PackageInfo, error) {
	fields := []string{"{{.ImportPath}}", "{{.Dir}}", "{{.Name}}"}
	output, err := runGo(ctx, moduleRoot, "list", "-f", strings.Join(fields, " "), packagePattern)
	if err != nil {
		return PackageInfo{}, err
	}
	parts := strings.SplitN(strings.TrimSpace(string(output)), " ", 3)
	if len(parts) != 3 {
		return PackageInfo{}, fmt.Errorf("unexpected go list output: %s", strings.TrimSpace(string(output)))
	}
	moduleOutput, err := runGo(ctx, moduleRoot, "list", "-m")
	if err != nil {
		return PackageInfo{}, err
	}
	return PackageInfo{
		ImportPath: parts[0],
		Dir:        parts[1],
		Name:       parts[2],
		ModulePath: strings.TrimSpace(string(moduleOutput)),
	}, nil
}

func runReflectionHelper(
	ctx context.Context,
	moduleRoot string,
	info PackageInfo,
	schemas []SchemaSpec,
	tools []DiscoveredToolSpec,
) (map[string]json.RawMessage, error) {
	if len(schemas) == 0 {
		return nil, fmt.Errorf("at least one schema spec is required")
	}
	reflectionHelperMu.Lock()
	defer reflectionHelperMu.Unlock()

	cleanup, err := writeBootstrapToolAccessors(info, tools)
	if err != nil {
		return nil, err
	}
	defer cleanup()
	cacheRoot := filepath.Join(moduleRoot, ".cache", "toolschema")
	if err := os.MkdirAll(cacheRoot, 0o755); err != nil {
		return nil, err
	}
	tempDir, err := os.MkdirTemp(cacheRoot, "reflect-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tempDir)
	source, err := renderReflectionHelperSource(moduleRoot, info, schemas)
	if err != nil {
		return nil, err
	}
	sourcePath := filepath.Join(tempDir, "main.go")
	if err := os.WriteFile(sourcePath, source, 0o644); err != nil {
		return nil, err
	}
	output, err := runGo(ctx, moduleRoot, "run", sourcePath)
	if err != nil {
		return nil, err
	}
	var reflected map[string]json.RawMessage
	if err := json.Unmarshal(output, &reflected); err != nil {
		return nil, err
	}
	return reflected, nil
}

func runGo(ctx context.Context, dir string, args ...string) ([]byte, error) {
	command := exec.CommandContext(ctx, "go", args...)
	command.Dir = dir
	command.Env = append(os.Environ(),
		"GOPATH="+filepath.Join(dir, ".cache", "go"),
		"GOCACHE="+filepath.Join(dir, ".cache", "go-build"),
	)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	err := command.Run()
	if err != nil {
		return nil, fmt.Errorf("go %s failed: %w\n%s", strings.Join(args, " "), err, stderr.String())
	}
	return stdout.Bytes(), nil
}

func renderReflectionHelperSource(moduleRoot string, info PackageInfo, schemas []SchemaSpec) ([]byte, error) {
	commentPath, err := filepath.Rel(moduleRoot, info.Dir)
	if err != nil {
		return nil, err
	}
	var buffer bytes.Buffer
	buffer.WriteString("package main\n\n")
	buffer.WriteString("import (\n")
	buffer.WriteString("\t\"encoding/json\"\n")
	buffer.WriteString("\t\"fmt\"\n")
	buffer.WriteString("\t\"os\"\n")
	buffer.WriteString("\tjsonschema \"github.com/invopop/jsonschema\"\n")
	buffer.WriteString("\ttarget \"")
	buffer.WriteString(info.ImportPath)
	buffer.WriteString("\"\n")
	buffer.WriteString(")\n\n")
	buffer.WriteString("func main() {\n")
	buffer.WriteString("\treflector := &jsonschema.Reflector{ExpandedStruct: true}\n")
	buffer.WriteString("\tif err := reflector.AddGoComments(\"")
	buffer.WriteString(info.ModulePath)
	buffer.WriteString("\", \"")
	buffer.WriteString(filepath.ToSlash(commentPath))
	buffer.WriteString("\"); err != nil {\n")
	buffer.WriteString("\t\tfmt.Fprintln(os.Stderr, err)\n\t\tos.Exit(1)\n\t}\n")
	buffer.WriteString("\tschemas := map[string]any{\n")
	for _, spec := range schemas {
		if _, err := ParseSchemaSpec(FormatSchemaSpec(spec)); err != nil {
			return nil, err
		}
		buffer.WriteString("\t\t\"")
		buffer.WriteString(spec.VariableName)
		buffer.WriteString("\": reflector.Reflect(&target.")
		buffer.WriteString(spec.TypeName)
		buffer.WriteString("{}),\n")
	}
	buffer.WriteString("\t}\n")
	buffer.WriteString("\tdata, err := json.Marshal(schemas)\n")
	buffer.WriteString("\tif err != nil {\n\t\tfmt.Fprintln(os.Stderr, err)\n\t\tos.Exit(1)\n\t}\n")
	buffer.WriteString("\tos.Stdout.Write(data)\n")
	buffer.WriteString("}\n")
	return buffer.Bytes(), nil
}
