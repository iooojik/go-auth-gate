package google

type Config struct {
	URL   string `yaml:"url"`
	AppID string `yaml:"appID"`
}

type Client struct {
	client HTTPClient
	cfg    Config
}

func New(cfg Config, client HTTPClient) *Client {
	c := &Client{
		client: client,
		cfg:    cfg,
	}

	return c
}
