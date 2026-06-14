package kanyerest_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tamnd/kanyerest-cli/kanyerest"
)

const fakeQuoteJSON = `{"quote":"I'm nice at ping pong"}`

func newTestClient(ts *httptest.Server) *kanyerest.Client {
	cfg := kanyerest.DefaultConfig()
	cfg.BaseURL = ts.URL
	cfg.Rate = 0
	return kanyerest.NewClient(cfg)
}

func TestRandomQuoteParsesQuote(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, fakeQuoteJSON)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	q, err := c.RandomQuote(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if q.Quote != "I'm nice at ping pong" {
		t.Errorf("Quote = %q, want \"I'm nice at ping pong\"", q.Quote)
	}
}

func TestRandomQuoteSendsUserAgent(t *testing.T) {
	var gotUA string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUA = r.Header.Get("User-Agent")
		_, _ = fmt.Fprint(w, fakeQuoteJSON)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	_, err := c.RandomQuote(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if gotUA == "" {
		t.Error("request carried no User-Agent")
	}
}

func TestRandomQuoteRetriesOn503(t *testing.T) {
	var hits int
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		if hits < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		_, _ = fmt.Fprint(w, fakeQuoteJSON)
	}))
	defer ts.Close()

	cfg := kanyerest.DefaultConfig()
	cfg.BaseURL = ts.URL
	cfg.Rate = 0
	cfg.Retries = 3
	c := kanyerest.NewClient(cfg)

	_, err := c.RandomQuote(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if hits != 3 {
		t.Errorf("server saw %d hits, want 3", hits)
	}
}

func TestRandomQuoteNon200Error(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	_, err := c.RandomQuote(context.Background())
	if err == nil {
		t.Error("expected error for 404, got nil")
	}
}
