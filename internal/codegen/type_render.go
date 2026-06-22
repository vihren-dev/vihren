package codegen

import (
	"go/types"
	"path"
	"sort"
)

type typeRenderer struct {
	imports *importSet
}

func (renderer *typeRenderer) render(typ types.Type) string {
	return types.TypeString(typ, renderer.imports.qualifier)
}

type importSpec struct {
	name       string
	importPath string
}

type importSet struct {
	currentPath string
	imports     map[string]importSpec
	usedNames   map[string]string
}

func newImportSet(current *types.Package, static map[string]string) *importSet {
	currentPath := ""
	if current != nil {
		currentPath = current.Path()
	}
	imports := &importSet{
		currentPath: currentPath,
		imports:     map[string]importSpec{},
		usedNames:   map[string]string{},
	}
	for name, importPath := range static {
		imports.imports[importPath] = importSpec{name: name, importPath: importPath}
		imports.usedNames[name] = importPath
	}
	return imports
}

func (imports *importSet) qualifier(pkg *types.Package) string {
	if pkg == nil || pkg.Path() == "" || pkg.Path() == imports.currentPath {
		return ""
	}
	if spec, ok := imports.imports[pkg.Path()]; ok {
		return spec.name
	}
	name := pkg.Name()
	if name == "" {
		name = path.Base(pkg.Path())
	}
	base := name
	for index := 2; ; index++ {
		if existingPath, ok := imports.usedNames[name]; !ok || existingPath == pkg.Path() {
			imports.usedNames[name] = pkg.Path()
			imports.imports[pkg.Path()] = importSpec{name: name, importPath: pkg.Path()}
			return name
		}
		name = base + intSuffix(index)
	}
}

func (imports *importSet) list() []importSpec {
	specs := make([]importSpec, 0, len(imports.imports))
	for _, spec := range imports.imports {
		specs = append(specs, spec)
	}
	sort.Slice(specs, func(i int, j int) bool {
		return specs[i].importPath < specs[j].importPath
	})
	return specs
}

func intSuffix(value int) string {
	const digits = "0123456789"
	if value == 0 {
		return "0"
	}
	var reversed []byte
	for value > 0 {
		reversed = append(reversed, digits[value%10])
		value = value / 10
	}
	for left, right := 0, len(reversed)-1; left < right; left, right = left+1, right-1 {
		reversed[left], reversed[right] = reversed[right], reversed[left]
	}
	return string(reversed)
}
