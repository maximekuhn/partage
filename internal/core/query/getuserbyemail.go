package query

import (
	"context"

	"github.com/maximekuhn/partage/internal/core/entity"
	"github.com/maximekuhn/partage/internal/core/store"
	"github.com/maximekuhn/partage/internal/core/valueobject"
)

type GetUserByEmailCommand struct {
	Email valueobject.Email
}

type GetUserByEmailCommandHandler struct {
	s store.UserStore
}

func NewGetUserByEmailCommandHandler(s store.UserStore) *GetUserByEmailCommandHandler {
	return &GetUserByEmailCommandHandler{s}
}

func (h *GetUserByEmailCommandHandler) Handle(ctx context.Context, cmd GetUserByEmailCommand) (*entity.User, bool, error) {
	return h.s.GetByEmail(ctx, cmd.Email)
}
