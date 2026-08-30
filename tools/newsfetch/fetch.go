package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// fetcher performs bounded HTTP GETs: HTTPS only, capped body, capped redirects,
// explicit User-Agent, modern TLS.
type fetcher struct {
	client *http.Client
	ua     string
	maxLen int64
}

// testTransport, when set, replaces the HTTP transport. Used only by tests to
// point the fetcher at an in-process TLS server.
var testTransport http.RoundTripper

func newFetcher(l Limits) *fetcher {
	f := &fetcher{
		ua:     l.UserAgent,
		maxLen: l.MaxResponseBytes,
		client: &http.Client{
			Timeout: time.Duration(l.HTTPTimeoutSeconds) * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 5 {
					return errors.New("stopped after 5 redirects")
				}
				if req.URL.Scheme != "https" {
					return fmt.Errorf("refusing redirect to non-https URL %q", req.URL)
				}
				return nil
			},
			Transport: &http.Transport{
				TLSClientConfig:     &tls.Config{MinVersion: tls.VersionTLS12},
				DisableKeepAlives:   true,
				MaxIdleConns:        2,
				TLSHandshakeTimeout: 10 * time.Second,
			},
		},
	}
	if testTransport != nil {
		f.client.Transport = testTransport
	}
	return f
}

func (f *fetcher) get(ctx context.Context, url string) ([]byte, error) {
	if !strings.HasPrefix(url, "https://") {
		return nil, fmt.Errorf("url is not https: %q", url)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", f.ua)
	req.Header.Set("Accept", "application/rss+xml, application/atom+xml, application/xml;q=0.9, application/json;q=0.9, text/xml;q=0.8, */*;q=0.5")

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %s", resp.Status)
	}

	// Read one byte past the cap so we can detect an over-size body.
	b, err := io.ReadAll(io.LimitReader(resp.Body, f.maxLen+1))
	if err != nil {
		return nil, err
	}
	if int64(len(b)) > f.maxLen {
		return nil, fmt.Errorf("response exceeds %d bytes", f.maxLen)
	}
	return b, nil
}
