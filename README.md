# SpectrumSec website

Company site for SpectrumSec — services, blog (Markdown in Git), and an
auto-aggregated threat-news section. Built with [Hugo](https://gohugo.io),
served by Caddy on a self-hosted VPS, deployed via signed release + pull.

## Status

| Phase | Scope | State |
|------:|-------|-------|
| 0 | Repo scaffold, config, CI skeleton | done |
| 1 | Design system, base templates, home + Why-Us, 6 demo services, sample blog, news rendering | done |
| 4–5 | `newsfetch` Go tool + `update-news.yml` scheduled PR workflow + CodeQL | done |
| 2–3 | Service/blog polish, front-matter schema enforcement | pending |
| 6 | Deploy: signed release, VPS puller, systemd units, `docs/vps-setup.md` | Caddyfile drafted |
| 7 | Umami (Podman pod) + GoAccess | pending |
| 8 | SEO/JSON-LD polish, a11y + Lighthouse gates, favicons, legal review | pending |

Full write-up: [docs/PROJECT.md](docs/PROJECT.md).

## Local development

Hugo is pinned in [`.tool-versions`](.tool-versions). Install that exact version
(via [mise](https://mise.jdx.dev): `mise install`), or download it and verify the
checksum (see [`.github/workflows/ci.yml`](.github/workflows/ci.yml) for the hash).

```sh
hugo server -D            # http://localhost:1313, drafts on
hugo --gc --minify        # production build into ./public
```

There is **no npm / node_modules**. CSS and JS are hand-authored, concatenated,
minified, and SRI-hashed by Hugo.

## Layout

```
config/_default/   Hugo config, split by concern. params.toml holds the
                   swappable company name + domain — nothing else hard-codes them.
config/feeds.yaml  Curated allowlist for the news aggregator.
content/           Markdown. services/, blog/<year>/<slug>/, news/, legal/.
data/              whyus.yaml, authors.yaml, news.json (news.json is generated).
assets/            css/, js/, fonts/ — processed by Hugo's pipeline.
layouts/           Templates. No third-party includes.
deploy/            Caddyfile + (coming) systemd units and pull-deploy script.
tools/newsfetch/   (coming) Go aggregator.
tests/             CI checks.
docs/              Contributor and ops docs.
```

## Security posture (summary)

- Strict CSP with **no `unsafe-inline`** anywhere on the public site; no external
  origins for scripts, styles, fonts, or images.
- All response headers defined in [`deploy/Caddyfile`](deploy/Caddyfile),
  CI-checked by [`tests/headers.test.sh`](tests/headers.test.sh).
- Toolchain pinned to exact versions; GitHub Actions pinned to commit SHAs;
  Renovate + Dependabot manage updates.
- Every change ships through a reviewed PR. See [SECURITY.md](SECURITY.md) and
  [docs/contributing.md](docs/contributing.md).

## Rebranding

Edit `[company]` and `[brand]` in `config/_default/params.toml` and the
`SITE_DOMAIN` in the deploy environment. Rebuild.
