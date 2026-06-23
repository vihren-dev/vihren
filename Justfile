set shell := ["bash", "-eu", "-o", "pipefail", "-c"]

test:
    @GOPATH="$PWD/.cache/go" GOCACHE="$PWD/.cache/go-build" go test ./... -timeout 120s

codegen ARGS="":
    @arg_string="{{ARGS}}"; arg_string="${arg_string#ARGS=}"; GOPATH="$PWD/.cache/go" GOCACHE="$PWD/.cache/go-build" go run ./cmd/vihren-gen $arg_string

clean-cache:
    @GOPATH="$PWD/.cache/go" GOCACHE="$PWD/.cache/go-build" go clean -cache -testcache
    @rm -rf .cache/go-tmp
    @mkdir -p .cache/go-tmp

run-codegenhello-embedded:
    @GOPATH="$PWD/.cache/go" GOCACHE="$PWD/.cache/go-build" go run ./examples/codegenhello/cmd/codegenhello-embedded

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
