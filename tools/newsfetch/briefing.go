package main

import (
	"fmt"
	"io"
	"time"
)

// briefingGroups orders the digest by how much leadership should care.
// An item lands in the first group it matches.
var briefingGroups = []struct {
	title string
	tags  []string
}{
	{"Actively exploited (CISA KEV)", []string{"KEV"}},
	{"Linked to ransomware campaigns", []string{"ransomware"}},
	{"National CERT alerts", []string{"alert"}},
	{"Critical new vulnerabilities (NVD)", []string{"critical"}},
	{"Guidance and threat intelligence", []string{"guidance", "threat-intel", "research"}},
}

func (it Item) hasAny(tags []string) bool {
	for _, have := range it.Tags {
		for _, want := range tags {
			if have == want {
				return true
			}
		}
	}
	return false
}

// writeBriefing prints a Markdown digest of the last `days` days of
// data/news.json — raw material for a human-written Board Briefing.
func writeBriefing(w io.Writer, storePath string, days int, now time.Time) error {
	st, err := loadStore(storePath)
	if err != nil {
		return err
	}
	cutoff := now.AddDate(0, 0, -days)

	grouped := make([][]Item, len(briefingGroups))
	total := 0
	for _, it := range st.Items {
		if it.publishedAt.Before(cutoff) {
			continue
		}
		for g, grp := range briefingGroups {
			if it.hasAny(grp.tags) {
				grouped[g] = append(grouped[g], it)
				total++
				break
			}
		}
	}

	fmt.Fprintf(w, "# Briefing raw material: %s – %s (%d items)\n\n",
		cutoff.Format("2 Jan 2006"), now.Format("2 Jan 2006"), total)
	fmt.Fprintln(w, "Facts only — rewrite for a non-technical reader before publishing.")
	for g, grp := range briefingGroups {
		if len(grouped[g]) == 0 {
			continue
		}
		fmt.Fprintf(w, "\n## %s (%d)\n\n", grp.title, len(grouped[g]))
		for _, it := range grouped[g] {
			fmt.Fprintf(w, "- %s — [%s](%s) · %s\n",
				it.publishedAt.Format("2 Jan"), it.Title, it.URL, it.SourceName)
		}
	}
	return nil
}
