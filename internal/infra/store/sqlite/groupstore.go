package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"strings"
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
	query := `
    SELECT pg.id,
           pg.owner, 
           pg.created_at,
           GROUP_CONCAT(pgu.user_id)
    FROM partage_group pg
    INNER JOIN partage_group_user pgu ON pg.id = pgu.group_id
    WHERE name = ?
    GROUP BY pg.id, pg.owner
    `
	var id uuid.UUID
	var owner uuid.UUID
	var created_at time.Time
	var membersConcat string

	err := s.db.QueryRowContext(ctx, query, name.String()).Scan(&id, &owner, &created_at, &membersConcat)
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

	memberIDs := strings.Split(membersConcat, ",")
	members := make([]valueobject.UserID, 0)
	for _, memberID := range memberIDs {
		id, err := valueobject.NewUserIDFromString(memberID)
		if err != nil {
			return nil, false, err
		}

		// ignore owner
		if id == ownerID {
			continue
		}
		members = append(members, id)
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
	query := `
    SELECT pg.id,
           pg.owner, 
           pg.created_at,
           pg.name, 
           GROUP_CONCAT(pgu.user_id)
    FROM partage_group pg
    INNER JOIN partage_group_user pgu ON pg.id = pgu.group_id
    WHERE pgu.user_id = ?
    GROUP BY pg.id, pg.owner
    `
	rows, err := s.db.QueryContext(ctx, query, userID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := make([]entity.Group, 0)
	for rows.Next() {
		var id uuid.UUID
		var owner uuid.UUID
		var created_at time.Time
		var name string
		var membersConcat string
		err = rows.Scan(&id, &owner, &created_at, &name, &membersConcat)
		if err != nil {
			return nil, err
		}
		groupID, err := valueobject.NewGroupID(id)
		if err != nil {
			return nil, err
		}

		ownerID, err := valueobject.NewUserID(owner)
		if err != nil {
			return nil, err
		}

		groupName, err := valueobject.NewGroupname(name)
		if err != nil {
			return nil, err
		}

		memberIDs := strings.Split(membersConcat, ",")
		members := make([]valueobject.UserID, 0)
		for _, memberID := range memberIDs {
			id, err := valueobject.NewUserIDFromString(memberID)
			if err != nil {
				return nil, err
			}

			// ignore owner
			if id == ownerID {
				continue
			}
			members = append(members, id)
		}

		g := entity.NewGroup(groupID, groupName, members, ownerID, created_at)
		groups = append(groups, *g)
	}

	return groups, nil
}

func (s *SQLiteGroupStore) FindByID(ctx context.Context, groupID valueobject.GroupID) (*entity.Group, bool, error) {
	query := `
    SELECT pg.name,
           pg.owner, 
           pg.created_at,
           GROUP_CONCAT(pgu.user_id)
    FROM partage_group pg
    INNER JOIN partage_group_user pgu ON pg.id = pgu.group_id
    WHERE pgu.group_id = ?
    GROUP BY pg.id, pg.owner
    `
	var name string
	var owner uuid.UUID
	var created_at time.Time
	var membersConcat string

	err := s.db.QueryRowContext(ctx, query, groupID.String()).Scan(&name, &owner, &created_at, &membersConcat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, err
	}

	ownerID, err := valueobject.NewUserID(owner)
	if err != nil {
		return nil, false, err
	}

	groupname, err := valueobject.NewGroupname(name)
	if err != nil {
		return nil, false, err
	}

	memberIDs := strings.Split(membersConcat, ",")
	members := make([]valueobject.UserID, 0)
	for _, memberID := range memberIDs {
		id, err := valueobject.NewUserIDFromString(memberID)
		if err != nil {
			return nil, false, err
		}

		// ignore owner
		if id == ownerID {
			continue
		}
		members = append(members, id)
	}

	g := entity.NewGroup(groupID, groupname, members, ownerID, created_at)
	return g, true, nil
}

func (s *SQLiteGroupStore) FindAllUsersInGroup(ctx context.Context, groupID valueobject.GroupID) ([]entity.Group, error) {
	return nil, errors.New("TODO: implement me")
}
