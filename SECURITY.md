# Security

## Reporting a vulnerability

Email `security@spectrumsec.eu`, encrypted to our
[PGP key](https://spectrumsec.eu/pgp-key.txt) where possible. Full policy:
<https://spectrumsec.eu/legal/security-policy/> and
<https://spectrumsec.eu/.well-known/security.txt>.

We acknowledge reports within 3 working days.

## How this repository is protected

- **Branch protection** on `main`: PR required, required status checks
  (`build`, `markdown`, `newsfetch`), linear history, no force-push.
- **Pinned everything**: Hugo + Go in `.tool-versions`; GitHub Actions by commit
  SHA; container images by digest. Renovate proposes updates; Dependabot backstops
  security advisories.
- **Least-privilege CI**: default `permissions: contents: read`. Only the
  news-update job gets `contents: write` + `pull-requests: write`.
- **No build-time network access**: `config/_default/security.toml` sets
  `security.http.urls = ['none']`. News data is committed JSON, never fetched
  during a site build.
- **Secret scanning + push protection** enabled on the repository.
- **CodeQL** scans the Go news tool.

## Handling the automated news PRs

`update-news.yml` opens a PR that touches only `data/news.json`. It may auto-merge
**only** when the diff is limited to that file and all checks pass. Any PR from the
bot that touches another path must be reviewed and is never auto-merged. Review the
source allowlist in `config/feeds.yaml` before adding a feed, and record the
licence in the PR.

## Deployment trust chain

`release.yml` builds the site, signs the tarball with `ssh-keygen -Y` (private
key held only in an Actions secret), and publishes a GitHub Release. The
deployment target **verifies the signature and SHA-256 before serving**, then
swaps a symlink atomically. There is no inbound access from CI to the server and
no deploy credential in GitHub. Full detail is in the private ops repo.
