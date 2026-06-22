package codegen

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"
)

// DiscoverConfig configures package loading for marker discovery.
type DiscoverConfig struct {
	Dir string
}

// Discover loads packages and returns their vihren marker records.
func Discover(config DiscoverConfig, patterns ...string) ([]DiscoveredPackage, []Diagnostic, error) {
	fileSet := token.NewFileSet()
	loaded, err := packages.Load(&packages.Config{
		Dir:   config.Dir,
		Fset:  fileSet,
		Mode:  packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo,
		Tests: false,
	}, patterns...)
	if err != nil {
		return nil, nil, err
	}
	var diagnostics []Diagnostic
	var discovered []DiscoveredPackage
	for _, pkg := range loaded {
		for _, pkgErr := range pkg.Errors {
			if isGeneratedBootstrapError(pkgErr.Msg) {
				continue
			}
			diagnostics = append(diagnostics, Diagnostic{Position: pkgErr.Pos, Message: pkgErr.Msg})
		}
		result, pkgDiagnostics := discoverPackage(fileSet, pkg)
		diagnostics = append(diagnostics, pkgDiagnostics...)
		if len(result.Markers) > 0 {
			discovered = append(discovered, result)
		}
	}
	sort.Slice(discovered, func(i int, j int) bool {
		return discovered[i].ImportPath < discovered[j].ImportPath
	})
	return discovered, diagnostics, nil
}

// isGeneratedBootstrapError suppresses missing generated-symbol errors so a
// workflow can call Activity proxies before vihren.gen.go exists.
func isGeneratedBootstrapError(message string) bool {
	if strings.HasPrefix(message, "# ") {
		return true
	}
	for _, generatedName := range []string{
		"Activity",
		"Register",
		"RegisterActivities",
		"RegisterWorkflows",
		"Client",
		"NewClient",
	} {
		if strings.Contains(message, "undefined: "+generatedName) {
			return true
		}
	}
	return false
}

// discoverPackage extracts markers from one loaded package.
func discoverPackage(fileSet *token.FileSet, pkg *packages.Package) (DiscoveredPackage, []Diagnostic) {
	result := DiscoveredPackage{
		Name:       pkg.Name,
		Dir:        packageDir(pkg),
		ImportPath: pkg.PkgPath,
		types:      pkg.Types,
	}
	var diagnostics []Diagnostic
	for _, file := range pkg.Syntax {
		for _, declaration := range file.Decls {
			fn, ok := declaration.(*ast.FuncDecl)
			if !ok || fn.Doc == nil {
				continue
			}
			for _, comment := range fn.Doc.List {
				kind, options, ok, markerDiagnostics := parseMarkerComment(fileSet, comment.Text, comment.Pos())
				diagnostics = append(diagnostics, markerDiagnostics...)
				if !ok {
					continue
				}
				marker, validationDiagnostics := markerFromFunction(fileSet, pkg, fn, kind, options, comment.Pos())
				diagnostics = append(diagnostics, validationDiagnostics...)
				result.Markers = append(result.Markers, marker)
			}
		}
	}
	diagnostics = append(diagnostics, validatePackageProxyNames(result.Markers)...)
	sort.Slice(result.Markers, func(i int, j int) bool {
		return result.Markers[i].Position < result.Markers[j].Position
	})
	return result, diagnostics
}

// markerFromFunction normalizes one marked function declaration.
func markerFromFunction(
	fileSet *token.FileSet,
	pkg *packages.Package,
	fn *ast.FuncDecl,
	kind MarkerKind,
	options MarkerOptions,
	markerPos token.Pos,
) (DiscoveredMarker, []Diagnostic) {
	position := fileSet.Position(markerPos).String()
	proxy := fn.Name.Name
	if options.Proxy != "" {
		proxy = options.Proxy
	}
	temporalName := options.Name
	if temporalName == "" {
		temporalName = pkg.PkgPath + "." + proxy
	}
	marker := DiscoveredMarker{
		Kind:         kind,
		PackageName:  pkg.Name,
		PackageDir:   packageDir(pkg),
		ImportPath:   pkg.PkgPath,
		FunctionName: fn.Name.Name,
		ReceiverName: receiverName(pkg, fn),
		ProxyName:    proxy,
		TemporalName: temporalName,
		ExplicitName: options.Name != "",
		Position:     position,
		Options:      options,
	}
	signature, ok := functionSignature(pkg, fn)
	if !ok {
		return marker, []Diagnostic{{Position: position, Message: "marked declaration has no function signature"}}
	}
	shape, diagnostics := validateSignature(kind, signature, position)
	marker.InputCount = shape.InputCount
	marker.HasOutput = shape.HasOutput
	marker.signature = signature
	return marker, diagnostics
}

// functionSignature returns the typed signature for an AST function declaration.
func functionSignature(pkg *packages.Package, fn *ast.FuncDecl) (*types.Signature, bool) {
	obj, ok := pkg.TypesInfo.Defs[fn.Name].(*types.Func)
	if !ok || obj == nil {
		return nil, false
	}
	signature, ok := obj.Type().(*types.Signature)
	return signature, ok
}

// receiverName returns the base receiver type name for a method.
func receiverName(pkg *packages.Package, fn *ast.FuncDecl) string {
	signature, ok := functionSignature(pkg, fn)
	if !ok || signature.Recv() == nil {
		return ""
	}
	receiverType := types.Unalias(signature.Recv().Type())
	if pointer, ok := receiverType.(*types.Pointer); ok {
		receiverType = types.Unalias(pointer.Elem())
	}
	if named, ok := receiverType.(*types.Named); ok {
		return named.Obj().Name()
	}
	return fmt.Sprintf("%s", receiverType)
}

// packageDir returns the first compiled Go file's package directory.
func packageDir(pkg *packages.Package) string {
	if len(pkg.GoFiles) == 0 {
		return ""
	}
	return filepath.Dir(pkg.GoFiles[0])
}

// validatePackageProxyNames reports duplicate package-local activity proxy names.
func validatePackageProxyNames(markers []DiscoveredMarker) []Diagnostic {
	var diagnostics []Diagnostic
	seen := map[string]DiscoveredMarker{}
	for _, marker := range markers {
		if marker.Kind != ActivityMarker {
			continue
		}
		existing, ok := seen[marker.ProxyName]
		if !ok {
			seen[marker.ProxyName] = marker
			continue
		}
		diagnostics = append(diagnostics, Diagnostic{
			Position: marker.Position,
			Message: fmt.Sprintf(
				"activity proxy name %q collides with %s at %s",
				marker.ProxyName,
				existing.FunctionName,
				existing.Position,
			),
		})
	}
	return diagnostics
}
