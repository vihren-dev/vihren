package explicitone

import "context"

// Input is the activity input.
type Input struct{}

// Output is the activity output.
type Output struct{}

// SharedOne uses a worker-wide explicit activity name.
//
//vihren:activity name=shared.operation
func SharedOne(ctx context.Context, in Input) (Output, error) {
	_ = ctx
	_ = in
	return Output{}, nil
}
