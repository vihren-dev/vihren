set shell := ["bash", "-eu", "-o", "pipefail", "-c"]

# Public developer command surface for Vihren. Private recipes live in an
# unpublished overlay, pulled in by the optional import at the bottom of this
# file; that overlay is absent in a public clone and is silently skipped there.

test:
    @go test ./... -timeout 120s

codegen ARGS="":
    @arg_string="{{ARGS}}"; arg_string="${arg_string#ARGS=}"; go run ./cmd/vihren-gen $arg_string

clean-cache:
    @go clean -cache -testcache
    @rm -rf .cache/go-tmp
    @mkdir -p .cache/go-tmp

worker-codegenhello ARGS="":
    @arg_string='{{ARGS}}'; arg_string="${arg_string#ARGS=}"; eval "go run ./examples/codegenhello/cmd/codegenhello-worker $arg_string"

run-codegenhello-embedded:
    @go run ./examples/codegenhello/cmd/codegenhello-embedded

workflow-codegenhello-start ARGS="":
    @arg_string='{{ARGS}}'; arg_string="${arg_string#ARGS=}"; eval "go run ./examples/codegenhello/cmd/codegenhello-start $arg_string"

# Install the tracked public/private workflow config into jj's per-repo config.
# Safe: backs up any existing non-workflow repo config before overwriting.
setup-jj:
    #!/usr/bin/env bash
    set -euo pipefail
    target="$(jj config path --repo)"
    mkdir -p "$(dirname "$target")"
    if [ -e "$target" ] && ! head -1 "$target" | grep -q 'vihren public/private workflow'; then
      backup="$target.bak.$(date +%s)"; cp "$target" "$backup"
      echo "backed up existing repo config -> $backup"
    fi
    cp "$PWD/jj/workflow.toml" "$target"
    echo "installed jj/workflow.toml -> $target"

# Private recipe overlay (site, valueflow, temporal-dev, audit, human-tui,
# wasmbusybox, demos, experiments). Optional import: silently skipped when the
# overlay is not present, e.g. in a public clone.
import? 'private/Justfile.just'
