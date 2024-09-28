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

type SQLiteGroupStore struct {
	db *sql.DB
}

func NewSQLiteGroupStore(db *sql.DB) *SQLiteGroupStore {
	return &SQLiteGroupStore{db}
}

func (s *SQLiteGroupStore) FindByName(ctx context.Context, name valueobject.Groupname) (*entity.Group, bool, error) {
	// TODO: handle members
	query := `
    SELECT id, owner, created_at FROM partage_group WHERE name = ?
    `
	var id uuid.UUID
	var owner uuid.UUID
	var created_at time.Time

	err := s.db.QueryRowContext(ctx, query, name.String()).Scan(&id, &owner, &created_at)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, err
	}

	groupID, err := valueobject.NewGroupID(id)
	if err != nil {
		return nil, false, err
	}

	ownerID, err := valueobject.NewUserID(owner)
	if err != nil {
		return nil, false, err
	}

	membersQuery := `
    SELECT user_id FROM partage_group_user WHERE group_id = ?
    `
	rows, err := s.db.QueryContext(ctx, membersQuery, groupID.String())
	defer func() { _ = rows.Close() }()

	if err != nil {
		return nil, false, err
	}

	members := make([]valueobject.UserID, 0)
	for rows.Next() {
		var userID uuid.UUID
		if err := rows.Scan(&userID); err != nil {
			return nil, false, err
		}

		memberID, err := valueobject.NewUserID(userID)
		if err != nil {
			return nil, false, err
		}

		members = append(members, memberID)
	}

	g := entity.NewGroup(groupID, name, members, ownerID, created_at)
	return g, true, nil
}

func (s *SQLiteGroupStore) Save(ctx context.Context, g *entity.Group) error {
	query := `
    INSERT INTO partage_group (id, name, owner, created_at) VALUES (?, ?, ?, ?)
    `
	_, err := s.db.ExecContext(ctx, query, g.ID.String(), g.Name.String(), g.Owner.String(), g.CreatedAt)

	if err != nil {
		return err
	}

	memberQuery := `
    INSERT INTO partage_group_user (group_id, user_id) VALUES (?, ?)
    `
	for _, memberID := range g.Members {
		_, err := s.db.ExecContext(ctx, memberQuery, g.ID.String(), memberID.String())
		if err != nil {
			return err
		}
	}
	return err
}
