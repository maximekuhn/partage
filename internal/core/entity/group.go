package entity

import (
	"time"

	"github.com/maximekuhn/partage/internal/core/valueobject"
)

type Group struct {
	ID        valueobject.GroupID
	Name      valueobject.Groupname
	Members   []valueobject.UserID
	Owner     valueobject.UserID
	CreatedAt time.Time
}

func NewGroup(
	id valueobject.GroupID,
	name valueobject.Groupname,
	members []valueobject.UserID,
	owner valueobject.UserID,
	createdAt time.Time,
) *Group {
	return &Group{id, name, members, owner, createdAt.UTC()}
}

// ContainsUser returns true if the user is either the owner or a member
func (g *Group) ContainsUser(userID valueobject.UserID) bool {
	if userID == g.Owner {
		return true
	}
	for _, member := range g.Members {
		if member == userID {
			return true
		}
	}
	return false
}
