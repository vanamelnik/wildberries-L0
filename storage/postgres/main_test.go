package postgres

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	pgUser        = "postgres"
	pgPassword    = "secret"
	pgDatabase    = "wildberries-l0-mock"
	containerPort = "5432/tcp"
)

var (
	pgMockStorage *Storage
)

type postgresContainer struct {
	testcontainers.Container
	mappedPort string
}

func (c postgresContainer) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", pgUser, pgPassword, c.mappedPort, pgDatabase)
}

func TestMain(m *testing.M) {
	ctx := context.Background()
	pgContainer, err := setupPostgresContainer(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer pgContainer.Terminate(ctx)
	if err := pgContainer.Start(ctx); err != nil {
		log.Fatal(err)
	}
	log.Printf("DSN: %s", pgContainer.GetDSN()) //TODO: remove debug
	pg, err := NewStorage(pgContainer.GetDSN())
	if err != nil {
		log.Fatal(err)
	}
	pgMockStorage = pg
	m.Run()
}

func setupPostgresContainer(ctx context.Context) (*postgresContainer, error) {
	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "postgres:14.1",
			Env: map[string]string{
				"POSTGRES_USER":     pgUser,
				"POSTGRES_PASSWORD": pgPassword,
				"POSTGRES_DB":       pgDatabase,
			},
			ExposedPorts: []string{containerPort},
			WaitingFor: wait.ForSQL(
				nat.Port(containerPort),
				"pgx",
				func(p nat.Port) string {
					return fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", pgUser, pgPassword, p.Port(), pgDatabase)
				},
			).Timeout(30 * time.Second),
			Name:       "postgres-mock",
			SkipReaper: true,
		},
		Started: true,
	}
	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return nil, err
	}
	log.Println("test conainer created")
	mappedPort, err := container.MappedPort(ctx, nat.Port(containerPort))
	if err != nil {
		return nil, err
	}
	return &postgresContainer{
		Container:  container,
		mappedPort: mappedPort.Port(),
	}, nil
}

func cleanOrdersTable(t *testing.T) {
	_, err := pgMockStorage.db.Exec(`DELETE FROM orders;`)
	require.NoErrorf(t, err, "could not delete all records from the orders table")
}
