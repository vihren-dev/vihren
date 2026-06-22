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
- the typed `Client.HelloWorkflow` / `Client.HelloWorkflowAsync` start methods.

The `cmd/` programs call the generated code directly:

- `codegenhello-embedded` — the whole example in one process: it starts an
  embedded Temporal server (`platform/embeddedtemporal`), hosts a worker via the
  generated `Register`, and starts the workflow through the generated client. No
  Docker, no daemon. This is the self-contained, blog-pasteable program.
- `codegenhello-worker` / `codegenhello-start` — the same logic split across a
  worker process and a starter process, for running against a real Temporal
  server.

Run the whole thing in one command, zero infrastructure:

```sh
just run-codegenhello-embedded
```

Or run the same command directly:

```sh
go run ./examples/codegenhello/cmd/codegenhello-embedded
```

Regenerate the code:

```sh
go generate ./examples/codegenhello
```

Run the split worker/starter against a Temporal server you already have running:

```sh
go run ./examples/codegenhello/cmd/codegenhello-worker
go run ./examples/codegenhello/cmd/codegenhello-start Ada
```
