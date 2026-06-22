package namesbeta

import "context"

// Input is the activity input.
type Input struct{}

// Output is the activity output.
type Output struct{}

// Shared uses a default import-path-scoped activity name.
//
//vihren:activity
func Shared(ctx context.Context, in Input) (Output, error) {
	_ = ctx
	_ = in
	return Output{}, nil
}
