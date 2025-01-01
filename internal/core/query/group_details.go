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
