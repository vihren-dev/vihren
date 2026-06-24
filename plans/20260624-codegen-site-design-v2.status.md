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
- `docs/private/blog/20260622-hello-world.md`
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
- Review follow-up: moved blog archive and post rendering onto shared `_default` templates, keeping the same header, footer, typography, navigation, and `Blog` terminology as the homepage.
- Review follow-up: reduced homepage code volume to one before/after comparison and one compact authored API example, with generated worker/client detail moved out of the homepage narrative.
- Review follow-up: split status wording between the experimental Codegen API and the local-development embedded Temporal runtime.
- Review follow-up: revised the launch post to frame Vihren Codegen as a standalone useful tool instead of incomplete scaffolding for a future agent framework.
- Fixed-header follow-up: made the shared header sticky at the top of the viewport on every page and added scroll padding so section anchors remain visible below the header.
- Laptop-header follow-up: relaxed forced heading wrapping and kept the shared header brand/nav on one row for laptop and tablet widths, only stacking the nav on narrow mobile screens.
- Syntax-highlight follow-up: converted homepage code snippets to Hugo Chroma output and added a focused light syntax stylesheet that also styles blog Markdown code fences.

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
- Review follow-up `just site-build` passed.
- Review follow-up legacy-term scan found no `Vihren Systems`, `Quickstart`, `Writing`, `Notes`, `Generated surface`, `Worker and client boundary`, `production runtime`, `agent development kit`, `actual agentic parts`, or `Back to notes` in source or built output.
- Review follow-up positive scan verified the built site contains `Blog`, `Vihren / Codegen`, `Guide`, `Codegen API`, `Embedded Temporal`, `"ComposeGreeting"`, `Activity.ComposeGreeting`, `go generate ./examples/codegenhello`, and `Back to blog`.
- Review follow-up `go generate ./examples/codegenhello` passed and `jj diff -- examples/codegenhello` produced no output.
- Review follow-up targeted example tests passed with `-timeout 60s`.
- Review follow-up `just run-codegenhello-embedded` and `go run ./examples/temporalhello/cmd/temporalhello-embedded` both printed `Hello, Ada`.
- Review follow-up `just test` passed.
- Review follow-up Headless Chrome viewport QA passed across 360x800, 390x844, 430x932, 600x960, 820x1180, 1024x768, 1366x768, 1440x900, and 1920x1080; blog archive and post pages used the shared site shell with no document-level horizontal overflow.
- Fixed-header follow-up `just site-build` passed.
- Fixed-header follow-up Headless Chrome QA passed after scrolling `/`, `/blog/`, and `/blog/20260622-hello-world/` at 390x640 and 1366x500: the shared header stayed pinned to viewport top with no document-level horizontal overflow in the standard viewport matrix.
- Laptop-header follow-up `just site-build` passed.
- Laptop-header follow-up Headless Chrome QA passed across 360x800, 390x844, 430x932, 600x960, 820x1180, 1024x768, 1366x768, 1440x900, and 1920x1080: no document-level horizontal overflow, laptop/tablet widths kept the brand/nav on one row, and homepage heading line-count checks reported one line for each visible heading. Blog archive and post h1 headings also measured as one line with the shared header on one row.
- Syntax-highlight follow-up `just site-build` passed.
- Syntax-highlight follow-up verified generated homepage and blog HTML contain Chroma `highlight`, `chroma`, keyword, string, function, and comment spans, and the bundled CSS contains the new syntax palette.
- Syntax-highlight follow-up Headless Chrome QA passed across the existing viewport matrix with no document-level horizontal overflow; visual screenshot QA confirmed the homepage code comparison renders with visible highlighting.

## Next Step

- None. DONE.
