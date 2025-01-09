package valueobject

import "github.com/google/uuid"

type GroupID struct {
	id uuid.UUID
}

func NewGroupID(id uuid.UUID) (GroupID, error) {
	i := GroupID{id}
	return i, nil
}

func NewGroupIDFromString(id string) (GroupID, error) {
	i := GroupID{uuid.Max}
	gid, err := uuid.Parse(id)
	if err != nil {
		return i, nil
	}
	return NewGroupID(gid)
}

func (i GroupID) String() string {
	return i.id.String()
}
