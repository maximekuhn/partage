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

	g := entity.NewGroup(groupID, name, make([]valueobject.UserID, 0), ownerID, created_at)
	return g, true, nil
}

func (s *SQLiteGroupStore) Save(ctx context.Context, g *entity.Group) error {
	// TODO: handle members
	query := `
    INSERT INTO partage_group (id, name, owner, created_at) VALUES (?, ?, ?, ?)
    `
	_, err := s.db.ExecContext(ctx, query, g.ID.String(), g.Name.String(), g.Owner.String(), g.CreatedAt)

	return err
}
