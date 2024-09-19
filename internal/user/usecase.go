package user

import (
	"context"

	"github.com/maximekuhn/partage/internal/common/valueobjects"
)

type Service struct {
	s   Store
	idp IDProvider
	dtp DatetimeProvider
}

func NewService(s Store, idp IDProvider, dtp DatetimeProvider) *Service {
	return &Service{s, idp, dtp}
}

// CreateUser returns the ID of the newly created user.
// If an error is returned, the ID must be considered invalid.
func (s *Service) CreateUser(
	ctx context.Context,
	nick Nickname,
	email valueobjects.Email,
) (ID, error) {
	id := s.idp.Provide()
	createdAt := s.dtp.Provide()
	u := NewUser(id, nick, email, createdAt)

	err := s.s.Save(ctx, u)
	if err != nil {
		return id, nil
	}

	return id, nil
}
