package entity

import (
	"time"

	"github.com/maximekuhn/partage/internal/core/valueobject"
)

type User struct {
	ID        valueobject.UserID
	Email     valueobject.Email
	Nickname  valueobject.Nickname
	CreatedAt time.Time
}

func NewUser(
	id valueobject.UserID,
	email valueobject.Email,
	nickname valueobject.Nickname,
	createdAt time.Time,
) *User {
	return &User{id, email, nickname, createdAt.UTC()}
}
