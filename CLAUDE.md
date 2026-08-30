# CLAUDE.md — working notes for this repo

## What this is
Company website for SpectrumSec (cyber security consultancy). Hugo static site,
self-hosted on a VPS behind Caddy, deployed by signed release + VPS-side pull.
All changes via PR.

## Hard rules
- **No npm / node_modules.** CSS and JS are hand-authored under `assets/`,
  bundled by Hugo (`resources.Concat | minify | fingerprint "sha384"`).
- **No inline `style=` or `<script>`.** Public CSP is `script-src 'self'` /
  `style-src 'self'` with no `unsafe-inline`. Use `assets/css/utilities.css`.
- **No external origins** on the public site — no CDN fonts, scripts, images.
- **Pin versions.** Hugo/Go in `.tool-versions`; Actions by SHA; images by digest.
- Company name/domain come only from `config/_default/params.toml` `[company]`.
- Raw HTML in Markdown is disabled (`markup.goldmark.renderer.unsafe = false`).
- Hugo build must be offline (`security.toml` → `security.http.urls = ['none']`).

## Hugo API notes (v0.165)
- Use `hugo.Data`, not `site.Data` (deprecated).
- Use `site.Language.Locale`, not `.LanguageCode`.
- JSON-LD inside `<script type="application/ld+json">` needs `| jsonify | safeJS`
  or html/template double-encodes it to a quoted string.
- CI builds with `--panicOnWarning`; deprecations fail the build.

## Local build
```
./bin/hugo.exe --gc --minify --panicOnWarning     # bin/ is gitignored
```

## Layout map
- `layouts/_default/` baseof, home, list, single, 404
- `layouts/services/single.html`, `layouts/contact/single.html`,
  `layouts/news/{list.html,rss.xml}`
- `layouts/partials/` head, header, footer, schema, analytics, icon(+sprite),
  pagination
- News is data-driven from `data/news.json` (schema v1: `{generated, items[]}`),
  not Markdown pages.

## Still to build (see README status table)
newsfetch Go tool + workflow; deploy scripts + systemd units + docs/vps-setup.md;
Umami Quadlet units under deploy/analytics/; CI link-check/a11y/Lighthouse gates;
favicons + OG default image; front-matter schema enforcement in templates.
