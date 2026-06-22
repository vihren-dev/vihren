// Package codegenhello is the minimal vihren-gen example: a developer writes
// one activity and one workflow, and the generator supplies all Temporal
// registration and typed-call plumbing (see vihren.gen.go). It is intentionally
// small enough to read end to end in a blog post.
package codegenhello

//go:generate go run ../../cmd/vihren-gen --manifest .cache/codegenhello.manifest.json ./examples/codegenhello

import (
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"
)

// DefaultTaskQueue is the Temporal task queue shared by the worker and the
// workflow starter in this example.
const DefaultTaskQueue = "vihren-codegenhello"

// GreetingActivities holds the process-local dependencies that generated
// activities run with.
type GreetingActivities struct {
	Prefix string
}

// GreetingInput is the typed input for both the workflow and the activity.
type GreetingInput struct {
	Name string
}

// GreetingOutput is the typed result returned to the caller.
type GreetingOutput struct {
	Message string
}

// ComposeGreeting builds a greeting. The marker tells vihren-gen to register it
// as a Temporal activity and to expose a type-checked proxy method.
//
//vihren:activity
func (activities *GreetingActivities) ComposeGreeting(ctx context.Context, in GreetingInput) (GreetingOutput, error) {
	_ = ctx
	return GreetingOutput{Message: fmt.Sprintf("%s, %s", activities.Prefix, in.Name)}, nil
}

// HelloWorkflow schedules the activity through the generated Activity proxy, so
// the activity name and argument types are checked at compile time.
//
//vihren:workflow
func HelloWorkflow(ctx workflow.Context, in GreetingInput) (GreetingOutput, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{StartToCloseTimeout: time.Second})
	return Activity.ComposeGreeting(ctx, in)
}
