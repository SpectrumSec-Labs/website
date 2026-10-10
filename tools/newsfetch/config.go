package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Config mirrors config/feeds.yaml.
type Config struct {
	Sources []Source `yaml:"sources"`
	Limits  Limits   `yaml:"limits"`
}

// Source is one entry in the feed allowlist.
type Source struct {
	ID       string `yaml:"id"`
	Name     string `yaml:"name"`
	URL      string `yaml:"url"`
	Kind     string `yaml:"kind"` // "rss", "atom" (both parsed as XML) or "json"
	Tag      string `yaml:"tag"`
	Homepage string `yaml:"homepage"`
	Enabled  bool   `yaml:"enabled"`

	// Lang is the source language (ISO 639-1). Non-English items get an
	// upper-cased language tag (e.g. "FR") so readers know before clicking.
	Lang string `yaml:"lang"`
	// Summary=false keeps only title + link, for sources whose licence does
	// not allow reproducing text on a commercial site (e.g. CC BY-NC).
	Summary *bool `yaml:"summary"`
	// WindowDays fills {start}/{end} placeholders in URL with a rolling
	// window ending now, for date-filtered APIs such as NVD.
	WindowDays int `yaml:"window_days"`
	// MaxStored caps how many items from this source data/news.json keeps,
	// so high-volume sources cannot crowd out the rest.
	MaxStored int `yaml:"max_stored"`
	// ExcludeTitle drops items whose title matches this RE2 regexp
	// (e.g. recurring podcast episodes in an otherwise useful feed).
	ExcludeTitle string `yaml:"exclude_title"`
}

func (s Source) keepSummary() bool { return s.Summary == nil || *s.Summary }

// Limits are global bounds applied to every fetch and to the output file.
type Limits struct {
	HTTPTimeoutSeconds int    `yaml:"http_timeout_seconds"`
	MaxResponseBytes   int64  `yaml:"max_response_bytes"`
	MaxItemsPerSource  int    `yaml:"max_items_per_source"`
	MaxStoredPerSource int    `yaml:"max_stored_per_source"`
	UserAgent          string `yaml:"user_agent"`
	RetentionDays      int    `yaml:"retention_days"`
	MaxItems           int    `yaml:"max_items"`
}

func loadConfig(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c Config
	dec := yaml.NewDecoder(strings.NewReader(string(b)))
	dec.KnownFields(true)
	if err := dec.Decode(&c); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}
	if err := c.validate(); err != nil {
		return nil, fmt.Errorf("invalid %s: %w", path, err)
	}
	c.applyDefaults()
	return &c, nil
}

func (c *Config) validate() error {
	if len(c.Sources) == 0 {
		return fmt.Errorf("no sources defined")
	}
	seen := map[string]bool{}
	for i, s := range c.Sources {
		switch {
		case s.ID == "":
			return fmt.Errorf("source %d: missing id", i)
		case seen[s.ID]:
			return fmt.Errorf("duplicate source id %q", s.ID)
		case !strings.HasPrefix(s.URL, "https://"):
			return fmt.Errorf("source %q: url must be https", s.ID)
		case s.Kind != "rss" && s.Kind != "atom" && s.Kind != "json":
			return fmt.Errorf("source %q: kind must be rss, atom or json", s.ID)
		case strings.Contains(s.URL, "{start}") != (s.WindowDays > 0):
			return fmt.Errorf("source %q: {start}/{end} in url and window_days must be set together", s.ID)
		}
		if s.ExcludeTitle != "" {
			if _, err := regexp.Compile(s.ExcludeTitle); err != nil {
				return fmt.Errorf("source %q: bad exclude_title: %w", s.ID, err)
			}
		}
		seen[s.ID] = true
	}
	return nil
}

func (c *Config) applyDefaults() {
	if c.Limits.HTTPTimeoutSeconds <= 0 {
		c.Limits.HTTPTimeoutSeconds = 10
	}
	if c.Limits.MaxResponseBytes <= 0 {
		c.Limits.MaxResponseBytes = 3 << 20
	}
	if c.Limits.MaxItemsPerSource <= 0 {
		c.Limits.MaxItemsPerSource = 40
	}
	if c.Limits.MaxStoredPerSource <= 0 {
		c.Limits.MaxStoredPerSource = 40
	}
	if c.Limits.RetentionDays <= 0 {
		c.Limits.RetentionDays = 90
	}
	if c.Limits.MaxItems <= 0 {
		c.Limits.MaxItems = 200
	}
	if c.Limits.UserAgent == "" {
		c.Limits.UserAgent = "SpectrumSecNewsBot/1.0 (+https://spectrumsec.eu/news/)"
	}
}

// storedCaps maps every configured source id (enabled or not) to its cap.
func (c *Config) storedCaps() map[string]int {
	caps := make(map[string]int, len(c.Sources))
	for _, s := range c.Sources {
		if s.MaxStored > 0 {
			caps[s.ID] = s.MaxStored
		} else {
			caps[s.ID] = c.Limits.MaxStoredPerSource
		}
	}
	return caps
}

// resolveURL fills the {start}/{end} placeholders (NVD date format, UTC).
func (s Source) resolveURL(now time.Time) string {
	if s.WindowDays <= 0 {
		return s.URL
	}
	const layout = "2006-01-02T15:04:05.000"
	start := now.UTC().AddDate(0, 0, -s.WindowDays).Format(layout)
	end := now.UTC().Format(layout)
	return strings.NewReplacer("{start}", start, "{end}", end).Replace(s.URL)
}

// enabledSources returns only the sources marked enabled, preserving order.
func (c *Config) enabledSources() []Source {
	out := make([]Source, 0, len(c.Sources))
	for _, s := range c.Sources {
		if s.Enabled {
			out = append(out, s)
		}
	}
	return out
}
