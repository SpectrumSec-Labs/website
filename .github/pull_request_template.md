<!-- Keep PRs small and single-purpose. Content PRs and infra PRs should not mix. -->

## What & why


## Type
- [ ] Content — blog / news / services copy
- [ ] Design / templates
- [ ] Build / CI / tooling
- [ ] Deploy / infrastructure
- [ ] Security

## Checklist
- [ ] `hugo --gc --minify` builds clean locally (or CI is green)
- [ ] No new external origins added to any page (if one is unavoidable, note it below — infra updates the CSP in the ops repo)
- [ ] No secrets, tokens, or internal hostnames in the diff
- [ ] Package / image / action versions are pinned (exact version or digest)
- [ ] Front matter matches the schema (content PRs) — see `tests/frontmatter-schema.mjs`

## CSP / header impact
<!-- "none", or list the directive change and the reason -->
none

## Screenshots
<!-- for visible changes; before/after if relevant -->
