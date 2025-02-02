package apple

import (
	"net/http"
)

type Config struct {
	tokenConfig TokenConfig `yaml:"token"`
	URL         string      `yaml:"url"`
	// KeyPath  string `yaml:"keyPath"`
	client HTTPClient `yaml:"-"`
}

type Client struct {
	cfg            Config
	tokenGenerator ClientSecretGenerator
}

// New creates a new instance of Client.
func New(cfg Config, tokenGenerator ClientSecretGenerator) *Client {
	if cfg.client == nil {
		cfg.client = http.DefaultClient
	}

	cl := &Client{
		cfg:            cfg,
		tokenGenerator: tokenGenerator,
	}

	return cl
}
