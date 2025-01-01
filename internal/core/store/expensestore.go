package store

import (
	"context"

	"github.com/maximekuhn/partage/internal/core/entity"
	"github.com/maximekuhn/partage/internal/core/valueobject"
)

type ExpenseStore interface {
	Save(ctx context.Context, e *entity.Expense) error
	GetAllForGroup(ctx context.Context, groupID valueobject.GroupID) ([]*entity.Expense, error)
}
