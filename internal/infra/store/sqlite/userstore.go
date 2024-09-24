package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/maximekuhn/partage/internal/core/common/valueobjects"
	"github.com/maximekuhn/partage/internal/core/user"
)

type SQLiteUserStore struct {
	db *sql.DB
}

type sqliteUser struct {
	id        string
	nickname  string
	email     string
	createdAt time.Time `db:"created_at"`
}

func NewSQLiteUserStore(db *sql.DB) *SQLiteUserStore {
	return &SQLiteUserStore{db}
}

func (s *SQLiteUserStore) Save(ctx context.Context, u *user.User) error {
	query := `
    INSERT INTO user (id, nickname, email, created_at) VALUES (?, ?, ?, ?)
    `
	_, err := s.db.ExecContext(ctx, query, u.ID.String(), u.Nick.String(), u.Email.String(), u.CreatedAt)
	return err
}

func (s *SQLiteUserStore) GetByID(ctx context.Context, id user.ID) (*user.User, bool, error) {
	query := `
    SELECT id, nickname, email, created_at FROM user WHERE id = ?
    `
	var su sqliteUser
	err := s.db.QueryRowContext(ctx, query, id.String()).Scan(&su.id, &su.nickname, &su.email, &su.createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, err
	}

	idUUID, err := uuid.Parse(su.id)
	if err != nil {
		return nil, false, err
	}
	id, err = user.NewID(idUUID)
	if err != nil {
		return nil, false, err
	}
	nn, err := user.NewNickname(su.nickname)
	if err != nil {
		return nil, false, err
	}
	em, err := valueobjects.NewEmail(su.email)
	if err != nil {
		return nil, false, err
	}

	u := user.NewUser(id, nn, em, su.createdAt)

	return u, true, nil
}
