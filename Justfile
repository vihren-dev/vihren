set shell := ["bash", "-eu", "-o", "pipefail", "-c"]

hugo-version := "v0.150.0"
hugo-bin := "private/site/.bin/hugo"

test:
    @GOPATH="$PWD/.cache/go" GOCACHE="$PWD/.cache/go-build" go test ./... -timeout 60s

site-hugo:
    @mkdir -p private/site/.bin .cache/go .cache/go-build .cache/go-tmp
    @if [ ! -x "{{hugo-bin}}" ] || ! "{{hugo-bin}}" version | grep -q "{{hugo-version}}"; then \
      GOBIN="$PWD/private/site/.bin" GOPATH="$PWD/.cache/go" GOCACHE="$PWD/.cache/go-build" GOTMPDIR="$PWD/.cache/go-tmp" \
        go install github.com/gohugoio/hugo@{{hugo-version}}; \
    fi

site-build: site-hugo
    @cd private/site; .bin/hugo --destination public --cleanDestinationDir --minify
    @test -f private/site/public/CNAME

site-serve: site-hugo
    @cd private/site; .bin/hugo server --destination public --bind 127.0.0.1 --baseURL http://127.0.0.1:1313/ --disableFastRender

codegen ARGS="":
    @arg_string="{{ARGS}}"; arg_string="${arg_string#ARGS=}"; GOPATH="$PWD/.cache/go" GOCACHE="$PWD/.cache/go-build" go run ./cmd/vihren-gen $arg_string

experiment-pp-publish-main-test:
    @GOPATH="$PWD/.cache/go" GOCACHE="$PWD/.cache/go-build" go test ./experiments/publicprivate/... -timeout 30s

clean-cache:
    @GOPATH="$PWD/.cache/go" GOCACHE="$PWD/.cache/go-build" go clean -cache -testcache
    @rm -rf .cache/go-tmp
    @mkdir -p .cache/go-tmp

valueflow-validate:
    @GOPATH="$PWD/.cache/go" GOCACHE="$PWD/.cache/go-build" go run ./cmd/vihren-valueflow validate --source-revision initial-source --value-revision initial-value

valueflow-answer ARGS="":
    @arg_string="{{ARGS}}"; arg_string="${arg_string#ARGS=}"; GOPATH="$PWD/.cache/go" GOCACHE="$PWD/.cache/go-build" go run ./cmd/vihren-valueflow answer $arg_string --source-revision initial-source --value-revision initial-value

valueflow-extract ARGS="":
    @arg_string="{{ARGS}}"; arg_string="${arg_string#ARGS=}"; scripts/valueflow-extract-facts $arg_string

temporal-doctor:
    @scripts/temporal-local-dev doctor

temporal-start:
    @scripts/temporal-local-dev start

temporal-ui-url:
    @scripts/temporal-local-dev ui-url

temporal-state-dir:
    @scripts/temporal-local-dev state-dir

temporal-live-check:
    @scripts/temporal-local-dev check

temporal-reset:
    @scripts/temporal-local-dev reset

worker-blah:
    @VIHREN_RUN_MODE="${VIHREN_RUN_MODE:-live}" GOPATH="$PWD/.cache/go" GOCACHE="$PWD/.cache/go-build" go run ./examples/blah/cmd/blah-worker

worker-humanwebchat ARGS="":
    @arg_string='{{ARGS}}'; arg_string="${arg_string#ARGS=}"; eval "VIHREN_RUN_MODE=\"${VIHREN_RUN_MODE:-live}\" GOPATH=\"$PWD/.cache/go\" GOCACHE=\"$PWD/.cache/go-build\" go run ./examples/humanwebchat/cmd/humanwebchat-worker $arg_string"

worker-codegenhello ARGS="":
    @arg_string='{{ARGS}}'; arg_string="${arg_string#ARGS=}"; eval "GOPATH=\"$PWD/.cache/go\" GOCACHE=\"$PWD/.cache/go-build\" go run ./examples/codegenhello/cmd/codegenhello-worker $arg_string"

run-codegenhello-embedded:
    @GOPATH="$PWD/.cache/go" GOCACHE="$PWD/.cache/go-build" go run ./examples/codegenhello/cmd/codegenhello-embedded

worker-audit-trial:
    @GOPATH="$PWD/.cache/go" GOCACHE="$PWD/.cache/go-build" go run ./cmd/vihren-audit-worker

