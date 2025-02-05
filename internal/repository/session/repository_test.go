package session_test

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	"github.com/jmoiron/sqlx"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type RepositoryTestSuite struct {
	suite.Suite

	db        *sqlx.DB
	container testcontainers.Container
}

func (s *RepositoryTestSuite) SetupSuite() {
	ctx := context.Background()

	var err error

	req := testcontainers.ContainerRequest{
		Image:        "mysql:8.0",
		ExposedPorts: []string{"3306/tcp"},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "rootpass",
			"MYSQL_DATABASE":      "testdb",
			"MYSQL_USER":          "testuser",
			"MYSQL_PASSWORD":      "testpass",
		},
		WaitingFor: wait.ForLog("port: 3306  MySQL Community Server"),
	}

	s.Require().NoError(err)

	s.container, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	s.Require().NoError(err)

	host, _ := s.container.Host(ctx)
	port, _ := s.container.MappedPort(ctx, "3306")

	dsn := fmt.Sprintf("testuser:testpass@tcp(%s:%s)/testdb?parseTime=true&multiStatements=true", host, port.Port())

	slog.Info("connecting", "dsn", dsn)

	s.db, err = sqlx.Connect("mysql", dsn)
	s.Require().NoError(err)
}

func (s *RepositoryTestSuite) SetupTest() {
	ctx := context.Background()

	schemas := [2]string{
		`create table users
(
    id         int auto_increment
        primary key,
    user_id    char(255)                 not null,
    created_at timestamp default CURRENT_TIMESTAMP not null
);`,
		`create table user_tokens
(
    id            int auto_increment
        primary key,
    user_id       char(255)                           not null,
    created_at    timestamp default CURRENT_TIMESTAMP not null,
    access_token  text                                not null,
    token_type    char(128)                           null,
    expires_in    int       default 3600              not null,
    refresh_token text                                not null,
    id_token      text                                not null,
    constraint user_tokens_pk
        unique (user_id)
);`,
	}

	for _, schema := range schemas {
		_, err := s.db.ExecContext(ctx, schema)
		s.Require().NoError(err)
	}
}

func (s *RepositoryTestSuite) TearDownTest() {
	// remove table.
}

func (s *RepositoryTestSuite) TearDownSuite() {
	s.Require().NoError(s.db.Close())
	s.Require().NoError(s.container.Terminate(context.Background()))
}

func TestRepositoryTestSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(RepositoryTestSuite))
}
