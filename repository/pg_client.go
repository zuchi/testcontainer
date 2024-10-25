package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type PGClient struct {
	conn *pgx.Conn
}

func NewPGClient(ctx context.Context, username, password, host, database string, port int) *PGClient {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", username, password, host, port, database)
	connect, err := pgx.Connect(ctx, connStr)
	if err != nil {
		fmt.Printf("Unable to connect to database: %v\n", err)
		return nil
	}

	return &PGClient{conn: connect}
}

func (c *PGClient) GetConn() *pgx.Conn {
	return c.conn
}

func (c *PGClient) Close(ctx context.Context) error {
	return c.conn.Close(ctx)
}
