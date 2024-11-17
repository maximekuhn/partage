package store

import (
	"context"

	"github.com/maximekuhn/partage/internal/core/entity"
)

type ExpenseStore interface {
	Save(ctx context.Context, e *entity.Expense) error
}
