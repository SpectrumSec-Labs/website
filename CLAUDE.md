# CLAUDE.md — working notes for this repo

## What this is
Public source for the SpectrumSec company website (cyber security consultancy).
Hugo static site, no npm. All changes via PR. Deploy config + infra docs are in
the separate PRIVATE repo `SpectrumSec-Labs/website-ops` (Caddy, systemd,
pull-deploy, PROJECT.md, vps-setup.md) — not here.

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

## Deployment
`release.yml` builds + signs (`ssh-keygen -Y`, key in `RELEASE_SIGNING_KEY`
secret) + publishes a GitHub Release. The VPS pulls and verifies it. All the
Caddy / systemd / pull-deploy config + `docs/PROJECT.md` + `docs/vps-setup.md`
live in the PRIVATE `SpectrumSec-Labs/website-ops` repo — the public key the
signature is checked against is `deploy/allowed_signers` there.

## Local tooling (Windows dev)
`bin/` holds hugo.exe, go/, node-v22.14.0-win-x64/, shellcheck.exe, caddy.exe —
all gitignored.

## Still to build (see website-ops/docs/PROJECT.md §8)
Phase 2-3 (front-matter schema enforcement, real copy); phase 7 (analytics —
lives in website-ops); phase 8 (favicons, OG image, lychee/pa11y/Lighthouse CI
gates, security-txt-bump.yml, template-generated robots/security.txt).
