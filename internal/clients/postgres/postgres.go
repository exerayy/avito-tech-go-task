package postgres

import (
	"context"
	"database/sql"
	"fmt"
)

type Client struct {
	pg *sql.DB
}

func Connect(dsn string) (*Client, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return &Client{pg: db}, nil
}

func (c *Client) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return c.pg.ExecContext(ctx, query, args...)
}

func (c *Client) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return c.pg.QueryContext(ctx, query, args...)
}

func (c *Client) Begin(ctx context.Context) (*sql.Tx, error) {
	return c.pg.Begin()
}

func (c *Client) Close() error {
	return c.pg.Close()
}
