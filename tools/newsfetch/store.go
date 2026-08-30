package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const storeSchema = 1

// Store is the on-disk shape of data/news.json.
type Store struct {
	Schema    int    `json:"schema"`
	Generated string `json:"generated"`
	Note      string `json:"note,omitempty"`
	Items     []Item `json:"items"`
}

func loadStore(path string) (*Store, error) {
	b, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return &Store{Schema: storeSchema}, nil
	}
	if err != nil {
		return nil, err
	}
	var s Store
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}
	// Rehydrate publishedAt for sorting/pruning.
	for i := range s.Items {
		if t, ok := parseDate(s.Items[i].Published); ok {
			s.Items[i].publishedAt = t
		}
	}
	return &s, nil
}

// marshal produces the deterministic byte form (2-space indent, trailing NL,
// HTML escaping off so URLs stay readable).
func (s *Store) marshal() ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// itemsEqual compares only the published item set (ignoring the Generated
// timestamp) so an unchanged run produces no write and no PR.
func itemsEqual(a, b []Item) bool {
	ab, _ := json.Marshal(a)
	bb, _ := json.Marshal(b)
	return bytes.Equal(ab, bb)
}

// writeAtomic writes via a temp file + rename in the same directory.
func writeAtomic(path string, data []byte) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".news-*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpName, path)
}

func nowUTC(override string) (time.Time, error) {
	if override == "" {
		return time.Now().UTC(), nil
	}
	t, err := time.Parse(time.RFC3339, override)
	if err != nil {
		return time.Time{}, fmt.Errorf("-now must be RFC3339: %w", err)
	}
	return t.UTC(), nil
}
