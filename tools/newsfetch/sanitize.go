package main

import (
	"crypto/sha256"
	"encoding/hex"
	"html"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var (
	tagRE   = regexp.MustCompile(`(?s)<[^>]*>`)
	wsRE    = regexp.MustCompile(`[\s\p{Zs}\x{feff}]+`)
	cveRE   = regexp.MustCompile(`CVE-\d{4}-\d{4,7}`)
	trackRE = regexp.MustCompile(`^(utm_|mc_|pk_|ref$|source$|fbclid$|gclid$)`)
)

// toPlainText strips HTML tags, unescapes entities, collapses whitespace, and
// truncates to max runes (0 = no limit). The result is display-safe regardless,
// because the site templates escape on output — this is about readability.
func toPlainText(s string, max int) string {
	s = tagRE.ReplaceAllString(s, " ")
	s = html.UnescapeString(s)
	s = wsRE.ReplaceAllString(s, " ")
	s = strings.TrimSpace(s)
	if max > 0 {
		s = truncateRunes(s, max)
	}
	return s
}

func truncateRunes(s string, max int) string {
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	cut := strings.TrimRight(string(r[:max]), " ,.;:—-")
	return cut + "…"
}

// canonicalURL lowercases the host, drops the fragment, and strips common
// tracking query parameters so the same article from two feeds dedups.
func canonicalURL(raw string) string {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || u.Host == "" {
		return strings.TrimSpace(raw)
	}
	u.Host = strings.ToLower(u.Host)
	u.Scheme = strings.ToLower(u.Scheme)
	u.Fragment = ""
	if q := u.Query(); len(q) > 0 {
		for k := range q {
			if trackRE.MatchString(k) {
				q.Del(k)
			}
		}
		u.RawQuery = q.Encode()
	}
	return u.String()
}

// itemID is a stable short id derived from the canonical URL.
func itemID(canonURL string) string {
	sum := sha256.Sum256([]byte(canonURL))
	return hex.EncodeToString(sum[:])[:16]
}

// dateLayouts covers the formats seen across the allowlisted feeds, including
// CISA's RFC822-with-weekday two-digit-year form ("Thu, 27 Aug 26 12:00:00 +0000").
var dateLayouts = []string{
	time.RFC3339,
	time.RFC3339Nano,
	time.RFC1123Z,
	time.RFC1123,
	time.RFC822Z,
	time.RFC822,
	"Mon, 02 Jan 06 15:04:05 -0700",
	"Mon, 2 Jan 06 15:04:05 -0700",
	"Mon, 02 Jan 06 15:04:05 MST",
	"Mon, 2 Jan 2006 15:04:05 -0700",
	"Mon, 2 Jan 2006 15:04:05 MST",
	"2006-01-02T15:04:05.000Z07:00",
	"2006-01-02T15:04:05.000",
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006-01-02",
}

// parseDate tries each known layout; ok is false if none match.
func parseDate(s string) (t time.Time, ok bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, false
	}
	for _, l := range dateLayouts {
		if parsed, err := time.Parse(l, s); err == nil {
			return parsed.UTC(), true
		}
	}
	return time.Time{}, false
}

// extractCVEs returns unique CVE identifiers found in text, upper-cased.
func extractCVEs(text string) []string {
	m := cveRE.FindAllString(strings.ToUpper(text), -1)
	if len(m) == 0 {
		return nil
	}
	return dedupeStrings(m)
}

func dedupeStrings(in []string) []string {
	seen := make(map[string]bool, len(in))
	out := make([]string, 0, len(in))
	for _, s := range in {
		if s == "" || seen[s] {
			continue
		}
		seen[s] = true
		out = append(out, s)
	}
	return out
}
