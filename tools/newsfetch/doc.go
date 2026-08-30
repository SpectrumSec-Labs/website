// Command newsfetch aggregates cyber-security advisories from a curated allowlist
// of government / CERT feeds (config/feeds.yaml) into data/news.json, which the
// Hugo site renders on /news/ and in an aggregated RSS feed.
//
// It is run by .github/workflows/update-news.yml on a schedule. The workflow
// opens a pull request only when the item set actually changes; when the diff is
// limited to data/news.json and CI is green it may auto-merge.
//
// Design constraints:
//   - Single dependency (gopkg.in/yaml.v3); everything else is the standard
//     library. Vendored so builds are hermetic and auditable.
//   - Network egress is bounded: HTTPS only, short timeout, response size cap,
//     explicit User-Agent, capped redirects.
//   - Summaries are reduced to plain text. Output is deterministic so unchanged
//     runs produce no diff.
//   - A run where every source fails never overwrites existing data.
package main
