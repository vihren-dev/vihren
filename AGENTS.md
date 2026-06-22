# AGENTS.md - Vihren Repository Rules

This repository is the from-scratch Go and Temporal rewrite of Slicer.

## Version Control

- This repository uses Jujutsu (`jj`) for version control.
- Prefer `jj` commands over `git` unless a task explicitly requires `git`.
- After each self-contained implementation step, update the paired status file
  in `./plans` and create a `jj` commit.

## Repository Split And Remotes

Vihren is maintained as one local Jujutsu repository with two logical release
lines and two GitHub remotes:

- `main` is the publishable upstream line. It tracks `main@public` and is
  pushed to the `main` branch of `git@github.com:vihren-dev/vihren.git`
  through the local Git remote named `public`.
- `main-private` is the downstream private line. It descends from `main`, tracks
  `main-private@private`, and is pushed to the `main-private` branch of
  `git@github.com:vihren-dev/vihren-private.git` through the local Git remote
  named `private`.

Only `main@public` and `main-private@private` are authoritative for release
work. Any remote branches named `public`, `internal`, or `main@private` are
legacy/stale helper refs unless a task explicitly says to clean them up.

Daily development normally happens on `main-private`. Public code, public docs,
and public examples must land on `main` first, then `main-private` must be
rebased on top of the updated `main` line. A separate workspace is preferred for
public work, for example:

```sh
jj workspace add ../vihren-public-work -r main
```

Use explicit refspec pushes when updating GitHub branches:

```sh
git push public refs/heads/main:refs/heads/main
git push private refs/heads/main-private:refs/heads/main-private
```

If `main-private` was intentionally rebased after a public change, update
`private/main-private` with `--force-with-lease` and the expected old
`private/main-private` commit. Do not use a blind force push.

## Public Vs Internal Files

The public v0.1 file allowlist is
`docs/internal/public-v0.1-manifest.md`. Treat that manifest as the source of
truth. The workflow document is `docs/internal/repository-structure.md`.

For v0.1, the public line is limited to:

- `cmd/vihren-gen`
- `internal/codegen`
- `internal/toolschema`, because `vihren-gen` imports it
- `platform/embeddedtemporal`
- `platform/embeddedtemporal/internal/litekit`
- `platform/blobref`
- `platform/toolcontract`
- `examples/codegenhello`
- root public support files listed in the manifest, including `README.md`,
  `go.mod`, `go.sum`, `Justfile`, and `docs/embedded-temporal.md`

Everything else is internal-only by default, including `plans/`, product docs,
architecture notes, internal workflow docs, experiments, non-v0.1 examples,
local scripts, IDE metadata, generated distribution artifacts, `AGENTS.md`, and
this `CLAUDE.md` symlink.

Before moving `main`, verify the public surface:

```sh
jj file list -r main | sort
go list ./...
go test ./... -timeout 120s
```

Also scan for internal-only paths and local-only dependencies before tagging or
pushing a release.

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
