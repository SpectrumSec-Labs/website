# newsfetch

Aggregates cyber-security advisories from the government / CERT allowlist in
`config/feeds.yaml` into `data/news.json`, which Hugo renders on `/news/` and as
an aggregated RSS feed.

## Run

```sh
cd tools/newsfetch
go run -mod=vendor . -feeds ../../config/feeds.yaml -out ../../data/news.json
```

Flags: `-feeds`, `-out`, `-now <RFC3339>` (freeze time, for tests), `-dry-run`
(print result to stdout, do not write).

Exit status is non-zero and the output file is left untouched if **every** source
fails. If the item set is unchanged the file is not rewritten (so the workflow
opens no PR).

## Design

- One dependency: `gopkg.in/yaml.v3`, vendored. Everything else is stdlib.
- HTTPS only, 10s timeout, 3 MiB response cap, ≤5 redirects, explicit UA.
- RSS and Atom are both parsed from XML (the `kind` hint only separates XML from
  JSON). JSON sources are matched by `id` — `cisa-kev` and `nvd-recent` each
  have a dedicated parser.
- Summaries are reduced to plain text. Output is deterministic: items sorted
  newest-first with id tie-break, tags sorted (`KEV` pinned first), 2-space
  indent, trailing newline.
- Retention (`retention_days`, `max_items` in `config/feeds.yaml`) is enforced on
  every run; existing items keep their original publish date and id.

## Tests

```sh
go test -mod=vendor ./...
```

Fixtures in `testdata/` cover RSS, Atom, and the KEV JSON shape; the end-to-end
test runs the whole flow against an in-process TLS server.

## Adding a source

Edit `config/feeds.yaml`. For an RSS/Atom feed that is usually all that is
needed. For a new JSON source, add a parser in `parse.go` keyed by the source
`id`. Record the licence in the PR — the allowlist is government / open data only.
