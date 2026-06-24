// Package temporalhello is the vanilla Temporal SDK counterpart to
// examples/codegenhello. It keeps the same business behavior while writing the
// activity name, workflow name, activity scheduling, worker registration, and
// client execution code by hand.
package temporalhello

import (
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"
)

const (
	// DefaultTaskQueue is the Temporal task queue shared by the worker and the
	// workflow starter in this example.
	DefaultTaskQueue = "vihren-temporalhello"

	// ComposeGreetingActivityName is the manual Temporal activity type.
	ComposeGreetingActivityName = "github.com/vihren-dev/vihren/examples/temporalhello.ComposeGreeting"

	// HelloWorkflowName is the manual Temporal workflow type.
	HelloWorkflowName = "github.com/vihren-dev/vihren/examples/temporalhello.HelloWorkflow"
)

// GreetingActivities holds the process-local dependencies that activities run
// with.
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

// ComposeGreeting builds a greeting using ordinary Temporal activity code.
func (activities *GreetingActivities) ComposeGreeting(ctx context.Context, in GreetingInput) (GreetingOutput, error) {
	_ = ctx
	return GreetingOutput{Message: fmt.Sprintf("%s, %s", activities.Prefix, in.Name)}, nil
}

// HelloWorkflow schedules ComposeGreeting through the Temporal SDK directly.
func HelloWorkflow(ctx workflow.Context, in GreetingInput) (GreetingOutput, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{StartToCloseTimeout: time.Second})

	var out GreetingOutput
	err := workflow.ExecuteActivity(ctx, ComposeGreetingActivityName, in).Get(ctx, &out)
	return out, err
}
