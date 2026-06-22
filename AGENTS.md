# AGENTS.md - Vihren Repository Rules

This repository is the from-scratch Go and Temporal rewrite of Slicer.

## Version Control

- This repository uses Jujutsu (`jj`) for version control.
- Prefer `jj` commands over `git` unless a task explicitly requires `git`.
- After each self-contained implementation step, update the paired status file
  in `./plans` and create a `jj` commit.

## Build And Commands

- Use the Go toolchain directly.
- Use `just` as the recorded developer command surface.
- All recurring development, test, lint, verification, and cleanup commands
  should be recorded in `Justfile`.
- Do not add Bazel files, Bazel dependencies, or Bazel commands to this repo.

## Testing

- Every test command must include an explicit timeout.
- Prefer fast deterministic tests that use fakes over external services.
- Temporal workflow tests should use the Temporal Go SDK test environment unless
  a test is explicitly marked as a live integration test.
- External systems such as Redis, GCS, databases, VCS providers, and LLM
  providers should be mocked or faked until a work package explicitly adds a
  live adapter.

## Demo Guidance

- `docs/how-to-write-a-good-demo.md` is the living guide for demos that show
  human-reviewable business value, not only implementation correctness.
- Keep the document updated from reviewer feedback after demos are run. The
  file starts empty until that feedback exists.

## Go Code Quality

- Keep package surfaces small and typed.
- Prefer functional-core logic behind imperative-shell adapters.
- Public packages, exported types, exported functions, and exported methods must
  have comments that explain their purpose.
- Avoid `any` unless the boundary genuinely requires untyped data and the reason
  is documented.
- Keep files focused; split files before they become difficult to scan.
- Do not let workflow code instantiate provider, database, storage, Redis, VCS,
  or LLM clients directly.

## Initial Architecture Direction

- `cmd/vihren-worker` should remain a thin process entry point.
- `internal/runtime` owns worker bootstrapping and dependency wiring.
- `internal/config` owns typed configuration loading.
- `internal/policy` owns versioned policy decisions.
- `internal/observability` owns logging, metrics, tracing, and event interfaces.
- `internal/temporaltest` owns shared Temporal SDK test helpers.
- Product behavior should be added behind explicit interfaces and contract
  tests, not broad shared clients.

## 18) Docstrings / comments

- Each class, struct, package, function, interface and so on should
  also contain a docstring which provides context for the reader to
  understand the motivation behind the object. Links to specs are
  welcome in the docstrings
