package contract

import "github.com/maximekuhn/partage/internal/core/valueobject"

type ExpenseIDProvider interface {
	Provide() valueobject.ExpenseID
}
