# Temporal Hello

The vanilla Temporal SDK counterpart to `examples/codegenhello`. It implements
the same greeting workflow without `vihren-gen`, generated files,
`//vihren:*` markers, or generated `Activity`, `Register`, and `NewClient`
helpers.

Hand-written workflow/activity code (`workflow.go`):

- `ComposeGreeting`, an ordinary Temporal activity method;
- `HelloWorkflow`, which schedules the activity with
  `workflow.ExecuteActivity(...).Get(...)`;
- explicit `ComposeGreetingActivityName` and `HelloWorkflowName` constants.

Manual worker registration (`cmd/temporalhello-embedded/main.go`):

- `RegisterActivityWithOptions(..., activity.RegisterOptions{Name: ...})`;
- `RegisterWorkflowWithOptions(..., workflow.RegisterOptions{Name: ...})`.

Manual client execution (`cmd/temporalhello-embedded/main.go`):

- `client.ExecuteWorkflow(..., HelloWorkflowName, GreetingInput{...})`;
- `WorkflowRun.Get(..., &GreetingOutput{})`.

Run the whole thing in one command, zero infrastructure:

```sh
go run ./examples/temporalhello/cmd/temporalhello-embedded
```

Expected output:

```text
Hello, Ada
```

Check the vanilla workflow path:

```sh
go test ./examples/temporalhello ./examples/temporalhello/cmd/temporalhello-embedded -timeout 60s
```

