package query

import (
	"context"
	"errors"

	"github.com/maximekuhn/partage/internal/core/entity"
	"github.com/maximekuhn/partage/internal/core/store"
	"github.com/maximekuhn/partage/internal/core/valueobject"
)

type GetGroupDetailsQuery struct {
	GroupID valueobject.GroupID
}

type GetGroupDetailsQueryOutput struct {
	Group   *entity.Group
	Members map[valueobject.UserID]*entity.User // guaranteed to contain all members of the group
}

type GetGroupDetailsQueryHandler struct {
	groupstore store.GroupStore
}

func NewGetGroupDetailsQueryHandler(groupstore store.GroupStore) *GetGroupDetailsQueryHandler {
	return &GetGroupDetailsQueryHandler{groupstore}
}

func (h *GetGroupDetailsQueryHandler) Handle(ctx context.Context, query GetGroupDetailsQuery) (*GetGroupDetailsQueryOutput, bool, error) {
	// note: this is not the most efficient implementation, as we will perform multiple queries.
	// In the future, we might use a aggregation to query directly a group with all members (and expenses, etc...)
	// For now, this if fine.

	_, found, err := h.groupstore.FindByID(ctx, query.GroupID)
	if err != nil {
		return nil, false, err
	}
	if !found {
		return nil, false, nil
	}

	return nil, false, errors.New("TODO: implement me")
}
