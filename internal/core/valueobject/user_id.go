package valueobject

import "github.com/google/uuid"

type UserID struct {
	id uuid.UUID
}

func NewUserID(id uuid.UUID) (UserID, error) {
	u := UserID{id}
	return u, nil
}

func NewUserIDFromString(id string) (UserID, error) {
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return UserID{}, err
	}
	return NewUserID(idUUID)

}

func (id UserID) String() string {
	return id.id.String()
}
