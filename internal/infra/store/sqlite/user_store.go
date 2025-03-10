package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/maximekuhn/partage/internal/core/entity"
	"github.com/maximekuhn/partage/internal/core/store"
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
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint") {
			return store.ErrUserStoreDuplicate
		}
		return err
	}
	return nil
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

func (s *SQLiteUserStore) GetByEmail(ctx context.Context, email valueobject.Email) (*entity.User, bool, error) {
	query := `
    SELECT id, nickname, created_at FROM user WHERE email = ?
    `
	var su sqliteUser
	err := s.db.QueryRowContext(ctx, query, email.String()).Scan(&su.id, &su.nickname, &su.createdAt)
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
	id, err := valueobject.NewUserID(idUUID)
	if err != nil {
		return nil, false, err
	}
	nn, err := valueobject.NewNickname(su.nickname)
	if err != nil {
		return nil, false, err
	}

	u := entity.NewUser(id, email, nn, su.createdAt)

	return u, true, nil
}

func (s *SQLiteUserStore) SelectAllInGroup(ctx context.Context, groupID valueobject.GroupID) ([]*entity.User, error) {
	query := `
    SELECT id, nickname, email, created_at
    FROM user u
    INNER JOIN partage_group_user pgu ON u.id = pgu.user_id
    WHERE pgu.group_id = ?
    `

	rows, err := s.db.QueryContext(ctx, query, groupID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := make([]*entity.User, 0)
	for rows.Next() {
		var id uuid.UUID
		var nickname string
		var email string
		var createdAt time.Time

		err = rows.Scan(&id, &nickname, &email, &createdAt)
		if err != nil {
			return nil, err
		}

		u, err := tryConvertUser(id, nickname, email, createdAt)
		if err != nil {
			return nil, err
		}
		members = append(members, u)
	}
	return members, nil
}

func tryConvertUser(id uuid.UUID, nickname, email string, createdAt time.Time) (*entity.User, error) {
	userID, err := valueobject.NewUserID(id)
	if err != nil {
		return nil, err
	}
	userNickname, err := valueobject.NewNickname(nickname)
	if err != nil {
		return nil, err
	}
	userEmail, err := valueobject.NewEmail(email)
	if err != nil {
		return nil, err
	}
	return entity.NewUser(userID, userEmail, userNickname, createdAt), nil
}
