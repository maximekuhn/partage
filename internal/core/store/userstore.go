package store

import (
	"context"

	"github.com/maximekuhn/partage/internal/core/entity"
	"github.com/maximekuhn/partage/internal/core/valueobject"
)

type UserStore interface {
	Save(ctx context.Context, u *entity.User) error
	GetByID(ctx context.Context, id valueobject.UserID) (*entity.User, bool, error)
}
