# AGENTS.md - Vihren Repository Rules

This repository contains Vihren - an Agent development kit based on Go and Temporal

## Version Control

- This repository uses Jujutsu (`jj`) for version control.
- Prefer `jj` commands over `git` unless a task explicitly requires `git`.
- After each self-contained implementation step, update the paired status file
  in `./private/plans` and create a `jj` commit.

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

Daily development normally happens on top of `main-private`. Pushing
to any remote is usually done manually so you shouldn't do it unless
explicitly asked. Moving the `main` and `main-private` bookmarks is
also done manually. You can normally only operate on bookmarks that
denote the feature work you are currently involved in. Anything else
must be explicitly requested from you.


## Public Vs Internal Files

To understand which files are public and which private use the
`public` and `private` JJ filesets.

## Build And Commands

- Use the Go toolchain directly.
- Use `just` as the recorded developer command surface.
- All recurring development, test, lint, verification, and cleanup commands
  should be recorded in `Justfile`.

## Testing

- Every test command must include an explicit timeout.
- Prefer fast deterministic tests that use fakes over external
  services. However every fake must have a corresponding live test
  which confirms that the fake's behaviour is the same as the
  behaviour of the live service.
- Temporal workflow tests should use the Temporal Go SDK test environment unless
  a test is explicitly marked as a live integration test.
- External systems such as Redis, GCS, databases, VCS providers, and
  LLM providers should be mocked or faked until a work package
  explicitly adds a live adapter. But we need a live test to verify
  the behaviour of the mock. If it is not possible to make a live
  test, you should consider yourself BLOCKED and you need to STOP
  work.

## Demo Guidance

- `docs/private/how-to-write-a-good-demo.md` is the living guide for demos that show
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
- Keep files focused; split files before they become difficult to scan
  (soft limit <200 lines, hard limit <400 lines for source
  files. Documentation files can be as large as needed).
- Do not let workflow code instantiate provider, database, storage, Redis, VCS,
  or LLM clients directly.

## 18) Docstrings / comments

- Each class, struct, package, function, interface and so on should
  also contain a docstring which provides context for the reader to
  understand the motivation behind the object. Links to specs are
  welcome in the docstrings.

  Any public-facing component must have docstrings which serve as
  documentation to the external user. They should follow good
  practices in Go community for writing package documentation
