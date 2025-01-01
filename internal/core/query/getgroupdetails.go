package query

import (
	"context"
	"errors"

	"github.com/maximekuhn/partage/internal/core/query/queryutils"
	"github.com/maximekuhn/partage/internal/core/store"
	"github.com/maximekuhn/partage/internal/core/valueobject"
)

type GetGroupDetailsQuery struct {
	GroupID valueobject.GroupID
}

type GetGroupDetailsQueryHandler struct {
	groupstore   store.GroupStore
	userstore    store.UserStore
	expensestore store.ExpenseStore
}

func NewGetGroupDetailsQueryHandler(groupstore store.GroupStore, userstore store.UserStore, expensestore store.ExpenseStore) *GetGroupDetailsQueryHandler {
	return &GetGroupDetailsQueryHandler{groupstore, userstore, expensestore}
}

func (h *GetGroupDetailsQueryHandler) Handle(ctx context.Context, query GetGroupDetailsQuery) (*GroupDetails, bool, error) {
	// note: this is not the most efficient implementation, as we will perform multiple queries.
	// In the future, we might use an aggregation to query directly a group with all members (and expenses, etc...)
	// For now, this if fine.

	g, found, err := h.groupstore.FindByID(ctx, query.GroupID)
	if err != nil {
		return nil, false, err
	}
	if !found {
		return nil, false, nil
	}

	members, err := h.userstore.SelectAllInGroup(ctx, query.GroupID)
	if err != nil {
		return nil, false, err
	}

	expenses, err := h.expensestore.GetAllForGroup(ctx, query.GroupID)
	if err != nil {
		return nil, false, err
	}

	owner := queryutils.GetGroupOwner(g, members)
	if owner == nil {
		return nil, false, errors.New("could not find group owner")
	}

	return NewGroupDetails(g.Name, members, *owner, g.CreatedAt, expenses), true, nil
}
