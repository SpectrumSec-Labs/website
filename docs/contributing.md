# Contributing

Every change ships through a pull request. No direct pushes to `main`.

## Before you start

Install the pinned Hugo version (`mise install`, or download + verify per
`.github/workflows/ci.yml`). Confirm `hugo --gc --minify` builds clean.

## Writing a blog post

1. Create `content/blog/<year>/<slug>/index.md` as a page bundle. Put images in
   the same folder.
2. Front matter (all required unless noted):

   ```yaml
   ---
   title: "…"
   slug: "short-stable-slug"
   date: 2026-09-01T09:00:00+03:00
   author: "spectrumsec"        # must be a key in data/authors.yaml
   summary: "50–200 characters."
   tags: ["at-least-one"]
   draft: true                  # flip to false when ready
   toc: true                    # optional
   ---
   ```

3. Raw HTML in Markdown is disabled by design. Use shortcodes if you need
   structure beyond CommonMark.
4. Open a PR labelled **Content**. `@spectrumsec/editors` reviews.

## Changing services

Edit `content/services/<slug>.md`. Keep `weight` spaced by 10s. `icon` must be a
symbol id defined in `layouts/partials/icons-sprite.html`.

## Design / template changes

- Never add an inline `style` attribute or `<script>` block — the CSP forbids
  them. Add a class to `assets/css/utilities.css` instead.
- No external origins. No web fonts from a CDN, no third-party scripts.
- If a change genuinely needs a CSP change, note it in the PR's "CSP / header
  impact" section — the header policy is enforced by Caddy in the private ops
  repo and someone from infra applies the matching change there.

## Review

`CODEOWNERS` routes `content/blog/` and `content/news/` to editors and the rest
to infra. `main` is protected: PR required, CI must pass (`build`, `markdown`,
`newsfetch`), linear history, no force-push. The automated `data/news.json` PR
auto-merges once CI is green.
