package sqlite

import (
	"context"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/maximekuhn/partage/internal/core/entity"
	"github.com/maximekuhn/partage/internal/core/valueobject"
)

func TestSaveGroup(t *testing.T) {
	db := CreateTmpDB()
	defer db.Close()

	s := NewSQLiteGroupStore(db)

	tests := []struct {
		title     string
		groupName string
		ownerID   uuid.UUID
		members   []uuid.UUID
	}{
		{
			title:     "Simple group with no members",
			groupName: "My awesome group where I am alone",
			ownerID:   uuid.New(),
			members:   []uuid.UUID{},
		},
		{
			title:     "A group with one member",
			groupName: "My awesome group with my only friend",
			ownerID:   uuid.New(),
			members:   []uuid.UUID{uuid.New()},
		},
		{
			title:     "A group with multiple members",
			groupName: "My awesome group with all my friends",
			ownerID:   uuid.New(),
			members: []uuid.UUID{
				uuid.New(),
				uuid.New(),
				uuid.New(),
				uuid.New(),
				uuid.New(),
				uuid.New(),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			group := createGroup(test.groupName, test.ownerID, test.members)
			err := s.Save(context.TODO(), group)
			if err != nil {
				t.Errorf("Save(): expected ok got error %v", err)
			}
		})
	}
}

func createGroup(name string, owner uuid.UUID, members []uuid.UUID) *entity.Group {
	id, err := valueobject.NewGroupID(uuid.New())
	if err != nil {
		panic(err)
	}
	ownerID, err := valueobject.NewUserID(owner)
	if err != nil {
		panic(err)
	}
	membersID := make([]valueobject.UserID, 0)
	for _, member := range members {
		memberID, err := valueobject.NewUserID(member)
		if err != nil {
			panic(err)
		}
		membersID = append(membersID, memberID)
	}
	groupname, err := valueobject.NewGroupname(name)
	if err != nil {
		panic(err)
	}
	createdAt := time.Now()
	return entity.NewGroup(id, groupname, membersID, ownerID, createdAt)
}

func TestFindGroupByName(t *testing.T) {
	db := CreateTmpDB()
	defer db.Close()

	s := NewSQLiteGroupStore(db)
	group := createGroup("My first group", uuid.New(), []uuid.UUID{
		uuid.New(),
		uuid.New(),
		uuid.New(),
		uuid.New(),
	})
	err := s.Save(context.TODO(), group)
	if err != nil {
		panic(err)
	}

	tests := []struct {
		title         string
		groupname     string
		shouldBeFound bool
	}{
		{
			title:         "Wrong name",
			groupname:     "My second group",
			shouldBeFound: false,
		},
		{
			title:         "Correct name",
			groupname:     "My first group",
			shouldBeFound: true,
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			groupname, err := valueobject.NewGroupname(test.groupname)
			if err != nil {
				panic(groupname)
			}
			g, found, err := s.FindByName(context.TODO(), groupname)
			if err != nil {
				t.Fatalf("FindByName(): should get ok got error %v", err)
			}

			if test.shouldBeFound && !found {
				t.Fatal("FindByName(): expected to find group got nothing")
			} else if !test.shouldBeFound && !found {
				// ok
				return
			}

			// sort members for comparison
			sort.Slice(g.Members, func(i, j int) bool {
				return g.Members[i].String() < g.Members[j].String()
			})

			sort.Slice(group.Members, func(i, j int) bool {
				return group.Members[i].String() < group.Members[j].String()
			})

			if !reflect.DeepEqual(g, group) {
				t.Fatalf("FindByName(): expected group %v got %v", group, g)
			}
		})
	}

}
