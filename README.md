# Vihren

[![Go Reference](https://pkg.go.dev/badge/github.com/vihren-dev/vihren.svg)](https://pkg.go.dev/github.com/vihren-dev/vihren)
[![Latest release](https://img.shields.io/github/v/release/vihren-dev/vihren?sort=semver)](https://github.com/vihren-dev/vihren/releases)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)

Vihren is a Go toolkit for building typed Temporal workflows with less
boilerplate. You write ordinary Go activities and workflows, mark them with
`//vihren:activity` and `//vihren:workflow`, and run `vihren-gen` to generate
Temporal registration, workflow-side activity proxies, and typed workflow
clients.

The current public release focuses on one concrete path:

- typed workflow and activity code generation;
- a small runnable `codegenhello` example;
- an embedded Temporal runtime for local demos and single-process development.

## Install

Requires Go 1.26 or newer.

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
- `NewClient`, a typed workflow client for starting workflows synchronously or
  asynchronously;
- typed workflow run handles, so async callers can await typed results while
  still reading Temporal workflow IDs.

In this example the compact `workflow.go` you write generates registration,
proxy, client, and run-handle code you never hand-write or hand-maintain.

## Run The Example

From this repository, run the self-contained example directly:

```sh
go run ./examples/codegenhello/cmd/codegenhello-embedded
```

It starts Temporal inside the process, starts a worker, runs the generated typed
workflow client, and prints:

```text
Hello, Ada
```

No Docker, daemon, or separate Temporal development server is required.

If you use `just`, the same command is available as:

```sh
just run-codegenhello-embedded
```

To verify generated code and tests from a checkout:

```sh
go generate ./examples/codegenhello
go test ./examples/codegenhello ./examples/codegenhello/cmd/codegenhello-embedded -timeout 60s
```

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

## License

Vihren is licensed under the [Apache License, Version 2.0](LICENSE). See the
[NOTICE](NOTICE) file for attribution.
