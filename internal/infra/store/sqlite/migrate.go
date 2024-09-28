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

	if err := createGroupTable(ctx, db); err != nil {
		return err
	}

	if err := createGroupUserAssociationTable(ctx, db); err != nil {
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

func createGroupTable(ctx context.Context, db *sql.DB) error {
	// NOTE: we are using 'partage_group' as table name because
	// 'group' is a reserved SQL keyword
	query := `
    CREATE TABLE IF NOT EXISTS partage_group (
        id TEXT PRIMARY KEY,
        name TEXT NOT NULL,
        owner TEXT NOT NULL,
        created_at DATE NOT NULL,
        UNIQUE(id, name)
    )
    `
	_, err := db.ExecContext(ctx, query)
	return err
}

func createGroupUserAssociationTable(ctx context.Context, db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS partage_group_user (
        group_id TEXT PRIMARY KEY,
        user_id TEXT NOT NULL
    )
    `
	_, err := db.ExecContext(ctx, query)
	return err
}
