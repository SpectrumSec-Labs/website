package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	var (
		feedsPath = flag.String("feeds", "config/feeds.yaml", "path to the feed allowlist")
		outPath   = flag.String("out", "data/news.json", "path to the news JSON to update")
		nowFlag   = flag.String("now", "", "override current time (RFC3339); for tests")
		dryRun    = flag.Bool("dry-run", false, "compute changes but do not write")
	)
	flag.Parse()

	log.SetFlags(0)
	log.SetPrefix("newsfetch: ")

	changed, err := run(*feedsPath, *outPath, *nowFlag, *dryRun)
	if err != nil {
		log.Fatal(err)
	}
	if changed {
		fmt.Println("news.json updated")
	} else {
		fmt.Println("no changes")
	}
}

// run returns whether the output file was (or would be) changed.
func run(feedsPath, outPath, nowOverride string, dryRun bool) (bool, error) {
	cfg, err := loadConfig(feedsPath)
	if err != nil {
		return false, err
	}
	maxPerSource = cfg.Limits.MaxItemsPerSource

	now, err := nowUTC(nowOverride)
	if err != nil {
		return false, err
	}

	prev, err := loadStore(outPath)
	if err != nil {
		return false, err
	}

	sources := cfg.enabledSources()
	if len(sources) == 0 {
		return false, fmt.Errorf("no enabled sources in %s", feedsPath)
	}

	f := newFetcher(cfg.Limits)
	var fresh []Item
	var okCount, failCount int

	for _, s := range sources {
		ctx, cancel := context.WithTimeout(context.Background(),
			time.Duration(cfg.Limits.HTTPTimeoutSeconds+5)*time.Second)
		body, err := f.get(ctx, s.URL)
		cancel()
		if err != nil {
			failCount++
			log.Printf("WARN source %q: fetch failed: %v", s.ID, err)
			continue
		}
		items, err := parseSource(s, body, now)
		if err != nil {
			failCount++
			log.Printf("WARN source %q: parse failed: %v", s.ID, err)
			continue
		}
		okCount++
		log.Printf("source %q: %d items", s.ID, len(items))
		fresh = append(fresh, items...)
	}

	// Never destroy existing data because the network had a bad day.
	if okCount == 0 {
		return false, fmt.Errorf("all %d sources failed; leaving %s untouched", failCount, outPath)
	}

	merged := mergeItems(prev.Items, fresh)
	merged = pruneAndSort(merged, now, cfg.Limits.RetentionDays, cfg.Limits.MaxItems)

	if itemsEqual(prev.Items, merged) {
		return false, nil
	}

	next := &Store{
		Schema:    storeSchema,
		Generated: now.Format(time.RFC3339),
		Items:     merged,
	}
	data, err := next.marshal()
	if err != nil {
		return false, err
	}

	log.Printf("summary: %d sources ok, %d failed, %d items (%d new-or-changed)",
		okCount, failCount, len(merged), len(merged)-countUnchanged(prev.Items, merged))

	if dryRun {
		os.Stdout.Write(data)
		return true, nil
	}
	if err := writeAtomic(outPath, data); err != nil {
		return false, err
	}
	return true, nil
}

func countUnchanged(prev, next []Item) int {
	old := make(map[string]bool, len(prev))
	for _, it := range prev {
		old[it.ID] = true
	}
	n := 0
	for _, it := range next {
		if old[it.ID] {
			n++
		}
	}
	return n
}
