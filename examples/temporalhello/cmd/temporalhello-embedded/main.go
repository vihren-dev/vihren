// Command temporalhello-embedded runs the vanilla Temporal counterpart to
// codegenhello in a single process: it starts an embedded Temporal server,
// manually registers the activity and workflow, starts the workflow through
// client.ExecuteWorkflow, and awaits the result through WorkflowRun.Get.
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"

	"github.com/vihren-dev/vihren/examples/temporalhello"
	"github.com/vihren-dev/vihren/platform/embeddedtemporal"
)

func main() {
	if err := run(context.Background(), os.Stdout); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, out io.Writer) error {
	// 1. Start Temporal inside this process. Ephemeral; gone when main returns.
	server, err := embeddedtemporal.Start(embeddedtemporal.WithClientOptions(client.Options{Logger: quietLogger{}}))
	if err != nil {
		return err
	}
	defer server.Close()

	// 2. Host a worker. Registration uses the Temporal SDK directly.
	if _, err := server.StartWorker(temporalhello.DefaultTaskQueue, func(registry worker.Registry) {
		activities := &temporalhello.GreetingActivities{Prefix: "Hello"}
		registry.RegisterActivityWithOptions(
			activities.ComposeGreeting,
			activity.RegisterOptions{Name: temporalhello.ComposeGreetingActivityName},
		)
		registry.RegisterWorkflowWithOptions(
			temporalhello.HelloWorkflow,
			workflow.RegisterOptions{Name: temporalhello.HelloWorkflowName},
		)
	}); err != nil {
		return err
	}

	// 3. Start and await the workflow through the vanilla Temporal client.
	run, err := server.Client().ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{ID: "temporalhello-embedded", TaskQueue: temporalhello.DefaultTaskQueue},
		temporalhello.HelloWorkflowName,
		temporalhello.GreetingInput{Name: "Ada"},
	)
	if err != nil {
		return err
	}
	var result temporalhello.GreetingOutput
	if err := run.Get(ctx, &result); err != nil {
		return err
	}
	_, err = fmt.Fprintln(out, result.Message)
	return err
}

// quietLogger keeps the Hello World output focused on the workflow result
// while still letting production users provide their own Temporal logger.
type quietLogger struct{}

func (quietLogger) Debug(string, ...interface{}) {}
func (quietLogger) Info(string, ...interface{})  {}
func (quietLogger) Warn(string, ...interface{})  {}
func (quietLogger) Error(string, ...interface{}) {}
