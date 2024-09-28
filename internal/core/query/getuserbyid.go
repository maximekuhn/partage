package query

import (
	"context"

	"github.com/maximekuhn/partage/internal/core/entity"
	"github.com/maximekuhn/partage/internal/core/store"
	"github.com/maximekuhn/partage/internal/core/valueobject"
)

type GetUserByIDCommand struct {
	ID valueobject.UserID
}

type GetUserByIDCommandHandler struct {
	s store.UserStore
}

func NewGetUserByIDCommandHandler(s store.UserStore) *GetUserByIDCommandHandler {
	return &GetUserByIDCommandHandler{s}
}

func (h *GetUserByIDCommandHandler) Handle(ctx context.Context, cmd GetUserByIDCommand) (*entity.User, bool, error) {
	return h.s.GetByID(ctx, cmd.ID)
}
