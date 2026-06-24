# 2026-06-24 Codegen Site Design V2 Plan

Issue: VIH-22

Spec sources:

- `docs/private/architecture/adr-20260624-vihren-codegen-design-v2.md`
- `private/assets/design-v2/vihren-codegen-redesign.html`
- `private/assets/design-v2/brand-spec.md`
- `examples/codegenhello/workflow.go`
- `examples/codegenhello/vihren.gen.go`
- `examples/codegenhello/cmd/codegenhello-embedded/main.go`
- `examples/temporalhello/workflow.go`
- `examples/temporalhello/cmd/temporalhello-embedded/main.go`

## Steps

1. Planning and source audit
   - Create this plan/status pair and verify the main checkout already has the VIH-19 `examples/temporalhello` counterpart.
   - Tests: `jj status`; `diff -qr examples/temporalhello ../vihren-vih-19-public/examples/temporalhello`.

2. Site shell and asset pipeline
   - Replace Hugo head/base/header/footer wiring with the light Codegen document shell.
   - Split CSS into `tokens.css`, `base.css`, and `components.css`.
   - Replace `site.js` with copy-only behavior in `copy.js`.
   - Tests: `just site-build`; inspect generated HTML for no Google Fonts and no old `site.js`/`site.css` references.

3. Homepage and snippets
   - Replace the homepage with the exported order: hero, guide, smallest useful example, status, notes.
   - Add `partials/code-comparison.html` for before/after snippets copied from canonical source excerpts.
   - Keep shortened snippets truthful and mark omissions where needed.
   - Tests: `just site-build`; `go generate ./examples/codegenhello`; `jj diff` to confirm no unintended generation diff.

4. Blog document styling
   - Restyle `/blog/` and post pages to the same light technical document system without removing routes, RSS, sitemap, or Markdown sources.
   - Update `private/site/content/blog/_index.md` for Codegen positioning.
   - Tests: `just site-build`; inspect generated `/blog/`, launch post, RSS, and sitemap outputs.

5. Required Go/example checks
   - Verify the website snippets still correspond to compiling examples.
   - Tests: `go test ./examples/codegenhello ./examples/codegenhello/cmd/codegenhello-embedded ./examples/temporalhello ./examples/temporalhello/cmd/temporalhello-embedded -timeout 60s`; `just run-codegenhello-embedded`; `go run ./examples/temporalhello/cmd/temporalhello-embedded`; confirm both embedded runs print `Hello, Ada`.

6. Visual QA
   - Run `just site-serve` and compare the built homepage against `private/assets/design-v2/vihren-codegen-redesign.html`.
   - Check 360x800, 390x844, 430x932, 600x960, 820x1180, 1024x768, 1366x768, 1440x900, and 1920x1080.
   - Tests: no horizontal scroll, no text overlap, legible code blocks, clean header/nav wrapping, visible focus states, and above-the-fold install/status content.

7. Practical full-suite check
   - Run `just test` if practical after targeted checks pass.
   - Tests: `just test` with the repository's built-in `-timeout 120s`.

## Completion Criteria

- The Hugo site matches the design-v2 light technical page direction.
- Homepage snippets are based on canonical example sources and the before/after comparison is truthful.
- Blog routes remain present and styled consistently.
- Required checks and visual QA are run or any blocker is reported.
- This plan and its paired status file are marked `DONE`.
