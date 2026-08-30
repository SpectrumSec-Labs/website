package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"sort"
	"strings"
	"time"
)

// parseSource turns a fetched payload into normalized items for one source.
// now is passed in so tests are deterministic and so a missing date can fall
// back to run time.
func parseSource(s Source, body []byte, now time.Time) ([]Item, error) {
	var (
		items []Item
		err   error
	)
	switch s.Kind {
	case "json":
		items, err = parseJSON(s, body, now)
	default: // rss / atom — try both regardless of the hint
		items, err = parseXML(s, body, now)
	}
	if err != nil {
		return nil, err
	}
	// Newest first, then cap per source.
	sort.SliceStable(items, func(i, j int) bool {
		return items[i].publishedAt.After(items[j].publishedAt)
	})
	if len(items) > maxPerSource {
		items = items[:maxPerSource]
	}
	return items, nil
}

// maxPerSource is set from config before parsing (avoids threading Limits
// through every call).
var maxPerSource = 40

// ---------------------------------------------------------------------------
// RSS / Atom
// ---------------------------------------------------------------------------

type rssDoc struct {
	Channel struct {
		Title string `xml:"title"`
		Items []struct {
			Title       string `xml:"title"`
			Link        string `xml:"link"`
			GUID        string `xml:"guid"`
			PubDate     string `xml:"pubDate"`
			DCDate      string `xml:"date"` // dc:date, namespace-agnostic
			Description string `xml:"description"`
		} `xml:"item"`
	} `xml:"channel"`
}

type atomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
}

type atomDoc struct {
	Title   string `xml:"title"`
	Entries []struct {
		Title     string     `xml:"title"`
		Links     []atomLink `xml:"link"`
		ID        string     `xml:"id"`
		Updated   string     `xml:"updated"`
		Published string     `xml:"published"`
		Summary   string     `xml:"summary"`
		Content   string     `xml:"content"`
	} `xml:"entry"`
}

func parseXML(s Source, body []byte, now time.Time) ([]Item, error) {
	var rss rssDoc
	if err := xml.Unmarshal(body, &rss); err == nil && len(rss.Channel.Items) > 0 {
		out := make([]Item, 0, len(rss.Channel.Items))
		for _, it := range rss.Channel.Items {
			date := firstNonEmpty(it.PubDate, it.DCDate)
			out = append(out, newItem(s, it.Title, it.Link, it.Description, date, now))
		}
		return out, nil
	}

	var atom atomDoc
	if err := xml.Unmarshal(body, &atom); err == nil && len(atom.Entries) > 0 {
		out := make([]Item, 0, len(atom.Entries))
		for _, e := range atom.Entries {
			link := pickAtomLink(e.Links)
			date := firstNonEmpty(e.Published, e.Updated)
			body := firstNonEmpty(e.Summary, e.Content)
			out = append(out, newItem(s, e.Title, link, body, date, now))
		}
		return out, nil
	}

	return nil, fmt.Errorf("no RSS <item> or Atom <entry> elements found")
}

func pickAtomLink(links []atomLink) string {
	var fallback string
	for _, l := range links {
		if l.Href == "" {
			continue
		}
		if fallback == "" {
			fallback = l.Href
		}
		if l.Rel == "" || l.Rel == "alternate" {
			return l.Href
		}
	}
	return fallback
}

// ---------------------------------------------------------------------------
// JSON sources (keyed by source id — each has a fixed, known shape)
// ---------------------------------------------------------------------------

func parseJSON(s Source, body []byte, now time.Time) ([]Item, error) {
	switch s.ID {
	case "cisa-kev":
		return parseKEV(s, body, now)
	case "nvd-recent":
		return parseNVD(s, body, now)
	default:
		return nil, fmt.Errorf("no JSON parser registered for source id %q", s.ID)
	}
}

type kevDoc struct {
	CatalogVersion  string `json:"catalogVersion"`
	Vulnerabilities []struct {
		CVEID                      string `json:"cveID"`
		VendorProject              string `json:"vendorProject"`
		Product                    string `json:"product"`
		VulnerabilityName          string `json:"vulnerabilityName"`
		DateAdded                  string `json:"dateAdded"`
		ShortDescription           string `json:"shortDescription"`
		KnownRansomwareCampaignUse string `json:"knownRansomwareCampaignUse"`
	} `json:"vulnerabilities"`
}

func parseKEV(s Source, body []byte, now time.Time) ([]Item, error) {
	var doc kevDoc
	if err := json.Unmarshal(body, &doc); err != nil {
		return nil, err
	}
	if len(doc.Vulnerabilities) == 0 {
		return nil, fmt.Errorf("KEV catalog has no vulnerabilities")
	}
	// Sort by dateAdded desc so the per-source cap keeps the most recent.
	sort.Slice(doc.Vulnerabilities, func(i, j int) bool {
		return doc.Vulnerabilities[i].DateAdded > doc.Vulnerabilities[j].DateAdded
	})
	limit := maxPerSource
	if limit > len(doc.Vulnerabilities) {
		limit = len(doc.Vulnerabilities)
	}

	out := make([]Item, 0, limit)
	for _, v := range doc.Vulnerabilities[:limit] {
		title := strings.TrimSpace(fmt.Sprintf("%s — %s %s: %s",
			v.CVEID, v.VendorProject, v.Product, v.VulnerabilityName))
		it := newItem(s, title,
			"https://nvd.nist.gov/vuln/detail/"+v.CVEID,
			v.ShortDescription, v.DateAdded, now)
		it.Tags = append(it.Tags, "KEV", v.CVEID)
		if strings.EqualFold(v.KnownRansomwareCampaignUse, "Known") {
			it.Tags = append(it.Tags, "ransomware")
		}
		it.finalizeTags()
		out = append(out, it)
	}
	return out, nil
}

type nvdDoc struct {
	Vulnerabilities []struct {
		CVE struct {
			ID           string `json:"id"`
			Published    string `json:"published"`
			Descriptions []struct {
				Lang  string `json:"lang"`
				Value string `json:"value"`
			} `json:"descriptions"`
		} `json:"cve"`
	} `json:"vulnerabilities"`
}

func parseNVD(s Source, body []byte, now time.Time) ([]Item, error) {
	var doc nvdDoc
	if err := json.Unmarshal(body, &doc); err != nil {
		return nil, err
	}
	out := make([]Item, 0, len(doc.Vulnerabilities))
	for _, v := range doc.Vulnerabilities {
		desc := ""
		for _, d := range v.CVE.Descriptions {
			if d.Lang == "en" {
				desc = d.Value
				break
			}
		}
		it := newItem(s, v.CVE.ID+": "+truncateRunes(toPlainText(desc, 0), 90),
			"https://nvd.nist.gov/vuln/detail/"+v.CVE.ID,
			desc, v.CVE.Published, now)
		it.Tags = append(it.Tags, "CVE", v.CVE.ID)
		it.finalizeTags()
		out = append(out, it)
	}
	return out, nil
}

func firstNonEmpty(vs ...string) string {
	for _, v := range vs {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
