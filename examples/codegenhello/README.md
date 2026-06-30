# Codegen Hello

The smallest possible `vihren-gen` example. The developer writes one activity
and one workflow; the generator writes every line of Temporal registration and
typed-call plumbing.

Hand-written code (`workflow.go`, ~50 lines):

- `ComposeGreeting`, tagged `//vihren:activity`;
- `HelloWorkflow`, tagged `//vihren:workflow`, which schedules the activity
  through the generated `Activity` proxy.

Generated code (`vihren.gen.go`):

- `Register(r, *GreetingActivities)` — the whole worker registration surface;
- the workflow-side `Activity` proxy (compile-time-checked activity calls);
- the typed `Client.HelloWorkflow` / `Client.HelloWorkflowAsync` start methods;
- `HelloWorkflowRun`, the typed async run handle with `Get(ctx)`, `GetID()`,
  and `GetRunID()`.

The `cmd/` programs call the generated code directly:

- `codegenhello-embedded` — the whole example in one process: it starts an
  embedded Temporal server (`platform/embeddedtemporal`), hosts a worker via the
  generated `Register`, and starts the workflow through the generated client. No
  Docker, no daemon. This is the self-contained, blog-pasteable program.
- `codegenhello-worker` / `codegenhello-start` — the same logic split across a
  worker process and a starter process, for running against a real Temporal
  server.

The sync client path returns the typed workflow result directly. The async path
returns `HelloWorkflowRun`, so callers can start the workflow, keep the Temporal
workflow ID/run ID, and later await the typed `GreetingOutput` with
`run.Get(ctx)`.

Run the whole thing in one command, zero infrastructure:

```sh
go run ./examples/codegenhello/cmd/codegenhello-embedded
```

Expected output:

```text
Hello, Ada
```

If you use `just`, the same command is available as:

```sh
just run-codegenhello-embedded
```

Regenerate the code:

```sh
go generate ./examples/codegenhello
```

Check the generated workflow path:

```sh
go test ./examples/codegenhello ./examples/codegenhello/cmd/codegenhello-embedded -timeout 60s
```

Run the split worker/starter against a Temporal server you already have running:

```sh
go run ./examples/codegenhello/cmd/codegenhello-worker
go run ./examples/codegenhello/cmd/codegenhello-start Ada
```
