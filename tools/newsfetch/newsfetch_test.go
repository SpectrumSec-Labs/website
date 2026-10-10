package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

var refNow = time.Date(2026, 8, 26, 0, 0, 0, 0, time.UTC)

func TestToPlainText(t *testing.T) {
	in := "<p>An attacker can <strong>bypass&nbsp;auth</strong>.</p>\n<p>Patch.</p>"
	got := toPlainText(in, 0)
	want := "An attacker can bypass auth . Patch."
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
	if toPlainText("abcdefghij", 5) != "abcde…" {
		t.Fatalf("truncation: got %q", toPlainText("abcdefghij", 5))
	}
}

func TestCanonicalURL(t *testing.T) {
	cases := map[string]string{
		"https://X.example.gov/a?utm_source=rss&ref=feed&id=7#frag": "https://x.example.gov/a?id=7",
		"https://e.gov/p/": "https://e.gov/p/",
		"not a url":        "not a url",
	}
	for in, want := range cases {
		if got := canonicalURL(in); got != want {
			t.Errorf("canonicalURL(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestParseDate(t *testing.T) {
	for _, s := range []string{
		"Mon, 24 Aug 2026 09:00:00 +0000",
		"2026-08-24T09:00:00Z",
		"2026-08-24",
		"2026-08-24T09:00:00.000",
	} {
		if _, ok := parseDate(s); !ok {
			t.Errorf("parseDate(%q) failed", s)
		}
	}
	if _, ok := parseDate("last tuesday"); ok {
		t.Error("expected parseDate to fail on garbage")
	}
}

func TestParseXML_RSS(t *testing.T) {
	maxPerSource = 10
	body := readFixture(t, "rss.xml")
	src := Source{ID: "example-cert", Name: "Example CERT", Kind: "rss", Tag: "advisory", Enabled: true}
	items, err := parseSource(src, body, refNow)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatalf("want 2 items, got %d", len(items))
	}
	first := items[0] // newest first
	if !strings.Contains(first.Title, "Example Router") {
		t.Errorf("unexpected title %q", first.Title)
	}
	if strings.Contains(first.Summary, "<") {
		t.Errorf("summary still has HTML: %q", first.Summary)
	}
	if first.URL != "https://cert.example.gov/advisories/2026-001" {
		t.Errorf("tracking params not stripped: %q", first.URL)
	}
	if !hasTag(first.Tags, "advisory") || !hasTag(first.Tags, "CVE-2026-12345") {
		t.Errorf("tags = %v", first.Tags)
	}
}

func TestParseXML_Atom(t *testing.T) {
	maxPerSource = 10
	body := readFixture(t, "atom.xml")
	src := Source{ID: "example-ncc", Name: "Example NCC", Kind: "atom", Tag: "guidance", Enabled: true}
	items, err := parseSource(src, body, refNow)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatalf("want 2, got %d", len(items))
	}
	if items[0].URL != "https://ncc.example.gov/reports/q3-loaders" {
		t.Errorf("expected alternate link, got %q", items[0].URL)
	}
}

func TestParseKEV(t *testing.T) {
	maxPerSource = 10
	body := readFixture(t, "kev.json")
	src := Source{ID: "cisa-kev", Name: "Example KEV", Kind: "json", Enabled: true}
	items, err := parseSource(src, body, refNow)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 3 {
		t.Fatalf("want 3, got %d", len(items))
	}
	top := items[0]
	if !hasTag(top.Tags, "KEV") || !hasTag(top.Tags, "CVE-2026-2222") || !hasTag(top.Tags, "ransomware") {
		t.Errorf("tags = %v", top.Tags)
	}
	if top.Tags[0] != "KEV" {
		t.Errorf("KEV should sort first: %v", top.Tags)
	}
	if top.URL != "https://nvd.nist.gov/vuln/detail/CVE-2026-2222" {
		t.Errorf("url = %q", top.URL)
	}
}

func TestMergeAndPrune(t *testing.T) {
	existing := []Item{
		{ID: "a", Title: "A", URL: "https://e.gov/a", Published: "2026-08-01T00:00:00Z",
			publishedAt: mustDate("2026-08-01T00:00:00Z"), Tags: []string{"advisory"}},
		{ID: "old", Title: "Old", URL: "https://e.gov/old", Published: "2026-01-01T00:00:00Z",
			publishedAt: mustDate("2026-01-01T00:00:00Z"), Tags: []string{"advisory"}},
	}
	fresh := []Item{
		{ID: "a", Title: "A", URL: "https://e.gov/a", Published: "2026-08-05T00:00:00Z",
			publishedAt: mustDate("2026-08-05T00:00:00Z"), Summary: "longer summary", Tags: []string{"KEV"}},
		{ID: "b", Title: "B", URL: "https://e.gov/b", Published: "2026-08-20T00:00:00Z",
			publishedAt: mustDate("2026-08-20T00:00:00Z"), Tags: []string{"guidance"}},
	}
	merged := pruneAndSort(mergeItems(existing, fresh), refNow, 90, 50, nil, 0)

	if len(merged) != 2 {
		t.Fatalf("want 2 after prune (old dropped), got %d: %+v", len(merged), merged)
	}
	if merged[0].ID != "b" {
		t.Errorf("expected newest (b) first, got %q", merged[0].ID)
	}
	var a Item
	for _, it := range merged {
		if it.ID == "a" {
			a = it
		}
	}
	if a.Published != "2026-08-01T00:00:00Z" {
		t.Errorf("existing publish date should be kept, got %q", a.Published)
	}
	if a.Summary != "longer summary" {
		t.Errorf("longer summary should be adopted, got %q", a.Summary)
	}
	if !hasTag(a.Tags, "advisory") || !hasTag(a.Tags, "KEV") {
		t.Errorf("tags should be unioned, got %v", a.Tags)
	}
}

func TestPerSourceCap(t *testing.T) {
	var items []Item
	for i := 0; i < 10; i++ {
		d := refNow.Add(-time.Duration(i) * time.Hour)
		items = append(items,
			Item{ID: "noisy" + strconv.Itoa(i), Source: "noisy", publishedAt: d},
			Item{ID: "quiet" + strconv.Itoa(i), Source: "quiet", publishedAt: d.Add(-30 * time.Minute)})
	}
	out := pruneAndSort(items, refNow, 90, 100, map[string]int{"noisy": 3}, 5)
	count := map[string]int{}
	for _, it := range out {
		count[it.Source]++
	}
	if count["noisy"] != 3 || count["quiet"] != 5 {
		t.Fatalf("caps not applied: %v", count)
	}
	if out[0].ID != "noisy0" {
		t.Errorf("newest item should survive the cap, got %q first", out[0].ID)
	}
}

func TestLangTagAndSummaryOptOut(t *testing.T) {
	off := false
	s := Source{ID: "cert-fr", Name: "CERT-FR", Lang: "fr", Tag: "alert", Summary: &off}
	it := newItem(s, "Vulnérabilité dans X (CVE-2026-1234)", "https://e.gouv.fr/a",
		"<p>Détails CVE-2026-5678</p>", "2026-08-20", refNow)
	if !hasTag(it.Tags, "FR") {
		t.Errorf("expected FR tag, got %v", it.Tags)
	}
	if it.Summary != "" {
		t.Errorf("summary should be dropped, got %q", it.Summary)
	}
	if !hasTag(it.Tags, "CVE-2026-5678") {
		t.Errorf("CVE ids from the summary should still become tags: %v", it.Tags)
	}
	en := newItem(Source{ID: "x", Lang: "en"}, "t", "https://e.gov/b", "s", "2026-08-20", refNow)
	if hasTag(en.Tags, "EN") || en.Summary != "s" {
		t.Errorf("English source should get no lang tag and keep its summary: %+v", en)
	}
}

func TestResolveURLAndValidation(t *testing.T) {
	s := Source{ID: "nvd", URL: "https://x.gov/cves?a={start}&b={end}", WindowDays: 3}
	got := s.resolveURL(refNow)
	want := "https://x.gov/cves?a=2026-08-23T00:00:00.000&b=2026-08-26T00:00:00.000"
	if got != want {
		t.Fatalf("resolveURL = %q, want %q", got, want)
	}
	bad := Config{Sources: []Source{{ID: "nvd", URL: "https://x.gov/?a={start}", Kind: "json"}}}
	if err := bad.validate(); err == nil {
		t.Error("placeholder without window_days should fail validation")
	}
}

func TestExcludeTitle(t *testing.T) {
	maxPerSource = 10
	src := Source{ID: "example-cert", Name: "Example CERT", Kind: "rss", ExcludeTitle: `^Guidance:`}
	items, err := parseSource(src, readFixture(t, "rss.xml"), refNow)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || strings.HasPrefix(items[0].Title, "Guidance:") {
		t.Fatalf("exclude_title not applied: %+v", items)
	}
	bad := Config{Sources: []Source{{ID: "x", URL: "https://e.gov/", Kind: "rss", ExcludeTitle: "("}}}
	if bad.validate() == nil {
		t.Error("invalid regexp should fail validation")
	}
}

func TestBriefing(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "news.json")
	st := `{"schema":1,"generated":"2026-08-26T00:00:00Z","items":[
	 {"id":"a","source":"cisa-kev","source_name":"CISA KEV","title":"KEV item","url":"https://e.gov/a","published":"2026-08-20T00:00:00Z","tags":["KEV","ransomware"]},
	 {"id":"b","source":"nvd-recent","source_name":"NVD","title":"Critical item","url":"https://e.gov/b","published":"2026-08-21T00:00:00Z","tags":["critical"]},
	 {"id":"c","source":"cisa-ics","source_name":"CISA ICS","title":"ICS noise","url":"https://e.gov/c","published":"2026-08-22T00:00:00Z","tags":["ics"]},
	 {"id":"d","source":"cisa-kev","source_name":"CISA KEV","title":"Too old","url":"https://e.gov/d","published":"2026-06-01T00:00:00Z","tags":["KEV"]}]}`
	os.WriteFile(path, []byte(st), 0o644)

	var buf strings.Builder
	if err := writeBriefing(&buf, path, 30, refNow); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	for _, want := range []string{"(2 items)", "## Actively exploited (CISA KEV) (1)", "[KEV item](https://e.gov/a)", "## Critical new vulnerabilities (NVD) (1)"} {
		if !strings.Contains(out, want) {
			t.Errorf("briefing missing %q:\n%s", want, out)
		}
	}
	for _, unwanted := range []string{"ICS noise", "Too old", "Linked to ransomware"} {
		if strings.Contains(out, unwanted) {
			t.Errorf("briefing should not contain %q:\n%s", unwanted, out)
		}
	}
}

