# 2026-06-24 Codegen Site Design V2 Status

Issue: VIH-22
Plan: `plans/20260624-codegen-site-design-v2.plan.md`
Status: In progress

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

## Checks

- `jj status` passed before planning and showed no local changes.
- `diff -qr examples/temporalhello ../vihren-vih-19-public/examples/temporalhello` produced no output.

## Next Step

- Commit the planning files, then rebuild the site shell and asset pipeline.
