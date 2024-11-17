package entity

import (
	"time"

	"github.com/maximekuhn/partage/internal/core/valueobject"
)

type Expense struct {
	ID           valueobject.ExpenseID
	Label        valueobject.ExpenseLabel
	PaidBy       valueobject.UserID
	Participants []valueobject.UserID // can contain the payer
	Amount       valueobject.Amount
	CreatedAt    time.Time
	GroupID      valueobject.GroupID
}

func NewExpense(
	id valueobject.ExpenseID,
	label valueobject.ExpenseLabel,
	paidBy valueobject.UserID,
	participants []valueobject.UserID,
	amount valueobject.Amount,
	createdAt time.Time,
	GroupID valueobject.GroupID,
) *Expense {
	return &Expense{id, label, paidBy, participants, amount, createdAt.UTC(), GroupID}
}