func TestRunEndToEnd(t *testing.T) {
	srv := fixtureServer(t)
	defer srv.Close()
	testTransport = srv.Client().Transport
	defer func() { testTransport = nil }()

	dir := t.TempDir()
	feeds := filepath.Join(dir, "feeds.yaml")
	raw, _ := os.ReadFile(filepath.Join("testdata", "feeds.yaml"))
	os.WriteFile(feeds, []byte(strings.ReplaceAll(string(raw), "{{BASE}}", srv.URL)), 0o644)
	out := filepath.Join(dir, "news.json")

	changed, err := run(feeds, out, refNow.Format(time.RFC3339), false)
	if err != nil {
		t.Fatal(err)
	}
	if !changed {
		t.Fatal("expected first run to write the file")
	}

	var store Store
	b, _ := os.ReadFile(out)
	if err := json.Unmarshal(b, &store); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if store.Schema != storeSchema || len(store.Items) == 0 {
		t.Fatalf("unexpected store: %+v", store)
	}
	if !strings.HasSuffix(string(b), "\n") {
		t.Error("output should end with a newline")
	}
	// The 2025 KEV entry must have been pruned by the 90-day window.
	for _, it := range store.Items {
		if strings.Contains(it.URL, "CVE-2025-0001") {
			t.Error("stale item was not pruned")
		}
	}

	// Second identical run must be a no-op.
	changed, err = run(feeds, out, refNow.Format(time.RFC3339), false)
	if err != nil {
		t.Fatal(err)
	}
	if changed {
		t.Error("expected second run to be a no-op")
	}
}

