# SpectrumSec website

Public source for the SpectrumSec company site — services, a blog fed from
Markdown in Git, and an auto-aggregated threat-news section. Built with
[Hugo](https://gohugo.io); no npm, no runtime JS framework.

Deployment configuration and infrastructure docs live in a separate private
repo (`SpectrumSec-Labs/website-ops`).

## Local development

Hugo is pinned in [`.tool-versions`](.tool-versions). Install that exact version
(via [mise](https://mise.jdx.dev): `mise install`) or download it and verify the
checksum (the hash is in [`.github/workflows/ci.yml`](.github/workflows/ci.yml)).

```sh
hugo server -D            # http://localhost:1313, drafts on
hugo --gc --minify        # production build into ./public
```

There is **no npm / node_modules**. CSS and JS are hand-authored, concatenated,
minified, and SRI-hashed by Hugo. The news aggregator is a small vendored Go
module: `cd tools/newsfetch && go run -mod=vendor . -dry-run`.

## Layout

```
config/_default/   Hugo config. params.toml [company] holds the swappable
                   name + domain — nothing else hard-codes them.
config/feeds.yaml  Curated allowlist for the news aggregator.
content/           Markdown: services/, blog/<year>/<slug>/, news/, legal/.
data/              whyus.yaml, authors.yaml, news.json (generated).
assets/            css/, js/, fonts/ — processed by Hugo's pipeline.
layouts/           Templates. No third-party includes.
tools/newsfetch/   Go aggregator (RSS/Atom + CISA-KEV/NVD) → data/news.json.
```

## Security posture

- Strict CSP with **no `unsafe-inline`**; no external origins for scripts,
  styles, fonts, or images. All assets SRI-hashed.
- Raw HTML in Markdown disabled; the Hugo build has no network access.
- Toolchain pinned exactly; GitHub Actions pinned to commit SHAs; Renovate +
  Dependabot manage updates.
- Every change ships through a reviewed PR — see [SECURITY.md](SECURITY.md) and
  [docs/contributing.md](docs/contributing.md).

## Workflows

| Workflow | Trigger | Purpose |
|---|---|---|
| `ci` | PR / push | Hugo build + sanity, markdown lint, `newsfetch` (gofmt/vet/test/build) |
| `codeql` | `tools/**` changes, weekly | Go static analysis |
| `update-news` | cron 6h | Aggregate feeds → PR to `data/news.json` (auto-merge on green) |
| `release` | push to `main` | Build, sign (`ssh-keygen -Y`), publish a GitHub Release the VPS pulls |

## Rebranding

Edit `[company]` and `[brand]` in `config/_default/params.toml` (and `SITE_DOMAIN`
in the ops repo). Rebuild.
