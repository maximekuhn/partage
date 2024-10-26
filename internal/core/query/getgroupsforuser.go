package query

import (
	"context"

	"github.com/maximekuhn/partage/internal/core/entity"
	"github.com/maximekuhn/partage/internal/core/store"
	"github.com/maximekuhn/partage/internal/core/valueobject"
)

type GetGroupsForUserQuery struct {
	UserID valueobject.UserID
}

type GetGroupsForUserQueryHandler struct {
	s store.GroupStore
}

func NewGetGroupsForUserQueryHandler(s store.GroupStore) *GetGroupsForUserQueryHandler {
	return &GetGroupsForUserQueryHandler{s}
}

func (h *GetGroupsForUserQueryHandler) Handle(ctx context.Context, query GetGroupsForUserQuery) ([]entity.Group, error) {
	return h.s.FindAllForUserID(ctx, query.UserID)
}
