package user

import (
	"context"
	"errors"
	"time"
)

var (
	ErrStoreDuplicate = errors.New("another user with the same ID or nickname already exists")
)

type Store interface {
	Save(ctx context.Context, u *User) error
	GetByID(ctx context.Context, id ID) (*User, bool, error)
}

type IDProvider interface {
	Provide() ID
}

type DatetimeProvider interface {
	Provide() time.Time
}
