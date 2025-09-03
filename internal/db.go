package internal

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func connectDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func executeQuery(ctx context.Context, db *sql.DB, query string) (*sql.Rows, error) {
	// Use context-aware query execution
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
