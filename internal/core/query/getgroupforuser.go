package query

import (
	"context"

	"github.com/maximekuhn/partage/internal/core/entity"
	"github.com/maximekuhn/partage/internal/core/store"
	"github.com/maximekuhn/partage/internal/core/valueobject"
)

type GetGroupForUserQuery struct {
	UserID valueobject.UserID
}

type GetGroupForUserQueryHandler struct {
	s store.GroupStore
}

func NewGetGroupForUserQueryHandler(s store.GroupStore) *GetGroupForUserQueryHandler {
	return &GetGroupForUserQueryHandler{s}
}

func (h *GetGroupForUserQueryHandler) Handle(ctx context.Context, query GetGroupForUserQuery) ([]entity.Group, error) {
	return h.s.FindAllForUserID(ctx, query.UserID)
}
