package sqlite

import (
	"context"
	"database/sql"

	"github.com/maximekuhn/partage/internal/auth"
)

type SQLiteAuthStore struct {
	db *sql.DB
}

func NewSQLiteAuthStore(db *sql.DB) *SQLiteAuthStore {
	return &SQLiteAuthStore{db}
}

func (s SQLiteAuthStore) Save(ctx context.Context, data auth.AuthData) error {
	query := `
    INSERT INTO auth (user_id, hashed_password) VALUES (?, ?)
    `
	_, err := s.db.ExecContext(ctx, query, data.UserID.String(), data.HashedPassword.Hash())
	return err
}
