package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/maximekuhn/partage/internal/core/valueobject"
)

type SQLiteInvitationStore struct {
	db *sql.DB
}

func NewSQLiteInvitationStore(db *sql.DB) *SQLiteInvitationStore {
	return &SQLiteInvitationStore{db}
}

func (s *SQLiteInvitationStore) Save(ctx context.Context, i valueobject.Invitation) error {
	query := `
    INSERT INTO invitation (group_id, user_id, created_at, updated_at) VALUES (?, ?, ?, ?)
    `
	_, err := s.db.ExecContext(ctx, query, i.GroupID.String(), i.UserID.String(), i.CreatedAt, i.UpdatedAt)
	return err
}

func (s *SQLiteInvitationStore) Update(ctx context.Context, i valueobject.Invitation) error {
	return errors.New("not implemented")
}

func (s *SQLiteInvitationStore) FindByInviteeID(
	ctx context.Context,
	inviteeID valueobject.UserID,
	groupID valueobject.GroupID,
) (valueobject.Invitation, bool, error) {
	return valueobject.Invitation{}, false, errors.New("not implemented")
}
