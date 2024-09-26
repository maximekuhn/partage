package valueobject

import "github.com/google/uuid"

type UserID struct {
	id uuid.UUID
}

func NewUserID(id uuid.UUID) (UserID, error) {
	u := UserID{id}
	return u, nil
}

func (id UserID) String() string {
	return id.id.String()
}
