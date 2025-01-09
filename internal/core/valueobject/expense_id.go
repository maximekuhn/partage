package valueobject

import "github.com/google/uuid"

type ExpenseID struct {
	id uuid.UUID
}

func NewExpenseID(id uuid.UUID) (ExpenseID, error) {
	return ExpenseID{id}, nil
}

func (e ExpenseID) String() string {
	return e.id.String()
}
