# Vihren

Vihren is a Go toolkit for building typed Temporal workflows with less
boilerplate. You write ordinary Go activities and workflows, mark them with
`//vihren:activity` and `//vihren:workflow`, and run `vihren-gen` to generate
Temporal registration, workflow-side activity proxies, and typed workflow
clients.

The v0.1 release focuses on one concrete path:

- typed workflow and activity code generation;
- a small runnable `codegenhello` example;
- an embedded Temporal runtime for local demos and single-process development.

## Install

```sh
go get github.com/vihren-dev/vihren
go install github.com/vihren-dev/vihren/cmd/vihren-gen@latest
```

## The Smallest Example

Write an activity and a workflow:

```go
package codegenhello

import (
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"
)

//go:generate vihren-gen ./...

const DefaultTaskQueue = "vihren-codegenhello"

type GreetingActivities struct {
	Prefix string
}

type GreetingInput struct {
	Name string
}

type GreetingOutput struct {
	Message string
}

// ComposeGreeting builds a greeting.
//
//vihren:activity
func (activities *GreetingActivities) ComposeGreeting(ctx context.Context, in GreetingInput) (GreetingOutput, error) {
	_ = ctx
	return GreetingOutput{Message: fmt.Sprintf("%s, %s", activities.Prefix, in.Name)}, nil
}

// HelloWorkflow calls the activity through the generated type-safe proxy.
//
//vihren:workflow
func HelloWorkflow(ctx workflow.Context, in GreetingInput) (GreetingOutput, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{StartToCloseTimeout: time.Second})
	return Activity.ComposeGreeting(ctx, in)
}
```

Run generation:

```sh
go generate ./...
```

The generated file contains:

- `Register`, which registers the workflow and activity on a Temporal worker;
- `Activity`, a workflow-safe activity proxy with compile-time checked
  arguments and result types;
- `NewClient`, a typed workflow client for starting and awaiting workflows.

## Run The Example

From this repository:

```sh
just run-codegenhello-embedded
```

Or without `just`:

```sh
go run ./examples/codegenhello/cmd/codegenhello-embedded
```

That command starts Temporal inside the process, starts a worker, runs the
generated typed workflow client, and prints:

```text
Hello, Ada
```

No Docker, daemon, or separate Temporal development server is required.

## Embedded Temporal

`platform/embeddedtemporal` runs a Temporal server inside the current process.
It is useful for examples, demos, and single-binary local development. It is not
a production Temporal deployment.

```go
server, err := embeddedtemporal.Start()
if err != nil {
	return err
}
defer server.Close()
```

See [docs/embedded-temporal.md](docs/embedded-temporal.md) and
[examples/codegenhello](examples/codegenhello).