func TestRunAllSourcesFailKeepsFile(t *testing.T) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer srv.Close()
	testTransport = srv.Client().Transport
	defer func() { testTransport = nil }()

	dir := t.TempDir()
	feeds := filepath.Join(dir, "feeds.yaml")
	raw, _ := os.ReadFile(filepath.Join("testdata", "feeds.yaml"))
	os.WriteFile(feeds, []byte(strings.ReplaceAll(string(raw), "{{BASE}}", srv.URL)), 0o644)
	out := filepath.Join(dir, "news.json")
	os.WriteFile(out, []byte(`{"schema":1,"generated":"2026-08-01T00:00:00Z","items":[]}`), 0o644)
	before, _ := os.ReadFile(out)

	if _, err := run(feeds, out, refNow.Format(time.RFC3339), false); err == nil {
		t.Fatal("expected an error when every source fails")
	}
	after, _ := os.ReadFile(out)
	if string(before) != string(after) {
		t.Error("news.json must not be modified when every source fails")
	}
}

// --- helpers ---------------------------------------------------------------

func fixtureServer(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	serve := func(name, ctype string) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", ctype)
			http.ServeFile(w, r, filepath.Join("testdata", name))
		}
	}
	mux.HandleFunc("/rss.xml", serve("rss.xml", "application/rss+xml"))
	mux.HandleFunc("/atom.xml", serve("atom.xml", "application/atom+xml"))
	mux.HandleFunc("/kev.json", serve("kev.json", "application/json"))
	return httptest.NewTLSServer(mux)
}

func readFixture(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func hasTag(tags []string, want string) bool {
	for _, t := range tags {
		if t == want {
			return true
		}
	}
	return false
}

func mustDate(s string) time.Time {
	t, ok := parseDate(s)
	if !ok {
		panic("bad date " + s)
	}
	return t
}
