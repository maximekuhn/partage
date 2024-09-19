package sqlite

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// CreateTmpDB returns a temporary db handler and apply migrations.
//
// If something wrong happens, the program crashes.
//
// The temporary db file will be deleted according to the OS settings.
func CreateTmpDB() *sql.DB {
	f, err := os.CreateTemp("", "test-db-*.sqlite3")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	db, err := sql.Open("sqlite3", f.Name())
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	if err = ApplyMigrations(db); err != nil {
		func() { _ = db.Close() }()
		panic(err)
	}

	return db
}
