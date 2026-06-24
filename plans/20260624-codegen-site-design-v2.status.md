# 2026-06-24 Codegen Site Design V2 Status

Issue: VIH-22
Plan: `plans/20260624-codegen-site-design-v2.plan.md`
Status: DONE

## Relevant Specs

- `docs/private/architecture/adr-20260624-vihren-codegen-design-v2.md`
- `private/assets/design-v2/vihren-codegen-redesign.html`
- `private/assets/design-v2/brand-spec.md`

## Relevant Source Files

- `private/site/hugo.toml`
- `private/site/layouts/_default/baseof.html`
- `private/site/layouts/_default/list.html`
- `private/site/layouts/_default/single.html`
- `private/site/layouts/index.html`
- `private/site/layouts/partials/head.html`
- `private/site/layouts/partials/header.html`
- `private/site/layouts/partials/footer.html`
- `private/site/layouts/partials/code-comparison.html`
- `private/site/layouts/partials/symbols.html`
- `private/site/assets/css/tokens.css`
- `private/site/assets/css/base.css`
- `private/site/assets/css/components.css`
- `private/site/assets/js/copy.js`
- `private/site/content/blog/_index.md`
- `examples/codegenhello/workflow.go`
- `examples/codegenhello/vihren.gen.go`
- `examples/codegenhello/cmd/codegenhello-embedded/main.go`
- `examples/temporalhello/workflow.go`
- `examples/temporalhello/cmd/temporalhello-embedded/main.go`

## Progress

- Created the plan/status pair before implementation work.
- Confirmed `examples/temporalhello` exists in this checkout and matches `../vihren-vih-19-public/examples/temporalhello`.
- Confirmed the starting working copy had no local changes.
- Rebuilt the Hugo site shell, header, footer, asset pipeline, homepage, code comparison partial, blog list, blog post template, light CSS, and copy-only JavaScript.
- Fixed the header seal `viewBox` found during mobile visual QA.
- Added footer long-link wrapping and visually verified the 360px footer.
- Compared rendered Hugo screenshots with `private/assets/design-v2/vihren-codegen-redesign.html` at 390x844 and 1440x900. The Hugo page keeps the design-v2 layout posture and wraps canonical content that the static prototype clipped.

## Checks

- `jj status` passed before planning and showed no local changes.
- `diff -qr examples/temporalhello ../vihren-vih-19-public/examples/temporalhello` produced no output.
- `just site-build` passed after the Hugo redesign.
- `go generate ./examples/codegenhello` passed after the Hugo redesign.
- `jj diff -- examples/codegenhello` produced no output after generation.
- Old-theme scan found no Google Fonts, nav toggle, tabs, counters, reveal hooks, dark Hugo highlight style, radial gradients, or backdrop blur in `private/site/layouts`, `private/site/assets`, or `private/site/public`.
- Visual QA found the header seal was missing its SVG `viewBox`, which made the mobile header too tall. Fixed before final QA.
- `go test ./examples/codegenhello ./examples/codegenhello/cmd/codegenhello-embedded ./examples/temporalhello ./examples/temporalhello/cmd/temporalhello-embedded -timeout 60s` passed.
- `just run-codegenhello-embedded` printed `Hello, Ada`.
- `go run ./examples/temporalhello/cmd/temporalhello-embedded` printed `Hello, Ada`.
- `just test` passed with the repository's `-timeout 120s` setting.
- Final `just site-build` passed after the header/footer visual fixes.
- `just site-serve` ran successfully at `http://127.0.0.1:1313/`.
- Headless Chrome viewport QA passed for 360x800, 390x844, 430x932, 600x960, 820x1180, 1024x768, 1366x768, 1440x900, and 1920x1080: no document-level horizontal scroll, minimum code font size 14px, header/nav inside viewport, install command and metadata above the fold, and copy fallback selecting the install command with an `aria-live` message.
- Verified generated `/blog/`, `/blog/index.xml`, `/sitemap.xml`, and `CNAME` exist after build.

## Next Step

- None. DONE.
