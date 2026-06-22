# Embedded Temporal

Date: 2026-06-22

`platform/embeddedtemporal` runs a Temporal server inside your own process, so a
single Go binary can host durable workflows with **zero external
infrastructure** — no Docker, no daemon, no `temporal server start-dev`. It is
built for simple agent development: write an agent, press run, watch a durable
workflow execute.

## One line to a running server

```go
server, err := embeddedtemporal.Start()
if err != nil {
    return err
}
defer server.Close()
```

`Start` returns a server that is ready to use. It is the same in-process Temporal
server the BusyBox demo uses, wrapped behind a small API and made to return
errors instead of aborting (the underlying `temporaltest` is test-oriented).

## Use it

```go
// Host a worker. A generated Register call drops straight in.
server.StartWorker(myapp.DefaultTaskQueue, func(r worker.Registry) {
    myapp.Register(r, &myapp.Activities{ /* deps */ })
})

// Start a workflow through the generated typed client.
out, err := myapp.NewClient(server.Client()).MyWorkflow(ctx, opts, in)
```

| Method | Purpose |
| --- | --- |
| `Start(opts...) (*Server, error)` | Launch the in-process server. |
| `(*Server) Client() client.Client` | Connected client for the default namespace. |
| `(*Server) StartWorker(taskQueue, register)` | Register and start a worker. |
| `(*Server) HostPort() string` | Frontend `host:port`, to dial from elsewhere. |
| `(*Server) Namespace() string` | The configured namespace. |
| `(*Server) Close()` | Stop workers, clients, server; release the DB lock. |

`Client`, `HostPort`, and `Namespace` expose the underlying primitives, so you
can build any client or worker the Temporal SDK allows. `StartWorker` is sugar
over `worker.New(server.Client(), ...)`; it is not the only way to construct a
worker.

## Persistence: a durable single-binary desktop app

By default the server is **ephemeral** — state lives in memory and disappears
when the process exits. Pass `WithDatabaseFile` to persist everything to a SQLite
file instead, so the app can start, do work, stop, and resume later:

```go
server, err := embeddedtemporal.Start(
    embeddedtemporal.WithDatabaseFile(filepath.Join(home, ".myagent", "state.db")),
    embeddedtemporal.WithNamespace("myagent"), // stable across restarts
)
defer server.Close()
```

Workflow histories, timers, and signals survive a restart: relaunch on the same
file and a long-running agent picks up exactly where it left off. This needs no
Docker, daemon, or separate process — the entire durable runtime is in your one
binary, which is ideal for desktop apps.

| Option | Effect |
| --- | --- |
| `WithDatabaseFile(path)` | Persist to a SQLite file (creates file + parent dir). |
| `WithNamespace(name)` | Use a stable namespace (default `"default"`). |

### Schema migration across upgrades

A persisted database holds real user data, so it must survive bumps of the
`go.temporal.io/server` dependency. On every startup the database schema is made
current:

- a **fresh** file gets the full schema plus a stamped schema version;
- an **existing** file is **migrated forward** by applying the server's versioned
  history-store migrations that are newer than the stamp.

Migration is a no-op when the database already matches the linked server. After a
server-dependency bump, the next launch upgrades the on-disk schema automatically
before the server opens it.

**Limits worth knowing.** Only the history store (the source of truth) is
migrated; the visibility store ships a static schema and is rebuildable, so it is
not auto-migrated. A real cross-version upgrade is exercised by Temporal's own
versioned SQL but cannot be fully end-to-end tested in this repo (faking an older
baseline conflicts with the bundled latest schema) — validate the upgrade path
whenever you bump `go.temporal.io/server`.

### One writer at a time

A SQLite-backed server assumes a single owning process. Persistent `Start`
acquires a lock (`<db>.lock`) and releases it on `Close`; a second instance on the
same file fails fast with a clear error. A crash can leave a stale lock that must
be removed manually (PID-liveness recovery is a future refinement).

## When to use which Temporal

- **Embedded ephemeral (`platform/embeddedtemporal`)** — fastest inner loop,
  self-contained binaries, demos, blog examples. State disappears on exit. Not
  for production.
- **Embedded persistent (`WithDatabaseFile`)** — single-user desktop durable
  agents: one binary, no dependencies, state survives restarts. Single writer;
  not for production multi-tenant use.
- **External Temporal server** — a real Temporal frontend and Web UI, useful
  when you want to inspect workflow history or run multi-process setups.

## Runnable example

`examples/codegenhello/cmd/codegenhello-embedded` is the whole codegenhello
example in one process. Run it with:

```sh
just run-codegenhello-embedded
```

Or without `just`:

```sh
go run ./examples/codegenhello/cmd/codegenhello-embedded
```
