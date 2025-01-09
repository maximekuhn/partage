package sqlite

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/maximekuhn/partage/internal/core/entity"
	"github.com/maximekuhn/partage/internal/core/store"
	"github.com/maximekuhn/partage/internal/core/valueobject"
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

func TestSaveDuplicate(t *testing.T) {
	db := CreateTmpDB()
	defer db.Close()

	s := NewSQLiteUserStore(db)

	u := createUser("toto", "toto@gmail.com")
	err := s.Save(context.TODO(), u)
	if err != nil {
		t.Errorf("Save(): expected no error got %v", err)
	}

	err = s.Save(context.TODO(), u)
	if err != store.ErrUserStoreDuplicate {
		t.Errorf("Save(): expected ErrUserStoreDuplicate got %v", err)
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
func createUser(nickname, email string) *entity.User {
	id, err := valueobject.NewUserID(uuid.New())
	if err != nil {
		panic(err)
	}
	nn, err := valueobject.NewNickname(nickname)
	if err != nil {
		panic(err)
	}
	em, err := valueobject.NewEmail(email)
	if err != nil {
		panic(err)
	}
	createdAt := time.Now()
	return entity.NewUser(id, em, nn, createdAt)
}
