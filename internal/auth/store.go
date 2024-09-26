package auth

import (
	"context"

	"github.com/maximekuhn/partage/internal/core/valueobject"
)

type AuthData struct {
	HashedPassword HashedPassword
	UserID         valueobject.UserID
}

type AuthStore interface {
	Save(ctx context.Context, data AuthData) error
}
