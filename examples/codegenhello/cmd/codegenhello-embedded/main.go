// Command codegenhello-embedded runs the whole codegenhello example in a single
// process: it starts an embedded Temporal server, hosts a worker for the
// generated activity and workflow, and starts the workflow through the generated
// typed client — no Docker, no daemon, no `just temporal-start`.
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"github.com/vihren-dev/vihren/examples/codegenhello"
	"github.com/vihren-dev/vihren/platform/embeddedtemporal"
)

func main() {
	if err := run(context.Background(), os.Stdout); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, out io.Writer) error {
	// 1. Start Temporal inside this process. Ephemeral; gone when main returns.
	server, err := embeddedtemporal.Start()
	if err != nil {
		return err
	}
	defer server.Close()

	// 2. Host a worker. The generated Register call is the whole registration.
	if _, err := server.StartWorker(codegenhello.DefaultTaskQueue, func(registry worker.Registry) {
		codegenhello.Register(registry, &codegenhello.GreetingActivities{Prefix: "Hello"})
	}); err != nil {
		return err
	}

	// 3. Start and await the workflow through the generated typed client.
	result, err := codegenhello.NewClient(server.Client()).HelloWorkflow(
		ctx,
		client.StartWorkflowOptions{ID: "codegenhello-embedded", TaskQueue: codegenhello.DefaultTaskQueue},
		codegenhello.GreetingInput{Name: "Ada"},
	)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(out, result.Message)
	return err
}
