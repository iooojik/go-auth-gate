package session

type Config struct {
	SqlDsn string `yaml:"sqlDsn"`
}

type Repository struct {
	client Conn
}

func New(client Conn) *Repository {
	r := &Repository{
		client: client,
	}

	return r
}
