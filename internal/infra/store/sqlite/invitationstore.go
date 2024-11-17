package sqlite

import (
	"context"
	"database/sql"

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
