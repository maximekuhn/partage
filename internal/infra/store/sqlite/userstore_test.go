package sqlite

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/maximekuhn/partage/internal/core/common/valueobjects"
	"github.com/maximekuhn/partage/internal/core/user"
)

func TestSave(t *testing.T) {
	db := CreateTmpDB()
	defer db.Close()

	s := NewSQLiteUserStore(db)

	u := createUser("toto", "toto@gmail.com")

	err := s.Save(context.TODO(), u)
	if err != nil {
		t.Errorf("Save(): expected no error got %v", err)
	}
}

func TestGet(t *testing.T) {
	db := CreateTmpDB()
	defer db.Close()

	s := NewSQLiteUserStore(db)

	u := createUser("toto", "toto@gmail.com")
	_ = s.Save(context.TODO(), u)

	uFound, found, err := s.GetByID(context.TODO(), u.ID)
	if err != nil {
		t.Errorf("GetByID(): expected no error got %v", err)
		return
	}
	if !found {
		t.Error("GetByID(): expected to find user got nothing")
		return
	}
	if !reflect.DeepEqual(u, uFound) {
		t.Errorf("GetByID(): expected %v got %v", u, uFound)
	}
}

// createUser returns a pointer to a User.
// Panics if one field could not be transformed into a valid value objects
func createUser(nickname, email string) *user.User {
	id, err := user.NewID(uuid.New())
	if err != nil {
		panic(err)
	}
	nn, err := user.NewNickname(nickname)
	if err != nil {
		panic(err)
	}
	em, err := valueobjects.NewEmail(email)
	if err != nil {
		panic(err)
	}
	createdAt := time.Now()
	return user.NewUser(id, nn, em, createdAt)
}
