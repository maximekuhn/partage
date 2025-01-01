package query

import (
	"time"

	"github.com/maximekuhn/partage/internal/core/entity"
	"github.com/maximekuhn/partage/internal/core/valueobject"
)

// GroupDetails represents a group with associated details
type GroupDetails struct {
	GroupName valueobject.Groupname
	Members   []*entity.User
	Owner     entity.User
	CreatedAt time.Time
	Expenses  []*entity.Expense
}

func NewGroupDetails(
	groupName valueobject.Groupname,
	members []*entity.User,
	owner entity.User,
	createdAt time.Time,
	expenses []*entity.Expense,
) *GroupDetails {
	return &GroupDetails{
		groupName,
		members,
		owner,
		createdAt.UTC(),
		expenses,
	}
}

// TotalAmount computes the total amount of all expenses in the group.
// This is a temporary function that will be moved to a command at some point.
// It assumes everything is in EURO.
func (g *GroupDetails) TotalAmount() float64 {
	total := 0.0
	for _, e := range g.Expenses {
		total += e.Amount.Amount()
	}
	return total
}
