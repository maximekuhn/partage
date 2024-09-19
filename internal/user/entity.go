package user

import (
	"time"

	"github.com/maximekuhn/partage/internal/common/valueobjects"
)

type User struct {
	ID        ID
	Nick      Nickname
	Email     valueobjects.Email
	CreatedAt time.Time
}

func NewUser(
	id ID,
	nick Nickname,
	email valueobjects.Email,
	createdAt time.Time,
) *User {
	return &User{id, nick, email, createdAt.UTC()}
}
