package google

type Config struct {
	URL string
}

type Client struct {
	client HTTPClient
	cfg    Config
}

func New() *Client {
	c := &Client{}

	return c
}
