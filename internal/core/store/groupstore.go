package store

import (
	"context"

	"github.com/maximekuhn/partage/internal/core/entity"
	"github.com/maximekuhn/partage/internal/core/valueobject"
)

type GroupStore interface {
	Save(ctx context.Context, g *entity.Group) error
	FindByName(ctx context.Context, name valueobject.Groupname) (*entity.Group, bool, error)
}
