package main

import (
	"sort"
	"strings"
	"time"
)

// Item is one aggregated news entry, serialized into data/news.json.
type Item struct {
	ID         string   `json:"id"`
	Source     string   `json:"source"`
	SourceName string   `json:"source_name"`
	Title      string   `json:"title"`
	URL        string   `json:"url"`
	Published  string   `json:"published"` // RFC3339 UTC
	Summary    string   `json:"summary,omitempty"`
	Tags       []string `json:"tags"`

	// publishedAt is the parsed form of Published, used for sorting/pruning.
	// Unexported, so it is never serialized; reconstructed on load.
	publishedAt time.Time
}

const (
	maxTitleRunes   = 200
	maxSummaryRunes = 280
	maxTags         = 6
)

// newItem builds a normalized Item from raw feed fields.
func newItem(s Source, title, rawURL, summary, rawDate string, now time.Time) Item {
	canon := canonicalURL(rawURL)

	when, ok := parseDate(rawDate)
	if !ok {
		when = now.UTC()
	}

	title = toPlainText(title, maxTitleRunes)
	summary = toPlainText(summary, maxSummaryRunes)

	it := Item{
		ID:          itemID(canon),
		Source:      s.ID,
		SourceName:  s.Name,
		Title:       title,
		URL:         canon,
		Published:   when.Format(time.RFC3339),
		Summary:     summary,
		publishedAt: when,
	}

	if s.Tag != "" {
		it.Tags = append(it.Tags, s.Tag)
	}
	for _, cve := range extractCVEs(title + " " + summary) {
		it.Tags = append(it.Tags, cve)
	}
	it.finalizeTags()
	return it
}

// finalizeTags trims, de-dupes and sorts tags for deterministic output, and
// caps the count. "KEV" is kept first when present because it is the signal
// readers scan for.
func (it *Item) finalizeTags() {
	cleaned := make([]string, 0, len(it.Tags))
	for _, t := range it.Tags {
		t = strings.TrimSpace(t)
		if t != "" {
			cleaned = append(cleaned, t)
		}
	}
	cleaned = dedupeStrings(cleaned)

	kev := false
	rest := cleaned[:0]
	for _, t := range cleaned {
		if t == "KEV" {
			kev = true
			continue
		}
		rest = append(rest, t)
	}
	sort.Strings(rest)

	out := make([]string, 0, maxTags)
	if kev {
		out = append(out, "KEV")
	}
	out = append(out, rest...)
	if len(out) > maxTags {
		out = out[:maxTags]
	}
	it.Tags = out
}

// valid reports whether an item has the minimum fields to be published.
func (it Item) valid() bool {
	return it.ID != "" && it.Title != "" && strings.HasPrefix(it.URL, "https://")
}

// mergeItems combines existing and freshly fetched items. Existing items win on
// identity (stable id, original publish date); tags are unioned; a longer
// summary from the new copy is adopted.
func mergeItems(existing, fresh []Item) []Item {
	byID := make(map[string]Item, len(existing)+len(fresh))
	order := make([]string, 0, len(existing)+len(fresh))

	add := func(it Item, isExisting bool) {
		if !it.valid() {
			return
		}
		if it.publishedAt.IsZero() {
			if t, ok := parseDate(it.Published); ok {
				it.publishedAt = t
			}
		}
		prev, seen := byID[it.ID]
		if !seen {
			byID[it.ID] = it
			order = append(order, it.ID)
			return
		}
		// Merge into the copy we already have.
		merged := prev
		merged.Tags = mergeTags(prev.Tags, it.Tags)
		if len(it.Summary) > len(merged.Summary) {
			merged.Summary = it.Summary
		}
		if isExisting {
			// keep existing publish date and source name
			merged.Published = it.Published
			merged.publishedAt = it.publishedAt
		}
		byID[it.ID] = merged
	}

	for _, it := range existing {
		add(it, true)
	}
	for _, it := range fresh {
		add(it, false)
	}

	out := make([]Item, 0, len(order))
	for _, id := range order {
		out = append(out, byID[id])
	}
	return out
}

func mergeTags(a, b []string) []string {
	it := Item{Tags: append(append([]string{}, a...), b...)}
	it.finalizeTags()
	return it.Tags
}

// pruneAndSort drops items older than the retention window, sorts newest-first
// (id as a stable tie-break), and caps the total.
func pruneAndSort(items []Item, now time.Time, retentionDays, max int) []Item {
	cutoff := now.AddDate(0, 0, -retentionDays)
	kept := make([]Item, 0, len(items))
	for _, it := range items {
		if it.publishedAt.Before(cutoff) {
			continue
		}
		kept = append(kept, it)
	}
	sort.SliceStable(kept, func(i, j int) bool {
		if kept[i].publishedAt.Equal(kept[j].publishedAt) {
			return kept[i].ID < kept[j].ID
		}
		return kept[i].publishedAt.After(kept[j].publishedAt)
	})
	if len(kept) > max {
		kept = kept[:max]
	}
	return kept
}
