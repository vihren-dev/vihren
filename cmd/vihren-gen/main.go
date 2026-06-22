package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/vihren-dev/vihren/internal/codegen"
	"github.com/vihren-dev/vihren/internal/toolschema"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	flags := flag.NewFlagSet("vihren-gen", flag.ContinueOnError)
	dryRun := flags.Bool("dry-run", false, "validate and render without writing files")
	generatedFileName := flags.String("output-file", "", "generated Go filename per package")
	manifestFileName := flags.String("manifest", "", "manifest filename at the module root")
	if err := flags.Parse(args); err != nil {
		return err
	}
	patterns := flags.Args()
	if len(patterns) == 0 {
		return fmt.Errorf("at least one package pattern is required")
	}
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	moduleRoot, err := toolschema.FindModuleRoot(cwd)
	if err != nil {
		return err
	}
	request := codegen.GenerateRequest{
		ModuleRoot:        moduleRoot,
		Patterns:          patterns,
		GeneratedFileName: *generatedFileName,
		ManifestFileName:  *manifestFileName,
	}
	result, diagnostics, err := codegen.Generate(request)
	if err != nil {
		return err
	}
	if len(diagnostics) > 0 {
		return fmt.Errorf("codegen diagnostics:\n%s", codegen.FormatDiagnostics(diagnostics))
	}
	if *dryRun {
		return nil
	}
	return codegen.WriteGenerated(request, result)
}
