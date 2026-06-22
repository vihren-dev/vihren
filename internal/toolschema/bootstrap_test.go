package toolschema

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderBootstrapToolAccessors(t *testing.T) {
	t.Parallel()
	data, err := renderBootstrapToolAccessors("fixturepkg", []DiscoveredToolSpec{{
		SpecVariableName: "ExampleToolSpec",
		AccessorName:     "ExampleTool",
	}})
	if err != nil {
		t.Fatalf("bootstrap render should succeed: %v", err)
	}
	output := string(data)
	for _, expected := range []string{
		`"github.com/vihren-dev/vihren/platform/toolcontract"`,
		"func ExampleTool() toolcontract.Tool",
		"toolcontract.NewSchemaDerivedTool(ExampleToolSpec.Name, ExampleToolSpec.Description, nil, nil)",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("bootstrap output missing %s: %s", expected, output)
		}
	}
}

func TestWriteBootstrapToolAccessorsSkipsExistingAccessor(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	source := []byte("package fixturepkg\n\nfunc ExampleTool() {}\n")
	if err := os.WriteFile(filepath.Join(dir, "tool.go"), source, 0o644); err != nil {
		t.Fatalf("write fixture source: %v", err)
	}
	cleanup, err := writeBootstrapToolAccessors(PackageInfo{Name: "fixturepkg", Dir: dir}, []DiscoveredToolSpec{{
		AccessorName: "ExampleTool",
	}})
	if err != nil {
		t.Fatalf("bootstrap write should succeed: %v", err)
	}
	defer cleanup()
	if _, err := os.Stat(filepath.Join(dir, "zz_toolschema_bootstrap_tmp.go")); !os.IsNotExist(err) {
		t.Fatalf("bootstrap file should not be written when accessor exists: %v", err)
	}
}
