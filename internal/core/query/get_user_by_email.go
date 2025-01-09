package query

import (
	"context"

	"github.com/maximekuhn/partage/internal/core/entity"
	"github.com/maximekuhn/partage/internal/core/store"
	"github.com/maximekuhn/partage/internal/core/valueobject"
)

type GetUserByEmailQuery struct {
	Email valueobject.Email
}

type GetUserByEmailQueryHandler struct {
	s store.UserStore
}

func NewGetUserByEmailCommandHandler(s store.UserStore) *GetUserByEmailQueryHandler {
	return &GetUserByEmailQueryHandler{s}
}

func (h *GetUserByEmailQueryHandler) Handle(ctx context.Context, cmd GetUserByEmailQuery) (*entity.User, bool, error) {
	return h.s.GetByEmail(ctx, cmd.Email)
}
