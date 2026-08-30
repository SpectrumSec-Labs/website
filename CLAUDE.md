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

## newsfetch (tools/newsfetch)
- Go 1.27, single vendored dep gopkg.in/yaml.v3. Build/test with `-mod=vendor`.
- Local dev: Go at bin/go/ (gitignored). `export GOROOT=.../bin/go` then `go.exe`.
- RSS+Atom via encoding/xml; CISA-KEV + NVD via JSON parsers keyed on source id.
- CISA feeds use RFC822 weekday + 2-digit year ("Thu, 27 Aug 26 ..."); date
  layouts live in sanitize.go.
- Deterministic output; unchanged item set => no write => no PR.
- gofmt check in CI must exclude vendor/.

## Still to build (see README status table + docs/PROJECT.md §8)
Phase 2-3 (front-matter schema enforcement, real copy); phase 6 (release.yml,
pull-deploy.sh, systemd units, docs/vps-setup.md); phase 7 (deploy/analytics/
Quadlet units, GoAccess); phase 8 (favicons, OG image, lychee/pa11y/Lighthouse
CI gates, security-txt-bump.yml, template-generated robots/security.txt).
