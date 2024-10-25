package tests

import (
	"context"
	"errors"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"testcontainer/repository"
	"testcontainer/use_case"
)

func spinUpDependencies(t *testing.T, ctx context.Context) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:17",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       "userdb",
			"POSTGRES_USER":     "pgsql",
			"POSTGRES_PASSWORD": "root",
		},
		WaitingFor: wait.ForAll(wait.ForListeningPort("5432/tcp")).
			WithStartupTimeoutDefault(60 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		t.Fatalf("could not start container: %v", err)
	}

	return container, nil
}

func migrateData(username, password, host, databaseName string, port int) error {
	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", username, password, host, port, databaseName)
	pgxConfig, err := pgx.ParseConfig(databaseUrl)
	if err != nil {
		log.Fatalf("Unable to parse PostgreSQL config: %v\n", err)
	}

	stdlib.RegisterConnConfig(pgxConfig)

	db := stdlib.OpenDB(*pgxConfig)
	defer db.Close()

	m, err := migrate.New("file://../migrations", databaseUrl)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func Test_IntegrationLoginTest(t *testing.T) {
	ctx := context.Background()

	container, err := spinUpDependencies(t, ctx)
	if err != nil {
		t.Fatalf("could not spin up dependencies: %v", err)
	}
	defer container.Terminate(ctx)

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("could not get host container: %v", err)
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("could not get port container: %v", err)
	}

	pgClient := repository.NewPGClient(ctx, "pgsql", "root", host, "userdb", port.Int())
	if pgClient == nil {
		t.Fatalf("couldnt create client for postgres database")
	}

	err = migrateData("pgsql", "root", host, "userdb", port.Int())
	if err != nil {
		t.Fatalf("could not migrate data: %v", err)
	}

	userRepo := repository.NewUserRepository(pgClient.GetConn())
	insertUC := use_case.NewInsertUC(userRepo)
	err = insertUC.InsertNewUser(ctx, use_case.User{
		Email:    "ecore@ecore.com",
		Password: "123123",
	})
	if err != nil {
		t.Fatalf("could not insert user: %v", err)
	}

	loginUc := use_case.NewLoginUc(userRepo)
	login, err := loginUc.Login(ctx, "ecore@ecore.com", "123123")
	if err != nil {
		t.Fatalf("could not login: %v", err)
	}

	if login.Email != "ecore@ecore.com" {
		t.Fatalf("invalid login")
	}
}
