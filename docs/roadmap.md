# Vihren Roadmap

Vihren is an open-source Go toolkit for building AI agents that are reliable
enough to run real, long-running work — the kind of process that can take days,
needs people to step in at the right moments, and has to keep going even if a
server restarts.

Most agent frameworks are great for demos but fragile in production: a crash
loses the agent's place, a multi-day wait means babysitting a process that must
not die, and running code the model wrote means standing up extra infrastructure
to do it safely. Vihren is built on [Temporal](https://temporal.io), a durable
execution engine, so an agent automatically survives crashes and restarts, can
pause for hours or days, and resumes exactly where it left off.

This page shows what Vihren can do today and where it's going.

> **This is a direction, not a promise.** We don't give dates, and the order can
> change. Some items will ship sooner, some later, and some may change shape or
> drop off entirely as we learn from the people building on Vihren. We share it so
> you can see where things are headed and tell us what matters most to you.

---

## Available today

The foundation that everything else builds on.

- **Typed Temporal workflows with less boilerplate.** Write ordinary Go, mark it
  with simple annotations, and Vihren generates the wiring for you — with full
  type-safety, so mistakes are caught at compile time instead of in production.
- **Run everything in one process.** An embedded Temporal runtime lets you build,
  run, and demo locally in a single binary — no Docker, no separate servers, no
  setup ceremony.

---

## In active development

The durable agent toolkit — turning the foundation above into a complete way to
build agents.

- **A durable agent loop.** An agent that runs for days and survives crashes,
  deploys, and restarts without losing its place. If the process dies forty steps
  into a run, it picks up at step forty-one.
- **Tools as plain Go functions.** Give your agent abilities by writing ordinary,
  typed functions — no special framework objects to learn.
- **Human-in-the-loop, done right.** Pause an agent to bring in the right people —
  an approver, an expert, an end customer — and wait as long as it takes before
  continuing. One run can involve several different people at different stages.
- **Safe code execution, no extra infrastructure.** Run code the model writes
  inside a secure in-process sandbox compiled into your app — so you don't have to
  build and operate a fleet of containers or virtual machines just to run it
  safely.
- **Memory you can actually read.** The agent keeps a plain, human-readable record
  of what it did and decided. It doubles as your audit trail — no opaque database
  to inspect.
- **Optional quality loops.** Wrap an agent in a "check and improve" pass to raise
  the quality of its output, using a reviewer model or a simple rule.

---

## Planned

Making Vihren production-ready and easy to operate and observe.

- **One-step connection to managed Temporal.** Connect securely to a managed
  Temporal service (such as Temporal Cloud) without hand-wiring credentials.
- **Smarter error handling.** Automatically stop retrying on errors that will
  never succeed — like a bad API key — instead of retrying forever.
- **Cloud storage for large, long-running state.** Store big payloads from
  long-running agents in cloud object storage.
- **Ready-made test doubles.** Ship fakes so you can test your agents quickly and
  cheaply, without calling real AI providers or spending money.
- **Built-in observability and evaluation.** Emit traces using open standards
  ([OpenTelemetry](https://opentelemetry.io)) so you can see every model call,
  tool call, human step, and sandbox run in the tools you already use — plus a way
  to evaluate agent quality repeatably as you change things.
- **Tool interoperability.** Connect to the growing ecosystem of tools that speak
  the [Model Context Protocol (MCP)](https://modelcontextprotocol.io), and expose
  your own Vihren tools to other systems the same way.

---

## Exploring

Ideas we're interested in, further out, and likely shaped by what users ask for.

- **Agent-to-agent collaboration** using emerging open standards.
- **Web and browser automation** for agents that need to act on websites.
- **Advanced planning and reusable skills** for harder, higher-stakes decisions.
- **Local and offline models** for teams with strict data-residency requirements.

---

## What Vihren is not trying to be

Being clear about scope keeps Vihren focused and dependable.

- **Not a no-code builder.** Vihren is for engineers who build with code.
- **Not a Python or JavaScript framework.** Vihren is Go-native, by design.
- **Not a hosted service.** You run Vihren on your own infrastructure; we don't
  operate it for you.
- **Not another dashboard to adopt.** Vihren emits open, standard telemetry so you
  can use the observability tools you already have, rather than learning a new one.
- **Not a black box.** Vihren makes Temporal easier to use — it doesn't hide it.
  You always understand and control what's running.

---

## Tell us what matters to you

This roadmap moves based on real needs. If something here would unblock you — or
something you need is missing — we'd like to hear about it. Open an issue or get in
touch.
