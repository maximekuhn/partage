package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/maximekuhn/partage/internal/auth"
	"github.com/maximekuhn/partage/internal/core/valueobject"
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

func (s *SQLiteAuthStore) GetByUserID(ctx context.Context, userID valueobject.UserID) (*auth.AuthData, bool, error) {
	query := `
    SELECT hashed_password FROM auth WHERE user_id = ?
    `
	var hash []byte
	err := s.db.QueryRowContext(ctx, query, userID.String()).Scan(&hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, err
	}
	data := auth.AuthData{
		HashedPassword: auth.NewHashedPassword(hash),
		UserID:         userID,
	}
	return &data, true, nil
}
