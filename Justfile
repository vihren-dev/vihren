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
