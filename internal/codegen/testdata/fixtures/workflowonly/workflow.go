package workflowonly

import "go.temporal.io/sdk/workflow"

// EchoInput is the workflow-only input.
type EchoInput struct {
	Value string
}

// EchoOutput is the workflow-only output.
type EchoOutput struct {
	Value string
}

// Echo echoes its input without activities.
//
//vihren:workflow
func Echo(ctx workflow.Context, in EchoInput) (EchoOutput, error) {
	_ = ctx
	return EchoOutput{Value: in.Value}, nil
}
