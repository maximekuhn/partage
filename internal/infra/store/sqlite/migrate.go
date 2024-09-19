package sqlite

import "database/sql"

func ApplyMigrations(db *sql.DB) error {
	return createUserTable(db)
}

func createUserTable(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS user (
        id TEXT PRIMARY KEY,
        nickname TEXT NOT NULL,
        email TEXT NOT NULL,
        created_at DATE NOT NULL,
        UNIQUE(id, nickname)
    )
    `
	_, err := db.Exec(query)
	return err
}
