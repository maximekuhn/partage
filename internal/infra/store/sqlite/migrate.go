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

	return createUserTable(db, ctx)
}

func createUserTable(db *sql.DB, ctx context.Context) error {
	query := `
    CREATE TABLE IF NOT EXISTS user (
        id TEXT PRIMARY KEY,
        nickname TEXT NOT NULL,
        email TEXT NOT NULL,
        created_at DATE NOT NULL,
        UNIQUE(id, nickname)
    )
    `
	_, err := db.ExecContext(ctx, query)
	return err
}
