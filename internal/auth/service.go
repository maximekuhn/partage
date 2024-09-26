package auth

import (
	"context"

	"github.com/maximekuhn/partage/internal/core/valueobject"
)

type AuthService struct {
	hasher *BcryptPasswordHasher
	store  AuthStore
}

func NewAuthService(hasher *BcryptPasswordHasher, store AuthStore) *AuthService {
	return &AuthService{hasher, store}
}

func (s *AuthService) Hash(p Password) (HashedPassword, error) {
	return s.hasher.hash(p)
}

func (s *AuthService) Verify(p Password, h HashedPassword) bool {
	return s.hasher.verify(p, h)
}

func (s *AuthService) Save(
	ctx context.Context,
	userID valueobject.UserID,
	h HashedPassword,
) error {
	data := AuthData{
		HashedPassword: h,
		UserID:         userID,
	}
	return s.store.Save(ctx, data)
}
