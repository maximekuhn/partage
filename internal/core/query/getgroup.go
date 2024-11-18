package query

import (
	"context"

	"github.com/maximekuhn/partage/internal/core/entity"
	"github.com/maximekuhn/partage/internal/core/store"
	"github.com/maximekuhn/partage/internal/core/valueobject"
)

type GetGroupQuery struct {
	GroupID valueobject.GroupID
}

type GetGroupQueryHandler struct {
	s store.GroupStore
}

func NewGetGroupQueryHandler(s store.GroupStore) *GetGroupQueryHandler {
	return &GetGroupQueryHandler{s}
}

func (h *GetGroupQueryHandler) Handle(ctx context.Context, query GetGroupQuery) (*entity.Group, bool, error) {
	return h.s.FindByID(ctx, query.GroupID)
}
