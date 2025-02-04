package session

type Config struct {
	SessionDuration int64 `yaml:"sessionDuration"`
}

type Repository struct {
	cfg    Config
	client Conn
}

func New(cfg Config, client Conn) *Repository {
	r := &Repository{
		cfg:    cfg,
		client: client,
	}

	return r
}
