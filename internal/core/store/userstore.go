package store

import (
	"context"
	"errors"

	"github.com/maximekuhn/partage/internal/core/entity"
	"github.com/maximekuhn/partage/internal/core/valueobject"
)

var (
	ErrUserStoreDuplicate = errors.New("another user with the same (id, nickname, email) already exists")
)

type UserStore interface {
	Save(ctx context.Context, u *entity.User) error
	GetByID(ctx context.Context, id valueobject.UserID) (*entity.User, bool, error)
}
