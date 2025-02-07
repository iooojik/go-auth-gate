package apple

import (
	"net/http"
)

type Config struct {
	TokenConfig TokenConfig `yaml:"token"`
	URL         string      `yaml:"url"`
	KeyPath     string      `yaml:"keyPath"`
}

type Client struct {
	cfg            Config
	tokenGenerator ClientSecretGenerator
	client         HTTPClient
}

// New creates a new instance of Client.
func New(cfg Config, tokenGenerator ClientSecretGenerator, client HTTPClient) *Client {
	if client == nil {
		client = http.DefaultClient
	}

	cl := &Client{
		cfg:            cfg,
		tokenGenerator: tokenGenerator,
		client:         client,
	}

	return cl
}
