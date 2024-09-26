package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/maximekuhn/partage/internal/core/entity"
	"github.com/maximekuhn/partage/internal/core/valueobject"
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

func (s *SQLiteUserStore) Save(ctx context.Context, u *entity.User) error {
	query := `
    INSERT INTO user (id, nickname, email, created_at) VALUES (?, ?, ?, ?)
    `
	_, err := s.db.ExecContext(ctx, query, u.ID.String(), u.Nickname.String(), u.Email.String(), u.CreatedAt)
	return err
}

func (s *SQLiteUserStore) GetByID(ctx context.Context, id valueobject.UserID) (*entity.User, bool, error) {
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
	id, err = valueobject.NewUserID(idUUID)
	if err != nil {
		return nil, false, err
	}
	nn, err := valueobject.NewNickname(su.nickname)
	if err != nil {
		return nil, false, err
	}
	em, err := valueobject.NewEmail(su.email)
	if err != nil {
		return nil, false, err
	}

	u := entity.NewUser(id, em, nn, su.createdAt)

	return u, true, nil
}
