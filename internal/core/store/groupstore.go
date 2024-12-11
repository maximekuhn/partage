package store

import (
	"context"

	"github.com/maximekuhn/partage/internal/core/entity"
	"github.com/maximekuhn/partage/internal/core/valueobject"
)

type GroupStore interface {
	Save(ctx context.Context, g *entity.Group) error
	FindByID(ctx context.Context, groupID valueobject.GroupID) (*entity.Group, bool, error)
	FindByName(ctx context.Context, name valueobject.Groupname) (*entity.Group, bool, error)
	FindAllForUserID(ctx context.Context, userID valueobject.UserID) ([]entity.Group, error)
	FindAllUsersInGroup(ctx context.Context, groupID valueobject.GroupID) ([]entity.Group, error)
}
