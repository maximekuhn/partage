package query

import (
	"context"

	"github.com/maximekuhn/partage/internal/core/entity"
	"github.com/maximekuhn/partage/internal/core/store"
	"github.com/maximekuhn/partage/internal/core/valueobject"
)

type GetUserByIDQuery struct {
	ID valueobject.UserID
}

type GetUserByIDQueryHandler struct {
	s store.UserStore
}

func NewGetUserByIDCommandHandler(s store.UserStore) *GetUserByIDQueryHandler {
	return &GetUserByIDQueryHandler{s}
}

func (h *GetUserByIDQueryHandler) Handle(ctx context.Context, cmd GetUserByIDQuery) (*entity.User, bool, error) {
	return h.s.GetByID(ctx, cmd.ID)
}
