package misc

import (
	"github.com/google/uuid"
	"github.com/maximekuhn/partage/internal/core/valueobject"
)

type ExpenseIDProviderProd struct{}

func (p *ExpenseIDProviderProd) Provide() valueobject.ExpenseID {
	id, err := valueobject.NewExpenseID(uuid.New())
	if err != nil {
		panic(err)
	}
	return id
}