workflow-blah-start ARGS="":
    @arg_string="{{ARGS}}"; arg_string="${arg_string#ARGS=}"; GOPATH="$PWD/.cache/go" GOCACHE="$PWD/.cache/go-build" go run ./examples/blah/cmd/blah-start $arg_string

workflow-humanwebchat-start ARGS="":
    @arg_string='{{ARGS}}'; arg_string="${arg_string#ARGS=}"; eval "GOPATH=\"$PWD/.cache/go\" GOCACHE=\"$PWD/.cache/go-build\" go run ./examples/humanwebchat/cmd/humanwebchat-start $arg_string"

workflow-codegenhello-start ARGS="":
    @arg_string='{{ARGS}}'; arg_string="${arg_string#ARGS=}"; eval "GOPATH=\"$PWD/.cache/go\" GOCACHE=\"$PWD/.cache/go-build\" go run ./examples/codegenhello/cmd/codegenhello-start $arg_string"

workflow-audit-trial-start ARGS="":
    @arg_string='{{ARGS}}'; arg_string="${arg_string#ARGS=}"; eval "GOPATH=\"$PWD/.cache/go\" GOCACHE=\"$PWD/.cache/go-build\" go run ./cmd/vihren-audit-start $arg_string"

wasmbusybox-materialize-assets:
    @mkdir -p examples/wasmbusybox/assets
    @out_path="$(nix build --no-link --print-out-paths ../wasm-sandbox#busybox-wasi-minimal)"; chmod u+w examples/wasmbusybox/assets/busybox.wasm examples/wasmbusybox/assets/wasm-sandbox.package.yaml 2>/dev/null || true; cp "$out_path/bin/busybox.wasm" examples/wasmbusybox/assets/busybox.wasm; cp "$out_path/wasm-sandbox.package.yaml" examples/wasmbusybox/assets/wasm-sandbox.package.yaml; chmod u+w examples/wasmbusybox/assets/busybox.wasm examples/wasmbusybox/assets/wasm-sandbox.package.yaml

wasmbusybox-build-all: wasmbusybox-materialize-assets
    @mkdir -p dist .cache/go-tmp
    @export GOPATH="$PWD/.cache/go"; \
      export GOCACHE="$PWD/.cache/go-build"; \
      export GOTMPDIR="$PWD/.cache/go-tmp"; \
      export CGO_ENABLED=0; \
      pkg="./examples/wasmbusybox/cmd/wasmbusybox-demo"; \
      tags="wasmbusybox_embed"; \
      ldflags="-s -w"; \
      for target in darwin/arm64 darwin/amd64 linux/amd64 linux/arm64 windows/amd64 windows/arm64; do \
        os="${target%/*}"; \
        arch="${target#*/}"; \
        ext=""; \
        if [ "$os" = "windows" ]; then ext=".exe"; fi; \
        out="dist/wasmbusybox-demo-$os-$arch$ext"; \
        echo "building $out"; \
        GOOS="$os" GOARCH="$arch" go build -p=1 -trimpath -tags "$tags" -ldflags "$ldflags" -o "$out" "$pkg"; \
      done

human-tui ARGS="":
    @arg_string='{{ARGS}}'; arg_string="${arg_string#ARGS=}"; eval "GOPATH=\"$PWD/.cache/go\" GOCACHE=\"$PWD/.cache/go-build\" go run ./cmd/vihren-human-tui $arg_string"

demo-audit-trial-scaffold:
    @GOPATH="$PWD/.cache/go" GOCACHE="$PWD/.cache/go-build" go run ./cmd/vihren-audit-start --dry-run --business-value "Shape-check the audit workflow." --repository-root .

demo-codegenhello:
    @bash scripts/codegenhello-demo

demo-step-1:
    @scripts/valueflow-demo step1

demo-step-2:
    @scripts/valueflow-demo step2

demo-step-3:
    @scripts/valueflow-demo step3

demo-live-codex-extraction INPUT="internal/valueflow/testdata/source-note.md":
    @VALUEFLOW_LIVE_INPUT="{{INPUT}}" scripts/valueflow-live-codex-demo

demo-step-4:
    @scripts/valueflow-demo step4

demo-step-5:
    @scripts/valueflow-demo step5

demo-step-6:
    @scripts/valueflow-demo step6

demo-step-7:
    @scripts/valueflow-demo step7

demo-final:
    @scripts/valueflow-demo final

demo-final-live INPUT="internal/valueflow/testdata/source-note.md":
    @VALUEFLOW_LIVE_INPUT="{{INPUT}}" scripts/valueflow-live-codex-demo
    @scripts/valueflow-demo final


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
