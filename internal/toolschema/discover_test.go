package toolschema

import "testing"

func TestDiscoverToolSpecsFindsFixturePackageSpec(t *testing.T) {
	t.Parallel()
	moduleRoot, err := FindModuleRoot(".")
	if err != nil {
		t.Fatalf("module root should resolve: %v", err)
	}
	info, err := loadPackageInfo(t.Context(), moduleRoot, "./internal/toolschema/testdata/fixturepkg")
	if err != nil {
		t.Fatalf("package info should load: %v", err)
	}
	tools, err := DiscoverToolSpecs(info)
	if err != nil {
		t.Fatalf("discover should succeed: %v", err)
	}
	if len(tools) != 1 {
		t.Fatalf("expected one discovered tool: %#v", tools)
	}
	tool := tools[0]
	if tool.SpecVariableName != "ExampleToolSpec" {
		t.Fatalf("unexpected spec variable: %#v", tool)
	}
	if tool.AccessorName != "ExampleTool" || tool.SingletonName != "exampleTool" {
		t.Fatalf("unexpected generated names: %#v", tool)
	}
	if tool.InputTypeName != "ExampleInput" || tool.OutputTypeName != "ExampleOutput" {
		t.Fatalf("unexpected type arguments: %#v", tool)
	}
}

func TestNewDiscoveredToolSpecRequiresSpecSuffix(t *testing.T) {
	t.Parallel()
	if _, err := newDiscoveredToolSpec("ExampleTool", "ExampleInput", "ExampleOutput"); err == nil {
		t.Fatal("expected missing Spec suffix to fail")
	}
}
