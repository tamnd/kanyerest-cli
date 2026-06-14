// Package kanyerest exposes the Kanye REST API as a kit Domain.
//
// A multi-domain host (ant) enables it with a single blank import:
//
//	import _ "github.com/tamnd/kanyerest-cli/kanyerest"
//
// The same Domain also builds the standalone kanyerest binary.
package kanyerest

import (
	"context"

	"github.com/tamnd/any-cli/kit"
	"github.com/tamnd/any-cli/kit/errs"
)

func init() { kit.Register(Domain{}) }

// Domain is the kanyerest driver.
type Domain struct{}

// Info describes the scheme, hosts, and binary identity.
func (Domain) Info() kit.DomainInfo {
	return kit.DomainInfo{
		Scheme: "kanyerest",
		Hosts:  []string{Host},
		Identity: kit.Identity{
			Binary: "kanyerest",
			Short:  "Random Kanye West quotes from api.kanye.rest",
			Long: `kanyerest fetches random Kanye West quotes from api.kanye.rest.
No API key or authentication required.`,
			Site: Host,
			Repo: "https://github.com/tamnd/kanyerest-cli",
		},
	}
}

// Register installs the client factory and all operations onto app.
func (Domain) Register(app *kit.App) {
	app.SetClient(newClient)

	kit.Handle(app, kit.OpMeta{
		Name:    "quote",
		Group:   "read",
		Summary: "Get a random Kanye West quote",
	}, quoteOp)
}

// newClient builds the client from host-resolved config.
func newClient(_ context.Context, cfg kit.Config) (any, error) {
	c := DefaultConfig()
	if cfg.UserAgent != "" {
		c.UserAgent = cfg.UserAgent
	}
	if cfg.Rate > 0 {
		c.Rate = cfg.Rate
	}
	if cfg.Retries > 0 {
		c.Retries = cfg.Retries
	}
	if cfg.Timeout > 0 {
		c.Timeout = cfg.Timeout
	}
	return NewClient(c), nil
}

// --- inputs ---

type quoteInput struct {
	Client *Client `kit:"inject"`
}

// --- handlers ---

func quoteOp(ctx context.Context, in quoteInput, emit func(Quote) error) error {
	q, err := in.Client.RandomQuote(ctx)
	if err != nil {
		return mapErr(err)
	}
	return emit(q)
}

// --- Resolver ---

// Classify turns an input into the canonical (type, id).
func (Domain) Classify(input string) (uriType, id string, err error) {
	if input == "" {
		return "", "", errs.Usage("empty kanyerest reference")
	}
	return "quote", input, nil
}

// Locate returns the live https URL for a (type, id).
func (Domain) Locate(uriType, id string) (string, error) {
	switch uriType {
	case "quote":
		return BaseURL + "/", nil
	default:
		return "", errs.Usage("kanyerest has no resource type %q", uriType)
	}
}

// mapErr converts library errors into kit error kinds.
func mapErr(err error) error {
	return err
}
