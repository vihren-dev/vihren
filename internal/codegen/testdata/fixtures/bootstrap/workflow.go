package bootstrap

import (
	"context"

	"go.temporal.io/sdk/workflow"
)

// Input is the bootstrap workflow input.
type Input struct {
	Value string
}

// Output is the bootstrap workflow output.
type Output struct {
	Value string
}

// Touch is a marked activity that will be called through generated Activity.
//
//vihren:activity
func Touch(ctx context.Context, in Input) (Output, error) {
	_ = ctx
	return Output{Value: in.Value}, nil
}

// Run calls the generated Activity proxy before vihren.gen.go exists.
//
//vihren:workflow
func Run(ctx workflow.Context, in Input) (Output, error) {
	return Activity.Touch(ctx, in)
}
