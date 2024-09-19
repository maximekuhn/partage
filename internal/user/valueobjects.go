package user

import "github.com/google/uuid"

type ID struct {
	id uuid.UUID
}

func NewID(id uuid.UUID) (ID, error) {
	i := ID{id}
	return i, nil
}

func (id ID) String() string {
	return id.id.String()
}

type Nickname struct {
	nickname string
}

func NewNickname(nickname string) (Nickname, error) {
	n := Nickname{nickname}
	return n, nil
}

func (n Nickname) String() string {
	return n.nickname
}
