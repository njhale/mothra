package bugzilla

import (
	"context"
	"encoding/json"
	"net/http"

	pbz "k8s.io/test-infra/prow/bugzilla"
)

type Bug pbz.Bug

type bugList struct {
	Bugs []*Bug `json:"bugs,omitempty"`
}

type Stream interface {
	Send(bug *Bug) error
}

type Query interface {
	URL() string
}

type AuthFunc func(req *http.Request) error

type Client struct {
	client       *http.Client
	authenticate AuthFunc
}

func (c *Client) List(ctx context.Context, query Query, res Stream) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, query.URL(), nil)
	if err != nil {
		return err
	}

	if err := c.authenticate(req); err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	body := resp.Body
	defer func() {
		if err := body.Close(); err != nil {
			// TODO(njhale): log error
		}
	}()

	// FIXME(njhale): better stream parsing
	dec := json.NewDecoder(body)
	var list bugList
	if err := dec.Decode(&list); err != nil {
		return err
	}

	for _, bug := range list.Bugs {
		if err := res.Send(bug); err != nil {
			return err
		}
	}

	return nil
}

type ClientConfig struct {
	Authenticate AuthFunc
}

func (c *ClientConfig) apply(options []ClientOption) {
	for _, option := range options {
		option(c)
	}
}

func defaultConfig() *ClientConfig {
	return &ClientConfig{
		Authenticate: func(req *http.Request) error {
			return nil
		},
	}
}

func NewClient(options ...ClientOption) (client *Client, err error) {
	config := defaultConfig()
	config.apply(options)

	client = &Client{
		client:       http.DefaultClient(),
		authenticate: config.Authenticate,
	}

	return
}

type ClientOption func(*ClientConfig)

func WithAPIKey(apiKey string) ClientOption {
	return func(c *ClientConfig) {
		c.Authenticate = func(req *http.Request) error {
			req.Header.Set("X-BUGZILLA-API-KEY", apiKey)
			values := req.URL.Query()
			values.Add("api_key", apiKey)
			req.URL.RawQuery = values.Encode()
			return nil
		}
	}
}
