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

		if userID == owner {
			continue
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
	_, err = s.db.ExecContext(ctx, memberQuery, g.ID.String(), g.Owner.String())
	if err != nil {
		return err
	}
	for _, memberID := range g.Members {
		_, err := s.db.ExecContext(ctx, memberQuery, g.ID.String(), memberID.String())
		if err != nil {
			return err
		}
	}
	return err
}

func (s *SQLiteGroupStore) FindAllForUserID(ctx context.Context, userID valueobject.UserID) ([]entity.Group, error) {
	idsQuery :=
		`SELECT group_id FROM partage_group_user WHERE user_id = ?`
	rows, err := s.db.QueryContext(ctx, idsQuery, userID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groupsID := make([]valueobject.GroupID, 0)
	for rows.Next() {
		var groupID string
		if err := rows.Scan(&groupID); err != nil {
			return nil, err
		}
		gid, err := valueobject.NewGroupIDFromString(groupID)
		if err != nil {
			return nil, err
		}
		groupsID = append(groupsID, gid)
	}

	groups := make([]entity.Group, 0)
	for _, id := range groupsID {
		// TODO: handle members
		groupQuery :=
			`SELECT name, owner, created_at FROM partage_group WHERE id = ?`

		var name string
		var owner uuid.UUID
		var createdAt time.Time

		if err := s.db.QueryRowContext(ctx, groupQuery, id.String()).Scan(&name, &owner, &createdAt); err != nil {
			return nil, err
		}

		gName, err := valueobject.NewGroupname(name)
		if err != nil {
			return nil, err
		}

		gOwner, err := valueobject.NewUserID(owner)
		if err != nil {
			return nil, err
		}

		g := entity.NewGroup(id, gName, make([]valueobject.UserID, 0), gOwner, createdAt)

		groups = append(groups, *g)
	}

	return groups, nil
}
