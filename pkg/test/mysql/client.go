package mysqltest

import (
	"context"
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var db *sqlx.DB

func SetupTestDB(t *testing.T, schema string) func() {
	ctx := context.Background()
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

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Could not start container: %s", err)
	}

	// Получаем порт
	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "3306")
	dsn := fmt.Sprintf("testuser:testpass@tcp(%s:%s)/testdb?parseTime=true", host, port)

	// Подключаемся к MySQL
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		t.Fatalf("Could not connect to database: %s", err)
	}

	// Создаем таблицы
	schema := `
	CREATE TABLE users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		apple_id VARCHAR(255) UNIQUE NOT NULL,
		session_duration INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		auth_type VARCHAR(50) NOT NULL
	);
	CREATE TABLE user_tokens (
		id INT AUTO_INCREMENT PRIMARY KEY,
		user_id INT NOT NULL,
		token TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE (user_id),
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);`
	_, err = db.Exec(schema)
	if err != nil {
		t.Fatalf("Could not create tables: %s", err)
	}

	return func() {
		_ = db.Close()
		_ = container.Terminate(ctx)
	}
}
