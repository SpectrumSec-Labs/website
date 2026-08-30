package main

import (
	"fmt"
	"os"
	"strings"

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
}

// Limits are global bounds applied to every fetch and to the output file.
type Limits struct {
	HTTPTimeoutSeconds int    `yaml:"http_timeout_seconds"`
	MaxResponseBytes   int64  `yaml:"max_response_bytes"`
	MaxItemsPerSource  int    `yaml:"max_items_per_source"`
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
