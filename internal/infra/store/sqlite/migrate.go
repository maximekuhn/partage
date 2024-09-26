package sqlite

import (
	"context"
	"database/sql"
	"time"
)

// ApplyMigrations creates all SQLite tables and apply configurations.
// An internal [context.Context] is used with a timeout of 10 seconds.
func ApplyMigrations(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := createUserTable(ctx, db); err != nil {
		return err
	}

	if err := createAuthTable(ctx, db); err != nil {
		return err
	}

	return nil
}

func createUserTable(ctx context.Context, db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS user (
        id TEXT PRIMARY KEY,
        nickname TEXT NOT NULL,
        email TEXT NOT NULL,
        created_at DATE NOT NULL,
        UNIQUE(email, nickname)
    )
    `
	_, err := db.ExecContext(ctx, query)
	return err
}

func createAuthTable(ctx context.Context, db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS auth (
        user_id TEXT PRIMARY KEY,
        hashed_password BLOB NOT NULL
    )
    `
	_, err := db.ExecContext(ctx, query)
	return err
}
