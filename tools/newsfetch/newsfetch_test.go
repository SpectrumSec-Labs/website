package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
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
	merged := pruneAndSort(mergeItems(existing, fresh), refNow, 90, 50)

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
