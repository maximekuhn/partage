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
